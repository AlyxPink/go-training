package main

import "fmt"

// Filter returns a new slice containing only even numbers from the input.
func Filter(nums []int) []int {
	// TODO: Implement filter to keep only even numbers
	// Hint: Create empty slice, range over input, append even numbers
	return nil
}

// Double returns a new slice with each element doubled.
func Double(nums []int) []int {
	// TODO: Implement doubling transformation
	// Hint: Create slice with same length, transform each element
	return nil
}

// RemoveAt removes the element at index i without preserving order.
// Returns the modified slice.
func RemoveAt(nums []int, i int) []int {
	// TODO: Implement removal
	// Hint: Swap with last element, then truncate
	return nil
}

func main() {
	// Test your implementations
	fmt.Println("Filter:", Filter([]int{1, 2, 3, 4, 5, 6}))
	fmt.Println("Double:", Double([]int{1, 2, 3}))
	fmt.Println("RemoveAt:", RemoveAt([]int{1, 2, 3, 4, 5}, 2))
}
