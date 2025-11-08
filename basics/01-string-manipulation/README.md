# Exercise 01: String Manipulation

**Difficulty**: ⭐ Beginner
**Estimated Time**: 30 minutes

## Learning Objectives

By completing this exercise, you will learn:
- How strings are represented in Go (immutable byte sequences)
- The difference between bytes and runes (UTF-8 handling)
- String iteration techniques
- Common string manipulation patterns
- When to use `strings` package vs manual manipulation

## Problem Description

Implement three string manipulation functions:

### 1. Reverse String
Write a function that reverses a string while properly handling Unicode characters.

**Example**:
```go
Reverse("hello") // returns "olleh"
Reverse("Hello, 世界") // returns "界世 ,olleH"
```

### 2. Palindrome Checker
Write a function that checks if a string is a palindrome (reads the same forwards and backwards).

**Example**:
```go
IsPalindrome("racecar") // returns true
IsPalindrome("hello") // returns false
IsPalindrome("A man, a plan, a canal: Panama") // returns true (bonus: ignore punctuation/case)
```

### 3. Character Counter
Write a function that counts the occurrences of each character in a string.

**Example**:
```go
CountChars("hello") // returns map[rune]int{'h':1, 'e':1, 'l':2, 'o':1}
```

## Requirements

- Handle Unicode characters correctly (use runes, not bytes)
- Do not use reverse functions from external libraries
- Write clean, idiomatic Go code
- All tests must pass

## Testing

Run the tests with:
```bash
go test -v
```

## Key Concepts

- **String**: Immutable sequence of bytes
- **Rune**: Alias for `int32`, represents a Unicode code point
- **UTF-8**: Variable-length encoding used by Go strings
- **Range loop**: When ranging over a string, you get rune values, not bytes
