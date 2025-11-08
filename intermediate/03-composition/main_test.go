package main

import (
	"io"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger("TEST")
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
	
	// Just test it doesn't panic
	logger.Log("test message")
	logger.Error("error message")
}

func TestService(t *testing.T) {
	svc := NewService("TestService")
	if svc == nil {
		t.Fatal("NewService returned nil")
	}
	
	// Test that Service has embedded Logger methods
	svc.Start()
	svc.Log("test")  // Should work via embedding
	svc.Stop()
}

func TestPerson(t *testing.T) {
	p := Person{FirstName: "John", LastName: "Doe"}
	want := "John Doe"
	if got := p.FullName(); got != want {
		t.Errorf("FullName() = %q, want %q", got, want)
	}
}

func TestEmployee(t *testing.T) {
	emp := NewEmployee("Jane", "Smith", 12345, "Engineering")
	if emp == nil {
		t.Fatal("NewEmployee returned nil")
	}
	
	// Test promoted method from Person
	want := "Jane Smith"
	if got := emp.FullName(); got != want {
		t.Errorf("FullName() = %q, want %q", got, want)
	}
	
	// Test Employee fields
	if emp.EmployeeID != 12345 {
		t.Errorf("EmployeeID = %d, want 12345", emp.EmployeeID)
	}
	if emp.Department != "Engineering" {
		t.Errorf("Department = %q, want \"Engineering\\", emp.Department)
	}
}

func TestBufferedReadWriteCloser(t *testing.T) {
	rwc := NewBufferedReadWriteCloser("test data")
	if rwc == nil {
		t.Fatal("NewBufferedReadWriteCloser returned nil")
	}
	
	// Test Read (from embedded Reader)
	buf := make([]byte, 4)
	n, err := rwc.Read(buf)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if n != 4 {
		t.Errorf("Read() = %d bytes, want 4", n)
	}
	if got := string(buf); got != "test" {
		t.Errorf("Read() = %q, want \"test\"", got)
	}
	
	// Test Close
	if err := rwc.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// Verify interface satisfaction
var _ ReadWriteCloser = (*BufferedReadWriteCloser)(nil)
