package main

import "fmt"

// FindMax returns the maximum value in the array.
func FindMax(arr [5]int) int {
	// TODO: Implement finding maximum value
	// Hint: Initialize with first element, then compare with rest
	return 0
}

// RotateRight rotates the array to the right by k positions.
func RotateRight(arr [5]int, k int) [5]int {
	// TODO: Implement array rotation
	// Hint: Use modulo to handle k > len(arr), create new array for result
	return [5]int{}
}

// FindDuplicates returns a slice of duplicate values in the array.
func FindDuplicates(arr [7]int) []int {
	// TODO: Implement duplicate finding
	// Hint: Use a map to track seen values, another to track duplicates
	return nil
}

func main() {
	// Test your implementations
	fmt.Println("FindMax:", FindMax([5]int{3, 7, 2, 9, 1}))
	fmt.Println("RotateRight:", RotateRight([5]int{1, 2, 3, 4, 5}, 2))
	fmt.Println("FindDuplicates:", FindDuplicates([7]int{1, 2, 3, 2, 4, 3, 5}))
}
