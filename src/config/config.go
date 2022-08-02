package config

import (
	_ "embed"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Port string `envconfig:"PORT" required:"false" default:"8080"`

	TwilioFromNumber string `envconfig:"TWILIO_FROM_NUMBER" required:"true"`
	TwilioAccountSID string `envconfig:"TWILIO_ACCOUNT_SID" required:"true"`
	TwilioAuthToken  string `envconfig:"TWILIO_AUTH_TOKEN" required:"true"`

	// ServiceHost is used in integration tests.
	ServiceHost string `envconfig:"SERVICE_HOST" required:"false" default:"127.0.0.1"`

	NPSApiKey string `envconfig:"NPS_API_KEY" required:"true"`
}

// LoadConfig loads environment variables with the prefix
func LoadConfig() (Configuration, error) {
	cfg := Configuration{}
	err := envconfig.Process(strings.ToUpper(""), &cfg)
	return cfg, err
}
