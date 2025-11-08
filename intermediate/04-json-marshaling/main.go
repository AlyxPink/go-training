package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// User represents a user with JSON serialization
type User struct {
	// TODO: Add correct JSON tags to these fields
	// Hint: Use `json:"field_name"` and `json:"field_name,omitempty"` for Age
	ID        int
	Username  string
	Email     string
	Age       int
	CreatedAt time.Time
}

// Timestamp wraps time.Time for custom JSON format
type Timestamp struct {
	time.Time
}

// MarshalJSON implements json.Marshaler
func (t Timestamp) MarshalJSON() ([]byte, error) {
	// TODO: Format as Unix timestamp (seconds since epoch)
	return nil, nil
}

// UnmarshalJSON implements json.Unmarshaler  
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// TODO: Parse Unix timestamp
	return nil
}

// Config represents application configuration
type Config struct {
	// TODO: Add Server field (nested struct with Host, Port)
	// TODO: Add Database field (nested struct with URL, MaxConns)
	// TODO: Add Debug field (bool)
}

// Validate checks if config is valid
func (c *Config) Validate() error {
	// TODO: Check required fields
	return nil
}

// Event with custom JSON representation
type Event struct {
	Type    string
	Payload map[string]interface{}
	Time    time.Time
}

// MarshalJSON creates custom JSON representation
func (e Event) MarshalJSON() ([]byte, error) {
	// TODO: Create custom format with type, data, timestamp
	return nil, nil
}

func main() {
	user := User{} // TODO: Initialize
	data, _ := json.Marshal(user)
	fmt.Println(string(data))
	
	var decoded User
	json.Unmarshal(data, &decoded)
	fmt.Printf("%+v\n", decoded)
}
