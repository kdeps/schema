package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// CoverageAnalyzer provides comprehensive test coverage analysis
type CoverageAnalyzer struct {
	coverage map[string]*FileCoverage
	summary  *CoverageSummary
	history  []*CoverageReport
	config   *CoverageConfig
	baseDir  string
	mu       sync.RWMutex
}

// CoverageConfig configures the coverage analyzer
type CoverageConfig struct {
	MinCoverage   float64  `json:"min_coverage"`
	CriticalPaths []string `json:"critical_paths"`
	ExcludePaths  []string `json:"exclude_paths"`
	IncludePaths  []string `json:"include_paths"`
	TrackHistory  bool     `json:"track_history"`
	HistorySize   int      `json:"history_size"`
	ExportPath    string   `json:"export_path"`
	GenerateHTML  bool     `json:"generate_html"`
	GenerateBadge bool     `json:"generate_badge"`
}

// FileCoverage represents coverage data for a single file
type FileCoverage struct {
	Path           string                       `json:"path"`
	TotalLines     int                          `json:"total_lines"`
	CoveredLines   int                          `json:"covered_lines"`
	UncoveredLines int                          `json:"uncovered_lines"`
	Coverage       float64                      `json:"coverage"`
	Functions      map[string]*FunctionCoverage `json:"functions"`
	Branches       map[string]*BranchCoverage   `json:"branches"`
	Lines          map[int]*LineCoverage        `json:"lines"`
	LastModified   time.Time                    `json:"last_modified"`
	Complexity     int                          `json:"complexity"`
	Risk           string                       `json:"risk"` // "low", "medium", "high", "critical"
}

// FunctionCoverage represents function-level coverage
type FunctionCoverage struct {
	Name         string  `json:"name"`
	StartLine    int     `json:"start_line"`
	EndLine      int     `json:"end_line"`
	TotalLines   int     `json:"total_lines"`
	CoveredLines int     `json:"covered_lines"`
	Coverage     float64 `json:"coverage"`
	Complexity   int     `json:"complexity"`
	Risk         string  `json:"risk"`
}

// BranchCoverage represents branch-level coverage
type BranchCoverage struct {
	Line            int     `json:"line"`
	TotalBranches   int     `json:"total_branches"`
	CoveredBranches int     `json:"covered_branches"`
	Coverage        float64 `json:"coverage"`
	Risk            string  `json:"risk"`
}

// LineCoverage represents line-level coverage
type LineCoverage struct {
	Line     int    `json:"line"`
	Covered  bool   `json:"covered"`
	Hits     int    `json:"hits"`
	Function string `json:"function"`
	Branch   bool   `json:"branch"`
	Risk     string `json:"risk"`
}

// CoverageSummary provides a summary of coverage data
type CoverageSummary struct {
	Timestamp        time.Time          `json:"timestamp"`
	TotalFiles       int                `json:"total_files"`
	CoveredFiles     int                `json:"covered_files"`
	UncoveredFiles   int                `json:"uncovered_files"`
	TotalLines       int                `json:"total_lines"`
	CoveredLines     int                `json:"covered_lines"`
	UncoveredLines   int                `json:"uncovered_lines"`
	TotalFunctions   int                `json:"total_functions"`
	CoveredFunctions int                `json:"covered_functions"`
	TotalBranches    int                `json:"total_branches"`
	CoveredBranches  int                `json:"covered_branches"`
	OverallCoverage  float64            `json:"overall_coverage"`
	FunctionCoverage float64            `json:"function_coverage"`
	BranchCoverage   float64            `json:"branch_coverage"`
	RiskBreakdown    map[string]int     `json:"risk_breakdown"`
	CriticalPaths    map[string]float64 `json:"critical_paths"`
	Trends           *CoverageTrends    `json:"trends"`
	Recommendations  []string           `json:"recommendations"`
}

// CoverageTrends represents coverage trends over time
type CoverageTrends struct {
	OverallTrend      string  `json:"overall_trend"` // "improving", "declining", "stable"
	ChangeRate        float64 `json:"change_rate"`
	DaysSinceLast     int     `json:"days_since_last"`
	PredictedCoverage float64 `json:"predicted_coverage"`
}

// CoverageReport represents a complete coverage report
type CoverageReport struct {
	ID          string                   `json:"id"`
	Timestamp   time.Time                `json:"timestamp"`
	Summary     *CoverageSummary         `json:"summary"`
	Files       map[string]*FileCoverage `json:"files"`
	Environment map[string]interface{}   `json:"environment"`
}

// NewCoverageAnalyzer creates a new coverage analyzer
func NewCoverageAnalyzer(baseDir string) *CoverageAnalyzer {
	return &CoverageAnalyzer{
		coverage: make(map[string]*FileCoverage),
		history:  make([]*CoverageReport, 0),
		config: &CoverageConfig{
			MinCoverage:   80.0,
			CriticalPaths: []string{"pkg/", "internal/", "cmd/"},
			ExcludePaths:  []string{"vendor/", "test/", "docs/", ".git/"},
			TrackHistory:  true,
			HistorySize:   30,
			ExportPath:    "coverage",
			GenerateHTML:  true,
			GenerateBadge: true,
		},
		baseDir: baseDir,
	}
}

// AddFileCoverage adds coverage data for a file
func (ca *CoverageAnalyzer) AddFileCoverage(path string, coverage *FileCoverage) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	ca.coverage[path] = coverage
}

// SetMinCoverage sets the minimum coverage threshold
func (ca *CoverageAnalyzer) SetMinCoverage(threshold float64) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	ca.config.MinCoverage = threshold
}

// AddCriticalPath adds a critical path for coverage analysis
func (ca *CoverageAnalyzer) AddCriticalPath(path string) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	ca.config.CriticalPaths = append(ca.config.CriticalPaths, path)
}

// AddExcludePath adds a path to exclude from coverage analysis
func (ca *CoverageAnalyzer) AddExcludePath(path string) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	ca.config.ExcludePaths = append(ca.config.ExcludePaths, path)
}

// AnalyzeCoverage analyzes the current coverage data
func (ca *CoverageAnalyzer) AnalyzeCoverage() *CoverageReport {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	report := &CoverageReport{
		ID:        fmt.Sprintf("coverage_%d", time.Now().Unix()),
		Timestamp: time.Now(),
		Files:     make(map[string]*FileCoverage),
		Environment: map[string]interface{}{
			"go_version": "1.24",
			"platform":   "darwin",
			"timestamp":  time.Now().Format(time.RFC3339),
		},
	}

	// Copy coverage data
	for path, fileCoverage := range ca.coverage {
		report.Files[path] = fileCoverage
	}

	// Generate summary
	report.Summary = ca.generateSummary(report.Files)

	// Store in history
	if ca.config.TrackHistory {
		ca.history = append(ca.history, report)
		if len(ca.history) > ca.config.HistorySize {
			ca.history = ca.history[1:]
		}
	}

	return report
}

// generateSummary generates a coverage summary
func (ca *CoverageAnalyzer) generateSummary(files map[string]*FileCoverage) *CoverageSummary {
	summary := &CoverageSummary{
		Timestamp:     time.Now(),
		RiskBreakdown: make(map[string]int),
		CriticalPaths: make(map[string]float64),
	}

	// Calculate totals
	for path, fileCoverage := range files {
		summary.TotalFiles++
		summary.TotalLines += fileCoverage.TotalLines
		summary.CoveredLines += fileCoverage.CoveredLines
		summary.UncoveredLines += fileCoverage.UncoveredLines

		if fileCoverage.Coverage > 0 {
			summary.CoveredFiles++
		} else {
			summary.UncoveredFiles++
		}

		// Count functions
		for _, function := range fileCoverage.Functions {
			summary.TotalFunctions++
			if function.Coverage > 0 {
				summary.CoveredFunctions++
			}
		}

		// Count branches
		for _, branch := range fileCoverage.Branches {
			summary.TotalBranches++
			if branch.Coverage > 0 {
				summary.CoveredBranches++
			}
		}

		// Risk breakdown
		summary.RiskBreakdown[fileCoverage.Risk]++

		// Critical paths
		for _, criticalPath := range ca.config.CriticalPaths {
			if strings.HasPrefix(path, criticalPath) {
				if _, exists := summary.CriticalPaths[criticalPath]; !exists {
					summary.CriticalPaths[criticalPath] = 0
				}
				summary.CriticalPaths[criticalPath] += fileCoverage.Coverage
			}
		}
	}

	// Calculate coverage percentages
	if summary.TotalLines > 0 {
		summary.OverallCoverage = float64(summary.CoveredLines) / float64(summary.TotalLines) * 100
	}
	if summary.TotalFunctions > 0 {
		summary.FunctionCoverage = float64(summary.CoveredFunctions) / float64(summary.TotalFunctions) * 100
	}
	if summary.TotalBranches > 0 {
		summary.BranchCoverage = float64(summary.CoveredBranches) / float64(summary.TotalBranches) * 100
	}

	// Calculate trends
	summary.Trends = ca.calculateTrends()

	// Generate recommendations
	summary.Recommendations = ca.generateRecommendations(summary)

	return summary
}

// calculateTrends calculates coverage trends
func (ca *CoverageAnalyzer) calculateTrends() *CoverageTrends {
	if len(ca.history) < 2 {
		return &CoverageTrends{
			OverallTrend:  "stable",
			ChangeRate:    0,
			DaysSinceLast: 0,
		}
	}

	// Get last two reports
	last := ca.history[len(ca.history)-1]
	previous := ca.history[len(ca.history)-2]

	change := last.Summary.OverallCoverage - previous.Summary.OverallCoverage
	daysDiff := int(last.Timestamp.Sub(previous.Timestamp).Hours() / 24)

	trend := "stable"
	if change > 1.0 {
		trend = "improving"
	} else if change < -1.0 {
		trend = "declining"
	}

	changeRate := change / float64(daysDiff)
	if daysDiff == 0 {
		changeRate = 0
	}

	// Predict future coverage
	predictedCoverage := last.Summary.OverallCoverage + (changeRate * 7) // 7 days prediction

	return &CoverageTrends{
		OverallTrend:      trend,
		ChangeRate:        changeRate,
		DaysSinceLast:     daysDiff,
		PredictedCoverage: predictedCoverage,
	}
}

// generateRecommendations generates coverage recommendations
func (ca *CoverageAnalyzer) generateRecommendations(summary *CoverageSummary) []string {
	var recommendations []string

	// Overall coverage recommendations
	if summary.OverallCoverage < ca.config.MinCoverage {
		recommendations = append(recommendations,
			fmt.Sprintf("Overall coverage (%.1f%%) is below minimum threshold (%.1f%%)",
				summary.OverallCoverage, ca.config.MinCoverage))
	}

	// Function coverage recommendations
	if summary.FunctionCoverage < 90.0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Function coverage (%.1f%%) is below recommended threshold (90%%)",
				summary.FunctionCoverage))
	}

	// Branch coverage recommendations
	if summary.BranchCoverage < 80.0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Branch coverage (%.1f%%) is below recommended threshold (80%%)",
				summary.BranchCoverage))
	}

	// Critical paths recommendations
	for path, coverage := range summary.CriticalPaths {
		if coverage < ca.config.MinCoverage {
			recommendations = append(recommendations,
				fmt.Sprintf("Critical path '%s' has low coverage (%.1f%%)", path, coverage))
		}
	}

	// Risk-based recommendations
	if summary.RiskBreakdown["critical"] > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Found %d files with critical risk - prioritize coverage for these files",
				summary.RiskBreakdown["critical"]))
	}

	if summary.RiskBreakdown["high"] > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Found %d files with high risk - consider adding more tests",
				summary.RiskBreakdown["high"]))
	}

	// Uncovered files recommendations
	if summary.UncoveredFiles > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Found %d uncovered files - consider adding tests for these files",
				summary.UncoveredFiles))
	}

	return recommendations
}

// GetUncoveredFiles returns files with no coverage
func (ca *CoverageAnalyzer) GetUncoveredFiles() []string {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	var uncovered []string
	for path, coverage := range ca.coverage {
		if coverage.Coverage == 0 {
			uncovered = append(uncovered, path)
		}
	}
	return uncovered
}

// GetLowCoverageFiles returns files with coverage below threshold
func (ca *CoverageAnalyzer) GetLowCoverageFiles(threshold float64) []string {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	var lowCoverage []string
	for path, coverage := range ca.coverage {
		if coverage.Coverage < threshold {
			lowCoverage = append(lowCoverage, path)
		}
	}
	return lowCoverage
}

// GetCriticalFiles returns files in critical paths
func (ca *CoverageAnalyzer) GetCriticalFiles() []string {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	var critical []string
	for path := range ca.coverage {
		for _, criticalPath := range ca.config.CriticalPaths {
			if strings.HasPrefix(path, criticalPath) {
				critical = append(critical, path)
				break
			}
		}
	}
	return critical
}

// GetHighRiskFiles returns files with high or critical risk
func (ca *CoverageAnalyzer) GetHighRiskFiles() []string {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	var highRisk []string
	for path, coverage := range ca.coverage {
		if coverage.Risk == "high" || coverage.Risk == "critical" {
			highRisk = append(highRisk, path)
		}
	}
	return highRisk
}

// ExportReport exports the coverage report
func (ca *CoverageAnalyzer) ExportReport(report *CoverageReport, format string) error {
	exportDir := filepath.Join(ca.baseDir, ca.config.ExportPath)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("coverage_report_%s.%s", report.ID, format)
	filepath := filepath.Join(exportDir, filename)

	switch format {
	case "json":
		return ca.exportJSON(report, filepath)
	case "html":
		return ca.exportHTML(report, filepath)
	case "csv":
		return ca.exportCSV(report, filepath)
	case "badge":
		return ca.exportBadge(report, filepath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// exportJSON exports report as JSON
func (ca *CoverageAnalyzer) exportJSON(report *CoverageReport, filepath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

// exportHTML exports report as HTML
func (ca *CoverageAnalyzer) exportHTML(report *CoverageReport, filepath string) error {
	html := ca.generateHTMLReport(report)
	return os.WriteFile(filepath, []byte(html), 0644)
}

// exportCSV exports report as CSV
func (ca *CoverageAnalyzer) exportCSV(report *CoverageReport, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	file.WriteString("File,Total Lines,Covered Lines,Uncovered Lines,Coverage %,Functions,Branches,Risk\n")

	// Write data
	for path, coverage := range report.Files {
		line := fmt.Sprintf("%s,%d,%d,%d,%.2f,%d,%d,%s\n",
			path,
			coverage.TotalLines,
			coverage.CoveredLines,
			coverage.UncoveredLines,
			coverage.Coverage,
			len(coverage.Functions),
			len(coverage.Branches),
			coverage.Risk,
		)
		file.WriteString(line)
	}

	return nil
}

// exportBadge exports coverage badge
func (ca *CoverageAnalyzer) exportBadge(report *CoverageReport, filepath string) error {
	coverage := report.Summary.OverallCoverage

	var color string

	if coverage >= 90 {
		color = "brightgreen"
	} else if coverage >= 80 {
		color = "green"
	} else if coverage >= 70 {
		color = "yellow"
	} else if coverage >= 60 {
		color = "orange"
	} else {
		color = "red"
	}

	badge := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="20">
  <linearGradient id="b" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="200" height="20" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <path fill="#555" d="M0 0h120v20H0z"/>
    <path fill="%s" d="M120 0h80v20h-80z"/>
    <path fill="url(#b)" d="M0 0h200v20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="60" y="15" fill="#010101" fill-opacity=".3">coverage</text>
    <text x="60" y="14">coverage</text>
    <text x="160" y="15" fill="#010101" fill-opacity=".3">%.1f%%</text>
    <text x="160" y="14">%.1f%%</text>
  </g>
</svg>`, color, coverage, coverage)

	return os.WriteFile(filepath, []byte(badge), 0644)
}

// generateHTMLReport generates HTML report content
func (ca *CoverageAnalyzer) generateHTMLReport(report *CoverageReport) string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Coverage Report - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .summary { background: #e7f3ff; padding: 15px; margin: 20px 0; border-radius: 5px; }
        .metric { display: inline-block; margin: 10px; padding: 10px; background: #f8f9fa; border-radius: 3px; }
        .file { margin: 10px 0; padding: 10px; border: 1px solid #ddd; border-radius: 3px; }
        .coverage-bar { width: 100%%; height: 20px; background: #f0f0f0; border-radius: 3px; overflow: hidden; }
        .coverage-fill { height: 100%%; background: linear-gradient(90deg, #ff4444 0%%, #ffaa00 50%%, #44ff44 100%%); }
        .risk-critical { border-left: 4px solid #dc3545; }
        .risk-high { border-left: 4px solid #fd7e14; }
        .risk-medium { border-left: 4px solid #ffc107; }
        .risk-low { border-left: 4px solid #28a745; }
        .recommendations { background: #fff3cd; padding: 15px; margin: 20px 0; border-radius: 5px; }
        table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Coverage Report</h1>
        <p><strong>Generated:</strong> %s</p>
        <p><strong>Overall Coverage:</strong> %.2f%%</p>
    </div>
    
    <div class="summary">
        <h2>Summary</h2>
        <div class="metric"><strong>Total Files:</strong> %d</div>
        <div class="metric"><strong>Covered Files:</strong> %d</div>
        <div class="metric"><strong>Uncovered Files:</strong> %d</div>
        <div class="metric"><strong>Total Lines:</strong> %d</div>
        <div class="metric"><strong>Covered Lines:</strong> %d</div>
        <div class="metric"><strong>Function Coverage:</strong> %.2f%%</div>
        <div class="metric"><strong>Branch Coverage:</strong> %.2f%%</div>
    </div>`,
		report.ID, report.Timestamp.Format(time.RFC3339), report.Summary.OverallCoverage,
		report.Summary.TotalFiles, report.Summary.CoveredFiles, report.Summary.UncoveredFiles,
		report.Summary.TotalLines, report.Summary.CoveredLines,
		report.Summary.FunctionCoverage, report.Summary.BranchCoverage)

	// Add recommendations
	if len(report.Summary.Recommendations) > 0 {
		html += `<div class="recommendations"><h3>Recommendations</h3><ul>`
		for _, rec := range report.Summary.Recommendations {
			html += fmt.Sprintf("<li>%s</li>", rec)
		}
		html += `</ul></div>`
	}

	// Add file details
	html += `<h2>File Details</h2>`
	for path, coverage := range report.Files {
		riskClass := fmt.Sprintf("risk-%s", coverage.Risk)
		html += fmt.Sprintf(`
    <div class="file %s">
        <h3>%s</h3>
        <div class="metric">Coverage: %.2f%%</div>
        <div class="metric">Lines: %d/%d</div>
        <div class="metric">Functions: %d</div>
        <div class="metric">Branches: %d</div>
        <div class="metric">Risk: %s</div>
        <div class="coverage-bar">
            <div class="coverage-fill" style="width: %.2f%%"></div>
        </div>
    </div>`,
			riskClass, path, coverage.Coverage,
			coverage.CoveredLines, coverage.TotalLines,
			len(coverage.Functions), len(coverage.Branches),
			coverage.Risk, coverage.Coverage)
	}

	html += `</body></html>`
	return html
}

// GetHistory returns coverage history
func (ca *CoverageAnalyzer) GetHistory() []*CoverageReport {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	history := make([]*CoverageReport, len(ca.history))
	copy(history, ca.history)
	return history
}

// CalculateRisk calculates risk level based on coverage and complexity
func (ca *CoverageAnalyzer) CalculateRisk(coverage float64, complexity int) string {
	if coverage < 50 || complexity > 10 {
		return "critical"
	} else if coverage < 70 || complexity > 7 {
		return "high"
	} else if coverage < 85 || complexity > 5 {
		return "medium"
	} else {
		return "low"
	}
}
