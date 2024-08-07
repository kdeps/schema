# Directory containing the .pkl files
PKL_DIR := deps/pkl
# Directory where the generated output will be stored
OUTPUT_DIR := pkg
# Command to process .pkl files
GEN_COMMAND := pkl-gen-go
# Get the current directory
CURRENT_DIR := $(shell pwd)
# Collect all .pkl files in PKL_DIR
PKL_FILES := $(wildcard $(PKL_DIR)/*.pkl)

# Generate output files in OUTPUT_DIR
generate:
	@pkl project resolve --root-dir $(CURRENT_DIR) --module-path $(PKL_DIR) $(PKL_DIR)

	@rm -rf $(OUTPUT_DIR)
	@for pkl in $(PKL_FILES); do \
		$(GEN_COMMAND) $$pkl --output-path $(OUTPUT_DIR) | sed 's;/github.com/kdeps/schema/pkg;;g'; \
	done
	@mv $(OUTPUT_DIR)/github.com/kdeps/schema/pkg/* $(OUTPUT_DIR)
	@rm -rf $(OUTPUT_DIR)/github.com
