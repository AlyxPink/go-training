package stringutil

import "testing"

func TestReverse(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
	}

	for _, tt := range tests {
		got := Reverse(tt.input)
		if got != tt.want {
			t.Errorf("Reverse(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	if !IsPalindrome("racecar") {
		t.Error("racecar should be palindrome")
	}
	if IsPalindrome("hello") {
		t.Error("hello should not be palindrome")
	}
}
