package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigMissingVars(t *testing.T) {
	assert := assert.New(t)

	_, err := LoadConfig()

	assert.EqualError(err, "required key TWILIO_FROM_NUMBER missing value")
}
