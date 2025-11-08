package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alyxpink/go-training/taskapi/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	
	err = runMigrations(db)
	require.NoError(t, err)
	
	return db
}

func TestCreateTask(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - implement handlers and models")

	db := setupTestDB(t)
	defer db.Close()

	store := models.NewTaskStore(db)
	router := setupRouter(store)
	
	payload := `{"title": "Test Task", "status": "pending", "priority": 3}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusCreated, rr.Code)
	
	var task models.Task
	err := json.NewDecoder(rr.Body).Decode(&task)
	require.NoError(t, err)
	assert.Equal(t, "Test Task", task.Title)
	assert.NotZero(t, task.ID)
}

func TestGetTask(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - implement handlers and models")

	db := setupTestDB(t)
	defer db.Close()
	
	store := models.NewTaskStore(db)
	router := setupRouter(store)
	
	// Create a task first
	task := &models.Task{Title: "Test", Status: "pending", Priority: 3}
	err := store.Create(task)
	require.NoError(t, err)
	
	// Get the task
	req := httptest.NewRequest("GET", "/tasks/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var retrieved models.Task
	err = json.NewDecoder(rr.Body).Decode(&retrieved)
	require.NoError(t, err)
	assert.Equal(t, "Test", retrieved.Title)
}

func TestListTasks(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - implement handlers and models")

	db := setupTestDB(t)
	defer db.Close()
	
	store := models.NewTaskStore(db)
	router := setupRouter(store)
	
	// Create some tasks
	for i := 0; i < 3; i++ {
		task := &models.Task{Title: "Task", Status: "pending", Priority: i + 1}
		store.Create(task)
	}
	
	req := httptest.NewRequest("GET", "/tasks", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, float64(3), response["total"])
}

func TestDeleteTask(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - implement handlers and models")

	db := setupTestDB(t)
	defer db.Close()
	
	store := models.NewTaskStore(db)
	router := setupRouter(store)
	
	// Create a task
	task := &models.Task{Title: "Delete Me", Status: "pending", Priority: 1}
	store.Create(task)
	
	// Delete it
	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusNoContent, rr.Code)
	
	// Verify it's gone
	_, err := store.GetByID(1)
	assert.Equal(t, models.ErrNotFound, err)
}
