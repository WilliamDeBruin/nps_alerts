package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/WilliamDeBruin/nps_alerts/src/nps"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type mockNpsClient struct {
	getAlertResponse *nps.AlertDetails
	getAlertErr      error
}

func (m *mockNpsClient) GetAlert(stateCode string) (*nps.AlertDetails, error) {
	return m.getAlertResponse, m.getAlertErr
}

func (m *mockNpsClient) SetTransport(rt http.RoundTripper) {}

type mockTwilioClient struct {
	sendMessageErr error
}

func (m *mockTwilioClient) SendMessage(to, message string) error {
	return m.sendMessageErr
}

func TestHealthHandler(t *testing.T) {
	assert := assert.New(t)

	r := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	s := Server{}

	s.HealthHandler(w, r)

	res := w.Result()

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	assert.Equal(string(data), "all is good!")
}

func TestIncomingSmsHelp(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "help")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    &mockNpsClient{},
		twilioClient: &mockTwilioClient{},
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(logs.All()[0].Message, "sent help message")
	assert.Equal(w.Result().StatusCode, http.StatusOK)
}

func TestIncomingSmsHelpErr(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "help")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient: &mockNpsClient{},
		twilioClient: &mockTwilioClient{
			sendMessageErr: errors.New("TEST_SEND_MESSAGE_ERR"),
		},
		logger: logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(logs.All()[0].Message, "TEST_SEND_MESSAGE_ERR")
	assert.Equal(w.Result().StatusCode, http.StatusInternalServerError)
}

func TestIncomingSmsAlert(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "alerts CA")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockNpsClient := &mockNpsClient{}
	mockNpsClient.getAlertResponse = &nps.AlertDetails{
		FullStateName:   "TEST_FULL_STATE_NAME",
		FullParkName:    "TEST_FULL_PARK_NAME",
		RecentAlertDate: "TEST_DATE",
		AlertHeader:     "TEST_HEADER",
		AlertMessage:    "TEST_MESSAGE",
		URL:             "TEST_URL",
	}

	mockTwilioClient := &mockTwilioClient{}
	mockTwilioClient.sendMessageErr = nil

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    mockNpsClient,
		twilioClient: mockTwilioClient,
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(logs.All()[0].Message, "alert response")
	assert.Equal(w.Result().StatusCode, http.StatusOK)
}

func TestIncomingSmsAlertSendMessageFail(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "alerts CA")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockNpsClient := &mockNpsClient{}
	mockNpsClient.getAlertResponse = &nps.AlertDetails{
		FullStateName:   "TEST_FULL_STATE_NAME",
		FullParkName:    "TEST_FULL_PARK_NAME",
		RecentAlertDate: "TEST_DATE",
		AlertHeader:     "TEST_HEADER",
		AlertMessage:    "TEST_MESSAGE",
		URL:             "TEST_URL",
	}

	mockTwilioClient := &mockTwilioClient{}
	mockTwilioClient.sendMessageErr = errors.New("TEST_FAILED_TO_SEND_MESSAGE")

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    mockNpsClient,
		twilioClient: mockTwilioClient,
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(logs.All()[1].Message, "TEST_FAILED_TO_SEND_MESSAGE")
	assert.Equal(w.Result().StatusCode, http.StatusInternalServerError)
}

func TestIncomingSmsAlertGetAlertFail(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "alerts CA")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockNpsClient := &mockNpsClient{}
	mockNpsClient.getAlertResponse = &nps.AlertDetails{
		FullStateName:   "TEST_FULL_STATE_NAME",
		FullParkName:    "TEST_FULL_PARK_NAME",
		RecentAlertDate: "TEST_DATE",
		AlertHeader:     "TEST_HEADER",
		AlertMessage:    "TEST_MESSAGE",
		URL:             "TEST_URL",
	}
	mockNpsClient.getAlertErr = errors.New("TEST_FAILED_TO_GET_ALERT")

	mockTwilioClient := &mockTwilioClient{}
	mockTwilioClient.sendMessageErr = nil

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    mockNpsClient,
		twilioClient: mockTwilioClient,
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(logs.All()[0].Message, "TEST_FAILED_TO_GET_ALERT")
	assert.Equal(w.Result().StatusCode, http.StatusInternalServerError)
}

func TestIncomingSmsAlertBadMessage(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "alerts new mexico")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockNpsClient := &mockNpsClient{}
	mockNpsClient.getAlertResponse = &nps.AlertDetails{
		FullStateName:   "TEST_FULL_STATE_NAME",
		FullParkName:    "TEST_FULL_PARK_NAME",
		RecentAlertDate: "TEST_DATE",
		AlertHeader:     "TEST_HEADER",
		AlertMessage:    "TEST_MESSAGE",
		URL:             "TEST_URL",
	}

	mockTwilioClient := &mockTwilioClient{}
	mockTwilioClient.sendMessageErr = nil

	core, _ := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    mockNpsClient,
		twilioClient: mockTwilioClient,
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(w.Result().StatusCode, http.StatusBadRequest)
}

func TestIncomingSmsAlertBadMessageSendFail(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "alerts new mexico")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockNpsClient := &mockNpsClient{}
	mockNpsClient.getAlertResponse = &nps.AlertDetails{
		FullStateName:   "TEST_FULL_STATE_NAME",
		FullParkName:    "TEST_FULL_PARK_NAME",
		RecentAlertDate: "TEST_DATE",
		AlertHeader:     "TEST_HEADER",
		AlertMessage:    "TEST_MESSAGE",
		URL:             "TEST_URL",
	}

	mockTwilioClient := &mockTwilioClient{}
	mockTwilioClient.sendMessageErr = errors.New("TEST_SEND_MESSAGE_FAIL_ERR")

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    mockNpsClient,
		twilioClient: mockTwilioClient,
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(logs.All()[0].Message, "TEST_SEND_MESSAGE_FAIL_ERR")
	assert.Equal(w.Result().StatusCode, http.StatusInternalServerError)
}

func TestIncomingSmsUnhandledBody(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "unhandled")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockNpsClient := &mockNpsClient{}
	mockNpsClient.getAlertResponse = &nps.AlertDetails{
		FullStateName:   "TEST_FULL_STATE_NAME",
		FullParkName:    "TEST_FULL_PARK_NAME",
		RecentAlertDate: "TEST_DATE",
		AlertHeader:     "TEST_HEADER",
		AlertMessage:    "TEST_MESSAGE",
		URL:             "TEST_URL",
	}

	mockTwilioClient := &mockTwilioClient{}
	mockTwilioClient.sendMessageErr = nil

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    mockNpsClient,
		twilioClient: mockTwilioClient,
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(logs.All()[0].Message, "unhandled text body")
	assert.Equal(w.Result().StatusCode, http.StatusBadRequest)
}

func TestIncomingSmsInvalidContentType(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencodeddd")

	s := Server{
		npsClient:    &mockNpsClient{},
		twilioClient: &mockTwilioClient{},
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(w.Result().StatusCode, http.StatusUnsupportedMediaType)
}

func TestIncomingSmsMissingFrom(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    &mockNpsClient{},
		twilioClient: &mockTwilioClient{},
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(logs.All()[0].Message, "missing field in request body: from")
}

func TestIncomingSmsMissingBody(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("from", "+123456789")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	s := Server{
		npsClient:    &mockNpsClient{},
		twilioClient: &mockTwilioClient{},
		logger:       logger,
	}

	s.IncomingSmsHandler(w, r)

	assert.Equal(w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(logs.All()[0].Message, "missing field in request body: body")
}
