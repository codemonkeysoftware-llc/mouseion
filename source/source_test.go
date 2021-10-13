package source_test

import (
	"net/http"
	"testing"

	"github.com/codemonkeysoftware/mouseion/source"
)

func TestSendData(t *testing.T) {
	logger := source.Logger{
		MouseionHost: "http://localhost:8069",
		HTTPClient:   &http.Client{},
		LogErrors:    true,
	}
	logger.Print("hello")
}
