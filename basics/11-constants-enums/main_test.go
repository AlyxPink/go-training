package main
import "testing"

func TestStatus(t *testing.T) {
	if StatusPending != 0 {
		t.Error("StatusPending should be 0")
	}
	if StatusActive != 1 {
		t.Error("StatusActive should be 1")
	}
}
