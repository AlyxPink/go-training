# Architectural Hints: REST API

## Database Setup

### Schema Design
```sql
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

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
```

### Database Connection
```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", filepath)
    if err != nil {
        return nil, err
    }
    
    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

## Model Layer

### Task Model with CRUD
```go
type TaskStore struct {
    db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
    return &TaskStore{db: db}
}

func (s *TaskStore) Create(task *Task) error {
    query := `
        INSERT INTO tasks (title, description, status, priority, due_date)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id, created_at, updated_at
    `
    
    err := s.db.QueryRow(
        query,
        task.Title,
        task.Description,
        task.Status,
        task.Priority,
        task.DueDate,
    ).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
    
    return err
}

func (s *TaskStore) GetByID(id int64) (*Task, error) {
    query := `
        SELECT id, title, description, status, priority, 
               created_at, updated_at, due_date
        FROM tasks
        WHERE id = ?
    `
    
    task := &Task{}
    err := s.db.QueryRow(query, id).Scan(
        &task.ID,
        &task.Title,
        &task.Description,
        &task.Status,
        &task.Priority,
        &task.CreatedAt,
        &task.UpdatedAt,
        &task.DueDate,
    )
    
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    
    return task, err
}

func (s *TaskStore) List(filters TaskFilters) ([]*Task, error) {
    query := "SELECT * FROM tasks WHERE 1=1"
    args := []interface{}{}
    
    if filters.Status != "" {
        query += " AND status = ?"
        args = append(args, filters.Status)
    }
    
    if filters.Priority > 0 {
        query += " AND priority = ?"
        args = append(args, filters.Priority)
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
        // Scan into task...
        tasks = append(tasks, task)
    }
    
    return tasks, rows.Err()
}
```

## Handler Layer

### Router Setup with Chi
```go
import "github.com/go-chi/chi/v5"

func NewRouter(store *TaskStore) *chi.Mux {
    r := chi.NewRouter()
    
    // Middleware
    r.Use(middleware.RequestLogger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.CORS)
    
    // Routes
    r.Route("/tasks", func(r chi.Router) {
        r.Get("/", listTasks(store))
        r.Post("/", createTask(store))
        r.Get("/{id}", getTask(store))
        r.Put("/{id}", updateTask(store))
        r.Delete("/{id}", deleteTask(store))
    })
    
    return r
}
```

### Handler Pattern
```go
func createTask(store *TaskStore) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req CreateTaskRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondError(w, http.StatusBadRequest, "invalid request body")
            return
        }
        
        // Validate
        if err := req.Validate(); err != nil {
            respondError(w, http.StatusBadRequest, err.Error())
            return
        }
        
        // Convert to model
        task := &Task{
            Title:       req.Title,
            Description: req.Description,
            Status:      req.Status,
            Priority:    req.Priority,
            DueDate:     req.DueDate,
        }
        
        // Create in database
        if err := store.Create(task); err != nil {
            respondError(w, http.StatusInternalServerError, "failed to create task")
            return
        }
        
        respondJSON(w, http.StatusCreated, task)
    }
}
```

### Response Helpers
```go
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(status)
    json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}
```

## Middleware Patterns

### Request Logger
```go
func RequestLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Generate request ID
        requestID := generateRequestID()
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        
        // Wrap response writer to capture status
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r.WithContext(ctx))
        
        log.Printf("[%s] %s %s %d %v",
            requestID,
            r.Method,
            r.URL.Path,
            wrapped.statusCode,
            time.Since(start),
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
    w.statusCode = code
    w.ResponseWriter.WriteHeader(code)
}
```

### Recovery Middleware
```go
func Recoverer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v\n%s", err, debug.Stack())
                respondError(w, http.StatusInternalServerError, "internal server error")
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

## Validation

### Request Validation
```go
type CreateTaskRequest struct {
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    Priority    int        `json:"priority"`
    DueDate     *time.Time `json:"due_date"`
}

func (r *CreateTaskRequest) Validate() error {
    if r.Title == "" {
        return errors.New("title is required")
    }
    
    if len(r.Title) > 200 {
        return errors.New("title must be less than 200 characters")
    }
    
    validStatuses := map[string]bool{
        "pending": true, "in_progress": true, "completed": true,
    }
    if r.Status != "" && !validStatuses[r.Status] {
        return errors.New("invalid status")
    }
    
    if r.Priority < 1 || r.Priority > 5 {
        return errors.New("priority must be between 1 and 5")
    }
    
    return nil
}
```

## Testing

### Integration Test Setup
```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatal(err)
    }
    
    return db
}

func TestCreateTask(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    store := NewTaskStore(db)
    router := NewRouter(store)
    
    payload := `{"title": "Test Task", "priority": 3}`
    req := httptest.NewRequest("POST", "/tasks", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
    
    var task Task
    json.NewDecoder(rr.Body).Decode(&task)
    assert.Equal(t, "Test Task", task.Title)
    assert.NotZero(t, task.ID)
}
```

## Common Patterns

### Error Handling
```go
var (
    ErrNotFound = errors.New("task not found")
    ErrInvalidInput = errors.New("invalid input")
)

// In handler
task, err := store.GetByID(id)
if err == ErrNotFound {
    respondError(w, http.StatusNotFound, "task not found")
    return
}
if err != nil {
    respondError(w, http.StatusInternalServerError, "internal error")
    return
}
```

### Transaction Pattern
```go
func (s *TaskStore) UpdateWithHistory(id int64, updates map[string]interface{}) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Update task
    if err := updateTask(tx, id, updates); err != nil {
        return err
    }
    
    // Create history entry
    if err := createHistory(tx, id, updates); err != nil {
        return err
    }
    
    return tx.Commit()
}
```
