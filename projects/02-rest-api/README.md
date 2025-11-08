# Project 2: Task Management REST API

**Difficulty**: ⭐⭐⭐⭐ | **Estimated Time**: 180 minutes

## Overview

Build a complete REST API for task management with CRUD operations, SQLite persistence, middleware chain, input validation, and proper error handling.

## Architecture

```
┌──────────┐
│  HTTP    │  (chi/mux router, middleware)
│  Layer   │
└────┬─────┘
     │
┌────▼─────┐
│ Handlers │  (request/response, validation)
└────┬─────┘
     │
┌────▼─────┐
│  Models  │  (business logic, CRUD operations)
└────┬─────┘
     │
┌────▼─────┐
│ Database │  (SQLite, migrations)
└──────────┘
```

## Features to Implement

### 1. CRUD Operations
- **CREATE**: `POST /tasks` - Create new task
- **READ**: `GET /tasks` - List all tasks
- **READ**: `GET /tasks/:id` - Get task by ID
- **UPDATE**: `PUT /tasks/:id` - Update task
- **DELETE**: `DELETE /tasks/:id` - Delete task

### 2. Task Model
```go
type Task struct {
    ID          int64     `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"` // pending, in_progress, completed
    Priority    int       `json:"priority"` // 1-5
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DueDate     *time.Time `json:"due_date,omitempty"`
}
```

### 3. Middleware Chain
- Request logging
- Request ID generation
- CORS headers
- Recovery from panics
- Content-Type validation
- Rate limiting (bonus)

### 4. Validation
- Required fields (title)
- Status enum validation
- Priority range (1-5)
- Due date format
- Title length (1-200 chars)

### 5. Error Handling
- Standardized error responses
- HTTP status codes
- Validation error details
- Database error handling

## API Specifications

### Create Task
```http
POST /tasks
Content-Type: application/json

{
  "title": "Implement user authentication",
  "description": "Add JWT-based auth",
  "status": "pending",
  "priority": 3,
  "due_date": "2024-12-31T23:59:59Z"
}

Response: 201 Created
{
  "id": 1,
  "title": "Implement user authentication",
  "description": "Add JWT-based auth",
  "status": "pending",
  "priority": 3,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z",
  "due_date": "2024-12-31T23:59:59Z"
}
```

### List Tasks
```http
GET /tasks?status=pending&priority=3

Response: 200 OK
{
  "tasks": [...],
  "total": 10,
  "page": 1,
  "per_page": 20
}
```

### Get Task
```http
GET /tasks/1

Response: 200 OK
{
  "id": 1,
  "title": "...",
  ...
}

Response: 404 Not Found
{
  "error": "task not found"
}
```

### Update Task
```http
PUT /tasks/1
Content-Type: application/json

{
  "status": "completed"
}

Response: 200 OK
{
  "id": 1,
  "status": "completed",
  ...
}
```

### Delete Task
```http
DELETE /tasks/1

Response: 204 No Content
```

## Requirements

### Database
- SQLite for simplicity
- Proper schema with indexes
- Migrations for schema changes
- Connection pooling
- Transaction support

### HTTP Layer
- RESTful design principles
- Proper status codes
- JSON request/response
- Error responses with details
- Query parameters for filtering

### Testing
- Unit tests for models
- Integration tests for handlers
- Table-driven tests
- Test database isolation
- HTTP test utilities

## Project Structure

```
02-rest-api/
├── README.md
├── HINTS.md
├── go.mod
├── main.go              # Server setup
├── handlers/
│   ├── tasks.go         # HTTP handlers (TODO)
│   └── errors.go        # Error responses
├── models/
│   ├── task.go          # Task model & CRUD (TODO)
│   └── database.go      # DB connection
├── middleware/
│   ├── logging.go       # Request logging (TODO)
│   ├── recovery.go      # Panic recovery
│   └── cors.go          # CORS headers
├── migrations/
│   └── 001_create_tasks.sql
├── main_test.go         # Integration tests
└── solution/
    ├── ARCHITECTURE.md
    └── [all files]
```

## Test Cases

```go
// Create task
POST /tasks {"title": "Test"} → 201, task with ID

// Get task
GET /tasks/1 → 200, task data
GET /tasks/999 → 404

// List tasks
GET /tasks → 200, array of tasks
GET /tasks?status=pending → 200, filtered tasks

// Update task
PUT /tasks/1 {"status": "completed"} → 200, updated task
PUT /tasks/999 {...} → 404

// Delete task
DELETE /tasks/1 → 204
DELETE /tasks/1 → 404 (already deleted)

// Validation
POST /tasks {} → 400, "title required"
POST /tasks {"title": "", ...} → 400, "title cannot be empty"
POST /tasks {"title": "x", "status": "invalid"} → 400, "invalid status"
POST /tasks {"title": "x", "priority": 10} → 400, "priority must be 1-5"
```

## Grading Criteria

- **Correctness** (30%): All endpoints work correctly
- **Database** (20%): Proper schema, migrations, transactions
- **Middleware** (15%): Complete middleware chain
- **Validation** (15%): Comprehensive input validation
- **Error Handling** (10%): Proper error responses
- **Code Quality** (10%): Clean, idiomatic Go

## Bonus Challenges

1. Add pagination to list endpoint
2. Implement task filtering by multiple criteria
3. Add sorting (by created_at, priority, etc.)
4. Implement soft deletes
5. Add search functionality
6. Implement rate limiting middleware
7. Add OpenAPI/Swagger documentation
8. Implement database migration tool

## Technical Concepts

1. **HTTP Routing**: chi/mux, route parameters
2. **Middleware**: composition, request/response manipulation
3. **JSON Handling**: encoding/decoding, custom marshaling
4. **SQL**: CRUD operations, prepared statements, transactions
5. **Validation**: input validation, error messages
6. **Testing**: httptest, test databases, table-driven tests

## Getting Started

1. Review HINTS.md for implementation guidance
2. Set up database schema
3. Implement Task model with CRUD methods
4. Build HTTP handlers
5. Add middleware chain
6. Write comprehensive tests

## Example Session

```bash
# Start server
$ go run main.go
Server listening on :8080

# Create task
$ curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Write tests", "priority": 2}'
{"id": 1, "title": "Write tests", "status": "pending", ...}

# List tasks
$ curl http://localhost:8080/tasks
{"tasks": [...], "total": 1}

# Update task
$ curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
{"id": 1, "status": "completed", ...}

# Delete task
$ curl -X DELETE http://localhost:8080/tasks/1
```

## Learning Outcomes

After completing this project, you will understand:
- RESTful API design principles
- HTTP middleware patterns
- SQL database integration
- Input validation strategies
- Error handling best practices
- Integration testing for HTTP services
