package config

import (
	_ "embed"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/WilliamDeBruin/nps_alerts/src/common"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Port string `envconfig:"PORT" required:"false" default:"8080"`

	TwilioFromNumber string `envconfig:"TWILIO_FROM_NUMBER" required:"true"`
	TwilioAccountSID string `envconfig:"TWILIO_ACCOUNT_SID" required:"true"`
	TwilioAuthToken  string `envconfig:"TWILIO_AUTH_TOKEN" required:"true"`

	// ServiceHost is used in integration tests.
	ServiceHost string `envconfig:"SERVICE_HOST" required:"false" default:"127.0.0.1"`
}

// ToJSON returns a json formatted version of the configuration with
// sensitive field redacted. To redact a field, add a json tag of `redact:"true"`
func (c *Configuration) ToJSON() string {
	if c == nil {
		return "null" // json for nil
	}
	ret := make(map[string]string)
	// We know c is of type struct, which allows us to safely use the reflect package to access
	// its members (along with the it's tags i.e. the metadata between the back-ticks)
	cfg := reflect.ValueOf(c).Elem()
	for i := 0; i < cfg.NumField(); i++ {
		name := cfg.Type().Field(i).Name
		val := fmt.Sprintf("%v", cfg.Field(i).Interface())
		tag := cfg.Type().Field(i).Tag.Get("redact")
		if redact, err := strconv.ParseBool(tag); err == nil && redact {
			val = "****"
		}
		ret[name] = val
	}
	return common.ToJSON(ret)
}

// LoadConfig loads environment variables with the prefix
func LoadConfig() (Configuration, error) {
	cfg := Configuration{}
	err := envconfig.Process(strings.ToUpper(""), &cfg)
	return cfg, err
}
