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

type Logger struct {
	client    *http.Client
	logErrors bool
	host      string
}

func NewLogger(host string, client *http.Client, logErrors bool) *Logger {
	return &Logger{
		host:      host,
		client:    client,
		logErrors: true,
	}
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
	entry := entry.LogEntry{Timestamp: time.Now(), Text: text, Tags: nil}
	entryJSON, err := json.Marshal(entry)
	if logger.logErrors {
		log.Println(err)
	}
	resp, err := logger.client.Post(logger.url(), "application/json", bytes.NewBuffer(entryJSON))
	if err != nil && logger.logErrors {
		log.Println(err)
	}
	if logger.logErrors && resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.StatusCode, body)
	}
}
func (logger *Logger) url() string {
	return logger.host + "/ingest.json"
}
