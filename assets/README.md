# PKL Assets Package

This package provides embedded PKL schema files and utilities for working with them.

## Features

### Embedded File System
All PKL schema files are embedded using Go's `embed.FS`, making them available at runtime without external dependencies.

### Reading Embedded Files
```go
// Read a PKL file from the embedded filesystem
data, err := assets.GetPKLFile("Tool.pkl")

// Read as a string
content, err := assets.GetPKLFileAsString("Workflow.pkl")

// List all PKL files
files, err := assets.ListPKLFiles()
```

### Temporary Directory Operations

The package provides functions to copy embedded assets to temporary directories with complete paths.

#### 1. Copy Assets to Temp Directory
```go
// Copy all embedded assets to a temporary directory
tempDir, err := assets.CopyAssetsToTempDir()
if err != nil {
    log.Fatal(err)
}
defer os.RemoveAll(tempDir) // Clean up when done

// Use the complete temp directory path
toolPath := filepath.Join(tempDir, "Tool.pkl")
// toolPath will be something like: /tmp/pkl-assets-123456789/Tool.pkl
```

**Key Features:**
- Creates a unique temporary directory with prefix `pkl-assets-*`
- Copies all PKL files maintaining directory structure
- Returns the **complete path** to the temporary directory
- Automatically creates all necessary subdirectories
- Preserves file permissions (0644 for files, 0755 for directories)

#### 2. Copy Assets with Package URL Conversion
```go
// Copy assets with package:// URLs converted to local paths
tempDir, err := assets.CopyAssetsToTempDirWithConversion()
if err != nil {
    log.Fatal(err)
}
defer os.RemoveAll(tempDir) // Clean up when done

// All PKL files will have package URLs converted to local paths
workflowPath := filepath.Join(tempDir, "Workflow.pkl")
```

**Key Features:**
- Same as `CopyAssetsToTempDir` but applies package URL conversion
- Converts `package://pkg.pkl-lang.org/...` to `external/...`
- Converts `package://schema.kdeps.com/core@...` to direct local paths
- Ensures offline compatibility

#### 3. Write Assets to Specific Directory
```go
// Write assets to a specific directory (creates if doesn't exist)
targetDir := "/path/to/my/pkl-files"
err := assets.WriteAssetsToDir(targetDir)
if err != nil {
    log.Fatal(err)
}

// Assets are now available at /path/to/my/pkl-files/
```

**Key Features:**
- Write to a specific directory path (not temp)
- Creates directory if it doesn't exist
- Useful for persistent storage or build processes
- No cleanup needed (you control the directory)

#### 4. Write Assets with Conversion to Specific Directory
```go
// Write assets with conversion to a specific directory
targetDir := "./local-pkl-files"
err := assets.WriteAssetsToDirWithConversion(targetDir)
if err != nil {
    log.Fatal(err)
}

// All PKL files in ./local-pkl-files/ use local paths
```

**Key Features:**
- Same as `WriteAssetsToDir` but with package URL conversion
- Perfect for creating offline-ready PKL file sets
- Useful for distribution or bundling

#### 5. Helper Function - Get Single File with Temp Dir
```go
// Get a single file with automatic temp dir management
content, tempDir, cleanup, err := assets.GetPKLFileFromTempDir("Tool.pkl")
if err != nil {
    log.Fatal(err)
}
defer cleanup() // Automatic cleanup

// Use the content
fmt.Println(content)

// Or access other files from the same temp dir
otherFile := filepath.Join(tempDir, "Resource.pkl")
```

**Key Features:**
- One-liner to get a file in a temp directory
- Returns content, temp dir path, and cleanup function
- All other files also available in the temp dir
- Cleanup is simple with the returned function

#### 6. Helper Function - Get Single File with Conversion
```go
// Get a single file with conversion and temp dir management
content, tempDir, cleanup, err := assets.GetPKLFileFromTempDirWithConversion("Workflow.pkl")
if err != nil {
    log.Fatal(err)
}
defer cleanup()

// Content has all package URLs converted to local paths
// All other files in tempDir also have conversions applied
```

**Key Features:**
- Same as `GetPKLFileFromTempDir` but with URL conversion
- Perfect for quick offline-ready access
- All files in temp dir are converted

### Package URL Conversion

Convert package:// URLs to local paths for offline usage:

```go
// Convert package URLs in content
converted := assets.ConvertPackageURLsToLocalPaths(content)

// Convert import/amends statements
converted = assets.ConvertImportStatements(content)

// Validate that no package URLs remain
isValid, remaining := assets.ValidateLocalPaths(content)
```

### Offline Compatibility

Ensure all PKL files use local paths:

```go
// Validate all PKL files
err := assets.EnsureOfflineCompatibility()
if err != nil {
    // Some files still contain package:// URLs
    log.Fatal(err)
}
```

## Directory Structure

When you copy assets to a temp directory, the structure will be:

```
/tmp/pkl-assets-123456789/    (complete temp directory path returned)
├── Tool.pkl
├── Resource.pkl
├── Workflow.pkl
├── LLM.pkl
├── ... (other PKL files)
└── external/
    ├── pkl-go/
    │   └── codegen/src/
    │       └── go.pkl
    └── pkl-pantry/
        └── packages/
            └── ...
```

## Usage Examples

### Example 1: Process PKL Files in Temp Directory
```go
package main

import (
    "log"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/kdeps/schema/assets"
)

func main() {
    // Copy assets to temp directory
    tempDir, err := assets.CopyAssetsToTempDir()
    if err != nil {
        log.Fatal(err)
    }
    defer os.RemoveAll(tempDir)

    // Process PKL files using the complete temp directory path
    toolFile := filepath.Join(tempDir, "Tool.pkl")
    cmd := exec.Command("pkl", "eval", toolFile)
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("PKL output: %s", output)
}
```

### Example 2: Offline PKL Processing
```go
package main

import (
    "log"
    "os"

    "github.com/kdeps/schema/assets"
)

func main() {
    // Copy assets with all package URLs converted to local paths
    tempDir, err := assets.CopyAssetsToTempDirWithConversion()
    if err != nil {
        log.Fatal(err)
    }
    defer os.RemoveAll(tempDir)

    log.Printf("Offline-ready PKL files available at: %s", tempDir)

    // All files in tempDir can now be processed without internet access
    // because all package:// URLs have been converted to local paths
}
```

## Testing

Run the test suite:

```bash
go test -v ./assets
```

Key test cases:
- `TestCopyAssetsToTempDir` - Verifies temp directory creation and asset copying
- `TestCopyAssetsToTempDirWithConversion` - Verifies conversion is applied
- `TestTempDirCleanup` - Verifies proper cleanup of temp directories
- All other existing tests for file reading and conversion

## Important Notes

1. **Always clean up temp directories** using `defer os.RemoveAll(tempDir)`
2. **Complete paths are returned** - the returned `tempDir` is ready to use
3. **Unique directories** - each call creates a new unique temp directory
4. **Thread-safe** - multiple goroutines can call these functions concurrently
5. **Error handling** - functions clean up on error automatically

## Which Function Should I Use?

Here's a quick guide to help you choose the right function:

| **Use Case** | **Function** | **Returns** |
|--------------|--------------|-------------|
| Need all files in a temp directory | `CopyAssetsToTempDir()` | Temp dir path |
| Need all files offline-ready in temp dir | `CopyAssetsToTempDirWithConversion()` | Temp dir path |
| Need all files in a specific permanent location | `WriteAssetsToDir(path)` | Error only |
| Need all files offline-ready in specific location | `WriteAssetsToDirWithConversion(path)` | Error only |
| Need one file quickly with temp dir | `GetPKLFileFromTempDir(file)` | Content, dir, cleanup, error |
| Need one file offline-ready with temp dir | `GetPKLFileFromTempDirWithConversion(file)` | Content, dir, cleanup, error |
| Just read a file from embed (in-memory) | `GetPKLFile(file)` | Bytes, error |
| Just read a file as string from embed | `GetPKLFileAsString(file)` | String, error |

**Quick Decision Tree:**
1. **Do you need files on disk?**
   - No → Use `GetPKLFile()` or `GetPKLFileAsString()` for in-memory access
   - Yes → Continue...

2. **Do you need offline compatibility (no package URLs)?**
   - Yes → Use a function with "Conversion" in the name
   - No → Use the regular version

3. **Do you want to control the directory path?**
   - Yes → Use `WriteAssetsToDir()` or `WriteAssetsToDirWithConversion()`
   - No (use temp dir) → Continue...

4. **Do you need just one file or all files?**
   - One file → Use `GetPKLFileFromTempDir()` or `GetPKLFileFromTempDirWithConversion()`
   - All files → Use `CopyAssetsToTempDir()` or `CopyAssetsToTempDirWithConversion()`

## API Reference

### Temporary Directory Functions

- `CopyAssetsToTempDir() (string, error)` - Copy assets to temp dir, returns complete path
- `CopyAssetsToTempDirWithConversion() (string, error)` - Copy to temp dir with URL conversion
- `WriteAssetsToDir(targetDir string) error` - Write assets to specific directory
- `WriteAssetsToDirWithConversion(targetDir string) error` - Write to directory with conversion
- `GetPKLFileFromTempDir(filename string) (content string, tempDir string, cleanup func(), err error)` - Get file with temp dir and cleanup
- `GetPKLFileFromTempDirWithConversion(filename string) (content string, tempDir string, cleanup func(), err error)` - Get file with conversion and temp dir

### Reading Functions

- `GetPKLFile(filename string) ([]byte, error)` - Read a PKL file from embedded FS
- `GetPKLFileAsString(filename string) (string, error)` - Read as string from embedded FS
- `GetPKLFileFromPKL(filename string) ([]byte, error)` - Read from pkl directory specifically
- `GetExternalFile(filename string) ([]byte, error)` - Read from external directory
- `GetPKLFileAsStringWithLocalPaths(filename string) (string, error)` - Read with conversion applied
- `GetPKLFileWithFullConversion(filename string) (string, error)` - Read with all conversions

### Listing Functions

- `ListPKLFiles() ([]string, error)` - List all PKL files
- `ListExternalFiles() ([]string, error)` - List external dependency files

### Conversion Functions

- `ConvertPackageURLsToLocalPaths(content string) string` - Convert package URLs to local paths
- `ConvertImportStatements(content string) string` - Convert import/amends statements
- `ValidateLocalPaths(content string) (bool, []string)` - Check for remaining package URLs
- `ConvertAllPKLFiles() (map[string]string, error)` - Convert all PKL files and return map

### Validation Functions

- `EnsureOfflineCompatibility() error` - Validate offline compatibility
- `ValidateAllPKLFiles() (map[string][]string, error)` - Check all files for package URLs
