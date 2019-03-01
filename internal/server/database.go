package server

import (
	"database/sql"
	"errors"
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
	if !isValidWebhookID(hash) {
		return errors.New("invalid webhook id")
	}

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

	sqlStatement := fmt.Sprintf("INSERT INTO qaas (hash, username) VALUES ($1, $2)")

	if _, err = tx.Exec(sqlStatement, hash, username); err != nil {
		return err
	}

	return err
}

type item struct {
	ID       int
	Hash     string
	Username string
}

func getRow(db *sql.DB, hash string) ([]item, error) {
	rows, err := db.Query("SELECT id, hash, username FROM qaas WHERE hash = $1", hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []item
	for rows.Next() {
		p := item{}
		if err := rows.Scan(&p.ID, &p.Hash, &p.Username); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return list, nil
}
