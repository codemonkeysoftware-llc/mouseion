package webserver

import (
	"context"
	"time"

	"github.com/codemonkeysoftware/mouseion/pkg/entry"
)

type EntryStorer interface {
	Save(ctx context.Context, logEntry *entry.LogEntry) error
	GetEntries(ctx context.Context, start, end time.Time, tags []string) ([]entry.LogEntry, error)
}
