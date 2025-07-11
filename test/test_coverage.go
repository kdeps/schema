package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

// TestCoverageReport represents a test coverage report
type TestCoverageReport struct {
	TotalTests      int                     `json:"total_tests"`
	CoveredTests    int                     `json:"covered_tests"`
	UncoveredTests  int                     `json:"uncovered_tests"`
	CoverageRate    float64                 `json:"coverage_rate"`
	TestCategories  map[string]TestCategory `json:"test_categories"`
	Recommendations []string                `json:"recommendations"`
}

// TestCategory represents a category of tests
type TestCategory struct {
	Name           string   `json:"name"`
	TotalTests     int      `json:"total_tests"`
	CoveredTests   int      `json:"covered_tests"`
	UncoveredTests int      `json:"uncovered_tests"`
	CoverageRate   float64  `json:"coverage_rate"`
	TestFiles      []string `json:"test_files"`
	MissingTests   []string `json:"missing_tests"`
}

// GenerateCoverageReport generates a comprehensive coverage report
func GenerateCoverageReport(t *testing.T) *TestCoverageReport {
	report := &TestCoverageReport{
		TestCategories: make(map[string]TestCategory),
	}

	// Analyze test categories
	categories := []struct {
		name        string
		testFiles   []string
		description string
	}{
		{
			name: "Resource_Readers",
			testFiles: []string{
				"agent_resource_reader_test.go",
				"pklres_resource_reader_test.go",
				"real_pklres_reader_test.go",
			},
			description: "Tests for resource reader implementations",
		},
		{
			name: "PKL_Integration",
			testFiles: []string{
				"pklres_integration_test.go",
				"integration_test.go",
			},
			description: "Tests for PKL integration scenarios",
		},
		{
			name: "Schema_Validation",
			testFiles: []string{
				"assets_test.go",
				"pklres_reader_test.go",
			},
			description: "Tests for schema validation and integration",
		},
		{
			name: "Performance_Benchmarks",
			testFiles: []string{
				"benchmark_test.go",
			},
			description: "Performance benchmarks and stress tests",
		},
		{
			name: "Utilities",
			testFiles: []string{
				"test_utils.go",
			},
			description: "Test utilities and helper functions",
		},
	}

	// Analyze each category
	for _, cat := range categories {
		category := analyzeTestCategory(cat.name, cat.testFiles, cat.description)
		report.TestCategories[cat.name] = category
		report.TotalTests += category.TotalTests
		report.CoveredTests += category.CoveredTests
		report.UncoveredTests += category.UncoveredTests
	}

	// Calculate overall coverage rate
	if report.TotalTests > 0 {
		report.CoverageRate = float64(report.CoveredTests) / float64(report.TotalTests) * 100
	}

	// Generate recommendations
	report.Recommendations = generateRecommendations(report)

	return report
}

// analyzeTestCategory analyzes a specific test category
func analyzeTestCategory(name string, testFiles []string, description string) TestCategory {
	category := TestCategory{
		Name:         name,
		TestFiles:    testFiles,
		MissingTests: []string{},
	}

	// Count test files that exist
	for _, testFile := range testFiles {
		if _, err := os.Stat(testFile); err == nil {
			category.TotalTests++
			category.CoveredTests++
		} else {
			category.MissingTests = append(category.MissingTests, testFile)
		}
	}

	// Calculate coverage rate
	if category.TotalTests > 0 {
		category.CoverageRate = float64(category.CoveredTests) / float64(category.TotalTests) * 100
	}

	return category
}

// generateRecommendations generates recommendations for improving test coverage
func generateRecommendations(report *TestCoverageReport) []string {
	var recommendations []string

	// Overall coverage recommendations
	if report.CoverageRate < 80 {
		recommendations = append(recommendations,
			fmt.Sprintf("Overall test coverage is %.1f%%. Aim for at least 80%% coverage.", report.CoverageRate))
	}

	// Category-specific recommendations
	for name, category := range report.TestCategories {
		if category.CoverageRate < 70 {
			recommendations = append(recommendations,
				fmt.Sprintf("Category '%s' has low coverage (%.1f%%). Consider adding more tests.", name, category.CoverageRate))
		}

		if len(category.MissingTests) > 0 {
			recommendations = append(recommendations,
				fmt.Sprintf("Category '%s' is missing test files: %s", name, strings.Join(category.MissingTests, ", ")))
		}
	}

	// Specific recommendations based on uncovered areas
	if report.UncoveredTests > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("There are %d uncovered test scenarios. Review and implement missing tests.", report.UncoveredTests))
	}

	// Performance testing recommendations
	if perfCategory, exists := report.TestCategories["Performance_Benchmarks"]; exists {
		if perfCategory.TotalTests < 3 {
			recommendations = append(recommendations,
				"Consider adding more performance benchmarks to ensure scalability.")
		}
	}

	// Integration testing recommendations
	if intCategory, exists := report.TestCategories["PKL_Integration"]; exists {
		if intCategory.TotalTests < 2 {
			recommendations = append(recommendations,
				"Add more integration tests to ensure end-to-end functionality.")
		}
	}

	return recommendations
}

// SaveCoverageReport saves the coverage report to a file
func SaveCoverageReport(report *TestCoverageReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal coverage report: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write coverage report: %v", err)
	}

	return nil
}

// PrintCoverageReport prints a formatted coverage report
func PrintCoverageReport(report *TestCoverageReport) {
	const (
		headerFormat = "%-20s %8s %8s %8s %10s"
		rowFormat    = "%-20s %8d %8d %8d %9.1f%%"
	)

	fmt.Println("🧪 Test Coverage Report")
	fmt.Println("=======================")
	fmt.Printf(headerFormat, "Category", "Total", "Covered", "Missing", "Coverage")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))

	for name, category := range report.TestCategories {
		fmt.Printf(rowFormat, name, category.TotalTests, category.CoveredTests, category.UncoveredTests, category.CoverageRate)
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf(rowFormat, "TOTAL", report.TotalTests, report.CoveredTests, report.UncoveredTests, report.CoverageRate)
	fmt.Println()

	if len(report.Recommendations) > 0 {
		fmt.Println("\n📋 Recommendations:")
		for i, rec := range report.Recommendations {
			fmt.Printf("  %d. %s\n", i+1, rec)
		}
	}
}

// TestCoverageAnalysis tests the coverage analysis functionality
func TestCoverageAnalysis(t *testing.T) {
	report := GenerateCoverageReport(t)

	// Verify report structure
	if report.TotalTests < 0 {
		t.Error("Total tests should be non-negative")
	}

	if report.CoverageRate < 0 || report.CoverageRate > 100 {
		t.Error("Coverage rate should be between 0 and 100")
	}

	// Save report
	err := SaveCoverageReport(report, "test_coverage_report.json")
	if err != nil {
		t.Errorf("Failed to save coverage report: %v", err)
	}

	// Print report
	PrintCoverageReport(report)

	t.Logf("✅ Coverage analysis completed - %d tests analyzed", report.TotalTests)
}

// AnalyzeTestFiles analyzes test files in the current directory
func AnalyzeTestFiles() ([]string, error) {
	var testFiles []string

	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "_test.go") {
			testFiles = append(testFiles, file.Name())
		}
	}

	return testFiles, nil
}

// GenerateTestTemplate generates a test template for a given category and test name
func GenerateTestTemplate(category, testName string) string {
	return fmt.Sprintf(`package test

import "testing"

func Test%s(t *testing.T) {
	// TODO: Implement test for %s
	t.Errorf("Test not implemented yet")
}
`, testName, category)
}

// CreateMissingTestFiles creates missing test files based on the coverage report
func CreateMissingTestFiles(report *TestCoverageReport) error {
	for categoryName, category := range report.TestCategories {
		for _, missingTest := range category.MissingTests {
			testName := strings.TrimSuffix(missingTest, "_test.go")
			testName = strings.Title(strings.ReplaceAll(testName, "_", " "))
			testName = strings.ReplaceAll(testName, " ", "")

			template := GenerateTestTemplate(categoryName, testName)
			filename := fmt.Sprintf("generated_%s_test.go", strings.ToLower(testName))

			if err := os.WriteFile(filename, []byte(template), 0644); err != nil {
				return fmt.Errorf("failed to create test file %s: %v", filename, err)
			}
		}
	}

	return nil
}

// CreateMissingTestFilesForAllCategories creates missing test files for all categories
func CreateMissingTestFilesForAllCategories(categories []string, testNames []string) error {
	for _, category := range categories {
		for _, testName := range testNames {
			template := GenerateTestTemplate(category, testName)
			filename := fmt.Sprintf("generated_%s_%s_test.go", strings.ToLower(category), strings.ToLower(testName))

			if err := os.WriteFile(filename, []byte(template), 0644); err != nil {
				return fmt.Errorf("failed to create test file %s: %v", filename, err)
			}
		}
	}

	return nil
}

func main() {
	flag.Parse()

	// Generate coverage report
	report := GenerateCoverageReport(nil)
	PrintCoverageReport(report)

	// Save report
	if err := SaveCoverageReport(report, "test_coverage_report.json"); err != nil {
		log.Fatalf("Failed to save coverage report: %v", err)
	}

	fmt.Println("✅ Coverage report saved to test_coverage_report.json")
}
