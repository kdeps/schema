#!/bin/bash

# Script to update PKL imports to use local dependencies based on versions.json
set -e

PKL_DIR="assets/pkl"
VERSIONS_FILE="versions.json"
EXTERNAL_DIR="$PKL_DIR/external"

# Check if versions.json exists
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "Error: $VERSIONS_FILE not found"
    exit 1
fi

# Check if external directory exists
if [ ! -d "$EXTERNAL_DIR" ]; then
    echo "Error: External directory $EXTERNAL_DIR not found"
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
                
                # Check if the external file exists
                external_file="$EXTERNAL_DIR/pkl-go/codegen/src/$file"
                if [ ! -f "$external_file" ]; then
                    echo "    Warning: External file $external_file not found, skipping..."
                    continue
                fi
                
                local_path="external/pkl-go/codegen/src/$file"
                remote_pattern_1="package://pkg\.pkl-lang\.org/pkl-go/pkl\.golang@[^\"]*#/$file"
                remote_pattern_2="package://pkg\.pkl-lang\.org/pkl-go/pkl\.golang@[^\"]*#/codegen/src/$file"
                
                # Count files that would be updated
                files_updated=0
                
                # Update imports in all PKL files (excluding external directory)
                find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" | while IFS= read -r pkl_file; do
                    if grep -q "$remote_pattern_1\|$remote_pattern_2" "$pkl_file" 2>/dev/null; then
                        echo "    Updating: $pkl_file"
                        sed -i '' \
                          -e "s|import \"$remote_pattern_1\"|import \"$local_path\"|g" \
                          -e "s|import \"$remote_pattern_2\"|import \"$local_path\"|g" \
                          -e "s|amends \"$remote_pattern_1\"|amends \"$local_path\"|g" \
                          -e "s|amends \"$remote_pattern_2\"|amends \"$local_path\"|g" \
                          "$pkl_file"
                        files_updated=$((files_updated + 1))
                    fi
                done
                
                if [ $files_updated -eq 0 ]; then
                    echo "    No files needed updating for $file"
                fi
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
                    
                    # Check if the external file exists
                    external_file="$EXTERNAL_DIR/pkl-pantry/packages/$pkg_name/$file"
                    if [ ! -f "$external_file" ]; then
                        echo "      Warning: External file $external_file not found, skipping..."
                        continue
                    fi
                    
                    local_path="external/pkl-pantry/packages/$pkg_name/$file"
                    remote_pattern="package://pkg\.pkl-lang\.org/pkl-pantry/$pkg_name@[^\"]*#/$file"
                    
                    # Count files that would be updated
                    files_updated=0
                    
                    # Update imports in all PKL files (excluding external directory)
                    find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" | while IFS= read -r pkl_file; do
                        if grep -q "$remote_pattern" "$pkl_file" 2>/dev/null; then
                            echo "      Updating: $pkl_file"
                            sed -i '' \
                              -e "s|import \"$remote_pattern\"|import \"$local_path\"|g" \
                              -e "s|amends \"$remote_pattern\"|amends \"$local_path\"|g" \
                              "$pkl_file"
                            files_updated=$((files_updated + 1))
                        fi
                    done
                    
                    if [ $files_updated -eq 0 ]; then
                        echo "      No files needed updating for $file"
                    fi
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
echo "Summary of updates:"
echo "Files with local imports (external/ path):"
find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" -exec grep -l "external/" {} \; | sort | wc -l | xargs echo "  Total files with local imports:"

echo ""
echo "Files still using remote imports:"
remote_files=$(find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" -exec grep -l "package://pkg\.pkl-lang\.org" {} \; 2>/dev/null | sort)
if [ -n "$remote_files" ]; then
    echo "$remote_files"
    echo ""
    echo "Run this script again if remote imports were found and need conversion."
else
    echo "  None found - all imports are using local paths!"
fi