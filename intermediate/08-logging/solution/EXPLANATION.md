# Logging

## Standard Logger
```go
logger := log.New(os.Stdout, "PREFIX: ", log.LstdFlags)
logger.Println("message")
```

Flags: `log.Ldate | log.Ltime | log.Lshortfile`

## Structured Logging (slog)
```go
slog.Info("event", "user", "alice", "action", "login")
```

Benefits:
- Machine parseable
- Queryable
- Consistent format

## Levels
INFO, WARN, ERROR, DEBUG - implement custom or use slog levels
