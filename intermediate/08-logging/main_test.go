package main

import (
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger("TEST")
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
	
	logger.Info("test")
	logger.Error("error")
}

func TestStructuredLog(t *testing.T) {
	StructuredLog("test", map[string]any{"key": "value"})
}
