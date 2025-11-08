package main

import (
	"errors"
	"fmt"
)

// ErrNegative is a sentinel error for negative numbers.
var ErrNegative = errors.New("negative number not allowed")

// Sqrt returns the square root of a number or an error if negative.
func Sqrt(x float64) (float64, error) {
	// TODO: Return error if x < 0, otherwise return approximate sqrt
	return 0, nil
}

// CustomError is a custom error type.
type CustomError struct {
	Op  string
	Err error
}

func (e *CustomError) Error() string {
	// TODO: Return formatted error message
	return ""
}

func (e *CustomError) Unwrap() error {
	// TODO: Return wrapped error
	return nil
}

func main() {
	result, err := Sqrt(-4)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
