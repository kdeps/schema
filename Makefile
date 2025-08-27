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

# Update versions and dependencies
update-deps:
	@echo "🚀 Updating versions and dependencies..."
	@./scripts/update_versions.sh
	@./scripts/update_pkl_version.sh
	@./scripts/download_deps.sh
	@./scripts/update_imports.sh
	@echo "✅ Dependencies updated successfully!"

# Setup offline dependencies (run once or when adding new dependencies)
setup-offline:
	@echo "🛠️  Setting up offline dependencies..."
	@./scripts/download_deps.sh
	@mkdir -p assets/pkl && cp deps/pkl/*.pkl assets/pkl/
	@mkdir -p assets/pkl/external
	@echo "🔗 Updating imports to use local paths..."
	@./scripts/update_imports.sh
	@echo "✅ Offline dependencies setup complete!"

# Generate output files in OUTPUT_DIR
generate: setup-offline
	@echo "📦 Starting PKL code generation..."
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

	@echo "🔗 Updating imports to use local paths..."
	@./scripts/update_imports.sh

	@echo "🔨 Testing Go build..."
	@if go build ./assets >/dev/null 2>&1; then \
		echo "✅ Go build test passed!"; \
	else \
		echo "❌ Go build test failed!"; \
		exit 1; \
	fi

	@echo "🎉 PKL code generation completed successfully!"

# Full update and generate (recommended for CI/CD)
generate-latest: update-deps generate

# Clean generated files and cached dependencies
clean:
	@echo "🧹 Cleaning generated files..."
	@rm -rf gen/
	@rm -rf assets/pkl/
	@echo "✅ Clean completed!"

# Clean everything including downloaded dependencies
clean-all:
	@echo "🧹 Cleaning all generated files and dependencies..."
	@rm -rf gen/
	@rm -rf assets/pkl/
	@rm -rf assets/pkl/external/
	@echo "✅ Full clean completed!"

# Show help
help:
	@echo "Available commands:"
	@echo "  make generate        - Generate PKL code with offline dependencies (default)"
	@echo "  make generate-latest - Update to latest versions and generate"
	@echo "  make setup-offline   - Setup offline dependencies only"
	@echo "  make update-deps     - Update versions and dependencies"
	@echo "  make clean           - Clean generated files"
	@echo "  make clean-all       - Clean all files including dependencies"
	@echo "  make help            - Show this help"

.PHONY: generate generate-latest setup-offline update-deps clean clean-all help
