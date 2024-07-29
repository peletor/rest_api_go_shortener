package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.new"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, errors.New(op + ": Failed to open database: " + err.Error())
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    alias TEXT NOT NULL UNIQUE,
		    url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
		`)

	if err != nil {
		return nil, errors.New(op + ": Failed to prepare table: " + err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, errors.New(op + ": Failed to create table: " + err.Error())
	}
	return &Storage{db: db}, nil
}
