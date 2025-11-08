package main
import "testing"

func TestFizzBuzz(t *testing.T) {
	if FizzBuzz(15) != "FizzBuzz" {
		t.Error("FizzBuzz(15) failed")
	}
	if FizzBuzz(3) != "Fizz" {
		t.Error("FizzBuzz(3) failed")
	}
}
