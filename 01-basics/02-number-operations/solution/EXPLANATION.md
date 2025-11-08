# Solution Explanation: Number Operations

## Key Go Concepts Demonstrated

### 1. For Loop Variants

Go has only one loop construct: `for`. But it's very versatile:

```go
// Classic for loop
for i := 0; i < 10; i++ {
    // ...
}

// While-style loop
for n > 0 {
    // ...
}

// Infinite loop
for {
    // ...
}
```

### 2. Integer Types and Operations

Go has multiple integer types with different sizes:
- `int`, `int8`, `int16`, `int32`, `int64` (signed)
- `uint`, `uint8`, `uint16`, `uint32`, `uint64` (unsigned)

**Integer division truncates**:
```go
5 / 2 == 2  // Not 2.5
```

**Integer overflow wraps around**:
```go
var x int8 = 127
x++  // x is now -128
```

### 3. Multiple Assignment

Go supports multiple assignment, which is perfect for swapping:

```go
prev, curr = curr, prev+curr  // No temporary variable needed
```

This evaluates all right-hand expressions before assignment, so you can safely swap.

### 4. Recursion in Go

Go supports recursion but doesn't optimize tail calls:

```go
// Recursive factorial
func Factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * Factorial(n-1)
}
```

**Important**: Deep recursion can cause stack overflow. Go's default stack size is small but grows dynamically.

### 5. Math Package

The `math` package provides mathematical functions:

```go
import "math"

math.Sqrt(float64(n))  // Square root
math.Pow(x, y)         // Power
math.Abs(x)            // Absolute value
```

## Algorithm Analysis

### Prime Checking
**Time Complexity**: O(√n)
**Space Complexity**: O(1)

**Optimizations**:
1. Check if n < 2 (not prime)
2. Check if n == 2 (only even prime)
3. Check if n is even (not prime)
4. Check odd divisors from 3 to √n

### Factorial
**Iterative**:
- Time: O(n)
- Space: O(1)

**Recursive**:
- Time: O(n)
- Space: O(n) for call stack

**Iterative is preferred** for better space efficiency.

### Fibonacci
**Iterative**:
- Time: O(n)
- Space: O(1)

**Recursive (naive)**:
- Time: O(2^n) - exponential!
- Space: O(n) for call stack

**Iterative is strongly preferred** for performance.

## Common Pitfalls

1. **Integer overflow**: Factorial grows very quickly
2. **Negative numbers**: Handle edge cases
3. **Off-by-one errors**: Be careful with loop bounds
4. **Inefficient recursion**: Naive Fibonacci is extremely slow

## Go Idioms Used

1. **Multiple return in loop**: `prev, curr = curr, prev+curr`
2. **Early return**: Check edge cases first, return early
3. **Type conversion**: `float64(n)` for math functions
4. **Increment operators**: `i++` is a statement, not an expression
5. **Modulo operator**: `n % i` for divisibility checking

## Performance Comparison

For Fibonacci(40):
- **Iterative**: ~1 microsecond
- **Recursive**: ~1 second (1,000,000x slower!)

This demonstrates why algorithm choice matters in Go.
