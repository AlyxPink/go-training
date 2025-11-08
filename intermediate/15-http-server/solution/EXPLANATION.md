# Solution Explanation: HTTP Server

## HTTP Handler Interface

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**Two ways to implement:**

1. **HandlerFunc**: Function that matches signature
```go
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello")
}
http.HandleFunc("/path", handler)
```

2. **Custom type**: Implement ServeHTTP method
```go
type MyHandler struct {}
func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
http.Handle("/path", &MyHandler{})
```

## Routing

### ServeMux
```go
mux := http.NewServeMux()
mux.HandleFunc("/users", handleUsers)
mux.HandleFunc("/users/", handleUser)  // Trailing slash = subtree
```

**Pattern matching:**
- Exact match: `/users`
- Subtree: `/users/` matches `/users/1`, `/users/2`, etc.

### Manual Routing
```go
func router(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/users":
        handleUsers(w, r)
    case "/products":
        handleProducts(w, r)
    default:
        http.NotFound(w, r)
    }
}
```

## Request Handling

### Read JSON Body
```go
var data MyStruct
if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}
```

### Write JSON Response
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(data)
```

### Method Routing
```go
switch r.Method {
case http.MethodGet:
    // Handle GET
case http.MethodPost:
    // Handle POST
default:
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
```

## Middleware Pattern

```go
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Before request
        log.Println("Before")
        
        next.ServeHTTP(w, r)  // Call next handler
        
        // After request
        log.Println("After")
    })
}
```

**Chain middleware:**
```go
handler := middleware1(middleware2(middleware3(finalHandler)))
```

### Common Middleware

**Logging:**
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

**Authentication:**
```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !isValid(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**CORS:**
```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}
```

## Thread Safety

**Problem:** Multiple requests access shared state

**Solution:** Use sync.RWMutex
```go
type Store struct {
    mu   sync.RWMutex
    data map[int]Item
}

func (s *Store) Get(id int) (Item, bool) {
    s.mu.RLock()         // Read lock
    defer s.mu.RUnlock()
    item, ok := s.data[id]
    return item, ok
}

func (s *Store) Set(id int, item Item) {
    s.mu.Lock()          // Write lock
    defer s.mu.Unlock()
    s.data[id] = item
}
```

## Status Codes

```go
w.WriteHeader(http.StatusOK)           // 200
w.WriteHeader(http.StatusCreated)      // 201
w.WriteHeader(http.StatusNoContent)    // 204
w.WriteHeader(http.StatusBadRequest)   // 400
w.WriteHeader(http.StatusNotFound)     // 404
w.WriteHeader(http.StatusInternalServerError) // 500
```

**Important:** Call WriteHeader before writing body

## Testing

### httptest Package
```go
req := httptest.NewRequest("GET", "/users", nil)
w := httptest.NewRecorder()

handler(w, req)

if w.Code != http.StatusOK {
    t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
}
```

## Best Practices

1. **Use middleware for cross-cutting concerns**: Logging, auth, CORS
2. **Lock shared state**: Use mutexes for concurrent access
3. **Set Content-Type**: Always set response content type
4. **Handle all methods**: Return 405 for unsupported methods
5. **Validate input**: Check request body and parameters
6. **Use http.Error**: For error responses
7. **Graceful shutdown**: Handle SIGINT/SIGTERM

## Graceful Shutdown

```go
srv := &http.Server{
    Addr:    ":8080",
    Handler: handler,
}

go func() {
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}()

// Wait for interrupt
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt)
<-c

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

## REST API Patterns

```
GET    /users       - List all
GET    /users/:id   - Get one
POST   /users       - Create
PUT    /users/:id   - Update
DELETE /users/:id   - Delete
```

## Further Exploration

- Third-party routers: gorilla/mux, chi, httprouter
- Validation libraries
- OpenAPI/Swagger documentation
- Rate limiting
- Request tracing
