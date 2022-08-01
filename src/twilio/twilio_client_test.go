package twilio

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type mockTwilioRestApi struct {
	mockCreateMessage func(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error)
}

func (m *mockTwilioRestApi) CreateMessage(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error) {
	return m.mockCreateMessage(params)
}

func TestNewClientSuccess(t *testing.T) {
	assert := assert.New(t)

	os.Setenv(accountSIDEnvKey, "TEST_SID")
	defer os.Unsetenv(accountSIDEnvKey)

	os.Setenv(authTokenEnvKey, "TEST_AUTH_TOKEN")
	defer os.Unsetenv(authTokenEnvKey)

	os.Setenv(fromPhoneEnvKey, "TEST_FROM_PHONE")
	defer os.Unsetenv(fromPhoneEnvKey)

	c, err := NewClient("TEST_NUMBER")

	assert.NotNil(c)
	assert.Nil(err)
}
func TestSendMessageSuccess(t *testing.T) {
	assert := assert.New(t)

	mockCreateMessage := func(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error) {
		return &openapi.ApiV2010Message{}, nil
	}

	c := &Client{
		api: &mockTwilioRestApi{
			mockCreateMessage: mockCreateMessage,
		},
	}

	err := c.SendMessage("123456", "TEST_MESSAGE")

	assert.Nil(err)
}

func TestSendMessageFail(t *testing.T) {
	assert := assert.New(t)

	mockCreateMessage := func(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error) {
		return &openapi.ApiV2010Message{}, errors.New("something went wrong!")
	}

	c := &Client{
		api: &mockTwilioRestApi{
			mockCreateMessage: mockCreateMessage,
		},
	}

	err := c.SendMessage("123456", "TEST_MESSAGE")

	assert.EqualError(err, "something went wrong!")
}

func TestSendMessageErrorCode(t *testing.T) {
	assert := assert.New(t)

	mockCreateMessage := func(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error) {
		errorMessage := "something else went wrong"
		errorCode := 12345
		return &openapi.ApiV2010Message{
			ErrorCode:    &errorCode,
			ErrorMessage: &errorMessage,
		}, nil
	}

	c := &Client{
		api: &mockTwilioRestApi{
			mockCreateMessage: mockCreateMessage,
		},
	}

	err := c.SendMessage("123456", "TEST_MESSAGE")

	assert.EqualError(err, "error in Twilio CreateMessage.\n\nError code: 12345\n\nError Message: something else went wrong")
}

func TestSendHelp(t *testing.T) {
	assert := assert.New(t)

	mockCreateMessage := func(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error) {
		return &openapi.ApiV2010Message{}, nil
	}

	c := &Client{
		api: &mockTwilioRestApi{
			mockCreateMessage: mockCreateMessage,
		},
	}

	err := c.SendHelp("123456")

	assert.Nil(err)
}
