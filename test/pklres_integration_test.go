package test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/kdeps/pkg/pklres"
)

// TestPklresIntegrationPKL loads PKL test cases from a PKL file and checks results using the real pklres reader
func TestPklresIntegrationPKL(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pklres-pkltest-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Recursively copy deps/pkl directory to tempDir
	srcDeps := filepath.Clean(filepath.Join("..", "deps", "pkl"))
	dstDeps := filepath.Join(tempDir, "deps", "pkl")
	if err := os.MkdirAll(dstDeps, 0755); err != nil {
		t.Fatalf("Failed to create deps/pkl dir: %v", err)
	}
	entries, err := os.ReadDir(srcDeps)
	if err != nil {
		t.Fatalf("Failed to read deps/pkl dir: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDeps, entry.Name())
		dst := filepath.Join(dstDeps, entry.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", dst, err)
		}
	}

	// Also copy test_pklres_integration.pkl to tempDir
	testFiles := []string{
		"test_pklres_integration.pkl",
		"exec_tests_pass.pkl",
		"python_tests_pass.pkl",
		"llm_tests_pass.pkl",
		"http_tests_pass.pkl",
		"data_tests_pass.pkl",
		"pklres_tests_pass.pkl",
		"all_tests_pass.pkl",
		"test_summary.pkl",
	}
	for _, fname := range testFiles {
		src := filepath.Join(".", fname)
		dst := filepath.Join(tempDir, fname)
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", src, err)
		}
		// Rewrite import paths for tempDir
		updated := strings.ReplaceAll(string(data), "../deps/pkl/", "deps/pkl/")
		if err := os.WriteFile(dst, []byte(updated), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", dst, err)
		}
	}

	// Create a temp DB for the pklres reader
	tempDB, err := os.CreateTemp("", "pklres-pkltest-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp DB: %v", err)
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		t.Fatalf("Failed to initialize pklres reader: %v", err)
	}

	// Register both pklres and agent resource readers
	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, pklresReader)
	if err != nil {
		t.Fatalf("Failed to create PKL evaluator: %v", err)
	}
	defer evaluator.Close()

	// Change working directory to tempDir for PKL import resolution
	origWD, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to chdir to tempDir: %v", err)
	}
	defer os.Chdir(origWD)

	// List of PKL files to check (new approach: direct value return)
	passFiles := []string{
		"exec_tests_pass.pkl",
		"python_tests_pass.pkl",
		"llm_tests_pass.pkl",
		"http_tests_pass.pkl",
		"data_tests_pass.pkl",
		"pklres_tests_pass.pkl",
		"all_tests_pass.pkl",
	}
	var failed []string
	for _, fname := range passFiles {
		source := pkl.FileSource(fname)
		var module map[string]interface{}
		if err := evaluator.EvaluateModule(context.Background(), source, &module); err != nil {
			failed = append(failed, fname+": error evaluating")
			t.Logf("Error evaluating %s: %v", fname, err)
			continue
		}
		val, ok := module["result"].(bool)
		if !ok || !val {
			failed = append(failed, fname)
		}
	}
	if len(failed) > 0 {
		t.Errorf("PKL integration test failures: %s", strings.Join(failed, ", "))
	}

	// Print test summary from test_summary.pkl
	source := pkl.FileSource("test_summary.pkl")
	var summaryModule map[string]interface{}
	if err := evaluator.EvaluateModule(context.Background(), source, &summaryModule); err == nil {
		if summary, ok := summaryModule["result"].(string); ok {
			t.Logf("Test Results:\n%s", summary)
		}
	}
}
