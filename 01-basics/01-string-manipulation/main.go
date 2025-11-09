package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

// Reverse returns the reversed version of the input string.
// It properly handles Unicode characters by working with runes.
func Reverse(s string) string {
	// TODO: Convert string to slice of runes
	// TODO: Reverse the slice using two pointers (swap from both ends)
	// TODO: Convert rune slice back to string and return
	runes := []rune(s)
	slices.Reverse(runes)
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome.
// It ignores case and non-alphanumeric characters.
func IsPalindrome(s string) bool {
	// TODO: Create slice to hold normalized runes
	// TODO: Convert to lowercase and filter oly letters and digits
	// TODO: Use two pointers to compare from both ends
	// TODO: Return false if any mismatch, true if all match
	var result string
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result += string(r)
		}
	}
	return result == Reverse(result)
}

// CountChars returns a map with the count of each character in the string.
func CountChars(s string) map[rune]int {
	// TODO: Create map to store character counts
	// TODO: Range over string (gives runes automatically)
	// TODO: Increment count for each rune in the map
	// TODO: Return the counts map
	count := map[rune]int{}
	for _, r := range s {
		count[r]++
	}
	return count
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
