package assets

import (
	"os"
	"path/filepath"
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

func TestConvertPackageURLsToLocalPaths(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "pkl-go package URL",
			input:    `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"`,
			expected: `import "external/pkl-go/codegen/src/go.pkl"`,
		},
		{
			name:     "pkl-pantry package URL",
			input:    `import "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"`,
			expected: `import "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"`,
		},
		{
			name:     "multiple package URLs",
			input:    `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"\nimport "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"`,
			expected: `import "external/pkl-go/codegen/src/go.pkl"\nimport "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"`,
		},
		{
			name:     "no package URLs",
			input:    `import "local/file.pkl"`,
			expected: `import "local/file.pkl"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertPackageURLsToLocalPaths(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestConvertImportStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "import statement with pkl-go",
			input:    `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"`,
			expected: `import "external/pkl-go/codegen/src/go.pkl"`,
		},
		{
			name:     "amends statement with pkl-pantry",
			input:    `amends "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"`,
			expected: `amends "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"`,
		},
		{
			name:     "mixed import and amends",
			input:    `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"\namends "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"`,
			expected: `import "external/pkl-go/codegen/src/go.pkl"\namends "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertImportStatements(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestValidateLocalPaths(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectValid    bool
		expectedMatches int
	}{
		{
			name:           "all local paths",
			input:          `import "external/pkl-go/codegen/src/go.pkl"\nimport "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"`,
			expectValid:    true,
			expectedMatches: 0,
		},
		{
			name:           "contains package URL",
			input:          `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"`,
			expectValid:    false,
			expectedMatches: 1,
		},
		{
			name:           "mixed paths",
			input:          `import "external/pkl-go/codegen/src/go.pkl"\nimport "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"`,
			expectValid:    false,
			expectedMatches: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isValid, matches := ValidateLocalPaths(tc.input)
			if isValid != tc.expectValid {
				t.Errorf("Expected valid: %v, got: %v", tc.expectValid, isValid)
			}
			if len(matches) != tc.expectedMatches {
				t.Errorf("Expected %d matches, got %d: %v", tc.expectedMatches, len(matches), matches)
			}
		})
	}
}

func TestGetPKLFileWithFullConversion(t *testing.T) {
	// Test that the conversion function works with actual PKL files
	content, err := GetPKLFileWithFullConversion("Tool.pkl")
	if err != nil {
		t.Fatalf("Failed to get PKL file with conversion: %v", err)
	}

	// Validate that no package URLs remain
	isValid, remaining := ValidateLocalPaths(content)
	if !isValid {
		t.Errorf("PKL file still contains package URLs after conversion: %v", remaining)
	}

	// Ensure the content is not empty
	if len(content) == 0 {
		t.Error("Converted content is empty")
	}
}

func TestValidateAllPKLFiles(t *testing.T) {
	// Test that all PKL files have been properly converted to local paths
	invalidFiles, err := ValidateAllPKLFiles()
	if err != nil {
		t.Fatalf("Failed to validate all PKL files: %v", err)
	}

	// All files should be valid (no package URLs)
	if len(invalidFiles) > 0 {
		t.Errorf("Found PKL files with remaining package URLs:")
		for filename, urls := range invalidFiles {
			t.Errorf("  %s: %v", filename, urls)
		}
	}
}

func TestConvertAllPKLFiles(t *testing.T) {
	// Test that all PKL files can be converted without errors
	convertedFiles, err := ConvertAllPKLFiles()
	if err != nil {
		t.Fatalf("Failed to convert all PKL files: %v", err)
	}

	// Should have converted some files
	if len(convertedFiles) == 0 {
		t.Error("No files were converted")
	}

	// Each converted file should be valid
	for filename, content := range convertedFiles {
		if len(content) == 0 {
			t.Errorf("Converted file %s is empty", filename)
		}

		isValid, remaining := ValidateLocalPaths(content)
		if !isValid {
			t.Errorf("Converted file %s still contains package URLs: %v", filename, remaining)
		}
	}
}

func TestEnsureOfflineCompatibility(t *testing.T) {
	// Test the high-level offline compatibility check
	err := EnsureOfflineCompatibility()
	if err != nil {
		t.Errorf("Offline compatibility check failed: %v", err)
	}
}

func TestCopyAssetsToTempDir(t *testing.T) {
	// Test copying assets to a temporary directory
	tempDir, err := CopyAssetsToTempDir()
	if err != nil {
		t.Fatalf("Failed to copy assets to temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Verify the temp directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("Temp directory was not created: %s", tempDir)
	}

	t.Logf("Assets copied to temp directory: %s", tempDir)

	// Verify some expected files exist
	expectedFiles := []string{
		"Tool.pkl",
		"Resource.pkl",
		"Workflow.pkl",
		"LLM.pkl",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(tempDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist in temp dir", filename)
		} else {
			t.Logf("✅ Found file: %s", filename)
		}
	}

	// Verify external directory structure exists
	externalDir := filepath.Join(tempDir, "external")
	if _, err := os.Stat(externalDir); os.IsNotExist(err) {
		t.Errorf("External directory does not exist in temp dir")
	} else {
		t.Logf("✅ External directory exists")
	}

	// Verify we can read a file from the temp directory
	toolPath := filepath.Join(tempDir, "Tool.pkl")
	data, err := os.ReadFile(toolPath)
	if err != nil {
		t.Fatalf("Failed to read Tool.pkl from temp dir: %v", err)
	}

	if len(data) == 0 {
		t.Error("Tool.pkl in temp dir is empty")
	}

	content := string(data)
	if !strings.Contains(content, "module") {
		t.Error("Tool.pkl in temp dir doesn't contain expected PKL content")
	}

	t.Logf("✅ Successfully copied and verified assets in temp directory")
}

func TestCopyAssetsToTempDirWithConversion(t *testing.T) {
	// Test copying assets with conversion
	tempDir, err := CopyAssetsToTempDirWithConversion()
	if err != nil {
		t.Fatalf("Failed to copy assets with conversion to temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Verify the temp directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("Temp directory was not created: %s", tempDir)
	}

	t.Logf("Assets copied with conversion to temp directory: %s", tempDir)

	// Read some PKL files and verify they don't contain package URLs
	testFiles := []string{"Tool.pkl", "Resource.pkl", "Workflow.pkl"}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read %s from temp dir: %v", filename, err)
			continue
		}

		content := string(data)
		isValid, remaining := ValidateLocalPaths(content)
		if !isValid {
			t.Errorf("File %s in temp dir still contains package URLs after conversion: %v", filename, remaining)
		} else {
			t.Logf("✅ File %s has no package URLs", filename)
		}
	}

	t.Logf("✅ Successfully copied and converted assets in temp directory")
}

func TestTempDirCleanup(t *testing.T) {
	// Test that we can create and clean up temp directories
	tempDir1, err := CopyAssetsToTempDir()
	if err != nil {
		t.Fatalf("Failed to create first temp dir: %v", err)
	}

	tempDir2, err := CopyAssetsToTempDir()
	if err != nil {
		t.Fatalf("Failed to create second temp dir: %v", err)
	}

	// Verify both directories exist and are different
	if tempDir1 == tempDir2 {
		t.Error("Two calls to CopyAssetsToTempDir returned the same directory")
	}

	if _, err := os.Stat(tempDir1); os.IsNotExist(err) {
		t.Error("First temp directory does not exist")
	}

	if _, err := os.Stat(tempDir2); os.IsNotExist(err) {
		t.Error("Second temp directory does not exist")
	}

	// Clean up first directory
	err = os.RemoveAll(tempDir1)
	if err != nil {
		t.Errorf("Failed to remove first temp dir: %v", err)
	}

	// Verify it's gone
	if _, err := os.Stat(tempDir1); !os.IsNotExist(err) {
		t.Error("First temp directory still exists after cleanup")
	}

	// Clean up second directory
	err = os.RemoveAll(tempDir2)
	if err != nil {
		t.Errorf("Failed to remove second temp dir: %v", err)
	}

	// Verify it's gone
	if _, err := os.Stat(tempDir2); !os.IsNotExist(err) {
		t.Error("Second temp directory still exists after cleanup")
	}

	t.Logf("✅ Temp directory cleanup test passed")
}

func TestWriteAssetsToDir(t *testing.T) {
	// Create a temp directory for testing
	testDir, err := os.MkdirTemp("", "test-write-assets-*")
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Write assets to the test directory
	err = WriteAssetsToDir(testDir)
	if err != nil {
		t.Fatalf("Failed to write assets to directory: %v", err)
	}

	t.Logf("Assets written to: %s", testDir)

	// Verify expected files exist
	expectedFiles := []string{
		"Tool.pkl",
		"Resource.pkl",
		"Workflow.pkl",
		"LLM.pkl",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(testDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", filename)
		} else {
			t.Logf("✅ Found file: %s", filename)
		}
	}

	// Verify external directory exists
	externalDir := filepath.Join(testDir, "external")
	if _, err := os.Stat(externalDir); os.IsNotExist(err) {
		t.Error("External directory does not exist")
	} else {
		t.Logf("✅ External directory exists")
	}

	t.Logf("✅ Successfully wrote assets to directory")
}

func TestWriteAssetsToDirWithConversion(t *testing.T) {
	// Create a temp directory for testing
	testDir, err := os.MkdirTemp("", "test-write-assets-conversion-*")
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Write assets with conversion
	err = WriteAssetsToDirWithConversion(testDir)
	if err != nil {
		t.Fatalf("Failed to write assets with conversion: %v", err)
	}

	t.Logf("Assets written with conversion to: %s", testDir)

	// Verify some PKL files don't contain package URLs
	testFiles := []string{"Tool.pkl", "Resource.pkl", "Workflow.pkl"}

	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read %s: %v", filename, err)
			continue
		}

		content := string(data)
		isValid, remaining := ValidateLocalPaths(content)
		if !isValid {
			t.Errorf("File %s still contains package URLs after conversion: %v", filename, remaining)
		} else {
			t.Logf("✅ File %s has no package URLs", filename)
		}
	}

	t.Logf("✅ Successfully wrote and validated converted assets")
}

func TestWriteAssetsToDirNonExistent(t *testing.T) {
	// Test writing to a non-existent directory (should be created)
	testDir := filepath.Join(os.TempDir(), "test-nested", "dir", "path")
	defer os.RemoveAll(filepath.Join(os.TempDir(), "test-nested"))

	err := WriteAssetsToDir(testDir)
	if err != nil {
		t.Fatalf("Failed to write assets to non-existent directory: %v", err)
	}

	// Verify the directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	} else {
		t.Logf("✅ Non-existent directory was created: %s", testDir)
	}

	// Verify files exist
	toolPath := filepath.Join(testDir, "Tool.pkl")
	if _, err := os.Stat(toolPath); os.IsNotExist(err) {
		t.Error("Tool.pkl was not written to the directory")
	} else {
		t.Logf("✅ Files written to created directory")
	}
}

func TestGetPKLFileFromTempDir(t *testing.T) {
	// Test the integrated helper function
	content, tempDir, cleanup, err := GetPKLFileFromTempDir("Tool.pkl")
	if err != nil {
		t.Fatalf("Failed to get PKL file from temp dir: %v", err)
	}
	defer cleanup()

	t.Logf("Got PKL file from temp dir: %s", tempDir)

	// Verify content is not empty
	if len(content) == 0 {
		t.Error("Content is empty")
	}

	// Verify content is valid PKL
	if !strings.Contains(content, "module") {
		t.Error("Content doesn't contain expected PKL syntax")
	}

	// Verify temp directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Temp directory doesn't exist")
	}

	t.Logf("✅ GetPKLFileFromTempDir works correctly")

	// Test cleanup
	cleanup()
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Temp directory was not cleaned up")
	} else {
		t.Logf("✅ Cleanup function works correctly")
	}
}

func TestGetPKLFileFromTempDirWithConversion(t *testing.T) {
	// Test the integrated helper function with conversion
	content, tempDir, cleanup, err := GetPKLFileFromTempDirWithConversion("Workflow.pkl")
	if err != nil {
		t.Fatalf("Failed to get PKL file from temp dir with conversion: %v", err)
	}
	defer cleanup()

	t.Logf("Got PKL file with conversion from temp dir: %s", tempDir)

	// Verify content is not empty
	if len(content) == 0 {
		t.Error("Content is empty")
	}

	// Verify no package URLs remain
	isValid, remaining := ValidateLocalPaths(content)
	if !isValid {
		t.Errorf("Content still contains package URLs: %v", remaining)
	} else {
		t.Logf("✅ No package URLs in content")
	}

	t.Logf("✅ GetPKLFileFromTempDirWithConversion works correctly")
}

func TestGetPKLFileFromTempDirNonExistent(t *testing.T) {
	// Test with a non-existent file
	_, _, cleanup, err := GetPKLFileFromTempDir("NonExistent.pkl")
	if err == nil {
		defer cleanup()
		t.Error("Expected error for non-existent file, got nil")
	} else {
		t.Logf("✅ Correctly returns error for non-existent file: %v", err)
	}
}

func TestIntegrationScenario(t *testing.T) {
	// Test a realistic integration scenario
	t.Run("scenario_1_write_and_process", func(t *testing.T) {
		// Create a working directory
		workDir, err := os.MkdirTemp("", "integration-test-*")
		if err != nil {
			t.Fatalf("Failed to create work directory: %v", err)
		}
		defer os.RemoveAll(workDir)

		// Write assets to the directory
		err = WriteAssetsToDirWithConversion(workDir)
		if err != nil {
			t.Fatalf("Failed to write assets: %v", err)
		}

		// Process multiple files
		files := []string{"Tool.pkl", "Resource.pkl", "Workflow.pkl"}
		for _, filename := range files {
			filePath := filepath.Join(workDir, filename)
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read %s: %v", filename, err)
				continue
			}

			// Verify each file
			content := string(data)
			if len(content) == 0 {
				t.Errorf("File %s is empty", filename)
			}

			isValid, _ := ValidateLocalPaths(content)
			if !isValid {
				t.Errorf("File %s contains package URLs", filename)
			}
		}

		t.Logf("✅ Integration scenario 1 passed")
	})

	t.Run("scenario_2_helper_function", func(t *testing.T) {
		// Use the helper function for quick access
		content, tempDir, cleanup, err := GetPKLFileFromTempDirWithConversion("LLM.pkl")
		if err != nil {
			t.Fatalf("Failed to get LLM.pkl: %v", err)
		}
		defer cleanup()

		// Verify we can also access other files in the same temp dir
		toolPath := filepath.Join(tempDir, "Tool.pkl")
		toolData, err := os.ReadFile(toolPath)
		if err != nil {
			t.Errorf("Failed to read Tool.pkl from same temp dir: %v", err)
		}

		if len(content) == 0 || len(toolData) == 0 {
			t.Error("Contents are empty")
		}

		t.Logf("✅ Integration scenario 2 passed")
	})

	t.Run("scenario_3_multiple_temp_dirs", func(t *testing.T) {
		// Test using multiple temp directories concurrently
		tempDir1, err := CopyAssetsToTempDir()
		if err != nil {
			t.Fatalf("Failed to create first temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir1)

		tempDir2, err := CopyAssetsToTempDirWithConversion()
		if err != nil {
			t.Fatalf("Failed to create second temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir2)

		// Verify both are different and both exist
		if tempDir1 == tempDir2 {
			t.Error("Temp directories should be different")
		}

		// Read from both
		file1 := filepath.Join(tempDir1, "Tool.pkl")
		file2 := filepath.Join(tempDir2, "Tool.pkl")

		data1, _ := os.ReadFile(file1)
		data2, _ := os.ReadFile(file2)

		if len(data1) == 0 || len(data2) == 0 {
			t.Error("Files are empty")
		}

		t.Logf("✅ Integration scenario 3 passed")
	})
}

func TestRedundantConversionBehavior(t *testing.T) {
	t.Run("simulate_missed_package_urls", func(t *testing.T) {
		// Simulate content that might have been missed by shell script conversion
		testCases := []struct {
			name        string
			input       string
			description string
		}{
			{
				name: "missed_pkl_go_import",
				input: `/// Tool module with missed conversion
@ModuleInfo { minPklVersion = "0.29.1" }
module Tool

import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"
import "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"

function test() = "example"`,
				description: "PKL file with one missed pkl-go import",
			},
			{
				name: "missed_pkl_pantry_import",
				input: `/// Resource module with missed conversion
@ModuleInfo { minPklVersion = "0.29.1" }
module Resource

import "external/pkl-go/codegen/src/go.pkl"
import "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"

class ResourceAction {}`,
				description: "PKL file with one missed pkl-pantry import",
			},
			{
				name: "multiple_missed_imports",
				input: `/// Module with multiple missed conversions
@ModuleInfo { minPklVersion = "0.29.1" }
module Test

import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"
import "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"
import "package://pkg.pkl-lang.org/pkl-pantry/org.openapis.v3@2.1.2#/Schema.pkl"

function multiTest() = "example"`,
				description: "PKL file with multiple missed imports",
			},
			{
				name: "missed_amends_statement",
				input: `/// Project with missed amends conversion
amends "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"

local version = "test"
package {
    name = "test-package"
    version = version
}`,
				description: "PKL file with missed amends statement",
			},
			{
				name: "mixed_converted_and_missed",
				input: `/// Mixed scenario - some converted, some missed
@ModuleInfo { minPklVersion = "0.29.1" }
module Mixed

// These were already converted correctly
import "external/pkl-go/codegen/src/go.pkl"
import "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"

// These were missed by the shell script
import "package://pkg.pkl-lang.org/pkl-pantry/org.openapis.v3@2.1.2#/Schema.pkl"
amends "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/Generator.pkl"

class MixedClass {}`,
				description: "PKL file with both converted and missed URLs",
			},
			{
				name: "schema_core_workflow_error",
				input: `/// This simulates the problematic test case
amends "package://schema.kdeps.com/core@1.0.0#/Workflow.pkl"

local testData = "example"
workflow {
    targetActionID = "test"
}`,
				description: "PKL file with schema.kdeps.com core URL (the problematic case)",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Logf("Testing: %s", tc.description)

				// 1. Verify the input contains package URLs (simulating missed conversion)
				isValidBefore, packageURLsBefore := ValidateLocalPaths(tc.input)
				if isValidBefore {
					t.Fatalf("Test case setup error: input should contain package URLs")
				}
				t.Logf("Found %d package URLs before conversion: %v", len(packageURLsBefore), packageURLsBefore)

				// 2. Apply redundant conversion
				converted := ConvertPackageURLsToLocalPaths(tc.input)
				t.Logf("Applied ConvertPackageURLsToLocalPaths")

				// 3. Apply import statement conversion for extra safety
				fullyConverted := ConvertImportStatements(converted)
				t.Logf("Applied ConvertImportStatements")

				// 4. Validate that all package URLs are now gone
				isValidAfter, remainingURLs := ValidateLocalPaths(fullyConverted)
				if !isValidAfter {
					t.Errorf("Conversion failed - remaining package URLs: %v", remainingURLs)
				} else {
					t.Logf("✅ All package URLs successfully converted")
				}

				// 5. Verify the converted content is still valid PKL
				if len(fullyConverted) == 0 {
					t.Error("Converted content is empty")
				}

				// 6. Check that local paths are properly formatted
				// For schema.kdeps.com URLs, we expect direct local paths (not external/)
				// For pkg.pkl-lang.org URLs, we expect external/ paths
				containsExternal := strings.Contains(fullyConverted, "external/")
				containsSchema := strings.Contains(tc.input, "schema.kdeps.com")
				
				if !containsSchema && !containsExternal {
					t.Error("Converted content should contain local external/ paths for pkg.pkl-lang.org URLs")
				} else if containsSchema && containsExternal {
					t.Error("Schema.kdeps.com URLs should be converted to direct local paths, not external/ paths")
				}

				// 7. Ensure no double conversion issues
				secondConversion := ConvertPackageURLsToLocalPaths(fullyConverted)
				if secondConversion != fullyConverted {
					t.Error("Second conversion changed already converted content - potential double conversion issue")
				}

				t.Logf("✅ Redundant conversion test passed for: %s", tc.name)
			})
		}
	})

	t.Run("conversion_robustness", func(t *testing.T) {
		// Test edge cases and robustness
		edgeCases := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "nested_quotes",
				input:    `description = "Uses package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"`,
				expected: `description = "Uses external/pkl-go/codegen/src/go.pkl"`,
			},
			{
				name:     "different_versions",
				input:    `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.10.0#/go.pkl"`,
				expected: `import "external/pkl-go/codegen/src/go.pkl"`,
			},
			{
				name:     "codegen_path_variation",
				input:    `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/codegen/src/go.pkl"`,
				expected: `import "external/pkl-go/codegen/src/go.pkl"`,
			},
			{
				name:     "multiple_same_line",
				input:    `// Both package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl and package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl`,
				expected: `// Both external/pkl-go/codegen/src/go.pkl and external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl`,
			},
			{
				name:     "schema_kdeps_core_url",
				input:    `amends "package://schema.kdeps.com/core@1.0.0#/Workflow.pkl"`,
				expected: `amends "Workflow.pkl"`,
			},
			{
				name:     "schema_kdeps_core_different_version",
				input:    `import "package://schema.kdeps.com/core@0.2.42#/Resource.pkl"`,
				expected: `import "Resource.pkl"`,
			},
			{
				name:     "mixed_schema_and_pkl_urls",
				input:    `import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"\namends "package://schema.kdeps.com/core@1.0.0#/Workflow.pkl"`,
				expected: `import "external/pkl-go/codegen/src/go.pkl"\namends "Workflow.pkl"`,
			},
		}

		for _, tc := range edgeCases {
			t.Run(tc.name, func(t *testing.T) {
				result := ConvertPackageURLsToLocalPaths(tc.input)
				if result != tc.expected {
					t.Errorf("Expected: %s\nGot: %s", tc.expected, result)
				}
			})
		}
	})

	t.Run("batch_validation_and_conversion", func(t *testing.T) {
		// Test that the batch functions would catch any missed conversions
		t.Log("Testing batch validation to catch any missed conversions")

		// This should pass since our files are already converted
		invalidFiles, err := ValidateAllPKLFiles()
		if err != nil {
			t.Fatalf("Batch validation failed: %v", err)
		}

		if len(invalidFiles) > 0 {
			t.Errorf("Found files with unconverted package URLs:")
			for filename, urls := range invalidFiles {
				t.Errorf("  %s: %v", filename, urls)
			}
		} else {
			t.Log("✅ All PKL files pass batch validation")
		}

		// Test batch conversion capability
		convertedFiles, err := ConvertAllPKLFiles()
		if err != nil {
			t.Fatalf("Batch conversion failed: %v", err)
		}

		if len(convertedFiles) == 0 {
			t.Error("Batch conversion returned no files")
		}

		t.Logf("✅ Successfully batch converted %d files", len(convertedFiles))
	})

	t.Run("integration_with_actual_files", func(t *testing.T) {
		// Test the conversion functions with actual PKL files
		testFiles := []string{"Tool.pkl", "Resource.pkl", "Workflow.pkl"}

		for _, filename := range testFiles {
			t.Run(filename, func(t *testing.T) {
				// Get original content
				_, err := GetPKLFileAsString(filename)
				if err != nil {
					t.Fatalf("Failed to read %s: %v", filename, err)
				}

				// Apply full conversion (should be idempotent since files are already converted)
				convertedContent, err := GetPKLFileWithFullConversion(filename)
				if err != nil {
					t.Fatalf("Failed to convert %s: %v", filename, err)
				}

				// Validate no package URLs remain
				isValid, remaining := ValidateLocalPaths(convertedContent)
				if !isValid {
					t.Errorf("File %s still contains package URLs after conversion: %v", filename, remaining)
				}

				// Since files should already be converted, content should be the same or very similar
				if len(convertedContent) == 0 {
					t.Errorf("Converted content for %s is empty", filename)
				}

				t.Logf("✅ File %s passes integration test", filename)
			})
		}
	})

	t.Run("error_resilience", func(t *testing.T) {
		// Test that functions handle edge cases gracefully
		testCases := []struct {
			name  string
			input string
		}{
			{"empty_string", ""},
			{"no_imports", "module Test\nclass Example {}"},
			{"already_local_paths", `import "external/pkl-go/codegen/src/go.pkl"`},
			{"malformed_urls", `import "package://invalid-url"`},
			{"mixed_content", "Some text\npackage://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl\nMore text"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// These should not panic or error
				result1 := ConvertPackageURLsToLocalPaths(tc.input)
				result2 := ConvertImportStatements(tc.input)
				isValid, _ := ValidateLocalPaths(tc.input)

				t.Logf("Input: %q", tc.input)
				t.Logf("ConvertPackageURLsToLocalPaths: %q", result1)
				t.Logf("ConvertImportStatements: %q", result2)
				t.Logf("ValidateLocalPaths: %v", isValid)
			})
		}
	})
}
