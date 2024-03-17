package main

import "database/sql"

type Store interface {
	// Feeds
	CreateFeed(f *Feed) (*Feed, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateFeed(f *Feed) (*Feed, error) {
	rows, err := s.db.Exec(`
		INSERT INTO feeds (title, url)
		VALUES (?, ?)
		`, f.Title, f.URL)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	f.ID = id
	return f, nil
}
