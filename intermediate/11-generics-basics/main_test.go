package main

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack[int]()
	
	if !stack.IsEmpty() {
		t.Error("New stack should be empty")
	}
	
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	
	if stack.IsEmpty() {
		t.Error("Stack should not be empty after pushes")
	}
	
	val, ok := stack.Pop()
	if !ok || val != 3 {
		t.Errorf("Pop() = %d, %v; want 3, true", val, ok)
	}
	
	val, ok = stack.Pop()
	if !ok || val != 2 {
		t.Errorf("Pop() = %d, %v; want 2, true", val, ok)
	}
	
	val, ok = stack.Pop()
	if !ok || val != 1 {
		t.Errorf("Pop() = %d, %v; want 1, true", val, ok)
	}
	
	if !stack.IsEmpty() {
		t.Error("Stack should be empty after all pops")
	}
	
	_, ok = stack.Pop()
	if ok {
		t.Error("Pop on empty stack should return false")
	}
}

func TestStack_String(t *testing.T) {
	stack := NewStack[string]()
	stack.Push("hello")
	stack.Push("world")
	
	val, ok := stack.Pop()
	if !ok || val != "world" {
		t.Errorf("Pop() = %q, want \"world\\", val)
	}
}

func TestQueue(t *testing.T) {
	queue := NewQueue[int]()
	
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	
	val, ok := queue.Dequeue()
	if !ok || val != 1 {
		t.Errorf("Dequeue() = %d, want 1", val)
	}
	
	val, ok = queue.Dequeue()
	if !ok || val != 2 {
		t.Errorf("Dequeue() = %d, want 2", val)
	}
	
	queue.Enqueue(4)
	
	val, ok = queue.Dequeue()
	if !ok || val != 3 {
		t.Errorf("Dequeue() = %d, want 3", val)
	}
	
	val, ok = queue.Dequeue()
	if !ok || val != 4 {
		t.Errorf("Dequeue() = %d, want 4", val)
	}
}

func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	doubled := Map(numbers, func(n int) int { return n * 2 })
	
	want := []int{2, 4, 6, 8}
	if len(doubled) != len(want) {
		t.Fatalf("Length = %d, want %d", len(doubled), len(want))
	}
	
	for i, v := range doubled {
		if v != want[i] {
			t.Errorf("doubled[%d] = %d, want %d", i, v, want[i])
		}
	}
}

func TestFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	
	want := []int{2, 4, 6}
	if len(evens) != len(want) {
		t.Fatalf("Length = %d, want %d", len(evens), len(want))
	}
	
	for i, v := range evens {
		if v != want[i] {
			t.Errorf("evens[%d] = %d, want %d", i, v, want[i])
		}
	}
}

func TestMax(t *testing.T) {
	if got := Max(1, 5, 3, 2); got != 5 {
		t.Errorf("Max(1,5,3,2) = %d, want 5", got)
	}
	
	if got := Max("apple", "zebra", "banana"); got != "zebra" {
		t.Errorf("Max(strings) = %q, want \"zebra\\", got)
	}
}
