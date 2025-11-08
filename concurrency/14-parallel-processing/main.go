package main

import "fmt"

// MapReduce processes data in parallel
// TODO: Implement map-reduce pattern
func MapReduce(data []int, mapFunc func(int) int, reduceFunc func(int, int) int, workers int) int {
	// TODO: Create channels for distribution
	// TODO: Start worker pool
	// TODO: Distribute data (map phase)
	// TODO: Collect and reduce results
	panic("not implemented")
}

// ParallelFilter filters data using multiple workers
// TODO: Implement parallel filter
func ParallelFilter(data []int, predicate func(int) bool, workers int) []int {
	// TODO: Distribute work to workers
	// TODO: Collect filtered results
	// TODO: Maintain order if needed
	panic("not implemented")
}

func main() {
	fmt.Println("Parallel Processing Examples")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Map-reduce: sum of squares
	result := MapReduce(
		data,
		func(x int) int { return x * x },  // Map: square
		func(a, b int) int { return a + b }, // Reduce: sum
		4, // workers
	)

	fmt.Printf("Sum of squares: %d\n", result)

	// Parallel filter
	evens := ParallelFilter(
		data,
		func(x int) bool { return x%2 == 0 },
		4,
	)

	fmt.Printf("Even numbers: %v\n", evens)
}
