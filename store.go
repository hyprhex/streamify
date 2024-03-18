package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)
	LoginUser(u *LoginUserRequest) (*LoginUserRequest, error)

	// Feeds
	CreateFeed(f *Feed) (*Feed, error)
	ListFeed() ([]Feed, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec(`
		INSERT INTO users (username, email, password) 
		VALUES (?, ?, ?)
		`, u.Username, u.Email, u.Password)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (s *Storage) LoginUser(u *LoginUserRequest) (*LoginUserRequest, error) {
	rows, err := s.db.Query(`
		SELECT username, password FROM users WHERE username = ?
	`, u.Username)
	if err != nil {
		return nil, err
	}

	userPW := u.Password
	for rows.Next() {
		err = rows.Scan(&u.Username, &u.Password)
		if err != nil {
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userPW))
	if err != nil {
		return nil, err
	}

	return u, nil
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

func (s *Storage) ListFeed() ([]Feed, error) {
	var feeds []Feed
	rows, err := s.db.Query(`
		SELECT * FROM feeds
		`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var f Feed
		if err := rows.Scan(&f.ID, &f.Title, &f.URL, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}

		feeds = append(feeds, f)
	}

	return feeds, nil
}
