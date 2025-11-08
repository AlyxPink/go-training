package main

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	count := 0
	for range Pipeline(numbers) {
		count++
	}
	if count == 0 {
		t.Error("Pipeline returned no results")
	}
}

func TestFanOutFanIn(t *testing.T) {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 10; i++ {
			input <- i
		}
	}()

	workers := FanOut(input, 3)
	if len(workers) != 3 {
		t.Errorf("FanOut created %d workers, want 3", len(workers))
	}

	output := FanIn(workers...)
	count := 0
	for range output {
		count++
	}
	if count != 10 {
		t.Errorf("FanIn produced %d results, want 10", count)
	}
}

func TestWorkerPool(t *testing.T) {
	jobs := make(chan Job, 10)
	for i := 1; i <= 10; i++ {
		jobs <- Job{ID: i, Data: i}
	}
	close(jobs)

	results := WorkerPool(jobs, 3)
	count := 0
	for range results {
		count++
	}
	if count != 10 {
		t.Errorf("WorkerPool processed %d jobs, want 10", count)
	}
}
