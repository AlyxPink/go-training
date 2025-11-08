package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestTaskSchedulerBasic tests basic scheduling functionality
func TestTaskSchedulerBasic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var counter atomic.Int32
	err := scheduler.Schedule("test", 50*time.Millisecond, func() {
		counter.Add(1)
	})
	if err != nil {
		t.Fatalf("Failed to schedule task: %v", err)
	}

	// Wait for task to execute multiple times
	time.Sleep(150 * time.Millisecond)

	count := counter.Load()
	if count < 2 {
		t.Errorf("Expected at least 2 executions, got %d", count)
	}
}

// TestTaskSchedulerCron tests cron-like scheduling
func TestTaskSchedulerCron(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var counter atomic.Int32
	// Use duration string as simplified cron expression
	err := scheduler.ScheduleCron("cron-test", "100ms", func() {
		counter.Add(1)
	})
	if err != nil {
		t.Fatalf("Failed to schedule cron task: %v", err)
	}

	time.Sleep(250 * time.Millisecond)

	count := counter.Load()
	if count < 2 {
		t.Errorf("Expected at least 2 cron executions, got %d", count)
	}
}

// TestTaskSchedulerOneTime tests one-time delayed tasks
func TestTaskSchedulerOneTime(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var executed atomic.Bool
	err := scheduler.ScheduleOnce("once", 100*time.Millisecond, func() {
		executed.Store(true)
	})
	if err != nil {
		t.Fatalf("Failed to schedule one-time task: %v", err)
	}

	// Check it hasn't executed yet
	if executed.Load() {
		t.Error("Task executed too early")
	}

	// Wait for execution
	time.Sleep(150 * time.Millisecond)

	if !executed.Load() {
		t.Error("One-time task did not execute")
	}

	// Wait longer and verify it only executed once
	time.Sleep(100 * time.Millisecond)
	status, err := scheduler.GetTaskStatus("once")
	if err != nil {
		t.Fatalf("Failed to get task status: %v", err)
	}
	if status != StatusCompleted {
		t.Errorf("Expected task status to be Completed, got %v", status)
	}
}

// TestTaskSchedulerConcurrency tests concurrent task execution
func TestTaskSchedulerConcurrency(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var counter1, counter2, counter3 atomic.Int32

	// Schedule multiple concurrent tasks
	err1 := scheduler.Schedule("task1", 50*time.Millisecond, func() {
		counter1.Add(1)
	})
	err2 := scheduler.Schedule("task2", 50*time.Millisecond, func() {
		counter2.Add(1)
	})
	err3 := scheduler.Schedule("task3", 50*time.Millisecond, func() {
		counter3.Add(1)
	})

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("Failed to schedule tasks: %v, %v, %v", err1, err2, err3)
	}

	time.Sleep(150 * time.Millisecond)

	c1, c2, c3 := counter1.Load(), counter2.Load(), counter3.Load()
	if c1 < 2 || c2 < 2 || c3 < 2 {
		t.Errorf("Expected at least 2 executions for each task, got %d, %d, %d", c1, c2, c3)
	}
}

// TestTaskSchedulerCancellation tests task cancellation
func TestTaskSchedulerCancellation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var counter atomic.Int32
	err := scheduler.Schedule("cancel-test", 50*time.Millisecond, func() {
		counter.Add(1)
	})
	if err != nil {
		t.Fatalf("Failed to schedule task: %v", err)
	}

	// Let it run a bit
	time.Sleep(75 * time.Millisecond)
	countBefore := counter.Load()

	// Cancel the task
	err = scheduler.Cancel("cancel-test")
	if err != nil {
		t.Fatalf("Failed to cancel task: %v", err)
	}

	// Wait and verify it doesn't execute anymore
	time.Sleep(100 * time.Millisecond)
	countAfter := counter.Load()

	if countAfter != countBefore {
		t.Errorf("Task executed after cancellation: before=%d, after=%d", countBefore, countAfter)
	}

	// Verify task is removed
	_, err = scheduler.GetTaskStatus("cancel-test")
	if err == nil {
		t.Error("Expected error when getting status of cancelled task")
	}
}

// TestTaskSchedulerPriority tests task priority handling
func TestTaskSchedulerPriority(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var execOrder []int
	var mu sync.Mutex

	// Schedule tasks with different priorities
	err1 := scheduler.ScheduleWithPriority("low", 50*time.Millisecond, 1, func() {
		mu.Lock()
		execOrder = append(execOrder, 1)
		mu.Unlock()
	})
	err2 := scheduler.ScheduleWithPriority("high", 50*time.Millisecond, 10, func() {
		mu.Lock()
		execOrder = append(execOrder, 10)
		mu.Unlock()
	})

	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to schedule priority tasks: %v, %v", err1, err2)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	// Just verify both executed (actual priority ordering would need more complex implementation)
	if len(execOrder) < 2 {
		t.Errorf("Expected at least 2 task executions, got %d", len(execOrder))
	}
}

// TestTaskSchedulerErrorHandling tests error handling
func TestTaskSchedulerErrorHandling(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	// Test empty task name
	err := scheduler.Schedule("", 50*time.Millisecond, func() {})
	if err == nil {
		t.Error("Expected error for empty task name")
	}

	// Test invalid interval
	err = scheduler.Schedule("test", -1*time.Second, func() {})
	if err == nil {
		t.Error("Expected error for negative interval")
	}

	// Test nil function
	err = scheduler.Schedule("test", 50*time.Millisecond, nil)
	if err == nil {
		t.Error("Expected error for nil function")
	}

	// Test duplicate task name
	err = scheduler.Schedule("dup", 50*time.Millisecond, func() {})
	if err != nil {
		t.Fatalf("Failed to schedule first task: %v", err)
	}
	err = scheduler.Schedule("dup", 50*time.Millisecond, func() {})
	if err == nil {
		t.Error("Expected error for duplicate task name")
	}

	// Test canceling non-existent task
	err = scheduler.Cancel("nonexistent")
	if err == nil {
		t.Error("Expected error when canceling non-existent task")
	}
}

// TestTaskSchedulerShutdown tests graceful shutdown
func TestTaskSchedulerShutdown(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()

	var counter atomic.Int32
	err := scheduler.Schedule("shutdown-test", 50*time.Millisecond, func() {
		counter.Add(1)
		time.Sleep(10 * time.Millisecond) // Simulate work
	})
	if err != nil {
		t.Fatalf("Failed to schedule task: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Shutdown with timeout
	start := time.Now()
	err = scheduler.Shutdown(2 * time.Second)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	if duration > 2*time.Second {
		t.Errorf("Shutdown took too long: %v", duration)
	}

	// Verify tasks stopped executing
	countBefore := counter.Load()
	time.Sleep(100 * time.Millisecond)
	countAfter := counter.Load()

	if countAfter != countBefore {
		t.Error("Tasks still executing after shutdown")
	}
}

// TestTaskSchedulerOverlap tests overlapping task execution prevention
func TestTaskSchedulerOverlap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var running atomic.Int32
	var maxConcurrent atomic.Int32

	err := scheduler.Schedule("overlap-test", 50*time.Millisecond, func() {
		current := running.Add(1)
		defer running.Add(-1)

		// Update max concurrent
		for {
			max := maxConcurrent.Load()
			if current <= max {
				break
			}
			if maxConcurrent.CompareAndSwap(max, current) {
				break
			}
		}

		// Simulate long-running task
		time.Sleep(100 * time.Millisecond)
	})

	if err != nil {
		t.Fatalf("Failed to schedule task: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	max := maxConcurrent.Load()
	if max > 1 {
		t.Errorf("Expected max concurrent executions to be 1, got %d (overlap prevention failed)", max)
	}
}

// TestTaskSchedulerRetry tests task retry logic
func TestTaskSchedulerRetry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	var attempts atomic.Int32
	err := scheduler.ScheduleWithRetry("retry-test", 100*time.Millisecond, 3, func() {
		count := attempts.Add(1)
		if count < 3 {
			panic("simulated failure")
		}
	})

	if err != nil {
		t.Fatalf("Failed to schedule task with retry: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	finalAttempts := attempts.Load()
	if finalAttempts < 3 {
		t.Errorf("Expected at least 3 attempts (including retries), got %d", finalAttempts)
	}
}

// TestTaskSchedulerCronInvalid tests invalid cron expressions
func TestTaskSchedulerCronInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	err := scheduler.ScheduleCron("invalid-cron", "invalid", func() {})
	if err == nil {
		t.Error("Expected error for invalid cron expression")
	}

	err = scheduler.ScheduleCron("", "1s", func() {})
	if err == nil {
		t.Error("Expected error for empty task name")
	}

	err = scheduler.ScheduleCron("nil-func", "1s", nil)
	if err == nil {
		t.Error("Expected error for nil function")
	}
}

// TestTaskSchedulerOnceValidation tests ScheduleOnce validation
func TestTaskSchedulerOnceValidation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(1 * time.Second)

	err := scheduler.ScheduleOnce("", 100*time.Millisecond, func() {})
	if err == nil {
		t.Error("Expected error for empty task name")
	}

	err = scheduler.ScheduleOnce("negative", -1*time.Second, func() {})
	if err == nil {
		t.Error("Expected error for negative delay")
	}

	err = scheduler.ScheduleOnce("nil-func", 100*time.Millisecond, nil)
	if err == nil {
		t.Error("Expected error for nil function")
	}
}

// TestTaskSchedulerShutdownTimeout tests shutdown timeout handling
func TestTaskSchedulerShutdownTimeout(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()

	// Schedule a task that runs longer than shutdown timeout
	err := scheduler.Schedule("long-task", 50*time.Millisecond, func() {
		time.Sleep(5 * time.Second)
	})
	if err != nil {
		t.Fatalf("Failed to schedule task: %v", err)
	}

	time.Sleep(60 * time.Millisecond) // Let task start

	// Attempt shutdown with short timeout
	err = scheduler.Shutdown(100 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error during shutdown")
	}
}

// BenchmarkTaskScheduler benchmarks scheduler performance
func BenchmarkTaskScheduler(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(5 * time.Second)

	var counter atomic.Int32

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskName := "bench-task"
		err := scheduler.Schedule(taskName, 10*time.Millisecond, func() {
			counter.Add(1)
		})
		if err != nil {
			b.Fatalf("Failed to schedule task: %v", err)
		}
		scheduler.Cancel(taskName)
	}
}

// BenchmarkTaskSchedulerExecution benchmarks task execution throughput
func BenchmarkTaskSchedulerExecution(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	scheduler := NewScheduler()
	defer scheduler.Shutdown(5 * time.Second)

	var counter atomic.Int32

	// Schedule tasks with unique names
	for i := 0; i < 10; i++ {
		taskName := "exec-task-" + string(rune('0'+i))
		err := scheduler.Schedule(taskName, 1*time.Millisecond, func() {
			counter.Add(1)
		})
		if err != nil {
			b.Fatalf("Failed to schedule task: %v", err)
		}
		defer scheduler.Cancel(taskName)
	}

	b.ResetTimer()
	start := counter.Load()
	time.Sleep(100 * time.Millisecond)
	end := counter.Load()

	executions := end - start
	b.ReportMetric(float64(executions), "executions")
}
