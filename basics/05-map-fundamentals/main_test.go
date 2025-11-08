package main

import (
	"reflect"
	"testing"
)

func TestWordFrequency(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]int
	}{
		{
			"simple words",
			[]string{"hello", "world", "hello"},
			map[string]int{"hello": 2, "world": 1},
		},
		{
			"empty slice",
			[]string{},
			map[string]int{},
		},
		{
			"single word repeated",
			[]string{"test", "test", "test"},
			map[string]int{"test": 3},
		},
		{
			"all unique",
			[]string{"a", "b", "c"},
			map[string]int{"a": 1, "b": 1, "c": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WordFrequency(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("WordFrequency(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInvertMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected map[int]string
	}{
		{
			"simple map",
			map[string]int{"a": 1, "b": 2, "c": 3},
			map[int]string{1: "a", 2: "b", 3: "c"},
		},
		{
			"empty map",
			map[string]int{},
			map[int]string{},
		},
		{
			"single entry",
			map[string]int{"x": 100},
			map[int]string{100: "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InvertMap(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InvertMap(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	tests := []struct {
		name     string
		m1       map[string]int
		m2       map[string]int
		expected map[string]int
	}{
		{
			"no overlap",
			map[string]int{"a": 1, "b": 2},
			map[string]int{"c": 3, "d": 4},
			map[string]int{"a": 1, "b": 2, "c": 3, "d": 4},
		},
		{
			"with overlap",
			map[string]int{"a": 1, "b": 2},
			map[string]int{"b": 3, "c": 4},
			map[string]int{"a": 1, "b": 3, "c": 4},
		},
		{
			"first empty",
			map[string]int{},
			map[string]int{"a": 1},
			map[string]int{"a": 1},
		},
		{
			"second empty",
			map[string]int{"a": 1},
			map[string]int{},
			map[string]int{"a": 1},
		},
		{
			"both empty",
			map[string]int{},
			map[string]int{},
			map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeMaps(tt.m1, tt.m2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeMaps(%v, %v) = %v, expected %v", tt.m1, tt.m2, result, tt.expected)
			}
		})
	}
}
