# Hints for Code Generation

## Hint 1: Parsing Go Files
```go
fset := token.NewFileSet()
node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
```

## Hint 2: Walking the AST
```go
ast.Inspect(node, func(n ast.Node) bool {
    switch x := n.(type) {
    case *ast.GenDecl:
        // Handle type declarations
    case *ast.FuncDecl:
        // Handle function declarations
    }
    return true
})
```

## Hint 3: Template Generation
```go
tmpl := template.Must(template.New("builder").Parse(builderTemplate))
tmpl.Execute(out, data)
```

## Hint 4: Go Generate
Add to source file:
```go
//go:generate go run generator.go -input=models.go -output=generated.go
```

## Hint 5: Type Information
```go
typeSpec := spec.(*ast.TypeSpec)
structType := typeSpec.Type.(*ast.StructType)
for _, field := range structType.Fields.List {
    fieldName := field.Names[0].Name
    fieldType := field.Type
}
```
