package main

import (
	"testing"
)

var testParts = []string{
	"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog",
}

func BenchmarkStringConcat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		StringConcat(testParts)
	}
}

func BenchmarkStringsBuilder(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		StringsBuilder(testParts)
	}
}

func BenchmarkBytesBuffer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BytesBuffer(testParts)
	}
}

// Map vs Slice benchmarks
var testMap = map[int]string{
	1: "one", 2: "two", 3: "three", 4: "four", 5: "five",
	6: "six", 7: "seven", 8: "eight", 9: "nine", 10: "ten",
}
var testKeys = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func BenchmarkMapLookup(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MapLookup(testMap, testKeys)
	}
}

var testSlice = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
var testTargets = []string{"one", "five", "ten"}

func BenchmarkSliceScan(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SliceScan(testSlice, testTargets)
	}
}

// Struct by value vs pointer
var testStruct = LargeStruct{ID: 1, Name: "test"}

func BenchmarkStructByValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ProcessStructByValue(testStruct)
	}
}

func BenchmarkStructByPointer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ProcessStructByPointer(&testStruct)
	}
}

// JSON marshaling benchmarks
var testUsers = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30},
	{ID: 2, Name: "Bob", Email: "bob@example.com", Age: 25},
	{ID: 3, Name: "Charlie", Email: "charlie@example.com", Age: 35},
}

func BenchmarkMarshalJSON(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MarshalJSON(testUsers)
	}
}

func BenchmarkMarshalJSONOptimized(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MarshalJSONOptimized(testUsers)
	}
}
