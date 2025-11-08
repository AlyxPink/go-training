package main

import (
	"errors"
	"strings"
	"testing"
)

func TestValidationError(t *testing.T) {
	fields := map[string]string{
		"email": "invalid email format",
		"age":   "must be positive",
	}

	err := NewValidationError(fields)
	if err == nil {
		t.Fatal("NewValidationError returned nil")
	}

	if len(err.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(err.Fields))
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "validation") {
		t.Errorf("Error message should mention validation: %s", errMsg)
	}
}

func TestFileError(t *testing.T) {
	underlying := errors.New("disk full")
	fileErr := NewFileError("write", "/tmp/test.txt", underlying)

	if fileErr == nil {
		t.Fatal("NewFileError returned nil")
	}

	if fileErr.Op != "write" {
		t.Errorf("Expected op 'write', got %s", fileErr.Op)
	}

	if fileErr.Path != "/tmp/test.txt" {
		t.Errorf("Expected path '/tmp/test.txt', got %s", fileErr.Path)
	}

	errMsg := fileErr.Error()
	if !strings.Contains(errMsg, "write") || !strings.Contains(errMsg, "/tmp/test.txt") {
		t.Errorf("Error message should include operation and path: %s", errMsg)
	}

	// Test unwrapping
	unwrapped := errors.Unwrap(fileErr)
	if unwrapped != underlying {
		t.Error("Unwrap should return underlying error")
	}
}

func TestErrorWrapping(t *testing.T) {
	err := ReadFile("/missing/file.txt")
	if err == nil {
		t.Fatal("Expected error for missing file")
	}

	// Check sentinel error
	if !errors.Is(err, ErrNotFound) {
		t.Error("Error should wrap ErrNotFound")
	}

	// Extract FileError
	var fileErr *FileError
	if !errors.As(err, &fileErr) {
		t.Error("Error should contain FileError")
	}

	if fileErr.Op != "read" {
		t.Errorf("Expected operation 'read', got %s", fileErr.Op)
	}
}

func TestPermissionError(t *testing.T) {
	err := ReadFile("/forbidden/file.txt")
	if err == nil {
		t.Fatal("Expected error for forbidden file")
	}

	if !errors.Is(err, ErrPermission) {
		t.Error("Error should wrap ErrPermission")
	}
}

func TestValidateFile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{"empty content", "", true},
		{"valid content", "valid file content", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFile(tt.content)
			if tt.shouldErr && err == nil {
				t.Error("Expected validation error")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.shouldErr {
				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Error("Validation error should be *ValidationError")
				}
			}
		})
	}
}

func TestProcessFile(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectErr    error
		expectErrMsg string
	}{
		{"valid file", "/valid/file.txt", nil, ""},
		{"missing file", "/missing/file.txt", ErrNotFound, "not found"},
		{"forbidden file", "/forbidden/file.txt", ErrPermission, "permission"},
		{"corrupt file", "/corrupt/file.txt", ErrInvalidFormat, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessFile(tt.path)

			if tt.expectErr == nil && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if tt.expectErr != nil {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if !errors.Is(err, tt.expectErr) {
					t.Errorf("Expected error %v, got %v", tt.expectErr, err)
				}
				if !strings.Contains(err.Error(), tt.expectErrMsg) {
					t.Errorf("Error message should contain '%s': %s", tt.expectErrMsg, err.Error())
				}
			}
		})
	}
}

func TestErrorWithStack(t *testing.T) {
	underlying := errors.New("base error")
	stackErr := NewErrorWithStack("operation failed", underlying)

	if stackErr == nil {
		t.Fatal("NewErrorWithStack returned nil")
	}

	// Test unwrapping
	unwrapped := errors.Unwrap(stackErr)
	if unwrapped != underlying {
		t.Error("Unwrap should return underlying error")
	}

	// Test stack trace
	trace := stackErr.StackTrace()
	if trace == "" {
		t.Error("StackTrace should return non-empty string")
	}

	// Stack trace should contain function names
	if !strings.Contains(trace, "TestErrorWithStack") {
		t.Errorf("Stack trace should contain test function name: %s", trace)
	}
}
