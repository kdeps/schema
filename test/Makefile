# Test suite Makefile for PKL Schema
# Provides targets for running various types of tests and generating reports

# ANSI color codes
GREEN=\033[1;32m
RED=\033[1;31m
YELLOW=\033[1;33m
BLUE=\033[1;34m
MAGENTA=\033[1;35m
CYAN=\033[1;36m
RESET=\033[0m

.PHONY: help test test-go test-pkl test-benchmark test-coverage test-integration test-all clean report test-advanced test-benchmark-system test-filtering test-coverage-analysis test-cli-simulation test-dev deps fmt lint security ci stats

# Default target
help:
	@echo "$(MAGENTA)🧪 PKL Schema Test Suite$(RESET)"
	@echo "$(MAGENTA)========================$(RESET)"
	@echo ""
	@echo "$(CYAN)Available targets:$(RESET)"
	@echo "  $(BLUE)test$(RESET)                - Run all tests, format code, lint, and generate report"
	@echo "  $(BLUE)test-go$(RESET)         - Run Go integration tests"
	@echo "  $(BLUE)test-pkl$(RESET)        - Run PKL CLI schema validation"
	@echo "  $(BLUE)test-benchmark$(RESET)  - Run performance benchmarks"
	@echo "  $(BLUE)test-coverage$(RESET)   - Generate test coverage report"
	@echo "  $(BLUE)test-integration$(RESET)- Run comprehensive integration tests"
	@echo "  $(BLUE)test-advanced$(RESET)   - Run advanced features tests"
	@echo "  $(BLUE)test-all$(RESET)        - Run all tests and generate reports"
	@echo "  $(BLUE)clean$(RESET)           - Clean temporary files and reports"
	@echo "  $(BLUE)report$(RESET)          - Generate comprehensive test report"
	@echo ""
	@echo "$(CYAN)Advanced Features:$(RESET)"
	@echo "  $(YELLOW)test-benchmark-system$(RESET)    - Test benchmark system functionality"
	@echo "  $(YELLOW)test-filtering$(RESET)          - Test advanced test filtering"
	@echo "  $(YELLOW)test-coverage-analysis$(RESET)  - Test coverage analysis system"
	@echo "  $(YELLOW)test-cli-simulation$(RESET)     - Test CLI-like functionality"
	@echo ""

# Run Go integration tests
test-go:
	@echo "$(BLUE)🔧 Running Go integration tests...$(RESET)"
	@go test -v -timeout=60s ./...
	@echo "$(GREEN)✅ Go tests completed$(RESET)"

# Run PKL CLI schema validation
test-pkl:
	@echo "$(BLUE)📋 Running PKL CLI schema validation...$(RESET)"
	@echo "Validating PKL files..."
	@for file in *.pkl; do \
		if [ -f "$$file" ]; then \
			echo "  Validating $$file..."; \
			pkl eval --allowed-resources "agent:/*,pklres:/*,session:/*,tool:/*,memory:/*,item:/*,prop:/*" "$$file" > /dev/null 2>&1 || echo "    $(YELLOW)⚠️  $$file has syntax issues$(RESET)"; \
		fi; \
	done
	@echo "$(GREEN)✅ PKL validation completed$(RESET)"

# Run comprehensive PKL CLI tests
test-pkl-comprehensive:
	@echo "$(BLUE)🧪 Running comprehensive PKL CLI tests...$(RESET)"
	@if [ -f "comprehensive_function_tests.pkl" ]; then \
		echo "  Running comprehensive function tests..."; \
		pkl eval --allowed-resources "agent:/*,pklres:/*,session:/*,tool:/*,memory:/*,item:/*,prop:/*" comprehensive_function_tests.pkl > /dev/null 2>&1 && echo "    $(GREEN)✅ Passed$(RESET)" || echo "    $(RED)❌ Failed$(RESET)"; \
	fi
	@if [ -f "test_functions.pkl" ]; then \
		echo "  Running test functions..."; \
		pkl eval --allowed-resources "agent:/*,pklres:/*,session:/*,tool:/*,memory:/*,item:/*,prop:/*" test_functions.pkl > /dev/null 2>&1 && echo "    $(GREEN)✅ Passed$(RESET)" || echo "    $(RED)❌ Failed$(RESET)"; \
	fi
	@echo "$(GREEN)✅ Comprehensive PKL tests completed$(RESET)"

# Run performance benchmarks
test-benchmark:
	@echo "$(BLUE)⚡ Running performance benchmarks...$(RESET)"
	@go test -v -bench=. -benchmem -timeout=300s ./...
	@echo "$(GREEN)✅ Benchmarks completed$(RESET)"

# Generate test coverage report
test-coverage:
	@echo "$(BLUE)📊 Generating test coverage report...$(RESET)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go test -v -run TestCoverageAnalysis ./...
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(RESET)"

# Run comprehensive integration tests
test-integration:
	@echo "$(BLUE)🔗 Running comprehensive integration tests...$(RESET)"
	@go test -v -run TestIntegrationSuite -timeout=120s ./...
	@echo "$(GREEN)✅ Integration tests completed$(RESET)"

# Run advanced features tests
test-advanced: test-benchmark-system test-filtering test-coverage-analysis test-cli-simulation
	@echo "$(MAGENTA)🚀 Advanced features tests completed!$(RESET)"

# Test benchmark system functionality
test-benchmark-system:
	@echo "$(YELLOW)⚡ Testing benchmark system...$(RESET)"
	@go test -v -run TestBenchmarkSystem -timeout=60s ./...
	@echo "$(GREEN)✅ Benchmark system tests completed$(RESET)"

# Test advanced test filtering
test-filtering:
	@echo "$(YELLOW)🔍 Testing advanced test filtering...$(RESET)"
	@go test -v -run TestTestFiltering -timeout=30s ./...
	@echo "$(GREEN)✅ Test filtering tests completed$(RESET)"

# Test coverage analysis system
test-coverage-analysis:
	@echo "$(YELLOW)📊 Testing coverage analysis system...$(RESET)"
	@go test -v -run TestCoverageAnalysis -timeout=30s ./...
	@echo "$(GREEN)✅ Coverage analysis tests completed$(RESET)"

# Test CLI-like functionality
test-cli-simulation:
	@echo "$(YELLOW)💻 Testing CLI simulation...$(RESET)"
	@go test -v -run TestAdvancedFeaturesCLI -timeout=60s ./...
	@echo "$(GREEN)✅ CLI simulation tests completed$(RESET)"

# Run all tests and generate reports
test-all: test-go test-pkl test-pkl-comprehensive test-benchmark test-coverage test-integration test-parallel test-analytics test-fixtures test-diagnostics test-advanced
	@echo "$(MAGENTA)🎉 All tests completed!$(RESET)"
	@echo "$(CYAN)📋 Reports generated:$(RESET)"
	@echo "  - $(GREEN)coverage.html$(RESET) (test coverage)"
	@echo "  - $(GREEN)test_coverage_report.json$(RESET) (detailed coverage)"
	@echo "  - $(GREEN)Benchmark results in test output$(RESET)"
	@echo "  - $(GREEN)Analytics reports in reports/ directory$(RESET)"
	@echo "  - $(GREEN)Advanced features reports in coverage/ and benchmarks/ directories$(RESET)"

# Clean temporary files and reports
clean:
	@echo "$(BLUE)🧹 Cleaning temporary files and reports...$(RESET)"
	@rm -f coverage.out coverage.html test_coverage_report.json
	@rm -f *.db
	@rm -rf pkl-test-*
	@rm -rf pklres-*-*.db
	@rm -rf coverage/
	@rm -rf benchmarks/
	@rm -rf reports/
	@echo "$(GREEN)✅ Cleanup completed$(RESET)"

# Generate comprehensive test report
report: test-coverage
	@echo "$(MAGENTA)📋 Generating comprehensive test report...$(RESET)"
	@echo "🧪 PKL Schema Test Report" > TEST_REPORT.md
	@echo "=========================" >> TEST_REPORT.md
	@echo "" >> TEST_REPORT.md
	@echo "Generated: $$(date)" >> TEST_REPORT.md
	@echo "" >> TEST_REPORT.md
	@echo "## Test Results" >> TEST_REPORT.md
	@echo "" >> TEST_REPORT.md
	@if [ -f "test_coverage_report.json" ]; then \
		echo "### Coverage Summary" >> TEST_REPORT.md; \
		echo "\`\`\`json" >> TEST_REPORT.md; \
		cat test_coverage_report.json >> TEST_REPORT.md; \
		echo "\`\`\`" >> TEST_REPORT.md; \
	fi
	@echo "" >> TEST_REPORT.md
	@echo "## Test Files" >> TEST_REPORT.md
	@echo "" >> TEST_REPORT.md
	@for file in *_test.go; do \
		if [ -f "$$file" ]; then \
			echo "- \`$$file\`" >> TEST_REPORT.md; \
		fi; \
	done
	@echo "" >> TEST_REPORT.md
	@echo "## PKL Test Files" >> TEST_REPORT.md
	@echo "" >> TEST_REPORT.md
	@for file in *.pkl; do \
		if [ -f "$$file" ]; then \
			echo "- \`$$file\`" >> TEST_REPORT.md; \
		fi; \
	done
	@echo "$(GREEN)✅ Test report generated: TEST_REPORT.md$(RESET)"

# Quick test (fastest tests only)
test-quick:
	@echo "$(BLUE)⚡ Running quick tests...$(RESET)"
	@go test -v -run TestCoverageAnalysis -timeout=30s ./...
	@echo "$(GREEN)✅ Quick tests completed$(RESET)"

# Run parallel execution tests
test-parallel:
	@echo "$(BLUE)🔄 Running parallel execution tests...$(RESET)"
	@go test -v -run TestParallelExecution -timeout=60s ./...
	@echo "$(GREEN)✅ Parallel tests completed$(RESET)"

# Run analytics tests
test-analytics:
	@echo "$(BLUE)📈 Running analytics tests...$(RESET)"
	@go test -v -run TestAnalytics -timeout=30s ./...
	@echo "$(GREEN)✅ Analytics tests completed$(RESET)"

# Run fixture management tests
test-fixtures:
	@echo "$(BLUE)🎭 Running fixture management tests...$(RESET)"
	@go test -v -run TestFixtures -timeout=30s ./...
	@echo "$(GREEN)✅ Fixture tests completed$(RESET)"

# Run diagnostic tests
test-diagnostics:
	@echo "$(BLUE)🔍 Running diagnostic tests...$(RESET)"
	@go test -v -run TestDiagnostics -timeout=30s ./...
	@echo "$(GREEN)✅ Diagnostic tests completed$(RESET)"

# Development mode (watch for changes)
test-dev:
	@echo "$(CYAN)👀 Development mode - watching for changes...$(RESET)"
	@echo "Press Ctrl+C to stop"
	@while true; do \
		make test-quick; \
		sleep 5; \
	done

# Install dependencies
deps:
	@echo "$(BLUE)📦 Installing dependencies...$(RESET)"
	@go mod tidy
	@go mod download
	@echo "$(GREEN)✅ Dependencies installed$(RESET)"

# Format code
fmt:
	@echo "$(BLUE)🎨 Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(RESET)"

# Lint code
lint:
	@echo "$(BLUE)🔍 Linting code...$(RESET)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found, skipping linting$(RESET)"; \
	fi
	@echo "$(GREEN)✅ Linting completed$(RESET)"

# Security scan
security:
	@echo "$(BLUE)🔒 Running security scan...$(RESET)"
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)⚠️  gosec not found, skipping security scan$(RESET)"; \
	fi
	@echo "$(GREEN)✅ Security scan completed$(RESET)"

# Full CI pipeline
ci: deps fmt lint security test-all report
	@echo "$(MAGENTA)🚀 CI pipeline completed successfully!$(RESET)"

# Show test statistics
stats:
	@echo "📈 Test Statistics"
	@echo "=================="
	@echo "Go test files: $$(ls -1 *_test.go 2>/dev/null | wc -l | tr -d ' ')"
	@echo "PKL test files: $$(ls -1 *.pkl 2>/dev/null | wc -l | tr -d ' ')"
	@echo "Total test files: $$(ls -1 *_test.go *.pkl 2>/dev/null | wc -l | tr -d ' ')"
	@if [ -f "test_coverage_report.json" ]; then \
		echo ""; \
		echo "Coverage report available: test_coverage_report.json"; \
	fi 