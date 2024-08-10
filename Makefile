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
			rm -rf $(OUTPUT_DIR)/github.com; \
		fi
