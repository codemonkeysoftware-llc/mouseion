package sqlstore

import (
	"context"
	"fmt"

	"github.com/codemonkeysoftware/mouseion/pkg/entry"
	"github.com/jmoiron/sqlx"
)

type SQLStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

type entryTag struct {
	ID  string `db:"id"`
	Tag string `db:"tag"`
}

func (sqlStore *SQLStore) Save(ctx context.Context, logEntry *entry.LogEntry) error {
	var id string
	tx := sqlStore.db.MustBeginTx(ctx, nil)
	defer tx.Rollback()
	err := tx.Get(&id, "INSERT INTO log_entry (timestamp, text) VALUES($1,$2) RETURNING id", logEntry.Timestamp, logEntry.Text)
	if err != nil {
		return fmt.Errorf("insert entries: %w", err)
	}
	if len(logEntry.Tags) > 0 {
		entryTags := make([]entryTag, len(logEntry.Tags))
		for i, tag := range logEntry.Tags {
			entryTags[i] = entryTag{
				ID:  id,
				Tag: tag,
			}
		}
		_, err = tx.NamedExec("INSERT INTO entry_tag (entry_id,tag) VALUES (:id,:tag)", entryTags)
		if err != nil {
			return fmt.Errorf("tag entries: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit %w", err)
	}
	return nil
}
