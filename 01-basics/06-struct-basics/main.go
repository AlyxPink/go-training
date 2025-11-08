package main

import "fmt"

// Person represents a person with basic information.
type Person struct {
	// TODO: Define fields: Name (string), Age (int), Email (string)
}

// NewPerson creates and returns a new Person.
func NewPerson(name string, age int, email string) *Person {
	// TODO: Return properly initialized Person
	return nil
}

// String returns a formatted string representation of the Person.
func (p Person) String() string {
	// TODO: Return formatted string like "Person{Name: Alice, Age: 30, Email: alice@example.com}"
	return ""
}

// IsAdult returns true if the person is 18 or older.
func (p Person) IsAdult() bool {
	// TODO: Check if Age >= 18
	return false
}

// Birthday increments the person's age by 1.
func (p *Person) Birthday() {
	// TODO: Increment age (note: pointer receiver needed to modify)
}

// Student embeds Person and adds academic information.
type Student struct {
	// TODO: Embed Person and add GPA field (float64)
}

// String returns a formatted string representation of the Student.
func (s Student) String() string {
	// TODO: Return formatted string with student info including GPA
	return ""
}

// IsHonorStudent returns true if GPA is 3.5 or higher.
func (s Student) IsHonorStudent() bool {
	// TODO: Check if GPA >= 3.5
	return false
}

func main() {
	// Test your implementations
	p := NewPerson("Alice", 25, "alice@example.com")
	fmt.Println(p)
	fmt.Println("Is adult:", p.IsAdult())

	// TODO: Uncomment after implementing Person and Student structs
	// s := Student{Person: Person{Name: "Bob", Age: 20, Email: "bob@university.edu"}, GPA: 3.8}
	// fmt.Println(s)
	// fmt.Println("Honor student:", s.IsHonorStudent())
}
