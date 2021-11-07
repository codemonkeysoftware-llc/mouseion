package sqlstore

import (
	"context"
	"fmt"
	"log"
	"time"

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
	err := tx.Get(&id, "INSERT INTO log_entry (timestamp, text) VALUES($1,$2) RETURNING id", logEntry.Timestamp.UnixNano(), logEntry.Text)
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

func (sqlStore *SQLStore) GetRange(ctx context.Context, start, end time.Time) ([]entry.LogEntry, error) {
	if start.IsZero() {
		start = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if end.IsZero() {
		end = time.Date(2030, 1, 1, 1, 1, 1, 1, time.UTC)
	}

	rows, err := sqlStore.db.QueryxContext(ctx, `
SELECT id, timestamp, text
FROM log_entry
WHERE timestamp >= $1
AND timestamp <= $2`,
		start.UnixNano(),
		end.UnixNano())
	if err != nil {
		return nil, fmt.Errorf("entries: %w", err)
	}
	entries := make([]entry.LogEntry, 0)
	for rows.Next() {
		entry := entry.LogEntry{}
		var timestamp int
		err := rows.Scan(&entry.ID, &timestamp, &entry.Text)
		if err != nil {
			return nil, err
		}
		entry.Timestamp = parseUnixUTC(timestamp)
		log.Println(entry)
		tags, err := sqlStore.GetTagsForEntry(ctx, entry.ID)
		if err != nil {
			return nil, fmt.Errorf("entry tags: %w", err)
		}
		entry.Tags = tags
		entries = append(entries, entry)
	}

	return entries, nil
}

func (sqlStore *SQLStore) GetTagsForEntry(ctx context.Context, entryID int) ([]string, error) {
	tags := make([]string, 0)
	err := sqlStore.db.SelectContext(ctx, &tags, "SELECT tag FROM entry_tag WHERE entry_id = $1", entryID)
	if err != nil {
		return nil, fmt.Errorf("entry tags: %w", err)
	}
	return tags, nil
}

func parseUnixUTC(unixnano int) time.Time {
	localTime := time.Unix(
		int64(unixnano)/int64(time.Second),
		int64(unixnano)%int64(time.Second),
	)

	_, offset := localTime.Zone()
	localTime.Add(time.Hour * time.Duration(offset))
	return localTime.UTC()
}
