#!/bin/bash

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
    find "/tmp/${PKL_PANTRY_DIR_NAME}" -name "*.pkl" -type f -exec cp {} "$PACKAGE_DIR/" \;
    
    # Cleanup temporary files for this package
    rm -rf "/tmp/${PKL_PANTRY_DIR_NAME}"
done

# Cleanup temporary files
rm -rf "/tmp/pkl-go-${PKL_GO_VERSION}"

echo "Dependencies downloaded successfully!"
echo "pkl-go repository in: $DEPS_DIR/pkl-go/"
echo "pkl-pantry repository in: $DEPS_DIR/pkl-pantry/"

# List downloaded pkl files
echo ""
echo "Downloaded PKL files:"
find "$DEPS_DIR" -name "*.pkl" -type f | head -10
echo "... and more"

# Show directory structure
echo ""
echo "Directory structure:"
ls -la "$DEPS_DIR/"