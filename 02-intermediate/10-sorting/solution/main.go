package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func SortByName(people []Person) {
	sort.Slice(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})
}

func SortByMultiple(people []Person) {
	sort.Slice(people, func(i, j int) bool {
		if people[i].Age != people[j].Age {
			return people[i].Age < people[j].Age
		}
		return people[i].Name < people[j].Name
	})
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
	
	SortByMultiple(people)
	fmt.Println("Multi-field:", people)
}
