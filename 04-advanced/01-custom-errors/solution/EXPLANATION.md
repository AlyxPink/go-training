# Custom Errors - Solution Explanation

## Key Concepts

### 1. Custom Error Types

Custom error types provide rich context beyond simple error messages:

```go
type FileError struct {
    Op        string    // What operation failed
    Path      string    // Which file was involved
    Err       error     // Underlying cause
    Timestamp time.Time // When it happened
}
```

**Benefits**:
- Structured error information
- Type-safe error handling
- Machine-readable error details
- Debuggability through context

### 2. Error Wrapping with `Unwrap()`

Implementing `Unwrap()` enables error chain traversal:

```go
func (e *FileError) Unwrap() error {
    return e.Err
}
```

This allows `errors.Is()` and `errors.As()` to inspect the entire error chain.

### 3. Sentinel Errors

Package-level errors for identity-based comparison:

```go
var ErrNotFound = errors.New("file not found")

// Usage
if errors.Is(err, ErrNotFound) {
    // Handle specifically
}
```

**When to use**:
- Common error conditions
- Error identity checks
- Simple errors without additional context

### 4. Error Inspection Patterns

**`errors.Is()`** - Check if error matches or wraps a target:
```go
if errors.Is(err, ErrPermission) {
    // Error chain contains ErrPermission
}
```

**`errors.As()`** - Extract specific error type from chain:
```go
var fileErr *FileError
if errors.As(err, &fileErr) {
    // Access fileErr.Path, fileErr.Op, etc.
}
```

### 5. Stack Trace Capture

Using `runtime.Callers()` to capture call stack:

```go
pc := make([]uintptr, 32)
n := runtime.Callers(2, pc) // Skip 2 frames
stack := pc[:n]
```

Convert to readable format with `runtime.CallersFrames()`:

```go
frames := runtime.CallersFrames(stack)
for {
    frame, more := frames.Next()
    // Format: frame.Function, frame.File, frame.Line
    if !more {
        break
    }
}
```

## Design Decisions

### Error Wrapping Strategy

Use `fmt.Errorf()` with `%w` verb for automatic wrapping:

```go
return fmt.Errorf("failed to process: %w", err)
```

This preserves the error chain while adding context.

### When to Create Custom Types

Create custom error types when:
- Errors need structured data (fields, codes, metadata)
- Callers need to inspect error details programmatically
- Error requires special formatting or handling

Use sentinel errors when:
- Simple identity check is sufficient
- No additional context is needed
- Error is package-scoped and well-known

### Performance Considerations

**Stack trace capture** is expensive:
- Only use for critical errors or debug builds
- Consider making it optional via build tags
- Cache formatted stack traces

**Error allocation**:
- Sentinel errors allocate once
- Custom error types allocate per occurrence
- Consider sync.Pool for high-frequency errors

## Production Patterns

### Error Boundaries

Define error handling boundaries in your application:

```go
// Internal errors with stack traces
func internalOperation() error {
    return NewErrorWithStack("critical failure", err)
}

// Public API converts to simpler errors
func PublicAPI() error {
    if err := internalOperation(); err != nil {
        // Strip internal details, return clean error
        return fmt.Errorf("operation failed: %w", ErrInternal)
    }
    return nil
}
```

### Error Codes

For APIs, consider adding error codes:

```go
type APIError struct {
    Code    string
    Message string
    Details map[string]interface{}
}
```

### Logging Integration

Errors with context integrate well with structured logging:

```go
var fileErr *FileError
if errors.As(err, &fileErr) {
    log.Error("file operation failed",
        "op", fileErr.Op,
        "path", fileErr.Path,
        "timestamp", fileErr.Timestamp)
}
```

## Common Pitfalls

1. **Breaking error chains**: Don't use `err.Error()` in `fmt.Errorf()` - use `%w`
2. **Comparing errors with `==`**: Use `errors.Is()` instead
3. **Type assertions on errors**: Use `errors.As()` instead
4. **Too many error types**: Start simple, add complexity as needed
5. **Exposing internal errors**: Wrap internal errors at API boundaries

## Testing Strategies

- Test error wrapping and unwrapping
- Verify sentinel error identity
- Validate custom error type extraction
- Check error message formatting
- Test error chain traversal

## Further Reading

- Go Blog: [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- Effective Go: Error handling
- Package `errors` documentation
- Dave Cheney: [Don't just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
