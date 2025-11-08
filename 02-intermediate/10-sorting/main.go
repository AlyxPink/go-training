package main

import (
	"fmt"
	"sort"
)

// Person represents a person
type Person struct {
	Name string
	Age  int
}

// ByAge implements sort.Interface for []Person by Age
type ByAge []Person

func (a ByAge) Len() int {
	// TODO: Return length
	return 0
}

func (a ByAge) Less(i, j int) bool {
	// TODO: Compare ages
	return false
}

func (a ByAge) Swap(i, j int) {
	// TODO: Swap elements
}

// SortByName sorts people by name
func SortByName(people []Person) {
	// TODO: Use sort.Slice
}

// SortByMultiple sorts by age then name
func SortByMultiple(people []Person) {
	// TODO: Multi-field sort
}

func main() {
	people := []Person{
		{"Bob", 30},
		{"Alice", 25},
		{"Charlie", 30},
	}
	
	sort.Sort(ByAge(people))
	fmt.Println("By age:", people)
	
	SortByName(people)
	fmt.Println("By name:", people)
}
