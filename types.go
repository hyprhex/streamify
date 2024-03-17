package main

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type Feed struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
