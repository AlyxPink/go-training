# Hints: Context Management

## Creating Contexts

```go
// Background context (never cancelled)
ctx := context.Background()

// With cancellation
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()

// With timeout
ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
defer cancel()

// With deadline
ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(5*time.Second))
defer cancel()
```

## Cancellation Pattern

```go
func Worker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Do work
		}
	}
}
```

## Context Values

```go
type key string
const userKey key = "user"

ctx := context.WithValue(parent, userKey, "alice")
user := ctx.Value(userKey).(string)
```
