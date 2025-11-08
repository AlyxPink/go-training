package main
import "testing"

func TestCounter(t *testing.T) {
	c := Counter()
	if c() != 1 {
		t.Error("First call should be 1")
	}
	if c() != 2 {
		t.Error("Second call should be 2")
	}
}
