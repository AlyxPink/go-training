package main

import (
	"strings"
	"testing"
)

func TestNewPerson(t *testing.T) {
	p := NewPerson("Alice", 30, "alice@example.com")
	
	if p == nil {
		t.Fatal("NewPerson returned nil")
	}
	
	if p.Name != "Alice" {
		t.Errorf("Name = %s, expected Alice", p.Name)
	}
	
	if p.Age != 30 {
		t.Errorf("Age = %d, expected 30", p.Age)
	}
	
	if p.Email != "alice@example.com" {
		t.Errorf("Email = %s, expected alice@example.com", p.Email)
	}
}

func TestPersonString(t *testing.T) {
	p := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}
	s := p.String()
	
	if !strings.Contains(s, "Alice") || !strings.Contains(s, "30") || !strings.Contains(s, "alice@example.com") {
		t.Errorf("String() = %s, should contain Name, Age, and Email", s)
	}
}

func TestIsAdult(t *testing.T) {
	tests := []struct {
		name     string
		age      int
		expected bool
	}{
		{"adult", 18, true},
		{"young adult", 25, true},
		{"minor", 17, false},
		{"child", 10, false},
		{"exactly 18", 18, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Person{Age: tt.age}
			result := p.IsAdult()
			if result != tt.expected {
				t.Errorf("IsAdult() with age %d = %v, expected %v", tt.age, result, tt.expected)
			}
		})
	}
}

func TestBirthday(t *testing.T) {
	p := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}
	p.Birthday()
	
	if p.Age != 26 {
		t.Errorf("After Birthday(), Age = %d, expected 26", p.Age)
	}
	
	// Test multiple birthdays
	p.Birthday()
	p.Birthday()
	
	if p.Age != 28 {
		t.Errorf("After 3 birthdays, Age = %d, expected 28", p.Age)
	}
}

func TestStudentEmbedding(t *testing.T) {
	s := Student{
		Person: Person{Name: "Bob", Age: 20, Email: "bob@university.edu"},
		GPA:    3.8,
	}
	
	// Test that we can access embedded Person fields
	if s.Name != "Bob" {
		t.Errorf("Student.Name = %s, expected Bob", s.Name)
	}
	
	if s.Age != 20 {
		t.Errorf("Student.Age = %d, expected 20", s.Age)
	}
	
	// Test that we can call embedded Person methods
	if !s.IsAdult() {
		t.Error("Student with Age 20 should be adult")
	}
}

func TestStudentString(t *testing.T) {
	s := Student{
		Person: Person{Name: "Bob", Age: 20, Email: "bob@university.edu"},
		GPA:    3.8,
	}
	
	str := s.String()
	
	if !strings.Contains(str, "Bob") || !strings.Contains(str, "3.8") {
		t.Errorf("Student String() = %s, should contain name and GPA", str)
	}
}

func TestIsHonorStudent(t *testing.T) {
	tests := []struct {
		name     string
		gpa      float64
		expected bool
	}{
		{"high gpa", 4.0, true},
		{"exactly 3.5", 3.5, true},
		{"slightly below", 3.4, false},
		{"low gpa", 2.0, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Student{GPA: tt.gpa}
			result := s.IsHonorStudent()
			if result != tt.expected {
				t.Errorf("IsHonorStudent() with GPA %.1f = %v, expected %v", tt.gpa, result, tt.expected)
			}
		})
	}
}
