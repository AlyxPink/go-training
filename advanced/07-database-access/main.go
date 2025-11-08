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
	result, err := r.db.Exec(
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Age,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *UserRepository) Get(id int64) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, age FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *User) error {
	_, err := r.db.Exec(
		"UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?",
		user.Name, user.Email, user.Age, user.ID,
	)
	return err
}

func (r *UserRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (r *UserRepository) List() ([]*User, error) {
	rows, err := r.db.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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
