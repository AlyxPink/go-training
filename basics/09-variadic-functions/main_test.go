package main
import "testing"

func TestSum(t *testing.T) {
	if Sum(1, 2, 3) != 6 {
		t.Error("Sum failed")
	}
	if Sum() != 0 {
		t.Error("Sum() should be 0")
	}
}
