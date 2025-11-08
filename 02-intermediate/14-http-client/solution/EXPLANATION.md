# Solution Explanation: HTTP Client

## Basic HTTP Client

### GET Request
```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()  // Important!

body, err := io.ReadAll(resp.Body)
```

**Always:**
- Check error
- Close response body
- Read body completely

### POST Request
```go
data := map[string]string{"key": "value"}
jsonData, _ := json.Marshal(data)

resp, err := http.Post(
    url,
    "application/json",
    bytes.NewBuffer(jsonData),
)
```

## Custom HTTP Client

### Why Custom Client?
```go
client := &http.Client{
    Timeout: 10 * time.Second,  // Prevent hanging
    Transport: &http.Transport{
        MaxIdleConns: 100,
        IdleConnTimeout: 90 * time.Second,
    },
}
```

**Benefits:**
- Control timeouts
- Connection pooling
- Custom transports

## Custom Requests

```go
req, err := http.NewRequest("GET", url, nil)
req.Header.Set("Authorization", "Bearer token")
req.Header.Set("User-Agent", "MyApp/1.0")

resp, err := client.Do(req)
```

**When to use:**
- Custom headers
- Non-standard methods
- Request body
- Context

## Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := client.Do(req)
```

**Benefits:**
- Cancellation
- Deadlines
- Request-scoped values

## Error Handling

### Check Status Code
```go
if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("bad status: %d", resp.StatusCode)
}
```

### Network Errors
```go
resp, err := client.Get(url)
if err != nil {
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        // Timeout error
    }
    return err
}
```

## JSON Handling

### Request
```go
data := MyStruct{Field: "value"}
jsonData, _ := json.Marshal(data)
req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
req.Header.Set("Content-Type", "application/json")
```

### Response
```go
var result MyStruct
json.NewDecoder(resp.Body).Decode(&result)
```

## Best Practices

1. **Always close response body**: Even on errors
2. **Set timeouts**: Prevent hanging forever
3. **Reuse clients**: Connection pooling
4. **Check status codes**: Don't assume success
5. **Use context**: For cancellation and deadlines
6. **Read body completely**: Or connections leak

## Common Patterns

### Retry Logic
```go
func GetWithRetry(url string, maxRetries int) (*http.Response, error) {
    for i := 0; i < maxRetries; i++ {
        resp, err := http.Get(url)
        if err == nil {
            return resp, nil
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return nil, fmt.Errorf("max retries exceeded")
}
```

### Rate Limiting
```go
type RateLimitedClient struct {
    client *http.Client
    limiter *rate.Limiter
}

func (c *RateLimitedClient) Get(url string) (*http.Response, error) {
    c.limiter.Wait(context.Background())
    return c.client.Get(url)
}
```

### API Client Wrapper
```go
type APIClient struct {
    baseURL string
    client  *http.Client
    token   string
}

func (c *APIClient) request(method, path string, body io.Reader) ([]byte, error) {
    req, _ := http.NewRequest(method, c.baseURL+path, body)
    req.Header.Set("Authorization", "Bearer "+c.token)
    
    resp, err := c.client.Do(req)
    // Handle response...
}
```

## Testing

### httptest Package
```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("response"))
}))
defer server.Close()

// Use server.URL for tests
```

**Benefits:**
- No external dependencies
- Fast tests
- Controlled responses
