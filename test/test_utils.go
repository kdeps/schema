package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/apple/pkl-go/pkl"
)

// TestMetrics tracks test execution metrics
type TestMetrics struct {
	TotalTests   int
	PassedTests  int
	FailedTests  int
	SkippedTests int
	StartTime    time.Time
	EndTime      time.Time
	TestResults  map[string]TestResult
	mu           sync.Mutex // protects TestResults and counters
}

// TestResult represents the result of a single test
type TestResult struct {
	Name     string
	Status   string // "PASS", "FAIL", "SKIP"
	Duration time.Duration
	Error    error
	Message  string
}

// TestSuite manages test execution and reporting
type TestSuite struct {
	metrics *TestMetrics
	logger  *TestLogger
	config  *TestConfig
}

// TestConfig holds test configuration
type TestConfig struct {
	Verbose       bool
	RetryCount    int
	RetryDelay    time.Duration
	Timeout       time.Duration
	Parallel      bool
	FilterPattern string
}

// TestLogger provides structured logging for tests
type TestLogger struct {
	verbose bool
}

// NewTestSuite creates a new test suite with default configuration
func NewTestSuite() *TestSuite {
	return &TestSuite{
		metrics: &TestMetrics{
			TestResults: make(map[string]TestResult),
		},
		logger: &TestLogger{verbose: true},
		config: &TestConfig{
			Verbose:       true,
			RetryCount:    3,
			RetryDelay:    time.Second,
			Timeout:       30 * time.Second,
			Parallel:      false,
			FilterPattern: "",
		},
	}
}

// RunTest executes a test with retry logic and metrics collection
func (ts *TestSuite) RunTest(t *testing.T, testName string, testFunc func(*testing.T) error) {
	startTime := time.Now()
	result := TestResult{
		Name:   testName,
		Status: "PASS",
	}

	// Execute test with retry logic
	var err error
	for attempt := 1; attempt <= ts.config.RetryCount; attempt++ {
		if attempt > 1 {
			ts.logger.Logf("Retrying test %s (attempt %d/%d)", testName, attempt, ts.config.RetryCount)
			time.Sleep(ts.config.RetryDelay)
		}

		err = testFunc(t)
		if err == nil {
			break
		}

		if attempt == ts.config.RetryCount {
			result.Status = "FAIL"
			result.Error = err
			result.Message = fmt.Sprintf("Test failed after %d attempts: %v", attempt, err)
		}
	}

	result.Duration = time.Since(startTime)
	ts.metrics.mu.Lock()
	ts.metrics.TestResults[testName] = result
	ts.metrics.mu.Unlock()

	// Update metrics
	ts.metrics.mu.Lock()
	ts.metrics.TotalTests++
	switch result.Status {
	case "PASS":
		ts.metrics.PassedTests++
	case "FAIL":
		ts.metrics.FailedTests++
	case "SKIP":
		ts.metrics.SkippedTests++
	}
	ts.metrics.mu.Unlock()

	// Log result
	ts.logger.Logf("Test %s: %s (%.2fs)", testName, result.Status, result.Duration.Seconds())
}

// GetMetrics returns the current test metrics
func (ts *TestSuite) GetMetrics() *TestMetrics {
	ts.metrics.EndTime = time.Now()
	ts.metrics.mu.Lock()
	defer ts.metrics.mu.Unlock()
	return ts.metrics
}

// PrintSummary prints a comprehensive test summary
func (ts *TestSuite) PrintSummary() {
	metrics := ts.GetMetrics()
	duration := metrics.EndTime.Sub(metrics.StartTime)

	const (
		headerFormat = "\n================================================================================\n"
		titleFormat  = "🧪 TEST SUITE SUMMARY\n"
		titleLine    = "================================================================================\n"
		footerFormat = "\n================================================================================\n"
	)

	fmt.Printf(headerFormat)
	fmt.Printf(titleFormat)
	fmt.Printf(titleLine)
	fmt.Printf("📊 EXECUTION METRICS:\n")
	fmt.Printf("   Total Tests: %d\n", metrics.TotalTests)
	fmt.Printf("   Passed: %d\n", metrics.PassedTests)
	fmt.Printf("   Failed: %d\n", metrics.FailedTests)
	fmt.Printf("   Skipped: %d\n", metrics.SkippedTests)
	fmt.Printf("   Success Rate: %.1f%%\n", float64(metrics.PassedTests)/float64(metrics.TotalTests)*100)
	fmt.Printf("   Duration: %.2fs\n", duration.Seconds())

	if metrics.FailedTests > 0 {
		fmt.Printf("\n❌ FAILED TESTS:\n")
		metrics.mu.Lock()
		for name, result := range metrics.TestResults {
			if result.Status == "FAIL" {
				fmt.Printf("   - %s: %s\n", name, result.Message)
			}
		}
		metrics.mu.Unlock()
	}

	fmt.Printf(footerFormat)
}

// Logf logs a formatted message if verbose mode is enabled
func (tl *TestLogger) Logf(format string, args ...interface{}) {
	if tl.verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// CreateTempPKLWorkspace creates a temporary workspace with all necessary PKL files
func CreateTempPKLWorkspace(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "pkl-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	// Copy deps/pkl directory
	depsDir := filepath.Join(tempDir, "deps", "pkl")
	if err := os.MkdirAll(depsDir, 0755); err != nil {
		cleanup()
		t.Fatalf("Failed to create deps/pkl dir: %v", err)
	}

	srcDeps := filepath.Clean(filepath.Join("..", "deps", "pkl"))
	entries, err := os.ReadDir(srcDeps)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to read deps/pkl dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDeps, entry.Name())
		dst := filepath.Join(depsDir, entry.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			cleanup()
			t.Fatalf("Failed to read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			cleanup()
			t.Fatalf("Failed to write %s: %v", dst, err)
		}
	}

	return tempDir, cleanup
}

// CopyPKLFile copies a PKL file to the temp directory with import path updates
func CopyPKLFile(t *testing.T, tempDir, fileName string) {
	src := filepath.Join(".", fileName)
	dst := filepath.Join(tempDir, fileName)

	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", src, err)
	}

	// Update import paths for tempDir
	updated := strings.ReplaceAll(string(data), "../deps/pkl/", "deps/pkl/")

	// Also copy any referenced test files
	if strings.Contains(string(data), "import \"./") {
		// Find all relative imports
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "import \"./") {
				// Extract the imported file name
				importLine := strings.TrimSpace(line)
				importFile := strings.TrimPrefix(importLine, "import \"./")
				importFile = strings.Split(importFile, "\"")[0] // Get everything before the first quote after the file name

				// Copy the imported file
				importSrc := filepath.Join(".", importFile)
				importDst := filepath.Join(tempDir, importFile)

				importData, err := os.ReadFile(importSrc)
				if err != nil {
					t.Fatalf("Failed to read imported file %s: %v", importSrc, err)
				}

				// Update import paths in the imported file
				importUpdated := strings.ReplaceAll(string(importData), "../deps/pkl/", "deps/pkl/")
				if err := os.WriteFile(importDst, []byte(importUpdated), 0644); err != nil {
					t.Fatalf("Failed to write imported file %s: %v", importDst, err)
				}
			}
		}
	}

	if err := os.WriteFile(dst, []byte(updated), 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", dst, err)
	}
}

// EvaluatePKLFile evaluates a PKL file and returns the result as a map
func EvaluatePKLFile(t *testing.T, evaluator pkl.Evaluator, fileName string) map[string]interface{} {
	source := pkl.FileSource(fileName)
	var module map[string]interface{}
	ctx := context.Background()
	if err := evaluator.EvaluateModule(ctx, source, &module); err != nil {
		t.Fatalf("Failed to evaluate %s: %v", fileName, err)
	}
	return module
}

// AssertTestResult checks if a test result is as expected
func AssertTestResult(t *testing.T, testName string, expected bool, actual bool, message string) {
	if expected != actual {
		t.Errorf("Test %s failed: %s (expected %v, got %v)", testName, message, expected, actual)
	}
}

// BenchmarkPKLEvaluation benchmarks PKL file evaluation
func BenchmarkPKLEvaluation(b *testing.B, evaluator pkl.Evaluator, fileName string) {
	source := pkl.FileSource(fileName)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var module map[string]interface{}
		if err := evaluator.EvaluateModule(ctx, source, &module); err != nil {
			b.Fatalf("Failed to evaluate %s: %v", fileName, err)
		}
	}
}
