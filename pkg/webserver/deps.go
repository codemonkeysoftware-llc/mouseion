package webserver

import (
	"context"

	"github.com/codemonkeysoftware/mouseion/pkg/entry"
)

type EntryStorer interface {
	Save(ctx context.Context, logEntry *entry.LogEntry) error
}
