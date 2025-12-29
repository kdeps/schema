#!/bin/bash

# Script to test the auto-update workflow logic locally
# This simulates what the GitHub Actions workflow does

set -e

echo "🧪 Testing Auto-update Workflow Logic"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ jq is not installed. Please install jq first.${NC}"
    exit 1
fi

echo "📦 Step 1: Fetching latest versions from GitHub..."
echo "---------------------------------------------------"

# Fetch latest PKL version
echo -n "  Fetching latest PKL version... "
LATEST_PKL=$(curl -s https://api.github.com/repos/apple/pkl/releases/latest | jq -r '.tag_name')
if [ -z "$LATEST_PKL" ] || [ "$LATEST_PKL" = "null" ]; then
    echo -e "${RED}❌ Failed${NC}"
    exit 1
fi
# Remove 'v' prefix
LATEST_PKL_NUM="${LATEST_PKL#v}"
echo -e "${GREEN}✓${NC} $LATEST_PKL_NUM"

# Fetch latest pkl-go version
echo -n "  Fetching latest pkl-go version... "
LATEST_PKL_GO=$(curl -s https://api.github.com/repos/apple/pkl-go/releases/latest | jq -r '.tag_name')
if [ -z "$LATEST_PKL_GO" ] || [ "$LATEST_PKL_GO" = "null" ]; then
    echo -e "${RED}❌ Failed${NC}"
    exit 1
fi
# Remove 'v' prefix
LATEST_PKL_GO_NUM="${LATEST_PKL_GO#v}"
echo -e "${GREEN}✓${NC} $LATEST_PKL_GO_NUM"

echo ""
echo "📋 Step 2: Reading current versions from versions.json..."
echo "---------------------------------------------------"

# Get current PKL version
echo -n "  Current PKL version... "
CURRENT_PKL=$(jq -r '.pkl.version' versions.json)
if [ -z "$CURRENT_PKL" ] || [ "$CURRENT_PKL" = "null" ]; then
    echo -e "${RED}❌ Failed to read${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} $CURRENT_PKL"

# Get current pkl-go version
echo -n "  Current pkl-go version... "
CURRENT_PKL_GO=$(jq -r '.dependencies."pkl-go".version' versions.json)
if [ -z "$CURRENT_PKL_GO" ] || [ "$CURRENT_PKL_GO" = "null" ]; then
    echo -e "${RED}❌ Failed to read${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} $CURRENT_PKL_GO"

echo ""
echo "🔍 Step 3: Comparing versions..."
echo "---------------------------------------------------"

PKL_UPDATE="false"
PKL_GO_UPDATE="false"

# Compare PKL versions
echo -n "  PKL: $CURRENT_PKL vs $LATEST_PKL_NUM... "
if [ "$LATEST_PKL_NUM" != "$CURRENT_PKL" ]; then
    echo -e "${YELLOW}UPDATE AVAILABLE${NC}"
    PKL_UPDATE="true"
else
    echo -e "${GREEN}UP TO DATE${NC}"
fi

# Compare pkl-go versions
echo -n "  pkl-go: $CURRENT_PKL_GO vs $LATEST_PKL_GO_NUM... "
if [ "$LATEST_PKL_GO_NUM" != "$CURRENT_PKL_GO" ]; then
    echo -e "${YELLOW}UPDATE AVAILABLE${NC}"
    PKL_GO_UPDATE="true"
else
    echo -e "${GREEN}UP TO DATE${NC}"
fi

echo ""
echo "📊 Summary"
echo "---------------------------------------------------"

if [ "$PKL_UPDATE" = "true" ] || [ "$PKL_GO_UPDATE" = "true" ]; then
    echo -e "${YELLOW}Updates available:${NC}"

    if [ "$PKL_UPDATE" = "true" ]; then
        echo "  • PKL: $CURRENT_PKL → $LATEST_PKL_NUM"
    fi

    if [ "$PKL_GO_UPDATE" = "true" ]; then
        echo "  • pkl-go: $CURRENT_PKL_GO → $LATEST_PKL_GO_NUM"
    fi

    echo ""
    echo "The workflow would:"
    echo "  1. Update versions.json"

    if [ "$PKL_GO_UPDATE" = "true" ]; then
        echo "  2. Update go.mod to github.com/apple/pkl-go@v$LATEST_PKL_GO_NUM"
        echo "  3. Run go mod tidy"
    fi

    if [ "$PKL_UPDATE" = "true" ]; then
        echo "  4. Update build.gradle.kts"
        echo "  5. Update minPklVersion in all .pkl files"
    fi

    echo "  6. Run scripts/download_deps.sh"
    echo "  7. Run scripts/fix_deps_imports.sh"
    echo "  8. Update embedded assets"
    echo "  9. Run Go tests"
    echo "  10. Create a PR with detailed changelog"

    echo ""
    echo -e "${YELLOW}Would you like to see what the PR title would be?${NC}"

    # Generate PR title
    COMPONENTS=()
    if [ "$PKL_UPDATE" = "true" ]; then
        COMPONENTS+=("PKL to $LATEST_PKL_NUM")
    fi
    if [ "$PKL_GO_UPDATE" = "true" ]; then
        COMPONENTS+=("pkl-go to $LATEST_PKL_GO_NUM")
    fi

    if [ ${#COMPONENTS[@]} -eq 1 ]; then
        PR_TITLE="chore: auto-update ${COMPONENTS[0]}"
    else
        PR_TITLE="chore: auto-update ${COMPONENTS[0]} and ${COMPONENTS[1]}"
    fi

    echo ""
    echo "PR Title: $PR_TITLE"
    echo ""

    echo -e "${GREEN}✅ Workflow logic test PASSED - Updates would be processed${NC}"
    exit 0
else
    echo -e "${GREEN}✅ All dependencies are up to date${NC}"
    echo ""
    echo "The workflow would:"
    echo "  1. Check versions"
    echo "  2. Determine no updates needed"
    echo "  3. Exit without creating a PR"
    echo ""
    echo -e "${GREEN}✅ Workflow logic test PASSED - No action needed${NC}"
    exit 0
fi
