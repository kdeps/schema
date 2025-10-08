package assets

import (
	"fmt"
	"os"
)

// ExampleCopyAssetsToTempDir demonstrates how to copy assets to a temporary directory.
//
// Usage:
//
//	tempDir, err := assets.CopyAssetsToTempDir()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer os.RemoveAll(tempDir)
//
//	// Now you can use tempDir as the base path for PKL operations
//	toolPath := filepath.Join(tempDir, "Tool.pkl")
func ExampleCopyAssetsToTempDir() {
	// Copy all embedded PKL assets to a temporary directory
	tempDir, err := CopyAssetsToTempDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Assets copied to: %s\n", tempDir)
	// Output example: Assets copied to: /tmp/pkl-assets-123456789
}

// ExampleCopyAssetsToTempDirWithConversion demonstrates how to copy assets
// with package URL to local path conversion applied.
//
// Usage:
//
//	tempDir, err := assets.CopyAssetsToTempDirWithConversion()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer os.RemoveAll(tempDir)
//
//	// All PKL files in tempDir will have package:// URLs converted to local paths
//	workflowPath := filepath.Join(tempDir, "Workflow.pkl")
func ExampleCopyAssetsToTempDirWithConversion() {
	// Copy all embedded PKL assets with conversion applied
	tempDir, err := CopyAssetsToTempDirWithConversion()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Assets copied with conversion to: %s\n", tempDir)
	// Output example: Assets copied with conversion to: /tmp/pkl-assets-987654321
}
