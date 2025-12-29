# Offline-Ready PKL Dependencies

This repository has been configured to work completely offline with local PKL dependencies.

## Overview

The repository now includes:
- **Version Configuration**: `versions.json` specifies dependency versions
- **Local Dependencies**: Complete repositories downloaded to `deps/pkl/external/`
- **Local Imports**: All PKL files use relative imports instead of remote package URLs
- **Embedded Assets**: Go assets include all dependencies for runtime use

## Dependencies

### pkl-go v0.11.1
- **Location**: `deps/pkl/external/pkl-go/`
- **Key Files**: `codegen/src/go.pkl` (Go annotations and code generation support)
- **Import Path**: `external/pkl-go/codegen/src/go.pkl`

### pkl-pantry v1.0.3
- **Location**: `deps/pkl/external/pkl-pantry/`
- **Current Package**: `pkl.experimental.uri` v1.0.3
  - **Files**: `URI.pkl` (URI utilities)
  - **Import Path**: `external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl`

#### Adding More pkl-pantry Packages

To add more packages from pkl-pantry, simply update `versions.json`:

```json
{
  "dependencies": {
    "pkl-pantry": {
      "packages": {
        "pkl.experimental.uri": {
          "version": "1.0.3",
          "files": ["URI.pkl"]
        },
        "pkl.toml": {
          "version": "1.0.2", 
          "files": ["toml.pkl"]
        },
        "pkl.experimental.net": {
          "version": "1.0.1",
          "files": ["ipv4.pkl", "ipv6.pkl", "u128.pkl"]
        }
      }
    }
  }
}
```

Then run:
1. `./scripts/download_deps.sh` - Downloads the complete repository
2. `./scripts/update_imports.sh` - Updates all import paths automatically

The script will handle any number of packages and files within pkl-pantry!

## File Changes

### Updated Import Statements
All PKL files now use local imports instead of remote package URLs:

**Before:**
```pkl
import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.10.0#/go.pkl"
import "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"
```

**After:**
```pkl
import "external/pkl-go/codegen/src/go.pkl"
import "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"
```

## Automated Dependency Updates

### GitHub Actions Workflow
The repository includes an automated workflow (`.github/workflows/auto-update-dependencies.yml`) that:
- **Runs Daily**: Checks for new PKL and pkl-go releases every day at 00:00 UTC
- **Manual Trigger**: Can be triggered manually via GitHub Actions UI
- **Automatic PRs**: Creates pull requests when updates are available

**What it does:**
1. Fetches latest versions from GitHub releases
2. Compares with current versions in `versions.json`
3. Updates all necessary files:
   - `versions.json` - Version references
   - `go.mod` and `go.sum` - Go dependencies
   - `build.gradle.kts` - Gradle PKL plugin version
   - All `.pkl` files - `minPklVersion` in `@ModuleInfo`
4. Downloads updated dependencies
5. Updates import paths
6. Updates embedded assets
7. Runs tests to verify compatibility
8. Creates a PR with detailed changelog and release notes

**PR Details:**
- Title: `chore: auto-update pkl to X.X.X and pkl-go to X.X.X`
- Labels: `dependencies`, `autoupdate`
- Body includes version changes, affected files, and release note links

This ensures the repository stays up-to-date with the latest PKL ecosystem releases automatically.

## Scripts

### Comprehensive Update (Recommended)
```bash
./scripts/update_all.sh
```
Runs all update operations in the correct order:
1. Fetches latest versions from GitHub APIs
2. Updates PKL version references in all files
3. Downloads updated dependencies
4. Updates import paths to local references
5. Updates embedded assets
6. Tests offline functionality and Go builds

### Individual Scripts

#### Update to Latest Versions
```bash
./scripts/update_versions.sh
```
Automatically fetches the latest versions from GitHub APIs:
- **PKL**: Latest release from `apple/pkl`
- **pkl-go**: Latest release from `apple/pkl-go` 
- **pkl-pantry packages**: Latest tags for each configured package
- Updates `versions.json` with fetched versions

#### Update PKL Version References
```bash
./scripts/update_pkl_version.sh
```
Updates `@ModuleInfo { minPklVersion = "..." }` in all PKL files to match the version specified in `versions.json`.

#### Download Dependencies
```bash
./scripts/download_deps.sh
```
Downloads the complete pkl-go and pkl-pantry repositories based on versions specified in `versions.json`.

#### Update Import Paths
```bash
./scripts/update_imports.sh
```
Dynamically reads `versions.json` and updates all PKL files to use local import paths instead of remote package URLs. The script:
- Processes each dependency defined in `versions.json`
- Uses the `files` array to determine which imports to replace
- Maps dependency names to their local directory structures
- Updates all import statements automatically

## Build Process

The `Makefile` has been fully integrated with the offline dependency system:

### Available Commands

```bash
# Generate PKL code with offline dependencies (recommended)
make generate

# Update to latest versions and generate
make generate-latest

# Setup offline dependencies only
make setup-offline

# Update versions and dependencies
make update-deps

# Clean generated files
make clean

# Clean all files including dependencies
make clean-all

# Show help
make help
```

### Default Workflow

1. **`make generate`** (default):
   - Sets up offline dependencies if needed
   - Generates PKL code to Go structs
   - Updates embedded assets
   - Tests offline functionality and Go builds
   - Provides comprehensive feedback with emojis

2. **`make generate-latest`** (CI/CD):
   - Updates to latest versions from GitHub APIs
   - Performs full generation process
   - Perfect for automated builds

### Process Steps

When you run `make generate`, it automatically:
1. 🛠️ Downloads pkl-go and pkl-pantry repositories
2. 🔄 Updates import paths to local references
3. 📦 Generates Go code from PKL files
4. 📁 Updates embedded assets in `assets/pkl/`
5. 🧪 Tests offline functionality
6. 🔨 Tests Go build compatibility
7. 🎉 Reports success or failure

## Go Asset Embedding

The `assets/pkl_assets.go` file embeds all PKL files and dependencies:
```go
//go:embed pkl
var PKLFS embed.FS
```

This includes:
- All main PKL files
- Complete pkl-go repository
- Complete pkl-pantry repository

## Benefits

1. **True Offline Operation**: No network requests required at build or runtime
2. **Version Control**: Dependencies are pinned to specific versions
3. **Self-Contained**: Everything needed is included in the repository
4. **Fast Builds**: No dependency downloading during build process

## Usage

All existing functionality works the same way, but now operates completely offline:
```go
import "github.com/kdeps/schema/assets"

content, err := assets.GetPKLFileAsString("Project.pkl")
```

## Version Updates

To update dependency versions:
1. Update `versions.json`
2. Run `./scripts/download_deps.sh`
3. Run `./scripts/update_imports.sh` if import paths change
4. Test with `pkl eval deps/pkl/Tool.pkl --no-cache`