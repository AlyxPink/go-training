package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Sentinel errors for common conditions
var (
	ErrNotFound      = errors.New("file not found")
	ErrPermission    = errors.New("permission denied")
	ErrInvalidFormat = errors.New("invalid file format")
)

// ValidationError represents validation failures with field-level details
type ValidationError struct {
	Fields    map[string]string
	Timestamp time.Time
	err       error
}

func NewValidationError(fields map[string]string) *ValidationError {
	return &ValidationError{
		Fields:    fields,
		Timestamp: time.Now(),
	}
}

func (e *ValidationError) Error() string {
	fieldList := make([]string, 0, len(e.Fields))
	for field, msg := range e.Fields {
		fieldList = append(fieldList, fmt.Sprintf("%s: %s", field, msg))
	}
	return fmt.Sprintf("validation failed (%d fields): %s", len(e.Fields), strings.Join(fieldList, ", "))
}

func (e *ValidationError) Unwrap() error {
	return e.err
}

// FileError represents file operation errors with context
type FileError struct {
	Op        string
	Path      string
	Err       error
	Timestamp time.Time
}

func NewFileError(op, path string, err error) *FileError {
	return &FileError{
		Op:        op,
		Path:      path,
		Err:       err,
		Timestamp: time.Now(),
	}
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// ErrorWithStack captures call stack for debugging
type ErrorWithStack struct {
	msg   string
	stack []uintptr
	err   error
}

func NewErrorWithStack(msg string, err error) *ErrorWithStack {
	// Capture stack trace (skip 2 frames: Callers and NewErrorWithStack)
	pc := make([]uintptr, 32)
	n := runtime.Callers(2, pc)
	return &ErrorWithStack{
		msg:   msg,
		stack: pc[:n],
		err:   err,
	}
}

func (e *ErrorWithStack) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *ErrorWithStack) Unwrap() error {
	return e.err
}

func (e *ErrorWithStack) StackTrace() string {
	if len(e.stack) == 0 {
		return ""
	}

	var sb strings.Builder
	frames := runtime.CallersFrames(e.stack)

	for {
		frame, more := frames.Next()
		fmt.Fprintf(&sb, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	return sb.String()
}

// ValidateFile validates file content
func ValidateFile(content string) error {
	// TODO: Create a map to store validation errors
	// TODO: Check if content is empty, add error to map
	// TODO: Check if content starts with "valid", add error to map
	// TODO: If there are errors, return NewValidationError
	// TODO: Return nil if validation passes
	panic("not implemented")
}

// ReadFile simulates file reading with various error conditions
func ReadFile(path string) (string, error) {
	// TODO: Check path for "missing", "forbidden", or "corrupt"
	// TODO: Return appropriate sentinel error wrapped in FileError
	// TODO: For valid paths, return "valid file content", nil
	// TODO: Use NewFileError to wrap errors with context
	panic("not implemented")
}

// ProcessFile demonstrates error handling and propagation
func ProcessFile(path string) error {
	// TODO: Call ReadFile and check for errors
	// TODO: Wrap any ReadFile error with fmt.Errorf using %w
	// TODO: Call ValidateFile on the content
	// TODO: Wrap any ValidateFile error with fmt.Errorf using %w
	// TODO: Return nil on success
	panic("not implemented")
}

func main() {
	testCases := []string{
		"/valid/file.txt",
		"/missing/file.txt",
		"/forbidden/file.txt",
		"/corrupt/file.txt",
	}

	for _, path := range testCases {
		fmt.Printf("Processing: %s\n", path)
		err := ProcessFile(path)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)

			// Demonstrate error inspection
			if errors.Is(err, ErrNotFound) {
				fmt.Println("  -> File not found")
			}
			if errors.Is(err, ErrPermission) {
				fmt.Println("  -> Permission denied")
			}

			// Extract custom error types
			var fileErr *FileError
			if errors.As(err, &fileErr) {
				fmt.Printf("  -> FileError: op=%s, path=%s\n", fileErr.Op, fileErr.Path)
			}

			var validationErr *ValidationError
			if errors.As(err, &validationErr) {
				fmt.Printf("  -> ValidationError: %d fields\n", len(validationErr.Fields))
			}
		} else {
			fmt.Println("  Success!")
		}
		fmt.Println()
	}

	// Demonstrate stack trace
	fmt.Println("Stack trace example:")
	stackErr := NewErrorWithStack("critical operation failed", errors.New("database connection lost"))
	fmt.Printf("Error: %v\n\nStack Trace:\n%s\n", stackErr, stackErr.StackTrace())
}
