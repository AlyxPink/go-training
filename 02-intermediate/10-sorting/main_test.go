package main

import (
	"sort"
	"testing"
)

func TestByAge(t *testing.T) {
	people := []Person{
		{"Bob", 30},
		{"Alice", 25},
		{"Charlie", 35},
	}
	
	sort.Sort(ByAge(people))
	
	if people[0].Age != 25 {
		t.Errorf("First person age = %d, want 25", people[0].Age)
	}
	if people[2].Age != 35 {
		t.Errorf("Last person age = %d, want 35", people[2].Age)
	}
}

func TestSortByName(t *testing.T) {
	people := []Person{
		{"Charlie", 30},
		{"Alice", 25},
		{"Bob", 35},
	}
	
	SortByName(people)
	
	if people[0].Name != "Alice" {
		t.Errorf("First person = %q, want Alice", people[0].Name)
	}
}

func TestSortByMultiple(t *testing.T) {
	people := []Person{
		{"Bob", 30},
		{"Alice", 30},
		{"Charlie", 25},
	}
	
	SortByMultiple(people)
	
	// Should be sorted by age, then name
	if people[0].Name != "Charlie" {
		t.Errorf("First = %q, want Charlie", people[0].Name)
	}
	if people[1].Name != "Alice" {
		t.Errorf("Second = %q, want Alice", people[1].Name)
	}
}
