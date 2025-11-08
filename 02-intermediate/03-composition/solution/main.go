package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Log(message string) {
	fmt.Printf("[%s] %s: %s\n", l.prefix, time.Now().Format("15:04:05"), message)
}

func (l *Logger) Error(message string) {
	fmt.Printf("[%s] ERROR: %s\n", l.prefix, message)
}

type Service struct {
	*Logger
	name string
}

func NewService(name string) *Service {
	return &Service{
		Logger: NewLogger("SERVICE"),
		name:   name,
	}
}

func (s *Service) Start() {
	s.Log(fmt.Sprintf("Service %s starting", s.name))
}

func (s *Service) Stop() {
	s.Log(fmt.Sprintf("Service %s stopping", s.name))
}

type Person struct {
	FirstName string
	LastName  string
}

func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

type Employee struct {
	Person
	EmployeeID int
	Department string
}

func NewEmployee(firstName, lastName string, id int, dept string) *Employee {
	return &Employee{
		Person:     Person{FirstName: firstName, LastName: lastName},
		EmployeeID: id,
		Department: dept,
	}
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
	return &BufferedReadWriteCloser{
		Reader: strings.NewReader(data),
	}
}

func (b *BufferedReadWriteCloser) Write(p []byte) (n int, err error) {
	return b.buffer.Write(p)
}

func (b *BufferedReadWriteCloser) Close() error {
	return nil
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
