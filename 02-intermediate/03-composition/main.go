package main

import (
	"fmt"
	"io"
	"strings"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	// TODO: Return new Logger with prefix
	panic("not implemented")
}

func (l *Logger) Log(message string) {
	// TODO: Print formatted log message with prefix and timestamp
	panic("not implemented")
}

func (l *Logger) Error(message string) {
	// TODO: Print error message with prefix
	panic("not implemented")
}

type Service struct {
	*Logger
	name string
}

func NewService(name string) *Service {
	// TODO: Return new Service embedding Logger with "SERVICE" prefix
	panic("not implemented")
}

func (s *Service) Start() {
	// TODO: Log service starting message
	panic("not implemented")
}

func (s *Service) Stop() {
	// TODO: Log service stopping message
	panic("not implemented")
}

type Person struct {
	FirstName string
	LastName  string
}

func (p Person) FullName() string {
	// TODO: Return full name as "FirstName LastName"
	panic("not implemented")
}

type Employee struct {
	Person
	EmployeeID int
	Department string
}

func NewEmployee(firstName, lastName string, id int, dept string) *Employee {
	// TODO: Return new Employee embedding Person with all fields
	panic("not implemented")
}

type ReadWriteCloser interface {
	io.Reader
	io.Writer
	io.Closer
}

type BufferedReadWriteCloser struct {
	*strings.Reader
	buffer strings.Builder
}

func NewBufferedReadWriteCloser(data string) *BufferedReadWriteCloser {
	// TODO: Return BufferedReadWriteCloser with strings.Reader embedded
	panic("not implemented")
}

func (b *BufferedReadWriteCloser) Write(p []byte) (n int, err error) {
	// TODO: Write to internal buffer
	panic("not implemented")
}

func (b *BufferedReadWriteCloser) Close() error {
	// TODO: Return nil (no-op close)
	panic("not implemented")
}

func main() {
	svc := NewService("API")
	svc.Start()
	svc.Log("Processing request")
	svc.Stop()

	emp := NewEmployee("John", "Doe", 12345, "Engineering")
	fmt.Println(emp.FullName())
	fmt.Println(emp.EmployeeID)
}
