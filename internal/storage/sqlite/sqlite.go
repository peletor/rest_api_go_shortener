package sqlite

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"rest_api_shortener/internal/storage"
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

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.saveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES(?, ?)")
	if err != nil {
		return 0, errors.New(op + ": Failed to prepare statement: " + err.Error())
	}

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		// TODO: refactor this
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, errors.New(op + ": Failed to save URL: " + storage.ErrURLExists.Error())
		}

		return 0, errors.New(op + ": Failed to save URL: " + err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New(op + ": Failed to get last inserted Id: " + err.Error())
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.getUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", errors.New(op + ": Failed to prepare statement: " + err.Error())
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", errors.New(op + ": Failed to get URL: " + err.Error())
	}

	return resUrl, nil
}
func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.deleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=?")

	if err != nil {
		return errors.New(op + ": Failed to prepare statement: " + err.Error())
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return errors.New(op + ": Failed to delete: " + err.Error())
	}

	return nil
}
