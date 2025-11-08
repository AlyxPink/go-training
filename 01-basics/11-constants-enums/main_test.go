package main

import "testing"

func TestStatusString(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{
			name:     "StatusPending returns Pending",
			status:   StatusPending,
			expected: "Pending",
		},
		{
			name:     "StatusActive returns Active",
			status:   StatusActive,
			expected: "Active",
		},
		{
			name:     "StatusCompleted returns Completed",
			status:   StatusCompleted,
			expected: "Completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Status.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestStatusValues(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected int
	}{
		{
			name:     "StatusPending has value 0",
			status:   StatusPending,
			expected: 0,
		},
		{
			name:     "StatusActive has value 1",
			status:   StatusActive,
			expected: 1,
		},
		{
			name:     "StatusCompleted has value 2",
			status:   StatusCompleted,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.status) != tt.expected {
				t.Errorf("Status value = %d, expected %d", tt.status, tt.expected)
			}
		})
	}
}

func TestStatusStringEquality(t *testing.T) {
	// Test that different status values produce different strings
	if StatusPending.String() == StatusActive.String() {
		t.Error("StatusPending and StatusActive should have different string representations")
	}
	if StatusActive.String() == StatusCompleted.String() {
		t.Error("StatusActive and StatusCompleted should have different string representations")
	}
	if StatusPending.String() == StatusCompleted.String() {
		t.Error("StatusPending and StatusCompleted should have different string representations")
	}
}

func TestStatusStringNotEmpty(t *testing.T) {
	// Ensure all status values return non-empty strings
	statuses := []Status{StatusPending, StatusActive, StatusCompleted}
	for _, status := range statuses {
		if status.String() == "" {
			t.Errorf("Status %d returned empty string", status)
		}
	}
}
