package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

type StringProcessor struct{}

func (p *StringProcessor) Process(items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = strings.ToUpper(item)
	}
	return result
}

type OptimizedStringProcessor struct {
	bufferPool *sync.Pool
}

func NewOptimizedStringProcessor() *OptimizedStringProcessor {
	return &OptimizedStringProcessor{
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (p *OptimizedStringProcessor) Process(items []string) []string {
	result := make([]string, len(items))

	for i, item := range items {
		buf := p.bufferPool.Get().(*bytes.Buffer)
		buf.Reset()

		buf.WriteString(strings.ToUpper(item))
		result[i] = buf.String()

		p.bufferPool.Put(buf)
	}

	return result
}

func DataAggregator(data [][]byte) []byte {
	var result []byte
	for _, chunk := range data {
		result = append(result, chunk...)
	}
	return result
}

func OptimizedDataAggregator(data [][]byte) []byte {
	// Calculate total size
	totalSize := 0
	for _, chunk := range data {
		totalSize += len(chunk)
	}

	// Preallocate
	result := make([]byte, 0, totalSize)
	for _, chunk := range data {
		result = append(result, chunk...)
	}

	return result
}

func SliceGrower(n int) []int {
	var result []int
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return result
}

func OptimizedSliceGrower(n int) []int {
	result := make([]int, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return result
}

type Item struct {
	ID    int
	Data  [1024]byte
	Value string
}

var itemPool = sync.Pool{
	New: func() interface{} {
		return &Item{}
	},
}

func GetItem() *Item {
	item := itemPool.Get().(*Item)
	// Reset fields
	item.ID = 0
	item.Value = ""
	for i := range item.Data {
		item.Data[i] = 0
	}
	return item
}

func PutItem(item *Item) {
	// Clear sensitive data
	item.Value = ""
	itemPool.Put(item)
}

func ProcessItems(n int) {
	for i := 0; i < n; i++ {
		item := GetItem()
		item.ID = i
		item.Value = fmt.Sprintf("item-%d", i)

		// Process item...

		PutItem(item)
	}
}

func main() {
	items := []string{"hello", "world", "golang", "performance"}

	processor := &StringProcessor{}
	result1 := processor.Process(items)
	fmt.Printf("Basic processor: %v\n", result1)

	optProcessor := NewOptimizedStringProcessor()
	result2 := optProcessor.Process(items)
	fmt.Printf("Optimized processor: %v\n", result2)

	data := [][]byte{
		[]byte("hello "),
		[]byte("world "),
		[]byte("from "),
		[]byte("golang"),
	}

	agg1 := DataAggregator(data)
	fmt.Printf("Basic aggregator: %s\n", agg1)

	agg2 := OptimizedDataAggregator(data)
	fmt.Printf("Optimized aggregator: %s\n", agg2)

	fmt.Println("Growing slices...")
	_ = SliceGrower(1000)
	_ = OptimizedSliceGrower(1000)

	fmt.Println("Processing items with pool...")
	ProcessItems(100)
}
