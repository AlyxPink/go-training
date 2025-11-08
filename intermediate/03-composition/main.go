package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Logger provides logging functionality
type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	// TODO: Return new Logger with prefix
	return nil
}

func (l *Logger) Log(message string) {
	// TODO: Print "[prefix] timestamp: message"
}

func (l *Logger) Error(message string) {
	// TODO: Print "[prefix] ERROR: message"
}

// Service embeds Logger and adds service functionality
type Service struct {
	// TODO: Embed Logger
	// TODO: Add name field
}

func NewService(name string) *Service {
	// TODO: Create Service with embedded Logger
	return nil
}

func (s *Service) Start() {
	// TODO: Log "Service <name> starting" using embedded Logger
}

func (s *Service) Stop() {
	// TODO: Log "Service <name> stopping"
}

// Person represents a person
type Person struct {
	FirstName string
	LastName  string
}

func (p Person) FullName() string {
	// TODO: Return "FirstName LastName"
	return ""
}

// Employee embeds Person and adds employee data
type Employee struct {
	// TODO: Embed Person
	// TODO: Add EmployeeID and Department fields
}

func NewEmployee(firstName, lastName string, id int, dept string) *Employee {
	// TODO: Create Employee with embedded Person
	return nil
}

// ReadWriteCloser combines multiple interfaces
type ReadWriteCloser interface {
	// TODO: Embed io.Reader, io.Writer, io.Closer
}

// BufferedReadWriteCloser implements ReadWriteCloser
type BufferedReadWriteCloser struct {
	// TODO: Embed types that provide needed interfaces
}

func NewBufferedReadWriteCloser(data string) *BufferedReadWriteCloser {
	// TODO: Create with strings.Reader embedded
	return nil
}

func (b *BufferedReadWriteCloser) Close() error {
	// TODO: Implement Close
	return nil
}

func main() {
	// Example usage - uncomment when implemented
	/*
	svc := NewService("API")
	if svc != nil {
		svc.Start()
		svc.Log("Processing request")
		svc.Stop()
	}

	emp := NewEmployee("John", "Doe", 12345, "Engineering")
	if emp != nil {
		fmt.Println(emp.FullName())  // Promoted from Person
		fmt.Println(emp.EmployeeID)
	}
	*/
	fmt.Println("Implement the TODOs and uncomment main() to see it work!")
}
