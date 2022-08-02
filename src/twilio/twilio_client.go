package twilio

import (
	"encoding/json"
	"fmt"

	"github.com/WilliamDeBruin/nps_alerts/src/nps"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	PhoneNumberEnvKey = "TWILIO_PHONE_NUMBER"
	accountSIDEnvKey  = "TWILIO_ACCOUNT_SID"
	authTokenEnvKey   = "TWILIO_AUTH_TOKEN"
	fromPhoneEnvKey   = "TWILIO_FROM_PHONE_NUMBER"

	helpMessage  = "Welcome to NPS alerts! Here is a list of commands:\n\nHelp: receive this help text\n\nAlerts {state}: Text \"alerts\" followed by the 2-letter state code of the state you would like to see alerts for"
	alertMessage = "Here is the most recent NPS %s alert from %s, published %s:\n\n%s\n\n%s\n\nFor a full list of NPS %s alerts, visit %s"
)

type TwilioRestClientApi interface {
	CreateMessage(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error)
}

type Client struct {
	api        TwilioRestClientApi
	fromNumber string
}

func NewClient(fromNumber string) (*Client, error) {

	if fromNumber == "" {
		return nil, fmt.Errorf("fromNumber cannot be empty")
	}

	client := twilio.NewRestClient()
	return &Client{
		api:        client.Api,
		fromNumber: fromNumber,
	}, nil
}

func (c *Client) SendMessage(to, message string) error {

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(c.fromNumber)
	params.SetBody(message)

	resp, err := c.api.CreateMessage(params)
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

func (c *Client) SendHelp(to string) error {
	return c.SendMessage(to, helpMessage)
}

func (c *Client) SendAlert(to string, params *nps.AlertDetails) error {
	message := fmt.Sprintf(alertMessage,
		params.FullStateName,
		params.FullParkName,
		params.RecentAlertDate,
		params.AlertHeader,
		params.AlertMessage,
		params.FullStateName,
		params.URL)

	return c.SendMessage(to, message)
}

func (c *Client) SendAlertErr(to string) error {
	return c.SendMessage(to, `I'm sorry, I couldn't understand your message. Please text "alerts {state}" for recent alerts`)
}
