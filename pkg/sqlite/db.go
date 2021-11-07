package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*
var migrations embed.FS

func Open(path string) *sqlx.DB {
	log.Printf("Opening DB at %s", path)
	db := sqlx.MustOpen("sqlite3", path)
	var version string
	db.Get(&version, "SELECT sqlite_version()")
	log.Printf("Using sqlite version %s", version)
	return db
}

func DoMigrations(db *sql.DB) error {
	log.Println("Running Migrations")
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}
	sourceDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	log.Println("Migrations Complete")
	return nil
}
