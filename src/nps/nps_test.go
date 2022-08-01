package nps

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAlert(t *testing.T) {
	assert := assert.New(t)
	os.Setenv(apiKeyEnvKey, "")
	defer os.Unsetenv(apiKeyEnvKey)

	c, _ := NewClient()

	assert.NotNil(c)
}
