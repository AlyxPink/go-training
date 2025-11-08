package main

import "fmt"

// Generator creates a channel and sends values
func Generator(values ...int) <-chan int {
	// TODO: Create channel, send values in goroutine, return channel
	return nil
}

// Square receives from in, squares values, sends to returned channel
func Square(in <-chan int) <-chan int {
	// TODO: Implement pipeline stage
	return nil
}

// FanOut distributes work from in to numWorkers workers
func FanOut(in <-chan int, numWorkers int, work func(int) int) []<-chan int {
	// TODO: Create numWorkers, each processing from in
	return nil
}

// FanIn merges multiple channels into one
func FanIn(channels ...<-chan int) <-chan int {
	// TODO: Merge all channels into single output
	return nil
}

func main() {
	fmt.Println("Channel Patterns")
}
