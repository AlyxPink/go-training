package main

import "fmt"

// IsPrime checks if a number is prime.
// A prime number is a natural number greater than 1 that has no positive divisors other than 1 and itself.
func IsPrime(n int) bool {
	// TODO: Implement prime checker
	// Hint: Numbers < 2 are not prime. Check divisibility up to sqrt(n)
	return false
}

// Factorial calculates the factorial of n (n!).
// The factorial of n is the product of all positive integers less than or equal to n.
func Factorial(n int) int {
	// TODO: Implement factorial calculation
	// Hint: 0! = 1, n! = n * (n-1)!
	return 0
}

// Fibonacci returns the nth Fibonacci number.
// The Fibonacci sequence: 0, 1, 1, 2, 3, 5, 8, 13, ...
func Fibonacci(n int) int {
	// TODO: Implement Fibonacci calculation
	// Hint: F(0) = 0, F(1) = 1, F(n) = F(n-1) + F(n-2)
	return 0
}

func main() {
	// Test your implementations
	fmt.Println("IsPrime(7):", IsPrime(7))
	fmt.Println("Factorial(5):", Factorial(5))
	fmt.Println("Fibonacci(6):", Fibonacci(6))
}
