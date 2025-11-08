package main

import (
	"fmt"
	"sync"
)

// Counter with mutex protection
type Counter struct {
	mu    sync.Mutex
	value int
}

// TODO: Implement Increment, Decrement, Value methods

// ReadWriteCounter with RWMutex
type ReadWriteCounter struct {
	mu    sync.RWMutex
	value int
}

// TODO: Implement methods using RWMutex

func main() {
	fmt.Println("Mutex examples")
}
