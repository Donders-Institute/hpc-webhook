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

func addRow(db *sql.DB, hash string, groupname string, username string) error {
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

	sqlStatement := fmt.Sprintf("INSERT INTO qaas (hash, groupname, username) VALUES ($1, $2, $3)")

	if _, err = tx.Exec(sqlStatement, hash, groupname, username); err != nil {
		return err
	}

	return err
}

func deleteRow(db *sql.DB, hash string, groupname string, username string) error {
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

	sqlStatement := fmt.Sprintf("DELETE FROM qaas (hash, groupname, username) VALUES ($1, $2, $3)")

	if _, err = tx.Exec(sqlStatement, hash, groupname, username); err != nil {
		return err
	}

	return err
}

type Item struct {
	ID        int    `json:"-"` // Do not output this one
	Hash      string `json:"hash"`
	Groupname string `json:"groupname"`
	Username  string `json:"username"`
}

// Find the rows with a specific hash (should be 1)
func getRow(db *sql.DB, hash string) ([]Item, error) {
	rows, err := db.Query("SELECT id, hash, groupname, username FROM qaas WHERE hash = $1", hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Item
	for rows.Next() {
		p := Item{}
		if err := rows.Scan(&p.ID, &p.Hash, &p.Groupname, &p.Username); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if rows.Err() != nil {
		return nil, err
	}
	if len(list) > 1 {
		return nil, fmt.Errorf("invalid getRow result: list should have length 1 but has length %d", len(list))
	}

	return list, nil
}

// Find the rows for a specific groupname, username
func getListRows(db *sql.DB, groupname string, username string) ([]Item, error) {
	rows, err := db.Query("SELECT id, hash, groupname, username FROM qaas groupname = $1, username = $2", groupname, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Item
	for rows.Next() {
		p := Item{}
		if err := rows.Scan(&p.ID, &p.Hash, &p.Groupname, &p.Username); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return list, nil
}
