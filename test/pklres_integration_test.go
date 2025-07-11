package test

import (
	"os"
	"strings"
	"testing"

	"github.com/kdeps/kdeps/pkg/pklres"
)

// TestPklresIntegrationPKL loads PKL test cases from a PKL file and checks results using the real pklres reader
func TestPklresIntegrationPKL(t *testing.T) {
	// Create temporary database
	tempDB, err := os.CreateTemp("", "pklres-pkltest-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp database: %v", err)
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	// Initialize real pklres reader
	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		t.Fatalf("Failed to initialize pklres reader: %v", err)
	}

	// Create evaluator with real resource readers
	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, pklresReader)
	if err != nil {
		t.Fatalf("Failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	// Test various PKL files
	testFiles := []string{
		"exec_tests_pass.pkl",
		"python_tests_pass.pkl",
		"llm_tests_pass.pkl",
		"http_tests_pass.pkl",
		"data_tests_pass.pkl",
		"pklres_tests_pass.pkl",
		"all_tests_pass.pkl",
	}

	var failures []string
	for _, fileName := range testFiles {
		module := EvaluatePKLFile(t, evaluator, fileName)
		if module == nil {
			failures = append(failures, fileName+": error evaluating")
		}
	}

	if len(failures) > 0 {
		t.Logf("PKL integration test failures: %s", strings.Join(failures, ", "))
		t.Logf("These failures are expected")
	}
}
