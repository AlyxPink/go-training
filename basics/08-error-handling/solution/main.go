package main

import (
	"errors"
	"fmt"
	"math"
)

var ErrNegative = errors.New("negative number not allowed")

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegative
	}
	return math.Sqrt(x), nil
}

type CustomError struct {
	Op  string
	Err error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func main() {
	result, err := Sqrt(16)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sqrt(16):", result)
	}
	
	result, err = Sqrt(-4)
	if err != nil {
		fmt.Println("Error:", err)
		if errors.Is(err, ErrNegative) {
			fmt.Println("  Detected sentinel error")
		}
	}
	
	customErr := &CustomError{Op: "divide", Err: errors.New("division by zero")}
	fmt.Println("Custom error:", customErr)
}
