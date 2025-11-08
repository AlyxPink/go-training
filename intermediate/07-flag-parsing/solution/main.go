package main

import (
	"flag"
	"fmt"
)

type CLI struct {
	Verbose bool
	Output  string
	Count   int
}

func ParseFlags() *CLI {
	cli := &CLI{}
	
	flag.BoolVar(&cli.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&cli.Verbose, "v", false, "Enable verbose output (short)")
	flag.StringVar(&cli.Output, "output", "stdout", "Output destination")
	flag.StringVar(&cli.Output, "o", "stdout", "Output destination (short)")
	flag.IntVar(&cli.Count, "count", 1, "Number of iterations")
	
	flag.Parse()
	
	return cli
}

func main() {
	cli := ParseFlags()
	fmt.Printf("Verbose: %v, Output: %s, Count: %d\n", 
		cli.Verbose, cli.Output, cli.Count)
}
