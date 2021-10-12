package entry

import "time"

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
}
