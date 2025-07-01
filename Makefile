# Directory containing the .pkl files
PKL_DIR := deps/pkl
# Directory where the generated output will be stored
OUTPUT_DIR := .
# Command to process .pkl files
GEN_COMMAND := pkl-gen-go
# Get the current directory
CURRENT_DIR := $(shell pwd)
# Collect all .pkl files in PKL_DIR
PKL_FILES := $(wildcard $(PKL_DIR)/*.pkl)

# Generate output files in OUTPUT_DIR
generate:
		@pkl project resolve --root-dir $(CURRENT_DIR) --module-path $(PKL_DIR) $(PKL_DIR)

		@if [ -d "$(OUTPUT_DIR)/gen" ]; then \
			rm -rf $(OUTPUT_DIR)/gen; \
		fi

		@for pkl in $(PKL_FILES); do \
			$(GEN_COMMAND) $$pkl --output-path $(OUTPUT_DIR) | sed 's;/github.com/kdeps/schema/gen;;g'; \
		done

		@if [ -d "$(OUTPUT_DIR)/github.com/kdeps/schema/gen" ]; then \
			mv $(OUTPUT_DIR)/github.com/kdeps/schema/gen $(OUTPUT_DIR); \
			mv $(OUTPUT_DIR)/github.com/kdeps/schema/deps/pkl $(OUTPUT_DIR); \
			rm -rf $(OUTPUT_DIR)/github.com; \
		fi

# Run the comprehensive PKL test suite
test:
		@echo "Running Comprehensive PKL Function Test Suite..."
		@pkl eval test/test_functions.pkl

# Run individual module tests
test-utils:
		@echo "Running Utils.pkl Unit Tests..."
		@pkl eval test/test_utils.pkl

# Run all tests (comprehensive + individual)
test-all: test test-utils
		@echo "All PKL tests completed successfully!"

# Run tests and generate Go code
test-and-generate: test-all generate

# Help target
help:
		@echo "Available targets:"
		@echo "  generate         - Generate Go code from PKL files"
		@echo "  test            - Run comprehensive PKL function test suite"
		@echo "  test-utils       - Run Utils.pkl unit tests"
		@echo "  test-all         - Run all test suites"
		@echo "  test-and-generate - Run all tests then generate Go code"
		@echo "  help            - Show this help message"

.PHONY: generate test test-utils test-all test-and-generate help
