package main

import (
	"testing"
	"time"
)

func TestRunTasks(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  int
	}{
		{
			name:  "empty tasks",
			tasks: []Task{},
			want:  0,
		},
		{
			name: "single task",
			tasks: []Task{
				{ID: 1, Duration: 10 * time.Millisecond},
			},
			want: 1,
		},
		{
			name: "multiple tasks",
			tasks: []Task{
				{ID: 1, Duration: 20 * time.Millisecond},
				{ID: 2, Duration: 10 * time.Millisecond},
				{ID: 3, Duration: 15 * time.Millisecond},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := RunTasks(tt.tasks)
			if len(results) != tt.want {
				t.Errorf("RunTasks() returned %d results, want %d", len(results), tt.want)
			}

			// Verify all task IDs are present
			seen := make(map[int]bool)
			for _, result := range results {
				if seen[result.ID] {
					t.Errorf("Duplicate task ID %d in results", result.ID)
				}
				seen[result.ID] = true
			}
		})
	}
}

func TestRunTasksConcurrency(t *testing.T) {
	// Create tasks that each take 100ms
	tasks := make([]Task, 10)
	for i := range tasks {
		tasks[i] = Task{ID: i + 1, Duration: 100 * time.Millisecond}
	}

	start := time.Now()
	results := RunTasks(tasks)
	elapsed := time.Since(start)

	if len(results) != 10 {
		t.Errorf("RunTasks() returned %d results, want 10", len(results))
	}

	// If truly concurrent, should complete in ~100ms, not 1000ms
	// Allow some overhead, but should be < 300ms
	if elapsed > 300*time.Millisecond {
		t.Errorf("RunTasks() took %v, expected concurrent execution in ~100ms", elapsed)
	}
}

func TestRunTasksWithChannel(t *testing.T) {
	tasks := []Task{
		{ID: 1, Duration: 20 * time.Millisecond},
		{ID: 2, Duration: 10 * time.Millisecond},
		{ID: 3, Duration: 15 * time.Millisecond},
	}

	results := RunTasksWithChannel(tasks)
	if len(results) != 3 {
		t.Errorf("RunTasksWithChannel() returned %d results, want 3", len(results))
	}

	// Verify all task IDs are present
	seen := make(map[int]bool)
	for _, result := range results {
		if seen[result.ID] {
			t.Errorf("Duplicate task ID %d in results", result.ID)
		}
		seen[result.ID] = true
	}
}

func TestRunTasksRaceCondition(t *testing.T) {
	// Run with: go test -race
	// This test will fail if there are race conditions
	tasks := make([]Task, 100)
	for i := range tasks {
		tasks[i] = Task{ID: i + 1, Duration: 1 * time.Millisecond}
	}

	results := RunTasks(tasks)
	if len(results) != 100 {
		t.Errorf("RunTasks() returned %d results, want 100", len(results))
	}
}

func BenchmarkRunTasks(b *testing.B) {
	tasks := []Task{
		{ID: 1, Duration: 1 * time.Millisecond},
		{ID: 2, Duration: 1 * time.Millisecond},
		{ID: 3, Duration: 1 * time.Millisecond},
		{ID: 4, Duration: 1 * time.Millisecond},
		{ID: 5, Duration: 1 * time.Millisecond},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunTasks(tasks)
	}
}

func BenchmarkRunTasksWithChannel(b *testing.B) {
	tasks := []Task{
		{ID: 1, Duration: 1 * time.Millisecond},
		{ID: 2, Duration: 1 * time.Millisecond},
		{ID: 3, Duration: 1 * time.Millisecond},
		{ID: 4, Duration: 1 * time.Millisecond},
		{ID: 5, Duration: 1 * time.Millisecond},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunTasksWithChannel(tasks)
	}
}