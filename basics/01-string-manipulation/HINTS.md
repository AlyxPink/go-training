# Hints for String Manipulation

## Level 1: Getting Started

- Remember that strings in Go are immutable - you cannot modify them in place
- When you need to build a string, consider using a slice of runes
- The `range` keyword over a string gives you runes, not bytes

## Level 2: String Reversal

- Convert the string to a slice of runes: `[]rune(s)`
- You can reverse a slice by swapping elements from both ends
- Use two pointers: one at the start, one at the end
- Convert back to string: `string(runes)`

## Level 3: Palindrome Logic

- A palindrome reads the same forwards and backwards
- You can reverse the string and compare, or use two pointers
- For case-insensitive comparison, use `strings.ToLower()`
- To ignore non-alphanumeric characters, use `unicode.IsLetter()` or `unicode.IsDigit()`

## Level 4: Character Counting

- Use a `map[rune]int` to store character counts
- Initialize the map with `make(map[rune]int)`
- Range over the string to get each rune
- Increment the count for each rune: `counts[r]++`

## Level 5: Edge Cases

Think about:
- Empty strings
- Strings with only one character
- Strings with special Unicode characters (emojis, Chinese characters, etc.)
- Strings with whitespace
- Case sensitivity
