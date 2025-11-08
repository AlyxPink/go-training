package main

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

// TestMapReduceBasic tests basic map-reduce functionality
func TestMapReduceBasic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5}

	// Sum of squares: 1 + 4 + 9 + 16 + 25 = 55
	result := MapReduce(
		data,
		func(x int) int { return x * x },
		func(a, b int) int { return a + b },
		2,
	)

	expected := 55
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// TestMapReduceSum tests sum reduction
func TestMapReduceSum(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Identity map, sum reduce
	result := MapReduce(
		data,
		func(x int) int { return x },
		func(a, b int) int { return a + b },
		4,
	)

	expected := 55 // 1+2+...+10
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// TestMapReduceProduct tests product reduction
func TestMapReduceProduct(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5}

	// Product: 1 * 2 * 3 * 4 * 5 = 120
	result := MapReduce(
		data,
		func(x int) int { return x },
		func(a, b int) int { return a * b },
		2,
	)

	expected := 120
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// TestMapReduceMax tests max reduction
func TestMapReduceMax(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{3, 7, 2, 9, 1, 5}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	result := MapReduce(
		data,
		func(x int) int { return x },
		max,
		3,
	)

	expected := 9
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// TestMapReduceEmpty tests empty input
func TestMapReduceEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r == "not implemented" {
				t.Skip("Function not implemented yet")
			}
			// Empty input might panic - that's acceptable
			t.Log("Empty input caused panic (acceptable)")
		}
	}()

	data := []int{}

	result := MapReduce(
		data,
		func(x int) int { return x },
		func(a, b int) int { return a + b },
		2,
	)

	t.Logf("Result for empty input: %d", result)
}

// TestMapReduceOneWorker tests with single worker
func TestMapReduceOneWorker(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5}

	result := MapReduce(
		data,
		func(x int) int { return x * 2 },
		func(a, b int) int { return a + b },
		1, // single worker
	)

	expected := 30 // (1+2+3+4+5)*2
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// TestMapReduceManyWorkers tests with many workers
func TestMapReduceManyWorkers(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := MapReduce(
		data,
		func(x int) int { return x },
		func(a, b int) int { return a + b },
		20, // more workers than data
	)

	expected := 55
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// TestParallelFilterBasic tests basic filtering
func TestParallelFilterBasic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter even numbers
	result := ParallelFilter(
		data,
		func(x int) bool { return x%2 == 0 },
		4,
	)

	expected := []int{2, 4, 6, 8, 10}

	// Sort both for comparison (order might vary)
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestParallelFilterOdd tests filtering odd numbers
func TestParallelFilterOdd(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := ParallelFilter(
		data,
		func(x int) bool { return x%2 == 1 },
		4,
	)

	expected := []int{1, 3, 5, 7, 9}
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestParallelFilterGreaterThan tests range filtering
func TestParallelFilterGreaterThan(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 5, 10, 15, 20, 25, 30}

	result := ParallelFilter(
		data,
		func(x int) bool { return x > 10 },
		3,
	)

	expected := []int{15, 20, 25, 30}
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestParallelFilterNone tests filtering with no matches
func TestParallelFilterNone(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 3, 5, 7, 9}

	result := ParallelFilter(
		data,
		func(x int) bool { return x%2 == 0 },
		2,
	)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v", result)
	}
}

// TestParallelFilterAll tests filtering with all matches
func TestParallelFilterAll(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{2, 4, 6, 8, 10}

	result := ParallelFilter(
		data,
		func(x int) bool { return x%2 == 0 },
		3,
	)

	sort.Ints(result)

	if !reflect.DeepEqual(result, data) {
		t.Errorf("Expected %v, got %v", data, result)
	}
}

// TestParallelFilterEmpty tests empty input
func TestParallelFilterEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{}

	result := ParallelFilter(
		data,
		func(x int) bool { return true },
		2,
	)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v", result)
	}
}

// TestParallelFilterOneWorker tests with single worker
func TestParallelFilterOneWorker(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := []int{1, 2, 3, 4, 5, 6}

	result := ParallelFilter(
		data,
		func(x int) bool { return x > 3 },
		1,
	)

	expected := []int{4, 5, 6}
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestParallelProcessingConcurrency tests concurrent execution
func TestParallelProcessingConcurrency(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}

	var mu sync.Mutex
	executionCount := 0

	// Track concurrent execution
	result := ParallelFilter(
		data,
		func(x int) bool {
			mu.Lock()
			executionCount++
			mu.Unlock()
			return x%2 == 0
		},
		10,
	)

	if len(result) != 50 {
		t.Errorf("Expected 50 even numbers, got %d", len(result))
	}

	if executionCount != 100 {
		t.Errorf("Expected 100 executions, got %d", executionCount)
	}
}

// TestMapReduceLargeData tests with larger dataset
func TestMapReduceLargeData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large data test in short mode")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	data := make([]int, 10000)
	for i := range data {
		data[i] = i + 1
	}

	result := MapReduce(
		data,
		func(x int) int { return x },
		func(a, b int) int { return a + b },
		8,
	)

	// Sum of 1 to 10000 = 10000 * 10001 / 2 = 50005000
	expected := 50005000
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// BenchmarkMapReduce benchmarks map-reduce performance
func BenchmarkMapReduce(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MapReduce(
			data,
			func(x int) int { return x * x },
			func(a, b int) int { return a + b },
			4,
		)
	}
}

// BenchmarkParallelFilter benchmarks parallel filter performance
func BenchmarkParallelFilter(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParallelFilter(
			data,
			func(x int) bool { return x%2 == 0 },
			4,
		)
	}
}

// BenchmarkWorkerScaling benchmarks performance with different worker counts
func BenchmarkWorkerScaling(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	workers := []int{1, 2, 4, 8, 16}

	for _, w := range workers {
		b.Run(string(rune(w)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				MapReduce(
					data,
					func(x int) int { return x * x },
					func(a, b int) int { return a + b },
					w,
				)
			}
		})
	}
}
