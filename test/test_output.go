package test

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
)

// TestResultSummary holds summary info for a test run
type TestResultSummary struct {
	Total     int
	Passed    int
	Failed    int
	Skipped   int
	Flaky     int
	SlowTests []string
	Failures  []string
	StartTime time.Time
	EndTime   time.Time
}

// PrintTestHeader prints a formatted test header
func PrintTestHeader(name string) {
	fmt.Printf("%s=== RUN   %s%s\n", ColorBlue, name, ColorReset)
}

// PrintTestResult prints a formatted test result
func PrintTestResult(name string, status string, duration time.Duration, details string) {
	color := ColorGreen
	switch status {
	case "PASS":
		color = ColorGreen
	case "FAIL":
		color = ColorRed
	case "SKIP":
		color = ColorYellow
	case "FLAKY":
		color = ColorCyan
	}
	fmt.Printf("%s--- %s: %s (%.2fs)%s\n", color, status, name, duration.Seconds(), ColorReset)
	if details != "" {
		fmt.Printf("%s      %s%s\n", ColorGray, details, ColorReset)
	}
}

// PrintTestSummary prints a summary of the test run
func PrintTestSummary(summary *TestResultSummary) {
	total := summary.Total
	passed := summary.Passed
	failed := summary.Failed
	skipped := summary.Skipped
	flaky := summary.Flaky
	dur := summary.EndTime.Sub(summary.StartTime)

	fmt.Printf("\n%sTest Summary:%s\n", ColorBlue, ColorReset)
	fmt.Printf("  Total:   %d\n", total)
	fmt.Printf("  Passed:  %s%d%s\n", ColorGreen, passed, ColorReset)
	fmt.Printf("  Failed:  %s%d%s\n", ColorRed, failed, ColorReset)
	fmt.Printf("  Skipped: %s%d%s\n", ColorYellow, skipped, ColorReset)
	fmt.Printf("  Flaky:   %s%d%s\n", ColorCyan, flaky, ColorReset)
	fmt.Printf("  Duration: %.2fs\n", dur.Seconds())

	if len(summary.SlowTests) > 0 {
		fmt.Printf("%s  Slow Tests:%s\n", ColorYellow, ColorReset)
		for _, name := range summary.SlowTests {
			fmt.Printf("    %s\n", name)
		}
	}
	if len(summary.Failures) > 0 {
		fmt.Printf("%s  Failures:%s\n", ColorRed, ColorReset)
		for _, name := range summary.Failures {
			fmt.Printf("    %s\n", name)
		}
	}
}

// ShouldColor returns true if output should be colorized
func ShouldColor() bool {
	return isatty(os.Stdout.Fd())
}

// isatty checks if the file descriptor is a terminal
func isatty(fd uintptr) bool {
	// Simple check for *nix systems
	return strings.Contains(os.Getenv("TERM"), "xterm") || strings.Contains(os.Getenv("TERM"), "color")
}
