package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Format string constants for linter compliance
const (
	testFormat               = "\n%d. Test: %s\n"
	errorFormat              = "   Error: %s\n"
	severityFormat           = "   Severity: %s\n"
	durationFormat           = "   Duration: %v\n"
	timestampFormat          = "   Timestamp: %s\n"
	suggestionFormat         = "     %d. %s\n"
	frameFormat              = "     %s\n"
	contextFormat            = "     %s: %v\n"
	successFormat            = "✅ Test %s passed in %v\n"
	errorAnalysisFormat      = "🐛 ERROR ANALYSIS: %s\n"
	errorAnalysisErrorFormat = "Error: %s\n"
	typeFormat               = "Type: %s\n"
	rootCauseFormat          = "Root Cause: %s\n"
	impactFormat             = "Impact: %s\n"
	reportFormat             = "📊 Debug report generated: %s\n"
	stepFormat               = "  %d. %s\n"
	contextAnalysisFormat    = "  %s: %v\n"
	severityReportFormat     = "   %s: %d\n"
)

// TestDiagnostic represents diagnostic information for test failures
type TestDiagnostic struct {
	TestName      string                 `json:"test_name"`
	Error         string                 `json:"error"`
	Stack         []string               `json:"stack"`
	Context       map[string]interface{} `json:"context"`
	Timestamp     time.Time              `json:"timestamp"`
	Duration      time.Duration          `json:"duration"`
	ResourceUsage ResourceUsage          `json:"resource_usage"`
	Suggestions   []string               `json:"suggestions"`
	Severity      string                 `json:"severity"` // "low", "medium", "high", "critical"
}

// ResourceUsage tracks resource consumption during tests
type ResourceUsage struct {
	MemoryMB     float64 `json:"memory_mb"`
	CPUPercent   float64 `json:"cpu_percent"`
	DiskUsageMB  float64 `json:"disk_usage_mb"`
	NetworkCalls int     `json:"network_calls"`
	FileOps      int     `json:"file_ops"`
}

// DiagnosticManager manages test diagnostics and error analysis
type DiagnosticManager struct {
	diagnostics []*TestDiagnostic
	baseDir     string
	debugMode   bool
}

// NewDiagnosticManager creates a new diagnostic manager
func NewDiagnosticManager(baseDir string) *DiagnosticManager {
	return &DiagnosticManager{
		diagnostics: make([]*TestDiagnostic, 0),
		baseDir:     baseDir,
		debugMode:   false,
	}
}

// EnableDebugMode enables debug mode for detailed diagnostics
func (dm *DiagnosticManager) EnableDebugMode() {
	dm.debugMode = true
}

// AddDiagnostic adds a diagnostic entry
func (dm *DiagnosticManager) AddDiagnostic(testName, errorMsg string, context map[string]interface{}) *TestDiagnostic {
	diagnostic := &TestDiagnostic{
		TestName:  testName,
		Error:     errorMsg,
		Stack:     dm.getStackTrace(),
		Context:   context,
		Timestamp: time.Now(),
		Severity:  dm.analyzeSeverity(errorMsg),
	}

	// Analyze the error and generate suggestions
	diagnostic.Suggestions = dm.generateSuggestions(errorMsg, context)

	dm.diagnostics = append(dm.diagnostics, diagnostic)
	return diagnostic
}

// getStackTrace gets the current stack trace
func (dm *DiagnosticManager) getStackTrace() []string {
	var stack []string
	for i := 1; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		// Skip test framework functions
		if strings.Contains(fn.Name(), "testing.") || strings.Contains(fn.Name(), "runtime.") {
			continue
		}

		stack = append(stack, fmt.Sprintf("%s:%d %s", filepath.Base(file), line, fn.Name()))
	}
	return stack
}

// analyzeSeverity analyzes the severity of an error
func (dm *DiagnosticManager) analyzeSeverity(errorMsg string) string {
	errorMsg = strings.ToLower(errorMsg)

	// Critical errors
	if strings.Contains(errorMsg, "panic") ||
		strings.Contains(errorMsg, "segmentation fault") ||
		strings.Contains(errorMsg, "out of memory") {
		return "critical"
	}

	// High severity errors
	if strings.Contains(errorMsg, "timeout") ||
		strings.Contains(errorMsg, "connection refused") ||
		strings.Contains(errorMsg, "permission denied") {
		return "high"
	}

	// Medium severity errors
	if strings.Contains(errorMsg, "not found") ||
		strings.Contains(errorMsg, "invalid") ||
		strings.Contains(errorMsg, "failed") {
		return "medium"
	}

	// Low severity errors
	return "low"
}

// generateSuggestions generates suggestions based on the error
func (dm *DiagnosticManager) generateSuggestions(errorMsg string, context map[string]interface{}) []string {
	var suggestions []string
	errorMsg = strings.ToLower(errorMsg)

	// Common error patterns and suggestions
	if strings.Contains(errorMsg, "timeout") {
		suggestions = append(suggestions,
			"Increase test timeout in TestConfig",
			"Check network connectivity",
			"Verify external service availability")
	}

	if strings.Contains(errorMsg, "permission denied") {
		suggestions = append(suggestions,
			"Check file permissions",
			"Run with appropriate privileges",
			"Verify directory access rights")
	}

	if strings.Contains(errorMsg, "not found") {
		suggestions = append(suggestions,
			"Verify file paths and existence",
			"Check import paths",
			"Ensure dependencies are installed")
	}

	if strings.Contains(errorMsg, "invalid") {
		suggestions = append(suggestions,
			"Validate input parameters",
			"Check data format and structure",
			"Verify configuration settings")
	}

	if strings.Contains(errorMsg, "connection refused") {
		suggestions = append(suggestions,
			"Check service availability",
			"Verify network configuration",
			"Ensure firewall settings allow connections")
	}

	if strings.Contains(errorMsg, "pkl-go") || strings.Contains(errorMsg, "invalid code for maps") {
		suggestions = append(suggestions,
			"This is a known pkl-go schema parsing issue",
			"Use PKL CLI for syntax validation",
			"Wait for pkl-go library updates",
			"Focus on resource reader functionality testing")
	}

	// Context-specific suggestions
	if context != nil {
		if resourceType, ok := context["resource_type"].(string); ok {
			suggestions = append(suggestions,
				fmt.Sprintf("Check %s resource configuration", resourceType),
				fmt.Sprintf("Verify %s resource reader implementation", resourceType))
		}

		if fileName, ok := context["file_name"].(string); ok {
			suggestions = append(suggestions,
				fmt.Sprintf("Validate PKL file: %s", fileName),
				"Check PKL syntax and imports")
		}
	}

	return suggestions
}

// SaveDiagnostics saves diagnostics to a file
func (dm *DiagnosticManager) SaveDiagnostics(filename string) error {
	diagnosticsDir := filepath.Join(dm.baseDir, "diagnostics")
	if err := os.MkdirAll(diagnosticsDir, 0755); err != nil {
		return fmt.Errorf("failed to create diagnostics directory: %v", err)
	}

	data, err := json.MarshalIndent(dm.diagnostics, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal diagnostics: %v", err)
	}

	filePath := filepath.Join(diagnosticsDir, filename)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write diagnostics: %v", err)
	}

	return nil
}

// PrintDiagnostics prints diagnostics to console
func (dm *DiagnosticManager) PrintDiagnostics() {
	if len(dm.diagnostics) == 0 {
		fmt.Printf("No diagnostics to report.\n")
		return
	}

	const headerLine = "================================================================================"
	const reportTitle = "🔍 TEST DIAGNOSTICS REPORT"

	fmt.Printf("\n%s\n", headerLine)
	fmt.Printf("%s\n", reportTitle)
	fmt.Printf("%s\n", headerLine)

	for i, diagnostic := range dm.diagnostics {
		fmt.Printf("\n%d. Test: %s\n", i+1, diagnostic.TestName)
		fmt.Printf("   Error: %s\n", diagnostic.Error)
		fmt.Printf("   Severity: %s\n", diagnostic.Severity)
		fmt.Printf("   Duration: %v\n", diagnostic.Duration)
		fmt.Printf("   Timestamp: %s\n", diagnostic.Timestamp.Format(time.RFC3339))

		if len(diagnostic.Suggestions) > 0 {
			fmt.Printf("   Suggestions:\n")
			for j, suggestion := range diagnostic.Suggestions {
				fmt.Printf("     %d. %s\n", j+1, suggestion)
			}
		}

		if dm.debugMode && len(diagnostic.Stack) > 0 {
			fmt.Printf("   Stack Trace:\n")
			for _, frame := range diagnostic.Stack {
				fmt.Printf("     %s\n", frame)
			}
		}

		if diagnostic.Context != nil {
			fmt.Printf("   Context:\n")
			for key, value := range diagnostic.Context {
				fmt.Printf("     %s: %v\n", key, value)
			}
		}
	}

	fmt.Printf("\n%s\n", headerLine)
}

// GetDiagnosticsBySeverity returns diagnostics filtered by severity
func (dm *DiagnosticManager) GetDiagnosticsBySeverity(severity string) []*TestDiagnostic {
	var filtered []*TestDiagnostic
	for _, diagnostic := range dm.diagnostics {
		if diagnostic.Severity == severity {
			filtered = append(filtered, diagnostic)
		}
	}
	return filtered
}

// GetDiagnosticsByTest returns diagnostics for a specific test
func (dm *DiagnosticManager) GetDiagnosticsByTest(testName string) []*TestDiagnostic {
	var filtered []*TestDiagnostic
	for _, diagnostic := range dm.diagnostics {
		if diagnostic.TestName == testName {
			filtered = append(filtered, diagnostic)
		}
	}
	return filtered
}

// ErrorAnalyzer provides advanced error analysis
type ErrorAnalyzer struct {
	diagnosticManager *DiagnosticManager
}

// NewErrorAnalyzer creates a new error analyzer
func NewErrorAnalyzer(diagnosticManager *DiagnosticManager) *ErrorAnalyzer {
	return &ErrorAnalyzer{
		diagnosticManager: diagnosticManager,
	}
}

// AnalyzeError performs detailed error analysis
func (ea *ErrorAnalyzer) AnalyzeError(testName string, err error, context map[string]interface{}) *ErrorAnalysis {
	analysis := &ErrorAnalysis{
		TestName:   testName,
		Error:      err.Error(),
		ErrorType:  ea.classifyError(err),
		RootCause:  ea.identifyRootCause(err),
		Impact:     ea.assessImpact(err),
		Resolution: ea.suggestResolution(err),
		Context:    context,
		Timestamp:  time.Now(),
	}

	// Add diagnostic entry
	ea.diagnosticManager.AddDiagnostic(testName, err.Error(), context)

	return analysis
}

// ErrorAnalysis represents detailed error analysis
type ErrorAnalysis struct {
	TestName   string                 `json:"test_name"`
	Error      string                 `json:"error"`
	ErrorType  string                 `json:"error_type"`
	RootCause  string                 `json:"root_cause"`
	Impact     string                 `json:"impact"`
	Resolution []string               `json:"resolution"`
	Context    map[string]interface{} `json:"context"`
	Timestamp  time.Time              `json:"timestamp"`
}

// classifyError classifies the type of error
func (ea *ErrorAnalyzer) classifyError(err error) string {
	errorMsg := strings.ToLower(err.Error())

	if strings.Contains(errorMsg, "timeout") {
		return "timeout"
	}
	if strings.Contains(errorMsg, "permission") {
		return "permission"
	}
	if strings.Contains(errorMsg, "not found") {
		return "not_found"
	}
	if strings.Contains(errorMsg, "invalid") {
		return "invalid"
	}
	if strings.Contains(errorMsg, "connection") {
		return "connection"
	}
	if strings.Contains(errorMsg, "pkl-go") {
		return "pkl_go_schema_parsing"
	}

	return "unknown"
}

// identifyRootCause identifies the root cause of the error
func (ea *ErrorAnalyzer) identifyRootCause(err error) string {
	errorMsg := strings.ToLower(err.Error())

	if strings.Contains(errorMsg, "timeout") {
		return "Network or service timeout"
	}
	if strings.Contains(errorMsg, "permission") {
		return "Insufficient file or system permissions"
	}
	if strings.Contains(errorMsg, "not found") {
		return "Missing file or resource"
	}
	if strings.Contains(errorMsg, "invalid") {
		return "Invalid input or configuration"
	}
	if strings.Contains(errorMsg, "connection") {
		return "Network connectivity issue"
	}
	if strings.Contains(errorMsg, "pkl-go") {
		return "pkl-go library schema parsing issue"
	}

	return "Unknown root cause"
}

// assessImpact assesses the impact of the error
func (ea *ErrorAnalyzer) assessImpact(err error) string {
	errorMsg := strings.ToLower(err.Error())

	if strings.Contains(errorMsg, "panic") || strings.Contains(errorMsg, "segmentation fault") {
		return "Critical - Test execution halted"
	}
	if strings.Contains(errorMsg, "timeout") || strings.Contains(errorMsg, "connection") {
		return "High - Test functionality impaired"
	}
	if strings.Contains(errorMsg, "not found") || strings.Contains(errorMsg, "invalid") {
		return "Medium - Test may work with corrections"
	}

	return "Low - Minor issue, test may still pass"
}

// suggestResolution suggests resolution steps
func (ea *ErrorAnalyzer) suggestResolution(err error) []string {
	return ea.diagnosticManager.generateSuggestions(err.Error(), nil)
}

// DebugHelper provides debugging assistance
type DebugHelper struct {
	diagnosticManager *DiagnosticManager
	errorAnalyzer     *ErrorAnalyzer
}

// NewDebugHelper creates a new debug helper
func NewDebugHelper(diagnosticManager *DiagnosticManager) *DebugHelper {
	return &DebugHelper{
		diagnosticManager: diagnosticManager,
		errorAnalyzer:     NewErrorAnalyzer(diagnosticManager),
	}
}

// DebugTest provides debugging information for a test
func (dh *DebugHelper) DebugTest(testName string, testFunc func() error) {
	start := time.Now()

	// Enable debug mode
	dh.diagnosticManager.EnableDebugMode()

	// Run test with error analysis
	err := testFunc()

	duration := time.Since(start)

	if err != nil {
		// Analyze error
		analysis := dh.errorAnalyzer.AnalyzeError(testName, err, map[string]interface{}{
			"duration":   duration,
			"debug_mode": true,
		})

		// Print detailed analysis
		dh.printErrorAnalysis(analysis)
	} else {
		fmt.Printf("✅ Test %s passed in %v\n", testName, duration)
	}
}

// printErrorAnalysis prints detailed error analysis
func (dh *DebugHelper) printErrorAnalysis(analysis *ErrorAnalysis) {
	const headerLine = "================================================================================"
	fmt.Printf("\n%s\n", headerLine)
	fmt.Printf("🐛 ERROR ANALYSIS: %s\n", analysis.TestName)
	fmt.Printf("%s\n", headerLine)
	fmt.Printf("Error: %s\n", analysis.Error)
	fmt.Printf("Type: %s\n", analysis.ErrorType)
	fmt.Printf("Root Cause: %s\n", analysis.RootCause)
	fmt.Printf("Impact: %s\n", analysis.Impact)
	fmt.Printf("Timestamp: %s\n", analysis.Timestamp.Format(time.RFC3339))

	if len(analysis.Resolution) > 0 {
		fmt.Printf("\nResolution Steps:\n")
		for i, step := range analysis.Resolution {
			fmt.Printf("  %d. %s\n", i+1, step)
		}
	}

	if analysis.Context != nil {
		fmt.Printf("\nContext:\n")
		for key, value := range analysis.Context {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Printf("\n%s\n", headerLine)
}

// GenerateDebugReport generates a comprehensive debug report
func (dh *DebugHelper) GenerateDebugReport(filename string) error {
	// Save diagnostics
	if err := dh.diagnosticManager.SaveDiagnostics(filename); err != nil {
		return err
	}

	// Print summary
	fmt.Printf("📊 Debug report generated: %s\n", filename)
	fmt.Printf("📋 Total diagnostics: %d\n", len(dh.diagnosticManager.diagnostics))

	// Print severity breakdown
	severities := []string{"critical", "high", "medium", "low"}
	for _, severity := range severities {
		count := len(dh.diagnosticManager.GetDiagnosticsBySeverity(severity))
		if count > 0 {
			fmt.Printf("   %s: %d\n", severity, count)
		}
	}

	return nil
}
