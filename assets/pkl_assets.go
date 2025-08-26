// Package assets makes the PKL schema available to downstream code/tests.
package assets

import (
	"embed"
	"fmt"
)

//go:embed pkl external
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
	path = fmt.Sprintf("external/%s", filename)
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

// GetPKLFile reads a specific PKL file from the pkl directory
func GetPKLFileFromPKL(filename string) ([]byte, error) {
	path := fmt.Sprintf("pkl/%s", filename)
	return PKLFS.ReadFile(path)
}

// GetExternalFile reads a specific file from the external directory
func GetExternalFile(filename string) ([]byte, error) {
	path := fmt.Sprintf("external/%s", filename)
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
	err := walkDir("external", &files)
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
