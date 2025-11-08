# Hints

## log package
```go
logger := log.New(os.Stdout, "PREFIX: ", log.LstdFlags)
logger.Println("message")
```

## slog (Go 1.21+)
```go
slog.Info("event", "key", "value")
```
