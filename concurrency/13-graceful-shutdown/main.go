package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Server represents a server that can shutdown gracefully
type Server struct {
	workers int
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewServer creates a new server
func NewServer(workers int) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start launches server workers
// TODO: Implement worker launching
func (s *Server) Start() {
	// TODO: Start worker goroutines
	// TODO: Each worker respects context cancellation
}

// Shutdown gracefully stops the server
// TODO: Implement graceful shutdown
func (s *Server) Shutdown(timeout time.Duration) error {
	// TODO: Cancel context
	// TODO: Wait for workers with timeout
	// TODO: Return error if timeout exceeded
	return nil
}

// worker simulates server work
func (s *Server) worker(id int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Printf("Worker %d: shutting down\n", id)
			return
		default:
			// Simulate work
			fmt.Printf("Worker %d: processing...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	fmt.Println("Graceful Shutdown Example")
	fmt.Println("Press Ctrl+C to trigger shutdown")

	server := NewServer(3)
	server.Start()

	// Setup signal handling
	// TODO: Catch SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	<-sigChan
	fmt.Println("\nShutdown signal received")

	// Graceful shutdown with timeout
	if err := server.Shutdown(5 * time.Second); err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Shutdown complete")
}
