package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// BenchmarkSystem provides comprehensive benchmarking capabilities
type BenchmarkSystem struct {
	benchmarks map[string]*BenchmarkSuite
	history    []*BenchmarkRun
	config     *BenchmarkConfig
	mu         sync.RWMutex
	baseDir    string
}

// BenchmarkConfig configures the benchmark system
type BenchmarkConfig struct {
	WarmupRuns     int           `json:"warmup_runs"`
	BenchmarkRuns  int           `json:"benchmark_runs"`
	Timeout        time.Duration `json:"timeout"`
	MemoryTracking bool          `json:"memory_tracking"`
	CPUProfiling   bool          `json:"cpu_profiling"`
	ExportPath     string        `json:"export_path"`
}

// BenchmarkSuite represents a collection of benchmarks
type BenchmarkSuite struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Benchmarks  map[string]*Benchmark  `json:"benchmarks"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Benchmark represents a single benchmark test
type Benchmark struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Function    func() error           `json:"-"`
	Setup       func() error           `json:"-"`
	Teardown    func() error           `json:"-"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// BenchmarkResult represents the result of a benchmark run
type BenchmarkResult struct {
	Name         string        `json:"name"`
	Duration     time.Duration `json:"duration"`
	MemoryUsage  uint64        `json:"memory_usage"`
	CPUUsage     float64       `json:"cpu_usage"`
	Iterations   int           `json:"iterations"`
	Operations   int           `json:"operations"`
	OpsPerSecond float64       `json:"ops_per_second"`
	MemoryPerOp  float64       `json:"memory_per_op"`
	Error        error         `json:"error,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
	Percentiles  Percentiles   `json:"percentiles"`
	Statistics   Statistics    `json:"statistics"`
}

// Percentiles represents percentile data
type Percentiles struct {
	P50  time.Duration `json:"p50"`
	P90  time.Duration `json:"p90"`
	P95  time.Duration `json:"p95"`
	P99  time.Duration `json:"p99"`
	P999 time.Duration `json:"p999"`
}

// Statistics represents statistical data
type Statistics struct {
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	StdDev float64 `json:"std_dev"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

// BenchmarkRun represents a complete benchmark run
type BenchmarkRun struct {
	ID          string                      `json:"id"`
	Timestamp   time.Time                   `json:"timestamp"`
	SuiteName   string                      `json:"suite_name"`
	Results     map[string]*BenchmarkResult `json:"results"`
	Environment map[string]interface{}      `json:"environment"`
	Summary     *BenchmarkSummary           `json:"summary"`
}

// BenchmarkSummary provides a summary of benchmark results
type BenchmarkSummary struct {
	TotalBenchmarks int           `json:"total_benchmarks"`
	Passed          int           `json:"passed"`
	Failed          int           `json:"failed"`
	TotalDuration   time.Duration `json:"total_duration"`
	AvgOpsPerSec    float64       `json:"avg_ops_per_sec"`
	BestPerformer   string        `json:"best_performer"`
	WorstPerformer  string        `json:"worst_performer"`
}

// NewBenchmarkSystem creates a new benchmark system
func NewBenchmarkSystem(baseDir string) *BenchmarkSystem {
	return &BenchmarkSystem{
		benchmarks: make(map[string]*BenchmarkSuite),
		history:    make([]*BenchmarkRun, 0),
		config: &BenchmarkConfig{
			WarmupRuns:     3,
			BenchmarkRuns:  10,
			Timeout:        30 * time.Second,
			MemoryTracking: true,
			CPUProfiling:   false,
			ExportPath:     "benchmarks",
		},
		baseDir: baseDir,
	}
}

// AddSuite adds a benchmark suite to the system
func (bs *BenchmarkSystem) AddSuite(name, description string) *BenchmarkSuite {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	suite := &BenchmarkSuite{
		Name:        name,
		Description: description,
		Benchmarks:  make(map[string]*Benchmark),
		Metadata:    make(map[string]interface{}),
	}
	bs.benchmarks[name] = suite
	return suite
}

// AddBenchmark adds a benchmark to a suite
func (suite *BenchmarkSuite) AddBenchmark(name, description string, fn func() error) *Benchmark {
	benchmark := &Benchmark{
		Name:        name,
		Description: description,
		Function:    fn,
		Metadata:    make(map[string]interface{}),
	}
	suite.Benchmarks[name] = benchmark
	return benchmark
}

// SetSetup sets the setup function for a benchmark
func (b *Benchmark) SetSetup(setup func() error) *Benchmark {
	b.Setup = setup
	return b
}

// SetTeardown sets the teardown function for a benchmark
func (b *Benchmark) SetTeardown(teardown func() error) *Benchmark {
	b.Teardown = teardown
	return b
}

// AddMetadata adds metadata to a benchmark
func (b *Benchmark) AddMetadata(key string, value interface{}) *Benchmark {
	b.Metadata[key] = value
	return b
}

// RunSuite runs all benchmarks in a suite
func (bs *BenchmarkSystem) RunSuite(ctx context.Context, suiteName string) (*BenchmarkRun, error) {
	bs.mu.RLock()
	suite, exists := bs.benchmarks[suiteName]
	bs.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("benchmark suite '%s' not found", suiteName)
	}

	run := &BenchmarkRun{
		ID:          fmt.Sprintf("bench_%d", time.Now().Unix()),
		Timestamp:   time.Now(),
		SuiteName:   suiteName,
		Results:     make(map[string]*BenchmarkResult),
		Environment: bs.getEnvironment(),
	}

	// Run each benchmark
	for name, benchmark := range suite.Benchmarks {
		select {
		case <-ctx.Done():
			return run, ctx.Err()
		default:
		}

		result, err := bs.runBenchmark(ctx, benchmark)
		if err != nil {
			result.Error = err
		}
		run.Results[name] = result
	}

	// Generate summary
	run.Summary = bs.generateSummary(run)

	// Store in history
	bs.mu.Lock()
	bs.history = append(bs.history, run)
	bs.mu.Unlock()

	return run, nil
}

// runBenchmark runs a single benchmark
func (bs *BenchmarkSystem) runBenchmark(ctx context.Context, benchmark *Benchmark) (*BenchmarkResult, error) {
	result := &BenchmarkResult{
		Name:      benchmark.Name,
		Timestamp: time.Now(),
	}

	// Setup
	if benchmark.Setup != nil {
		if err := benchmark.Setup(); err != nil {
			return result, fmt.Errorf("setup failed: %w", err)
		}
	}

	// Teardown
	if benchmark.Teardown != nil {
		defer benchmark.Teardown()
	}

	// Warmup runs
	for i := 0; i < bs.config.WarmupRuns; i++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		if err := benchmark.Function(); err != nil {
			return result, fmt.Errorf("warmup run %d failed: %w", i+1, err)
		}
	}

	// Benchmark runs
	var durations []time.Duration
	var memoryUsages []uint64
	var cpuUsages []float64

	for i := 0; i < bs.config.BenchmarkRuns; i++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		start := time.Now()
		startMem := bs.getMemoryUsage()
		startCPU := bs.getCPUUsage()

		if err := benchmark.Function(); err != nil {
			return result, fmt.Errorf("benchmark run %d failed: %w", i+1, err)
		}

		duration := time.Since(start)
		endMem := bs.getMemoryUsage()
		endCPU := bs.getCPUUsage()

		durations = append(durations, duration)
		memoryUsages = append(memoryUsages, endMem-startMem)
		cpuUsages = append(cpuUsages, endCPU-startCPU)
	}

	// Calculate statistics
	result.Duration = bs.calculateAverageDuration(durations)
	result.MemoryUsage = bs.calculateAverageMemory(memoryUsages)
	result.CPUUsage = bs.calculateAverageCPU(cpuUsages)
	result.Iterations = len(durations)
	result.Operations = len(durations)
	result.OpsPerSecond = float64(result.Operations) / result.Duration.Seconds()
	result.MemoryPerOp = float64(result.MemoryUsage) / float64(result.Operations)
	result.Percentiles = bs.calculatePercentiles(durations)
	result.Statistics = bs.calculateStatistics(durations)

	return result, nil
}

// calculateAverageDuration calculates the average duration
func (bs *BenchmarkSystem) calculateAverageDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

// calculateAverageMemory calculates the average memory usage
func (bs *BenchmarkSystem) calculateAverageMemory(memoryUsages []uint64) uint64 {
	if len(memoryUsages) == 0 {
		return 0
	}

	var total uint64
	for _, m := range memoryUsages {
		total += m
	}
	return total / uint64(len(memoryUsages))
}

// calculateAverageCPU calculates the average CPU usage
func (bs *BenchmarkSystem) calculateAverageCPU(cpuUsages []float64) float64 {
	if len(cpuUsages) == 0 {
		return 0
	}

	var total float64
	for _, c := range cpuUsages {
		total += c
	}
	return total / float64(len(cpuUsages))
}

// calculatePercentiles calculates percentile data
func (bs *BenchmarkSystem) calculatePercentiles(durations []time.Duration) Percentiles {
	if len(durations) == 0 {
		return Percentiles{}
	}

	// Sort durations
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	return Percentiles{
		P50:  bs.percentile(sorted, 0.50),
		P90:  bs.percentile(sorted, 0.90),
		P95:  bs.percentile(sorted, 0.95),
		P99:  bs.percentile(sorted, 0.99),
		P999: bs.percentile(sorted, 0.999),
	}
}

// percentile calculates the nth percentile
func (bs *BenchmarkSystem) percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	index := int(p * float64(len(sorted)-1))
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

// calculateStatistics calculates statistical data
func (bs *BenchmarkSystem) calculateStatistics(durations []time.Duration) Statistics {
	if len(durations) == 0 {
		return Statistics{}
	}

	// Convert to float64 for calculations
	values := make([]float64, len(durations))
	for i, d := range durations {
		values[i] = float64(d.Nanoseconds())
	}

	// Calculate mean
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calculate standard deviation
	var variance float64
	for _, v := range values {
		variance += math.Pow(v-mean, 2)
	}
	variance /= float64(len(values))
	stdDev := math.Sqrt(variance)

	// Find min and max
	min := values[0]
	max := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Calculate median
	sort.Float64s(values)
	median := values[len(values)/2]
	if len(values)%2 == 0 {
		median = (values[len(values)/2-1] + values[len(values)/2]) / 2
	}

	return Statistics{
		Mean:   mean,
		Median: median,
		StdDev: stdDev,
		Min:    min,
		Max:    max,
	}
}

// generateSummary generates a summary of benchmark results
func (bs *BenchmarkSystem) generateSummary(run *BenchmarkRun) *BenchmarkSummary {
	summary := &BenchmarkSummary{
		TotalBenchmarks: len(run.Results),
	}

	var totalDuration time.Duration
	var totalOpsPerSec float64
	var bestPerformer string
	var worstPerformer string
	var bestOpsPerSec float64
	var worstOpsPerSec float64

	for name, result := range run.Results {
		if result.Error != nil {
			summary.Failed++
		} else {
			summary.Passed++
			totalDuration += result.Duration
			totalOpsPerSec += result.OpsPerSecond

			if bestPerformer == "" || result.OpsPerSecond > bestOpsPerSec {
				bestPerformer = name
				bestOpsPerSec = result.OpsPerSecond
			}
			if worstPerformer == "" || result.OpsPerSecond < worstOpsPerSec {
				worstPerformer = name
				worstOpsPerSec = result.OpsPerSecond
			}
		}
	}

	summary.TotalDuration = totalDuration
	if summary.Passed > 0 {
		summary.AvgOpsPerSec = totalOpsPerSec / float64(summary.Passed)
	}
	summary.BestPerformer = bestPerformer
	summary.WorstPerformer = worstPerformer

	return summary
}

// getEnvironment gets current environment information
func (bs *BenchmarkSystem) getEnvironment() map[string]interface{} {
	env := make(map[string]interface{})
	env["timestamp"] = time.Now().Format(time.RFC3339)
	env["go_version"] = "1.24" // This would be dynamically determined
	env["platform"] = "darwin"
	env["arch"] = "amd64"
	return env
}

// getMemoryUsage gets current memory usage (placeholder)
func (bs *BenchmarkSystem) getMemoryUsage() uint64 {
	// This would use runtime.ReadMemStats in a real implementation
	return 0
}

// getCPUUsage gets current CPU usage (placeholder)
func (bs *BenchmarkSystem) getCPUUsage() float64 {
	// This would use runtime.ReadCPUStats in a real implementation
	return 0
}

// ExportResults exports benchmark results
func (bs *BenchmarkSystem) ExportResults(run *BenchmarkRun, format string) error {
	exportDir := filepath.Join(bs.baseDir, bs.config.ExportPath)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("benchmark_%s_%s.%s", run.SuiteName, run.ID, format)
	filepath := filepath.Join(exportDir, filename)

	switch format {
	case "json":
		return bs.exportJSON(run, filepath)
	case "csv":
		return bs.exportCSV(run, filepath)
	case "html":
		return bs.exportHTML(run, filepath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// exportJSON exports results as JSON
func (bs *BenchmarkSystem) exportJSON(run *BenchmarkRun, filepath string) error {
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

// exportCSV exports results as CSV
func (bs *BenchmarkSystem) exportCSV(run *BenchmarkRun, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	file.WriteString("Benchmark,Duration (ns),Memory Usage (bytes),CPU Usage,Ops/sec,Memory/Op,P50,P90,P95,P99,P999,Error\n")

	// Write data
	for name, result := range run.Results {
		errorStr := ""
		if result.Error != nil {
			errorStr = result.Error.Error()
		}

		line := fmt.Sprintf("%s,%d,%d,%.2f,%.2f,%.2f,%d,%d,%d,%d,%d,%s\n",
			name,
			result.Duration.Nanoseconds(),
			result.MemoryUsage,
			result.CPUUsage,
			result.OpsPerSecond,
			result.MemoryPerOp,
			result.Percentiles.P50.Nanoseconds(),
			result.Percentiles.P90.Nanoseconds(),
			result.Percentiles.P95.Nanoseconds(),
			result.Percentiles.P99.Nanoseconds(),
			result.Percentiles.P999.Nanoseconds(),
			errorStr,
		)
		file.WriteString(line)
	}

	return nil
}

// exportHTML exports results as HTML
func (bs *BenchmarkSystem) exportHTML(run *BenchmarkRun, filepath string) error {
	html := bs.generateHTMLReport(run)
	return os.WriteFile(filepath, []byte(html), 0644)
}

// generateHTMLReport generates HTML report content
func (bs *BenchmarkSystem) generateHTMLReport(run *BenchmarkRun) string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Benchmark Report - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .summary { background: #e7f3ff; padding: 15px; margin: 20px 0; border-radius: 5px; }
        .benchmark { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .metric { display: inline-block; margin: 10px; padding: 10px; background: #f8f9fa; border-radius: 3px; }
        .error { background: #f8d7da; color: #721c24; padding: 10px; border-radius: 3px; }
        table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Benchmark Report</h1>
        <p><strong>Suite:</strong> %s</p>
        <p><strong>Generated:</strong> %s</p>
    </div>
    
    <div class="summary">
        <h2>Summary</h2>
        <div class="metric"><strong>Total Benchmarks:</strong> %d</div>
        <div class="metric"><strong>Passed:</strong> %d</div>
        <div class="metric"><strong>Failed:</strong> %d</div>
        <div class="metric"><strong>Total Duration:</strong> %v</div>
        <div class="metric"><strong>Average Ops/sec:</strong> %.2f</div>
        <div class="metric"><strong>Best Performer:</strong> %s</div>
        <div class="metric"><strong>Worst Performer:</strong> %s</div>
    </div>`,
		run.SuiteName, run.SuiteName, run.Timestamp.Format(time.RFC3339),
		run.Summary.TotalBenchmarks, run.Summary.Passed, run.Summary.Failed,
		run.Summary.TotalDuration, run.Summary.AvgOpsPerSec,
		run.Summary.BestPerformer, run.Summary.WorstPerformer)

	// Add benchmark details
	for name, result := range run.Results {
		html += fmt.Sprintf(`
    <div class="benchmark">
        <h3>%s</h3>`, name)

		if result.Error != nil {
			html += fmt.Sprintf(`<div class="error"><strong>Error:</strong> %s</div>`, result.Error.Error())
		} else {
			html += fmt.Sprintf(`
        <div class="metric"><strong>Duration:</strong> %v</div>
        <div class="metric"><strong>Memory Usage:</strong> %d bytes</div>
        <div class="metric"><strong>CPU Usage:</strong> %.2f%%</div>
        <div class="metric"><strong>Ops/sec:</strong> %.2f</div>
        <div class="metric"><strong>Memory/Op:</strong> %.2f bytes</div>
        
        <h4>Percentiles</h4>
        <table>
            <tr><th>P50</th><th>P90</th><th>P95</th><th>P99</th><th>P999</th></tr>
            <tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>
        </table>
        
        <h4>Statistics</h4>
        <table>
            <tr><th>Mean</th><th>Median</th><th>Std Dev</th><th>Min</th><th>Max</th></tr>
            <tr><td>%.2f ns</td><td>%.2f ns</td><td>%.2f ns</td><td>%.2f ns</td><td>%.2f ns</td></tr>
        </table>`,
				result.Duration, result.MemoryUsage, result.CPUUsage,
				result.OpsPerSecond, result.MemoryPerOp,
				result.Percentiles.P50, result.Percentiles.P90, result.Percentiles.P95,
				result.Percentiles.P99, result.Percentiles.P999,
				result.Statistics.Mean, result.Statistics.Median, result.Statistics.StdDev,
				result.Statistics.Min, result.Statistics.Max)
		}

		html += `</div>`
	}

	html += `</body></html>`
	return html
}

// GetHistory returns benchmark history
func (bs *BenchmarkSystem) GetHistory() []*BenchmarkRun {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	history := make([]*BenchmarkRun, len(bs.history))
	copy(history, bs.history)
	return history
}

// AnalyzeTrends analyzes benchmark trends over time
func (bs *BenchmarkSystem) AnalyzeTrends() map[string]*BenchmarkTrend {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	trends := make(map[string]*BenchmarkTrend)

	// Group results by benchmark name
	benchmarkResults := make(map[string][]*BenchmarkResult)
	for _, run := range bs.history {
		for name, result := range run.Results {
			benchmarkResults[name] = append(benchmarkResults[name], result)
		}
	}

	// Analyze trends for each benchmark
	for name, results := range benchmarkResults {
		if len(results) < 2 {
			continue
		}

		trend := &BenchmarkTrend{
			BenchmarkName: name,
			DataPoints:    len(results),
			FirstRun:      results[0].Timestamp,
			LastRun:       results[len(results)-1].Timestamp,
		}

		// Calculate trend direction
		firstAvg := results[0].Duration.Nanoseconds()
		lastAvg := results[len(results)-1].Duration.Nanoseconds()

		if lastAvg < firstAvg {
			trend.Direction = "improving"
			trend.Improvement = float64(firstAvg-lastAvg) / float64(firstAvg) * 100
		} else if lastAvg > firstAvg {
			trend.Direction = "regressing"
			trend.Regression = float64(lastAvg-firstAvg) / float64(firstAvg) * 100
		} else {
			trend.Direction = "stable"
		}

		trends[name] = trend
	}

	return trends
}

// BenchmarkTrend represents a benchmark trend over time
type BenchmarkTrend struct {
	BenchmarkName string    `json:"benchmark_name"`
	Direction     string    `json:"direction"`
	Improvement   float64   `json:"improvement"`
	Regression    float64   `json:"regression"`
	DataPoints    int       `json:"data_points"`
	FirstRun      time.Time `json:"first_run"`
	LastRun       time.Time `json:"last_run"`
}
