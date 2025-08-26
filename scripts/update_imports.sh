#!/bin/bash

# Script to update PKL imports to use local dependencies based on versions.json
set -e

PKL_DIR="assets/pkl"
VERSIONS_FILE="versions.json"

# Check if versions.json exists
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "Error: $VERSIONS_FILE not found"
    exit 1
fi

echo "Updating PKL imports in $PKL_DIR to use local dependencies based on $VERSIONS_FILE..."

# Process each dependency from versions.json
jq -r '.dependencies | to_entries[] | @base64' "$VERSIONS_FILE" | while IFS= read -r row; do
    dependency=$(echo "$row" | base64 --decode)
    
    dep_name=$(echo "$dependency" | jq -r '.key')
    dep_version=$(echo "$dependency" | jq -r '.value.version')
    dep_url=$(echo "$dependency" | jq -r '.value.url')
    
    echo "Processing dependency: $dep_name v$dep_version"
    
    case "$dep_name" in
        "pkl-go")
            # Process pkl-go files (simple structure)
            echo "$dependency" | jq -r '.value.files[]' | while IFS= read -r file; do
                echo "  Updating imports for file: $file"
                
                local_path="../external/pkl-go/codegen/src/$file"
                remote_pattern="package://pkg\.pkl-lang\.org/pkl-go/pkl\.golang@[^\"]*#/$file"
                
                # Update imports in all PKL files
                find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" -exec sed -i '' \
                  "s|import \"$remote_pattern\"|import \"$local_path\"|g" {} \;
            done
            ;;
        "pkl-pantry")
            # Process pkl-pantry packages (nested structure)
            echo "$dependency" | jq -r '.value.packages | to_entries[] | @base64' | while IFS= read -r pkg_row; do
                package=$(echo "$pkg_row" | base64 --decode)
                
                pkg_name=$(echo "$package" | jq -r '.key')
                pkg_version=$(echo "$package" | jq -r '.value.version')
                
                echo "  Processing package: $pkg_name v$pkg_version"
                
                # Process each file in the package
                echo "$package" | jq -r '.value.files[]' | while IFS= read -r file; do
                    echo "    Updating imports for file: $file"
                    
                    # Convert package name to path (e.g., pkl.experimental.uri -> pkl/experimental/uri)
                    pkg_path=$(echo "$pkg_name" | tr '.' '/')
                    local_path="../external/pkl-pantry/packages/$pkg_name/$file"
                    remote_pattern="package://pkg\.pkl-lang\.org/pkl-pantry/$pkg_name@[^\"]*#/$file"
                    
                    # Update imports in all PKL files
                    find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" -exec sed -i '' \
                      "s|import \"$remote_pattern\"|import \"$local_path\"|g" {} \;
                done
            done
            ;;
        *)
            echo "  Warning: Unknown dependency $dep_name, skipping..."
            continue
            ;;
    esac
done

echo "Import updates completed!"

# Show what was changed
echo ""
echo "Files with updated imports:"
find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" -exec grep -l "external/" {} \; | sort