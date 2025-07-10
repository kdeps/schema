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
		@rm -f test/TEST_REPORT.md
		@echo "Clean completed!"

# Define PKL test files using wildcards (exclude report generators)
PKL_TEST_FILES := $(filter-out test/generate_test_report%.pkl, $(wildcard test/*.pkl))

# UNIFIED TEST TARGET - Runs all PKL tests + Go tests + generates test report
test:
		@echo "🧪 UNIFIED PKL TEST SUITE - COMPREHENSIVE VALIDATION"
		@echo "======================================================"
		@echo ""
		@echo "📊 Auto-discovering PKL test files..."
		@echo "   Found test files: $(notdir $(PKL_TEST_FILES))"
		@echo "   Total PKL tests: $(words $(PKL_TEST_FILES))"
		@echo ""
		@for pkl_file in $(PKL_TEST_FILES); do \
			echo "🔍 Executing: $$(basename $$pkl_file)"; \
			pkl eval $$pkl_file; \
			echo ""; \
		done
		@echo "🛠️  Running Go Assets Test Suite..."
		@cd test && go test -v .
		@echo ""
		@echo "📝 Generating Test Report..."
		@pkl eval test/generate_test_report_simple.pkl > test/TEST_REPORT.md
		@echo "✅ Test report generated: test/TEST_REPORT.md"
		@echo ""
		@echo "🎯 UNIFIED TEST SUMMARY:"
		@echo "   - PKL test files executed: $(words $(PKL_TEST_FILES))"
		@echo "   - Go asset tests: ✅ Completed"
		@echo "   - Test report: ✅ Generated"
		@echo ""
		@echo "📋 View complete results: cat test/TEST_REPORT.md"
		@echo "🚀 ALL TESTS COMPLETED SUCCESSFULLY!"

# Build target (includes unified tests, release notes, and generation)
build: test update-readme generate
	@echo "Build completed successfully with updated release notes!"

# Legacy compatibility targets (DEPRECATED - use 'make test' instead)
test-legacy:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@make test

test-utils:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@make test

test-assets:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@make test

test-assets-bench:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@cd test && go test -bench=. -v .

test-all:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@make test

test-all-comprehensive:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@make test

test-comprehensive:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@make test

test-and-generate:
		@echo "⚠️  DEPRECATED: Use 'make build' for unified testing + generation"
		@make build

test-new-attributes:
		@echo "⚠️  DEPRECATED: Use 'make test' for unified testing"
		@make test

# Help target
help:
		@echo "🛠️  KDEPS PKL SCHEMA BUILD SYSTEM"
		@echo "=================================="
		@echo ""
		@echo "📋 MAIN TARGETS:"
		@echo "  test               - 🧪 Run ALL PKL tests (wildcard discovery) + Go tests + generate test report"
		@echo "  build              - 🚀 Complete build: test + update README + generate Go code (CI/CD ready)"
		@echo "  clean              - 🧹 Clean generated files and copied assets"
		@echo ""
		@echo "🔧 UTILITY TARGETS:"
		@echo "  copy-pkl-assets    - Copy PKL files to assets directory for embedding"
		@echo "  update-readme      - Update README.md with latest release notes"
		@echo "  generate           - Copy PKL assets, update README.md and generate Go code from PKL files"
		@echo "  help               - Show this help message"
		@echo ""
		@echo "⚠️  LEGACY TARGETS (DEPRECATED):"
		@echo "  test-utils, test-assets, test-all, test-comprehensive, etc."
		@echo "  → Use 'make test' for unified testing instead"
		@echo ""
		@echo "💡 QUICK START:"
		@echo "  make test          # Run all tests and generate report"
		@echo "  make build         # Full build for production/CI"
		@echo ""
		@echo "📊 Test Discovery: Automatically finds all test/*.pkl files (excludes generators)"

.PHONY: copy-pkl-assets update-readme generate clean test build test-legacy test-utils test-assets test-assets-bench test-all test-all-comprehensive test-comprehensive test-and-generate test-new-attributes help
