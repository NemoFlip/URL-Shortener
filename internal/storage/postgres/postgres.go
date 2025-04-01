package postgres

import (
	"RESTProject/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(dataSourceName string) (*Storage, error) {
	const fn = "storage.postgres.NewStorage"

	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &Storage{DB: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	const fn = "storage.postgres.SaveURL"
	query := "INSERT INTO url (url, alias) VALUES ($1, $2)"

	_, err := s.DB.Exec(query, urlToSave, alias)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // code 23505 — violation of UNIQUE
			return fmt.Errorf("%s: %w", fn, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.postgres.GetURL"
	query := "SELECT url FROM url WHERE alias = $1"
	var url string
	err := s.DB.QueryRow(query, alias).Scan(&url)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const fn = "storage.postgres.DeleteURL"
	query := "DELETE FROM url WHERE alias = $1"
	res, err := s.DB.Exec(query, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", fn, err)
	}
	if rowsAffected == 0 {
		return storage.ErrURLNotFound
	}
	return nil
}

func (s *Storage) UpdateURL(alias, newURL string) error {
	const fn = "storage.postgres.UpdateURL"
	query := "UPDATE url SET url = $1 WHERE alias = $2"

	res, err := s.DB.Exec(query, newURL, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", fn, err)
	}
	if rowsAffected == 0 {
		return storage.ErrURLNotFound
	}
	return nil
}
