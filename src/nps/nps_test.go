package nps

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	responseBody any
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}

	responseString, _ := json.Marshal(m.responseBody)

	response.Header.Set("Content-Type", "application/json")
	response.Body = ioutil.NopCloser(strings.NewReader(string(responseString)))
	return response, nil
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	c, err := NewClient("TEST_KEY")

	assert.NotNil(c)
	assert.Nil(err)

	c, err = NewClient("")

	assert.Nil(c)
	assert.EqualError(err, "apiKey cannot be empty")
}

func TestGetAlertInvalidStateCode(t *testing.T) {
	assert := assert.New(t)

	mockTransport := &mockTransport{}

	c, _ := NewClient("TEST_KEY")
	c.SetTransport(mockTransport)

	details, err := c.GetAlert("MV")

	assert.Nil(details)
	assert.EqualError(err, "state code MV is not a valid state code")
}

func TestGetAlertSuccess(t *testing.T) {
	assert := assert.New(t)

	mockTransport := &mockTransport{}

	c, _ := NewClient("TEST_KEY")
	c.SetTransport(mockTransport)

	mockTransport.responseBody = alertResponse{
		Total: "1",
		Limit: "1",
		Start: "0",
		Data: []npsAlert{
			{
				ID:              "TEST_ID",
				URL:             "TEST_URL",
				Title:           "TEST_TITLE",
				ParkCode:        "yell", // valid park code == yellowstone natl
				Description:     "TEST_DESCRIPTION",
				Category:        "TEST_CATEGORY",
				LastIndexedDate: "2022-08-02 12:34:45.6",
			},
		},
	}

	details, err := c.GetAlert("MT")

	assert.Equal(details, &AlertDetails{
		FullStateName:   "Montana",
		FullParkName:    "Yellowstone",
		RecentAlertDate: "2022-08-02 12:34:45.6",
		AlertHeader:     "TEST_TITLE",
		AlertMessage:    "TEST_DESCRIPTION",
		URL:             "https://www.nps.gov/planyourvisit/alerts.htm?s=MT&p=1&v=0",
	})
	assert.Nil(err)
}

func TestGetAlertInvalidParkCode(t *testing.T) {
	assert := assert.New(t)

	mockTransport := &mockTransport{}

	c, _ := NewClient("TEST_KEY")
	c.SetTransport(mockTransport)

	mockTransport.responseBody = alertResponse{
		Total: "1",
		Limit: "1",
		Start: "0",
		Data: []npsAlert{
			{
				ID:              "TEST_ID",
				URL:             "TEST_URL",
				Title:           "TEST_TITLE",
				ParkCode:        "INVALID_CODE", // valid park code == yellowstone natl
				Description:     "TEST_DESCRIPTION",
				Category:        "TEST_CATEGORY",
				LastIndexedDate: "2022-08-02 12:34:45.6",
			},
		},
	}

	details, err := c.GetAlert("MT")

	assert.Nil(details)
	assert.EqualError(err, "cannot find details for park code INVALID_CODE")
}
