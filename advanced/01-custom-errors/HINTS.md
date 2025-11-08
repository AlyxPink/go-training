# Hints for Custom Errors

## Hint 1: Custom Error Type Structure

Custom errors should embed context and optionally wrap another error:

```go
type ValidationError struct {
    Fields map[string]string // field -> error message
    err    error             // wrapped error
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed: %d fields", len(e.Fields))
}
```

## Hint 2: Error Wrapping

Implement `Unwrap()` to support error chain inspection:

```go
func (e *ValidationError) Unwrap() error {
    return e.err
}
```

## Hint 3: Sentinel Errors

Define package-level sentinel errors for common conditions:

```go
var (
    ErrNotFound   = errors.New("file not found")
    ErrPermission = errors.New("permission denied")
)
```

## Hint 4: Error Context

Add context when wrapping errors using `fmt.Errorf()`:

```go
if err != nil {
    return fmt.Errorf("failed to process %s: %w", filename, err)
}
```

## Hint 5: Stack Trace Capture

For debugging, capture the call stack at error creation:

```go
type ErrorWithStack struct {
    msg   string
    stack []uintptr
}

func NewErrorWithStack(msg string) *ErrorWithStack {
    pc := make([]uintptr, 32)
    n := runtime.Callers(2, pc)
    return &ErrorWithStack{msg: msg, stack: pc[:n]}
}
```
