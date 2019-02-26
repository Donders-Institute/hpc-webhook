package server

import (
	"database/sql"
	"fmt"
	"log"

	// Postgres driver
	_ "github.com/lib/pq"
)

// InitDB initializes the database
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db, err
}

func addRow(db *sql.DB, hash string, username string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	sqlStatement := fmt.Sprintf("INSERT INTO qaas (hash, username) VALUES (?, ?)")

	if _, err = tx.Exec(sqlStatement, hash, username); err != nil {
		return err
	}

	return err
}
