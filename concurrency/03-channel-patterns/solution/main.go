package main

import (
	"fmt"
	"sync"
)

func Pipeline(numbers []int) <-chan int {
	// Stage 1: Generate
	gen := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for _, n := range numbers {
				ch <- n
			}
		}()
		return ch
	}

	// Stage 2: Square
	sq := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for v := range in {
				out <- v * v
			}
		}()
		return out
	}

	// Stage 3: Filter even
	filter := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for v := range in {
				if v%2 == 0 {
					out <- v
				}
			}
		}()
		return out
	}

	return filter(sq(gen()))
}

func FanOut(in <-chan int, n int) []<-chan int {
	channels := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		ch := make(chan int)
		channels[i] = ch

		go func(out chan int) {
			defer close(out)
			for v := range in {
				out <- v * 2
			}
		}(ch)
	}
	return channels
}

func FanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(channels))
	for _, ch := range channels {
		ch := ch
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func WorkerPool(jobs <-chan Job, numWorkers int) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- Result{
					JobID:  job.ID,
					Output: job.Data * job.Data,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

type Job struct {
	ID   int
	Data int
}

type Result struct {
	JobID  int
	Output int
}

func main() {
	fmt.Println("Channel Patterns Demo")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("Pipeline results:")
	for result := range Pipeline(numbers) {
		fmt.Printf("%d ", result)
	}
	fmt.Println()
}
