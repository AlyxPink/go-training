package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestBuggyCounter tests the counter with race conditions
func TestBuggyCounter(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes in BuggyCounter - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	counter := &BuggyCounter{}
	iterations := 1000
	var wg sync.WaitGroup

	// Launch multiple goroutines incrementing counter
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()

	// With races, this will likely not equal iterations
	// The test should PASS but race detector should catch issues
	value := counter.Value()
	if value > int64(iterations) {
		t.Errorf("Counter value %d exceeds iterations %d (impossible without race)", value, iterations)
	}

	t.Logf("Counter value: %d (expected: %d without races)", value, iterations)
}

// TestBuggyCounterConcurrentReads tests concurrent reads and writes
func TestBuggyCounterConcurrentReads(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	counter := &BuggyCounter{}
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			counter.Increment()
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	// Reader goroutines
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 50; j++ {
				_ = counter.Value() // Race detector should catch this
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	<-done
	t.Logf("Final counter value: %d", counter.Value())
}

// TestBuggyMapWriter tests concurrent map writes
func TestBuggyMapWriter(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	// This test may panic due to concurrent map writes
	// We'll catch the panic to make the test pass
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Caught panic (expected with buggy code): %v", r)
		}
	}()

	m := BuggyMapWriter()

	// Map should have 10 entries if no races
	t.Logf("Map size: %d (expected: 10 without races)", len(m))

	// Verify some keys exist
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		if _, ok := m[key]; ok {
			t.Logf("Found %s", key)
		}
	}
}

// TestBuggyMapWriterRace explicitly tests for race detection
func TestBuggyMapWriterRace(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Caught panic from concurrent map write: %v", r)
		}
	}()

	// Run multiple times to increase chance of race
	for attempt := 0; attempt < 5; attempt++ {
		_ = BuggyMapWriter()
	}
}

// TestBuggySliceAppend tests concurrent slice appends
func TestBuggySliceAppend(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Caught panic (possible with slice growth races): %v", r)
		}
	}()

	s := BuggySliceAppend()

	// Slice should have 10 elements if no races
	if len(s) != 10 {
		t.Logf("Slice length: %d (expected: 10 without races)", len(s))
	}

	t.Logf("Slice: %v", s)

	// Check for duplicates or missing values (signs of races)
	seen := make(map[int]bool)
	for _, v := range s {
		if seen[v] {
			t.Logf("Duplicate value detected: %d", v)
		}
		seen[v] = true
	}

	if len(seen) != len(s) {
		t.Logf("Duplicates found: %d unique values in %d elements", len(seen), len(s))
	}
}

// TestBuggySliceAppendMultiple runs multiple iterations
func TestBuggySliceAppendMultiple(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Caught panic: %v", r)
		}
	}()

	for i := 0; i < 10; i++ {
		s := BuggySliceAppend()
		if len(s) != 10 {
			t.Logf("Iteration %d: length %d (expected 10)", i, len(s))
		}
	}
}

// TestBuggyLoopCapture tests loop variable capture
func TestBuggyLoopCapture(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	// Capture output by redirecting fmt.Println
	// This test just ensures the function runs without panic

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BuggyLoopCapture panicked: %v", r)
		}
	}()

	BuggyLoopCapture()

	// Sleep to let goroutines finish
	time.Sleep(200 * time.Millisecond)

	t.Log("BuggyLoopCapture completed (check output for race effects)")
}

// TestLoopCaptureValues tests what values are actually printed
func TestLoopCaptureValues(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	// Run the function and wait
	done := make(chan bool)

	go func() {
		BuggyLoopCapture()
		time.Sleep(200 * time.Millisecond)
		done <- true
	}()

	<-done
	t.Log("Loop capture test completed - with race, may print same value multiple times")
}

// TestRaceConditionProbability tests likelihood of race detection
func TestRaceConditionProbability(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	// Run multiple iterations to increase race detection probability
	iterations := 100

	t.Run("Counter", func(t *testing.T) {
		incorrect := 0
		for i := 0; i < iterations; i++ {
			counter := &BuggyCounter{}
			var wg sync.WaitGroup

			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					counter.Increment()
				}()
			}

			wg.Wait()
			if counter.Value() != int64(100) {
				incorrect++
			}
		}

		if incorrect > 0 {
			t.Logf("Race detected: %d/%d iterations had incorrect count", incorrect, iterations)
		} else {
			t.Log("No races detected in this run (may need -race flag)")
		}
	})
}

// BenchmarkBuggyCounter benchmarks counter performance
func BenchmarkBuggyCounter(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping buggy code benchmark with race detector")
	}

	counter := &BuggyCounter{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})

	b.Logf("Final counter value: %d (expected: %d)", counter.Value(), b.N)
}

// BenchmarkBuggyMapWriter benchmarks map writer
func BenchmarkBuggyMapWriter(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping buggy code benchmark with race detector")
	}

	defer func() {
		if r := recover(); r != nil {
			b.Logf("Caught panic: %v", r)
		}
	}()

	for i := 0; i < b.N; i++ {
		_ = BuggyMapWriter()
	}
}

// TestConcurrentReadWrite tests mixed read/write operations
func TestConcurrentReadWrite(t *testing.T) {
	t.Fatal("TODO: Implement race condition fixes - currently has intentional race conditions")

	if testing.Short() {
		t.Skip("skipping buggy code test with race detector")
	}

	counter := &BuggyCounter{}
	done := make(chan bool)
	var wg sync.WaitGroup

	// Start writers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}

	// Start readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				value := counter.Value()
				if value < 0 {
					t.Errorf("Reader %d: got negative value %d", id, value)
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		done <- true
	}()

	<-done
	t.Logf("Final counter value: %d (expected: 500 without races)", counter.Value())
}

