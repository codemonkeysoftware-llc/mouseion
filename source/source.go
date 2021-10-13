package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/codemonkeysoftware/mouseion/pkg/entry"
)

// Logger sends log messages to MouseionHost using HTTPClient. MouseionHost
// is the only required field.
type Logger struct {
	HTTPClient   *http.Client
	LogErrors    bool
	MouseionHost string
}

func (logger *Logger) Print(v ...interface{}) {
	logger.send(fmt.Sprint(v...))
}

func (logger *Logger) Printf(format string, v ...interface{}) {
	logger.send(fmt.Sprintf(format, v...))
}

func (logger *Logger) Println(v ...interface{}) {
	logger.Print(v...)
}

func (logger *Logger) send(text string) {
	err := send(logger.HTTPClient, logger.MouseionHost, text)
	if err != nil && logger.LogErrors {
		log.Println(err)
	}
}

func send(client *http.Client, host string, text string) error {
	if client != nil {
		client = http.DefaultClient
	}
	entry := entry.LogEntry{Timestamp: time.Now(), Text: text, Tags: nil}
	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	resp, err := client.Post(jsonURL(host), "application/json", bytes.NewBuffer(entryJSON))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%d %s", resp.StatusCode, body)
	}
	return nil
}
func jsonURL(host string) string {
	return host + "/ingest.json"
}
