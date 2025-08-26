package assets

import (
	"fmt"
	"strings"
	"testing"
)

func TestDebugEmbeddedFiles(t *testing.T) {
	// List all external files and look for go.pkl
	files, err := ListExternalFiles()
	if err != nil {
		t.Fatalf("Failed to list external files: %v", err)
	}

	fmt.Printf("Total external files: %d\n", len(files))

	// Look for go.pkl files
	goPklFiles := []string{}
	for _, file := range files {
		if strings.Contains(file, "go.pkl") {
			goPklFiles = append(goPklFiles, file)
		}
	}

	fmt.Printf("go.pkl files found: %v\n", goPklFiles)

	// Try to read the first go.pkl file if it exists
	if len(goPklFiles) > 0 {
		data, err := GetExternalFile(goPklFiles[0])
		if err != nil {
			t.Errorf("Failed to read %s: %v", goPklFiles[0], err)
		} else {
			fmt.Printf("Successfully read %s (%d bytes)\n", goPklFiles[0], len(data))
		}
	}

	// Test reading go.pkl directly
	data, err := GetExternalFile("pkl-go/codegen/src/go.pkl")
	if err != nil {
		t.Errorf("Failed to read pkl-go/codegen/src/go.pkl: %v", err)
	} else {
		fmt.Printf("Successfully read pkl-go/codegen/src/go.pkl (%d bytes)\n", len(data))
	}

	// List first 10 files for debugging
	fmt.Println("First 10 external files:")
	for i, file := range files {
		if i >= 10 {
			break
		}
		fmt.Printf("  %s\n", file)
	}
}
