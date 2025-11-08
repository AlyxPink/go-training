# Solution Explanation: String Manipulation

## Key Go Concepts Demonstrated

### 1. Strings vs Runes vs Bytes

In Go, strings are **immutable sequences of bytes**. However, Go uses UTF-8 encoding, which means a single character can be represented by 1-4 bytes.

```go
s := "Hello, 世界"
fmt.Println(len(s))           // 13 (bytes)
fmt.Println(len([]rune(s)))   // 9 (characters/runes)
```

**Important**: When working with strings that may contain Unicode characters, always use `rune` (alias for `int32`), not `byte`.

### 2. String Reversal Idiom

The standard Go approach to reversing a string:

```go
runes := []rune(s)                          // Convert to rune slice
for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
    runes[i], runes[j] = runes[j], runes[i] // Swap
}
return string(runes)                        // Convert back
```

**Why this works**:
- Converting to `[]rune` ensures each Unicode character is one element
- Two-pointer swap is an efficient O(n/2) algorithm
- Multiple assignment `i, j = i+1, j-1` is idiomatic Go

### 3. Range Over Strings

When you range over a string, you get **runes**, not bytes:

```go
for i, r := range "Hello, 世界" {
    // i is the byte index
    // r is the rune (Unicode code point)
}
```

This is why `range` is perfect for character counting.

### 4. Map Initialization and Zero Values

Go maps have a useful property: accessing a non-existent key returns the **zero value** for the value type.

```go
counts := make(map[rune]int)
counts['a']++  // Works! counts['a'] starts at 0, becomes 1
```

This eliminates the need for "key exists" checks in counting scenarios.

### 5. Unicode Package

The `unicode` package provides useful functions for character classification:

```go
unicode.IsLetter(r)  // Checks if rune is a letter
unicode.IsDigit(r)   // Checks if rune is a digit
unicode.IsSpace(r)   // Checks if rune is whitespace
```

These work correctly for all Unicode characters, not just ASCII.

### 6. Strings Package

The `strings` package provides common string operations:

```go
strings.ToLower(s)   // Convert to lowercase
strings.ToUpper(s)   // Convert to uppercase
strings.TrimSpace(s) // Remove leading/trailing whitespace
```

## Algorithm Analysis

### Reverse Function
- **Time Complexity**: O(n) where n is the number of runes
- **Space Complexity**: O(n) for the rune slice
- **Why not in-place?**: Strings are immutable in Go

### IsPalindrome Function
- **Time Complexity**: O(n) for normalization + O(n) for comparison = O(n)
- **Space Complexity**: O(n) for normalized slice
- **Alternative**: Could use Reverse() but that creates two extra slices

### CountChars Function
- **Time Complexity**: O(n) single pass through string
- **Space Complexity**: O(k) where k is the number of unique characters
- **Map efficiency**: Hash maps provide O(1) average case for insert/lookup

## Common Pitfalls

1. **Using bytes instead of runes**: `s[i]` gives you a byte, not a character
2. **String mutability**: You cannot do `s[0] = 'H'` in Go
3. **Length confusion**: `len(s)` returns bytes, not character count
4. **Range index**: The index in `for i, r := range s` is the byte position, not character position

## Go Idioms Used

1. **Multiple return values in loop**: `i, j = i+1, j-1`
2. **Swap idiom**: `a, b = b, a` (no temporary variable needed)
3. **Map zero value**: Leveraging `counts[r]++` without checking existence
4. **Conversion syntax**: `[]rune(s)` and `string(runes)`
5. **Table-driven tests**: The test file demonstrates this Go testing pattern
