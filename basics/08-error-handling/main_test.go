package main

import (
	"errors"
	"testing"
)

func TestSqrt(t *testing.T) {
	result, err := Sqrt(16)
	if err != nil {
		t.Errorf("Sqrt(16) returned error: %v", err)
	}
	if result < 3.9 || result > 4.1 {
		t.Errorf("Sqrt(16) = %f, expected ~4", result)
	}
	
	_, err = Sqrt(-4)
	if err == nil {
		t.Error("Sqrt(-4) should return error")
	}
	if !errors.Is(err, ErrNegative) {
		t.Error("Sqrt(-4) should return ErrNegative")
	}
}

func TestCustomError(t *testing.T) {
	baseErr := errors.New("base error")
	customErr := &CustomError{Op: "test", Err: baseErr}
	
	if customErr.Error() == "" {
		t.Error("CustomError.Error() should not be empty")
	}
	
	if errors.Unwrap(customErr) != baseErr {
		t.Error("CustomError.Unwrap() should return base error")
	}
}
