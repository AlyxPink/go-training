package main

import "fmt"

// Person represents a person with basic information.
type Person struct {
	Name  string
	Age   int
	Email string
}

// NewPerson creates and returns a new Person.
// This is the constructor pattern in Go.
func NewPerson(name string, age int, email string) *Person {
	return &Person{
		Name:  name,
		Age:   age,
		Email: email,
	}
}

// String returns a formatted string representation of the Person.
// This implements the fmt.Stringer interface.
func (p Person) String() string {
	return fmt.Sprintf("Person{Name: %s, Age: %d, Email: %s}", p.Name, p.Age, p.Email)
}

// IsAdult returns true if the person is 18 or older.
// Uses value receiver since we don't need to modify the struct.
func (p Person) IsAdult() bool {
	return p.Age >= 18
}

// Birthday increments the person's age by 1.
// Uses pointer receiver to modify the struct.
func (p *Person) Birthday() {
	p.Age++
}

// Student embeds Person and adds academic information.
type Student struct {
	Person // Embedded struct - Student "has-a" Person
	GPA    float64
}

// String returns a formatted string representation of the Student.
// This overrides the embedded Person's String method.
func (s Student) String() string {
	return fmt.Sprintf("Student{Name: %s, Age: %d, Email: %s, GPA: %.2f}", 
		s.Name, s.Age, s.Email, s.GPA)
}

// IsHonorStudent returns true if GPA is 3.5 or higher.
func (s Student) IsHonorStudent() bool {
	return s.GPA >= 3.5
}

func main() {
	// Test Person
	fmt.Println("=== Person ===")
	p := NewPerson("Alice", 25, "alice@example.com")
	fmt.Println(p)
	fmt.Println("Is adult:", p.IsAdult())
	
	fmt.Println("\nBefore birthday:", p)
	p.Birthday()
	fmt.Println("After birthday:", p)
	
	// Test Student
	fmt.Println("\n=== Student ===")
	s := Student{
		Person: Person{Name: "Bob", Age: 20, Email: "bob@university.edu"},
		GPA:    3.8,
	}
	fmt.Println(s)
	fmt.Println("Honor student:", s.IsHonorStudent())
	
	// Accessing embedded fields
	fmt.Println("\nAccessing embedded Person fields:")
	fmt.Println("Name:", s.Name)
	fmt.Println("Is adult:", s.IsAdult())
	
	// Modifying through embedded methods
	s.Birthday()
	fmt.Println("After birthday:", s)
}
