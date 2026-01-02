## Auto-update Dependencies

This PR automatically updates PKL and pkl-go dependencies to their latest versions.

### PKL Update

**Version:** `0.29.1` → `0.30.2`

**Changes:**
- Updated `versions.json`
- Updated `build.gradle.kts`
- Updated `minPklVersion` in all .pkl files (deps/pkl and assets/pkl)
- Downloaded updated PKL dependencies
- Updated import paths
- Updated embedded assets

**Release Notes:** https://github.com/apple/pkl/releases/tag/v0.30.2

### pkl-go Update

**Version:** `0.11.1` → `0.12.1`

**Changes:**
- Updated `versions.json`
- Updated `go.mod` and `go.sum`
- Downloaded updated pkl-go dependencies
- Updated import paths
- Updated embedded assets

**Release Notes:** https://github.com/apple/pkl-go/releases/tag/v0.12.1


### Test Results

All Go tests passed successfully after the update.


---
*This PR was automatically created by the auto-update-dependencies workflow.*
