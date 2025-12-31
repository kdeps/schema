#!/bin/bash

# Script to fix imports in deps/pkl/*.pkl files to use package URIs
set -e

DEPS_DIR="deps/pkl"
VERSIONS_FILE="versions.json"

echo "Fixing imports in $DEPS_DIR/*.pkl files to use package URIs..."

# Get pkl-go version
PKL_GO_VERSION=$(jq -r '.dependencies."pkl-go".version' "$VERSIONS_FILE")

# Fix pkl-go imports
echo "Fixing pkl-go imports to version $PKL_GO_VERSION..."
find "$DEPS_DIR" -name "*.pkl" -type f -exec sed -i "s|import \"external/pkl-go/codegen/src/go.pkl\"|import \"package://pkg.pkl-lang.org/pkl-go/pkl.golang@$PKL_GO_VERSION#/go.pkl\"|g" {} \;

# Fix pkl-pantry imports
echo "Fixing pkl-pantry imports..."
find "$DEPS_DIR" -name "*.pkl" -type f -exec sed -i "s|import \"external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl\"|import \"package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl\"|g" {} \;

echo "Import fixes completed!"
echo ""
echo "Files with updated imports:"
find "$DEPS_DIR" -name "*.pkl" -type f -exec grep -l "package://" {} \;
