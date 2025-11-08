package main

import "testing"

// TODO: Write table-driven test for Add
func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		// TODO: Add test cases
		{"positive numbers", 2, 3, 5},
		{"with zero", 5, 0, 5},
		{"negative numbers", -2, -3, -5},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TODO: Write table-driven test for Multiply
func TestMultiply(t *testing.T) {
	// TODO: Similar to TestAdd
}

// TODO: Write table-driven test for IsPalindrome
func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// TODO: Add test cases
		{"simple palindrome", "racecar", true},
		{"not palindrome", "hello", false},
		{"single char", "a", true},
		{"empty string", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.input); got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TODO: Write benchmark for Fibonacci
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}

// TODO: Write benchmark for IsPalindrome
func BenchmarkIsPalindrome(b *testing.B) {
	s := "racecar"
	for i := 0; i < b.N; i++ {
		IsPalindrome(s)
	}
}

// Helper function example
func assertEqual(t *testing.T, got, want int) {
	t.Helper()  // Marks as helper for better error messages
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
