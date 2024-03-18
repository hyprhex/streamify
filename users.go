package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	errUsernameRequired = errors.New("username is required")
	errEmailRequired    = errors.New("email is required")
	errPasswordRequired = errors.New("password is required")
)

type UsersService struct {
	store Store
}

func NewUsersService(s Store) *UsersService {
	return &UsersService{store: s}
}

func (s *UsersService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleRegisterUser).Methods("POST")
	r.HandleFunc("/users/login", s.handleLoginUser).Methods("POST")
}

func (s *UsersService) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Error reading request body"})
		return
	}

	defer r.Body.Close()

	user := new(User)
	err = json.Unmarshal(body, user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	err = validateUserPayload(user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Incorrect username or password"})
		return
	}

	user.Password = string(hash)

	u, err := s.store.CreateUser(user)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error Creating user"})
		return
	}

	WriteJSON(w, http.StatusCreated, u)
}

func (s *UsersService) handleLoginUser(w http.ResponseWriter, r *http.Request) {
}

func validateUserPayload(user *User) error {
	if user.Username == "" {
		return errUsernameRequired
	}

	if user.Email == "" {
		return errEmailRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}
