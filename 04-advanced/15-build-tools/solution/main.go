// +build !minimal

package main

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var (
	EnableLogging   = true
	EnableProfiling = true
	EnableDebug     = true
)

func buildInfo() string {
	return fmt.Sprintf(`Version: %s
Build Time: %s
Git Commit: %s
Go Version: %s`, Version, BuildTime, GitCommit, runtime.Version())
}

func platformInfo() string {
	return fmt.Sprintf(`OS: %s
Architecture: %s
CPUs: %d`, runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
}

func main() {
	fmt.Println("=== Build Information ===")
	fmt.Println(buildInfo())

	fmt.Println("\n=== Platform Information ===")
	fmt.Println(platformInfo())

	fmt.Println("\n=== Features ===")
	fmt.Printf("Logging: %v\n", EnableLogging)
	fmt.Printf("Profiling: %v\n", EnableProfiling)
	fmt.Printf("Debug: %v\n", EnableDebug)
}
