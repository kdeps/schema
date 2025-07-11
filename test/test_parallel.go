package test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestDependency represents a test dependency
type TestDependency struct {
	Name         string        `json:"name"`
	Dependencies []string      `json:"dependencies"`
	Resources    []string      `json:"resources"`
	Timeout      time.Duration `json:"timeout"`
	Priority     int           `json:"priority"`
}

// TestExecutionPlan represents a plan for parallel test execution
type TestExecutionPlan struct {
	Tests        map[string]*TestDependency `json:"tests"`
	Dependencies map[string][]string        `json:"dependencies"`
	Resources    map[string]int             `json:"resources"`
	MaxParallel  int                        `json:"max_parallel"`
}

// ParallelTestExecutor manages parallel test execution
type ParallelTestExecutor struct {
	plan         *TestExecutionPlan
	executor     *TestSuite
	resourcePool *ResourcePool
	results      map[string]*TestResult
	mu           sync.RWMutex
}

// NewParallelTestExecutor creates a new parallel test executor
func NewParallelTestExecutor(plan *TestExecutionPlan, executor *TestSuite) *ParallelTestExecutor {
	return &ParallelTestExecutor{
		plan:         plan,
		executor:     executor,
		resourcePool: NewResourcePool(plan.Resources),
		results:      make(map[string]*TestResult),
	}
}

// Execute executes tests in parallel according to the plan
func (pte *ParallelTestExecutor) Execute(ctx context.Context) error {
	// Build dependency graph
	graph := pte.buildDependencyGraph()

	// Execute tests in dependency order
	return pte.executeInOrder(ctx, graph)
}

// buildDependencyGraph builds a dependency graph for test execution
func (pte *ParallelTestExecutor) buildDependencyGraph() map[string][]string {
	graph := make(map[string][]string)

	// Initialize graph
	for testName := range pte.plan.Tests {
		graph[testName] = []string{}
	}

	// Add dependencies
	for testName, test := range pte.plan.Tests {
		for _, dep := range test.Dependencies {
			if _, exists := graph[dep]; exists {
				graph[dep] = append(graph[dep], testName)
			}
		}
	}

	return graph
}

// executeInOrder executes tests in dependency order with parallel execution
func (pte *ParallelTestExecutor) executeInOrder(ctx context.Context, graph map[string][]string) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, pte.plan.MaxParallel)
	errors := make(chan error, len(pte.plan.Tests))
	graphMu := sync.Mutex{} // Protect graph updates

	// Track completed tests
	completed := make(map[string]bool)
	completedMu := sync.Mutex{}

	for len(graph) > 0 {
		readyTests := pte.findReadyTests(graph)
		if len(readyTests) == 0 {
			// No ready tests, check if we're stuck
			if len(completed) == 0 {
				return fmt.Errorf("no tests can be executed - possible circular dependency")
			}
			break
		}

		// Execute ready tests in parallel
		for _, testName := range readyTests {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case semaphore <- struct{}{}:
				wg.Add(1)
				go func(name string) {
					defer wg.Done()
					defer func() { <-semaphore }()

					if err := pte.executeTest(ctx, name); err != nil {
						select {
						case errors <- fmt.Errorf("test %s failed: %w", name, err):
						default:
						}
					}

					// Mark as completed
					completedMu.Lock()
					completed[name] = true
					completedMu.Unlock()

					// Update graph with mutex protection
					graphMu.Lock()
					pte.updateGraph(graph, []string{name})
					graphMu.Unlock()
				}(testName)
			}
		}

		// Wait for current batch to complete
		wg.Wait()

		// Check for errors
		select {
		case err := <-errors:
			return err
		default:
		}
	}

	return nil
}

// findReadyTests finds tests that are ready to execute (no dependencies)
func (pte *ParallelTestExecutor) findReadyTests(graph map[string][]string) []string {
	var ready []string

	for testName, dependencies := range graph {
		if len(dependencies) == 0 {
			ready = append(ready, testName)
		}
	}

	return ready
}

// updateGraph updates the dependency graph after test execution
func (pte *ParallelTestExecutor) updateGraph(graph map[string][]string, completed []string) {
	// Remove completed tests from all dependency lists
	for _, completedTest := range completed {
		for testName, deps := range graph {
			newDeps := make([]string, 0)
			for _, dep := range deps {
				if dep != completedTest {
					newDeps = append(newDeps, dep)
				}
			}
			graph[testName] = newDeps
		}
		
		// Remove the completed test from the graph itself
		delete(graph, completedTest)
	}
}

// executeTest executes a single test with resource management
func (pte *ParallelTestExecutor) executeTest(ctx context.Context, testName string) error {
	test := pte.plan.Tests[testName]

	// Acquire resources
	resources := pte.resourcePool.Acquire(test.Resources)
	defer pte.resourcePool.Release(resources)

	// Create test context with timeout
	testCtx, cancel := context.WithTimeout(ctx, test.Timeout)
	defer cancel()

	// Execute test
	start := time.Now()
	err := pte.executor.RunTestWithContext(testCtx, testName, func(t *testing.T) error {
		// This would be replaced with actual test execution
		return nil
	})
	duration := time.Since(start)

	// Record result
	pte.mu.Lock()
	pte.results[testName] = &TestResult{
		Name:     testName,
		Status:   pte.getStatus(err),
		Duration: duration,
		Error:    err,
	}
	pte.mu.Unlock()

	return err
}

// getStatus determines test status from error
func (pte *ParallelTestExecutor) getStatus(err error) string {
	if err == nil {
		return "PASS"
	}
	return "FAIL"
}

// GetResults returns all test results
func (pte *ParallelTestExecutor) GetResults() map[string]*TestResult {
	pte.mu.RLock()
	defer pte.mu.RUnlock()

	results := make(map[string]*TestResult)
	for k, v := range pte.results {
		results[k] = v
	}
	return results
}

// ResourcePool manages resource allocation for parallel tests
type ResourcePool struct {
	resources map[string]int
	available map[string]int
	mu        sync.Mutex
}

// NewResourcePool creates a new resource pool
func NewResourcePool(resources map[string]int) *ResourcePool {
	available := make(map[string]int)
	for k, v := range resources {
		available[k] = v
	}

	return &ResourcePool{
		resources: resources,
		available: available,
	}
}

// Acquire acquires resources for test execution
func (rp *ResourcePool) Acquire(required []string) []string {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	var acquired []string

	for _, resource := range required {
		if rp.available[resource] > 0 {
			rp.available[resource]--
			acquired = append(acquired, resource)
		}
	}

	return acquired
}

// Release releases resources back to the pool
func (rp *ResourcePool) Release(resources []string) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	for _, resource := range resources {
		if rp.available[resource] < rp.resources[resource] {
			rp.available[resource]++
		}
	}
}

// TestScheduler manages test scheduling and execution
type TestScheduler struct {
	executor *ParallelTestExecutor
	plan     *TestExecutionPlan
}

// NewTestScheduler creates a new test scheduler
func NewTestScheduler(executor *TestSuite) *TestScheduler {
	return &TestScheduler{
		plan: &TestExecutionPlan{
			Tests:        make(map[string]*TestDependency),
			Dependencies: make(map[string][]string),
			Resources:    make(map[string]int),
			MaxParallel:  4,
		},
	}
}

// AddTest adds a test to the execution plan
func (ts *TestScheduler) AddTest(name string, dependencies []string, resources []string, timeout time.Duration, priority int) {
	ts.plan.Tests[name] = &TestDependency{
		Name:         name,
		Dependencies: dependencies,
		Resources:    resources,
		Timeout:      timeout,
		Priority:     priority,
	}

	// Update resource requirements
	for _, resource := range resources {
		if ts.plan.Resources[resource] == 0 {
			ts.plan.Resources[resource] = 1
		}
	}
}

// SetMaxParallel sets the maximum number of parallel tests
func (ts *TestScheduler) SetMaxParallel(max int) {
	ts.plan.MaxParallel = max
}

// SetResourceLimit sets the limit for a specific resource
func (ts *TestScheduler) SetResourceLimit(resource string, limit int) {
	ts.plan.Resources[resource] = limit
}

// Schedule schedules and executes tests
func (ts *TestScheduler) Schedule(ctx context.Context, executor *TestSuite) error {
	ts.executor = NewParallelTestExecutor(ts.plan, executor)
	return ts.executor.Execute(ctx)
}

// GetResults returns test execution results
func (ts *TestScheduler) GetResults() map[string]*TestResult {
	if ts.executor == nil {
		return make(map[string]*TestResult)
	}
	return ts.executor.GetResults()
}

// TestExecutionOptimizer optimizes test execution order
type TestExecutionOptimizer struct {
	plan *TestExecutionPlan
}

// NewTestExecutionOptimizer creates a new test execution optimizer
func NewTestExecutionOptimizer(plan *TestExecutionPlan) *TestExecutionOptimizer {
	return &TestExecutionOptimizer{
		plan: plan,
	}
}

// Optimize optimizes the test execution plan
func (teo *TestExecutionOptimizer) Optimize() *TestExecutionPlan {
	optimized := &TestExecutionPlan{
		Tests:        make(map[string]*TestDependency),
		Dependencies: make(map[string][]string),
		Resources:    make(map[string]int),
		MaxParallel:  teo.plan.MaxParallel,
	}

	// Copy tests and sort by priority
	var tests []*TestDependency
	for _, test := range teo.plan.Tests {
		tests = append(tests, test)
	}

	// Sort by priority (higher priority first)
	for i := 0; i < len(tests)-1; i++ {
		for j := i + 1; j < len(tests); j++ {
			if tests[i].Priority < tests[j].Priority {
				tests[i], tests[j] = tests[j], tests[i]
			}
		}
	}

	// Rebuild plan with optimized order
	for _, test := range tests {
		optimized.Tests[test.Name] = test
	}

	// Copy resources
	for k, v := range teo.plan.Resources {
		optimized.Resources[k] = v
	}

	return optimized
}

// TestExecutionMonitor monitors test execution progress
type TestExecutionMonitor struct {
	plan     *TestExecutionPlan
	results  map[string]*TestResult
	progress chan TestProgress
	mu       sync.RWMutex
}

// TestProgress represents test execution progress
type TestProgress struct {
	TestName   string        `json:"test_name"`
	Status     string        `json:"status"`
	Progress   float64       `json:"progress"`
	Duration   time.Duration `json:"duration"`
	TotalTests int           `json:"total_tests"`
	Completed  int           `json:"completed"`
	Passed     int           `json:"passed"`
	Failed     int           `json:"failed"`
}

// NewTestExecutionMonitor creates a new test execution monitor
func NewTestExecutionMonitor(plan *TestExecutionPlan) *TestExecutionMonitor {
	return &TestExecutionMonitor{
		plan:     plan,
		results:  make(map[string]*TestResult),
		progress: make(chan TestProgress, 100),
	}
}

// UpdateProgress updates test execution progress
func (tem *TestExecutionMonitor) UpdateProgress(testName, status string, duration time.Duration) {
	tem.mu.Lock()
	defer tem.mu.Unlock()

	// Update results
	tem.results[testName] = &TestResult{
		Name:     testName,
		Status:   status,
		Duration: duration,
	}

	// Calculate progress
	total := len(tem.plan.Tests)
	completed := len(tem.results)
	passed := 0
	failed := 0

	for _, result := range tem.results {
		if result.Status == "PASS" {
			passed++
		} else if result.Status == "FAIL" {
			failed++
		}
	}

	progress := float64(completed) / float64(total) * 100

	// Send progress update
	select {
	case tem.progress <- TestProgress{
		TestName:   testName,
		Status:     status,
		Progress:   progress,
		Duration:   duration,
		TotalTests: total,
		Completed:  completed,
		Passed:     passed,
		Failed:     failed,
	}:
	default:
		// Channel full, skip update
	}
}

// GetProgressChannel returns the progress channel
func (tem *TestExecutionMonitor) GetProgressChannel() <-chan TestProgress {
	return tem.progress
}

// GetResults returns current test results
func (tem *TestExecutionMonitor) GetResults() map[string]*TestResult {
	tem.mu.RLock()
	defer tem.mu.RUnlock()

	results := make(map[string]*TestResult)
	for k, v := range tem.results {
		results[k] = v
	}
	return results
}

// RunTestWithContext runs a test with context support
func (ts *TestSuite) RunTestWithContext(ctx context.Context, testName string, testFunc func(*testing.T) error) error {
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

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = testFunc(&testing.T{})
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

	// Update metrics
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

	return err
}

// TestParallelExecution demonstrates parallel test execution
func TestParallelExecution(t *testing.T) {
	// Create test suite
	suite := NewTestSuite()

	// Create scheduler
	scheduler := NewTestScheduler(suite)

	// Add tests with dependencies
	scheduler.AddTest("setup", []string{}, []string{"database"}, 30*time.Second, 1)
	scheduler.AddTest("test1", []string{"setup"}, []string{"api"}, 10*time.Second, 2)
	scheduler.AddTest("test2", []string{"setup"}, []string{"api"}, 10*time.Second, 2)
	scheduler.AddTest("test3", []string{"test1", "test2"}, []string{"file"}, 5*time.Second, 3)
	scheduler.AddTest("cleanup", []string{"test3"}, []string{"database"}, 5*time.Second, 4)

	// Set resource limits
	scheduler.SetResourceLimit("database", 1)
	scheduler.SetResourceLimit("api", 2)
	scheduler.SetResourceLimit("file", 1)

	// Set parallel limit
	scheduler.SetMaxParallel(3)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Execute tests
	err := scheduler.Schedule(ctx, suite)
	if err != nil {
		t.Errorf("Parallel execution failed: %v", err)
	}

	// Get results
	results := scheduler.GetResults()
	if len(results) != 5 {
		t.Errorf("Expected 5 test results, got %d", len(results))
	}

	// Print results
	for testName, result := range results {
		t.Logf("Test %s: %s (%.2fs)", testName, result.Status, result.Duration.Seconds())
	}
}
