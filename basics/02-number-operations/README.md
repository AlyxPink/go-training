# Exercise 02: Number Operations

**Difficulty**: ⭐ Beginner
**Estimated Time**: 35 minutes

## Learning Objectives

By completing this exercise, you will learn:
- Basic arithmetic operations in Go
- Loop constructs (for loop variants)
- Conditional logic
- Recursion fundamentals
- Integer overflow behavior
- Mathematical algorithm implementation

## Problem Description

Implement three number operation functions:

### 1. Prime Checker
Write a function that determines if a number is prime.

**Example**:
```go
IsPrime(7)  // returns true
IsPrime(10) // returns false
IsPrime(1)  // returns false
```

### 2. Factorial
Write a function that calculates the factorial of a number (n!).

**Example**:
```go
Factorial(5)  // returns 120 (5 * 4 * 3 * 2 * 1)
Factorial(0)  // returns 1
```

### 3. Fibonacci Sequence
Write a function that returns the nth Fibonacci number.

**Example**:
```go
Fibonacci(6)  // returns 8 (0, 1, 1, 2, 3, 5, 8)
Fibonacci(0)  // returns 0
```

## Requirements

- Implement both iterative and recursive solutions where applicable
- Handle edge cases (0, 1, negative numbers)
- Write efficient algorithms
- All tests must pass

## Testing

Run the tests with:
```bash
go test -v
```

## Key Concepts

- **Prime Number**: A number > 1 with no positive divisors other than 1 and itself
- **Factorial**: The product of all positive integers ≤ n
- **Fibonacci**: Each number is the sum of the two preceding ones
- **Iteration vs Recursion**: Different approaches with tradeoffs
