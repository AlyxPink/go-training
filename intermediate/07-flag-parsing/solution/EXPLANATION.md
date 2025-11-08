# Flag Parsing

## Basic Pattern
```go
var name = flag.String("name", "default", "help")
flag.Parse()
fmt.Println(*name)  // Note: pointer
```

## Or with Var
```go
var s string
flag.StringVar(&s, "name", "default", "help")
flag.Parse()
fmt.Println(s)  // No pointer
```

## Custom Types
Implement flag.Value interface:
```go
type MyFlag struct { val string }
func (f *MyFlag) String() string { return f.val }
func (f *MyFlag) Set(s string) error { f.val = s; return nil }
```
