# 🧪 Enhanced PKL Schema Test Suite

A comprehensive, production-ready test suite for the kdeps PKL schema with advanced features including parallel execution, real-time monitoring, analytics, and CI/CD integration.

## 🚀 Key Features

### Core Testing Infrastructure
- **Thread-safe test execution** with parallel and distributed testing
- **Comprehensive test categorization** (unit, integration, performance, security)
- **Advanced retry logic** for flaky test detection and recovery
- **Real-time performance monitoring** with resource tracking
- **Intelligent test scheduling** with dependency management

### Analytics & Reporting
- **Multi-format analytics export** (JSON, Markdown, HTML, CSV)
- **Historical trend analysis** with regression detection
- **Real-time monitoring dashboard** with metrics visualization
- **Performance benchmarking** with detailed profiling
- **Coverage analysis** with quality gates

### CI/CD Integration
- **GitHub Actions workflows** for automated testing and deployment
- **Multi-platform testing** (Linux, macOS, multiple Go versions)
- **Security scanning** with gosec integration
- **Quality gates** with coverage and performance thresholds
- **Automated deployment** with staging and production environments

## 📋 Quick Start

### Prerequisites
```bash
# Install Go 1.24+
go version

# Install PKL 0.28.2+
pkl --version

# Install additional tools
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
```

### Basic Usage
```bash
# Run all tests
cd schema/test
make test

# Run specific test categories
make test-unit
make test-integration
make test-performance
make test-parallel

# Generate analytics reports
make analytics-export
make analytics-view

# Run CI/CD pipeline locally
make ci-test
make ci-coverage
make ci-benchmark
make ci-security
```

## 🧪 Test Categories

### 1. Unit Tests
```bash
# Run unit tests with coverage
go test -v -coverprofile=coverage.out ./...

# Run with race detection
go test -v -race ./...

# Run specific test
go test -v -run TestComprehensiveSuite ./...
```

### 2. Integration Tests
```bash
# Run integration tests
go test -v -tags=integration -timeout=10m ./...

# Run PKL schema tests
pkl eval comprehensive_function_tests.pkl
pkl eval null_safety_tests.pkl
pkl eval state_management_tests.pkl
```

### 3. Performance Tests
```bash
# Run benchmarks
go test -v -bench=. -benchmem ./...

# Run specific benchmarks
go test -v -bench=BenchmarkPKL -benchmem ./...

# Generate performance report
make benchmark-report
```

### 4. Parallel Tests
```bash
# Run parallel tests
go test -v -parallel=8 ./...

# Run distributed tests
make test-distributed
```

## 📊 Analytics & Monitoring

### Real-time Monitoring
```bash
# Start monitoring dashboard
make monitoring-start

# View metrics
make monitoring-metrics

# Stop monitoring
make monitoring-stop
```

### Analytics Export
```bash
# Export analytics in multiple formats
make analytics-export

# View historical trends
make analytics-view

# Generate custom reports
make report-custom
```

### Performance Analysis
```bash
# Run performance benchmarks
make benchmark

# Analyze performance trends
make benchmark-analyze

# Generate performance report
make benchmark-report
```

## 🔧 Advanced Configuration

### Test Configuration
```go
// Custom test configuration
config := &TestConfig{
    Verbose:       true,
    RetryCount:    3,
    RetryDelay:    time.Second,
    Timeout:       30 * time.Second,
    Parallel:      true,
    FilterPattern: "Test.*",
}
```

### Monitoring Configuration
```go
// Custom monitoring configuration
monitor := NewMonitoringSystem(&MonitoringConfig{
    Enabled:           true,
    MonitorInterval:   5 * time.Second,
    AlertThresholds:   &AlertThresholds{
        CPUUsagePercent:    90,
        MemoryUsagePercent: 80,
        TestDurationMs:     5000,
        FailureRatePercent: 10,
    },
    MetricsRetention:  24 * time.Hour,
    EnableProfiling:   true,
})
```

### Analytics Configuration
```go
// Custom analytics configuration
analytics := NewTestAnalytics("reports", &AnalyticsConfig{
    HistorySize:         100,
    TrendWindow:         30 * 24 * time.Hour,
    RegressionThreshold: 0.1,
    ReportFormats:       []string{"json", "markdown", "html"},
    ExportPath:          "reports",
})
```

## 🚀 CI/CD Integration

### GitHub Actions Workflows

#### Test Suite (`test.yml`)
- **Unit Tests**: Multi-platform testing with race detection
- **Integration Tests**: PKL schema validation and resource testing
- **Performance Tests**: Benchmark execution and analysis
- **Parallel Tests**: Concurrent test execution validation
- **Security Scan**: gosec security vulnerability scanning
- **Quality Gates**: Coverage and performance thresholds

#### Deployment (`deploy.yml`)
- **Build**: Multi-platform Docker image building
- **Staging**: Automated staging deployment with smoke tests
- **Production**: Production deployment with health checks
- **Rollback**: Automatic rollback on deployment failure

### Local CI Simulation
```bash
# Run full CI pipeline locally
make ci-full

# Run individual CI stages
make ci-test
make ci-coverage
make ci-benchmark
make ci-security
make ci-quality
```

## 📈 Performance Optimization

### Caching Strategy
- **Test result caching** for faster re-runs
- **Resource pooling** for efficient parallel execution
- **Intelligent test scheduling** based on dependencies
- **Performance trend analysis** for optimization insights

### Resource Management
- **Memory usage monitoring** with automatic cleanup
- **CPU utilization tracking** with alerting
- **Network call optimization** with connection pooling
- **File operation batching** for improved I/O performance

## 🔍 Troubleshooting

### Common Issues

#### Concurrent Map Writes
```bash
# Error: fatal error: concurrent map writes
# Solution: All maps are now thread-safe with sync.Map and mutexes
```

#### Test Timeouts
```bash
# Increase timeout for slow tests
go test -v -timeout=5m ./...

# Or configure in TestConfig
config.Timeout = 5 * time.Minute
```

#### Memory Issues
```bash
# Monitor memory usage
make monitoring-metrics

# Check for memory leaks
go test -v -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### Debug Mode
```bash
# Enable debug mode for detailed logging
DEBUG=true go test -v ./...

# Run with diagnostic system
make test-diagnostics
```

## 📚 API Reference

### Core Types
```go
// TestSuite manages test execution
type TestSuite struct {
    metrics *TestMetrics
    logger  *TestLogger
    config  *TestConfig
}

// TestMetrics tracks execution metrics
type TestMetrics struct {
    TotalTests   int
    PassedTests  int
    FailedTests  int
    SkippedTests int
    StartTime    time.Time
    EndTime      time.Time
    TestResults  map[string]TestResult
    mu           sync.Mutex
}

// TestResult represents individual test results
type TestResult struct {
    Name     string
    Status   string // "PASS", "FAIL", "SKIP"
    Duration time.Duration
    Error    error
    Message  string
}
```

### Key Functions
```go
// Run test with retry logic
func (ts *TestSuite) RunTest(t *testing.T, testName string, testFunc func(*testing.T) error)

// Run test with context support
func (ts *TestSuite) RunTestWithContext(ctx context.Context, testName string, testFunc func(*testing.T) error) error

// Get test metrics
func (ts *TestSuite) GetMetrics() *TestMetrics

// Print test summary
func (ts *TestSuite) PrintSummary()
```

## 🤝 Contributing

### Adding New Tests
1. Create test file following naming convention: `test_*.go`
2. Use `TestSuite` for structured test execution
3. Add appropriate tags for categorization
4. Include performance benchmarks for new features
5. Update documentation and examples

### Test Guidelines
- **Thread Safety**: All shared data must be protected
- **Error Handling**: Comprehensive error scenarios
- **Performance**: Include benchmarks for performance-critical code
- **Documentation**: Clear test descriptions and examples
- **Maintainability**: Modular and reusable test components

### Code Quality
- **Coverage**: Maintain >80% test coverage
- **Performance**: No significant performance regressions
- **Security**: Pass all security scans
- **Documentation**: Keep documentation up-to-date

## 📄 License

This test suite is part of the kdeps project and follows the same license terms.

## 🆘 Support

For issues and questions:
1. Check the troubleshooting section
2. Review existing GitHub issues
3. Create a new issue with detailed information
4. Include test output and environment details

---

**Last Updated**: $(date)
**Version**: Enhanced Test Suite v2.0
**Compatibility**: Go 1.24+, PKL 0.28.2+ 