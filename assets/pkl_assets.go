// Package assets makes the PKL schema available to downstream code/tests.
package assets

import (
	"embed"
	"fmt"
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

	return content
}

// ConvertImportStatements converts import and amends statements from package URLs to local paths
func ConvertImportStatements(content string) string {
	// Convert import statements
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

	// Convert amends statements
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
	packageURLPattern := regexp.MustCompile(`package://pkg\.pkl-lang\.org/[^"]+`)
	matches := packageURLPattern.FindAllString(content, -1)
	return len(matches) == 0, matches
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
