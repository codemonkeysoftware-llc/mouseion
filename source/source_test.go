package source_test

import (
	"net/http"
	"testing"

	"github.com/codemonkeysoftware/mouseion/source"
)

func TestSendData(t *testing.T) {
	logger := source.NewLogger("http://localhost:8069", &http.Client{}, true)
	logger.Print("hello")
}
