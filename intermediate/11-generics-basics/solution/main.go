package main

import "fmt"

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func Filter[T any](slice []T, pred func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

func Max[T comparable](values ...T) T {
	if len(values) == 0 {
		var zero T
		return zero
	}
	
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func main() {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	
	for !stack.IsEmpty() {
		val, _ := stack.Pop()
		fmt.Println(val)
	}
	
	numbers := []int{1, 2, 3, 4}
	doubled := Map(numbers, func(n int) int { return n * 2 })
	fmt.Println("Doubled:", doubled)
	
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("Evens:", evens)
	
	fmt.Println("Max:", Max(1, 5, 3, 2))
}
