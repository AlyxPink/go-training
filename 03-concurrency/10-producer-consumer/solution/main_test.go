package main

import (
	"testing"
	"time"
)

func TestProducerConsumer(t *testing.T) {
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

func TestProducerConsumerNoRace(t *testing.T) {
	pc := NewProducerConsumer(5)

	// Start multiple producers and consumers
	for i := 0; i < 5; i++ {
		pc.StartProducer(i, 10)
	}

	for i := 0; i < 3; i++ {
		pc.StartConsumer(i)
	}

	time.Sleep(3 * time.Second)
	pc.Shutdown()
}
