package main

import (
	// TODO: Uncomment for flag parsing
	// "flag"
	"fmt"
)

// CLI represents command-line configuration
type CLI struct {
	Verbose bool
	Output  string
	Count   int
}

// ParseFlags parses command-line flags
func ParseFlags() *CLI {
	// TODO: Define flags and parse
	return nil
}

func main() {
	cli := ParseFlags()
	if cli != nil {
		fmt.Printf("Verbose: %v, Output: %s, Count: %d\n", 
			cli.Verbose, cli.Output, cli.Count)
	}
}
