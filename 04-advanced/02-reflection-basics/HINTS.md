# Hints for Reflection Basics

## Hint 1: Getting Reflection Values

Use `reflect.ValueOf()` and `reflect.TypeOf()` to inspect values:

```go
v := reflect.ValueOf(data)
t := reflect.TypeOf(data)

// Handle pointers
if v.Kind() == reflect.Ptr {
    v = v.Elem()
    t = t.Elem()
}
```

## Hint 2: Iterating Struct Fields

Loop through struct fields and access their properties:

```go
for i := 0; i < v.NumField(); i++ {
    field := v.Field(i)
    fieldType := t.Field(i)

    name := fieldType.Name
    tag := fieldType.Tag.Get("validate")
    value := field.Interface()
}
```

## Hint 3: Parsing Struct Tags

Split tag values to extract validation rules:

```go
tag := "required,min=5,max=100"
rules := strings.Split(tag, ",")
for _, rule := range rules {
    if strings.Contains(rule, "=") {
        parts := strings.SplitN(rule, "=", 2)
        key, value := parts[0], parts[1]
    }
}
```

## Hint 4: Setting Struct Fields

Ensure field is settable before modification:

```go
if field.CanSet() {
    field.SetString("new value")
    field.SetInt(42)
}
```

## Hint 5: Type Conversion

Convert between reflect.Value and concrete types:

```go
// To interface{}
value := field.Interface()

// To specific type
if field.Kind() == reflect.String {
    str := field.String()
}

// Type assertion
if intVal, ok := value.(int); ok {
    // use intVal
}
```
