# Exercise 01: Custom Errors

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Create custom error types with context and stack traces
- Implement error wrapping and unwrapping
- Use `errors.Is()` and `errors.As()` for error inspection
- Build error trees for complex error handling
- Design sentinel errors for specific conditions

## Problem Description

Build a file processing system that demonstrates advanced error handling patterns. The system should:

1. Define custom error types for different failure modes (NotFound, Permission, Validation, etc.)
2. Implement error wrapping to preserve error chains
3. Support error inspection using `errors.Is()` and `errors.As()`
4. Include optional stack trace capture for debugging
5. Provide rich error context (operation, file path, timestamp)

## Requirements

- Custom error types must implement the `error` interface
- Errors should be wrapped with context as they propagate up the call stack
- Sentinel errors for common conditions (ErrNotFound, ErrPermission, etc.)
- Error formatting with `%+v` should show detailed context
- Support error equality checks and type assertions

## Example Usage

```go
err := ProcessFile("/path/to/file.txt")
if err != nil {
    // Check for specific error type
    if errors.Is(err, ErrNotFound) {
        // Handle file not found
    }

    // Extract custom error details
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Validation failed: %v\n", validationErr.Fields)
    }
}
```

## Test Coverage

- Error wrapping and unwrapping
- Sentinel error comparison
- Custom error type extraction
- Error formatting
- Stack trace capture (optional)
