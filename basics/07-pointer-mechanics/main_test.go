package main

import "testing"

func TestSwap(t *testing.T) {
	a, b := 10, 20
	Swap(&a, &b)
	if a != 20 || b != 10 {
		t.Errorf("Swap failed: a=%d, b=%d", a, b)
	}
}

func TestUpdateAge(t *testing.T) {
	p := &Person{Name: "Alice", Age: 25}
	UpdateAge(p, 30)
	if p.Age != 30 {
		t.Errorf("UpdateAge failed: got %d", p.Age)
	}
	
	UpdateAge(nil, 30) // Should not panic
}

func TestNilSafeIncrement(t *testing.T) {
	var nilPtr *int
	if result := NilSafeIncrement(nilPtr); result != 1 {
		t.Errorf("NilSafeIncrement(nil) = %d, expected 1", result)
	}
	
	val := 5
	if result := NilSafeIncrement(&val); result != 6 {
		t.Errorf("NilSafeIncrement(&5) = %d, expected 6", result)
	}
}
