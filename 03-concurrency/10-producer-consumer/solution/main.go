package main

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	ID   int
	Data string
}

type ProducerConsumer struct {
	buffer chan Item
	done   chan struct{}
	wg     sync.WaitGroup
}

func NewProducerConsumer(bufferSize int) *ProducerConsumer {
	return &ProducerConsumer{
		buffer: make(chan Item, bufferSize),
		done:   make(chan struct{}),
	}
}

func (pc *ProducerConsumer) StartProducer(id int, numItems int) {
	pc.wg.Add(1)
	go func() {
		defer pc.wg.Done()

		for i := 0; i < numItems; i++ {
			item := Item{
				ID:   i,
				Data: fmt.Sprintf("Producer %d item %d", id, i),
			}

			select {
			case pc.buffer <- item:
				fmt.Printf("Producer %d: sent item %d\n", id, i)
			case <-pc.done:
				fmt.Printf("Producer %d: shutting down\n", id)
				return
			}

			time.Sleep(100 * time.Millisecond)
		}

		fmt.Printf("Producer %d: finished\n", id)
	}()
}

func (pc *ProducerConsumer) StartConsumer(id int) {
	pc.wg.Add(1)
	go func() {
		defer pc.wg.Done()

		// Capture done channel to avoid race with Shutdown
		done := pc.done

		for {
			select {
			case item := <-pc.buffer:
				fmt.Printf("Consumer %d: processing item %d\n", id, item.ID)
				time.Sleep(200 * time.Millisecond)
			case <-done:
				// Drain remaining items
				for {
					select {
					case item := <-pc.buffer:
						fmt.Printf("Consumer %d: draining item %d\n", id, item.ID)
					default:
						fmt.Printf("Consumer %d: shutting down\n", id)
						return
					}
				}
			}
		}
	}()
}

func (pc *ProducerConsumer) Shutdown() {
	fmt.Println("Initiating shutdown...")
	close(pc.done)
	pc.wg.Wait()
	close(pc.buffer)
	fmt.Println("Shutdown complete")
}

func main() {
	fmt.Println("Producer-Consumer Pattern")

	pc := NewProducerConsumer(10)

	for i := 0; i < 3; i++ {
		pc.StartProducer(i, 5)
	}

	for i := 0; i < 2; i++ {
		pc.StartConsumer(i)
	}

	time.Sleep(2 * time.Second)
	pc.Shutdown()
}
