package main
import "testing"

func TestStringToInt(t *testing.T) {
	val, err := StringToInt("123")
	if err != nil || val != 123 {
		t.Error("StringToInt failed")
	}
}

func TestIntToString(t *testing.T) {
	if IntToString(123) != "123" {
		t.Error("IntToString failed")
	}
}
