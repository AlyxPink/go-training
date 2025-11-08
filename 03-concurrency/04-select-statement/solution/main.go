package main

import (
	"fmt"
	"time"
)

func Multiplex(ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					continue
				}
				out <- v
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					continue
				}
				out <- v
			}
		}
	}()

	return out
}

func Timeout(ch <-chan string, timeout time.Duration) (string, bool) {
	select {
	case msg := <-ch:
		return msg, true
	case <-time.After(timeout):
		return "", false
	}
}

func NonBlockingSend(ch chan int, value int) bool {
	select {
	case ch <- value:
		return true
	default:
		return false
	}
}

func Worker(jobs <-chan int, quit <-chan struct{}) {
	for {
		select {
		case job := <-jobs:
			fmt.Printf("Processing job %d\n", job)
			time.Sleep(100 * time.Millisecond)
		case <-quit:
			fmt.Println("Worker quitting")
			return
		}
	}
}

func main() {
	fmt.Println("Select Statement Examples")

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

	ch := make(chan string, 1)
	ch <- "quick"
	if msg, ok := Timeout(ch, 1*time.Second); ok {
		fmt.Println("Received:", msg)
	}

	chSlow := make(chan string)
	if _, ok := Timeout(chSlow, 100*time.Millisecond); !ok {
		fmt.Println("Timeout occurred")
	}
}
