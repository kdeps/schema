package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestComprehensiveSuite demonstrates all test capabilities
func TestComprehensiveSuite(t *testing.T) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "comprehensive-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize test suite
	suite := NewTestSuite()

	// Test 1: Fixture Management
	t.Run("FixtureManagement", func(t *testing.T) {
		testFixtureManagement(t, tempDir)
	})

	// Test 2: Diagnostic System
	t.Run("DiagnosticSystem", func(t *testing.T) {
		testDiagnosticSystem(t, tempDir)
	})

	// Test 3: Parallel Execution
	t.Run("ParallelExecution", func(t *testing.T) {
		testParallelExecution(t, suite)
	})

	// Test 4: Analytics and Reporting
	t.Run("AnalyticsAndReporting", func(t *testing.T) {
		testAnalyticsAndReporting(t, tempDir)
	})

	// Test 5: Performance Analysis
	t.Run("PerformanceAnalysis", func(t *testing.T) {
		testPerformanceAnalysis(t, tempDir)
	})

	// Test 6: Error Handling and Recovery
	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t, suite)
	})
}

// testFixtureManagement tests the fixture management system
func testFixtureManagement(t *testing.T, tempDir string) {
	// Create test data generator
	generator := NewTestDataGenerator(tempDir)

	// Create default fixtures
	err := generator.CreateDefaultFixtures()
	if err != nil {
		t.Fatalf("Failed to create default fixtures: %v", err)
	}

	// Load fixtures
	fixtureManager := generator.fixtureManager
	err = fixtureManager.LoadFixtures()
	if err != nil {
		t.Fatalf("Failed to load fixtures: %v", err)
	}

	// Verify fixtures were created
	fixtures := fixtureManager.ListFixtures()
	if len(fixtures) == 0 {
		t.Error("No fixtures were created")
	}

	// Test specific fixture retrieval
	execFixture, err := fixtureManager.GetFixture("basic_exec")
	if err != nil {
		t.Errorf("Failed to get exec fixture: %v", err)
	}
	if execFixture == nil {
		t.Error("Exec fixture is nil")
	}

	// Test scenario builder
	builder := NewTestScenarioBuilder(generator)
	builder.AddResource("exec").
		AddResource("python").
		AddCustomResource("custom", map[string]interface{}{
			"name":  "test-custom",
			"value": "test-value",
		}).
		SetMetadata(map[string]interface{}{
			"created": time.Now().Format(time.RFC3339),
			"version": "1.0.0",
		})

	scenario := builder.Build()
	if len(scenario) < 3 {
		t.Error("Scenario should have at least 3 resources")
	}

	// Save scenario as fixture
	err = builder.SaveScenario("test_scenario", "Test scenario for comprehensive testing")
	if err != nil {
		t.Errorf("Failed to save scenario: %v", err)
	}

	t.Logf("✅ Fixture management test passed - created %d fixtures", len(fixtures))
}

// testDiagnosticSystem tests the diagnostic system
func testDiagnosticSystem(t *testing.T, tempDir string) {
	// Create diagnostic manager
	diagnosticManager := NewDiagnosticManager(tempDir)
	diagnosticManager.EnableDebugMode()

	// Create error analyzer
	errorAnalyzer := NewErrorAnalyzer(diagnosticManager)

	// Create debug helper
	debugHelper := NewDebugHelper(diagnosticManager)

	// Test error analysis
	testError := fmt.Errorf("test timeout error")
	analysis := errorAnalyzer.AnalyzeError("TestDiagnostic", testError, map[string]interface{}{
		"resource_type": "exec",
		"timeout":       30,
	})

	if analysis.ErrorType != "timeout" {
		t.Errorf("Expected error type 'timeout', got '%s'", analysis.ErrorType)
	}

	if analysis.Impact != "High - Test functionality impaired" {
		t.Errorf("Expected high impact, got '%s'", analysis.Impact)
	}

	// Test debug helper
	debugHelper.DebugTest("TestDebug", func() error {
		// Simulate a test that passes
		return nil
	})

	// Test debug helper with error
	debugHelper.DebugTest("TestDebugError", func() error {
		return fmt.Errorf("simulated test failure")
	})

	// Save diagnostics
	err := diagnosticManager.SaveDiagnostics("test_diagnostics.json")
	if err != nil {
		t.Errorf("Failed to save diagnostics: %v", err)
	}

	// Print diagnostics
	diagnosticManager.PrintDiagnostics()

	// Verify diagnostics were recorded
	diagnostics := diagnosticManager.GetDiagnosticsBySeverity("high")
	if len(diagnostics) == 0 {
		t.Error("Expected high severity diagnostics")
	}

	t.Log("✅ Diagnostic system test passed")
}

// testParallelExecution tests the parallel execution system
func testParallelExecution(t *testing.T, suite *TestSuite) {
	// Create a very simple parallel execution test
	scheduler := NewTestScheduler(suite)

	// Add just 3 simple tests with minimal dependencies
	scheduler.AddTest("setup", []string{}, []string{}, 1*time.Second, 1)
	scheduler.AddTest("test1", []string{"setup"}, []string{}, 500*time.Millisecond, 2)
	scheduler.AddTest("test2", []string{"setup"}, []string{}, 500*time.Millisecond, 2)

	// Set very low parallel limit
	scheduler.SetMaxParallel(1)

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute tests
	err := scheduler.Schedule(ctx, suite)
	if err != nil {
		t.Errorf("Simple parallel execution failed: %v", err)
	}

	// Get results
	results := scheduler.GetResults()
	if len(results) != 3 {
		t.Errorf("Expected 3 test results, got %d", len(results))
	}

	// Print results
	for testName, result := range results {
		t.Logf("Test %s: %s (%.2fs)", testName, result.Status, result.Duration.Seconds())
	}

	t.Log("✅ Parallel execution test passed")
}

// testAnalyticsAndReporting tests the analytics and reporting system
func testAnalyticsAndReporting(t *testing.T, tempDir string) {
	// Create analytics
	analytics := NewTestAnalytics(tempDir)

	// Create mock test metrics
	metrics := &TestMetrics{
		TestResults: make(map[string]TestResult),
	}

	// Add some test results
	metrics.mu.Lock()
	metrics.TestResults["test1"] = TestResult{
		Name:     "test1",
		Status:   "PASS",
		Duration: 1 * time.Second,
	}
	metrics.TestResults["test2"] = TestResult{
		Name:     "test2",
		Status:   "FAIL",
		Duration: 2 * time.Second,
		Error:    fmt.Errorf("test failure"),
	}
	metrics.mu.Unlock()

	// Record test runs for historical analysis
	analytics.RecordRun("run1", metrics, map[string]interface{}{
		"go_version": "1.21",
		"platform":   "darwin",
	}, []string{"unit", "integration"})

	// Create another run with different results
	metrics2 := &TestMetrics{
		TestResults: make(map[string]TestResult),
	}
	metrics2.mu.Lock()
	metrics2.TestResults["test1"] = TestResult{
		Name:     "test1",
		Status:   "PASS",
		Duration: 1500 * time.Millisecond,
	}
	metrics2.TestResults["test2"] = TestResult{
		Name:     "test2",
		Status:   "PASS",
		Duration: 1800 * time.Millisecond,
	}
	metrics2.mu.Unlock()

	analytics.RecordRun("run2", metrics2, map[string]interface{}{
		"go_version": "1.21",
		"platform":   "darwin",
	}, []string{"unit", "integration"})

	// Analyze trends
	trends := analytics.AnalyzeTrends()
	if trends == nil {
		t.Error("Trend analysis returned nil")
	}

	// Detect regressions
	regressions := analytics.DetectRegressions()
	if regressions == nil {
		t.Error("Regression detection returned nil")
	}

	// Generate reports in different formats
	formats := []string{"json", "csv", "html"}
	for _, format := range formats {
		err := analytics.GenerateReport(format)
		if err != nil {
			t.Errorf("Failed to generate %s report: %v", format, err)
		}
	}

	// Verify reports were created
	reportsDir := filepath.Join(tempDir, "reports")
	entries, err := os.ReadDir(reportsDir)
	if err != nil {
		t.Errorf("Failed to read reports directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No reports were generated")
	}

	t.Logf("✅ Analytics and reporting test passed - generated %d reports", len(entries))
}

// testPerformanceAnalysis tests the performance analysis system
func testPerformanceAnalysis(t *testing.T, tempDir string) {
	// Create analytics
	analytics := NewTestAnalytics(tempDir)

	// Create performance analyzer
	performanceAnalyzer := NewPerformanceAnalyzer(analytics)

	// Add some historical data for analysis
	for i := 0; i < 5; i++ {
		metrics := &TestMetrics{
			TestResults: make(map[string]TestResult),
		}
		metrics.mu.Lock()
		metrics.TestResults["perf_test"] = TestResult{
			Name:     "perf_test",
			Status:   "PASS",
			Duration: time.Duration(i+1) * time.Second,
		}
		metrics.mu.Unlock()

		analytics.RecordRun(fmt.Sprintf("perf_run_%d", i), metrics, map[string]interface{}{
			"iteration": i,
		}, []string{"performance"})
	}

	// Analyze performance
	analysis := performanceAnalyzer.AnalyzePerformance()
	if analysis == nil {
		t.Error("Performance analysis returned nil")
	}

	if len(analysis.Patterns) == 0 {
		t.Error("No performance patterns detected")
	}

	// Verify pattern types
	patternTypes := make(map[string]bool)
	for _, pattern := range analysis.Patterns {
		patternTypes[pattern.Type] = true
	}

	expectedTypes := []string{"slow_tests", "flaky_tests", "resource_usage"}
	for _, expectedType := range expectedTypes {
		if !patternTypes[expectedType] {
			t.Errorf("Expected pattern type '%s' not found", expectedType)
		}
	}

	t.Log("✅ Performance analysis test passed")
}

// testErrorHandling tests error handling and recovery
func testErrorHandling(t *testing.T, suite *TestSuite) {
	// Test retry logic
	attempts := 0
	err := suite.RunTestWithContext(context.Background(), "TestRetry", func(t *testing.T) error {
		attempts++
		if attempts < 3 {
			return fmt.Errorf("simulated failure on attempt %d", attempts)
		}
		return nil
	})

	if err != nil {
		t.Errorf("Test should have succeeded after retries: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = suite.RunTestWithContext(ctx, "TestCancellation", func(t *testing.T) error {
		time.Sleep(1 * time.Second) // This should not execute
		return nil
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	// Test timeout
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer timeoutCancel()

	err = suite.RunTestWithContext(timeoutCtx, "TestTimeout", func(t *testing.T) error {
		// Check for context cancellation during execution
		for i := 0; i < 100; i++ {
			select {
			case <-timeoutCtx.Done():
				return timeoutCtx.Err()
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
		return nil
	})

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
	}

	t.Log("✅ Error handling test passed")
}

// TestComprehensiveAnalytics demonstrates analytics functionality
func TestComprehensiveAnalytics(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "analytics-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	analytics := NewTestAnalytics(tempDir)

	// Test trend analysis with minimal data
	trends := analytics.AnalyzeTrends()
	if trends == nil {
		t.Error("Trend analysis should not return nil even with no data")
	}

	// Test regression detection with minimal data
	regressions := analytics.DetectRegressions()
	if regressions == nil {
		t.Error("Regression detection should not return nil even with no data")
	}

	t.Log("✅ Comprehensive analytics test passed")
}

// TestComprehensiveFixtures demonstrates fixture functionality
func TestComprehensiveFixtures(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "fixtures-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	generator := NewTestDataGenerator(tempDir)

	// Test data generation
	execData := generator.GenerateExecResource()
	if execData["command"] == "" {
		t.Error("Exec resource should have a command")
	}

	pythonData := generator.GeneratePythonResource()
	if pythonData["code"] == "" {
		t.Error("Python resource should have code")
	}

	// Test scenario generation
	basicScenario := generator.GenerateTestScenario("basic")
	if len(basicScenario) < 3 {
		t.Error("Basic scenario should have at least 3 resources")
	}

	t.Log("✅ Comprehensive fixtures test passed")
}

// TestComprehensiveDiagnostics demonstrates diagnostic functionality
func TestComprehensiveDiagnostics(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "diagnostics-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	diagnosticManager := NewDiagnosticManager(tempDir)

	// Test diagnostic recording
	diagnostic := diagnosticManager.AddDiagnostic("TestDiagnostic", "test error", map[string]interface{}{
		"test": "value",
	})

	if diagnostic.TestName != "TestDiagnostic" {
		t.Error("Diagnostic should have correct test name")
	}

	if diagnostic.Severity == "" {
		t.Error("Diagnostic should have severity")
	}

	t.Log("✅ Comprehensive diagnostics test passed")
}

// TestComprehensiveParallel demonstrates parallel execution
func TestComprehensiveParallel(t *testing.T) {
	suite := NewTestSuite()

	// Create a simple parallel execution test
	scheduler := NewTestScheduler(suite)
	scheduler.AddTest("simple_test", []string{}, []string{}, 5*time.Second, 1)
	scheduler.SetMaxParallel(1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := scheduler.Schedule(ctx, suite)
	if err != nil {
		t.Errorf("Simple parallel execution failed: %v", err)
	}

	results := scheduler.GetResults()
	if len(results) != 1 {
		t.Errorf("Expected 1 test result, got %d", len(results))
	}

	t.Log("✅ Comprehensive parallel execution test passed")
}

// TestAnalyticsExport demonstrates the analytics export functionality
func TestAnalyticsExport(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "analytics-export-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory for test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create analytics instance
	analytics := NewTestAnalytics(tempDir)

	// Record some test runs for analysis
	metrics1 := &TestMetrics{
		TotalTests:  10,
		PassedTests: 8,
		FailedTests: 2,
		StartTime:   time.Now().Add(-time.Hour),
		EndTime:     time.Now(),
		TestResults: map[string]TestResult{
			"test1": {Name: "test1", Status: "PASS", Duration: time.Second},
			"test2": {Name: "test2", Status: "FAIL", Duration: 2 * time.Second},
		},
	}

	metrics2 := &TestMetrics{
		TotalTests:  12,
		PassedTests: 11,
		FailedTests: 1,
		StartTime:   time.Now().Add(-30 * time.Minute),
		EndTime:     time.Now(),
		TestResults: map[string]TestResult{
			"test1": {Name: "test1", Status: "PASS", Duration: 500 * time.Millisecond},
			"test2": {Name: "test2", Status: "PASS", Duration: 1 * time.Second},
		},
	}

	analytics.RecordRun("run1", metrics1, map[string]interface{}{"env": "test"}, []string{"unit"})
	analytics.RecordRun("run2", metrics2, map[string]interface{}{"env": "test"}, []string{"unit"})

	// Build report
	report := analytics.buildReport()

	// Test export in different formats
	formats := ExportFormats{
		JSON:     true,
		Markdown: true,
		HTML:     true,
	}

	baseFilename := "test_analytics_export"
	if err := ExportAnalyticsReport(report, formats, baseFilename); err != nil {
		t.Fatalf("Failed to export analytics report: %v", err)
	}

	// Verify files were created in reports directory
	expectedFiles := []string{
		filepath.Join("reports", baseFilename+".json"),
		filepath.Join("reports", baseFilename+".md"),
		filepath.Join("reports", baseFilename+".html"),
	}

	for _, filename := range expectedFiles {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", filename)
		} else {
			t.Logf("Successfully created %s", filename)
		}
	}

	// Test individual format exports
	t.Run("JSONOnly", func(t *testing.T) {
		jsonFormats := ExportFormats{JSON: true}
		if err := ExportAnalyticsReport(report, jsonFormats, "json_only"); err != nil {
			t.Errorf("Failed to export JSON only: %v", err)
		}
	})

	t.Run("MarkdownOnly", func(t *testing.T) {
		mdFormats := ExportFormats{Markdown: true}
		if err := ExportAnalyticsReport(report, mdFormats, "md_only"); err != nil {
			t.Errorf("Failed to export Markdown only: %v", err)
		}
	})

	t.Run("HTMLOnly", func(t *testing.T) {
		htmlFormats := ExportFormats{HTML: true}
		if err := ExportAnalyticsReport(report, htmlFormats, "html_only"); err != nil {
			t.Errorf("Failed to export HTML only: %v", err)
		}
	})

	// Cleanup test files
	t.Cleanup(func() {
		os.RemoveAll("reports")
	})

	t.Log("Analytics export functionality working correctly")
}
