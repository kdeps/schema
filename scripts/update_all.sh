#!/bin/bash

# Comprehensive script to update all versions and dependencies
set -e

echo "🚀 Starting comprehensive update process..."
echo ""

# Step 1: Update versions to latest
echo "📦 Step 1: Fetching latest versions..."
./scripts/update_versions.sh
echo ""

# Step 2: Update PKL version references in files
echo "🏷️  Step 2: Updating PKL version references in files..."
./scripts/update_pkl_version.sh
echo ""

# Step 3: Download updated dependencies
echo "⬇️  Step 3: Downloading updated dependencies..."
./scripts/download_deps.sh
echo ""

# Step 4: Update import paths
echo "🔄 Step 4: Updating import paths..."
./scripts/update_imports.sh
echo ""

# Step 5: Update assets
echo "📁 Step 5: Updating embedded assets..."
mkdir -p assets/pkl && cp deps/pkl/*.pkl assets/pkl/
cp -r deps/pkl/external assets/pkl/
echo "Assets updated successfully!"
echo ""

# Step 6: Test offline functionality
echo "🧪 Step 6: Testing offline functionality..."
if pkl eval deps/pkl/Tool.pkl --no-cache --format json >/dev/null 2>&1; then
    echo "✅ Offline functionality test passed!"
else
    echo "❌ Offline functionality test failed!"
    exit 1
fi
echo ""

# Step 7: Test Go build
echo "🔨 Step 7: Testing Go build..."
if go build ./assets >/dev/null 2>&1; then
    echo "✅ Go build test passed!"
else
    echo "❌ Go build test failed!"
    exit 1
fi
echo ""

echo "🎉 All updates completed successfully!"
echo ""

# Show final status
echo "📊 Final Status:"
echo "PKL version: $(jq -r '.pkl.version' versions.json)"
echo "pkl-go version: $(jq -r '.dependencies["pkl-go"].version' versions.json)"
jq -r '.dependencies["pkl-pantry"].packages | to_entries[] | "pkl-pantry/" + .key + ": " + .value.version' versions.json
echo ""
echo "✨ Your repository is now up-to-date and fully offline-ready!"