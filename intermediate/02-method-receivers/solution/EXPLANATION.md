# Solution Explanation: Method Receivers

## Receiver Selection Decision Matrix

| Type | Methods | Receiver Choice | Reason |
|------|---------|----------------|---------|
| Counter | All methods | Pointer (`*Counter`) | Mutates state |
| Point | Distance | Value (`Point`) | Read-only, small type |
| Point | Translate, String | Pointer (`*Point`) | Translate mutates, String for consistency |
| Configuration | All methods | Pointer (`*Configuration`) | Large struct, avoid copying |
| Temperature | ToFahrenheit, IsFreezing | Value (`Temperature`) | Read-only operations on small type |
| Temperature | Warm | Pointer (`*Temperature`) | Mutates value |

## Design Decisions

### Counter - Pointer Receivers

```go
func (c *Counter) Increment() {
	c.count++  // Modifies receiver
}
```

**Why pointer?**
- Methods modify internal state
- Consistency: all methods use pointer receivers
- Caller sees mutations: `c.Increment()` actually changes `c`

**Common mistake:**
```go
func (c Counter) Increment() {  // Wrong!
	c.count++  // Modifies a copy, caller won't see change
}
```

### Point - Mixed Approach (But Consistent Pointer)

```go
func (p Point) Distance(other Point) float64 {  // Could be value
	// Read-only operation
}

func (p *Point) Translate(dx, dy int) {  // Must be pointer
	p.X += dx  // Modifies receiver
}
```

**Why pointer for all?**
- Translate requires pointer (mutates)
- Distance could be value, but pointer for consistency
- Go idiom: use same receiver type for all methods on a type

**Alternative (not recommended):**
```go
// Mixing receivers is confusing
func (p Point) Distance(other Point) float64   // value
func (p *Point) Translate(dx, dy int)          // pointer
// This works but breaks consistency guideline
```

### Configuration - Always Pointer

```go
type Configuration struct {
	Host    string
	Port    int
	Timeout int
	Debug   bool
	// Imagine 20 more fields...
}

func (c *Configuration) Validate() bool {
	// Even though read-only, use pointer to avoid copying large struct
}
```

**Why pointer for read-only methods?**
- Large structs are expensive to copy
- Passing by value copies all fields
- Pointer is 8 bytes regardless of struct size

**Performance impact:**
```go
// Value receiver: copies entire struct on every call
func (c Configuration) Validate() bool  // Copies ~100 bytes

// Pointer receiver: copies only pointer
func (c *Configuration) Validate() bool  // Copies 8 bytes
```

### Temperature - Mixed Receivers (Intentional)

```go
type Temperature int  // Just an int, very small

func (t Temperature) ToFahrenheit() float64 {
	return float64(t)*9/5 + 32  // No mutation
}

func (t *Temperature) Warm(degrees int) {
	*t += Temperature(degrees)  // Mutation requires pointer
}
```

**Why mix receivers here?**
- Base type is small (int = 8 bytes)
- Read-only methods can use value (no copy overhead)
- Mutating methods must use pointer
- This is an exception to consistency rule for small, simple types

## Go Idioms Demonstrated

### 1. Receiver Type Consistency

**Guideline:** Use the same receiver type for all methods on a type

```go
// Good: consistent pointer receivers
type Person struct { name string }
func (p *Person) SetName(name string) { p.name = name }
func (p *Person) GetName() string     { return p.name }

// Less ideal: mixed receivers (though it works)
func (p *Person) SetName(name string) { p.name = name }
func (p Person) GetName() string      { return p.name }
```

### 2. When to Break Consistency

**Exception:** Small, primitive-like types (like Temperature)

```go
type Age int

func (a Age) String() string    { return fmt.Sprintf("%d years", a) }
func (a *Age) Increment()       { *a++ }
```

This is acceptable because Age is essentially an int.

### 3. Method Sets and Interfaces

**Important rule:**
- Type `T` has methods with receiver `T`
- Type `*T` has methods with receiver `T` or `*T`

```go
type Incrementer interface {
	Increment()
}

var _ Incrementer = (*Counter)(nil)  // ✓ *Counter implements Incrementer
var _ Incrementer = Counter{}         // ✗ Counter does NOT implement Incrementer
```

### 4. Pointer Receiver Necessity

**Always use pointer when:**
1. Method modifies receiver
2. Receiver is large struct (avoid copy overhead)
3. Consistency with other pointer-receiver methods
4. Need to implement interface that requires pointer receiver

**Can use value when:**
1. Method is read-only
2. Type is small (int, small struct)
3. Want immutability semantics

## Performance Considerations

### Copy Overhead

```go
// Assume 100-byte struct
type Large struct { /* many fields */ }

// Every call copies 100 bytes
func (l Large) Process() {}

// Every call copies 8 bytes (pointer)
func (l *Large) Process() {}
```

### Method Calls

```go
p := Point{X: 1, Y: 2}
p.Distance(...)  // If Distance has value receiver: no allocation
                  // If Distance has pointer receiver: compiler may allocate

// To call pointer-receiver method on value:
// Compiler automatically takes address: (&p).Distance(...)
```

## Common Mistakes

### 1. Forgetting to Mutate

```go
func (t Temperature) Warm(degrees int) {  // Wrong!
	t += Temperature(degrees)  // Modifies copy, not original
}
```

### 2. Copying Large Structs Unnecessarily

```go
func (c Configuration) ApplyDefaults() {  // Wasteful
	// Copies entire Configuration struct
}
```

### 3. Inconsistent Receivers

```go
type User struct { name string }
func (u *User) SetName(n string) { u.name = n }  // pointer
func (u User) GetName() string   { return u.name }  // value
// Works but confusing and inconsistent
```

## Testing Insights

**Verify mutation:**
```go
c := Counter{}
c.Increment()
// Must check c.Value(), not return value of Increment()
```

**Test both pointer and value:**
```go
// These should both work if receiver is pointer
p := &Point{X: 1, Y: 2}
p.Translate(1, 1)

p2 := Point{X: 1, Y: 2}
p2.Translate(1, 1)  // Compiler takes address automatically
```

## Further Exploration

- Study method sets impact on interface satisfaction
- Explore escape analysis: when do values escape to heap?
- Learn about `go vet` checks for receiver consistency
- Understand pointer receiver implications for concurrency
