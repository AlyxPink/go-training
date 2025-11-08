package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/alyxpink/go-training/taskapi/models"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	store *models.TaskStore
}

func NewTaskHandler(store *models.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

func (r *CreateTaskRequest) Validate() error {
	// Check required fields
	if r.Title == "" {
		return errors.New("title is required")
	}

	// Check title length
	if len(r.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}

	// Validate status enum if provided
	if r.Status != "" {
		validStatuses := map[string]bool{
			"pending":     true,
			"in_progress": true,
			"completed":   true,
		}
		if !validStatuses[r.Status] {
			return errors.New("status must be one of: pending, in_progress, completed")
		}
	} else {
		// Set default status
		r.Status = "pending"
	}

	// Validate priority range
	if r.Priority != 0 && (r.Priority < 1 || r.Priority > 5) {
		return errors.New("priority must be between 1 and 5")
	}

	// Set default priority if not provided
	if r.Priority == 0 {
		r.Priority = 3
	}

	return nil
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	if err := h.store.Create(task); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	task, err := h.store.GetByID(id)
	if err == models.ErrNotFound {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priorityStr := r.URL.Query().Get("priority")

	var priority int
	if priorityStr != "" {
		var err error
		priority, err = strconv.Atoi(priorityStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid priority")
			return
		}
	}

	tasks, err := h.store.List(status, priority)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"tasks": tasks,
		"total": len(tasks),
	})
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *int       `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

func (r *UpdateTaskRequest) Validate() error {
	// Validate title if provided
	if r.Title != nil {
		if *r.Title == "" {
			return errors.New("title cannot be empty")
		}
		if len(*r.Title) > 200 {
			return errors.New("title must be less than 200 characters")
		}
	}

	// Validate status if provided
	if r.Status != nil {
		validStatuses := map[string]bool{
			"pending":     true,
			"in_progress": true,
			"completed":   true,
		}
		if !validStatuses[*r.Status] {
			return errors.New("status must be one of: pending, in_progress, completed")
		}
	}

	// Validate priority if provided
	if r.Priority != nil {
		if *r.Priority < 1 || *r.Priority > 5 {
			return errors.New("priority must be between 1 and 5")
		}
	}

	return nil
}

func (r *UpdateTaskRequest) ToMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if r.Title != nil {
		updates["title"] = *r.Title
	}
	if r.Description != nil {
		updates["description"] = *r.Description
	}
	if r.Status != nil {
		updates["status"] = *r.Status
	}
	if r.Priority != nil {
		updates["priority"] = *r.Priority
	}
	if r.DueDate != nil {
		updates["due_date"] = *r.DueDate
	}

	return updates
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	updates := req.ToMap()
	task, err := h.store.Update(id, updates)
	if err == models.ErrNotFound {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	if err := h.store.Delete(id); err == models.ErrNotFound {
		respondError(w, http.StatusNotFound, "task not found")
		return
	} else if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
