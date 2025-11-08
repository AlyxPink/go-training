package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUser_JSON(t *testing.T) {
	user := User{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Age:      30,
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	
	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	
	if decoded.ID != user.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, user.ID)
	}
	if decoded.Username != user.Username {
		t.Errorf("Username = %q, want %q", decoded.Username, user.Username)
	}
}

func TestUser_OmitEmpty(t *testing.T) {
	user := User{ID: 1, Username: "test", Email: "test@example.com"}
	data, _ := json.Marshal(user)
	
	// Age should not be in JSON when zero
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	if _, exists := m["age"]; exists {
		t.Error("age should be omitted when zero")
	}
}

func TestTimestamp_JSON(t *testing.T) {
	ts := Timestamp{Time: time.Unix(1609459200, 0)}
	
	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	
	var decoded Timestamp
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	
	if !decoded.Time.Equal(ts.Time) {
		t.Errorf("Time = %v, want %v", decoded.Time, ts.Time)
	}
}

func TestConfig_JSON(t *testing.T) {
	configJSON := `{
		"server": {"host": "localhost", "port": 8080},
		"database": {"url": "postgres://localhost", "max_connections": 10},
		"debug": true
	}`
	
	var cfg Config
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	
	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate error: %v", err)
	}
}
