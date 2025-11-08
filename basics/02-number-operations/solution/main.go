package main

import (
	"fmt"
	"math"
)

// IsPrime checks if a number is prime.
// A prime number is a natural number greater than 1 that has no positive divisors other than 1 and itself.
func IsPrime(n int) bool {
	// Numbers less than 2 are not prime
	if n < 2 {
		return false
	}

	// 2 is the only even prime number
	if n == 2 {
		return true
	}

	// Even numbers are not prime
	if n%2 == 0 {
		return false
	}

	// Check odd divisors up to sqrt(n)
	// If n has a divisor greater than sqrt(n), it must also have one smaller than sqrt(n)
	sqrtN := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrtN; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

// Factorial calculates the factorial of n (n!) using iteration.
// The factorial of n is the product of all positive integers less than or equal to n.
func Factorial(n int) int {
	// Base case: 0! = 1
	if n == 0 {
		return 1
	}

	// Iterative approach
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}

	return result
}

// FactorialRecursive calculates factorial using recursion.
func FactorialRecursive(n int) int {
	// Base case
	if n <= 1 {
		return 1
	}

	// Recursive case: n! = n * (n-1)!
	return n * FactorialRecursive(n-1)
}

// Fibonacci returns the nth Fibonacci number using iteration.
// The Fibonacci sequence: 0, 1, 1, 2, 3, 5, 8, 13, ...
func Fibonacci(n int) int {
	// Base cases
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	// Iterative approach with two variables
	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}

	return curr
}

// FibonacciRecursive returns the nth Fibonacci number using recursion.
// Note: This is inefficient for large n (exponential time complexity).
func FibonacciRecursive(n int) int {
	// Base cases
	if n <= 1 {
		return n
	}

	// Recursive case: F(n) = F(n-1) + F(n-2)
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func main() {
	// Test your implementations
	fmt.Println("Prime checking:")
	fmt.Println("IsPrime(7):", IsPrime(7))
	fmt.Println("IsPrime(10):", IsPrime(10))
	fmt.Println()

	fmt.Println("Factorial:")
	fmt.Println("Factorial(5):", Factorial(5))
	fmt.Println("FactorialRecursive(5):", FactorialRecursive(5))
	fmt.Println()

	fmt.Println("Fibonacci:")
	fmt.Println("Fibonacci(6):", Fibonacci(6))
	fmt.Println("FibonacciRecursive(6):", FibonacciRecursive(6))
}
