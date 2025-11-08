# Hints

```go
var name = flag.String("name", "default", "description")
flag.Parse()
fmt.Println(*name)
```

Custom types implement flag.Value interface.
