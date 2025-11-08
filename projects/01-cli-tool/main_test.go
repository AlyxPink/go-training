package main

import (
	"encoding/json"
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
