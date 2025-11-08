package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/alyxpink/go-training/taskapi/handlers"
	"github.com/alyxpink/go-training/taskapi/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: Initialize database
	db, err := initDB("tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: Create task store
	store := models.NewTaskStore(db)

	// TODO: Setup router with middleware
	r := setupRouter(store)

	// TODO: Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func initDB(filepath string) (*sql.DB, error) {
	// TODO: Open SQLite database
	// TODO: Run migrations
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	// TODO: Create tasks table
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		priority INTEGER NOT NULL DEFAULT 3,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		due_date DATETIME,
		CONSTRAINT status_check CHECK (status IN ('pending', 'in_progress', 'completed')),
		CONSTRAINT priority_check CHECK (priority BETWEEN 1 AND 5)
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	`

	_, err := db.Exec(schema)
	return err
}

func setupRouter(store *models.TaskStore) *chi.Mux {
	r := chi.NewRouter()

	// TODO: Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// TODO: Define routes
	h := handlers.NewTaskHandler(store)
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})

	return r
}
