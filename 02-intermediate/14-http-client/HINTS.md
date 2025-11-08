# Hints: HTTP Client

## Basic GET
```go
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    return err
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
```

## POST with JSON
```go
data := map[string]string{"key": "value"}
jsonData, _ := json.Marshal(data)
resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
```

## Custom Client
```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("Authorization", "Bearer token")
resp, err := client.Do(req)
```
