package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// User represents a user in the database
type User struct {
	ID    int64
	Name  string
	Email string
	Age   int
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// InitDB initializes the database schema
func InitDB(dbPath string) (*sql.DB, error) {
	// TODO: Open SQLite database
	// TODO: Create users table if not exists
	// CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT, age INTEGER)
	return nil, nil
}

// Create inserts a new user
func (r *UserRepository) Create(user *User) error {
	// TODO: Insert user into database
	// TODO: Use prepared statement
	// TODO: Get last insert ID and set user.ID
	return nil
}

// Get retrieves a user by ID
func (r *UserRepository) Get(id int64) (*User, error) {
	// TODO: Query user by ID
	// TODO: Scan result into User struct
	// TODO: Return sql.ErrNoRows if not found
	return nil, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *User) error {
	// TODO: Update user in database
	// TODO: Use prepared statement
	return nil
}

// Delete removes a user by ID
func (r *UserRepository) Delete(id int64) error {
	// TODO: Delete user from database
	return nil
}

// List retrieves all users
func (r *UserRepository) List() ([]*User, error) {
	// TODO: Query all users
	// TODO: Scan all rows into slice
	return nil, nil
}

// Transaction demonstrates transaction usage
func (r *UserRepository) Transaction(fn func(*sql.Tx) error) error {
	// TODO: Begin transaction
	// TODO: Execute function
	// TODO: Commit on success, rollback on error
	return nil
}

// CreateMultiple creates multiple users in a transaction
func (r *UserRepository) CreateMultiple(users []*User) error {
	// TODO: Use Transaction to insert multiple users atomically
	// TODO: All should succeed or all fail
	return nil
}

func main() {
	db, err := InitDB("users.db")
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Create user
	user := &User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}
	if err := repo.Create(user); err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Get user
	retrieved, err := repo.Get(user.ID)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}
	fmt.Printf("Retrieved user: %+v\n", retrieved)

	// Update user
	retrieved.Age = 31
	if err := repo.Update(retrieved); err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}

	// List all users
	users, err := repo.List()
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		return
	}
	fmt.Printf("All users: %+v\n", users)

	// Create multiple users in transaction
	newUsers := []*User{
		{Name: "Bob", Email: "bob@example.com", Age: 25},
		{Name: "Charlie", Email: "charlie@example.com", Age: 35},
	}
	if err := repo.CreateMultiple(newUsers); err != nil {
		fmt.Printf("Error creating users: %v\n", err)
		return
	}
	fmt.Println("Created multiple users successfully")
}
