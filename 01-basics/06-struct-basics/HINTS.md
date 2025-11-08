# Hints for Struct Basics

## Level 1: Struct Definition

```go
type Person struct {
    Name  string
    Age   int
    Email string
}
```

## Level 2: Initialization

```go
// Struct literal
p := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

// Zero value
var p Person  // Name: "", Age: 0, Email: ""

// Constructor pattern
func NewPerson(name string, age int, email string) *Person {
    return &Person{Name: name, Age: age, Email: email}
}
```

## Level 3: Methods

```go
// Value receiver - cannot modify struct
func (p Person) IsAdult() bool {
    return p.Age >= 18
}

// Pointer receiver - can modify struct
func (p *Person) Birthday() {
    p.Age++
}
```

## Level 4: Embedding

```go
type Student struct {
    Person  // Embedded struct
    GPA float64
}

// Access embedded fields
s := Student{Person: Person{Name: "Bob"}, GPA: 3.8}
fmt.Println(s.Name)  // Can access Person.Name directly
```

## Level 5: String Method

Implement `String()` method to make your struct printable:

```go
func (p Person) String() string {
    return fmt.Sprintf("Person{Name: %s, Age: %d, Email: %s}", 
        p.Name, p.Age, p.Email)
}
```
