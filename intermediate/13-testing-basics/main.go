package main

import (
	"strings"
)

// Add returns sum of two integers
func Add(a, b int) int {
	return a + b
}

// Multiply returns product of two integers
func Multiply(a, b int) int {
	return a * b
}

// IsPalindrome checks if string is palindrome
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

// Fibonacci returns nth fibonacci number
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func main() {}
