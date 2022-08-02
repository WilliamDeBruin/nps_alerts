package server

import (
	"testing"

	"github.com/WilliamDeBruin/nps_alerts/src/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)

	cfg := &config.Configuration{
		TwilioFromNumber: "+123456789",
		NPSApiKey:        "TEST_KEY",
	}
	logger := zaptest.NewLogger(t)

	s, err := NewServer(cfg, logger)

	assert.NotNil(s)
	assert.Nil(err)
}

func TestNewServerMissingParams(t *testing.T) {
	assert := assert.New(t)

	cfg := &config.Configuration{}
	logger := zaptest.NewLogger(t)

	s, err := NewServer(cfg, logger)

	assert.Nil(s)
	assert.EqualError(err, "error initializing twilio client: fromNumber cannot be empty")

	cfg = &config.Configuration{
		TwilioFromNumber: "+123456789",
	}

	s, err = NewServer(cfg, logger)

	assert.Nil(s)
	assert.EqualError(err, "error initializing nps client: apiKey cannot be empty")

}

// func TestListen(t *testing.T) {
// 	assert := assert.New(t)

// 	cfg := &config.Configuration{}
// 	logger := zaptest.NewLogger(t)

// 	s, err := NewServer(cfg, logger)

// }
