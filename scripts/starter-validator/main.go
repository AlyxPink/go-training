package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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
	verbose    bool
	outputFile string
	noColor    bool
	showProgress bool
)

type TestResult struct {
	Exercise string
	Status   string // "pass", "fail", "warn"
	Message  string
	Logs     []string
	Duration time.Duration
}

type ExerciseTester struct {
	repoRoot string
	printMu  sync.Mutex
	verbose  bool
}

func NewExerciseTester(repoRoot string, verbose bool) *ExerciseTester {
	return &ExerciseTester{
		repoRoot: repoRoot,
		verbose:  verbose,
	}
}

func (et *ExerciseTester) printResult(result TestResult) {
	et.printMu.Lock()
	defer et.printMu.Unlock()

	if result.Status == "pass" {
		// Success: single line (unless verbose)
		if et.verbose {
			for _, log := range result.Logs {
				fmt.Println(log)
			}
		}
		fmt.Println(successStyle.Render("✅") + " " + result.Exercise)
		if et.verbose {
			fmt.Println()
		}
	} else {
		// Failure/Warning: show all logs
		for _, log := range result.Logs {
			fmt.Println(log)
		}
		if result.Status == "fail" {
			fmt.Println(failStyle.Render(fmt.Sprintf("  ❌ %s: %s", result.Exercise, result.Message)))
		} else if result.Status == "warn" {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  ⚠️  %s: %s", result.Exercise, result.Message)))
		}
		fmt.Println()
	}
}

func (et *ExerciseTester) testExercise(exercise string) TestResult {
	start := time.Now()
	result := TestResult{
		Exercise: exercise,
		Logs:     make([]string, 0),
	}

	exercisePath := filepath.Join(et.repoRoot, exercise)

	// Log progress (verbose, will be cleared on success)
	result.Logs = append(result.Logs, fmt.Sprintf("🔄 Testing %s...", exercise))

	// Check go.mod exists
	if _, err := os.Stat(filepath.Join(exercisePath, "go.mod")); os.IsNotExist(err) {
		result.Status = "fail"
		result.Message = "Missing go.mod"
		result.Logs = append(result.Logs, "  ❌ Missing go.mod")
		result.Duration = time.Since(start)
		return result
	}

	// Download dependencies
	result.Logs = append(result.Logs, "  📥 Downloading dependencies...")
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = exercisePath
	if err := cmd.Run(); err != nil {
		result.Status = "fail"
		result.Message = "Dependency download failed"
		result.Logs = append(result.Logs, "  ❌ Dependency download failed")
		result.Duration = time.Since(start)
		return result
	}

	// Verify dependencies
	result.Logs = append(result.Logs, "  🔍 Verifying dependencies...")
	cmd = exec.Command("go", "mod", "verify")
	cmd.Dir = exercisePath
	if err := cmd.Run(); err != nil {
		result.Status = "fail"
		result.Message = "Dependency verification failed"
		result.Logs = append(result.Logs, "  ❌ Dependency verification failed")
		result.Duration = time.Since(start)
		return result
	}

	// Compile starter code
	result.Logs = append(result.Logs, "  🔨 Compiling starter code...")
	cmd = exec.Command("go", "build", "-v", "./...")
	cmd.Dir = exercisePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Status = "fail"
		result.Message = "Starter compilation failed"
		result.Logs = append(result.Logs, "  ❌ Compilation failed")
		if len(output) > 0 {
			result.Logs = append(result.Logs, "")
			result.Logs = append(result.Logs, "Compilation output:")
			result.Logs = append(result.Logs, string(output))
		}
		result.Duration = time.Since(start)
		return result
	}

	// Run tests with timeout
	result.Logs = append(result.Logs, "  🧪 Running tests (30s timeout)...")
	cmd = exec.Command("go", "test", "-v", "./...")
	cmd.Dir = exercisePath

	// Set timeout
	type testOutput struct {
		output []byte
		err    error
	}
	done := make(chan testOutput, 1)
	go func() {
		output, err := cmd.CombinedOutput()
		done <- testOutput{output, err}
	}()

	select {
	case testOut := <-done:
		if et.verbose && len(testOut.output) > 0 {
			result.Logs = append(result.Logs, "")
			result.Logs = append(result.Logs, "Test output:")
			result.Logs = append(result.Logs, string(testOut.output))
		}
		if testOut.err == nil {
			// Tests passed - this is good
			result.Status = "pass"
			result.Message = ""
		} else {
			// Tests failed - this is expected for starter code
			result.Status = "pass"
			result.Message = ""
		}
	case <-time.After(30 * time.Second):
		cmd.Process.Kill()
		result.Status = "warn"
		result.Message = "Tests timed out after 30s"
		result.Logs = append(result.Logs, "  ⏱️  Tests timed out (likely hanging/infinite loop)")
	}

	result.Duration = time.Since(start)
	return result
}

func findExercises(repoRoot string) ([]string, error) {
	exercises := make([]string, 0)

	// Expected top-level directories for exercises
	validPrefixes := []string{"basics", "intermediate", "advanced", "concurrency", "projects"}

	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip unwanted directories
		if info.IsDir() {
			name := info.Name()
			relPath, _ := filepath.Rel(repoRoot, path)

			// Skip backup, vendor, hidden dirs, scripts, bin, etc.
			if strings.Contains(name, "backup") || name == "vendor" ||
				name == ".git" || strings.HasPrefix(name, ".") ||
				name == "bin" || name == "scripts" || name == "claudedocs" {
				return filepath.SkipDir
			}

			// Skip if not under a valid prefix
			if relPath != "." {
				parts := strings.Split(relPath, string(filepath.Separator))
				if len(parts) > 0 {
					isValid := false
					for _, prefix := range validPrefixes {
						if parts[0] == prefix {
							isValid = true
							break
						}
					}
					if !isValid {
						return filepath.SkipDir
					}
				}
			}
		}

		if info.Name() == "go.mod" {
			relPath, err := filepath.Rel(repoRoot, filepath.Dir(path))
			if err != nil {
				return err
			}

			// Skip solution directories - they should pass all tests
			if strings.Contains(relPath, "solution") {
				return nil
			}

			// Only include paths with numbers (exercise directories)
			parts := strings.Split(relPath, string(filepath.Separator))
			for _, part := range parts {
				if len(part) > 0 && part[0] >= '0' && part[0] <= '9' {
					exercises = append(exercises, relPath)
					break
				}
			}
		}

		return nil
	})

	sort.Strings(exercises)
	return exercises, err
}

func main() {
	// Parse flags
	flag.BoolVar(&verbose, "v", false, "Show detailed test output (verbose mode)")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed test output (verbose mode)")
	flag.StringVar(&outputFile, "output", "", "Write failed exercises to file")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&showProgress, "progress", true, "Show live progress updates")
	flag.Parse()

	// Disable colors if requested
	if noColor {
		successStyle = lipgloss.NewStyle()
		failStyle = lipgloss.NewStyle()
		warnStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
		titleStyle = lipgloss.NewStyle()
	}

	// When running from scripts/starter-validator, go up two levels to project root
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
	fmt.Println(titleStyle.Render("  CI Starter Code Validation (Go Edition)"))
	fmt.Println(titleStyle.Render("═══════════════════════════════════════════════════════"))
	fmt.Println()

	// Find exercises
	fmt.Println(dimStyle.Render("🔍 Finding exercises..."))
	exercises, err := findExercises(absRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding exercises: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nTesting %s exercises in parallel (10 concurrent)...\n",
		titleStyle.Render(fmt.Sprintf("%d", len(exercises))))
	fmt.Println(dimStyle.Render("Each exercise: download deps → verify → compile → test"))
	fmt.Println()

	startTime := time.Now()

	// Run tests concurrently
	tester := NewExerciseTester(absRoot, verbose)
	results := make(chan TestResult, len(exercises))
	semaphore := make(chan struct{}, 10) // 10 concurrent workers

	// Progress tracking
	var completed, total int32 = 0, int32(len(exercises))
	var progressMu sync.Mutex
	stopProgress := make(chan bool)

	// Start progress monitor if enabled
	if showProgress && !verbose {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-stopProgress:
					return
				case <-ticker.C:
					progressMu.Lock()
					current := completed
					progressMu.Unlock()
					fmt.Printf("\r%s Testing... %d/%d exercises completed",
						dimStyle.Render("🔄"),
						current,
						total)
				}
			}
		}()
	}

	var wg sync.WaitGroup
	for _, exercise := range exercises {
		wg.Add(1)
		go func(ex string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			result := tester.testExercise(ex)
			results <- result

			// Print result immediately (synchronized)
			tester.printResult(result)

			// Update progress
			progressMu.Lock()
			completed++
			progressMu.Unlock()
		}(exercise)
	}

	// Wait for all tests to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	allResults := make([]TestResult, 0, len(exercises))
	for result := range results {
		allResults = append(allResults, result)
	}

	// Stop progress monitor
	if showProgress && !verbose {
		stopProgress <- true
		fmt.Print("\r" + strings.Repeat(" ", 80) + "\r") // Clear progress line
	}

	elapsed := time.Since(startTime)

	// Print summary
	fmt.Println()
	fmt.Println(titleStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(titleStyle.Render("SUMMARY"))
	fmt.Println(titleStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	passed := 0
	failed := 0
	warned := 0
	failedExercises := make([]TestResult, 0)
	warnedExercises := make([]TestResult, 0)

	for _, result := range allResults {
		switch result.Status {
		case "pass":
			passed++
		case "fail":
			failed++
			failedExercises = append(failedExercises, result)
		case "warn":
			warned++
			warnedExercises = append(warnedExercises, result)
		}
	}

	// Show failures
	if failed > 0 {
		fmt.Println(failStyle.Render(fmt.Sprintf("❌ FAILED (%d):", failed)))
		for _, result := range failedExercises {
			fmt.Printf("  - %s - %s\n", result.Exercise, result.Message)
		}
		fmt.Println()
	}

	// Show warnings
	if warned > 0 {
		fmt.Println(warnStyle.Render(fmt.Sprintf("⚠️  WARNINGS (%d):", warned)))
		for _, result := range warnedExercises {
			fmt.Printf("  - %s - %s\n", result.Exercise, result.Message)
		}
		fmt.Println()
	}

	fmt.Printf("Total:    %d\n", len(exercises))
	fmt.Printf("Passed:   %s\n", successStyle.Render(fmt.Sprintf("%d", passed)))
	fmt.Printf("Failed:   %s\n", failStyle.Render(fmt.Sprintf("%d", failed)))
	fmt.Printf("Warnings: %s\n", warnStyle.Render(fmt.Sprintf("%d", warned)))
	fmt.Printf("Time:     %s\n", dimStyle.Render(elapsed.Round(time.Second).String()))
	fmt.Println()

	// Write output file if requested
	if outputFile != "" && (failed > 0 || warned > 0) {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		} else {
			defer f.Close()
			if failed > 0 {
				fmt.Fprintln(f, "# Failed Exercises")
				for _, result := range failedExercises {
					fmt.Fprintf(f, "%s - %s\n", result.Exercise, result.Message)
				}
			}
			if warned > 0 {
				if failed > 0 {
					fmt.Fprintln(f, "")
				}
				fmt.Fprintln(f, "# Warned Exercises")
				for _, result := range warnedExercises {
					fmt.Fprintf(f, "%s - %s\n", result.Exercise, result.Message)
				}
			}
			fmt.Printf("Output written to: %s\n", outputFile)
			fmt.Println()
		}
	}

	if failed == 0 {
		fmt.Println(successStyle.Render("✅ All starter code validations passed"))
		os.Exit(0)
	} else {
		fmt.Println(failStyle.Render("❌ Some starter validations failed"))
		os.Exit(1)
	}
}
