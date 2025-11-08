# Exercise 15: HTTP Server

**Difficulty**: ⭐⭐⭐ Intermediate-Advanced
**Estimated Time**: 75 minutes

## Learning Objectives

- Build HTTP servers with net/http
- Implement REST API endpoints
- Use http.Handler and http.HandlerFunc
- Create middleware patterns
- Handle routing and requests

## Problem Description

Build a simple REST API server with middleware.

### Requirements

1. **GET /users** - List all users
2. **GET /users/:id** - Get user by ID
3. **POST /users** - Create new user
4. **DELETE /users/:id** - Delete user
5. **Middleware** - Logging, authentication

### Expected Behavior

```bash
curl http://localhost:8080/users
curl http://localhost:8080/users/1
curl -X POST http://localhost:8080/users -d '{"name":"Alice"}'
curl -X DELETE http://localhost:8080/users/1
```

## Testing
```bash
go test -v
go run main.go  # Start server
```
