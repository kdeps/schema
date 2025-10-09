// Package assets makes the PKL schema available to downstream code/tests.
package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed pkl
var PKLFS embed.FS

// GetPKLFile reads a specific PKL file from the embedded filesystem
func GetPKLFile(filename string) ([]byte, error) {
	// Try pkl directory first
	path := fmt.Sprintf("pkl/%s", filename)
	data, err := PKLFS.ReadFile(path)
	if err == nil {
		return data, nil
	}

	// Try external directory
	path = fmt.Sprintf("pkl/external/%s", filename)
	return PKLFS.ReadFile(path)
}

// GetPKLFileAsString reads a specific PKL file from the embedded filesystem and returns it as a string
func GetPKLFileAsString(filename string) (string, error) {
	data, err := GetPKLFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetPKLFileAsStringWithLocalPaths reads a PKL file and converts any remaining package URLs to local paths
func GetPKLFileAsStringWithLocalPaths(filename string) (string, error) {
	content, err := GetPKLFileAsString(filename)
	if err != nil {
		return "", err
	}
	return ConvertPackageURLsToLocalPaths(content), nil
}

// GetPKLFile reads a specific PKL file from the pkl directory
func GetPKLFileFromPKL(filename string) ([]byte, error) {
	path := fmt.Sprintf("pkl/%s", filename)
	return PKLFS.ReadFile(path)
}

// GetExternalFile reads a specific file from the external directory
func GetExternalFile(filename string) ([]byte, error) {
	path := fmt.Sprintf("pkl/external/%s", filename)
	return PKLFS.ReadFile(path)
}

// ListPKLFiles returns all PKL files in the pkl directory
func ListPKLFiles() ([]string, error) {
	var files []string
	entries, err := PKLFS.ReadDir("pkl")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 4 && entry.Name()[len(entry.Name())-4:] == ".pkl" {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// ListExternalFiles returns all files in the external directory (recursive)
func ListExternalFiles() ([]string, error) {
	var files []string
	err := walkDir("pkl/external", &files)
	return files, err
}

// walkDir recursively walks through the embedded filesystem
func walkDir(dir string, files *[]string) error {
	entries, err := PKLFS.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := fmt.Sprintf("%s/%s", dir, entry.Name())
		if entry.IsDir() {
			err := walkDir(path, files)
			if err != nil {
				return err
			}
		} else {
			*files = append(*files, path)
		}
	}
	return nil
}

// ConvertPackageURLsToLocalPaths converts package:// URLs to local relative paths
func ConvertPackageURLsToLocalPaths(content string) string {
	// Convert pkl-go package URLs
	pklGoPattern := regexp.MustCompile(`package://pkg\.pkl-lang\.org/pkl-go/pkl\.golang@[^"#]*#/(?:codegen/src/)?([^"]+)`)
	content = pklGoPattern.ReplaceAllString(content, `external/pkl-go/codegen/src/$1`)

	// Convert pkl-pantry package URLs
	pklPantryPattern := regexp.MustCompile(`package://pkg\.pkl-lang\.org/pkl-pantry/([^@"]+)@[^"#]*#/([^"]+)`)
	content = pklPantryPattern.ReplaceAllString(content, `external/pkl-pantry/packages/$1/$2`)

	// Convert schema.kdeps.com core package URLs to local paths
	// This handles URLs like package://schema.kdeps.com/core@1.0.0#/Workflow.pkl
	schemaCorePattern := regexp.MustCompile(`package://schema\.kdeps\.com/core@[^"#]*#/([^"]+)`)
	content = schemaCorePattern.ReplaceAllString(content, `$1`)

	return content
}

// ConvertImportStatements converts import and amends statements from package URLs to local paths
func ConvertImportStatements(content string) string {
	// Convert import statements for pkg.pkl-lang.org
	importPattern := regexp.MustCompile(`(import\s+)"(package://pkg\.pkl-lang\.org/[^"]+)"`)
	content = importPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := importPattern.FindStringSubmatch(match)
		if len(parts) == 3 {
			localPath := ConvertPackageURLsToLocalPaths(`"` + parts[2] + `"`)
			localPath = strings.Trim(localPath, `"`)
			return parts[1] + `"` + localPath + `"`
		}
		return match
	})

	// Convert amends statements for pkg.pkl-lang.org
	amendsPattern := regexp.MustCompile(`(amends\s+)"(package://pkg\.pkl-lang\.org/[^"]+)"`)
	content = amendsPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := amendsPattern.FindStringSubmatch(match)
		if len(parts) == 3 {
			localPath := ConvertPackageURLsToLocalPaths(`"` + parts[2] + `"`)
			localPath = strings.Trim(localPath, `"`)
			return parts[1] + `"` + localPath + `"`
		}
		return match
	})

	// Convert import statements for schema.kdeps.com
	schemaImportPattern := regexp.MustCompile(`(import\s+)"(package://schema\.kdeps\.com/[^"]+)"`)
	content = schemaImportPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := schemaImportPattern.FindStringSubmatch(match)
		if len(parts) == 3 {
			localPath := ConvertPackageURLsToLocalPaths(`"` + parts[2] + `"`)
			localPath = strings.Trim(localPath, `"`)
			return parts[1] + `"` + localPath + `"`
		}
		return match
	})

	// Convert amends statements for schema.kdeps.com
	schemaAmendsPattern := regexp.MustCompile(`(amends\s+)"(package://schema\.kdeps\.com/[^"]+)"`)
	content = schemaAmendsPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := schemaAmendsPattern.FindStringSubmatch(match)
		if len(parts) == 3 {
			localPath := ConvertPackageURLsToLocalPaths(`"` + parts[2] + `"`)
			localPath = strings.Trim(localPath, `"`)
			return parts[1] + `"` + localPath + `"`
		}
		return match
	})

	return content
}

// GetPKLFileWithFullConversion reads a PKL file and applies all available conversion methods
func GetPKLFileWithFullConversion(filename string) (string, error) {
	content, err := GetPKLFileAsString(filename)
	if err != nil {
		return "", err
	}
	
	// Apply redundant conversions to ensure all package URLs are converted
	content = ConvertPackageURLsToLocalPaths(content)
	content = ConvertImportStatements(content)
	
	return content, nil
}

// ValidateLocalPaths checks if a PKL file content has any remaining package:// URLs
func ValidateLocalPaths(content string) (bool, []string) {
	// Check for pkg.pkl-lang.org URLs
	pklLangPattern := regexp.MustCompile(`package://pkg\.pkl-lang\.org/[^"\s]+`)
	pklLangMatches := pklLangPattern.FindAllString(content, -1)
	
	// Check for schema.kdeps.com URLs  
	schemaPattern := regexp.MustCompile(`package://schema\.kdeps\.com/[^"\s]+`)
	schemaMatches := schemaPattern.FindAllString(content, -1)
	
	// Combine all matches
	allMatches := append(pklLangMatches, schemaMatches...)
	
	return len(allMatches) == 0, allMatches
}

// ValidateAllPKLFiles checks all PKL files for remaining package:// URLs
func ValidateAllPKLFiles() (map[string][]string, error) {
	files, err := ListPKLFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list PKL files: %w", err)
	}

	invalidFiles := make(map[string][]string)
	
	for _, filename := range files {
		content, err := GetPKLFileAsString(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filename, err)
		}
		
		isValid, matches := ValidateLocalPaths(content)
		if !isValid {
			invalidFiles[filename] = matches
		}
	}
	
	return invalidFiles, nil
}

// ConvertAllPKLFiles applies package URL to local path conversion to all PKL files in memory
func ConvertAllPKLFiles() (map[string]string, error) {
	files, err := ListPKLFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list PKL files: %w", err)
	}

	convertedFiles := make(map[string]string)
	
	for _, filename := range files {
		content, err := GetPKLFileWithFullConversion(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s: %w", filename, err)
		}
		convertedFiles[filename] = content
	}
	
	return convertedFiles, nil
}

// EnsureOfflineCompatibility is a high-level function that validates and ensures all PKL files use local paths
// This provides redundant conversion as a safety net for any missed remote package URLs
func EnsureOfflineCompatibility() error {
	// First, validate current state
	invalidFiles, err := ValidateAllPKLFiles()
	if err != nil {
		return fmt.Errorf("failed to validate PKL files: %w", err)
	}

	// If any files still have package URLs, report them
	if len(invalidFiles) > 0 {
		var details strings.Builder
		details.WriteString("Found PKL files with remaining package URLs:\n")
		for filename, urls := range invalidFiles {
			details.WriteString(fmt.Sprintf("  %s: %v\n", filename, urls))
		}
		return fmt.Errorf("offline compatibility check failed:\n%s", details.String())
	}

	return nil
}

// CopyAssetsToTempDir copies all embedded PKL assets to a temporary directory
// and returns the complete path to the temporary directory.
func CopyAssetsToTempDir() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "pkl-assets-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Copy all files from the embedded filesystem to the temp directory
	err = fs.WalkDir(PKLFS, "pkl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path in the temp directory
		relPath, err := filepath.Rel("pkl", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		destPath := filepath.Join(tempDir, relPath)

		// If it's a directory, create it
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read the file from embedded FS
		data, err := PKLFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write the file to the temp directory
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		// Clean up on error
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to copy assets: %w", err)
	}

	return tempDir, nil
}

// CopyAssetsToTempDirWithConversion copies all embedded PKL assets to a temporary directory
// with package URL to local path conversion applied, and returns the complete path.
func CopyAssetsToTempDirWithConversion() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "pkl-assets-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Copy all files from the embedded filesystem to the temp directory
	err = fs.WalkDir(PKLFS, "pkl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path in the temp directory
		relPath, err := filepath.Rel("pkl", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		destPath := filepath.Join(tempDir, relPath)

		// If it's a directory, create it
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read the file from embedded FS
		data, err := PKLFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Apply conversion if it's a .pkl file
		content := string(data)
		if filepath.Ext(path) == ".pkl" {
			content = ConvertPackageURLsToLocalPaths(content)
			content = ConvertImportStatements(content)
		}

		// Write the file to the temp directory
		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		// Clean up on error
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to copy assets with conversion: %w", err)
	}

	return tempDir, nil
}

// WriteAssetsToDir writes all embedded PKL assets to the specified directory.
// If the directory doesn't exist, it will be created.
// Returns the complete path to the directory.
func WriteAssetsToDir(targetDir string) error {
	// Create the target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Copy all files from the embedded filesystem to the target directory
	err := fs.WalkDir(PKLFS, "pkl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path in the target directory
		relPath, err := filepath.Rel("pkl", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		destPath := filepath.Join(targetDir, relPath)

		// If it's a directory, create it
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read the file from embedded FS
		data, err := PKLFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write the file to the target directory
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to write assets to directory: %w", err)
	}

	return nil
}

// WriteAssetsToDirWithConversion writes all embedded PKL assets to the specified directory
// with package URL to local path conversion applied.
// If the directory doesn't exist, it will be created.
func WriteAssetsToDirWithConversion(targetDir string) error {
	// Create the target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Copy all files from the embedded filesystem to the target directory
	err := fs.WalkDir(PKLFS, "pkl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path in the target directory
		relPath, err := filepath.Rel("pkl", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		destPath := filepath.Join(targetDir, relPath)

		// If it's a directory, create it
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read the file from embedded FS
		data, err := PKLFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Apply conversion if it's a .pkl file
		content := string(data)
		if filepath.Ext(path) == ".pkl" {
			content = ConvertPackageURLsToLocalPaths(content)
			content = ConvertImportStatements(content)
		}

		// Write the file to the target directory
		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to write assets to directory with conversion: %w", err)
	}

	return nil
}

// GetPKLFileFromTempDir is a helper function that combines CopyAssetsToTempDir
// and reading a specific file. It returns the file content and a cleanup function.
// The caller should defer the cleanup function to remove the temp directory.
func GetPKLFileFromTempDir(filename string) (content string, tempDir string, cleanup func(), err error) {
	tempDir, err = CopyAssetsToTempDir()
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to copy assets to temp dir: %w", err)
	}

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	filePath := filepath.Join(tempDir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		cleanup()
		return "", "", nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return string(data), tempDir, cleanup, nil
}

// GetPKLFileFromTempDirWithConversion is similar to GetPKLFileFromTempDir but
// applies package URL conversion to all files in the temp directory.
func GetPKLFileFromTempDirWithConversion(filename string) (content string, tempDir string, cleanup func(), err error) {
	tempDir, err = CopyAssetsToTempDirWithConversion()
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to copy assets to temp dir with conversion: %w", err)
	}

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	filePath := filepath.Join(tempDir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		cleanup()
		return "", "", nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return string(data), tempDir, cleanup, nil
}
