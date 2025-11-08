package main

import (
	"flag"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	
	// Test basic parsing
	cli := ParseFlags()
	if cli == nil {
		t.Fatal("ParseFlags returned nil")
	}
}
