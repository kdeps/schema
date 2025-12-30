#!/usr/bin/env bash

# Script to download external PKL dependencies for offline use
set -e

# Require bash 4.0+ for associative arrays
if ((BASH_VERSINFO[0] < 4)); then
    echo "Warning: Bash 4.0+ required for advanced conflict resolution"
    echo "Using simple conflict detection mode"
    SIMPLE_MODE=true
else
    SIMPLE_MODE=false
fi

DEPS_DIR="assets/pkl/external"
VERSIONS_FILE="versions.json"

# Check if versions.json exists
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "Error: $VERSIONS_FILE not found"
    exit 1
fi

# Generic case-insensitive conflict resolution
# This will automatically detect and rename conflicting files
# Format: "relative_path:new_name"
if [ "$SIMPLE_MODE" = false ]; then
    declare -A RENAME_MAP
fi

# Function to detect case-insensitive conflicts
detect_conflicts() {
    local dir="$1"
    local seen_file="/tmp/seen_files.txt"

    rm -f "$seen_file"
    touch "$seen_file"

    find "$dir" -name "*.pkl" -type f | sort | while read -r file; do
        filename=$(basename "$file")
        filename_lower=$(echo "$filename" | tr '[:upper:]' '[:lower:]')
        rel_path="${file#$dir/}"

        # Check if lowercase version already seen
        existing=$(grep "^${filename_lower}:" "$seen_file" 2>/dev/null | head -1 | cut -d: -f2-)

        if [ -n "$existing" ]; then
            # Conflict detected!
            echo "⚠️  Case conflict detected:"
            echo "   Existing: $existing"
            echo "   Conflict: $rel_path"

            # Generate new name by prefixing with parent directory name
            parent_dir=$(dirname "$rel_path")
            parent_name=$(basename "$parent_dir")
            new_name="${parent_name}_${filename}"

            # Store rename mapping
            echo "$rel_path:$new_name" >> /tmp/rename_map.txt
        else
            # Record this file (lowercase:fullpath)
            echo "${filename_lower}:${rel_path}" >> "$seen_file"
        fi
    done

    rm -f "$seen_file"
}

# Function to apply renames
apply_renames() {
    local base_dir="$1"

    if [ ! -f /tmp/rename_map.txt ]; then
        return 0
    fi

    while IFS=: read -r old_path new_name; do
        old_file="$base_dir/$old_path"
        new_file="$(dirname "$old_file")/$new_name"

        if [ -f "$old_file" ]; then
            echo "   Renaming: $(basename "$old_path") → $new_name"
            mv "$old_file" "$new_file"

            # Store mapping for reference fixing (to a file for bash 3 compat)
            echo "$old_path|$(dirname "$old_path")/$new_name" >> /tmp/rename_mappings.txt
        fi
    done < /tmp/rename_map.txt

    rm -f /tmp/rename_map.txt
}

# Function to fix references in PKL files
fix_references() {
    local base_dir="$1"

    if [ ! -f /tmp/rename_mappings.txt ]; then
        echo "No files were renamed, skipping reference updates"
        return 0
    fi

    echo ""
    echo "Fixing references to renamed files..."

    # Build sed script once for all replacements
    local sed_script="/tmp/sed_replacements.txt"
    rm -f "$sed_script"

    # Build list of old filenames for grep pattern
    local grep_pattern="/tmp/grep_pattern.txt"
    rm -f "$grep_pattern"

    while IFS='|' read -r old_path new_path; do
        old_name=$(basename "$old_path")
        new_name=$(basename "$new_path")

        # Add to grep pattern (for filtering files that need updates)
        echo "$old_name" >> "$grep_pattern"

        # Add sed replacement rules
        echo "s|import \".*/${old_name}\"|import \"${new_path}\"|g" >> "$sed_script"
        echo "s|\".*/${old_name}\"|\"${new_path}\"|g" >> "$sed_script"
    done < /tmp/rename_mappings.txt

    # Build combined grep pattern (match any old filename)
    local pattern=$(paste -sd '|' "$grep_pattern")

    # Find files that actually contain references to renamed files
    # This pre-filters to avoid processing files that don't need updates
    local files_to_update="/tmp/files_to_update.txt"
    find "$base_dir" -name "*.pkl" -type f -print0 | \
        xargs -0 grep -l -E "$pattern" 2>/dev/null > "$files_to_update" || true

    if [ ! -s "$files_to_update" ]; then
        echo "   No files need reference updates"
        rm -f /tmp/rename_mappings.txt "$sed_script" "$grep_pattern" "$files_to_update"
        echo "✓ Reference fixing complete"
        return 0
    fi

    # Apply all replacements to filtered files in one pass
    local count=0
    while read -r pkl_file; do
        # Apply sed script (cross-platform)
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '' -f "$sed_script" "$pkl_file"
        else
            # Linux
            sed -i -f "$sed_script" "$pkl_file"
        fi
        echo "   Updated references in: $(basename "$pkl_file")"
        ((count++))
    done < "$files_to_update"

    # Cleanup
    rm -f /tmp/rename_mappings.txt "$sed_script" "$grep_pattern" "$files_to_update"

    echo "✓ Reference fixing complete ($count files updated)"
}

# Read versions from JSON file  
PKL_GO_VERSION=$(jq -r '.dependencies."pkl-go".version' "$VERSIONS_FILE")

echo "Downloading PKL dependencies for offline use..."
echo "pkl-go version: $PKL_GO_VERSION"

# Clean up any existing external dependencies
rm -rf "$DEPS_DIR"
mkdir -p "$DEPS_DIR"

# Download pkl-go entire repository
echo "Downloading pkl-go complete repository..."
PKL_GO_URL="https://github.com/apple/pkl-go/archive/v${PKL_GO_VERSION}.tar.gz"
curl -L "$PKL_GO_URL" | tar -xz -C /tmp/

# Copy only PKL files from pkl-go repository
echo "Copying pkl-go PKL files..."
mkdir -p "$DEPS_DIR/pkl-go"
find "/tmp/pkl-go-${PKL_GO_VERSION}" -name "*.pkl" -type f -exec sh -c 'rel_path="${1#/tmp/pkl-go-'"${PKL_GO_VERSION}"'/}" && mkdir -p "$2/$(dirname "$rel_path")" && cp "$1" "$2/$rel_path"' _ {} "$DEPS_DIR/pkl-go" \;

# Download all pkl-pantry packages
echo "Downloading pkl-pantry packages..."
mkdir -p "$DEPS_DIR/pkl-pantry/packages"

# Get all pkl-pantry packages from versions.json
PKL_PANTRY_PACKAGES=$(jq -r '.dependencies."pkl-pantry".packages | keys[]' "$VERSIONS_FILE")

for package in $PKL_PANTRY_PACKAGES; do
    echo "Processing package: $package"
    
    # Get version for this package
    VERSION=$(jq -r ".dependencies.\"pkl-pantry\".packages.\"$package\".version" "$VERSIONS_FILE")
    
    # Create package directory
    PACKAGE_DIR="$DEPS_DIR/pkl-pantry/packages/$package"
    mkdir -p "$PACKAGE_DIR"
    
    # Download package
    PKL_PANTRY_TAG="${package}@${VERSION}"
    PKL_PANTRY_URL="https://github.com/apple/pkl-pantry/archive/${PKL_PANTRY_TAG}.tar.gz"
    
    echo "  Downloading $package@$VERSION..."
    curl -L "$PKL_PANTRY_URL" | tar -xz -C /tmp/
    
    # Copy only PKL files (GitHub replaces @ with - in directory name)
    PKL_PANTRY_DIR_NAME="pkl-pantry-$(echo "${PKL_PANTRY_TAG}" | tr '@' '-')"

    # Copy all PKL files
    find "/tmp/${PKL_PANTRY_DIR_NAME}" -name "*.pkl" -type f -exec cp {} "$PACKAGE_DIR/" \;

    # Cleanup temporary files for this package
    rm -rf "/tmp/${PKL_PANTRY_DIR_NAME}"
done

# Cleanup temporary files
rm -rf "/tmp/pkl-go-${PKL_GO_VERSION}"

echo ""
echo "Dependencies downloaded successfully!"
echo "pkl-go repository in: $DEPS_DIR/pkl-go/"
echo "pkl-pantry repository in: $DEPS_DIR/pkl-pantry/"

# Step 1: Detect case-insensitive conflicts
echo ""
echo "Checking for case-insensitive filename conflicts..."
rm -f /tmp/rename_map.txt
detect_conflicts "$DEPS_DIR"

# Step 2: Apply renames if conflicts were found
if [ -f /tmp/rename_map.txt ]; then
    echo ""
    echo "Applying renames to resolve conflicts..."
    apply_renames "$DEPS_DIR"

    # Step 3: Fix all references to renamed files
    fix_references "$DEPS_DIR"
else
    echo "✓ No case-insensitive conflicts detected"
fi

# List downloaded pkl files
echo ""
echo "Downloaded PKL files:"
find "$DEPS_DIR" -name "*.pkl" -type f | head -10
echo "... and more"

# Show directory structure
echo ""
echo "Directory structure:"
ls -la "$DEPS_DIR/"