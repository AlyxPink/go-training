# Hints for Number Operations

## Level 1: Getting Started

- Think about edge cases first: what should happen with 0, 1, or negative numbers?
- The `for` loop in Go is the only loop construct but is very versatile
- Integer division in Go truncates: `5 / 2 == 2`

## Level 2: Prime Number Logic

- Numbers less than 2 are not prime
- You only need to check divisors up to sqrt(n)
- If n is divisible by any number from 2 to sqrt(n), it's not prime
- Use modulo operator `%` to check divisibility

## Level 3: Factorial Approaches

- Factorial of 0 and 1 is 1 (base cases)
- Iterative: multiply numbers from 1 to n
- Recursive: n! = n * (n-1)!
- Watch out for integer overflow with large numbers

## Level 4: Fibonacci Strategies

- First two Fibonacci numbers are 0 and 1
- Iterative approach: maintain two variables for previous values
- Recursive approach: F(n) = F(n-1) + F(n-2)
- Naive recursion is very slow for large n (exponential time)

## Level 5: Optimization Tips

- For prime checking, you can skip even numbers after checking 2
- For Fibonacci, iterative is much faster than naive recursion
- Consider using `uint` for non-negative results
- The `math` package has `math.Sqrt()` for prime checking
