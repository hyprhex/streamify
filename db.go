package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(cfg mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")

	return &MySQLStorage{
		db: db,
	}
}

func (s *MySQLStorage) Init() (*sql.DB, error) {
	if err := s.createUsersTable(); err != nil {
		return nil, err
	}

	if err := s.createFeedsTable(); err != nil {
		return nil, err
	}

	if err := s.createUsersFeedsTable(); err != nil {
		return nil, err
	}

	if err := s.createPostsTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *MySQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		username VARCHAR(255) NOT NULL, 
		email VARCHAR(255) NOT NULL, 
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		UNIQUE (username, email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createFeedsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS feeds (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		title VARCHAR(255) NOT NULL, 
		url VARCHAR(255) NOT NULL, 
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createUsersFeedsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users_feeds (
		user_id INT UNSIGNED NOT NULL,
		feed_id INT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (user_id, feed_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createPostsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		feed_id INT UNSIGNED NOT NULL,
		title VARCHAR(255) NOT NULL, 
		url VARCHAR(255) NOT NULL, 
		description TEXT,
		pub_date TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}
