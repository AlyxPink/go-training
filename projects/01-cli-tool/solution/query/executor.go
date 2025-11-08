package query

import (
	"fmt"
)

func (f *FieldSelect) Execute(data interface{}) (interface{}, error) {
	// TODO: Implement field selection
	// Hint: Type assert to map[string]interface{}, return field value

	// Special case: handle .length on arrays
	if f.Field == "length" {
		switch v := data.(type) {
		case []interface{}:
			return len(v), nil
		case string:
			return len(v), nil
		}
	}

	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot select field from non-object (got %T)", data)
	}

	val, exists := m[f.Field]
	if !exists {
		return nil, nil
	}

	return val, nil
}

func (a *ArrayIndex) Execute(data interface{}) (interface{}, error) {
	// TODO: Implement array indexing
	// Hint: Type assert to []interface{}, check bounds, return element
	arr, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot index non-array (got %T)", data)
	}

	if a.Index < 0 || a.Index >= len(arr) {
		return nil, fmt.Errorf("array index out of bounds: %d (length: %d)", a.Index, len(arr))
	}

	return arr[a.Index], nil
}

func (a *ArrayIterate) Execute(data interface{}) (interface{}, error) {
	// TODO: Implement array iteration
	// Hint: Type assert to []interface{}, return as-is or process each element
	arr, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot iterate over non-array (got %T)", data)
	}

	return arr, nil
}

func (l *LengthOp) Execute(data interface{}) (interface{}, error) {
	// TODO: Implement length operation
	// Hint: Handle arrays, maps, strings
	switch v := data.(type) {
	case []interface{}:
		return len(v), nil
	case map[string]interface{}:
		return len(v), nil
	case string:
		return len(v), nil
	default:
		return nil, fmt.Errorf("cannot get length of %T", data)
	}
}

func (p *Pipe) Execute(data interface{}) (interface{}, error) {
	// TODO: Implement pipe operation
	// Hint: Execute left, then execute right on result
	result, err := p.Left.Execute(data)
	if err != nil {
		return nil, err
	}

	return p.Right.Execute(result)
}
