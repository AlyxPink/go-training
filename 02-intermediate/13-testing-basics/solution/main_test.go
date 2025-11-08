package main

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed", -2, 3, 1},
		{"zero", 0, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"racecar", true},
		{"RaceCar", true},
		{"hello", false},
		{"a", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := IsPalindrome(tt.input)
			if got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(5, 3)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome("racecar")
	}
}
