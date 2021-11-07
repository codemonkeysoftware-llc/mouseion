package sqlstore_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/codemonkeysoftware/mouseion/pkg/entry"
	sqlstore "github.com/codemonkeysoftware/mouseion/pkg/entry/sql_store"
	"github.com/codemonkeysoftware/mouseion/pkg/sqlite"
	"github.com/codemonkeysoftware/mouseion/pkg/testhelpers"
)

func TestEntryLifecycle(t *testing.T) {
	db := sqlite.Open(fmt.Sprintf("%s_test.db", testhelpers.RandString(12)))
	err := sqlite.DoMigrations(db.DB)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
	}()
	store := sqlstore.New(db)
	entryDate := time.Date(2011, 11, 11, 11, 11, 11, 11, time.UTC)
	logEntry := entry.LogEntry{
		Timestamp: entryDate,
		Text:      "peppero day",
		Tags:      []string{"snacks", "dates"},
	}
	err = store.Save(context.Background(), &logEntry)
	if err != nil {
		t.Fatal(err)
	}

	badEntry1 := entry.LogEntry{
		Timestamp: entryDate.Add(-time.Minute),
		Text:      "peppero day1",
		Tags:      []string{"snacks", "dates"},
	}
	err = store.Save(context.Background(), &badEntry1)
	if err != nil {
		t.Fatal(err)
	}

	badEntry2 := entry.LogEntry{
		Timestamp: entryDate.Add(time.Minute),
		Text:      "peppero day2",
		Tags:      []string{"snacks", "dates"},
	}
	err = store.Save(context.Background(), &badEntry2)
	if err != nil {
		t.Fatal(err)
	}

	results, err := store.GetRange(
		context.Background(),
		entryDate.Add(-time.Second*10),
		entryDate.Add(time.Second*10),
	)

	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Text != logEntry.Text {
		log.Fatalf("result text was %s", results[0].Text)
	}
	if results[0].Timestamp != logEntry.Timestamp {
		log.Fatalf("expected timestamp %s, got %s", logEntry.Timestamp, results[0].Timestamp)
	}

	results, err = store.GetRange(
		context.Background(),
		time.Time{},
		time.Time{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func BenchmarkGetRange(b *testing.B) {
	originalWriter := log.Default().Writer()
	log.Default().SetOutput(io.Discard)
	defer func() {
		log.Default().SetOutput(originalWriter)
	}()
	b.Run("100 x 10", func(b *testing.B) {
		const (
			entryCount = 100
			tagCount   = 10
		)
		dbName := fmt.Sprintf("%s_test.db", testhelpers.RandString(12))
		db := sqlite.Open(dbName)
		err := sqlite.DoMigrations(db.DB)
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			db.Close()
			os.Remove(dbName)
		}()
		store := sqlstore.New(db)

		for i := 0; i < entryCount; i++ {
			logEntry := entry.LogEntry{
				Timestamp: time.Now(),
				Text:      "logentry",
				Tags:      make([]string, 0, tagCount),
			}
			for n := 0; n < tagCount; n++ {
				logEntry.Tags = append(logEntry.Tags, fmt.Sprintf("tag %d", n))
			}
			store.Save(context.Background(), &logEntry)
		}
		b.ResetTimer()
		_, err = store.GetRange(context.Background(), time.Time{}, time.Time{})
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
	})
}
