package test

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kdeps/kdeps/pkg/pklres"
)

// TestIntegrationSuite runs all integration tests with comprehensive reporting
func TestIntegrationSuite(t *testing.T) {
	suite := NewTestSuite()
	suite.metrics.StartTime = time.Now()

	// Run all test categories
	t.Run("Resource_Readers", func(t *testing.T) {
		suite.RunTest(t, "Agent_Resource_Reader", testAgentResourceReader)
		suite.RunTest(t, "Pklres_Resource_Reader", testPklresResourceReader)
		suite.RunTest(t, "Real_Pklres_Reader", testRealPklresReader)
	})

	t.Run("PKL_Integration", func(t *testing.T) {
		suite.RunTest(t, "PKL_File_Evaluation", testPKLFileEvaluation)
		suite.RunTest(t, "PKL_Resource_Integration", testPKLResourceIntegration)
		suite.RunTest(t, "PKL_Complex_Workflows", testPKLComplexWorkflows)
	})

	t.Run("Schema_Validation", func(t *testing.T) {
		suite.RunTest(t, "Schema_Compatibility", testSchemaCompatibility)
		suite.RunTest(t, "Resource_Type_Validation", testResourceTypeValidation)
		suite.RunTest(t, "Import_Path_Resolution", testImportPathResolution)
	})

	t.Run("Performance_Tests", func(t *testing.T) {
		suite.RunTest(t, "Resource_Reader_Performance", testResourceReaderPerformance)
		suite.RunTest(t, "PKL_Evaluation_Performance", testPKLEvaluationPerformance)
		suite.RunTest(t, "Concurrent_Operations", testConcurrentOperations)
	})

	// Print comprehensive test summary
	suite.PrintSummary()
}

// testAgentResourceReader tests the agent resource reader functionality
func testAgentResourceReader(t *testing.T) error {
	reader := &AgentResourceReader{}

	// Test basic agent resolution
	uri, _ := url.Parse("agent:/test-action")
	result, err := reader.Read(*uri)
	if err != nil {
		return err
	}

	// Verify result structure
	if result == nil {
		return fmt.Errorf("expected non-nil result from agent reader")
	}

	return nil
}

// testPklresResourceReader tests the pklres resource reader functionality
func testPklresResourceReader(t *testing.T) error {
	reader := &PklresResourceReader{}

	// Test get operation
	uri, _ := url.Parse("pklres:/test-id?type=exec&key=command&op=get")
	result, err := reader.Read(*uri)
	if err != nil {
		return err
	}

	if result == nil {
		return fmt.Errorf("expected non-nil result from pklres reader")
	}

	// Test set operation
	setURI, _ := url.Parse("pklres:/test-id?type=exec&key=command&op=set&value=echo%20hello")
	_, err = reader.Read(*setURI)
	if err != nil {
		return err
	}

	return nil
}

// testRealPklresReader tests the real pklres reader with database
func testRealPklresReader(t *testing.T) error {
	// Create temporary database
	tempDB, err := os.CreateTemp("", "pklres-test-*.db")
	if err != nil {
		return err
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	// Initialize real pklres reader
	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		return err
	}

	// Test set operation
	setURI, _ := url.Parse("pklres:/real-test-id?op=set&type=exec&key=command&value=echo%20real%20test")
	_, err = pklresReader.Read(*setURI)
	if err != nil {
		return err
	}

	// Test get operation
	getURI, _ := url.Parse("pklres:/real-test-id?op=get&type=exec&key=command")
	result, err := pklresReader.Read(*getURI)
	if err != nil {
		return err
	}

	if result == nil {
		return fmt.Errorf("expected non-nil result from real pklres reader")
	}

	return nil
}

// testPKLFileEvaluation tests PKL file evaluation with various file types
func testPKLFileEvaluation(t *testing.T) error {
	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, &PklresResourceReader{})
	if err != nil {
		return err
	}
	defer evaluator.Close()

	// Test different PKL file types
	testFiles := []string{
		"exec_tests_pass.pkl",
		"python_tests_pass.pkl",
		"llm_tests_pass.pkl",
		"http_tests_pass.pkl",
		"data_tests_pass.pkl",
		"pklres_tests_pass.pkl",
		"all_tests_pass.pkl",
		"test_summary.pkl",
	}

	for _, fileName := range testFiles {
		module := EvaluatePKLFile(t, evaluator, fileName)
		if module == nil {
			return fmt.Errorf("failed to evaluate %s", fileName)
		}
	}

	return nil
}

// testPKLResourceIntegration tests PKL integration with resource readers
func testPKLResourceIntegration(t *testing.T) error {
	// Create temporary workspace
	tempDir, cleanup := CreateTempPKLWorkspace(t)
	defer cleanup()

	// Copy test files
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

	for _, fileName := range testFiles {
		CopyPKLFile(t, tempDir, fileName)
	}

	// Create evaluator with real resource readers
	tempDB, err := os.CreateTemp("", "pklres-integration-*.db")
	if err != nil {
		return err
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		return err
	}

	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, pklresReader)
	if err != nil {
		return err
	}
	defer evaluator.Close()

	// Test integration file
	integrationFile := filepath.Join(tempDir, "test_pklres_integration.pkl")
	module := EvaluatePKLFile(t, evaluator, integrationFile)
	if module == nil {
		return fmt.Errorf("failed to evaluate integration file")
	}

	return nil
}

// testPKLComplexWorkflows tests complex PKL workflows
func testPKLComplexWorkflows(t *testing.T) error {
	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, &PklresResourceReader{})
	if err != nil {
		return err
	}
	defer evaluator.Close()

	// Test complex workflow scenarios
	workflows := []struct {
		name     string
		fileName string
	}{
		{"Multi_Resource_Workflow", "all_tests_pass.pkl"},
		{"Test_Summary_Workflow", "test_summary.pkl"},
	}

	for _, workflow := range workflows {
		module := EvaluatePKLFile(t, evaluator, workflow.fileName)
		if module == nil {
			return fmt.Errorf("failed to evaluate workflow %s", workflow.name)
		}
	}

	return nil
}

// testSchemaCompatibility tests schema compatibility
func testSchemaCompatibility(t *testing.T) error {
	// Test that all PKL files are compatible with the current schema
	testFiles := []string{
		"exec_tests_pass.pkl",
		"python_tests_pass.pkl",
		"llm_tests_pass.pkl",
		"http_tests_pass.pkl",
		"data_tests_pass.pkl",
		"pklres_tests_pass.pkl",
		"all_tests_pass.pkl",
		"test_summary.pkl",
	}

	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, &PklresResourceReader{})
	if err != nil {
		return err
	}
	defer evaluator.Close()

	for _, fileName := range testFiles {
		module := EvaluatePKLFile(t, evaluator, fileName)
		if module == nil {
			return fmt.Errorf("schema compatibility test failed for %s", fileName)
		}
	}

	return nil
}

// testResourceTypeValidation tests resource type validation
func testResourceTypeValidation(t *testing.T) error {
	// Test that all resource types are properly validated
	resourceTypes := []string{"exec", "python", "llm", "http", "data"}

	for _, resourceType := range resourceTypes {
		// Test with mock readers
		reader := &PklresResourceReader{}
		uri, _ := url.Parse(fmt.Sprintf("pklres:/test-id?type=%s&key=test&op=get", resourceType))
		_, err := reader.Read(*uri)
		if err != nil {
			return fmt.Errorf("resource type validation failed for %s: %v", resourceType, err)
		}
	}

	return nil
}

// testImportPathResolution tests import path resolution
func testImportPathResolution(t *testing.T) error {
	// Create temporary workspace
	tempDir, cleanup := CreateTempPKLWorkspace(t)
	defer cleanup()

	// Test that import paths are properly resolved
	testFiles := []string{
		"exec_tests_pass.pkl",
		"python_tests_pass.pkl",
		"llm_tests_pass.pkl",
		"http_tests_pass.pkl",
		"data_tests_pass.pkl",
		"pklres_tests_pass.pkl",
	}

	for _, fileName := range testFiles {
		CopyPKLFile(t, tempDir, fileName)
	}

	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, &PklresResourceReader{})
	if err != nil {
		return err
	}
	defer evaluator.Close()

	// Test each file with updated import paths
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		module := EvaluatePKLFile(t, evaluator, filePath)
		if module == nil {
			return fmt.Errorf("import path resolution failed for %s", fileName)
		}
	}

	return nil
}

// testResourceReaderPerformance tests resource reader performance
func testResourceReaderPerformance(t *testing.T) error {
	// Test performance with multiple operations
	reader := &PklresResourceReader{}

	start := time.Now()
	for i := 0; i < 100; i++ {
		uri, _ := url.Parse(fmt.Sprintf("pklres:/perf-test-%d?type=exec&key=command&op=get", i))
		_, err := reader.Read(*uri)
		if err != nil {
			return err
		}
	}
	duration := time.Since(start)

	// Ensure performance is reasonable (less than 1 second for 100 operations)
	if duration > time.Second {
		return fmt.Errorf("resource reader performance test failed: %v for 100 operations", duration)
	}

	return nil
}

// testPKLEvaluationPerformance tests PKL evaluation performance
func testPKLEvaluationPerformance(t *testing.T) error {
	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, &PklresResourceReader{})
	if err != nil {
		return err
	}
	defer evaluator.Close()

	// Test performance with multiple file evaluations
	start := time.Now()
	for i := 0; i < 10; i++ {
		module := EvaluatePKLFile(t, evaluator, "test_summary.pkl")
		if module == nil {
			return fmt.Errorf("PKL evaluation failed on iteration %d", i)
		}
	}
	duration := time.Since(start)

	// Ensure performance is reasonable (less than 5 seconds for 10 evaluations)
	if duration > 5*time.Second {
		return fmt.Errorf("PKL evaluation performance test failed: %v for 10 evaluations", duration)
	}

	return nil
}

// testConcurrentOperations tests concurrent resource operations
func testConcurrentOperations(t *testing.T) error {
	// Create temporary database for concurrent testing
	tempDB, err := os.CreateTemp("", "pklres-concurrent-*.db")
	if err != nil {
		return err
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		return err
	}

	// Test concurrent set operations
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			uri, _ := url.Parse(fmt.Sprintf("pklres:/concurrent-test-%d?op=set&type=exec&key=command&value=echo%%20concurrent%%20%d", id, id))
			_, err := pklresReader.Read(*uri)
			done <- err
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < 10; i++ {
		if err := <-done; err != nil {
			return fmt.Errorf("concurrent operation failed: %v", err)
		}
	}

	return nil
}
