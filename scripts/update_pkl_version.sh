#!/bin/bash

# Script to update minPklVersion in all PKL files based on versions.json
set -e

PKL_DIR="deps/pkl"
VERSIONS_FILE="versions.json"

# Check if versions.json exists
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "Error: $VERSIONS_FILE not found"
    exit 1
fi

# Get PKL version from versions.json
PKL_VERSION=$(jq -r '.pkl.version' "$VERSIONS_FILE")

if [ "$PKL_VERSION" = "null" ] || [ -z "$PKL_VERSION" ]; then
    echo "Error: No PKL version found in $VERSIONS_FILE"
    exit 1
fi

echo "Updating PKL version references to $PKL_VERSION..."

# Update minPklVersion in all PKL files
find "$PKL_DIR" -name "*.pkl" -type f ! -path "*/external/*" -exec sed -i '' \
    "s/@ModuleInfo { minPklVersion = \"[^\"]*\" }/@ModuleInfo { minPklVersion = \"$PKL_VERSION\" }/g" {} \;

echo "PKL version references updated!"

# Show what was changed
echo ""
echo "Files with updated PKL version:"
grep -r "minPklVersion = \"$PKL_VERSION\"" "$PKL_DIR" --include="*.pkl" | cut -d: -f1 | sort | uniq