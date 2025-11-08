package main

import (
	"fmt"
	"math"
)

// FindMax returns the maximum value in the array.
func FindMax(arr [5]int) int {
	// Initialize with the first element
	max := arr[0]

	// Compare with remaining elements
	for i := 1; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
		}
	}

	return max
}

// FindMaxAlternative uses math.MinInt as initial value.
func FindMaxAlternative(arr [5]int) int {
	max := math.MinInt

	for _, val := range arr {
		if val > max {
			max = val
		}
	}

	return max
}

// RotateRight rotates the array to the right by k positions.
func RotateRight(arr [5]int, k int) [5]int {
	n := len(arr)
	
	// Normalize k to be within array bounds
	k = k % n
	
	// Handle k == 0 case
	if k == 0 {
		return arr
	}

	// Create result array
	var result [5]int

	// Copy elements to their new positions
	for i := 0; i < n; i++ {
		// New position is (i + k) % n
		newPos := (i + k) % n
		result[newPos] = arr[i]
	}

	return result
}

// FindDuplicates returns a slice of duplicate values in the array.
func FindDuplicates(arr [7]int) []int {
	// Map to track how many times we've seen each value
	seen := make(map[int]int)

	// Count occurrences
	for _, val := range arr {
		seen[val]++
	}

	// Collect duplicates (values that appear more than once)
	var duplicates []int
	for val, count := range seen {
		if count > 1 {
			duplicates = append(duplicates, val)
		}
	}

	// Return empty slice if no duplicates
	if duplicates == nil {
		return []int{}
	}

	return duplicates
}

func main() {
	// Test your implementations
	fmt.Println("FindMax:", FindMax([5]int{3, 7, 2, 9, 1}))
	fmt.Println("FindMaxAlternative:", FindMaxAlternative([5]int{3, 7, 2, 9, 1}))
	fmt.Println()

	fmt.Println("RotateRight by 2:", RotateRight([5]int{1, 2, 3, 4, 5}, 2))
	fmt.Println("RotateRight by 0:", RotateRight([5]int{1, 2, 3, 4, 5}, 0))
	fmt.Println()

	fmt.Println("FindDuplicates:", FindDuplicates([7]int{1, 2, 3, 2, 4, 3, 5}))
	fmt.Println("No duplicates:", FindDuplicates([7]int{1, 2, 3, 4, 5, 6, 7}))
}
