package main

import (
	"log"

	sqlstore "github.com/codemonkeysoftware/mouseion/pkg/entry/sql_store"
	"github.com/codemonkeysoftware/mouseion/pkg/sqlite"
	"github.com/codemonkeysoftware/mouseion/pkg/webserver"
)

func main() {
	db := sqlite.Open("./mouseion.db")
	defer db.Close()
	err := sqlite.DoMigrations(db.DB)
	if err != nil {
		log.Println(err)
	}
	entryService := sqlstore.New(db)
	webserver.New(entryService).Start()
}
