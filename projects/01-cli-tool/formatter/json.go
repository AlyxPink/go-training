package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Formatter interface {
	Format(data interface{}) (string, error)
}

type JSONFormatter struct {
	Compact bool
}

func (f *JSONFormatter) Format(data interface{}) (string, error) {
	// TODO: Implement JSON formatting
	// Hint: Use json.MarshalIndent for pretty, json.Marshal for compact
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	
	if !f.Compact {
		enc.SetIndent("", "  ")
	}
	
	if err := enc.Encode(data); err != nil {
		return "", fmt.Errorf("encoding JSON: %w", err)
	}
	
	return b.String(), nil
}

type RawFormatter struct{}

func (f *RawFormatter) Format(data interface{}) (string, error) {
	// TODO: Implement raw formatting (no quotes for strings)
	// Hint: If string, return as-is; otherwise JSON encode
	if s, ok := data.(string); ok {
		return s + "\n", nil
	}
	
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	
	return string(b) + "\n", nil
}
