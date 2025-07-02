// Package assets makes the PKL schema available to downstream code/tests.
//
// This package provides embedded PKL schema files and utilities for working with them
// in Go applications and tests. All PKL files from deps/pkl/ are embedded at build time
// and can be accessed without external file dependencies.
//
// # Test Usage Examples
//
// ## Basic PKL Workspace Setup for Tests
//
//	func TestPKLWorkflow(t *testing.T) {
//	    // Setup PKL workspace with all schema files
//	    workspace, err := assets.SetupPKLWorkspaceInTmpDir()
//	    if err != nil {
//	        t.Fatalf("Failed to setup PKL workspace: %v", err)
//	    }
//	    defer workspace.Cleanup() // Important: clean up temp files
//
//	    // Get absolute paths for PKL imports
//	    workflowPath := workspace.GetImportPath("Workflow.pkl")
//	    resourcePath := workspace.GetImportPath("Resource.pkl")
//
//	    // Use these paths in your PKL files or evaluation code
//	    t.Logf("Workflow PKL available at: %s", workflowPath)
//	}
//
// ## PKL File with Absolute Imports
//
// Create a test PKL file that imports the schema:
//
//	// test_workflow.pkl
//	import "/absolute/path/from/workspace/Workflow.pkl" as Workflow
//	import "/absolute/path/from/workspace/Resource.pkl" as Resource
//
//	// All cross-references work because files are in same directory
//	myWorkflow = new Workflow {
//	    AgentID = "test-agent"
//	    Description = "Test workflow"
//	    Version = "1.0.0"
//	    TargetActionID = "test-action"
//	    Workflows {}
//	    Settings = new Workflow.Project.Settings {
//	        AgentSettings = new Workflow.Docker.DockerSettings {}
//	        APISettings = new Workflow.APIServer.APIServerSettings {}
//	    }
//	}
//
// ## Advanced Test Example with PKL Evaluation
//
//	func TestPKLEvaluation(t *testing.T) {
//	    workspace, err := assets.SetupPKLWorkspaceInTmpDir()
//	    if err != nil {
//	        t.Fatalf("Setup failed: %v", err)
//	    }
//	    defer workspace.Cleanup()
//
//	    // Create test PKL content with imports
//	    testContent := fmt.Sprintf(`
//	        import "%s" as Workflow
//	        import "%s" as LLM
//
//	        testWorkflow = new Workflow {
//	            AgentID = "test"
//	            Description = "Test"
//	            Version = "1.0.0"
//	            TargetActionID = "test-action"
//	            Workflows {}
//	            Settings = new Workflow.Project.Settings {
//	                AgentSettings = new Workflow.Docker.DockerSettings {}
//	                APISettings = new Workflow.APIServer.APIServerSettings {}
//	            }
//	        }
//
//	        testLLM = (LLM.resource("test-llm")).Response
//	    `, workspace.GetImportPath("Workflow.pkl"), workspace.GetImportPath("LLM.pkl"))
//
//	    // Write test file
//	    testFile := filepath.Join(workspace.Directory, "test.pkl")
//	    err = os.WriteFile(testFile, []byte(testContent), 0644)
//	    if err != nil {
//	        t.Fatalf("Failed to write test file: %v", err)
//	    }
//
//	    // Now evaluate with pkl command or library
//	    // pkl eval test.pkl
//	}
//
// ## Available PKL Schema Files
//
// All files from deps/pkl/ are available:
//   - Workflow.pkl - Workflow definitions and validation
//   - Resource.pkl - Resource actions and configurations
//   - LLM.pkl - Language model configurations
//   - APIServer.pkl - API server settings
//   - Docker.pkl - Docker container settings
//   - And 16 more schema files...
//
// ## LLM Integration Notes
//
// When using with LLM tools or code generation:
package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed pkl/*.pkl
var PKLFS embed.FS

// GetPKLFile reads a specific PKL file from the embedded filesystem.
// The filename should not include the path (e.g., "Workflow.pkl", not "pkl/Workflow.pkl")
func GetPKLFile(filename string) ([]byte, error) {
	return PKLFS.ReadFile(filepath.Join("pkl", filename))
}

// GetPKLFileAsString reads a specific PKL file as a string from the embedded filesystem.
func GetPKLFileAsString(filename string) (string, error) {
	data, err := GetPKLFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetSourcePKLPath returns the absolute path to the original PKL file in the source directory.
// This is useful during development but won't work in deployed binaries.
func GetSourcePKLPath(filename string) (string, error) {
	// Get the directory where this Go file is located
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// Navigate from assets/ to deps/pkl/
	assetDir := filepath.Dir(currentFile)
	repoRoot := filepath.Dir(assetDir)
	pklPath := filepath.Join(repoRoot, "deps", "pkl", filename)

	// Check if the file exists
	if _, err := os.Stat(pklPath); err != nil {
		return "", fmt.Errorf("PKL file not found at %s: %w", pklPath, err)
	}

	return filepath.Abs(pklPath)
}

// ExtractPKLFileToTemp extracts an embedded PKL file to a temporary location and returns the absolute path.
// The caller is responsible for cleaning up the temporary file.
func ExtractPKLFileToTemp(filename string) (string, error) {
	data, err := GetPKLFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded file %s: %w", filename, err)
	}

	// Create a temporary file in system temp directory
	tmpDir := os.TempDir()
	tempFile, err := os.CreateTemp(tmpDir, fmt.Sprintf("pkl_%s_*.pkl", strings.TrimSuffix(filename, ".pkl")))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file in %s: %w", tmpDir, err)
	}
	defer tempFile.Close()

	// Write the embedded content to the temp file
	if _, err := tempFile.Write(data); err != nil {
		os.Remove(tempFile.Name()) // Clean up on error
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

// ExtractPKLFileWithName extracts an embedded PKL file to a directory preserving the original filename.
// If dir is empty, uses system tmpdir. Returns the absolute path of the extracted file.
func ExtractPKLFileWithName(filename, dir string) (string, error) {
	data, err := GetPKLFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded file %s: %w", filename, err)
	}

	if dir == "" {
		dir = os.TempDir()
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Preserve original filename
	outputPath := filepath.Join(dir, filename)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", outputPath, err)
	}

	return outputPath, nil
}

// ExtractAllPKLFilesToDir extracts all embedded PKL files to the specified directory.
// If dir is empty, uses a temporary directory in system tmpdir. Returns the directory path.
func ExtractAllPKLFilesToDir(dir string) (string, error) {
	if dir == "" {
		tmpDir := os.TempDir()
		var err error
		dir, err = os.MkdirTemp(tmpDir, "pkl_extracted_*")
		if err != nil {
			return "", fmt.Errorf("failed to create temp directory in %s: %w", tmpDir, err)
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	files, err := ListPKLFiles()
	if err != nil {
		return "", fmt.Errorf("failed to list PKL files: %w", err)
	}

	for _, filename := range files {
		data, err := GetPKLFile(filename)
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %w", filename, err)
		}

		outputPath := filepath.Join(dir, filename)
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			return "", fmt.Errorf("failed to write %s: %w", outputPath, err)
		}
	}

	return dir, nil
}

// SetupPKLWorkspace extracts all PKL files to a directory and returns a PKLWorkspace
// for easy access to absolute paths suitable for PKL imports.
// If dir is empty, creates a temporary directory in system tmpdir.
func SetupPKLWorkspace(dir string) (*PKLWorkspace, error) {
	extractedDir, err := ExtractAllPKLFilesToDir(dir)
	if err != nil {
		return nil, err
	}

	return &PKLWorkspace{
		Directory: extractedDir,
		isTemp:    dir == "", // Track if we created a temp directory
	}, nil
}

// SetupPKLWorkspaceInTmpDir creates a PKL workspace specifically in the system temporary directory.
// This is a convenience function that explicitly uses tmpdir.
func SetupPKLWorkspaceInTmpDir() (*PKLWorkspace, error) {
	tmpDir := os.TempDir()
	workspaceDir, err := os.MkdirTemp(tmpDir, "pkl_workspace_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace in tmpdir %s: %w", tmpDir, err)
	}

	extractedDir, err := ExtractAllPKLFilesToDir(workspaceDir)
	if err != nil {
		os.RemoveAll(workspaceDir) // Clean up on error
		return nil, err
	}

	return &PKLWorkspace{
		Directory: extractedDir,
		isTemp:    true,
	}, nil
}

// PKLWorkspace represents an extracted PKL workspace with all schema files.
type PKLWorkspace struct {
	Directory string
	isTemp    bool
}

// GetTmpDir returns the system temporary directory path.
// This respects TMPDIR environment variable on Unix systems.
func GetTmpDir() string {
	return os.TempDir()
}

// IsInTmpDir checks if the workspace directory is located in the system temporary directory.
func (w *PKLWorkspace) IsInTmpDir() bool {
	tmpDir := os.TempDir()
	return strings.HasPrefix(w.Directory, tmpDir)
}

// IsTemporary returns true if this workspace was created as a temporary directory
// and will be cleaned up when Cleanup() is called.
func (w *PKLWorkspace) IsTemporary() bool {
	return w.isTemp
}

// GetAbsolutePath returns the absolute path for a PKL file suitable for imports.
// Example: workspace.GetAbsolutePath("Workflow.pkl") -> "/tmp/pkl_123/Workflow.pkl"
func (w *PKLWorkspace) GetAbsolutePath(filename string) string {
	return filepath.Join(w.Directory, filename)
}

// GetImportPath returns a path suitable for PKL import statements.
// Example: workspace.GetImportPath("Workflow.pkl") -> "/tmp/pkl_123/Workflow.pkl"
func (w *PKLWorkspace) GetImportPath(filename string) string {
	return w.GetAbsolutePath(filename)
}

// ListFiles returns all PKL files available in the workspace.
func (w *PKLWorkspace) ListFiles() ([]string, error) {
	entries, err := os.ReadDir(w.Directory)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pkl") {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// Cleanup removes the workspace directory if it was created as a temporary directory.
// Call this in defer or when done with the workspace.
func (w *PKLWorkspace) Cleanup() error {
	if w.isTemp {
		return os.RemoveAll(w.Directory)
	}
	return nil
}

// GetPKLFilePath returns a path-like identifier for an embedded PKL file.
// This is not a real filesystem path but can be used for identification/logging.
func GetPKLFilePath(filename string) string {
	return fmt.Sprintf("embed://pkl/%s", filename)
}

// ListPKLFiles returns a list of all embedded PKL files.
func ListPKLFiles() ([]string, error) {
	var files []string

	err := fs.WalkDir(PKLFS, "pkl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".pkl") {
			files = append(files, d.Name())
		}

		return nil
	})

	return files, err
}

// ValidatePKLFiles checks that all expected PKL files are present in the embedded filesystem.
func ValidatePKLFiles() error {
	expectedFiles := []string{
		"APIServer.pkl",
		"APIServerRequest.pkl",
		"APIServerResponse.pkl",
		"Data.pkl",
		"Docker.pkl",
		"Document.pkl",
		"Exec.pkl",
		"HTTP.pkl",
		"Item.pkl",
		"Kdeps.pkl",
		"LLM.pkl",
		"Memory.pkl",
		"Project.pkl",
		"Python.pkl",
		"Resource.pkl",
		"Session.pkl",
		"Skip.pkl",
		"Tool.pkl",
		"Utils.pkl",
		"WebServer.pkl",
		"Workflow.pkl",
	}

	availableFiles, err := ListPKLFiles()
	if err != nil {
		return fmt.Errorf("failed to list PKL files: %w", err)
	}

	// Convert to map for faster lookup
	available := make(map[string]bool)
	for _, file := range availableFiles {
		available[file] = true
	}

	// Check for missing files
	var missing []string
	for _, expected := range expectedFiles {
		if !available[expected] {
			missing = append(missing, expected)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing PKL files: %v", missing)
	}

	return nil
}

// GetFileSystem returns the embedded filesystem for advanced usage.
func GetFileSystem() embed.FS {
	return PKLFS
}

// PKLFileExists checks if a specific PKL file exists in the embedded filesystem.
func PKLFileExists(filename string) bool {
	_, err := GetPKLFile(filename)
	return err == nil
}
