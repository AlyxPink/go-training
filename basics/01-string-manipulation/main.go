package main

import (
	"fmt"
)

// Reverse returns the reversed version of the input string.
// It properly handles Unicode characters by working with runes.
func Reverse(s string) string {
	// TODO: Implement string reversal
	// Hint: Convert to []rune, reverse the slice, convert back to string
	return ""
}

// IsPalindrome checks if a string is a palindrome.
// It ignores case and non-alphanumeric characters.
func IsPalindrome(s string) bool {
	// TODO: Implement palindrome check
	// Hint: You can use Reverse() or compare characters from both ends
	return false
}

// CountChars returns a map with the count of each character in the string.
func CountChars(s string) map[rune]int {
	// TODO: Implement character counting
	// Hint: Use a map[rune]int and range over the string
	return nil
}

func main() {
	// Test your implementations here
	fmt.Println("Reverse('hello'):", Reverse("hello"))
	fmt.Println("IsPalindrome('racecar'):", IsPalindrome("racecar"))
	fmt.Println("CountChars('hello'):", CountChars("hello"))
}
