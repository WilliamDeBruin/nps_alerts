package twilio

import (
	"encoding/json"
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioRestClientApi interface {
	CreateMessage(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error)
}

type fetcher struct {
	API        TwilioRestClientApi
	fromNumber string
}

type Client interface {
	SendMessage(to, message string) error
}

func NewClient(fromNumber string) (Client, error) {

	if fromNumber == "" {
		return nil, fmt.Errorf("fromNumber cannot be empty")
	}

	client := twilio.NewRestClient()
	return &fetcher{
		API:        client.Api,
		fromNumber: fromNumber,
	}, nil
}

func (c *fetcher) SendMessage(to, message string) error {

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(c.fromNumber)
	params.SetBody(message)

	resp, err := c.API.CreateMessage(params)
	if err != nil {
		return err
	}

	if resp.ErrorCode != nil {
		return fmt.Errorf("error in Twilio CreateMessage.\n\nError code: %d\n\nError Message: %s", *resp.ErrorCode, *resp.ErrorMessage)
	}

	response, err := json.Marshal(*resp)

	fmt.Println("Response: " + string(response))

	if err != nil {
		return err
	}

	return nil
}
