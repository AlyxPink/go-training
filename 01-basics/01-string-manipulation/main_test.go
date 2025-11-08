package main

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"single character", "a", "a"},
		{"simple word", "hello", "olleh"},
		{"unicode characters", "Hello, 世界", "界世 ,olleH"},
		{"emoji", "Hello 👋 World", "dlroW 👋 olleH"},
		{"palindrome", "racecar", "racecar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", true},
		{"single character", "a", true},
		{"simple palindrome", "racecar", true},
		{"not palindrome", "hello", false},
		{"case insensitive palindrome", "RaceCar", true},
		{"palindrome with spaces", "race car", true},
		{"complex palindrome", "A man a plan a canal Panama", true},
		{"palindrome with punctuation", "A man, a plan, a canal: Panama", true},
		{"numeric palindrome", "12321", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[rune]int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: map[rune]int{},
		},
		{
			name:  "single character",
			input: "a",
			expected: map[rune]int{
				'a': 1,
			},
		},
		{
			name:  "repeated characters",
			input: "hello",
			expected: map[rune]int{
				'h': 1,
				'e': 1,
				'l': 2,
				'o': 1,
			},
		},
		{
			name:  "unicode characters",
			input: "Hello, 世界",
			expected: map[rune]int{
				'H': 1,
				'e': 1,
				'l': 2,
				'o': 1,
				',': 1,
				' ': 1,
				'世': 1,
				'界': 1,
			},
		},
		{
			name:  "case sensitive",
			input: "AaBbCc",
			expected: map[rune]int{
				'A': 1,
				'a': 1,
				'B': 1,
				'b': 1,
				'C': 1,
				'c': 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountChars(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CountChars(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
