package main
import "testing"

func TestSumRecursive(t *testing.T) {
	if SumRecursive(5) != 15 {
		t.Error("SumRecursive(5) should be 15")
	}
	if SumRecursive(0) != 0 {
		t.Error("SumRecursive(0) should be 0")
	}
}
