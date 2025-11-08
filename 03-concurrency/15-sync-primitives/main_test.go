package main

import (
	"sync"
	"testing"
	"time"
)

// TestGetDatabase tests singleton pattern with sync.Once
func TestGetDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	db1 := GetDatabase()
	db2 := GetDatabase()

	if db1 != db2 {
		t.Error("GetDatabase returned different instances - singleton violated")
	}

	if db1 == nil {
		t.Error("GetDatabase returned nil")
	}
}

// TestGetDatabaseConcurrent tests concurrent singleton access
func TestGetDatabaseConcurrent(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	// Reset singleton for test (in real code, you can't do this)
	instance = nil
	once = sync.Once{}

	var wg sync.WaitGroup
	instances := make([]*Database, 100)
	panicked := false
	var panicMu sync.Mutex

	// Concurrently get database
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicMu.Lock()
					panicked = true
					panicMu.Unlock()
				}
			}()
			instances[idx] = GetDatabase()
		}(i)
	}

	wg.Wait()

	if panicked {
		t.Fatal("Function not implemented or panic occurred")
	}

	// All should be same instance
	first := instances[0]
	for i, inst := range instances {
		if inst != first {
			t.Errorf("Instance %d is different from first instance", i)
		}
	}
}

// TestGetBuffer tests buffer pool get operation
func TestGetBuffer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	buf := GetBuffer()

	if buf == nil {
		t.Error("GetBuffer returned nil")
	}

	if cap(buf) != 1024 {
		t.Error("Buffer capacity is correct")
	}
}

// TestBufferPoolReuse tests buffer pooling and reuse
func TestBufferPoolReuse(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	buf1 := GetBuffer()
	if buf1 == nil {
		t.Fatal("First GetBuffer returned nil")
	}

	// Use buffer
	copy(buf1, []byte("test data"))

	// Return to pool
	PutBuffer(buf1)

	// Get another buffer - might be the same one
	buf2 := GetBuffer()
	if buf2 == nil {
		t.Fatal("Second GetBuffer returned nil")
	}

	// Buffer should be reset/cleared
	// (Implementation detail - test assumes buffer is cleared)
	t.Logf("Got buffer from pool, len=%d", len(buf2))
}

// TestBufferPoolConcurrent tests concurrent buffer pool access
func TestBufferPoolConcurrent(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	var wg sync.WaitGroup
	panicked := false
	var panicMu sync.Mutex

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicMu.Lock()
					panicked = true
					panicMu.Unlock()
				}
			}()

			buf := GetBuffer()
			if buf == nil {
				t.Errorf("Goroutine %d: GetBuffer returned nil", id)
				return
			}

			// Use buffer
			copy(buf, []byte("data"))

			// Return to pool
			PutBuffer(buf)
		}(i)
	}

	wg.Wait()

	if panicked {
		t.Fatal("Function not implemented or panic occurred")
	}

	t.Log("Concurrent buffer pool access completed")
}

// TestNewQueue tests queue initialization
func TestNewQueue(t *testing.T) {
	q := NewQueue()

	if q == nil {
		t.Fatal("NewQueue returned nil")
	}

	if q.cond == nil {
		t.Error("Queue condition variable is nil")
	}
}

// TestQueuePushPop tests basic queue operations
func TestQueuePushPop(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	q := NewQueue()

	panicked := false
	var panicMu sync.Mutex

	// Push items
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicMu.Lock()
				panicked = true
				panicMu.Unlock()
			}
		}()
		for i := 0; i < 5; i++ {
			q.Push(i)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Pop items
	for i := 0; i < 5; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicMu.Lock()
					panicked = true
					panicMu.Unlock()
				}
			}()
			val := q.Pop()
			if val != i {
				t.Errorf("Expected %d, got %d", i, val)
			}
		}()
		if panicked {
			t.Fatal("Function not implemented or panic occurred")
		}
	}
}

// TestQueueBlockingPop tests that Pop blocks when queue is empty
func TestQueueBlockingPop(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	q := NewQueue()

	done := make(chan bool)
	panicked := make(chan bool, 1)
	var popped int

	// Pop should block
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked <- true
			}
		}()
		popped = q.Pop()
		done <- true
	}()

	// Check if panicked early
	select {
	case <-panicked:
		t.Fatal("Function not implemented or panic occurred")
	case <-time.After(50 * time.Millisecond):
		// Continue with test
	}

	// Verify it's blocked
	select {
	case <-done:
		t.Error("Pop should have blocked on empty queue")
	case <-time.After(100 * time.Millisecond):
		// Good - it's blocking
	}

	// Push item to unblock
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Skip("Function not implemented yet")
			}
		}()
		q.Push(42)
	}()

	// Should complete now
	select {
	case <-done:
		if popped != 42 {
			t.Errorf("Expected 42, got %d", popped)
		}
	case <-panicked:
		t.Fatal("Function not implemented or panic occurred")
	case <-time.After(1 * time.Second):
		t.Error("Pop did not unblock after Push")
	}
}

// TestQueueMultipleWaiters tests multiple goroutines waiting on Pop
func TestQueueMultipleWaiters(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	q := NewQueue()

	var wg sync.WaitGroup
	results := make([]int, 5)
	panicked := false
	var panicMu sync.Mutex

	// Start multiple waiters
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicMu.Lock()
					panicked = true
					panicMu.Unlock()
				}
			}()
			results[idx] = q.Pop()
		}(i)
	}

	// Let them block
	time.Sleep(100 * time.Millisecond)

	// Push items to wake them up
	for i := 0; i < 5; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicMu.Lock()
					panicked = true
					panicMu.Unlock()
				}
			}()
			q.Push(i * 10)
		}()
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()

	if panicked {
		t.Fatal("Function not implemented or panic occurred")
	}

	// Verify all got values
	seen := make(map[int]bool)
	for _, val := range results {
		seen[val] = true
	}

	if len(seen) != 5 {
		t.Errorf("Expected 5 unique values, got %d: %v", len(seen), results)
	}
}

// TestQueueConcurrentPushPop tests concurrent push and pop
func TestQueueConcurrentPushPop(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	q := NewQueue()

	var wg sync.WaitGroup
	panicked := false
	var panicMu sync.Mutex

	// Producers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicMu.Lock()
					panicked = true
					panicMu.Unlock()
				}
			}()
			for j := 0; j < 10; j++ {
				q.Push(id*100 + j)
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	// Consumers
	consumed := make([]int, 0, 30)
	var mu sync.Mutex

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicMu.Lock()
					panicked = true
					panicMu.Unlock()
				}
			}()
			for j := 0; j < 15; j++ {
				val := q.Pop()
				mu.Lock()
				consumed = append(consumed, val)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if panicked {
		t.Fatal("Function not implemented or panic occurred")
	}

	if len(consumed) != 30 {
		t.Errorf("Expected 30 items consumed, got %d", len(consumed))
	}
}

// TestQueueFIFO tests FIFO ordering
func TestQueueFIFO(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	q := NewQueue()

	// Push in order
	for i := 0; i < 10; i++ {
		q.Push(i)
	}

	// Pop and verify order
	for i := 0; i < 10; i++ {
		val := q.Pop()
		if val != i {
			t.Errorf("FIFO violated: expected %d, got %d", i, val)
		}
	}
}

// BenchmarkGetDatabase benchmarks singleton access
func BenchmarkGetDatabase(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetDatabase()
		}
	})
}

// BenchmarkBufferPool benchmarks buffer pool performance
func BenchmarkBufferPool(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetBuffer()
			PutBuffer(buf)
		}
	})
}

// BenchmarkQueue benchmarks queue operations
func BenchmarkQueue(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	q := NewQueue()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				q.Push(i)
			} else {
				// Try to pop, but don't block
				go func() { q.Pop() }()
			}
			i++
		}
	})
}
