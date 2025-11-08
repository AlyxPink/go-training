package mathutil

import "testing"

func TestAdd(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Error("2 + 3 should be 5")
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{5, 3, 5},
		{3, 5, 5},
		{5, 5, 5},
	}

	for _, tt := range tests {
		got := Max(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
