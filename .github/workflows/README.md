# GitHub Actions Workflows

## Auto-update Dependencies

**File:** `auto-update-dependencies.yml`

### Purpose
Automatically monitors and updates PKL and pkl-go dependencies to their latest versions.

### Schedule
- **Automatic:** Daily at 00:00 UTC
- **Manual:** Can be triggered via GitHub Actions UI

### What It Updates

#### PKL
- Source: [apple/pkl](https://github.com/apple/pkl) releases
- Files updated:
  - `versions.json` - PKL version reference
  - `build.gradle.kts` - Gradle plugin version
  - All `.pkl` files - `minPklVersion` in `@ModuleInfo`
  - Dependencies and embedded assets

#### pkl-go
- Source: [apple/pkl-go](https://github.com/apple/pkl-go) releases
- Files updated:
  - `versions.json` - pkl-go version reference
  - `go.mod` and `go.sum` - Go module dependencies
  - Dependencies and embedded assets

### Workflow Steps

1. **Check for Updates**
   - Fetches latest releases from GitHub API
   - Compares with current versions in `versions.json`
   - Determines if updates are needed

2. **Apply Updates** (if needed)
   - Updates `versions.json`
   - Updates `go.mod` (for pkl-go)
   - Updates `build.gradle.kts` (for PKL)
   - Updates all `.pkl` files with new `minPklVersion`
   - Runs dependency download scripts
   - Updates import paths
   - Regenerates embedded assets

3. **Run Tests**
   - Executes Go tests to verify compatibility
   - Ensures no breaking changes

4. **Create Pull Request**
   - Generates detailed PR with changelog
   - Includes release notes links
   - Adds labels: `dependencies`, `autoupdate`

### Manual Trigger

To manually trigger the workflow:

1. Go to Actions tab in GitHub
2. Select "Auto-update Dependencies"
3. Click "Run workflow"
4. Select branch (usually `main`)
5. Click "Run workflow"

### PR Format

**Title:**
```
chore: auto-update pkl to X.X.X and pkl-go to X.X.X
```

**Labels:**
- `dependencies`
- `autoupdate`

**Body includes:**
- Version changes summary
- List of all modified files
- Links to release notes
- Test results

### Troubleshooting

**Workflow fails at dependency download:**
- Check network connectivity in GitHub Actions
- Verify GitHub API rate limits
- Check if pkl/pkl-go repositories are accessible

**Tests fail after update:**
- Review the PR to see what changed
- Check release notes for breaking changes
- May need manual intervention to fix compatibility

**No PR created despite new versions:**
- Check workflow logs for errors
- Verify `GITHUB_TOKEN` has proper permissions
- Ensure `versions.json` format is correct

### Dependencies

The workflow requires:
- Ubuntu latest runner
- Go 1.24.4+
- Java 21 (Temurin)
- Gradle (via setup-gradle action)
- jq (for JSON parsing)
- pkl CLI (downloaded during workflow)

### Permissions Required

```yaml
permissions:
  contents: write
  pull-requests: write
```

### Related Files

- `versions.json` - Dependency version configuration
- `scripts/download_deps.sh` - Downloads PKL dependencies
- `scripts/fix_deps_imports.sh` - Updates import paths
- `OFFLINE.md` - Offline dependency documentation

## Release

**File:** `release.yaml`

### Purpose
Builds and releases schema packages, generates documentation.

### Triggers
- Push to tags
- Push to main/master branch
- Pull requests to main/master

See workflow file for detailed configuration.

## Test PKLDoc

**File:** `test-pkldoc.yaml`

### Purpose
Tests PKL documentation generation.

See workflow file for detailed configuration.
