package main

import (
	"fmt"
	"time"
)

// Multiplex merges two channels into one using select
// TODO: Implement channel multiplexing
func Multiplex(ch1, ch2 <-chan int) <-chan int {
	// TODO: Create output channel
	// TODO: Use select to read from either input
	// TODO: Handle both channels closing

	// Return closed channel to prevent infinite blocking in tests
	out := make(chan int)
	close(out)
	return out
}

// Timeout performs an operation with timeout
// TODO: Implement timeout pattern
func Timeout(ch <-chan string, timeout time.Duration) (string, bool) {
	// TODO: Use select with time.After
	// TODO: Return value and true if received, or "", false on timeout
	return "", false
}

// NonBlockingSend attempts non-blocking send
// TODO: Implement non-blocking send
func NonBlockingSend(ch chan int, value int) bool {
	// TODO: Use select with default case
	// TODO: Return true if sent, false if would block
	return false
}

// Worker processes jobs until quit signal
// TODO: Implement cancellable worker
func Worker(jobs <-chan int, quit <-chan struct{}) {
	// TODO: Use select to handle jobs or quit
	// TODO: Process jobs until quit received
}

func main() {
	fmt.Println("Select Statement Examples")

	// Multiplex example
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		defer close(ch1)
		for i := 0; i < 3; i++ {
			ch1 <- i
		}
	}()
	go func() {
		defer close(ch2)
		for i := 10; i < 13; i++ {
			ch2 <- i
		}
	}()

	fmt.Println("Multiplexed values:")
	for v := range Multiplex(ch1, ch2) {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// Timeout example
	ch := make(chan string, 1)
	ch <- "quick"
	if msg, ok := Timeout(ch, 1*time.Second); ok {
		fmt.Println("Received:", msg)
	}

	ch2Slow := make(chan string)
	if _, ok := Timeout(ch2Slow, 100*time.Millisecond); !ok {
		fmt.Println("Timeout occurred")
	}
}
