package main

import "fmt"

// Swap swaps the values of two integers using pointers.
func Swap(a, b *int) {
	// TODO: Swap the values that a and b point to
}

// UpdatePerson modifies a Person struct via pointer.
type Person struct {
	Name string
	Age  int
}

func UpdateAge(p *Person, newAge int) {
	// TODO: Update the age of the person
}

// NilSafe safely increments a pointer to int, handling nil.
func NilSafeIncrement(p *int) int {
	// TODO: Return incremented value, or 1 if p is nil
	return 0
}

func main() {
	x, y := 10, 20
	Swap(&x, &y)
	fmt.Printf("After swap: x=%d, y=%d\n", x, y)
}
