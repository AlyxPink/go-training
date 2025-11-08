package main

import "fmt"

// WordFrequency counts how many times each word appears in the input slice.
func WordFrequency(words []string) map[string]int {
	// Initialize the map
	freq := make(map[string]int)

	// Count each word
	// Zero value for int is 0, so this works for new keys
	for _, word := range words {
		freq[word]++
	}

	return freq
}

// InvertMap swaps keys and values. If multiple keys have the same value,
// only one will be kept (map iteration order is non-deterministic).
func InvertMap(m map[string]int) map[int]string {
	// Create inverted map
	inverted := make(map[int]string)

	// Swap keys and values
	for key, value := range m {
		inverted[value] = key
	}

	return inverted
}

// MergeMaps merges two maps. Values from the second map override the first.
func MergeMaps(m1, m2 map[string]int) map[string]int {
	// Create result map
	result := make(map[string]int)

	// Copy all entries from first map
	for k, v := range m1 {
		result[k] = v
	}

	// Copy all entries from second map (overwrites duplicates)
	for k, v := range m2 {
		result[k] = v
	}

	return result
}

// MergeMapsInPlace modifies the first map (more memory efficient).
func MergeMapsInPlace(m1, m2 map[string]int) map[string]int {
	// Add/overwrite entries from m2 into m1
	for k, v := range m2 {
		m1[k] = v
	}

	return m1
}

func main() {
	// Test your implementations
	fmt.Println("WordFrequency:", WordFrequency([]string{"hello", "world", "hello"}))
	fmt.Println()

	fmt.Println("InvertMap:", InvertMap(map[string]int{"a": 1, "b": 2, "c": 3}))
	fmt.Println()

	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	fmt.Println("MergeMaps:", MergeMaps(m1, m2))
	fmt.Println("Original maps unchanged:", m1, m2)
}
