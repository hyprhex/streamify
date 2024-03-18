package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	errTitleRequired = errors.New("title is required")
	errURLRequired   = errors.New("URL is required")
)

type FeedsService struct {
	store Store
}

func NewFeedsServices(s Store) *FeedsService {
	return &FeedsService{store: s}
}

func (s *FeedsService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/feeds", s.handleCreateFeed).Methods("POST")
	r.HandleFunc("/feeds", s.handleListFeed).Methods("GET")
}

func (s *FeedsService) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	defer r.Body.Close()

	feed := new(Feed)
	err = json.Unmarshal(body, feed)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	err = validateFeedPayload(feed)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	f, err := s.store.CreateFeed(feed)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating feed"})
		return
	}

	WriteJSON(w, http.StatusCreated, f)
}

func (s *FeedsService) handleListFeed(w http.ResponseWriter, r *http.Request) {
	f, err := s.store.ListFeed()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusOK, f)
}

func validateFeedPayload(feed *Feed) error {
	if feed.Title == "" {
		return errTitleRequired
	}

	if feed.URL == "" {
		return errURLRequired
	}

	return nil
}
