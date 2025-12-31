#!/bin/bash

# Script to convert remote package URLs to local paths for offline compatibility
set -e

DEPS_DIR="deps/pkl"
VERSIONS_FILE="versions.json"

echo "Converting package URLs to local paths in $DEPS_DIR/*.pkl files for offline use..."

# Get pkl-go version
PKL_GO_VERSION=$(jq -r '.dependencies."pkl-go".version' "$VERSIONS_FILE")

# Function to replace in file (cross-platform compatible)
replace_in_file() {
    local file=$1
    local pattern=$2
    local replacement=$3

    sed "s|$pattern|$replacement|g" "$file" > "$file.tmp"
    mv "$file.tmp" "$file"
}

# Convert pkl-go package URLs to local paths (handles any version)
echo "Converting pkl-go imports (version $PKL_GO_VERSION)..."
find "$DEPS_DIR" -name "*.pkl" -type f | while read -r file; do
    replace_in_file "$file" "package://pkg.pkl-lang.org/pkl-go/pkl.golang@[^#]*#/go.pkl" "external/pkl-go/codegen/src/go.pkl"
done

# Convert pkl-pantry package URLs to local paths (handles any version)
echo "Converting pkl-pantry imports..."
find "$DEPS_DIR" -name "*.pkl" -type f | while read -r file; do
    replace_in_file "$file" "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@[^#]*#/URI.pkl" "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"
done

echo "Import conversion completed!"
echo ""

# Verify no package URLs remain
REMAINING_URLS=$(find "$DEPS_DIR" -name "*.pkl" -type f -exec grep -l "package://" {} \; 2>/dev/null | wc -l)

if [ "$REMAINING_URLS" -eq 0 ]; then
    echo "✅ All package URLs successfully converted to local paths"
else
    echo "⚠️  Warning: $REMAINING_URLS file(s) still contain package URLs:"
    find "$DEPS_DIR" -name "*.pkl" -type f -exec grep -l "package://" {} \;
    exit 1
fi
