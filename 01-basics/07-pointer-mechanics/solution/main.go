package main

import "fmt"

func Swap(a, b *int) {
	*a, *b = *b, *a
}

type Person struct {
	Name string
	Age  int
}

func UpdateAge(p *Person, newAge int) {
	if p != nil {
		p.Age = newAge
	}
}

func NilSafeIncrement(p *int) int {
	if p == nil {
		return 1
	}
	return *p + 1
}

func main() {
	x, y := 10, 20
	fmt.Printf("Before: x=%d, y=%d\n", x, y)
	Swap(&x, &y)
	fmt.Printf("After: x=%d, y=%d\n", x, y)
	
	p := &Person{Name: "Alice", Age: 25}
	UpdateAge(p, 26)
	fmt.Println("Updated person:", p)
	
	var nilPtr *int
	fmt.Println("NilSafeIncrement(nil):", NilSafeIncrement(nilPtr))
	val := 5
	fmt.Println("NilSafeIncrement(&5):", NilSafeIncrement(&val))
}
