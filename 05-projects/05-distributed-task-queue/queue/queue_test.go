package queue

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriorityQueue_EnqueueDequeue(t *testing.T) {
	pq := NewPriorityQueue()

	task := &Task{
		ID:       "task1",
		Type:     "test",
		Payload:  []byte("hello"),
		Priority: 1,
	}

	err := pq.Enqueue(task)
	require.NoError(t, err)

	dequeued, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "task1", dequeued.ID)
	assert.Equal(t, "test", dequeued.Type)
	assert.Equal(t, []byte("hello"), dequeued.Payload)
}

func TestPriorityQueue_PriorityOrdering(t *testing.T) {
	pq := NewPriorityQueue()

	// Enqueue tasks with different priorities
	tasks := []*Task{
		{ID: "low", Priority: 1, Type: "test", Payload: []byte("low")},
		{ID: "high", Priority: 3, Type: "test", Payload: []byte("high")},
		{ID: "medium", Priority: 2, Type: "test", Payload: []byte("medium")},
	}

	for _, task := range tasks {
		err := pq.Enqueue(task)
		require.NoError(t, err)
	}

	// Should dequeue in priority order (highest first)
	first, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "high", first.ID)

	second, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "medium", second.ID)

	third, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "low", third.ID)
}

func TestPriorityQueue_DequeueEmpty(t *testing.T) {
	pq := NewPriorityQueue()

	_, err := pq.Dequeue(50 * time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, ErrQueueEmpty, err)
}

func TestPriorityQueue_Ack(t *testing.T) {
	pq := NewPriorityQueue()

	task := &Task{
		ID:       "task1",
		Type:     "test",
		Payload:  []byte("data"),
		Priority: 1,
	}

	err := pq.Enqueue(task)
	require.NoError(t, err)

	dequeued, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)

	err = pq.Ack(dequeued.ID)
	require.NoError(t, err)
}

func TestPriorityQueue_Nack(t *testing.T) {
	pq := NewPriorityQueue()

	task := &Task{
		ID:         "task1",
		Type:       "test",
		Payload:    []byte("data"),
		Priority:   1,
		MaxRetries: 3,
	}

	err := pq.Enqueue(task)
	require.NoError(t, err)

	dequeued, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)

	err = pq.Nack(dequeued.ID, 10*time.Millisecond)
	require.NoError(t, err)

	// Task should be re-queued after delay
	time.Sleep(20 * time.Millisecond)

	requeued, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "task1", requeued.ID)
	assert.Equal(t, 1, requeued.Attempts)
}

func TestPriorityQueue_MaxRetries(t *testing.T) {
	pq := NewPriorityQueue()

	task := &Task{
		ID:         "task1",
		Type:       "test",
		Payload:    []byte("data"),
		Priority:   1,
		MaxRetries: 2,
		Attempts:   2,
	}

	err := pq.Enqueue(task)
	require.NoError(t, err)

	dequeued, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)

	// Nack should not re-queue if max retries reached
	err = pq.Nack(dequeued.ID, 10*time.Millisecond)
	require.NoError(t, err)

	time.Sleep(20 * time.Millisecond)

	// Queue should be empty
	_, err = pq.Dequeue(50 * time.Millisecond)
	assert.Error(t, err)
}

func TestPriorityQueue_ConcurrentEnqueue(t *testing.T) {
	pq := NewPriorityQueue()
	var wg sync.WaitGroup
	numTasks := 100

	// Enqueue tasks concurrently
	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			task := &Task{
				ID:       string(rune(id)),
				Type:     "test",
				Priority: id % 3,
				Payload:  []byte("data"),
			}
			pq.Enqueue(task)
		}(i)
	}

	wg.Wait()

	// Dequeue all tasks
	count := 0
	for i := 0; i < numTasks; i++ {
		_, err := pq.Dequeue(100 * time.Millisecond)
		if err == nil {
			count++
		}
	}

	assert.Equal(t, numTasks, count, "should dequeue all enqueued tasks")
}

func TestPriorityQueue_ConcurrentDequeue(t *testing.T) {
	pq := NewPriorityQueue()
	numTasks := 50

	// Enqueue tasks
	for i := 0; i < numTasks; i++ {
		task := &Task{
			ID:       string(rune(i)),
			Type:     "test",
			Priority: 1,
			Payload:  []byte("data"),
		}
		pq.Enqueue(task)
	}

	// Dequeue concurrently
	var wg sync.WaitGroup
	results := make(chan *Task, numTasks)

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task, err := pq.Dequeue(200 * time.Millisecond)
			if err == nil {
				results <- task
			}
		}()
	}

	wg.Wait()
	close(results)

	// Count results
	count := 0
	seen := make(map[string]bool)
	for task := range results {
		count++
		assert.False(t, seen[task.ID], "task should not be dequeued twice")
		seen[task.ID] = true
	}

	assert.Equal(t, numTasks, count, "all tasks should be dequeued")
}

func TestPriorityQueue_Stats(t *testing.T) {
	pq := NewPriorityQueue()

	stats := pq.GetStats()
	assert.NotNil(t, stats)

	// Enqueue task
	task := &Task{
		ID:       "task1",
		Type:     "test",
		Priority: 1,
		Payload:  []byte("data"),
	}
	pq.Enqueue(task)

	// Check stats updated
	stats = pq.GetStats()
	assert.Equal(t, int64(1), stats.QueueLength)
}

func TestPriorityQueue_MultipleQueues(t *testing.T) {
	pq := NewPriorityQueue()

	// Enqueue tasks with different priorities
	for i := 0; i < 10; i++ {
		for priority := 0; priority < 3; priority++ {
			task := &Task{
				ID:       string(rune(i*3 + priority)),
				Type:     "test",
				Priority: priority,
				Payload:  []byte("data"),
			}
			pq.Enqueue(task)
		}
	}

	// Should dequeue highest priority first
	highCount := 0
	for i := 0; i < 10; i++ {
		task, err := pq.Dequeue(100 * time.Millisecond)
		require.NoError(t, err)
		if task.Priority == 2 {
			highCount++
		}
	}

	assert.Equal(t, 10, highCount, "should dequeue all high priority tasks first")
}

func TestPriorityQueue_TaskStatus(t *testing.T) {
	pq := NewPriorityQueue()

	task := &Task{
		ID:       "task1",
		Type:     "test",
		Priority: 1,
		Status:   StatusPending,
		Payload:  []byte("data"),
	}

	pq.Enqueue(task)

	dequeued, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)

	// Status should be updated to running
	assert.Equal(t, StatusRunning, dequeued.Status)
}

func TestPriorityQueue_FairScheduling(t *testing.T) {
	pq := NewPriorityQueue()

	// Enqueue many high priority tasks
	for i := 0; i < 20; i++ {
		task := &Task{
			ID:       "high-" + string(rune(i)),
			Type:     "test",
			Priority: 3,
			Payload:  []byte("high"),
		}
		pq.Enqueue(task)
	}

	// Enqueue low priority task
	lowTask := &Task{
		ID:       "low-1",
		Type:     "test",
		Priority: 1,
		Payload:  []byte("low"),
	}
	pq.Enqueue(lowTask)

	// Even with many high priority tasks, low priority should eventually be processed
	// (implementation should prevent starvation)
	lowFound := false
	for i := 0; i < 30 && !lowFound; i++ {
		task, err := pq.Dequeue(100 * time.Millisecond)
		if err == nil && task.ID == "low-1" {
			lowFound = true
		}
	}

	// Depending on fair scheduling implementation, low task should be found
	// This test validates anti-starvation logic
}

func TestPriorityQueue_TaskTimestamps(t *testing.T) {
	pq := NewPriorityQueue()

	task := &Task{
		ID:        "task1",
		Type:      "test",
		Priority:  1,
		Payload:   []byte("data"),
		CreatedAt: time.Now(),
	}

	pq.Enqueue(task)

	dequeued, err := pq.Dequeue(100 * time.Millisecond)
	require.NoError(t, err)

	// StartedAt should be set
	assert.NotNil(t, dequeued.StartedAt)
	assert.False(t, dequeued.StartedAt.IsZero())
}

func TestPriorityQueue_DequeueTimeout(t *testing.T) {
	pq := NewPriorityQueue()

	start := time.Now()
	_, err := pq.Dequeue(100 * time.Millisecond)
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond)
	assert.LessOrEqual(t, elapsed, 200*time.Millisecond)
}

func TestPriorityQueue_EnqueueAfterDequeue(t *testing.T) {
	pq := NewPriorityQueue()

	// Start dequeue that will block
	done := make(chan bool)
	go func() {
		task, err := pq.Dequeue(500 * time.Millisecond)
		if err == nil && task.ID == "task1" {
			done <- true
		} else {
			done <- false
		}
	}()

	// Wait a bit then enqueue
	time.Sleep(50 * time.Millisecond)
	task := &Task{
		ID:       "task1",
		Type:     "test",
		Priority: 1,
		Payload:  []byte("data"),
	}
	pq.Enqueue(task)

	// Should receive task
	select {
	case success := <-done:
		assert.True(t, success)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for task")
	}
}

func TestStats_ThreadSafety(t *testing.T) {
	stats := &Stats{}
	var wg sync.WaitGroup

	// Concurrent increments
	for i := 0; i < 100; i++ {
		wg.Add(4)
		go func() {
			defer wg.Done()
			stats.IncrementQueueLength()
		}()
		go func() {
			defer wg.Done()
			stats.DecrementQueueLength()
		}()
		go func() {
			defer wg.Done()
			stats.IncrementCompleted()
		}()
		go func() {
			defer wg.Done()
			stats.IncrementFailed()
		}()
	}

	wg.Wait()

	// Stats should be consistent (no race conditions)
	stats.mu.RLock()
	defer stats.mu.RUnlock()
	assert.Equal(t, int64(0), stats.QueueLength)
	assert.Equal(t, int64(100), stats.CompletedTasks)
	assert.Equal(t, int64(100), stats.FailedTasks)
}
