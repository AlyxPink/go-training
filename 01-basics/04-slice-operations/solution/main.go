package main

import "fmt"

// Filter returns a new slice containing only even numbers from the input.
func Filter(nums []int) []int {
	// Create an empty slice to hold results
	// Using capacity hint can improve performance
	result := make([]int, 0, len(nums)/2)

	// Iterate and append even numbers
	for _, num := range nums {
		if num%2 == 0 {
			result = append(result, num)
		}
	}

	return result
}

// FilterAlternative using append to nil slice.
func FilterAlternative(nums []int) []int {
	var result []int // nil slice
	
	for _, num := range nums {
		if num%2 == 0 {
			result = append(result, num)
		}
	}
	
	// Return empty slice instead of nil for consistency
	if result == nil {
		return []int{}
	}
	
	return result
}

// Double returns a new slice with each element doubled.
func Double(nums []int) []int {
	// Pre-allocate slice with exact size needed
	result := make([]int, len(nums))

	// Transform each element
	for i, num := range nums {
		result[i] = num * 2
	}

	return result
}

// DoubleAlternative using append (less efficient due to potential reallocations).
func DoubleAlternative(nums []int) []int {
	var result []int
	
	for _, num := range nums {
		result = append(result, num*2)
	}
	
	if result == nil {
		return []int{}
	}
	
	return result
}

// RemoveAt removes the element at index i without preserving order.
// Returns the modified slice.
func RemoveAt(nums []int, i int) []int {
	// Bounds check
	if i < 0 || i >= len(nums) {
		return nums
	}

	// Swap with last element
	nums[i] = nums[len(nums)-1]

	// Truncate last element
	return nums[:len(nums)-1]
}

// RemoveAtPreserveOrder removes element at index i while preserving order.
func RemoveAtPreserveOrder(nums []int, i int) []int {
	if i < 0 || i >= len(nums) {
		return nums
	}

	// Shift elements left
	return append(nums[:i], nums[i+1:]...)
}

func main() {
	// Test your implementations
	fmt.Println("Filter even numbers:", Filter([]int{1, 2, 3, 4, 5, 6}))
	fmt.Println("FilterAlternative:", FilterAlternative([]int{1, 2, 3, 4, 5, 6}))
	fmt.Println()

	fmt.Println("Double:", Double([]int{1, 2, 3}))
	fmt.Println("DoubleAlternative:", DoubleAlternative([]int{1, 2, 3}))
	fmt.Println()

	s := []int{1, 2, 3, 4, 5}
	fmt.Println("Original:", s)
	s = RemoveAt(s, 2)
	fmt.Println("After RemoveAt(2):", s)
	
	s2 := []int{1, 2, 3, 4, 5}
	s2 = RemoveAtPreserveOrder(s2, 2)
	fmt.Println("RemoveAtPreserveOrder(2):", s2)
}
