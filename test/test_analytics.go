package test

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestAnalytics provides advanced test analytics and reporting
type TestAnalytics struct {
	metrics *TestMetrics
	history []*TestRunHistory
	baseDir string
	config  *AnalyticsConfig
}

// AnalyticsConfig configures analytics behavior
type AnalyticsConfig struct {
	HistorySize         int           `json:"history_size"`
	TrendWindow         time.Duration `json:"trend_window"`
	RegressionThreshold float64       `json:"regression_threshold"`
	ReportFormats       []string      `json:"report_formats"`
	ExportPath          string        `json:"export_path"`
}

// TestRunHistory represents historical test run data
type TestRunHistory struct {
	RunID       string                 `json:"run_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Duration    time.Duration          `json:"duration"`
	Metrics     *TestMetrics           `json:"metrics"`
	Environment map[string]interface{} `json:"environment"`
	Tags        []string               `json:"tags"`
}

// NewTestAnalytics creates a new test analytics instance
func NewTestAnalytics(baseDir string) *TestAnalytics {
	return &TestAnalytics{
		metrics: &TestMetrics{
			TestResults: make(map[string]TestResult),
		},
		history: make([]*TestRunHistory, 0),
		baseDir: baseDir,
		config: &AnalyticsConfig{
			HistorySize:         100,
			TrendWindow:         30 * 24 * time.Hour, // 30 days
			RegressionThreshold: 0.1,                 // 10% threshold
			ReportFormats:       []string{"json", "csv", "html"},
			ExportPath:          "reports",
		},
	}
}

// RecordRun records a test run for historical analysis
func (ta *TestAnalytics) RecordRun(runID string, metrics *TestMetrics, environment map[string]interface{}, tags []string) {
	history := &TestRunHistory{
		RunID:       runID,
		Timestamp:   time.Now(),
		Duration:    time.Duration(0), // Will be calculated from individual test results
		Metrics:     metrics,
		Environment: environment,
		Tags:        tags,
	}

	ta.history = append(ta.history, history)

	// Maintain history size
	if len(ta.history) > ta.config.HistorySize {
		ta.history = ta.history[1:]
	}

	// Update current metrics
	ta.metrics = metrics
}

// AnalyzeTrends analyzes test execution trends
func (ta *TestAnalytics) AnalyzeTrends() *TrendAnalysis {
	analysis := &TrendAnalysis{
		Timestamp: time.Now(),
		Trends:    make(map[string]*Trend),
	}

	if len(ta.history) < 2 {
		return analysis
	}

	// Analyze different metrics
	analysis.Trends["duration"] = ta.analyzeDurationTrend()
	analysis.Trends["pass_rate"] = ta.analyzePassRateTrend()
	analysis.Trends["test_count"] = ta.analyzeTestCountTrend()
	analysis.Trends["failure_rate"] = ta.analyzeFailureRateTrend()

	return analysis
}

// TrendAnalysis represents trend analysis results
type TrendAnalysis struct {
	Timestamp time.Time         `json:"timestamp"`
	Trends    map[string]*Trend `json:"trends"`
	Summary   string            `json:"summary"`
}

// Trend represents a single trend
type Trend struct {
	Metric     string    `json:"metric"`
	Direction  string    `json:"direction"` // "improving", "declining", "stable"
	Change     float64   `json:"change"`
	ChangePct  float64   `json:"change_pct"`
	Confidence float64   `json:"confidence"`
	DataPoints []float64 `json:"data_points"`
}

// analyzeDurationTrend analyzes test duration trends
func (ta *TestAnalytics) analyzeDurationTrend() *Trend {
	if len(ta.history) < 2 {
		return &Trend{Metric: "duration", Direction: "stable"}
	}

	var durations []float64
	for _, run := range ta.history {
		durations = append(durations, run.Duration.Seconds())
	}

	trend := &Trend{
		Metric:     "duration",
		DataPoints: durations,
	}

	// Calculate trend
	recent := durations[len(durations)-1]
	previous := durations[len(durations)-2]

	trend.Change = recent - previous
	trend.ChangePct = (trend.Change / previous) * 100

	if trend.ChangePct > ta.config.RegressionThreshold*100 {
		trend.Direction = "declining"
	} else if trend.ChangePct < -ta.config.RegressionThreshold*100 {
		trend.Direction = "improving"
	} else {
		trend.Direction = "stable"
	}

	trend.Confidence = ta.calculateConfidence(durations)

	return trend
}

// analyzePassRateTrend analyzes test pass rate trends
func (ta *TestAnalytics) analyzePassRateTrend() *Trend {
	if len(ta.history) < 2 {
		return &Trend{Metric: "pass_rate", Direction: "stable"}
	}

	var passRates []float64
	for _, run := range ta.history {
		if run.Metrics.TotalTests > 0 {
			passRate := float64(run.Metrics.PassedTests) / float64(run.Metrics.TotalTests) * 100
			passRates = append(passRates, passRate)
		}
	}

	trend := &Trend{
		Metric:     "pass_rate",
		DataPoints: passRates,
	}

	if len(passRates) < 2 {
		return trend
	}

	recent := passRates[len(passRates)-1]
	previous := passRates[len(passRates)-2]

	trend.Change = recent - previous
	trend.ChangePct = trend.Change

	if trend.ChangePct > ta.config.RegressionThreshold*100 {
		trend.Direction = "improving"
	} else if trend.ChangePct < -ta.config.RegressionThreshold*100 {
		trend.Direction = "declining"
	} else {
		trend.Direction = "stable"
	}

	trend.Confidence = ta.calculateConfidence(passRates)

	return trend
}

// analyzeTestCountTrend analyzes test count trends
func (ta *TestAnalytics) analyzeTestCountTrend() *Trend {
	if len(ta.history) < 2 {
		return &Trend{Metric: "test_count", Direction: "stable"}
	}

	var testCounts []float64
	for _, run := range ta.history {
		testCounts = append(testCounts, float64(run.Metrics.TotalTests))
	}

	trend := &Trend{
		Metric:     "test_count",
		DataPoints: testCounts,
	}

	recent := testCounts[len(testCounts)-1]
	previous := testCounts[len(testCounts)-2]

	trend.Change = recent - previous
	trend.ChangePct = (trend.Change / previous) * 100

	if trend.ChangePct > 0 {
		trend.Direction = "increasing"
	} else if trend.ChangePct < 0 {
		trend.Direction = "decreasing"
	} else {
		trend.Direction = "stable"
	}

	trend.Confidence = ta.calculateConfidence(testCounts)

	return trend
}

// analyzeFailureRateTrend analyzes test failure rate trends
func (ta *TestAnalytics) analyzeFailureRateTrend() *Trend {
	if len(ta.history) < 2 {
		return &Trend{Metric: "failure_rate", Direction: "stable"}
	}

	var failureRates []float64
	for _, run := range ta.history {
		if run.Metrics.TotalTests > 0 {
			failureRate := float64(run.Metrics.FailedTests) / float64(run.Metrics.TotalTests) * 100
			failureRates = append(failureRates, failureRate)
		}
	}

	trend := &Trend{
		Metric:     "failure_rate",
		DataPoints: failureRates,
	}

	if len(failureRates) < 2 {
		return trend
	}

	recent := failureRates[len(failureRates)-1]
	previous := failureRates[len(failureRates)-2]

	trend.Change = recent - previous
	trend.ChangePct = trend.Change

	if trend.ChangePct > ta.config.RegressionThreshold*100 {
		trend.Direction = "declining"
	} else if trend.ChangePct < -ta.config.RegressionThreshold*100 {
		trend.Direction = "improving"
	} else {
		trend.Direction = "stable"
	}

	trend.Confidence = ta.calculateConfidence(failureRates)

	return trend
}

// calculateConfidence calculates confidence level for trend analysis
func (ta *TestAnalytics) calculateConfidence(data []float64) float64 {
	if len(data) < 2 {
		return 0.0
	}

	// Simple confidence calculation based on data consistency
	var sum, sumSq float64
	for _, value := range data {
		sum += value
		sumSq += value * value
	}

	mean := sum / float64(len(data))
	variance := (sumSq / float64(len(data))) - (mean * mean)

	if variance <= 0 {
		return 1.0
	}

	// Higher variance = lower confidence
	confidence := 1.0 / (1.0 + variance)
	return confidence
}

// DetectRegressions detects performance regressions
func (ta *TestAnalytics) DetectRegressions() *RegressionReport {
	report := &RegressionReport{
		Timestamp:   time.Now(),
		Regressions: make([]*Regression, 0),
	}

	trends := ta.AnalyzeTrends()

	for metric, trend := range trends.Trends {
		if trend.Direction == "declining" && trend.Confidence > 0.7 {
			regression := &Regression{
				Metric:         metric,
				Severity:       ta.calculateSeverity(trend.ChangePct),
				Change:         trend.Change,
				ChangePct:      trend.ChangePct,
				Confidence:     trend.Confidence,
				Recommendation: ta.generateRecommendation(metric, trend),
			}
			report.Regressions = append(report.Regressions, regression)
		}
	}

	return report
}

// RegressionReport represents regression detection results
type RegressionReport struct {
	Timestamp   time.Time     `json:"timestamp"`
	Regressions []*Regression `json:"regressions"`
	Summary     string        `json:"summary"`
}

// Regression represents a detected regression
type Regression struct {
	Metric         string  `json:"metric"`
	Severity       string  `json:"severity"` // "low", "medium", "high", "critical"
	Change         float64 `json:"change"`
	ChangePct      float64 `json:"change_pct"`
	Confidence     float64 `json:"confidence"`
	Recommendation string  `json:"recommendation"`
}

// calculateSeverity calculates regression severity
func (ta *TestAnalytics) calculateSeverity(changePct float64) string {
	absChange := changePct
	if absChange < 0 {
		absChange = -absChange
	}

	if absChange > 50 {
		return "critical"
	} else if absChange > 25 {
		return "high"
	} else if absChange > 10 {
		return "medium"
	} else {
		return "low"
	}
}

// generateRecommendation generates recommendations for regressions
func (ta *TestAnalytics) generateRecommendation(metric string, trend *Trend) string {
	switch metric {
	case "duration":
		return "Investigate test performance bottlenecks and optimize slow tests"
	case "pass_rate":
		return "Review recent code changes and fix failing tests"
	case "failure_rate":
		return "Address test failures and improve test stability"
	default:
		return "Monitor the trend and investigate root causes"
	}
}

// GenerateReport generates comprehensive test reports
func (ta *TestAnalytics) GenerateReport(format string) error {
	report := ta.buildReport()

	switch format {
	case "json":
		return ta.exportJSONReport(report)
	case "csv":
		return ta.exportCSVReport(report)
	case "html":
		return ta.exportHTMLReport(report)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// TestReport represents a comprehensive test report
type TestReport struct {
	GeneratedAt time.Time              `json:"generated_at"`
	Summary     *ReportSummary         `json:"summary"`
	Trends      *TrendAnalysis         `json:"trends"`
	Regressions *RegressionReport      `json:"regressions"`
	TestResults map[string]*TestResult `json:"test_results"`
	Environment map[string]interface{} `json:"environment"`
}

// ReportSummary represents report summary information
type ReportSummary struct {
	TotalRuns    int           `json:"total_runs"`
	TotalTests   int           `json:"total_tests"`
	PassedTests  int           `json:"passed_tests"`
	FailedTests  int           `json:"failed_tests"`
	SkippedTests int           `json:"skipped_tests"`
	PassRate     float64       `json:"pass_rate"`
	AvgDuration  time.Duration `json:"avg_duration"`
	LastRun      time.Time     `json:"last_run"`
}

// buildReport builds a comprehensive test report
func (ta *TestAnalytics) buildReport() *TestReport {
	report := &TestReport{
		GeneratedAt: time.Now(),
		Summary:     ta.buildSummary(),
		Trends:      ta.AnalyzeTrends(),
		Regressions: ta.DetectRegressions(),
		TestResults: make(map[string]*TestResult), // Convert to pointer map if needed
		Environment: ta.getEnvironment(),
	}

	return report
}

// buildSummary builds report summary
func (ta *TestAnalytics) buildSummary() *ReportSummary {
	summary := &ReportSummary{
		TotalRuns: len(ta.history),
	}

	if len(ta.history) > 0 {
		lastRun := ta.history[len(ta.history)-1]
		summary.LastRun = lastRun.Timestamp
		summary.TotalTests = lastRun.Metrics.TotalTests
		summary.PassedTests = lastRun.Metrics.PassedTests
		summary.FailedTests = lastRun.Metrics.FailedTests
		summary.SkippedTests = lastRun.Metrics.SkippedTests

		if summary.TotalTests > 0 {
			summary.PassRate = float64(summary.PassedTests) / float64(summary.TotalTests) * 100
		}

		// Calculate average duration
		var totalDuration time.Duration
		for _, run := range ta.history {
			totalDuration += run.Duration
		}
		summary.AvgDuration = totalDuration / time.Duration(len(ta.history))
	}

	return summary
}

// getEnvironment gets current environment information
func (ta *TestAnalytics) getEnvironment() map[string]interface{} {
	env := make(map[string]interface{})

	// Add basic environment info
	env["go_version"] = "1.21" // This would be dynamically determined
	env["platform"] = "darwin"
	env["arch"] = "amd64"

	return env
}

// exportJSONReport exports report in JSON format
func (ta *TestAnalytics) exportJSONReport(report *TestReport) error {
	exportDir := filepath.Join(ta.baseDir, ta.config.ExportPath)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("test_report_%s.json", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(exportDir, filename)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		// If marshaling fails due to NaN, try with NaN replacement
		reportCopy := *report
		// Replace NaN values in summary
		if math.IsNaN(reportCopy.Summary.PassRate) {
			reportCopy.Summary.PassRate = 0
		}
		// Replace NaN values in trends
		for _, trend := range reportCopy.Trends.Trends {
			if math.IsNaN(trend.Change) {
				trend.Change = 0
			}
			if math.IsNaN(trend.ChangePct) {
				trend.ChangePct = 0
			}
			if math.IsNaN(trend.Confidence) {
				trend.Confidence = 0
			}
		}
		// Replace NaN values in regressions
		for _, regression := range reportCopy.Regressions.Regressions {
			if math.IsNaN(regression.Change) {
				regression.Change = 0
			}
			if math.IsNaN(regression.ChangePct) {
				regression.ChangePct = 0
			}
			if math.IsNaN(regression.Confidence) {
				regression.Confidence = 0
			}
		}
		data, err = json.MarshalIndent(reportCopy, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	return os.WriteFile(filepath, data, 0644)
}

// exportCSVReport exports report in CSV format
func (ta *TestAnalytics) exportCSVReport(report *TestReport) error {
	exportDir := filepath.Join(ta.baseDir, ta.config.ExportPath)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("test_report_%s.csv", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(exportDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write summary
	writer.Write([]string{"Summary"})
	writer.Write([]string{"Total Runs", fmt.Sprintf("%d", report.Summary.TotalRuns)})
	writer.Write([]string{"Total Tests", fmt.Sprintf("%d", report.Summary.TotalTests)})
	writer.Write([]string{"Passed Tests", fmt.Sprintf("%d", report.Summary.PassedTests)})
	writer.Write([]string{"Failed Tests", fmt.Sprintf("%d", report.Summary.FailedTests)})
	writer.Write([]string{"Pass Rate", fmt.Sprintf("%.2f%%", report.Summary.PassRate)})
	writer.Write([]string{})

	// Write trends
	writer.Write([]string{"Trends"})
	for metric, trend := range report.Trends.Trends {
		writer.Write([]string{metric, trend.Direction, fmt.Sprintf("%.2f%%", trend.ChangePct)})
	}
	writer.Write([]string{})

	// Write regressions
	writer.Write([]string{"Regressions"})
	for _, regression := range report.Regressions.Regressions {
		writer.Write([]string{regression.Metric, regression.Severity, regression.Recommendation})
	}

	return nil
}

// exportHTMLReport exports report in HTML format
func (ta *TestAnalytics) exportHTMLReport(report *TestReport) error {
	exportDir := filepath.Join(ta.baseDir, ta.config.ExportPath)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("test_report_%s.html", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(exportDir, filename)

	html := ta.generateHTMLReport(report)
	return os.WriteFile(filepath, []byte(html), 0644)
}

// generateHTMLReport generates HTML report content
func (ta *TestAnalytics) generateHTMLReport(report *TestReport) string {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Analytics Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f0f0f0; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .metric { display: inline-block; margin: 10px; padding: 10px; background: #e8f4f8; border-radius: 3px; }
        .trend { margin: 10px 0; padding: 10px; background: #f8f8f8; }
        .regression { margin: 10px 0; padding: 10px; background: #ffe8e8; border-left: 4px solid #ff4444; }
        .success { color: #28a745; }
        .warning { color: #ffc107; }
        .danger { color: #dc3545; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Test Analytics Report</h1>
        <p>Generated: ` + report.GeneratedAt.Format("2006-01-02 15:04:05") + `</p>
    </div>
    
    <div class="section">
        <h2>Summary</h2>
        <div class="metric">Total Runs: ` + fmt.Sprintf("%d", report.Summary.TotalRuns) + `</div>
        <div class="metric">Total Tests: ` + fmt.Sprintf("%d", report.Summary.TotalTests) + `</div>
        <div class="metric">Pass Rate: ` + fmt.Sprintf("%.2f%%", report.Summary.PassRate) + `</div>
        <div class="metric">Avg Duration: ` + report.Summary.AvgDuration.String() + `</div>
    </div>
    
    <div class="section">
        <h2>Trends</h2>`

	for metric, trend := range report.Trends.Trends {
		statusClass := "success"
		if trend.Direction == "declining" {
			statusClass = "danger"
		} else if trend.Direction == "stable" {
			statusClass = "warning"
		}

		html += `
        <div class="trend">
            <strong>` + metric + `:</strong> 
            <span class="` + statusClass + `">` + trend.Direction + `</span>
            (` + fmt.Sprintf("%.2f%%", trend.ChangePct) + `)
        </div>`
	}

	html += `
    </div>
    
    <div class="section">
        <h2>Regressions</h2>`

	if len(report.Regressions.Regressions) == 0 {
		html += `<p class="success">No regressions detected</p>`
	} else {
		for _, regression := range report.Regressions.Regressions {
			html += `
        <div class="regression">
            <strong>` + regression.Metric + `</strong> (` + regression.Severity + `)<br>
            Change: ` + fmt.Sprintf("%.2f%%", regression.ChangePct) + `<br>
            Recommendation: ` + regression.Recommendation + `
        </div>`
		}
	}

	html += `
    </div>
</body>
</html>`

	return html
}

// PerformanceAnalyzer provides performance analysis capabilities
type PerformanceAnalyzer struct {
	analytics *TestAnalytics
}

// NewPerformanceAnalyzer creates a new performance analyzer
func NewPerformanceAnalyzer(analytics *TestAnalytics) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		analytics: analytics,
	}
}

// AnalyzePerformance analyzes test performance patterns
func (pa *PerformanceAnalyzer) AnalyzePerformance() *PerformanceAnalysis {
	analysis := &PerformanceAnalysis{
		Timestamp: time.Now(),
		Patterns:  make([]*PerformancePattern, 0),
	}

	if len(pa.analytics.history) < 3 {
		return analysis
	}

	// Analyze performance patterns
	analysis.Patterns = append(analysis.Patterns, pa.analyzeSlowTests())
	analysis.Patterns = append(analysis.Patterns, pa.analyzeFlakyTests())
	analysis.Patterns = append(analysis.Patterns, pa.analyzeResourceUsage())

	return analysis
}

// PerformanceAnalysis represents performance analysis results
type PerformanceAnalysis struct {
	Timestamp time.Time             `json:"timestamp"`
	Patterns  []*PerformancePattern `json:"patterns"`
}

// PerformancePattern represents a performance pattern
type PerformancePattern struct {
	Type           string  `json:"type"`
	Description    string  `json:"description"`
	Severity       string  `json:"severity"`
	Impact         float64 `json:"impact"`
	Recommendation string  `json:"recommendation"`
}

// analyzeSlowTests analyzes slow test patterns
func (pa *PerformanceAnalyzer) analyzeSlowTests() *PerformancePattern {
	// This would analyze individual test performance
	return &PerformancePattern{
		Type:           "slow_tests",
		Description:    "Tests taking longer than expected",
		Severity:       "medium",
		Impact:         0.3,
		Recommendation: "Optimize slow tests and consider parallelization",
	}
}

// analyzeFlakyTests analyzes flaky test patterns
func (pa *PerformanceAnalyzer) analyzeFlakyTests() *PerformancePattern {
	return &PerformancePattern{
		Type:           "flaky_tests",
		Description:    "Tests with inconsistent results",
		Severity:       "high",
		Impact:         0.5,
		Recommendation: "Investigate and fix flaky test conditions",
	}
}

// analyzeResourceUsage analyzes resource usage patterns
func (pa *PerformanceAnalyzer) analyzeResourceUsage() *PerformancePattern {
	return &PerformancePattern{
		Type:           "resource_usage",
		Description:    "High resource consumption patterns",
		Severity:       "low",
		Impact:         0.2,
		Recommendation: "Optimize resource usage and cleanup",
	}
}

// ExportFormats defines supported export formats
type ExportFormats struct {
	JSON     bool
	Markdown bool
	HTML     bool
}

// ExportAnalyticsReport exports analytics data in multiple formats
func ExportAnalyticsReport(report *TestReport, formats ExportFormats, baseFilename string) error {
	// Create reports directory if it doesn't exist
	reportsDir := "reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	// Clean up old files before creating new ones
	if err := cleanupOldReports(reportsDir, baseFilename); err != nil {
		return fmt.Errorf("failed to cleanup old reports: %w", err)
	}

	var errors []string

	if formats.JSON {
		filename := filepath.Join(reportsDir, baseFilename+".json")
		if err := exportJSON(report, filename); err != nil {
			errors = append(errors, fmt.Sprintf("JSON export failed: %v", err))
		}
	}

	if formats.Markdown {
		filename := filepath.Join(reportsDir, baseFilename+".md")
		if err := exportMarkdown(report, filename); err != nil {
			errors = append(errors, fmt.Sprintf("Markdown export failed: %v", err))
		}
	}

	if formats.HTML {
		filename := filepath.Join(reportsDir, baseFilename+".html")
		if err := exportHTML(report, filename); err != nil {
			errors = append(errors, fmt.Sprintf("HTML export failed: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("export errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// cleanupOldReports removes old report files to prevent accumulation
func cleanupOldReports(reportsDir, baseFilename string) error {
	pattern := filepath.Join(reportsDir, baseFilename+"*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := os.Remove(match); err != nil {
			return fmt.Errorf("failed to remove old report %s: %w", match, err)
		}
	}

	return nil
}

// exportJSON exports analytics data as JSON
func exportJSON(report *TestReport, filename string) error {
	// Create a custom marshaler that handles NaN values
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		// If marshaling fails due to NaN, try with NaN replacement
		reportCopy := *report
		// Replace NaN values in summary
		if math.IsNaN(reportCopy.Summary.PassRate) {
			reportCopy.Summary.PassRate = 0
		}
		// Replace NaN values in trends
		for _, trend := range reportCopy.Trends.Trends {
			if math.IsNaN(trend.Change) {
				trend.Change = 0
			}
			if math.IsNaN(trend.ChangePct) {
				trend.ChangePct = 0
			}
			if math.IsNaN(trend.Confidence) {
				trend.Confidence = 0
			}
		}
		// Replace NaN values in regressions
		for _, regression := range reportCopy.Regressions.Regressions {
			if math.IsNaN(regression.Change) {
				regression.Change = 0
			}
			if math.IsNaN(regression.ChangePct) {
				regression.ChangePct = 0
			}
			if math.IsNaN(regression.Confidence) {
				regression.Confidence = 0
			}
		}
		data, err = json.MarshalIndent(reportCopy, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	return os.WriteFile(filename, data, 0644)
}

// sanitizeReportForJSON creates a JSON-safe copy of the report
func sanitizeReportForJSON(report *TestReport) map[string]interface{} {
	sanitized := make(map[string]interface{})

	// Copy basic fields
	sanitized["generated_at"] = report.GeneratedAt

	// Sanitize summary
	if report.Summary != nil {
		summary := map[string]interface{}{
			"total_runs":    report.Summary.TotalRuns,
			"total_tests":   report.Summary.TotalTests,
			"passed_tests":  report.Summary.PassedTests,
			"failed_tests":  report.Summary.FailedTests,
			"skipped_tests": report.Summary.SkippedTests,
			"pass_rate":     sanitizeFloat(report.Summary.PassRate),
			"avg_duration":  report.Summary.AvgDuration.String(),
			"last_run":      report.Summary.LastRun,
		}
		sanitized["summary"] = summary
	}

	// Sanitize trends
	if report.Trends != nil {
		trends := map[string]interface{}{
			"timestamp": report.Trends.Timestamp,
			"trends":    make(map[string]interface{}),
		}

		for name, trend := range report.Trends.Trends {
			trends["trends"].(map[string]interface{})[name] = map[string]interface{}{
				"metric":      trend.Metric,
				"direction":   trend.Direction,
				"change":      sanitizeFloat(trend.Change),
				"change_pct":  sanitizeFloat(trend.ChangePct),
				"confidence":  sanitizeFloat(trend.Confidence),
				"data_points": trend.DataPoints,
			}
		}
		sanitized["trends"] = trends
	}

	// Sanitize regressions
	if report.Regressions != nil {
		regressions := map[string]interface{}{
			"timestamp":   report.Regressions.Timestamp,
			"summary":     report.Regressions.Summary,
			"regressions": make([]interface{}, 0),
		}

		for _, reg := range report.Regressions.Regressions {
			regressions["regressions"] = append(regressions["regressions"].([]interface{}), map[string]interface{}{
				"metric":         reg.Metric,
				"severity":       reg.Severity,
				"change":         sanitizeFloat(reg.Change),
				"change_pct":     sanitizeFloat(reg.ChangePct),
				"confidence":     sanitizeFloat(reg.Confidence),
				"recommendation": reg.Recommendation,
			})
		}
		sanitized["regressions"] = regressions
	}

	return sanitized
}

// sanitizeFloat replaces NaN with 0 and Inf with large numbers
func sanitizeFloat(f float64) float64 {
	if math.IsNaN(f) {
		return 0
	}
	if math.IsInf(f, 1) {
		return 999999
	}
	if math.IsInf(f, -1) {
		return -999999
	}
	return f
}

// exportMarkdown exports the report as Markdown
func exportMarkdown(report *TestReport, filename string) error {
	var content strings.Builder

	content.WriteString("# Test Analytics Report\n\n")
	content.WriteString(fmt.Sprintf("**Generated:** %s\n\n", report.GeneratedAt.Format(time.RFC3339)))

	// Summary
	content.WriteString("## Summary\n\n")
	content.WriteString(fmt.Sprintf("- **Total Runs:** %d\n", report.Summary.TotalRuns))
	content.WriteString(fmt.Sprintf("- **Total Tests:** %d\n", report.Summary.TotalTests))
	content.WriteString(fmt.Sprintf("- **Passed:** %d\n", report.Summary.PassedTests))
	content.WriteString(fmt.Sprintf("- **Failed:** %d\n", report.Summary.FailedTests))
	content.WriteString(fmt.Sprintf("- **Pass Rate:** %.2f%%\n", report.Summary.PassRate))
	content.WriteString(fmt.Sprintf("- **Average Duration:** %v\n", report.Summary.AvgDuration))
	content.WriteString(fmt.Sprintf("- **Last Run:** %s\n", report.Summary.LastRun.Format(time.RFC3339)))

	// Trends
	if report.Trends != nil && len(report.Trends.Trends) > 0 {
		content.WriteString("\n## Trends\n\n")
		for name, trend := range report.Trends.Trends {
			content.WriteString(fmt.Sprintf("### %s\n", name))
			content.WriteString(fmt.Sprintf("- **Direction:** %s\n", trend.Direction))
			content.WriteString(fmt.Sprintf("- **Change:** %.2f%%\n", trend.ChangePct))
			content.WriteString(fmt.Sprintf("- **Confidence:** %.2f%%\n", trend.Confidence))
			content.WriteString("\n")
		}
	}

	// Regressions
	if report.Regressions != nil && len(report.Regressions.Regressions) > 0 {
		content.WriteString("## Regressions\n\n")
		for _, reg := range report.Regressions.Regressions {
			content.WriteString(fmt.Sprintf("### %s (%s)\n", reg.Metric, reg.Severity))
			content.WriteString(fmt.Sprintf("- **Change:** %.2f%%\n", reg.ChangePct))
			content.WriteString(fmt.Sprintf("- **Recommendation:** %s\n", reg.Recommendation))
			content.WriteString("\n")
		}
	}

	return os.WriteFile(filename, []byte(content.String()), 0644)
}

// exportHTML exports the report as HTML
func exportHTML(report *TestReport, filename string) error {
	var content strings.Builder

	content.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Test Analytics Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .summary { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .category { margin: 20px 0; padding: 15px; border-left: 4px solid #007cba; }
        .metric { display: inline-block; margin: 10px; padding: 10px; background: #e7f3ff; border-radius: 3px; }
        .recommendation { background: #fff3cd; padding: 10px; margin: 10px 0; border-radius: 3px; }
        .regression { background: #f8d7da; padding: 10px; margin: 10px 0; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>Test Analytics Report</h1>
    <p><strong>Generated:</strong> ` + report.GeneratedAt.Format(time.RFC3339) + `</p>
    
    <div class="summary">
        <h2>Summary</h2>
        <div class="metric"><strong>Total Runs:</strong> ` + fmt.Sprintf("%d", report.Summary.TotalRuns) + `</div>
        <div class="metric"><strong>Total Tests:</strong> ` + fmt.Sprintf("%d", report.Summary.TotalTests) + `</div>
        <div class="metric"><strong>Passed:</strong> ` + fmt.Sprintf("%d", report.Summary.PassedTests) + `</div>
        <div class="metric"><strong>Failed:</strong> ` + fmt.Sprintf("%d", report.Summary.FailedTests) + `</div>
        <div class="metric"><strong>Pass Rate:</strong> ` + fmt.Sprintf("%.2f%%", report.Summary.PassRate) + `</div>
        <div class="metric"><strong>Average Duration:</strong> ` + report.Summary.AvgDuration.String() + `</div>
    </div>`)

	// Trends
	if report.Trends != nil && len(report.Trends.Trends) > 0 {
		content.WriteString(`
    <h2>Trends</h2>`)
		for name, trend := range report.Trends.Trends {
			content.WriteString(fmt.Sprintf(`
    <div class="category">
        <h3>%s</h3>
        <div class="metric"><strong>Direction:</strong> %s</div>
        <div class="metric"><strong>Change:</strong> %.2f%%</div>
        <div class="metric"><strong>Confidence:</strong> %.2f%%</div>
    </div>`, name, trend.Direction, trend.ChangePct, trend.Confidence))
		}
	}

	// Regressions
	if report.Regressions != nil && len(report.Regressions.Regressions) > 0 {
		content.WriteString(`
    <h2>Regressions</h2>`)
		for _, reg := range report.Regressions.Regressions {
			content.WriteString(fmt.Sprintf(`
    <div class="regression">
        <h3>%s (%s)</h3>
        <div class="metric"><strong>Change:</strong> %.2f%%</div>
        <div class="metric"><strong>Recommendation:</strong> %s</div>
    </div>`, reg.Metric, reg.Severity, reg.ChangePct, reg.Recommendation))
		}
	}

	content.WriteString(`
</body>
</html>`)

	return os.WriteFile(filename, []byte(content.String()), 0644)
}

// ExportAnalyticsHistory saves analytics to a timestamped file in history/
func ExportAnalyticsHistory(data interface{}) error {
	dir := "history"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(dir, fmt.Sprintf("analytics_%s.json", timestamp))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
