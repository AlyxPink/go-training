package main

import (
	"fmt"
	"math"
)

// Counter tracks an integer count
type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Value() int {
	return c.count
}

func (c *Counter) Reset() {
	c.count = 0
}

// Point represents a 2D coordinate
type Point struct {
	X, Y int
}

func (p Point) Distance(other Point) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func (p *Point) Translate(dx, dy int) {
	p.X += dx
	p.Y += dy
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
}

// Configuration represents app configuration
type Configuration struct {
	Host    string
	Port    int
	Timeout int
	Debug   bool
}

func (c *Configuration) Validate() bool {
	return c.Host != "" && c.Port > 0
}

func (c *Configuration) ApplyDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.Timeout == 0 {
		c.Timeout = 30
	}
}

// Temperature represents temperature in Celsius
type Temperature int

func (t Temperature) ToFahrenheit() float64 {
	return float64(t)*9/5 + 32
}

func (t Temperature) IsFreezing() bool {
	return t <= 0
}

func (t *Temperature) Warm(degrees int) {
	*t += Temperature(degrees)
}

func main() {
	c := Counter{}
	c.Increment()
	c.Increment()
	fmt.Println("Counter:", c.Value())

	p := Point{X: 0, Y: 0}
	p.Translate(3, 4)
	fmt.Println(p)
	fmt.Println("Distance:", p.Distance(Point{X: 0, Y: 0}))

	cfg := Configuration{}
	cfg.ApplyDefaults()
	fmt.Println("Valid:", cfg.Validate())

	temp := Temperature(0)
	fmt.Println("Freezing:", temp.IsFreezing())
	temp.Warm(10)
	fmt.Println("Temperature:", temp)
}
