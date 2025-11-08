package main

import (
	"testing"
)

var testStrings = []string{"hello", "world", "golang", "performance", "optimization"}

func BenchmarkStringProcessor(b *testing.B) {
	p := &StringProcessor{}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		p.Process(testStrings)
	}
}

func BenchmarkOptimizedStringProcessor(b *testing.B) {
	p := NewOptimizedStringProcessor()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		p.Process(testStrings)
	}
}

var testData = [][]byte{
	[]byte("hello "),
	[]byte("world "),
	[]byte("from "),
	[]byte("golang "),
	[]byte("performance "),
}

func BenchmarkDataAggregator(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		DataAggregator(testData)
	}
}

func BenchmarkOptimizedDataAggregator(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		OptimizedDataAggregator(testData)
	}
}

func BenchmarkSliceGrower(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SliceGrower(1000)
	}
}

func BenchmarkOptimizedSliceGrower(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		OptimizedSliceGrower(1000)
	}
}

func BenchmarkProcessItems(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ProcessItems(100)
	}
}
