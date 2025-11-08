package main

import (
	"math"
	"testing"
)

func TestCounter(t *testing.T) {
	t.Run("increment", func(t *testing.T) {
		c := Counter{}
		c.Increment()
		if got := c.Value(); got != 1 {
			t.Errorf("Value() = %d, want 1", got)
		}
		c.Increment()
		if got := c.Value(); got != 2 {
			t.Errorf("Value() = %d, want 2", got)
		}
	})

	t.Run("reset", func(t *testing.T) {
		c := Counter{count: 10}
		c.Reset()
		if got := c.Value(); got != 0 {
			t.Errorf("Value() after Reset() = %d, want 0", got)
		}
	})
}

func TestPoint(t *testing.T) {
	t.Run("distance", func(t *testing.T) {
		p1 := Point{X: 0, Y: 0}
		p2 := Point{X: 3, Y: 4}
		got := p1.Distance(p2)
		want := 5.0
		if math.Abs(got-want) > 0.0001 {
			t.Errorf("Distance() = %f, want %f", got, want)
		}
	})

	t.Run("translate", func(t *testing.T) {
		p := Point{X: 1, Y: 2}
		p.Translate(2, 3)
		if p.X != 3 || p.Y != 5 {
			t.Errorf("After Translate(2,3): Point{%d, %d}, want Point{3, 5}", p.X, p.Y)
		}
	})

	t.Run("string", func(t *testing.T) {
		p := Point{X: 3, Y: 4}
		got := p.String()
		want := "Point(3, 4)"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestConfiguration(t *testing.T) {
	t.Run("validate valid config", func(t *testing.T) {
		cfg := Configuration{Host: "localhost", Port: 8080}
		if !cfg.Validate() {
			t.Error("Validate() = false, want true for valid config")
		}
	})

	t.Run("validate invalid config", func(t *testing.T) {
		cfg := Configuration{Host: "", Port: 0}
		if cfg.Validate() {
			t.Error("Validate() = true, want false for invalid config")
		}
	})

	t.Run("apply defaults", func(t *testing.T) {
		cfg := Configuration{}
		cfg.ApplyDefaults()
		if cfg.Host != "localhost" {
			t.Errorf("Host = %q, want \"localhost\\", cfg.Host)
		}
		if cfg.Port != 8080 {
			t.Errorf("Port = %d, want 8080", cfg.Port)
		}
		if cfg.Timeout != 30 {
			t.Errorf("Timeout = %d, want 30", cfg.Timeout)
		}
	})
}

func TestTemperature(t *testing.T) {
	t.Run("to fahrenheit", func(t *testing.T) {
		temp := Temperature(0)
		got := temp.ToFahrenheit()
		want := 32.0
		if got != want {
			t.Errorf("ToFahrenheit() = %f, want %f", got, want)
		}

		temp = Temperature(100)
		got = temp.ToFahrenheit()
		want = 212.0
		if got != want {
			t.Errorf("ToFahrenheit() = %f, want %f", got, want)
		}
	})

	t.Run("is freezing", func(t *testing.T) {
		tests := []struct {
			temp Temperature
			want bool
		}{
			{Temperature(-5), true},
			{Temperature(0), true},
			{Temperature(1), false},
			{Temperature(20), false},
		}

		for _, tt := range tests {
			if got := tt.temp.IsFreezing(); got != tt.want {
				t.Errorf("Temperature(%d).IsFreezing() = %v, want %v", tt.temp, got, tt.want)
			}
		}
	})

	t.Run("warm", func(t *testing.T) {
		temp := Temperature(10)
		temp.Warm(5)
		if temp != 15 {
			t.Errorf("After Warm(5): Temperature = %d, want 15", temp)
		}
	})
}
