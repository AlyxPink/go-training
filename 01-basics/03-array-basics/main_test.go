package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindMax(t *testing.T) {
	tests := []struct {
		name     string
		input    [5]int
		expected int
	}{
		{"positive numbers", [5]int{3, 7, 2, 9, 1}, 9},
		{"negative numbers", [5]int{-5, -2, -8, -1, -10}, -1},
		{"mixed numbers", [5]int{-3, 7, -2, 9, -1}, 9},
		{"all same", [5]int{5, 5, 5, 5, 5}, 5},
		{"max at start", [5]int{10, 1, 2, 3, 4}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindMax(tt.input)
			if result != tt.expected {
				t.Errorf("FindMax(%v) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRotateRight(t *testing.T) {
	tests := []struct {
		name     string
		input    [5]int
		k        int
		expected [5]int
	}{
		{"rotate by 2", [5]int{1, 2, 3, 4, 5}, 2, [5]int{4, 5, 1, 2, 3}},
		{"rotate by 0", [5]int{1, 2, 3, 4, 5}, 0, [5]int{1, 2, 3, 4, 5}},
		{"rotate by 5", [5]int{1, 2, 3, 4, 5}, 5, [5]int{1, 2, 3, 4, 5}},
		{"rotate by 1", [5]int{1, 2, 3, 4, 5}, 1, [5]int{5, 1, 2, 3, 4}},
		{"rotate by 7", [5]int{1, 2, 3, 4, 5}, 7, [5]int{4, 5, 1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RotateRight(tt.input, tt.k)
			if result != tt.expected {
				t.Errorf("RotateRight(%v, %d) = %v, expected %v", tt.input, tt.k, result, tt.expected)
			}
		})
	}
}

func TestFindDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    [7]int
		expected []int
	}{
		{"two duplicates", [7]int{1, 2, 3, 2, 4, 3, 5}, []int{2, 3}},
		{"no duplicates", [7]int{1, 2, 3, 4, 5, 6, 7}, []int{}},
		{"all duplicates", [7]int{1, 1, 1, 1, 1, 1, 1}, []int{1}},
		{"one duplicate", [7]int{1, 2, 3, 4, 5, 6, 1}, []int{1}},
		{"multiple same", [7]int{1, 2, 2, 2, 3, 3, 3}, []int{2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindDuplicates(tt.input)
			
			// Sort both slices for comparison since order doesn't matter
			sort.Ints(result)
			sort.Ints(tt.expected)
			
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindDuplicates(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
