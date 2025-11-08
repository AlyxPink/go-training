package main

import (
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", -5, 3, -2},
		{"zero", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestToUppercase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"simple string", "hello", "HELLO", false},
		{"mixed case", "Hello World", "HELLO WORLD", false},
		{"already uppercase", "HELLO", "HELLO", false},
		{"with numbers", "hello123", "HELLO123", false},
		{"empty string", "", "", false},
		{"special chars", "hello!@#", "HELLO!@#", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUppercase(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUppercase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUppercase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewDataPoint(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		value   float64
		dpName  string
		wantErr bool
	}{
		{"valid datapoint", 1, 42.5, "test", false},
		{"zero values", 0, 0.0, "", false},
		{"negative id", -1, 100.0, "negative", false},
		{"long name", 1, 1.0, strings.Repeat("a", 100), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDataPoint(tt.id, tt.value, tt.dpName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDataPoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.ID != tt.id {
					t.Errorf("DataPoint.ID = %d, want %d", got.ID, tt.id)
				}
				if got.Value != tt.value {
					t.Errorf("DataPoint.Value = %f, want %f", got.Value, tt.value)
				}
				// Name might be truncated if too long
				expectedName := tt.dpName
				if len(expectedName) > 49 {
					expectedName = expectedName[:49]
				}
				if got.Name != expectedName {
					t.Errorf("DataPoint.Name = %q, want %q", got.Name, expectedName)
				}
			}
		})
	}
}

func TestSumArray(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int
		want    int
	}{
		{"simple array", []int{1, 2, 3, 4, 5}, 15},
		{"empty array", []int{}, 0},
		{"single element", []int{42}, 42},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"mixed signs", []int{-5, 10, -3, 8}, 10},
		{"zeros", []int{0, 0, 0}, 0},
		{"large array", []int{100, 200, 300, 400, 500}, 1500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SumArray(tt.numbers)
			if got != tt.want {
				t.Errorf("SumArray(%v) = %d, want %d", tt.numbers, got, tt.want)
			}
		})
	}
}

func TestStringProcessor(t *testing.T) {
	processor := NewStringProcessor()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple string", "hello world", false},
		{"empty string", "", false},
		{"special characters", "test!@#$%", false},
		{"unicode", "hello 世界", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processor.Process(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == "" && tt.input != "" {
				t.Error("StringProcessor.Process() returned empty string for non-empty input")
			}
		})
	}
}

func TestCGOMemoryManagement(t *testing.T) {
	// Test that multiple conversions don't leak memory
	// This is a basic sanity test - use valgrind for real leak detection
	t.Run("multiple conversions", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			_, err := ToUppercase("test string")
			if err != nil {
				t.Fatalf("ToUppercase failed on iteration %d: %v", i, err)
			}
		}
	})

	t.Run("multiple datapoints", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			_, err := NewDataPoint(i, float64(i), "test")
			if err != nil {
				t.Fatalf("NewDataPoint failed on iteration %d: %v", i, err)
			}
		}
	})

	t.Run("multiple array sums", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		for i := 0; i < 1000; i++ {
			sum := SumArray(numbers)
			if sum != 15 {
				t.Fatalf("SumArray returned wrong result on iteration %d: %d", i, sum)
			}
		}
	})
}

// Benchmark tests
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

func BenchmarkToUppercase(b *testing.B) {
	input := "hello world this is a test string"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToUppercase(input)
	}
}

func BenchmarkSumArray(b *testing.B) {
	numbers := make([]int, 100)
	for i := range numbers {
		numbers[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumArray(numbers)
	}
}

func BenchmarkNewDataPoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDataPoint(i, float64(i), "benchmark test")
	}
}
