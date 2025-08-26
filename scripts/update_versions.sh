#!/bin/bash

# Script to update versions.json with the latest available versions
set -e

VERSIONS_FILE="versions.json"

echo "Updating versions.json with latest available versions..."

# Function to get latest GitHub release version
get_latest_github_release() {
    local repo=$1
    curl -s "https://api.github.com/repos/$repo/releases/latest" | jq -r '.tag_name' | sed 's/^v//'
}

# Function to get latest GitHub tag version
get_latest_github_tag() {
    local repo=$1
    curl -s "https://api.github.com/repos/$repo/tags" | jq -r '.[0].name' | sed 's/^v//'
}

# Function to get latest pkl-pantry package version
get_latest_pantry_package_version() {
    local package=$1
    curl -s "https://api.github.com/repos/apple/pkl-pantry/tags" | \
        jq -r ".[] | select(.name | startswith(\"$package@\")) | .name" | \
        head -1 | sed "s/^$package@//"
}

echo "Fetching latest versions..."

# Get latest PKL version
echo "  Checking PKL version..."
PKL_LATEST=$(get_latest_github_release "apple/pkl")
if [ -n "$PKL_LATEST" ] && [ "$PKL_LATEST" != "null" ]; then
    echo "    Latest PKL version: $PKL_LATEST"
    jq ".pkl.version = \"$PKL_LATEST\"" "$VERSIONS_FILE" > "$VERSIONS_FILE.tmp" && mv "$VERSIONS_FILE.tmp" "$VERSIONS_FILE"
else
    echo "    Could not fetch PKL version, keeping current"
fi

# Get latest pkl-go version
echo "  Checking pkl-go version..."
PKL_GO_LATEST=$(get_latest_github_release "apple/pkl-go")
if [ -n "$PKL_GO_LATEST" ] && [ "$PKL_GO_LATEST" != "null" ]; then
    echo "    Latest pkl-go version: $PKL_GO_LATEST"
    jq ".dependencies[\"pkl-go\"].version = \"$PKL_GO_LATEST\" | .dependencies[\"pkl-go\"].url = \"package://pkg.pkl-lang.org/pkl-go/pkl.golang@$PKL_GO_LATEST\"" "$VERSIONS_FILE" > "$VERSIONS_FILE.tmp" && mv "$VERSIONS_FILE.tmp" "$VERSIONS_FILE"
else
    echo "    Could not fetch pkl-go version, keeping current"
fi

# Get latest pkl-pantry package versions
echo "  Checking pkl-pantry packages..."
jq -r '.dependencies["pkl-pantry"].packages | keys[]' "$VERSIONS_FILE" | while read -r package; do
    echo "    Checking $package..."
    PACKAGE_LATEST=$(get_latest_pantry_package_version "$package")
    if [ -n "$PACKAGE_LATEST" ] && [ "$PACKAGE_LATEST" != "null" ]; then
        echo "      Latest $package version: $PACKAGE_LATEST"
        jq ".dependencies[\"pkl-pantry\"].packages[\"$package\"].version = \"$PACKAGE_LATEST\"" "$VERSIONS_FILE" > "$VERSIONS_FILE.tmp" && mv "$VERSIONS_FILE.tmp" "$VERSIONS_FILE"
    else
        echo "      Could not fetch $package version, keeping current"
    fi
done

echo "Versions updated successfully!"

# Show updated versions
echo ""
echo "Current versions in versions.json:"
echo "PKL: $(jq -r '.pkl.version' "$VERSIONS_FILE")"
echo "pkl-go: $(jq -r '.dependencies["pkl-go"].version' "$VERSIONS_FILE")"
jq -r '.dependencies["pkl-pantry"].packages | to_entries[] | "pkl-pantry/" + .key + ": " + .value.version' "$VERSIONS_FILE"