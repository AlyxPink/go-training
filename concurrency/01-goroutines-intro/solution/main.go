package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID       int
	Duration time.Duration
}

type TaskResult struct {
	ID        int
	Message   string
	Completed time.Time
}

// RunTasks executes all tasks concurrently using mutex for result collection
func RunTasks(tasks []Task) []TaskResult {
	if len(tasks) == 0 {
		return []TaskResult{}
	}

	var (
		results []TaskResult
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	// Launch goroutine for each task
	for _, task := range tasks {
		wg.Add(1) // Increment BEFORE launching goroutine

		// Capture task value for closure
		task := task

		go func() {
			defer wg.Done() // Ensure Done() called even on panic

			// Process task
			result := processTask(task)

			// Safely append to shared slice
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return results
}

// processTask simulates task execution
func processTask(task Task) TaskResult {
	// Simulate work
	time.Sleep(task.Duration)

	return TaskResult{
		ID:        task.ID,
		Message:   fmt.Sprintf("Task %d completed after %v", task.ID, task.Duration),
		Completed: time.Now(),
	}
}

// RunTasksWithChannel executes tasks concurrently using channels for result collection
func RunTasksWithChannel(tasks []Task) []TaskResult {
	if len(tasks) == 0 {
		return []TaskResult{}
	}

	// Buffered channel to avoid blocking
	resultChan := make(chan TaskResult, len(tasks))
	var wg sync.WaitGroup

	// Launch goroutine for each task
	for _, task := range tasks {
		wg.Add(1)
		task := task // Capture for closure

		go func() {
			defer wg.Done()
			result := processTask(task)
			resultChan <- result // Send to channel
		}()
	}

	// Close channel after all goroutines complete
	// Must be in separate goroutine to avoid deadlock
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from channel
	results := make([]TaskResult, 0, len(tasks))
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func main() {
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