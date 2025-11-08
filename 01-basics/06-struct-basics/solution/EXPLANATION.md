# Solution Explanation: Struct Basics

## Key Go Concepts Demonstrated

### 1. Struct Definition

Structs are typed collections of fields:

```go
type Person struct {
    Name  string
    Age   int
    Email string
}
```

**Key points**:
- Fields are named
- Each field has a type
- Exported fields start with capital letter
- Unexported (private) fields start with lowercase

### 2. Struct Initialization

Multiple ways to create structs:

```go
// Zero value
var p Person  // All fields get zero values

// Struct literal (named fields - recommended)
p := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

// Struct literal (positional - fragile)
p := Person{"Alice", 30, "alice@example.com"}

// Partial initialization
p := Person{Name: "Alice"}  // Age: 0, Email: ""

// Pointer to struct
p := &Person{Name: "Alice", Age: 30}
```

### 3. Methods and Receivers

Methods are functions with a receiver:

```go
// Value receiver - receives a copy
func (p Person) IsAdult() bool {
    return p.Age >= 18
}

// Pointer receiver - receives a pointer
func (p *Person) Birthday() {
    p.Age++  // Modifies the original
}
```

**When to use pointer receivers**:
- Method needs to modify the receiver
- Struct is large (avoid copying)
- Consistency (if any method uses pointer, all should)

**When to use value receivers**:
- Method doesn't modify the receiver
- Struct is small
- Immutability is desired

### 4. Constructor Pattern

Go doesn't have constructors, but we use factory functions:

```go
func NewPerson(name string, age int, email string) *Person {
    return &Person{
        Name:  name,
        Age:   age,
        Email: email,
    }
}
```

**Advantages**:
- Validation logic
- Default values
- Private fields initialization
- Returns pointer (common pattern)

### 5. fmt.Stringer Interface

Implementing `String()` makes your type printable:

```go
func (p Person) String() string {
    return fmt.Sprintf("Person{Name: %s, Age: %d}", p.Name, p.Age)
}

p := Person{Name: "Alice", Age: 30}
fmt.Println(p)  // Automatically calls String()
```

### 6. Struct Embedding (Composition)

Go uses composition instead of inheritance:

```go
type Student struct {
    Person  // Embedded field
    GPA float64
}

s := Student{
    Person: Person{Name: "Bob", Age: 20},
    GPA: 3.8,
}

// Can access embedded fields directly
fmt.Println(s.Name)     // Same as s.Person.Name
fmt.Println(s.IsAdult()) // Can call embedded methods
```

**How embedding works**:
- Embedded fields are promoted to the outer struct
- Can access as if they were fields of outer struct
- Can still access via explicit path: `s.Person.Name`
- Methods of embedded type are promoted too

### 7. Method Overriding

Outer type methods shadow embedded type methods:

```go
// Person's String method
func (p Person) String() string {
    return fmt.Sprintf("Person{...}")
}

// Student's String method overrides Person's
func (s Student) String() string {
    return fmt.Sprintf("Student{...}")
}

s := Student{...}
fmt.Println(s)  // Calls Student.String(), not Person.String()
```

## Common Patterns

### 1. Constructor with Validation

```go
func NewPerson(name string, age int, email string) (*Person, error) {
    if age < 0 {
        return nil, errors.New("age cannot be negative")
    }
    if name == "" {
        return nil, errors.New("name is required")
    }
    return &Person{Name: name, Age: age, Email: email}, nil
}
```

### 2. Builder Pattern

```go
type PersonBuilder struct {
    person Person
}

func (b *PersonBuilder) WithName(name string) *PersonBuilder {
    b.person.Name = name
    return b
}

func (b *PersonBuilder) WithAge(age int) *PersonBuilder {
    b.person.Age = age
    return b
}

func (b *PersonBuilder) Build() Person {
    return b.person
}

// Usage
p := (&PersonBuilder{}).
    WithName("Alice").
    WithAge(30).
    Build()
```

### 3. Option Pattern

```go
type PersonOption func(*Person)

func WithEmail(email string) PersonOption {
    return func(p *Person) {
        p.Email = email
    }
}

func NewPerson(name string, age int, opts ...PersonOption) *Person {
    p := &Person{Name: name, Age: age}
    for _, opt := range opts {
        opt(p)
    }
    return p
}

// Usage
p := NewPerson("Alice", 30, WithEmail("alice@example.com"))
```

## Common Pitfalls

1. **Forgetting pointer receiver**:
```go
// WRONG - doesn't modify original
func (p Person) Birthday() {
    p.Age++  // Modifies copy
}

// RIGHT
func (p *Person) Birthday() {
    p.Age++  // Modifies original
}
```

2. **Mixing receiver types**:
```go
// Inconsistent - avoid this
func (p Person) IsAdult() bool { ... }
func (p *Person) String() string { ... }

// Better - use pointer consistently
func (p *Person) IsAdult() bool { ... }
func (p *Person) String() string { ... }
```

3. **Embedded field conflicts**:
```go
type A struct { X int }
type B struct { X int }
type C struct {
    A
    B
}

c := C{}
// c.X  // Ambiguous! Compile error
c.A.X  // Must be explicit
```

4. **Comparing structs with uncomparable fields**:
```go
type Person struct {
    Name string
    Tags []string  // Slice is not comparable
}

p1 := Person{Name: "Alice"}
p2 := Person{Name: "Alice"}
// p1 == p2  // Compile error: slice is not comparable
```

## Go Idioms Used

1. **Constructor pattern**: `NewType()` functions
2. **Pointer receivers for mutation**: Methods that modify use `*T`
3. **Value receivers for queries**: Methods that read use `T`
4. **Struct embedding over inheritance**: Composition pattern
5. **fmt.Stringer interface**: `String()` method for custom printing
6. **Named struct literals**: `Person{Name: "Alice", Age: 30}`

## Performance Considerations

**Value vs Pointer**:
```go
// Small struct - value receiver OK
type Point struct { X, Y int }
func (p Point) String() string { ... }

// Large struct - pointer receiver better
type LargeStruct struct { data [1000]int }
func (l *LargeStruct) Process() { ... }  // Avoid copying
```

**Zero Values**:
- Structs are initialized to zero values (no null pointers to fields)
- Design types where zero value is useful when possible
- Example: `bytes.Buffer` zero value is ready to use
