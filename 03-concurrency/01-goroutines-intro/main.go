package main

import (
	"fmt"
	// TODO: Uncomment for WaitGroup
	// "sync"
	"time"
)

// Task represents a unit of work to be executed concurrently
type Task struct {
	ID       int
	Duration time.Duration
}

// TaskResult represents the outcome of a completed task
type TaskResult struct {
	ID        int
	Message   string
	Completed time.Time
}

// RunTasks executes all tasks concurrently and returns results
// TODO: Implement concurrent task execution using goroutines
func RunTasks(tasks []Task) []TaskResult {
	// TODO: Create a slice to store results
	// TODO: Create a WaitGroup for synchronization
	// TODO: Create a mutex to protect the results slice (or use a channel)

	// TODO: Launch a goroutine for each task
	// Remember to:
	// - Call wg.Add(1) BEFORE launching goroutine
	// - Pass task by value or use proper closure capture
	// - Call defer wg.Done() at start of goroutine
	// - Protect shared state when appending results

	// TODO: Wait for all goroutines to complete

	// TODO: Return collected results
	return nil
}

// processTask simulates task execution
// TODO: Implement task processing logic
func processTask(task Task) TaskResult {
	// TODO: Simulate work by sleeping for task.Duration
	// TODO: Create and return a TaskResult with appropriate message
	return TaskResult{}
}

// RunTasksWithChannel executes tasks concurrently using channels for result collection
// TODO: Implement an alternative version using channels instead of mutex
func RunTasksWithChannel(tasks []Task) []TaskResult {
	// TODO: Create a buffered channel for results
	// TODO: Create a WaitGroup

	// TODO: Launch goroutines that send results to channel

	// TODO: Close channel after all goroutines complete
	// Hint: Use a separate goroutine to close channel after wg.Wait()

	// TODO: Collect results from channel

	return nil
}

func main() {
	// Example tasks
	tasks := []Task{
		{ID: 1, Duration: 100 * time.Millisecond},
		{ID: 2, Duration: 50 * time.Millisecond},
		{ID: 3, Duration: 75 * time.Millisecond},
		{ID: 4, Duration: 200 * time.Millisecond},
		{ID: 5, Duration: 150 * time.Millisecond},
	}

	fmt.Println("Running tasks concurrently with mutex...")
	start := time.Now()
	results := RunTasks(tasks)
	elapsed := time.Since(start)

	fmt.Printf("Completed %d tasks in %v\n", len(results), elapsed)
	for _, result := range results {
		fmt.Printf("  %s\n", result.Message)
	}

	fmt.Println("\nRunning tasks concurrently with channel...")
	start = time.Now()
	results = RunTasksWithChannel(tasks)
	elapsed = time.Since(start)

	fmt.Printf("Completed %d tasks in %v\n", len(results), elapsed)
	for _, result := range results {
		fmt.Printf("  %s\n", result.Message)
	}
}
