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
@ModuleInfo { minPklVersion = "0.29.0" }
module Tool

import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.1#/go.pkl"
import "external/pkl-pantry/packages/pkl.experimental.uri/URI.pkl"

function test() = "example"`,
				description: "PKL file with one missed pkl-go import",
			},
			{
				name: "missed_pkl_pantry_import",
				input: `/// Resource module with missed conversion
@ModuleInfo { minPklVersion = "0.29.0" }
module Resource

import "external/pkl-go/codegen/src/go.pkl"
import "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"

class ResourceAction {}`,
				description: "PKL file with one missed pkl-pantry import",
			},
			{
				name: "multiple_missed_imports",
				input: `/// Module with multiple missed conversions
@ModuleInfo { minPklVersion = "0.29.0" }
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
@ModuleInfo { minPklVersion = "0.29.0" }
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
				if !strings.Contains(fullyConverted, "external/") {
					t.Error("Converted content should contain local external/ paths")
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
