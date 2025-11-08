package handlers

import (
	"encoding/json"
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
	// TODO: Implement validation
	// Hint: Check required fields, value ranges, enum values
	return nil
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode request, validate, create task, respond
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
	// TODO: Parse ID from URL, get task, respond
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
	// TODO: Parse query parameters, filter tasks, respond
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

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse ID, decode updates, update task, respond
	respondError(w, http.StatusNotImplemented, "update not implemented")
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse ID, delete task, respond 204
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
