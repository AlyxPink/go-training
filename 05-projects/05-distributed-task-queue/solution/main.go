package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alyxpink/go-training/taskqueue/queue"
	"github.com/alyxpink/go-training/taskqueue/worker"
)

var (
	workers = flag.Int("workers", 5, "Number of workers")
	mode    = flag.String("mode", "worker", "Mode: worker or producer")
)

func main() {
	flag.Parse()

	// Create queue
	q := queue.NewPriorityQueue()

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	if *mode == "worker" {
		// Start worker pool
		pool := worker.NewWorkerPool(q, *workers)

		// Register handlers
		pool.RegisterHandler("process", processTaskHandler)
		pool.RegisterHandler("email", emailTaskHandler)

		pool.Start(ctx)
		log.Println("Worker pool started")

		<-ctx.Done()
		pool.Stop()
	} else {
		// Producer mode - enqueue sample tasks
		log.Println("Producer mode - enqueueing sample tasks")

		// Enqueue some sample tasks
		for i := 0; i < 10; i++ {
			task := &queue.Task{
				ID:       "task-" + string(rune(i)),
				Type:     "process",
				Priority: i % 5,
				Payload:  []byte("sample data"),
			}
			if err := q.Enqueue(task); err != nil {
				log.Printf("Failed to enqueue task: %v", err)
			}
		}

		log.Println("Enqueued 10 sample tasks")
		stats := q.GetStats()
		log.Printf("Queue stats: Length=%d, Completed=%d, Failed=%d",
			stats.QueueLength, stats.CompletedTasks, stats.FailedTasks)
	}
}

func processTaskHandler(payload []byte) ([]byte, error) {
	// Implement actual task processing
	log.Printf("Processing task: %s", string(payload))
	return []byte("processed"), nil
}

func emailTaskHandler(payload []byte) ([]byte, error) {
	// Implement email sending
	log.Printf("Sending email: %s", string(payload))
	return []byte("sent"), nil
}
