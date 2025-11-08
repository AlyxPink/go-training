package main

import "fmt"

// WordFrequency counts how many times each word appears in the input slice.
func WordFrequency(words []string) map[string]int {
	// TODO: Implement word frequency counter
	// Hint: Create map, range over words, increment count for each word
	return nil
}

// InvertMap swaps keys and values. If multiple keys have the same value,
// only one will be kept (map iteration order is non-deterministic).
func InvertMap(m map[string]int) map[int]string {
	// TODO: Implement map inversion
	// Hint: Create new map with swapped types, range and swap k/v
	return nil
}

// MergeMaps merges two maps. Values from the second map override the first.
func MergeMaps(m1, m2 map[string]int) map[string]int {
	// TODO: Implement map merging
	// Hint: Create result map, copy from m1, then copy from m2
	return nil
}

func main() {
	// Test your implementations
	fmt.Println("WordFrequency:", WordFrequency([]string{"hello", "world", "hello"}))
	fmt.Println("InvertMap:", InvertMap(map[string]int{"a": 1, "b": 2}))
	fmt.Println("MergeMaps:", MergeMaps(
		map[string]int{"a": 1, "b": 2},
		map[string]int{"b": 3, "c": 4},
	))
}
