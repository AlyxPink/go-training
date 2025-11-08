package main

import "fmt"

// SendNumbers sends numbers 1 through n to the channel
// TODO: Implement number generation
func SendNumbers(ch chan<- int, n int) {
	// TODO: Send numbers 1 to n
	// TODO: Close channel when done
}

// ReceiveSum receives all numbers from channel and returns their sum
// TODO: Implement sum calculation
func ReceiveSum(ch <-chan int) int {
	// TODO: Receive all values until channel closed
	// TODO: Return sum
	return 0
}

// BufferedPipeline demonstrates buffered channel usage
// TODO: Implement a pipeline with buffered channels
func BufferedPipeline(numbers []int) []int {
	// TODO: Create buffered channel
	// TODO: Send numbers in goroutine
	// TODO: Transform numbers (e.g., square them)
	// TODO: Collect results
	return nil
}

// GenerateSquares creates a channel that emits squares of 1 to n
// TODO: Implement generator pattern
func GenerateSquares(n int) <-chan int {
	// TODO: Create channel
	// TODO: Launch goroutine to send values
	// TODO: Close channel after sending all values
	// TODO: Return channel
	return nil
}

func main() {
	fmt.Println("Channel Basics Examples")

	// Example 1: Simple send/receive
	ch := make(chan int)
	go SendNumbers(ch, 10)
	sum := ReceiveSum(ch)
	fmt.Printf("Sum of 1-10: %d\n", sum)

	// Example 2: Buffered pipeline
	numbers := []int{1, 2, 3, 4, 5}
	squares := BufferedPipeline(numbers)
	fmt.Printf("Squares: %v\n", squares)

	// Example 3: Generator pattern
	fmt.Println("Squares from generator:")
	for square := range GenerateSquares(5) {
		fmt.Printf("%d ", square)
	}
	fmt.Println()
}
