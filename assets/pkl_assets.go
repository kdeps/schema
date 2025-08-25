// Package assets makes the PKL schema available to downstream code/tests.
package assets

import (
	"embed"
	"fmt"
)

//go:embed pkl/*.pkl
var PKLFS embed.FS

// GetPKLFile reads a specific PKL file from the embedded filesystem
func GetPKLFile(filename string) ([]byte, error) {
	path := fmt.Sprintf("pkl/%s", filename)
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
