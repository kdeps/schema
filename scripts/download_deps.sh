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

# Function to fix references to a renamed file
fix_references_for_file() {
    local base_dir="$1"
    local old_name="$2"
    local new_path="$3"

    # Find files that reference the old filename
    local files_with_refs
    files_with_refs=$(find "$base_dir" -name "*.pkl" -type f -print0 | \
        xargs -0 grep -l "$old_name" 2>/dev/null) || true

    if [ -z "$files_with_refs" ]; then
        return 0
    fi

    # Update references in each file
    echo "$files_with_refs" | while IFS= read -r pkl_file; do
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '' "s|import \".*/${old_name}\"|import \"${new_path}\"|g" "$pkl_file"
            sed -i '' "s|\".*/${old_name}\"|\"${new_path}\"|g" "$pkl_file"
        else
            # Linux
            sed -i "s|import \".*/${old_name}\"|import \"${new_path}\"|g" "$pkl_file"
            sed -i "s|\".*/${old_name}\"|\"${new_path}\"|g" "$pkl_file"
        fi
        echo "   Updated references in: $(basename "$pkl_file")"
    done
}

# Function to detect and resolve case-insensitive conflicts
detect_and_resolve_conflicts() {
    local dir="$1"
    local seen_file="/tmp/seen_files.txt"
    local conflicts_found=false

    rm -f "$seen_file"
    touch "$seen_file"

    find "$dir" -name "*.pkl" -type f | sort | while read -r file; do
        filename=$(basename "$file")
        filename_lower=$(echo "$filename" | tr '[:upper:]' '[:lower:]')
        rel_path="${file#$dir/}"

        # Check if lowercase version already seen
        existing=$(grep "^${filename_lower}:" "$seen_file" 2>/dev/null | head -1 | cut -d: -f2-)

        if [ -n "$existing" ]; then
            # Conflict detected! Resolve immediately
            echo "⚠️  Case conflict detected:"
            echo "   Existing: $existing"
            echo "   Conflict: $rel_path"

            # Generate new name by prefixing with parent directory name
            parent_dir=$(dirname "$rel_path")
            parent_name=$(basename "$parent_dir")
            new_name="${parent_name}_${filename}"

            old_file="$file"
            new_file="$(dirname "$file")/$new_name"
            new_rel_path="$(dirname "$rel_path")/$new_name"

            # Rename the file immediately
            echo "   Renaming: $filename → $new_name"
            mv "$old_file" "$new_file"

            # Fix references immediately
            fix_references_for_file "$dir" "$filename" "$new_rel_path"

            # Record the new file (lowercase:fullpath) so it's not detected as conflict
            echo "${filename_lower}:${new_rel_path}" >> "$seen_file"
            conflicts_found=true
        else
            # Record this file (lowercase:fullpath)
            echo "${filename_lower}:${rel_path}" >> "$seen_file"
        fi
    done

    rm -f "$seen_file"

    if [ "$conflicts_found" = false ]; then
        return 1
    fi
    return 0
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

# Detect and resolve case-insensitive conflicts
echo ""
echo "Checking for case-insensitive filename conflicts..."
if detect_and_resolve_conflicts "$DEPS_DIR"; then
    echo ""
    echo "✓ All conflicts resolved"
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