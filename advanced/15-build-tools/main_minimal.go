// +build minimal

package main

import "fmt"

// Minimal build with features disabled
var (
	Version   = "minimal"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var (
	EnableLogging   = false
	EnableProfiling = false
	EnableDebug     = false
)

func buildInfo() string {
	return fmt.Sprintf("Version: %s (minimal build)", Version)
}

func platformInfo() string {
	return "Minimal build - platform info disabled"
}

func main() {
	fmt.Println("=== Minimal Build ===")
	fmt.Println(buildInfo())
}
