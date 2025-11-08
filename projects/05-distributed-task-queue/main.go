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

	// TODO: Create queue
	q := queue.NewPriorityQueue()

	// TODO: Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	if *mode == "worker" {
		// TODO: Start worker pool
		pool := worker.NewWorkerPool(q, *workers)
		
		// Register handlers
		pool.RegisterHandler("process", processTaskHandler)
		pool.RegisterHandler("email", emailTaskHandler)
		
		pool.Start(ctx)
		log.Println("Worker pool started")
		
		<-ctx.Done()
		pool.Stop()
	} else {
		// TODO: Producer mode - enqueue sample tasks
		log.Println("Producer mode not fully implemented")
	}
}

func processTaskHandler(payload []byte) ([]byte, error) {
	// TODO: Implement actual task processing
	log.Printf("Processing task: %s", string(payload))
	return []byte("processed"), nil
}

func emailTaskHandler(payload []byte) ([]byte, error) {
	// TODO: Implement email sending
	log.Printf("Sending email: %s", string(payload))
	return []byte("sent"), nil
}
