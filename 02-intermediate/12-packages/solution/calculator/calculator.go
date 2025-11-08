package calculator

import "example.com/packages/calculator/internal/ops"

func Add(a, b int) int {
	return ops.Add(a, b)
}

func Subtract(a, b int) int {
	return ops.Subtract(a, b)
}
