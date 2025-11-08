package query

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Query struct {
	nodes []QueryNode
}

type QueryNode interface {
	Execute(data interface{}) (interface{}, error)
}

// FieldSelect represents .field
type FieldSelect struct {
	Field string
}

// ArrayIndex represents [0]
type ArrayIndex struct {
	Index int
}

// ArrayIterate represents []
type ArrayIterate struct{}

// Pipe represents |
type Pipe struct {
	Left, Right QueryNode
}

// LengthOp represents length
type LengthOp struct{}

func Parse(queryStr string) (*Query, error) {
	queryStr = strings.TrimSpace(queryStr)
	if queryStr == "" {
		return nil, fmt.Errorf("empty query")
	}

	// TODO: Implement query parser
	// Hint: Tokenize, then build AST
	// For now, handle basic cases
	
	nodes, err := parseExpression(queryStr)
	if err != nil {
		return nil, err
	}

	return &Query{nodes: nodes}, nil
}

func parseExpression(expr string) ([]QueryNode, error) {
	// TODO: Implement full parser
	// This is a simplified version for basic cases
	
	var nodes []QueryNode
	i := 0

	for i < len(expr) {
		if expr[i] == '.' {
			i++
			// Parse field name
			end := i
			for end < len(expr) && (unicode.IsLetter(rune(expr[end])) || unicode.IsDigit(rune(expr[end])) || expr[end] == '_') {
				end++
			}
			if end > i {
				nodes = append(nodes, &FieldSelect{Field: expr[i:end]})
				i = end
			}
		} else if expr[i] == '[' {
			i++
			if i < len(expr) && expr[i] == ']' {
				// Array iterate
				nodes = append(nodes, &ArrayIterate{})
				i++
			} else {
				// Array index
				end := i
				for end < len(expr) && unicode.IsDigit(rune(expr[end])) {
					end++
				}
				if end > i {
					index, err := strconv.Atoi(expr[i:end])
					if err != nil {
						return nil, fmt.Errorf("invalid array index: %s", expr[i:end])
					}
					nodes = append(nodes, &ArrayIndex{Index: index})
					i = end
					if i < len(expr) && expr[i] == ']' {
						i++
					}
				}
			}
		} else if expr[i] == ' ' {
			i++
		} else {
			// Check for "length"
			if strings.HasPrefix(expr[i:], "length") {
				nodes = append(nodes, &LengthOp{})
				i += 6
			} else {
				return nil, fmt.Errorf("unexpected character at position %d: %c", i, expr[i])
			}
		}
	}

	return nodes, nil
}

func (q *Query) Execute(data interface{}) (interface{}, error) {
	result := data
	for _, node := range q.nodes {
		var err error
		result, err = node.Execute(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
