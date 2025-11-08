package main

import (
	"fmt"
	"sync"
)

// Singleton demonstrates sync.Once
var (
	instance *Database
	once     sync.Once
)

type Database struct {
	connections int
}

// GetDatabase returns singleton instance
// TODO: Use sync.Once for initialization
func GetDatabase() *Database {
	// TODO: Implement singleton pattern
	return nil
}

// BufferPool demonstrates sync.Pool
var bufferPool = sync.Pool{
	New: func() interface{} {
		// TODO: Return new buffer
		return nil
	},
}

// GetBuffer gets a buffer from pool
// TODO: Implement buffer pooling
func GetBuffer() []byte {
	// TODO: Get from pool
	// TODO: Reset buffer
	return nil
}

// PutBuffer returns buffer to pool
func PutBuffer(buf []byte) {
	// TODO: Clear buffer
	// TODO: Put back in pool
}

// Queue demonstrates sync.Cond
type Queue struct {
	items []int
	mu    sync.Mutex
	cond  *sync.Cond
}

// NewQueue creates a queue with condition variable
func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push adds item and signals waiters
// TODO: Implement with Cond
func (q *Queue) Push(item int) {
	// TODO: Lock
	// TODO: Add item
	// TODO: Signal waiting goroutine
	// TODO: Unlock
}

// Pop waits for and removes item
// TODO: Implement with Cond.Wait
func (q *Queue) Pop() int {
	// TODO: Lock
	// TODO: Wait while empty
	// TODO: Pop item
	// TODO: Unlock
	return 0
}

func main() {
	fmt.Println("Advanced Sync Primitives")

	// Singleton example
	db1 := GetDatabase()
	db2 := GetDatabase()
	fmt.Printf("Same instance: %v\n", db1 == db2)

	// Pool example
	buf := GetBuffer()
	fmt.Printf("Buffer size: %d\n", len(buf))
	PutBuffer(buf)

	// Cond example
	q := NewQueue()
	go func() {
		for i := 0; i < 5; i++ {
			q.Push(i)
		}
	}()

	for i := 0; i < 5; i++ {
		fmt.Printf("Popped: %d\n", q.Pop())
	}
}
