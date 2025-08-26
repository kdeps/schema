#!/bin/bash

# Script to download external PKL dependencies for offline use
set -e

DEPS_DIR="deps/pkl/external"
VERSIONS_FILE="versions.json"

# Check if versions.json exists
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "Error: $VERSIONS_FILE not found"
    exit 1
fi

# Read versions from JSON file  
PKL_GO_VERSION=$(jq -r '.dependencies."pkl-go".version' "$VERSIONS_FILE")
PKL_PANTRY_VERSION=$(jq -r '.dependencies."pkl-pantry".version' "$VERSIONS_FILE")

echo "Downloading PKL dependencies for offline use..."
echo "pkl-go version: $PKL_GO_VERSION"
echo "pkl-pantry version: $PKL_PANTRY_VERSION"

# Clean up any existing external dependencies
rm -rf "$DEPS_DIR"
mkdir -p "$DEPS_DIR"

# Download pkl-go entire repository
echo "Downloading pkl-go complete repository..."
PKL_GO_URL="https://github.com/apple/pkl-go/archive/v${PKL_GO_VERSION}.tar.gz"
curl -L "$PKL_GO_URL" | tar -xz -C /tmp/

# Copy entire pkl-go repository
echo "Copying pkl-go repository..."
cp -r "/tmp/pkl-go-${PKL_GO_VERSION}" "$DEPS_DIR/pkl-go"

# Download pkl-pantry experimental.uri package
echo "Downloading pkl-pantry experimental.uri package..."
PKL_PANTRY_TAG="pkl.experimental.uri@${PKL_PANTRY_VERSION}"
PKL_PANTRY_URL="https://github.com/apple/pkl-pantry/archive/${PKL_PANTRY_TAG}.tar.gz"
curl -L "$PKL_PANTRY_URL" | tar -xz -C /tmp/

# Copy pkl-pantry repository (GitHub replaces @ with - in directory name)
echo "Copying pkl-pantry repository..."
PKL_PANTRY_DIR_NAME="pkl-pantry-$(echo "${PKL_PANTRY_TAG}" | tr '@' '-')"
cp -r "/tmp/${PKL_PANTRY_DIR_NAME}" "$DEPS_DIR/pkl-pantry"

# Cleanup temporary files
rm -rf "/tmp/pkl-go-${PKL_GO_VERSION}" "/tmp/${PKL_PANTRY_DIR_NAME}"

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