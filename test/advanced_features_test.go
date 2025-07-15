package test

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestAdvancedFeatures demonstrates all the advanced testing features
func TestAdvancedFeatures(t *testing.T) {
	t.Run("BenchmarkSystem", testBenchmarkSystem)
	t.Run("TestFiltering", testTestFiltering)
	t.Run("CoverageAnalysis", testCoverageAnalysis)
	t.Run("Integration", testAdvancedFeaturesIntegration)
}

// testBenchmarkSystem tests the benchmark system
func testBenchmarkSystem(t *testing.T) {
	benchmarkSystem := NewBenchmarkSystem(".")

	// Create a benchmark suite
	suite := benchmarkSystem.AddSuite("performance", "Performance benchmarks for core functionality")

	// Add benchmarks
	suite.AddBenchmark("fast_operation", "Test fast operations", func() error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}).SetSetup(func() error {
		// Setup code
		return nil
	}).SetTeardown(func() error {
		// Teardown code
		return nil
	}).AddMetadata("category", "unit")

	suite.AddBenchmark("medium_operation", "Test medium operations", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}).AddMetadata("category", "integration")

	suite.AddBenchmark("slow_operation", "Test slow operations", func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}).AddMetadata("category", "e2e")

	// Run benchmarks
	ctx := context.Background()
	run, err := benchmarkSystem.RunSuite(ctx, "performance")
	if err != nil {
		t.Fatalf("Failed to run benchmark suite: %v", err)
	}

	// Verify results
	if run.Summary.TotalBenchmarks != 3 {
		t.Errorf("Expected 3 benchmarks, got %d", run.Summary.TotalBenchmarks)
	}

	if run.Summary.Passed != 3 {
		t.Errorf("Expected 3 passed benchmarks, got %d", run.Summary.Passed)
	}

	// Export results
	err = benchmarkSystem.ExportResults(run, "json")
	if err != nil {
		t.Errorf("Failed to export benchmark results: %v", err)
	}

	// Analyze trends
	trends := benchmarkSystem.AnalyzeTrends()
	if len(trends) == 0 {
		t.Log("No trends available yet (first run)")
	}

	t.Logf("✅ Benchmark system test passed - ran %d benchmarks", run.Summary.TotalBenchmarks)
}

// testTestFiltering tests the test filtering system
func testTestFiltering(t *testing.T) {
	// Create test data
	tests := []*TestInfo{
		{
			Name:       "TestFastUnit",
			Tags:       []string{"unit", "fast"},
			Category:   "unit",
			Status:     "pass",
			Duration:   1 * time.Millisecond,
			Priority:   1,
			FlakyScore: 0.0,
		},
		{
			Name:       "TestSlowIntegration",
			Tags:       []string{"integration", "slow"},
			Category:   "integration",
			Status:     "pass",
			Duration:   5 * time.Second,
			Priority:   2,
			FlakyScore: 0.1,
		},
		{
			Name:       "TestFlakyE2E",
			Tags:       []string{"e2e", "flaky", "slow"},
			Category:   "e2e",
			Status:     "fail",
			Duration:   10 * time.Second,
			Priority:   3,
			FlakyScore: 0.8,
		},
		{
			Name:       "TestPerformance",
			Tags:       []string{"performance", "benchmark"},
			Category:   "performance",
			Status:     "pass",
			Duration:   2 * time.Second,
			Priority:   1,
			FlakyScore: 0.0,
			Performance: &PerformanceMetrics{
				OpsPerSecond: 1500,
				MemoryUsage:  50 * 1024 * 1024,
			},
		},
	}

	// Test quick tests filter
	quickFilter := NewQuickTestsFilter()
	quickTests := quickFilter.FilterTests(tests)
	if len(quickTests) != 1 {
		t.Errorf("Expected 1 quick test, got %d", len(quickTests))
	}

	// Test slow tests filter
	slowFilter := NewSlowTestsFilter()
	slowTests := slowFilter.FilterTests(tests)
	if len(slowTests) != 2 {
		t.Errorf("Expected 2 slow tests, got %d", len(slowTests))
	}

	// Test flaky tests filter
	flakyFilter := NewFlakyTestsFilter()
	flakyTests := flakyFilter.FilterTests(tests)
	if len(flakyTests) != 1 {
		t.Errorf("Expected 1 flaky test, got %d", len(flakyTests))
	}

	// Test unit tests filter
	unitFilter := NewUnitTestsFilter()
	unitTests := unitFilter.FilterTests(tests)
	if len(unitTests) != 1 {
		t.Errorf("Expected 1 unit test, got %d", len(unitTests))
	}

	// Test custom filter
	customFilter := NewTestFilter().
		AddCustomFilter("high_performance", "Tests with high performance", func(test *TestInfo) bool {
			return test.Performance != nil && test.Performance.OpsPerSecond > 1000
		}, nil)

	highPerfTests := customFilter.FilterTests(tests)
	if len(highPerfTests) != 1 {
		t.Errorf("Expected 1 high performance test, got %d", len(highPerfTests))
	}

	// Test complex filter
	complexFilter := NewTestFilter().
		AddTags("performance").
		SetDurationRange(nil, &[]time.Duration{3 * time.Second}[0]).
		SetPriorityRange(1, 2)

	complexTests := complexFilter.FilterTests(tests)
	if len(complexTests) != 1 {
		t.Errorf("Expected 1 test matching complex criteria, got %d", len(complexTests))
	}

	t.Logf("✅ Test filtering test passed - tested %d filters", 6)
}

// testCoverageAnalysis tests the coverage analysis system
func testCoverageAnalysis(t *testing.T) {
	coverageAnalyzer := NewCoverageAnalyzer(".")

	// Add sample coverage data
	coverageAnalyzer.AddFileCoverage("pkg/core/engine.go", &FileCoverage{
		Path:           "pkg/core/engine.go",
		TotalLines:     100,
		CoveredLines:   85,
		UncoveredLines: 15,
		Coverage:       85.0,
		Complexity:     5,
		Risk:           "low",
		Functions: map[string]*FunctionCoverage{
			"NewEngine": {
				Name:         "NewEngine",
				StartLine:    10,
				EndLine:      30,
				TotalLines:   20,
				CoveredLines: 20,
				Coverage:     100.0,
				Complexity:   2,
				Risk:         "low",
			},
			"Process": {
				Name:         "Process",
				StartLine:    35,
				EndLine:      80,
				TotalLines:   45,
				CoveredLines: 40,
				Coverage:     88.9,
				Complexity:   8,
				Risk:         "medium",
			},
		},
		Branches: map[string]*BranchCoverage{
			"10": {
				Line:            10,
				TotalBranches:   2,
				CoveredBranches: 2,
				Coverage:        100.0,
				Risk:            "low",
			},
		},
	})

	coverageAnalyzer.AddFileCoverage("internal/utils/helper.go", &FileCoverage{
		Path:           "internal/utils/helper.go",
		TotalLines:     50,
		CoveredLines:   30,
		UncoveredLines: 20,
		Coverage:       60.0,
		Complexity:     12,
		Risk:           "high",
		Functions: map[string]*FunctionCoverage{
			"HelperFunction": {
				Name:         "HelperFunction",
				StartLine:    5,
				EndLine:      45,
				TotalLines:   40,
				CoveredLines: 25,
				Coverage:     62.5,
				Complexity:   15,
				Risk:         "high",
			},
		},
	})

	coverageAnalyzer.AddFileCoverage("docs/README.md", &FileCoverage{
		Path:           "docs/README.md",
		TotalLines:     20,
		CoveredLines:   0,
		UncoveredLines: 20,
		Coverage:       0.0,
		Complexity:     1,
		Risk:           "low",
	})

	// Analyze coverage
	report := coverageAnalyzer.AnalyzeCoverage()

	// Verify summary
	if report.Summary.TotalFiles != 3 {
		t.Errorf("Expected 3 files, got %d", report.Summary.TotalFiles)
	}

	if report.Summary.OverallCoverage < 60.0 || report.Summary.OverallCoverage > 70.0 {
		t.Errorf("Expected coverage around 65%%, got %.2f%%", report.Summary.OverallCoverage)
	}

	// Test file queries
	uncoveredFiles := coverageAnalyzer.GetUncoveredFiles()
	if len(uncoveredFiles) != 1 {
		t.Errorf("Expected 1 uncovered file, got %d", len(uncoveredFiles))
	}

	lowCoverageFiles := coverageAnalyzer.GetLowCoverageFiles(70.0)
	if len(lowCoverageFiles) != 2 {
		t.Errorf("Expected 2 low coverage files, got %d", len(lowCoverageFiles))
	}

	criticalFiles := coverageAnalyzer.GetCriticalFiles()
	if len(criticalFiles) != 2 {
		t.Errorf("Expected 2 critical files, got %d", len(criticalFiles))
	}

	highRiskFiles := coverageAnalyzer.GetHighRiskFiles()
	if len(highRiskFiles) != 1 {
		t.Errorf("Expected 1 high risk file, got %d", len(highRiskFiles))
	}

	// Export report
	err := coverageAnalyzer.ExportReport(report, "json")
	if err != nil {
		t.Errorf("Failed to export coverage report: %v", err)
	}

	// Test risk calculation
	risk := coverageAnalyzer.CalculateRisk(90.0, 3)
	if risk != "low" {
		t.Errorf("Expected low risk, got %s", risk)
	}

	risk = coverageAnalyzer.CalculateRisk(30.0, 15)
	if risk != "critical" {
		t.Errorf("Expected critical risk, got %s", risk)
	}

	t.Logf("✅ Coverage analysis test passed - analyzed %d files", report.Summary.TotalFiles)
}

// testAdvancedFeaturesIntegration tests the integration of all advanced features
func testAdvancedFeaturesIntegration(t *testing.T) {
	// Create all systems
	benchmarkSystem := NewBenchmarkSystem(".")
	coverageAnalyzer := NewCoverageAnalyzer(".")
	testFilter := NewTestFilter()

	// Create a comprehensive test suite
	suite := benchmarkSystem.AddSuite("integration", "Integration tests for advanced features")

	// Add benchmarks with different characteristics
	suite.AddBenchmark("filter_performance", "Test filtering performance", func() error {
		// Simulate filtering operation
		tests := make([]*TestInfo, 1000)
		for i := 0; i < 1000; i++ {
			tests[i] = &TestInfo{
				Name:     fmt.Sprintf("Test%d", i),
				Tags:     []string{"unit", "integration", "e2e"},
				Category: "unit",
				Status:   "pass",
				Duration: time.Duration(i%100) * time.Millisecond,
				Priority: i % 5,
			}
		}

		// Apply multiple filters using the testFilter instance
		testFilter.AddTags("unit").
			SetDurationRange(nil, &[]time.Duration{50 * time.Millisecond}[0]).
			SetPriorityRange(1, 3)

		filtered := testFilter.FilterTests(tests)
		if len(filtered) == 0 {
			return fmt.Errorf("no tests matched filter criteria")
		}
		return nil
	}).AddMetadata("category", "performance")

	suite.AddBenchmark("coverage_analysis", "Test coverage analysis performance", func() error {
		// Simulate coverage analysis
		for i := 0; i < 100; i++ {
			coverageAnalyzer.AddFileCoverage(fmt.Sprintf("file%d.go", i), &FileCoverage{
				Path:         fmt.Sprintf("file%d.go", i),
				TotalLines:   100 + i,
				CoveredLines: 80 + i,
				Coverage:     float64(80+i) / float64(100+i) * 100,
				Complexity:   i % 10,
				Risk:         "low",
			})
		}

		report := coverageAnalyzer.AnalyzeCoverage()
		if report.Summary.TotalFiles != 100 {
			return fmt.Errorf("expected 100 files, got %d", report.Summary.TotalFiles)
		}
		return nil
	}).AddMetadata("category", "analysis")

	// Run benchmarks
	ctx := context.Background()
	run, err := benchmarkSystem.RunSuite(ctx, "integration")
	if err != nil {
		t.Fatalf("Failed to run integration benchmark suite: %v", err)
	}

	// Verify results
	if run.Summary.TotalBenchmarks != 2 {
		t.Errorf("Expected 2 benchmarks, got %d", run.Summary.TotalBenchmarks)
	}

	if run.Summary.Passed != 2 {
		t.Errorf("Expected 2 passed benchmarks, got %d", run.Summary.Passed)
	}

	// Test combined filtering and coverage analysis
	coverageReport := coverageAnalyzer.AnalyzeCoverage()

	// Export all results
	err = benchmarkSystem.ExportResults(run, "json")
	if err != nil {
		t.Errorf("Failed to export benchmark results: %v", err)
	}

	err = coverageAnalyzer.ExportReport(coverageReport, "html")
	if err != nil {
		t.Errorf("Failed to export coverage report: %v", err)
	}

	// Test trend analysis
	trends := benchmarkSystem.AnalyzeTrends()
	history := coverageAnalyzer.GetHistory()

	t.Logf("✅ Integration test passed:")
	t.Logf("   - Ran %d benchmarks", run.Summary.TotalBenchmarks)
	t.Logf("   - Analyzed %d files for coverage", coverageReport.Summary.TotalFiles)
	t.Logf("   - Generated %d trend data points", len(trends))
	t.Logf("   - Stored %d coverage reports", len(history))
}

// TestAdvancedFeaturesCLI tests CLI-like functionality
func TestAdvancedFeaturesCLI(t *testing.T) {
	t.Run("BenchmarkCLI", testBenchmarkCLI)
	t.Run("FilterCLI", testFilterCLI)
	t.Run("CoverageCLI", testCoverageCLI)
}

// testBenchmarkCLI simulates CLI benchmark commands
func testBenchmarkCLI(t *testing.T) {
	benchmarkSystem := NewBenchmarkSystem(".")

	// Simulate "benchmark run performance" command
	suite := benchmarkSystem.AddSuite("performance", "Performance benchmarks")
	suite.AddBenchmark("test1", "Test 1", func() error {
		time.Sleep(1 * time.Millisecond)
		return nil
	})
	suite.AddBenchmark("test2", "Test 2", func() error {
		time.Sleep(2 * time.Millisecond)
		return nil
	})

	run, err := benchmarkSystem.RunSuite(context.Background(), "performance")
	if err != nil {
		t.Fatalf("Benchmark run failed: %v", err)
	}

	// Simulate "benchmark export json" command
	err = benchmarkSystem.ExportResults(run, "json")
	if err != nil {
		t.Errorf("Benchmark export failed: %v", err)
	}

	// Simulate "benchmark trends" command
	trends := benchmarkSystem.AnalyzeTrends()
	t.Logf("Benchmark trends: %d data points", len(trends))

	t.Log("✅ Benchmark CLI simulation passed")
}

// testFilterCLI simulates CLI filter commands
func testFilterCLI(t *testing.T) {
	// Simulate "test filter --tags unit,fast" command
	filter := NewTestFilter().AddTags("unit", "fast")

	// Simulate "test filter --exclude slow" command
	filter.ExcludeTags("slow")

	// Simulate "test filter --category unit" command
	filter.AddCategories("unit")

	// Simulate "test filter --duration-max 5s" command
	maxDuration := 5 * time.Second
	filter.SetDurationRange(nil, &maxDuration)

	tests := []*TestInfo{
		{Name: "Test1", Tags: []string{"unit", "fast"}, Category: "unit", Duration: 1 * time.Millisecond},
		{Name: "Test2", Tags: []string{"unit", "slow"}, Category: "unit", Duration: 10 * time.Second},
		{Name: "Test3", Tags: []string{"integration"}, Category: "integration", Duration: 2 * time.Second},
	}

	filtered := filter.FilterTests(tests)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered test, got %d", len(filtered))
	}

	t.Log("✅ Filter CLI simulation passed")
}

// testCoverageCLI simulates CLI coverage commands
func testCoverageCLI(t *testing.T) {
	coverageAnalyzer := NewCoverageAnalyzer(".")

	// Simulate "coverage analyze" command
	coverageAnalyzer.AddFileCoverage("test.go", &FileCoverage{
		Path:         "test.go",
		TotalLines:   100,
		CoveredLines: 80,
		Coverage:     80.0,
		Risk:         "low",
	})

	report := coverageAnalyzer.AnalyzeCoverage()

	// Simulate "coverage export html" command
	err := coverageAnalyzer.ExportReport(report, "html")
	if err != nil {
		t.Errorf("Coverage export failed: %v", err)
	}

	// Simulate "coverage show uncovered" command
	uncovered := coverageAnalyzer.GetUncoveredFiles()
	t.Logf("Uncovered files: %d", len(uncovered))

	// Simulate "coverage show low" command
	lowCoverage := coverageAnalyzer.GetLowCoverageFiles(90.0)
	t.Logf("Low coverage files: %d", len(lowCoverage))

	t.Log("✅ Coverage CLI simulation passed")
}
