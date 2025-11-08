package main

import (
	"fmt"
	"sync/atomic"
)

type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Decrement() {
	atomic.AddInt64(&c.value, -1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func main() {
	fmt.Println("Atomic operations")
	c := &AtomicCounter{}
	c.Increment()
	fmt.Println("Counter:", c.Value())
}
