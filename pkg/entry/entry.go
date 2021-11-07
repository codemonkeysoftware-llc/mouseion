package entry

import "time"

type LogEntry struct {
	ID        int       `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Text      string    `json:"text" db:"text"`
	Tags      []string  `json:"tags"`
}
