package main

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"invalid.email", false},
		{"@example.com", false},
		{"user@example", false},
	}
	
	for _, tt := range tests {
		if got := ValidateEmail(tt.email); got != tt.valid {
			t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.valid)
		}
	}
}

func TestExtractPhoneNumbers(t *testing.T) {
	text := "Call 555-123-4567 or 555-987-6543 for help"
	phones := ExtractPhoneNumbers(text)
	
	if len(phones) != 2 {
		t.Errorf("Got %d phone numbers, want 2", len(phones))
	}
}

func TestParseLogLine(t *testing.T) {
	line := "[2024-01-01 12:00:00] INFO: Server started"
	ts, level, msg, ok := ParseLogLine(line)
	
	if !ok {
		t.Fatal("ParseLogLine failed")
	}
	
	if ts != "2024-01-01 12:00:00" {
		t.Errorf("timestamp = %q", ts)
	}
	if level != "INFO" {
		t.Errorf("level = %q", level)
	}
	if msg != "Server started" {
		t.Errorf("message = %q", msg)
	}
}
