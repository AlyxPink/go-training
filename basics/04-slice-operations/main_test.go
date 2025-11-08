package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"mixed numbers", []int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}},
		{"no even numbers", []int{1, 3, 5, 7}, []int{}},
		{"all even numbers", []int{2, 4, 6, 8}, []int{2, 4, 6, 8}},
		{"empty slice", []int{}, []int{}},
		{"single even", []int{2}, []int{2}},
		{"single odd", []int{1}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDouble(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"simple numbers", []int{1, 2, 3}, []int{2, 4, 6}},
		{"with zero", []int{0, 1, 2}, []int{0, 2, 4}},
		{"negative numbers", []int{-1, -2, -3}, []int{-2, -4, -6}},
		{"empty slice", []int{}, []int{}},
		{"single element", []int{5}, []int{10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Double(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Double(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveAt(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		index int
	}{
		{"remove middle", []int{1, 2, 3, 4, 5}, 2},
		{"remove first", []int{1, 2, 3, 4, 5}, 0},
		{"remove last", []int{1, 2, 3, 4, 5}, 4},
		{"two elements", []int{1, 2}, 0},
		{"single element", []int{1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to track what was removed
			original := make([]int, len(tt.input))
			copy(original, tt.input)
			removedValue := original[tt.index]

			result := RemoveAt(tt.input, tt.index)

			// Check length decreased by 1
			if len(result) != len(original)-1 {
				t.Errorf("RemoveAt(%v, %d) length = %d, expected %d", 
					original, tt.index, len(result), len(original)-1)
			}

			// Check removed value is not in result
			for _, v := range result {
				if v == removedValue {
					// Count occurrences in original
					countOriginal := 0
					countResult := 0
					for _, ov := range original {
						if ov == removedValue {
							countOriginal++
						}
					}
					for _, rv := range result {
						if rv == removedValue {
							countResult++
						}
					}
					// Result should have one less occurrence
					if countResult != countOriginal-1 {
						t.Errorf("RemoveAt did not remove the correct element")
					}
					break
				}
			}

			// Check all other elements are present
			sortedOriginal := make([]int, len(original))
			copy(sortedOriginal, original)
			sort.Ints(sortedOriginal)

			sortedResult := make([]int, len(result))
			copy(sortedResult, result)
			sort.Ints(sortedResult)

			// Remove one occurrence of the removed value from sorted original
			for i, v := range sortedOriginal {
				if v == removedValue {
					sortedOriginal = append(sortedOriginal[:i], sortedOriginal[i+1:]...)
					break
				}
			}

			if !reflect.DeepEqual(sortedResult, sortedOriginal) {
				t.Errorf("RemoveAt(%v, %d) = %v, elements don't match", 
					original, tt.index, result)
			}
		})
	}
}
