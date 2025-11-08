# Hints: HTTP Server

## Basic Server
```go
http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello")
})
http.ListenAndServe(":8080", nil)
```

## Handler Interface
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

## Middleware Pattern
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

## JSON Response
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
```
