# Hints: CGO Integration

## CGO Basics

### Import Statement
```go
/*
#include <stdlib.h>
#include <string.h>

int add(int a, int b) {
    return a + b;
}
*/
import "C"
```

### Type Conversions
- Go int ↔ C.int
- Go string → C.CString (must free!)
- C char* → Go string with C.GoString()
- Go []byte ↔ C arrays

### Memory Management
```go
// Always free C memory
cstr := C.CString("hello")
defer C.free(unsafe.Pointer(cstr))
```

## Common Patterns

### String Conversion
```go
// Go to C
goStr := "hello"
cStr := C.CString(goStr)
defer C.free(unsafe.Pointer(cStr))

// C to Go
cResult := C.some_c_function()
goResult := C.GoString(cResult)
```

### Array Conversion
```go
// Go slice to C array
goSlice := []int32{1, 2, 3}
cArray := (*C.int)(unsafe.Pointer(&goSlice[0]))

// C array to Go slice
cArr := C.malloc(C.size_t(len) * C.size_t(unsafe.Sizeof(C.int(0))))
defer C.free(cArr)
```

### Struct Handling
```go
type GoStruct struct {
    Field1 int32
    Field2 float64
}

// Pass to C
cStruct := (*C.struct_name)(unsafe.Pointer(&goStruct))
```

## Important Notes

1. **Thread Safety**: C code runs outside Go's runtime
2. **Performance**: CGO calls have overhead (~50ns)
3. **Build Tags**: Use `// +build cgo` for CGO-specific code
4. **Cross Compilation**: CGO makes cross-compilation harder
5. **Memory**: Go's GC doesn't track C memory

## Error Handling

```go
result := C.some_function()
if result == C.NULL {
    return errors.New("C function failed")
}
```

## Build Commands

```bash
# Normal build
go build

# Disable CGO
CGO_ENABLED=0 go build

# Cross compile (more complex with CGO)
GOOS=linux GOARCH=amd64 go build
```
