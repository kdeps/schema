# PKL Assets - Quick Reference Card

## 🚀 Quick Start (Choose One)

```go
// 1️⃣  Read from embed (in-memory, existing)
content, _ := assets.GetPKLFileAsString("Tool.pkl")

// 2️⃣  Get one file in temp dir (new, easiest)
content, _, cleanup, _ := assets.GetPKLFileFromTempDir("Tool.pkl")
defer cleanup()

// 3️⃣  Get all files in temp dir (new)
tempDir, _ := assets.CopyAssetsToTempDir()
defer os.RemoveAll(tempDir)

// 4️⃣  Write to specific directory (new)
assets.WriteAssetsToDir("/path/to/output")
```

## 📋 All Functions at a Glance

| Function | What It Does | Returns |
|----------|--------------|---------|
| **Temp Dir - Auto Location** |||
| `CopyAssetsToTempDir()` | All files → temp dir | temp path |
| `CopyAssetsToTempDirWithConversion()` | All files offline → temp dir | temp path |
| `GetPKLFileFromTempDir(file)` | One file → temp dir + cleanup | content, path, cleanup |
| `GetPKLFileFromTempDirWithConversion(file)` | One file offline → temp dir + cleanup | content, path, cleanup |
| **Specific Dir - You Control** |||
| `WriteAssetsToDir(path)` | All files → your dir | error |
| `WriteAssetsToDirWithConversion(path)` | All files offline → your dir | error |
| **Embedded FS - In Memory** |||
| `GetPKLFile(file)` | Read file from embed | bytes |
| `GetPKLFileAsString(file)` | Read file as string | string |
| **Utility** |||
| `ListPKLFiles()` | List all PKL files | []string |
| `ConvertPackageURLsToLocalPaths(content)` | Convert URLs | string |
| `ValidateLocalPaths(content)` | Check for URLs | bool, []string |

## 🎯 Decision Matrix

```
Need files on disk?
├─ No  → Use GetPKLFile() or GetPKLFileAsString()
└─ Yes → Need offline (no package URLs)?
         ├─ No  → Control directory?
         │        ├─ No  → CopyAssetsToTempDir()
         │        └─ Yes → WriteAssetsToDir(path)
         └─ Yes → Control directory?
                  ├─ No  → CopyAssetsToTempDirWithConversion()
                  └─ Yes → WriteAssetsToDirWithConversion(path)

Need just one file quickly?
└─ Use GetPKLFileFromTempDir(file) or GetPKLFileFromTempDirWithConversion(file)
```

## 💡 Common Patterns

### Pattern 1: Quick File Access
```go
content, _, cleanup, err := assets.GetPKLFileFromTempDir("Workflow.pkl")
if err != nil {
    return err
}
defer cleanup()
// Use content...
```

### Pattern 2: Process All Files
```go
tempDir, err := assets.CopyAssetsToTempDir()
if err != nil {
    return err
}
defer os.RemoveAll(tempDir)

files, _ := filepath.Glob(filepath.Join(tempDir, "*.pkl"))
for _, file := range files {
    // Process each file...
}
```

### Pattern 3: Offline Bundle Creation
```go
func createBundle(outputDir string) error {
    return assets.WriteAssetsToDirWithConversion(outputDir)
}
```

### Pattern 4: CLI Tool
```go
func runPKL(filename string) error {
    tempDir, err := assets.CopyAssetsToTempDir()
    if err != nil {
        return err
    }
    defer os.RemoveAll(tempDir)

    cmd := exec.Command("pkl", "eval", filepath.Join(tempDir, filename))
    return cmd.Run()
}
```

### Pattern 5: Testing
```go
func TestWorkflow(t *testing.T) {
    content, tempDir, cleanup, err := assets.GetPKLFileFromTempDir("Test.pkl")
    require.NoError(t, err)
    defer cleanup()

    // Modify file in tempDir if needed
    modified := filepath.Join(tempDir, "Modified.pkl")
    os.WriteFile(modified, []byte(content+"extra"), 0644)

    // Test...
}
```

## ⚡ Performance Tips

1. **Reuse temp directories** - Don't create new ones unnecessarily
2. **Use helpers** - `GetPKLFileFromTempDir()` is optimized for single file access
3. **Defer cleanup** - Always use `defer` to avoid leaks
4. **Cache when possible** - If processing multiple times, create once

## ⚠️ Important Notes

✅ **Always return complete paths** - All functions return full absolute paths
✅ **Always cleanup** - Use `defer os.RemoveAll(tempDir)` or `defer cleanup()`
✅ **Thread-safe** - Safe to call from multiple goroutines
✅ **Auto-create dirs** - WriteAssetsToDir creates directories as needed
✅ **Backwards compatible** - All old code still works

## 🔍 Quick Troubleshooting

**Q: Files not found in temp dir?**
A: Make sure you're using the returned `tempDir` path, not a relative path.

**Q: Temp directories filling up disk?**
A: Make sure you're calling `defer os.RemoveAll(tempDir)` or `defer cleanup()`.

**Q: Package URLs still in files?**
A: Use the "WithConversion" variant of the function.

**Q: Permission errors?**
A: Temp directories are created with 0755, files with 0644. Check your umask.

## 📚 Full Documentation

For complete documentation, see:
- `assets/README.md` - Comprehensive guide
- `INTEGRATION_SUMMARY.md` - Integration details
- `TEMPDIR_FEATURE.md` - Feature documentation

## 🧪 Example Code

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/kdeps/schema/assets"
)

func main() {
    // Example 1: Read from embed
    content, err := assets.GetPKLFileAsString("Tool.pkl")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Content length:", len(content))

    // Example 2: Use temp directory
    tempDir, err := assets.CopyAssetsToTempDir()
    if err != nil {
        log.Fatal(err)
    }
    defer os.RemoveAll(tempDir)

    fmt.Println("Temp dir:", tempDir)
    files, _ := filepath.Glob(filepath.Join(tempDir, "*.pkl"))
    fmt.Println("Files:", len(files))

    // Example 3: Helper function
    content2, _, cleanup, err := assets.GetPKLFileFromTempDir("Workflow.pkl")
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()

    fmt.Println("Workflow content length:", len(content2))

    // Example 4: Create offline bundle
    if err := assets.WriteAssetsToDirWithConversion("./output"); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Bundle created in ./output")
}
```

---

**For more examples and detailed documentation, see `assets/README.md`**
