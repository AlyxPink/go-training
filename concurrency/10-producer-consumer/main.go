package main

import (
	"fmt"
	"sync"
	"time"
)

// Item represents work to be processed
type Item struct {
	ID   int
	Data string
}

// ProducerConsumer coordinates producers and consumers
type ProducerConsumer struct {
	buffer  chan Item
	done    chan struct{}
	wg      sync.WaitGroup
}

// NewProducerConsumer creates a new system
// TODO: Initialize with buffer size
func NewProducerConsumer(bufferSize int) *ProducerConsumer {
	// TODO: Implement this function
	panic("not implemented")
}

// StartProducer launches a producer goroutine
// TODO: Implement producer
func (pc *ProducerConsumer) StartProducer(id int, numItems int) {
	// TODO: Produce items
	// TODO: Handle shutdown signal
	// TODO: Send items to buffer
	panic("not implemented")
}

// StartConsumer launches a consumer goroutine
// TODO: Implement consumer
func (pc *ProducerConsumer) StartConsumer(id int) {
	// TODO: Receive items from buffer
	// TODO: Process items
	// TODO: Handle shutdown
	panic("not implemented")
}

// Shutdown gracefully stops the system
// TODO: Implement graceful shutdown
func (pc *ProducerConsumer) Shutdown() {
	// TODO: Signal shutdown
	// TODO: Wait for goroutines
	// TODO: Close channels
	panic("not implemented")
}

func main() {
	fmt.Println("Producer-Consumer Pattern")

	pc := NewProducerConsumer(10)

	// Start producers
	for i := 0; i < 3; i++ {
		pc.StartProducer(i, 5)
	}

	// Start consumers
	for i := 0; i < 2; i++ {
		pc.StartConsumer(i)
	}

	time.Sleep(2 * time.Second)
	pc.Shutdown()
}
