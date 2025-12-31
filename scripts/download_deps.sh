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
PKL_GO_FILES=$(jq -r '.dependencies."pkl-go".files[]' "$VERSIONS_FILE")

echo "Downloading PKL dependencies for offline use..."
echo "pkl-go version: $PKL_GO_VERSION"

# Clean up any existing external dependencies
rm -rf "$DEPS_DIR"
mkdir -p "$DEPS_DIR"

# Download pkl-go repository
echo "Downloading pkl-go repository..."
PKL_GO_URL="https://github.com/apple/pkl-go/archive/v${PKL_GO_VERSION}.tar.gz"
curl -L "$PKL_GO_URL" | tar -xz -C /tmp/

# Copy only specified PKL files from pkl-go repository
echo "Copying specified pkl-go PKL files..."
mkdir -p "$DEPS_DIR/pkl-go"

for file in $PKL_GO_FILES; do
    # Find the file in the downloaded repository
    SOURCE_FILE=$(find "/tmp/pkl-go-${PKL_GO_VERSION}" -type f -name "$(basename "$file")" -path "*/$file" | head -1)

    if [ -z "$SOURCE_FILE" ]; then
        echo "  Warning: File $file not found in pkl-go repository"
        continue
    fi

    # Determine relative path from pkl-go root
    rel_path="${SOURCE_FILE#/tmp/pkl-go-${PKL_GO_VERSION}/}"

    # Create target directory structure
    TARGET_FILE="$DEPS_DIR/pkl-go/$rel_path"
    mkdir -p "$(dirname "$TARGET_FILE")"

    cp "$SOURCE_FILE" "$TARGET_FILE"
    echo "  Copied: $rel_path"
done

# Download all pkl-pantry packages
echo "Downloading pkl-pantry packages..."
mkdir -p "$DEPS_DIR/pkl-pantry/packages"

# Get all pkl-pantry packages from versions.json
PKL_PANTRY_PACKAGES=$(jq -r '.dependencies."pkl-pantry".packages | keys[]' "$VERSIONS_FILE")

for package in $PKL_PANTRY_PACKAGES; do
    echo "Processing package: $package"

    # Get version for this package
    VERSION=$(jq -r ".dependencies.\"pkl-pantry\".packages.\"$package\".version" "$VERSIONS_FILE")

    # Get files list for this package
    FILES=$(jq -r ".dependencies.\"pkl-pantry\".packages.\"$package\".files[]" "$VERSIONS_FILE")

    # Create package directory
    PACKAGE_DIR="$DEPS_DIR/pkl-pantry/packages/$package"
    mkdir -p "$PACKAGE_DIR"

    # Download package
    PKL_PANTRY_TAG="${package}@${VERSION}"
    PKL_PANTRY_URL="https://github.com/apple/pkl-pantry/archive/${PKL_PANTRY_TAG}.tar.gz"

    echo "  Downloading $package@$VERSION..."
    curl -L "$PKL_PANTRY_URL" | tar -xz -C /tmp/

    # Copy only specified PKL files (GitHub replaces @ with - in directory name)
    PKL_PANTRY_DIR_NAME="pkl-pantry-$(echo "${PKL_PANTRY_TAG}" | tr '@' '-')"

    # Copy each specified file, preserving directory structure for subdirectories
    for file in $FILES; do
        # Find the file in the downloaded package
        SOURCE_FILE=$(find "/tmp/${PKL_PANTRY_DIR_NAME}" -type f -name "$(basename "$file")" -path "*/$file" | head -1)

        if [ -z "$SOURCE_FILE" ]; then
            echo "  Warning: File $file not found in package"
            continue
        fi

        # Determine target path
        if [[ "$file" == */* ]]; then
            # File is in a subdirectory, preserve structure
            TARGET_FILE="$PACKAGE_DIR/$file"
            mkdir -p "$(dirname "$TARGET_FILE")"
        else
            # File is in root of package
            TARGET_FILE="$PACKAGE_DIR/$file"
        fi

        cp "$SOURCE_FILE" "$TARGET_FILE"
        echo "    Copied: $file"
    done

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