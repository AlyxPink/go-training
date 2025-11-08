package main

import (
	"fmt"
	"sync"
)

var (
	instance *Database
	once     sync.Once
)

type Database struct {
	connections int
}

func GetDatabase() *Database {
	once.Do(func() {
		fmt.Println("Creating database instance")
		instance = &Database{connections: 10}
	})
	return instance
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func GetBuffer() []byte {
	buf := bufferPool.Get().([]byte)
	return buf[:0] // Reset length
}

func PutBuffer(buf []byte) {
	if cap(buf) == 1024 {
		bufferPool.Put(buf)
	}
}

type Queue struct {
	items []int
	mu    sync.Mutex
	cond  *sync.Cond
}

func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Push(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *Queue) Pop() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait()
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func main() {
	fmt.Println("Advanced Sync Primitives")

	db1 := GetDatabase()
	db2 := GetDatabase()
	fmt.Printf("Same instance: %v\n", db1 == db2)

	buf := GetBuffer()
	fmt.Printf("Buffer size: %d\n", cap(buf))
	PutBuffer(buf)

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
