package server

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertHandler(t *testing.T) {
	assert := assert.New(t)

	data := url.Values{}
	data.Set("body", "alerts CA")
	data.Set("from", "+12407439754")

	r := httptest.NewRequest("POST", "http://example.com/", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s := Server{}

	s.IncomingSmsHandler(w, r)

	assert.True(true)
}
