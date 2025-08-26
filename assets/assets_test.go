package assets

import (
	"strings"
	"testing"
)

func TestGetPKLFile(t *testing.T) {
	// Test reading a PKL file from the pkl directory
	data, err := GetPKLFile("Tool.pkl")
	if err != nil {
		t.Fatalf("Failed to read Tool.pkl: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Tool.pkl is empty")
	}

	content := string(data)
	if !strings.Contains(content, "module") {
		t.Fatal("Tool.pkl doesn't contain expected content")
	}
}

func TestGetPKLFileFromPKL(t *testing.T) {
	// Test reading a PKL file specifically from the pkl directory
	data, err := GetPKLFileFromPKL("Resource.pkl")
	if err != nil {
		t.Fatalf("Failed to read Resource.pkl: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Resource.pkl is empty")
	}

	content := string(data)
	if !strings.Contains(content, "module") {
		t.Fatal("Resource.pkl doesn't contain expected content")
	}
}

func TestGetExternalFile(t *testing.T) {
	// Test reading a file from the external directory
	data, err := GetExternalFile("pkl-go/codegen/src/go.pkl")
	if err != nil {
		t.Fatalf("Failed to read external go.pkl: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("external go.pkl is empty")
	}

	content := string(data)
	if !strings.Contains(content, "module") {
		t.Fatal("external go.pkl doesn't contain expected content")
	}
}

func TestGetPKLFileAsString(t *testing.T) {
	// Test reading a PKL file as string
	content, err := GetPKLFileAsString("Workflow.pkl")
	if err != nil {
		t.Fatalf("Failed to read Workflow.pkl as string: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("Workflow.pkl string is empty")
	}

	if !strings.Contains(content, "module") {
		t.Fatal("Workflow.pkl string doesn't contain expected content")
	}
}

func TestListPKLFiles(t *testing.T) {
	// Test listing PKL files
	files, err := ListPKLFiles()
	if err != nil {
		t.Fatalf("Failed to list PKL files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No PKL files found")
	}

	// Check for expected files
	expectedFiles := []string{"Tool.pkl", "Resource.pkl", "Workflow.pkl", "LLM.pkl"}
	foundFiles := make(map[string]bool)
	for _, file := range files {
		foundFiles[file] = true
	}

	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("Expected file %s not found in PKL files", expected)
		}
	}

	t.Logf("Found %d PKL files: %v", len(files), files)
}

func TestListExternalFiles(t *testing.T) {
	// Test listing external files
	files, err := ListExternalFiles()
	if err != nil {
		t.Fatalf("Failed to list external files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No external files found")
	}

	// Check for expected external files
	expectedPatterns := []string{
		"pkl-go/codegen/src/go.pkl",
		"pkl-pantry/packages/",
	}

	foundPatterns := make(map[string]bool)
	for _, file := range files {
		for _, pattern := range expectedPatterns {
			if strings.Contains(file, pattern) {
				foundPatterns[pattern] = true
			}
		}
	}

	for _, pattern := range expectedPatterns {
		if !foundPatterns[pattern] {
			t.Errorf("Expected pattern %s not found in external files", pattern)
		}
	}

	t.Logf("Found %d external files", len(files))
}

func TestGetPKLFileFallback(t *testing.T) {
	// Test that GetPKLFile falls back to external directory
	// This should work because go.pkl exists in external but not in pkl
	data, err := GetPKLFile("pkl-go/codegen/src/go.pkl")
	if err != nil {
		t.Fatalf("Failed to read go.pkl with fallback: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("go.pkl is empty")
	}

	content := string(data)
	if !strings.Contains(content, "module") {
		t.Fatal("go.pkl doesn't contain expected content")
	}
}

func TestSpecificExternalFiles(t *testing.T) {
	// Test specific external files that should exist
	testCases := []struct {
		name     string
		filepath string
	}{
		{"pkl-go go.pkl", "pkl-go/codegen/src/go.pkl"},
		{"pkl-pantry URI.pkl", "pkl-pantry/packages/pkl.experimental.uri/URI.pkl"},
		{"pkl-pantry JsonSchema.pkl", "pkl-pantry/packages/org.json_schema/JsonSchema.pkl"},
		{"pkl-pantry Schema.pkl", "pkl-pantry/packages/org.openapis.v3/Schema.pkl"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := GetExternalFile(tc.filepath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", tc.filepath, err)
			}

			if len(data) == 0 {
				t.Fatalf("%s is empty", tc.filepath)
			}

			content := string(data)
			if !strings.Contains(content, "module") {
				t.Fatalf("%s doesn't contain expected content", tc.filepath)
			}
		})
	}
}

func TestPKLFileContentValidation(t *testing.T) {
	// Test that PKL files contain valid PKL syntax
	testFiles := []string{"Tool.pkl", "Resource.pkl", "Workflow.pkl", "LLM.pkl"}

	for _, filename := range testFiles {
		t.Run(filename, func(t *testing.T) {
			content, err := GetPKLFileAsString(filename)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", filename, err)
			}

			// Basic PKL syntax validation
			if !strings.Contains(content, "module") {
				t.Errorf("%s doesn't contain 'module' keyword", filename)
			}

			if !strings.Contains(content, "{") || !strings.Contains(content, "}") {
				t.Errorf("%s doesn't contain expected PKL syntax", filename)
			}
		})
	}
}

func TestExternalFileContentValidation(t *testing.T) {
	// Test that external files contain valid content
	testFiles := []struct {
		name     string
		filepath string
	}{
		{"go.pkl", "pkl-go/codegen/src/go.pkl"},
		{"URI.pkl", "pkl-pantry/packages/pkl.experimental.uri/URI.pkl"},
	}

	for _, tc := range testFiles {
		t.Run(tc.name, func(t *testing.T) {
			content, err := GetExternalFile(tc.filepath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", tc.filepath, err)
			}

			contentStr := string(content)
			if !strings.Contains(contentStr, "module") {
				t.Errorf("%s doesn't contain 'module' keyword", tc.filepath)
			}
		})
	}
}

func BenchmarkGetPKLFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetPKLFile("Tool.pkl")
		if err != nil {
			b.Fatalf("Failed to read Tool.pkl: %v", err)
		}
	}
}

func BenchmarkListPKLFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ListPKLFiles()
		if err != nil {
			b.Fatalf("Failed to list PKL files: %v", err)
		}
	}
}
