#!/usr/bin/env bash

# Script to download external PKL dependencies for offline use
set -e

DEPS_DIR="assets/pkl/external"
VERSIONS_FILE="versions.json"

# Check if versions.json exists
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "Error: $VERSIONS_FILE not found"
    exit 1
fi

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
curl -sL "$PKL_GO_URL" | tar -xz -C /tmp/

# Copy only PKL files from pkl-go repository
echo "Copying pkl-go PKL files..."
mkdir -p "$DEPS_DIR/pkl-go"
find "/tmp/pkl-go-${PKL_GO_VERSION}" -name "*.pkl" -type f -exec sh -c 'rel_path="${1#/tmp/pkl-go-'"${PKL_GO_VERSION}"'/}" && mkdir -p "$2/$(dirname "$rel_path")" && cp "$1" "$2/$rel_path"' _ {} "$DEPS_DIR/pkl-go" \;

# Download all pkl-pantry packages
# Note: pkl-pantry uses a monorepo structure where all packages share the same release
echo ""
echo "Downloading pkl-pantry packages..."
mkdir -p "$DEPS_DIR/pkl-pantry/packages"

# Get a version tag to download (all packages share releases in the monorepo)
FIRST_PACKAGE=$(jq -r '.dependencies."pkl-pantry".packages | keys[0]' "$VERSIONS_FILE")
PKL_PANTRY_VERSION=$(jq -r ".dependencies.\"pkl-pantry\".packages.\"$FIRST_PACKAGE\".version" "$VERSIONS_FILE")
PKL_PANTRY_TAG="${FIRST_PACKAGE}@${PKL_PANTRY_VERSION}"
PKL_PANTRY_URL="https://github.com/apple/pkl-pantry/archive/${PKL_PANTRY_TAG}.tar.gz"
PKL_PANTRY_DIR_NAME="pkl-pantry-$(echo "${PKL_PANTRY_TAG}" | tr '@' '-')"

echo "Downloading pkl-pantry monorepo (${PKL_PANTRY_TAG})..."
curl -sL "$PKL_PANTRY_URL" | tar -xz -C /tmp/

# Copy files for each package from the monorepo
PKL_PANTRY_PACKAGES=$(jq -r '.dependencies."pkl-pantry".packages | keys[]' "$VERSIONS_FILE")

for package in $PKL_PANTRY_PACKAGES; do
    echo ""
    echo "Processing package: $package"

    # Create package directory
    PACKAGE_DIR="$DEPS_DIR/pkl-pantry/packages/$package"
    mkdir -p "$PACKAGE_DIR"

    # Copy only the files specified in versions.json
    jq -r ".dependencies.\"pkl-pantry\".packages.\"$package\".files[]" "$VERSIONS_FILE" | while IFS= read -r file; do
        SOURCE_PATH="/tmp/${PKL_PANTRY_DIR_NAME}/packages/${package}/${file}"

        # Handle files in subdirectories (e.g., "internal/Type.pkl")
        if [[ "$file" == *"/"* ]]; then
            mkdir -p "$PACKAGE_DIR/$(dirname "$file")"
            DEST_PATH="$PACKAGE_DIR/$file"
        else
            DEST_PATH="$PACKAGE_DIR/$file"
        fi

        if [ -f "$SOURCE_PATH" ]; then
            cp "$SOURCE_PATH" "$DEST_PATH"
            echo "  ✓ $file"
        else
            echo "  ⚠ Warning: Could not find $file"
        fi
    done
done

# Cleanup temporary files
rm -rf "/tmp/${PKL_PANTRY_DIR_NAME}"
rm -rf "/tmp/pkl-go-${PKL_GO_VERSION}"

echo ""
echo "✓ Dependencies downloaded successfully!"
echo ""
echo "Summary:"
echo "  pkl-go:     $(find "$DEPS_DIR/pkl-go" -name "*.pkl" -type f | wc -l | tr -d ' ') files"
echo "  pkl-pantry: $(find "$DEPS_DIR/pkl-pantry/packages" -name "*.pkl" -type f | wc -l | tr -d ' ') files"
