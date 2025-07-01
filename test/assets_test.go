package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kdeps/schema/assets"
)

// TestPKLFileEmbedding tests that PKL files are properly embedded
func TestPKLFileEmbedding(t *testing.T) {
	t.Run("PKLFileExists", func(t *testing.T) {
		testFiles := []string{"Workflow.pkl", "Resource.pkl", "LLM.pkl", "APIServer.pkl"}
		for _, filename := range testFiles {
			if !assets.PKLFileExists(filename) {
				t.Errorf("PKL file %s should exist in embedded filesystem", filename)
			}
		}
	})

	t.Run("GetPKLFile", func(t *testing.T) {
		data, err := assets.GetPKLFile("Workflow.pkl")
		if err != nil {
			t.Fatalf("Failed to get Workflow.pkl: %v", err)
		}
		if len(data) == 0 {
			t.Error("Workflow.pkl should not be empty")
		}

		// Check for expected content
		content := string(data)
		if !strings.Contains(content, "Workflow Management") {
			t.Error("Workflow.pkl should contain 'Workflow Management'")
		}
	})

	t.Run("GetPKLFileAsString", func(t *testing.T) {
		content, err := assets.GetPKLFileAsString("Utils.pkl")
		if err != nil {
			t.Fatalf("Failed to get Utils.pkl as string: %v", err)
		}
		if len(content) == 0 {
			t.Error("Utils.pkl should not be empty")
		}
	})

	t.Run("ListPKLFiles", func(t *testing.T) {
		files, err := assets.ListPKLFiles()
		if err != nil {
			t.Fatalf("Failed to list PKL files: %v", err)
		}
		if len(files) != 21 {
			t.Errorf("Expected 21 PKL files, got %d", len(files))
		}

		// Check for key files
		expectedFiles := map[string]bool{
			"Workflow.pkl": false, "Resource.pkl": false, "LLM.pkl": false,
			"APIServer.pkl": false, "Docker.pkl": false,
		}
		for _, file := range files {
			if _, exists := expectedFiles[file]; exists {
				expectedFiles[file] = true
			}
		}
		for file, found := range expectedFiles {
			if !found {
				t.Errorf("Expected file %s not found in list", file)
			}
		}
	})

	t.Run("ValidatePKLFiles", func(t *testing.T) {
		if err := assets.ValidatePKLFiles(); err != nil {
			t.Errorf("PKL file validation failed: %v", err)
		}
	})
}

// TestTmpDirFunctionality tests tmpdir-related functions
func TestTmpDirFunctionality(t *testing.T) {
	t.Run("GetTmpDir", func(t *testing.T) {
		tmpDir := assets.GetTmpDir()
		if tmpDir == "" {
			t.Error("GetTmpDir should return a non-empty path")
		}

		// Verify it's a real directory
		if _, err := os.Stat(tmpDir); err != nil {
			t.Errorf("TmpDir %s should be accessible: %v", tmpDir, err)
		}
	})

	t.Run("ExtractPKLFileToTemp", func(t *testing.T) {
		tempPath, err := assets.ExtractPKLFileToTemp("Workflow.pkl")
		if err != nil {
			t.Fatalf("Failed to extract PKL file to temp: %v", err)
		}
		defer os.Remove(tempPath)

		// Verify file exists
		if _, err := os.Stat(tempPath); err != nil {
			t.Errorf("Temp file should exist at %s: %v", tempPath, err)
		}

		// Verify it's in tmpdir
		if !strings.HasPrefix(tempPath, assets.GetTmpDir()) {
			t.Errorf("Temp file should be in tmpdir %s, got %s", assets.GetTmpDir(), tempPath)
		}

		// Verify content
		content, err := os.ReadFile(tempPath)
		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}
		if !strings.Contains(string(content), "Workflow Management") {
			t.Error("Temp file should contain original PKL content")
		}
	})

	t.Run("ExtractPKLFileWithName", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(assets.GetTmpDir(), "pkl_test_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		extractedPath, err := assets.ExtractPKLFileWithName("Workflow.pkl", tempDir)
		if err != nil {
			t.Fatalf("Failed to extract PKL file with name: %v", err)
		}

		// Verify filename is preserved
		expectedPath := filepath.Join(tempDir, "Workflow.pkl")
		if extractedPath != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, extractedPath)
		}

		// Verify file exists
		if _, err := os.Stat(extractedPath); err != nil {
			t.Errorf("Extracted file should exist: %v", err)
		}

		// Verify filename
		if filepath.Base(extractedPath) != "Workflow.pkl" {
			t.Errorf("Filename should be preserved as Workflow.pkl, got %s", filepath.Base(extractedPath))
		}
	})
}

// TestPKLWorkspace tests workspace functionality
func TestPKLWorkspace(t *testing.T) {
	t.Run("SetupPKLWorkspace", func(t *testing.T) {
		workspace, err := assets.SetupPKLWorkspace("")
		if err != nil {
			t.Fatalf("Failed to setup PKL workspace: %v", err)
		}
		defer workspace.Cleanup()

		// Test workspace properties
		if workspace.Directory == "" {
			t.Error("Workspace directory should not be empty")
		}
		if !workspace.IsTemporary() {
			t.Error("Workspace should be temporary when created with empty dir")
		}
		if !workspace.IsInTmpDir() {
			t.Error("Workspace should be in tmpdir")
		}

		// Test file access
		workflowPath := workspace.GetAbsolutePath("Workflow.pkl")
		if _, err := os.Stat(workflowPath); err != nil {
			t.Errorf("Workflow.pkl should exist in workspace: %v", err)
		}

		// Test import path
		importPath := workspace.GetImportPath("Resource.pkl")
		if !strings.HasSuffix(importPath, "Resource.pkl") {
			t.Errorf("Import path should end with Resource.pkl, got %s", importPath)
		}

		// Test file listing
		files, err := workspace.ListFiles()
		if err != nil {
			t.Fatalf("Failed to list workspace files: %v", err)
		}
		if len(files) != 21 {
			t.Errorf("Expected 21 files in workspace, got %d", len(files))
		}
	})

	t.Run("SetupPKLWorkspaceInTmpDir", func(t *testing.T) {
		workspace, err := assets.SetupPKLWorkspaceInTmpDir()
		if err != nil {
			t.Fatalf("Failed to setup PKL workspace in tmpdir: %v", err)
		}
		defer workspace.Cleanup()

		// Verify it's in tmpdir
		if !workspace.IsInTmpDir() {
			t.Error("Workspace should be in tmpdir")
		}
		if !workspace.IsTemporary() {
			t.Error("Workspace should be temporary")
		}

		// Verify all files have original names
		files, err := workspace.ListFiles()
		if err != nil {
			t.Fatalf("Failed to list files: %v", err)
		}

		expectedFiles := []string{"Workflow.pkl", "Resource.pkl", "LLM.pkl"}
		for _, expected := range expectedFiles {
			found := false
			for _, file := range files {
				if file == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected file %s not found in workspace", expected)
			}
		}
	})

	t.Run("WorkspaceCleanup", func(t *testing.T) {
		workspace, err := assets.SetupPKLWorkspaceInTmpDir()
		if err != nil {
			t.Fatalf("Failed to setup workspace: %v", err)
		}

		workspaceDir := workspace.Directory

		// Verify directory exists
		if _, err := os.Stat(workspaceDir); err != nil {
			t.Errorf("Workspace directory should exist before cleanup: %v", err)
		}

		// Cleanup
		if err := workspace.Cleanup(); err != nil {
			t.Errorf("Cleanup should not fail: %v", err)
		}

		// Verify directory is removed
		if _, err := os.Stat(workspaceDir); !os.IsNotExist(err) {
			t.Error("Workspace directory should be removed after cleanup")
		}
	})

	t.Run("CustomDirectory", func(t *testing.T) {
		customDir := filepath.Join(assets.GetTmpDir(), "custom_pkl_test")
		defer os.RemoveAll(customDir)

		workspace, err := assets.SetupPKLWorkspace(customDir)
		if err != nil {
			t.Fatalf("Failed to setup workspace with custom dir: %v", err)
		}
		defer workspace.Cleanup()

		// Should not be temporary since we provided a custom directory
		if workspace.IsTemporary() {
			t.Error("Workspace should not be temporary when custom dir is provided")
		}

		// Should still be in our custom location
		if !strings.Contains(workspace.Directory, "custom_pkl_test") {
			t.Errorf("Workspace should be in custom directory, got %s", workspace.Directory)
		}
	})
}

// TestExtractAllPKLFiles tests bulk extraction
func TestExtractAllPKLFiles(t *testing.T) {
	t.Run("ExtractToTempDir", func(t *testing.T) {
		extractedDir, err := assets.ExtractAllPKLFilesToDir("")
		if err != nil {
			t.Fatalf("Failed to extract all PKL files: %v", err)
		}
		defer os.RemoveAll(extractedDir)

		// Verify it's in tmpdir
		if !strings.HasPrefix(extractedDir, assets.GetTmpDir()) {
			t.Errorf("Extracted dir should be in tmpdir %s, got %s", assets.GetTmpDir(), extractedDir)
		}

		// Verify all files exist with original names
		expectedFiles := []string{"Workflow.pkl", "Resource.pkl", "LLM.pkl", "APIServer.pkl"}
		for _, filename := range expectedFiles {
			filePath := filepath.Join(extractedDir, filename)
			if _, err := os.Stat(filePath); err != nil {
				t.Errorf("File %s should exist at %s: %v", filename, filePath, err)
			}
		}

		// Count total files
		entries, err := os.ReadDir(extractedDir)
		if err != nil {
			t.Fatalf("Failed to read extracted directory: %v", err)
		}

		pklCount := 0
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".pkl") {
				pklCount++
			}
		}
		if pklCount != 21 {
			t.Errorf("Expected 21 PKL files, found %d", pklCount)
		}
	})

	t.Run("ExtractToCustomDir", func(t *testing.T) {
		customDir := filepath.Join(assets.GetTmpDir(), "custom_extract_test")
		defer os.RemoveAll(customDir)

		extractedDir, err := assets.ExtractAllPKLFilesToDir(customDir)
		if err != nil {
			t.Fatalf("Failed to extract to custom dir: %v", err)
		}

		if extractedDir != customDir {
			t.Errorf("Should return the custom directory path, got %s", extractedDir)
		}

		// Verify files exist
		workflowPath := filepath.Join(customDir, "Workflow.pkl")
		if _, err := os.Stat(workflowPath); err != nil {
			t.Errorf("Workflow.pkl should exist in custom dir: %v", err)
		}
	})
}

// TestGetPKLFilePath tests path identifiers
func TestGetPKLFilePath(t *testing.T) {
	path := assets.GetPKLFilePath("Workflow.pkl")
	expected := "embed://pkl/Workflow.pkl"
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

// BenchmarkPKLFileOperations benchmarks key operations
func BenchmarkPKLFileOperations(b *testing.B) {
	b.Run("GetPKLFile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := assets.GetPKLFile("Workflow.pkl")
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ListPKLFiles", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := assets.ListPKLFiles()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("SetupWorkspace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			workspace, err := assets.SetupPKLWorkspaceInTmpDir()
			if err != nil {
				b.Fatal(err)
			}
			workspace.Cleanup()
		}
	})
}
