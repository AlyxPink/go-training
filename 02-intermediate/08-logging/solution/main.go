package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	logger *log.Logger
	level  int
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags),
		level:  0,
	}
}

func (l *Logger) Info(msg string) {
	l.logger.Println("INFO:", msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Println("ERROR:", msg)
}

func StructuredLog(msg string, attrs map[string]any) {
	args := []any{msg}
	for k, v := range attrs {
		args = append(args, k, v)
	}
	slog.Info(args[0].(string), args[1:]...)
}

func main() {
	logger := NewLogger("APP")
	logger.Info("Application started")
	logger.Error("An error occurred")
	
	StructuredLog("user login", map[string]any{
		"user_id": 123,
		"ip": "192.168.1.1",
	})
}
