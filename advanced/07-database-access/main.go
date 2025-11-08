package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    int64
	Name  string
	Email string
	Age   int
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return db, nil
}

func (r *UserRepository) Create(user *User) error {
	// TODO: Execute INSERT query with user.Name, user.Email, user.Age
	// TODO: Check for errors
	// TODO: Get LastInsertId from result
	// TODO: Set user.ID to the returned id
	// TODO: Return nil on success
	panic("not implemented")
}

func (r *UserRepository) Get(id int64) (*User, error) {
	// TODO: Create new User struct
	// TODO: QueryRow with SELECT query using id parameter
	// TODO: Scan result into user fields
	// TODO: Return user pointer and error
	panic("not implemented")
}

func (r *UserRepository) Update(user *User) error {
	// TODO: Execute UPDATE query with user fields
	// TODO: Use user.ID in WHERE clause
	// TODO: Return error from Exec
	panic("not implemented")
}

func (r *UserRepository) Delete(id int64) error {
	// TODO: Execute DELETE query with id parameter
	// TODO: Return error from Exec
	panic("not implemented")
}

func (r *UserRepository) List() ([]*User, error) {
	// TODO: Query all users from database
	// TODO: Create slice to hold users
	// TODO: Iterate rows with rows.Next()
	// TODO: Scan each row into User struct
	// TODO: Append to users slice
	// TODO: Check rows.Err() after iteration
	// TODO: Return users slice and error
	panic("not implemented")
}

func (r *UserRepository) Transaction(fn func(*sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) CreateMultiple(users []*User) error {
	return r.Transaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, user := range users {
			result, err := stmt.Exec(user.Name, user.Email, user.Age)
			if err != nil {
				return err
			}

			id, err := result.LastInsertId()
			if err != nil {
				return err
			}
			user.ID = id
		}

		return nil
	})
}

func main() {
	db, err := InitDB("users.db")
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

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

	retrieved, err := repo.Get(user.ID)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}
	fmt.Printf("Retrieved user: %+v\n", retrieved)

	retrieved.Age = 31
	if err := repo.Update(retrieved); err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}

	users, err := repo.List()
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		return
	}
	fmt.Printf("All users: %+v\n", users)

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
