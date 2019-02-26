package server

import (
	"database/sql"
	"log"

	// Postgres driver
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
