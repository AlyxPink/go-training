package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// User represents a user
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserStore manages users
type UserStore struct {
	mu    sync.RWMutex
	users map[int]User
	nextID int
}

// NewUserStore creates a new user store
func NewUserStore() *UserStore {
	// TODO: Initialize store
	return nil
}

// GetAll returns all users
func (s *UserStore) GetAll() []User {
	// TODO: Return all users (thread-safe)
	return nil
}

// Get returns user by ID
func (s *UserStore) Get(id int) (User, bool) {
	// TODO: Get user (thread-safe)
	return User{}, false
}

// Create adds a new user
func (s *UserStore) Create(name string) User {
	// TODO: Create user (thread-safe)
	return User{}
}

// Delete removes a user
func (s *UserStore) Delete(id int) bool {
	// TODO: Delete user (thread-safe)
	return false
}

// Server wraps HTTP server
type Server struct {
	store *UserStore
}

// NewServer creates a new server
func NewServer() *Server {
	// TODO: Initialize server
	return nil
}

// HandleUsers handles /users endpoint
func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle GET (list) and POST (create)
}

// HandleUser handles /users/:id endpoint
func (s *Server) HandleUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle GET (get by ID) and DELETE
}

// LoggingMiddleware logs requests
func LoggingMiddleware(next http.Handler) http.Handler {
	// TODO: Log method and path, then call next
	return nil
}

func main() {
	server := NewServer()
	if server == nil {
		fmt.Println("Implement NewServer first!")
		return
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/users", server.HandleUsers)
	mux.HandleFunc("/users/", server.HandleUser)
	
	handler := LoggingMiddleware(mux)
	
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
