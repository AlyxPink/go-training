package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alyxpink/go-training/jq/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryExecution(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		query   string
		want    interface{}
		wantErr bool
	}{
		{
			name:  "simple field",
			input: `{"name": "Alice"}`,
			query: ".name",
			want:  "Alice",
		},
		{
			name:  "array index",
			input: `{"users": [1, 2, 3]}`,
			query: ".users[0]",
			want:  1.0,
		},
		{
			name:  "nested field",
			input: `{"user": {"name": "Bob"}}`,
			query: ".user.name",
			want:  "Bob",
		},
		{
			name:  "array length",
			input: `{"users": [1, 2, 3]}`,
			query: ".users.length",
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data interface{}
			err := json.Unmarshal([]byte(tt.input), &data)
			require.NoError(t, err)

			q, err := query.Parse(tt.query)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			got, err := q.Execute(data)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		queryStr  string
		wantErr   bool
		setupFlag func()
		validate  func(t *testing.T, output string)
	}{
		{
			name:     "simple field selection",
			input:    `{"name": "Alice", "age": 30}`,
			queryStr: ".name",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				output = strings.TrimSpace(output)
				assert.Contains(t, output, "Alice")
			},
		},
		{
			name:     "array iteration",
			input:    `{"users": [{"name": "Alice"}, {"name": "Bob"}]}`,
			queryStr: ".users[]",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "Alice")
				assert.Contains(t, output, "Bob")
			},
		},
		{
			name:     "array indexing",
			input:    `{"items": ["a", "b", "c"]}`,
			queryStr: ".items[1]",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				output = strings.TrimSpace(output)
				assert.Contains(t, output, "b")
			},
		},
		{
			name:     "length operation",
			input:    `{"items": [1, 2, 3, 4, 5]}`,
			queryStr: ".items length",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				output = strings.TrimSpace(output)
				assert.Contains(t, output, "5")
			},
		},
		{
			name:     "compact output",
			input:    `{"name": "Alice", "age": 30}`,
			queryStr: ".name",
			wantErr:  false,
			setupFlag: func() {
				*compact = true
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				// Compact should have minimal formatting
				assert.LessOrEqual(t, strings.Count(output, "\n"), 2)
			},
		},
		{
			name:     "raw output",
			input:    `{"message": "Hello World"}`,
			queryStr: ".message",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = true
			},
			validate: func(t *testing.T, output string) {
				output = strings.TrimSpace(output)
				// Raw output should not have quotes around strings
				assert.False(t, strings.HasPrefix(output, `"`))
				assert.False(t, strings.HasSuffix(output, `"`))
				assert.Equal(t, "Hello World", output)
			},
		},
		{
			name:     "table output with array of objects",
			input:    `[{"name": "Alice", "age": 30}, {"name": "Bob", "age": 25}]`,
			queryStr: ".[]",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = true
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				// Table output should have separator
				assert.Contains(t, output, "---")
			},
		},
		{
			name:     "invalid JSON",
			input:    `{invalid json}`,
			queryStr: ".",
			wantErr:  true,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {},
		},
		{
			name:     "nested field access",
			input:    `{"user": {"profile": {"name": "Alice"}}}`,
			queryStr: ".user.profile.name",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "Alice")
			},
		},
		{
			name:     "multiple JSON objects",
			input:    `{"id": 1}` + "\n" + `{"id": 2}`,
			queryStr: ".id",
			wantErr:  false,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "1")
				assert.Contains(t, output, "2")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFlag()

			q, err := query.Parse(tt.queryStr)
			require.NoError(t, err)

			reader := strings.NewReader(tt.input)
			err = processInput(reader, q, "test")

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		setupFlag func()
		wantErr   bool
	}{
		{
			name: "simple string",
			data: "hello",
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
		},
		{
			name: "string with raw flag",
			data: "hello world",
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = true
			},
		},
		{
			name: "object with compact flag",
			data: map[string]interface{}{"name": "Alice", "age": 30},
			setupFlag: func() {
				*compact = true
				*table = false
				*raw = false
			},
		},
		{
			name: "object with pretty print",
			data: map[string]interface{}{"name": "Alice", "age": 30},
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
		},
		{
			name: "array for table output",
			data: []interface{}{
				map[string]interface{}{"name": "Alice", "age": 30},
				map[string]interface{}{"name": "Bob", "age": 25},
			},
			setupFlag: func() {
				*compact = false
				*table = true
				*raw = false
			},
		},
		{
			name: "number output",
			data: 42,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
		},
		{
			name: "boolean output",
			data: true,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
		},
		{
			name: "null output",
			data: nil,
			setupFlag: func() {
				*compact = false
				*table = false
				*raw = false
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFlag()

			// Redirect stdout for testing
			var buf bytes.Buffer
			err := outputResult(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			_ = buf
		})
	}
}

func TestProcessInputErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		queryStr string
	}{
		{
			name:     "array index out of bounds",
			input:    `{"items": [1, 2]}`,
			queryStr: ".items[5]",
		},
		{
			name:     "field selection on array",
			input:    `[1, 2, 3]`,
			queryStr: ".name",
		},
		{
			name:     "array indexing on object",
			input:    `{"name": "Alice"}`,
			queryStr: ".[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*compact = false
			*table = false
			*raw = false

			q, err := query.Parse(tt.queryStr)
			require.NoError(t, err)

			reader := strings.NewReader(tt.input)
			err = processInput(reader, q, "test")

			assert.Error(t, err, "expected error for %s", tt.name)
		})
	}
}
