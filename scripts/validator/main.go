package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	failStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
)

// Command-line flags
var (
	runStarter   bool
	runSolutions bool
	verbose      bool
	outputFile   string
	noColor      bool
)

func main() {
	// Parse flags
	flag.BoolVar(&runStarter, "starter", false, "Validate starter code only")
	flag.BoolVar(&runSolutions, "solutions", false, "Validate solutions only")
	flag.BoolVar(&verbose, "v", false, "Show detailed test output (verbose mode)")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed test output (verbose mode)")
	flag.StringVar(&outputFile, "output", "", "Write failed exercises to file")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.Parse()

	// If neither specified, run both
	runAll := !runStarter && !runSolutions
	if runAll {
		runStarter = true
		runSolutions = true
	}

	// Disable colors if requested
	if noColor {
		successStyle = lipgloss.NewStyle()
		failStyle = lipgloss.NewStyle()
		warnStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
		titleStyle = lipgloss.NewStyle()
	}

	// Get repo root
	repoRoot := filepath.Join("..", "..")
	if flag.NArg() > 0 {
		repoRoot = flag.Arg(0)
	}

	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving repo root: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("═══════════════════════════════════════════════════════"))
	fmt.Println(titleStyle.Render("  Go Training Unified Validator"))
	fmt.Println(titleStyle.Render("═══════════════════════════════════════════════════════"))
	fmt.Println()

	startTime := time.Now()
	var starterExitCode, solutionExitCode int

	// Run starter validation
	if runStarter {
		fmt.Println(titleStyle.Render("🔄 Running Starter Code Validation..."))
		fmt.Println()
		starterExitCode = runValidator("starter-validator", absRoot)
		fmt.Println()
	}

	// Run solution validation
	if runSolutions {
		fmt.Println(titleStyle.Render("🔄 Running Solution Code Validation..."))
		fmt.Println()
		solutionExitCode = runValidator("solution-validator", absRoot)
		fmt.Println()
	}

	elapsed := time.Since(startTime)

	// Print unified summary
	fmt.Println(titleStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(titleStyle.Render("UNIFIED SUMMARY"))
	fmt.Println(titleStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	if runStarter {
		if starterExitCode == 0 {
			fmt.Println(successStyle.Render("✅") + " Starter Code: PASSED")
		} else {
			fmt.Println(failStyle.Render("❌") + " Starter Code: FAILED")
		}
	}

	if runSolutions {
		if solutionExitCode == 0 {
			fmt.Println(successStyle.Render("✅") + " Solutions: PASSED")
		} else {
			fmt.Println(failStyle.Render("❌") + " Solutions: FAILED")
		}
	}

	fmt.Println()
	fmt.Printf("Total Time: %s\n", dimStyle.Render(elapsed.Round(time.Second).String()))
	fmt.Println()

	// Exit with failure if any validation failed
	if starterExitCode != 0 || solutionExitCode != 0 {
		fmt.Println(failStyle.Render("❌ Some validations failed"))
		os.Exit(1)
	} else {
		fmt.Println(successStyle.Render("✅ All validations passed"))
		os.Exit(0)
	}
}

func runValidator(name string, repoRoot string) int {
	// Build path to validator binary
	validatorPath := filepath.Join("..", name, name)
	if _, err := os.Stat(validatorPath); os.IsNotExist(err) {
		// Try without directory prefix (for when running from repo root)
		validatorPath = filepath.Join("scripts", name, name)
		if _, err := os.Stat(validatorPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: %s binary not found\n", name)
			fmt.Fprintf(os.Stderr, "Please build it first: cd scripts/%s && go build\n", name)
			return 1
		}
	}

	// Build arguments
	args := []string{}
	if verbose {
		args = append(args, "-v")
	}
	if noColor {
		args = append(args, "-no-color")
	}
	if outputFile != "" {
		// Append validator name to output file
		outFile := outputFile
		if ext := filepath.Ext(outputFile); ext != "" {
			base := outputFile[:len(outputFile)-len(ext)]
			outFile = fmt.Sprintf("%s_%s%s", base, name, ext)
		} else {
			outFile = fmt.Sprintf("%s_%s", outputFile, name)
		}
		args = append(args, "-output", outFile)
	}
	args = append(args, repoRoot)

	// Run validator
	cmd := exec.Command(validatorPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
