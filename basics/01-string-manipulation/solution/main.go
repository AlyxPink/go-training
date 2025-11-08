package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Reverse returns the reversed version of the input string.
// It properly handles Unicode characters by working with runes.
func Reverse(s string) string {
	// Convert string to slice of runes to handle Unicode correctly
	runes := []rune(s)

	// Reverse the slice using two pointers
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert back to string
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome.
// It ignores case and non-alphanumeric characters.
func IsPalindrome(s string) bool {
	// Normalize the string: lowercase and keep only alphanumeric characters
	var normalized []rune
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			normalized = append(normalized, r)
		}
	}

	// Check if normalized string is a palindrome using two pointers
	for i, j := 0, len(normalized)-1; i < j; i, j = i+1, j-1 {
		if normalized[i] != normalized[j] {
			return false
		}
	}

	return true
}

// CountChars returns a map with the count of each character in the string.
func CountChars(s string) map[rune]int {
	// Initialize the map
	counts := make(map[rune]int)

	// Range over string gives us runes, which handles Unicode properly
	for _, r := range s {
		counts[r]++
	}

	return counts
}

func main() {
	// Test your implementations here
	fmt.Println("Reverse('hello'):", Reverse("hello"))
	fmt.Println("Reverse('Hello, 世界'):", Reverse("Hello, 世界"))
	fmt.Println()

	fmt.Println("IsPalindrome('racecar'):", IsPalindrome("racecar"))
	fmt.Println("IsPalindrome('hello'):", IsPalindrome("hello"))
	fmt.Println("IsPalindrome('A man, a plan, a canal: Panama'):", IsPalindrome("A man, a plan, a canal: Panama"))
	fmt.Println()

	fmt.Println("CountChars('hello'):", CountChars("hello"))
	fmt.Println("CountChars('Hello, 世界'):", CountChars("Hello, 世界"))
}
