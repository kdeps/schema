# Directory containing the .pkl files
PKL_DIR := deps/pkl
# Directory where the generated output will be stored
OUTPUT_DIR := .
# Assets directory for embedding PKL files
ASSETS_PKL_DIR := assets/pkl
# Command to process .pkl files
GEN_COMMAND := pkl-gen-go
# Get the current directory
CURRENT_DIR := $(shell pwd)
# Collect all .pkl files in PKL_DIR
PKL_FILES := $(wildcard $(PKL_DIR)/*.pkl)

# Copy PKL files to assets directory for embedding
copy-pkl-assets:
		@echo "Copying PKL files to assets directory for embedding..."
		@mkdir -p $(ASSETS_PKL_DIR)
		@cp $(PKL_DIR)/*.pkl $(ASSETS_PKL_DIR)/
		@echo "PKL files copied to $(ASSETS_PKL_DIR)"

# Update README.md with latest release notes
update-readme:
		@echo "Updating README.md with latest release notes..."
		@chmod +x scripts/generate_release_notes.sh
		@scripts/generate_release_notes.sh > README.md
		@echo "README.md updated successfully!"

# Generate output files in OUTPUT_DIR (now includes README update and PKL asset copying)
generate: update-readme copy-pkl-assets
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

# Clean generated files and copied assets
clean:
		@echo "Cleaning generated files and assets..."
		@rm -rf $(OUTPUT_DIR)/gen
		@rm -rf $(ASSETS_PKL_DIR)
		@echo "Clean completed!"

# Run the comprehensive PKL test suite
test:
		@echo "Running Comprehensive PKL Function Test Suite..."
		@pkl eval test/test_functions.pkl

# Run individual module tests
test-utils:
		@echo "Running Utils.pkl Unit Tests..."
		@pkl eval test/test_utils.pkl

# Run Go assets tests
test-assets:
		@echo "Running Go Assets Test Suite..."
		@cd test && go test -v .

# Run Go assets tests with benchmarks
test-assets-bench:
		@echo "Running Go Assets Benchmarks..."
		@cd test && go test -bench=. -v .

# Run all PKL tests (comprehensive + individual)
test-all: test test-utils
		@echo "All PKL tests completed successfully!"

# Run all tests including Go assets tests
test-all-comprehensive: test-all test-assets
		@echo "All PKL and Go assets tests completed successfully!"

# Build target (includes tests, release notes, and generation)
build: test-all-comprehensive update-readme generate
	@echo "Build completed successfully with updated release notes!"

# Run tests and generate Go code (now includes README update and PKL asset copying)
test-and-generate: test-all-comprehensive generate

# Run newly added attributes tests
test-new-attributes:
		@echo "Running new attributes tests..."
		@pkl eval test/test_new_attributes.pkl
		@echo "New attributes tests completed!"

# Help target
help:
		@echo "Available targets:"
		@echo "  copy-pkl-assets    - Copy PKL files to assets directory for embedding"
		@echo "  update-readme      - Update README.md with latest release notes"
		@echo "  generate           - Copy PKL assets, update README.md and generate Go code from PKL files"
		@echo "  clean             - Clean generated files and copied assets"
		@echo "  test              - Run comprehensive PKL function test suite"
		@echo "  test-utils         - Run Utils.pkl unit tests"
		@echo "  test-assets        - Run Go assets test suite"
		@echo "  test-assets-bench  - Run Go assets tests with benchmarks"
		@echo "  test-all           - Run all PKL test suites"
		@echo "  test-all-comprehensive - Run all PKL and Go assets tests"
		@echo "  build             - Run tests, update README.md, and generate Go code (recommended for CI/CD)"
		@echo "  test-and-generate  - Run all tests, copy PKL assets, update README.md, then generate Go code"
		@echo "  test-new-attributes - Run tests for newly added attributes"
		@echo "  help              - Show this help message"

.PHONY: copy-pkl-assets update-readme generate clean test test-utils test-assets test-assets-bench test-all test-all-comprehensive build test-and-generate test-new-attributes help
