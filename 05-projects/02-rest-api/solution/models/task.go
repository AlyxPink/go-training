package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNotFound     = errors.New("task not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) Create(task *Task) error {
	// Set default status if not provided
	if task.Status == "" {
		task.Status = "pending"
	}

	// Set default priority if not provided
	if task.Priority == 0 {
		task.Priority = 3
	}

	query := `
		INSERT INTO tasks (title, description, status, priority, due_date)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(query, task.Title, task.Description, task.Status, task.Priority, task.DueDate)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = id

	// Fetch created_at and updated_at timestamps
	return s.db.QueryRow("SELECT created_at, updated_at FROM tasks WHERE id = ?", id).
		Scan(&task.CreatedAt, &task.UpdatedAt)
}

func (s *TaskStore) GetByID(id int64) (*Task, error) {
	task := &Task{}
	query := `
		SELECT id, title, description, status, priority, created_at, updated_at, due_date
		FROM tasks WHERE id = ?
	`

	err := s.db.QueryRow(query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status,
		&task.Priority, &task.CreatedAt, &task.UpdatedAt, &task.DueDate,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return task, err
}

func (s *TaskStore) List(status string, priority int) ([]*Task, error) {
	query := "SELECT id, title, description, status, priority, created_at, updated_at, due_date FROM tasks WHERE 1=1"
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if priority > 0 {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status,
			&task.Priority, &task.CreatedAt, &task.UpdatedAt, &task.DueDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (s *TaskStore) Update(id int64, updates map[string]interface{}) (*Task, error) {
	// Check if task exists first
	_, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Build dynamic UPDATE query
	setClauses := []string{}
	args := []interface{}{}

	for key, value := range updates {
		switch key {
		case "title", "description", "status", "priority", "due_date":
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
			args = append(args, value)
		}
	}

	if len(setClauses) == 0 {
		// No valid fields to update, just return the existing task
		return s.GetByID(id)
	}

	// Always update the updated_at timestamp
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	args = append(args, id)

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Return the updated task
	return s.GetByID(id)
}

func (s *TaskStore) Delete(id int64) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
