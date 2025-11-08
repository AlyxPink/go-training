package main
import "testing"

func TestSafeDivide(t *testing.T) {
	result, err := SafeDivide(10, 2)
	if err != nil || result != 5 {
		t.Error("SafeDivide(10,2) failed")
	}
	
	_, err = SafeDivide(10, 0)
	if err == nil {
		t.Error("SafeDivide(10,0) should error")
	}
}
