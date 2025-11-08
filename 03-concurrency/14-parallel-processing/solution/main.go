package main

import (
	"fmt"
	"sync"
)

func MapReduce(data []int, mapFunc func(int) int, reduceFunc func(int, int) int, workers int) int {
	// Create channels
	jobs := make(chan int, len(data))
	results := make(chan int, len(data))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range jobs {
				results <- mapFunc(num)
			}
		}()
	}

	// Send jobs
	for _, num := range data {
		jobs <- num
	}
	close(jobs)

	// Close results when done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Reduce results
	result := 0
	first := true
	for r := range results {
		if first {
			result = r
			first = false
		} else {
			result = reduceFunc(result, r)
		}
	}

	return result
}

func ParallelFilter(data []int, predicate func(int) bool, workers int) []int {
	type result struct {
		index int
		value int
		ok    bool
	}

	jobs := make(chan struct {
		index int
		value int
	}, len(data))

	results := make(chan result, len(data))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				ok := predicate(job.value)
				results <- result{job.index, job.value, ok}
			}
		}()
	}

	// Send jobs
	for i, v := range data {
		jobs <- struct {
			index int
			value int
		}{i, v}
	}
	close(jobs)

	// Close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect filtered results
	filtered := make([]int, 0)
	for r := range results {
		if r.ok {
			filtered = append(filtered, r.value)
		}
	}

	return filtered
}

func main() {
	fmt.Println("Parallel Processing Examples")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := MapReduce(
		data,
		func(x int) int { return x * x },
		func(a, b int) int { return a + b },
		4,
	)

	fmt.Printf("Sum of squares: %d\n", result)

	evens := ParallelFilter(
		data,
		func(x int) bool { return x%2 == 0 },
		4,
	)

	fmt.Printf("Even numbers: %v\n", evens)
}
