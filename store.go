package main

import "database/sql"

type Store interface {
	// Feeds
	CreateFeed() error
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateFeed() error {
	return nil
}
