package main

import "fmt"

// Stack is a generic stack
type Stack[T any] struct {
	// TODO: Add slice to store items
}

// NewStack creates a new stack
func NewStack[T any]() *Stack[T] {
	// TODO: Initialize stack
	return nil
}

// Push adds item to stack
func (s *Stack[T]) Push(item T) {
	// TODO: Add to stack
}

// Pop removes and returns top item
func (s *Stack[T]) Pop() (T, bool) {
	// TODO: Remove from stack, return zero value and false if empty
	var zero T
	return zero, false
}

// IsEmpty returns true if stack is empty
func (s *Stack[T]) IsEmpty() bool {
	// TODO: Check if empty
	return true
}

// Queue is a generic queue
type Queue[T any] struct {
	// TODO: Add slice for queue
}

// NewQueue creates a new queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Initialize queue
	return nil
}

// Enqueue adds item to back of queue
func (q *Queue[T]) Enqueue(item T) {
	// TODO: Add to queue
}

// Dequeue removes and returns front item
func (q *Queue[T]) Dequeue() (T, bool) {
	// TODO: Remove from front
	var zero T
	return zero, false
}

// Map applies function to each element
func Map[T any, U any](slice []T, fn func(T) U) []U {
	// TODO: Transform slice
	return nil
}

// Filter returns elements matching predicate
func Filter[T any](slice []T, pred func(T) bool) []T {
	// TODO: Filter slice
	return nil
}

// Max returns maximum of comparable values
func Max[T comparable](values ...T) T {
	// TODO: Find maximum
	var zero T
	return zero
}

func main() {
	// Stack example
	stack := NewStack[int]()
	if stack != nil {
		stack.Push(1)
		stack.Push(2)
		stack.Push(3)
		
		for !stack.IsEmpty() {
			val, _ := stack.Pop()
			fmt.Println(val)
		}
	}
	
	// Map example
	numbers := []int{1, 2, 3, 4}
	doubled := Map(numbers, func(n int) int { return n * 2 })
	fmt.Println("Doubled:", doubled)
	
	// Filter example
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("Evens:", evens)
}
