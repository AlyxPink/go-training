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
	fields := make(map[string]string)

	if content == "" {
		fields["content"] = "must not be empty"
	}

	if !strings.HasPrefix(content, "valid") {
		fields["format"] = "must start with 'valid'"
	}

	if len(fields) > 0 {
		return NewValidationError(fields)
	}

	return nil
}

// ReadFile simulates file reading with various error conditions
func ReadFile(path string) (string, error) {
	var baseErr error

	// Simulate different error conditions based on path
	if strings.Contains(path, "missing") {
		baseErr = ErrNotFound
	} else if strings.Contains(path, "forbidden") {
		baseErr = ErrPermission
	} else if strings.Contains(path, "corrupt") {
		baseErr = ErrInvalidFormat
	} else {
		// Return valid content
		return "valid file content", nil
	}

	// Wrap with FileError for context
	return "", NewFileError("read", path, baseErr)
}

// ProcessFile demonstrates error handling and propagation
func ProcessFile(path string) error {
	// Read file
	content, err := ReadFile(path)
	if err != nil {
		// Wrap with additional context
		return fmt.Errorf("failed to process file: %w", err)
	}

	// Validate content
	if err := ValidateFile(content); err != nil {
		// Wrap validation error with context
		return fmt.Errorf("file %s validation failed: %w", path, err)
	}

	return nil
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
