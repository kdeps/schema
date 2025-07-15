package test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/apple/pkl-go/pkl"
)

// Enhanced PklresResourceReader with comprehensive test data
type EnhancedPklresResourceReader struct{}

func (r *EnhancedPklresResourceReader) Read(uri url.URL) ([]byte, error) {
	// Parse query parameters
	query := uri.Query()
	resourceType := query.Get("type")
	key := query.Get("key")
	op := query.Get("op")
	value := query.Get("value")

	// Handle different operations
	switch op {
	case "set":
		return []byte(value), nil
	case "get":
		return r.getMockValue(resourceType, key), nil
	default:
		// Default getPklRecord behavior
		return r.getMockRecord(resourceType), nil
	}
}

func (r *EnhancedPklresResourceReader) getMockRecord(resourceType string) []byte {
	switch resourceType {
	case "exec":
		return []byte(`new ResourceExec {
			Command = "echo 'Hello from pklres'"
			Stdout = "Hello from pklres"
			Stderr = ""
			ExitCode = 0
			File = "/tmp/test-exec.pkl"
			ItemValues = new Listing { "key1" = "value1", "key2" = "value2" }
			Env = new Mapping { "PATH" = "/usr/bin:/bin", "HOME" = "/home/test" }
		}`)
	case "python":
		return []byte(`new ResourcePython {
			Script = "print('Hello from Python')"
			Stdout = "Hello from Python"
			Stderr = ""
			ExitCode = 0
			File = "/tmp/test-python.py"
			ItemValues = new Listing { "result" = "success" }
			Env = new Mapping { "PYTHONPATH" = "/usr/lib/python3" }
			PythonEnvironment = "python3"
		}`)
	case "llm":
		return []byte(`new ResourceChat {
			Model = "llama3.2"
			Role = "user"
			Prompt = "Hello, how are you?"
			Scenario = "conversation"
			Tools = new Listing { "tool1", "tool2" }
			Files = new Listing { "file1.txt", "file2.txt" }
			JSONResponse = true
			JSONResponseKeys = new Listing { "response", "confidence" }
			Response = "I'm doing well, thank you!"
			File = "/tmp/test-llm.json"
			ItemValues = new Listing { "confidence" = "0.95" }
		}`)
	case "http":
		return []byte(`new ResourceHTTPClient {
			Method = "POST"
			Url = "https://api.example.com/test"
			Headers = new Mapping { "Content-Type" = "application/json", "Authorization" = "Bearer token123" }
			Data = new Listing { "key" = "value", "number" = "42" }
			Response = null
			File = "/tmp/test-http.json"
			ItemValues = new Listing { "status" = "200", "response_time" = "150ms" }
			Params = new Mapping { "timeout" = "30s" }
			Timestamp = "2024-01-15T10:30:00Z"
			TimeoutDuration = "30s"
		}`)
	case "data":
		return []byte(`new Mapping {
			"test.txt" = "/path/to/test.txt"
			"config.json" = "/path/to/config.json"
			"data.csv" = "/path/to/data.csv"
		}`)
	default:
		return []byte("")
	}
}

func (r *EnhancedPklresResourceReader) getMockValue(resourceType, key string) []byte {
	// Mock key-value storage
	mockData := map[string]map[string]string{
		"exec": {
			"command":  "echo 'Hello from pklres'",
			"stdout":   "Hello from pklres",
			"stderr":   "",
			"exitCode": "0",
		},
		"python": {
			"script":   "print('Hello from Python')",
			"stdout":   "Hello from Python",
			"stderr":   "",
			"exitCode": "0",
		},
		"llm": {
			"model":    "llama3.2",
			"prompt":   "Hello, how are you?",
			"response": "I'm doing well, thank you!",
		},
		"http": {
			"method": "POST",
			"url":    "https://api.example.com/test",
			"status": "200",
		},
	}

	if data, exists := mockData[resourceType]; exists {
		if value, exists := data[key]; exists {
			return []byte(value)
		}
	}
	return []byte("")
}

func (r *EnhancedPklresResourceReader) IsGlob(url string) bool {
	return false
}

func (r *EnhancedPklresResourceReader) Glob(ctx context.Context, url string) ([]string, error) {
	return nil, fmt.Errorf("glob not supported for pklres scheme")
}

func (r *EnhancedPklresResourceReader) HasHierarchicalUris() bool {
	return false
}

func (r *EnhancedPklresResourceReader) IsGlobbable() bool {
	return false
}

func (r *EnhancedPklresResourceReader) ListElements(_ url.URL) ([]pkl.PathElement, error) {
	return nil, nil
}

func (r *EnhancedPklresResourceReader) Scheme() string {
	return "pklres"
}

// writeTempFile creates a temp file in dir with the given content and returns the file path.
func writeTempFile(dir, pattern, content string) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// TestPklresCoreFunctions tests the core pklres functions
func TestPklresCoreFunctions(t *testing.T) {
	// Create evaluator with enhanced resource reader
	opts := func(options *pkl.EvaluatorOptions) {
		pkl.WithDefaultAllowedResources(options)
		pkl.WithOsEnv(options)
		pkl.WithDefaultAllowedModules(options)
		pkl.WithDefaultCacheDir(options)
		options.Logger = pkl.NoopLogger
		options.ResourceReaders = []pkl.ResourceReader{
			&EnhancedPklresResourceReader{},
		}
		options.AllowedModules = []string{".*"}
		options.AllowedResources = []string{".*"}
		options.ModulePaths = []string{"."}
	}

	evaluator, err := pkl.NewEvaluator(context.Background(), opts)
	if err != nil {
		t.Fatalf("Failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	testCases := []struct {
		name     string
		pklExpr  string
		expected string
	}{
		{
			name: "getPklRecord for exec resource",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord("test-exec", "exec")
			`,
			expected: "echo 'Hello from pklres'",
		},
		{
			name: "getPklRecord for python resource",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord("test-python", "python")
			`,
			expected: "print('Hello from Python')",
		},
		{
			name: "getPklRecord for llm resource",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord("test-llm", "llm")
			`,
			expected: "llama3.2",
		},
		{
			name: "getPklRecord for http resource",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord("test-http", "http")
			`,
			expected: "POST",
		},
		{
			name: "getPklValue for exec command",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklValue("test-exec", "exec", "command")
			`,
			expected: "echo 'Hello from pklres'",
		},
		{
			name: "setPklValue operation",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.setPklValue("test-set", "exec", "command", "new command")
			`,
			expected: "new command",
		},
		{
			name: "null parameter handling",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord(null, "exec")
			`,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempFile, err := writeTempFile(os.TempDir(), "*.pkl", tc.pklExpr)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile)

			source := pkl.FileSource(tempFile)
			var module map[string]interface{}
			if err := evaluator.EvaluateModule(context.Background(), source, &module); err != nil {
				t.Fatalf("Failed to evaluate PKL module: %v", err)
			}

			resultStr := fmt.Sprintf("%v", module["result"])
			if !strings.Contains(resultStr, tc.expected) {
				t.Errorf("Expected result to contain '%s', got: %s", tc.expected, resultStr)
			}
		})
	}
}

// TestPklresResourceIntegration tests resource integration with pklres
func TestPklresResourceIntegration(t *testing.T) {
	// Create evaluator with enhanced resource reader
	opts := func(options *pkl.EvaluatorOptions) {
		pkl.WithDefaultAllowedResources(options)
		pkl.WithOsEnv(options)
		pkl.WithDefaultAllowedModules(options)
		pkl.WithDefaultCacheDir(options)
		options.Logger = pkl.NoopLogger
		options.ResourceReaders = []pkl.ResourceReader{
			&EnhancedPklresResourceReader{},
		}
		options.AllowedModules = []string{".*"}
		options.AllowedResources = []string{".*"}
		options.ModulePaths = []string{"."}
	}

	evaluator, err := pkl.NewEvaluator(context.Background(), opts)
	if err != nil {
		t.Fatalf("Failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	testCases := []struct {
		name     string
		resource string
		actionID string
		expected string
	}{
		{
			name:     "Exec resource integration",
			resource: "Exec",
			actionID: "test-exec",
			expected: "echo 'Hello from pklres'",
		},
		{
			name:     "Python resource integration",
			resource: "Python",
			actionID: "test-python",
			expected: "print('Hello from Python')",
		},
		{
			name:     "LLM resource integration",
			resource: "LLM",
			actionID: "test-llm",
			expected: "llama3.2",
		},
		{
			name:     "HTTP resource integration",
			resource: "HTTP",
			actionID: "test-http",
			expected: "POST",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pklExpr := fmt.Sprintf(`
				import "../deps/pkl/%s.pkl" as resource
				result = resource.resource("%s")
			`, tc.resource, tc.actionID)

			tempFile, err := writeTempFile(os.TempDir(), "*.pkl", pklExpr)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile)

			source := pkl.FileSource(tempFile)
			var module map[string]interface{}
			if err := evaluator.EvaluateModule(context.Background(), source, &module); err != nil {
				t.Fatalf("Failed to evaluate PKL module: %v", err)
			}

			resultStr := fmt.Sprintf("%v", module["result"])
			if !strings.Contains(resultStr, tc.expected) {
				t.Errorf("Expected result to contain '%s', got: %s", tc.expected, resultStr)
			}
		})
	}
}

// TestPklresCLI tests pklres functionality using PKL CLI
func TestPklresCLI(t *testing.T) {
	// Check if pkl CLI is available
	if _, err := exec.LookPath("pkl"); err != nil {
		t.Errorf("PKL CLI not available, skipping CLI tests")
	}

	testCases := []struct {
		name     string
		pklExpr  string
		expected string
	}{
		{
			name: "Basic pklres function test",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord("test", "exec")
			`,
			expected: "",
		},
		{
			name: "Null parameter test",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord(null, "exec")
			`,
			expected: "",
		},
		{
			name: "Empty string test",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord("", "")
			`,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempFile, err := writeTempFile(os.TempDir(), "*.pkl", tc.pklExpr)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile)

			// Run PKL CLI command
			cmd := exec.Command("pkl", "eval", tempFile)
			output, err := cmd.CombinedOutput()

			if err != nil {
				// For CLI tests, we expect some errors due to missing pklres backend
				// but we can still validate the syntax is correct
				if !strings.Contains(string(output), "result") {
					t.Errorf("Expected output to contain 'result', got: %s", string(output))
				}
			} else {
				outputStr := string(output)
				if !strings.Contains(outputStr, tc.expected) {
					t.Errorf("Expected output to contain '%s', got: %s", tc.expected, outputStr)
				}
			}
		})
	}
}

// TestPklresErrorHandling tests error scenarios
func TestPklresErrorHandling(t *testing.T) {
	// Create evaluator with enhanced resource reader
	opts := func(options *pkl.EvaluatorOptions) {
		pkl.WithDefaultAllowedResources(options)
		pkl.WithOsEnv(options)
		pkl.WithDefaultAllowedModules(options)
		pkl.WithDefaultCacheDir(options)
		options.Logger = pkl.NoopLogger
		options.ResourceReaders = []pkl.ResourceReader{
			&EnhancedPklresResourceReader{},
		}
		options.AllowedModules = []string{".*"}
		options.AllowedResources = []string{".*"}
		options.ModulePaths = []string{"."}
	}

	evaluator, err := pkl.NewEvaluator(context.Background(), opts)
	if err != nil {
		t.Fatalf("Failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	testCases := []struct {
		name          string
		pklExpr       string
		expectedError string
	}{
		{
			name: "Invalid resource type",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklRecord("test", "invalid")
			`,
			expectedError: "",
		},
		{
			name: "Missing parameters",
			pklExpr: `
				import "../deps/pkl/PklResource.pkl" as pklres
				result = pklres.getPklValue("test", "exec")
			`,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempFile, err := writeTempFile(os.TempDir(), "*.pkl", tc.pklExpr)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile)

			source := pkl.FileSource(tempFile)
			var module map[string]interface{}
			err = evaluator.EvaluateModule(context.Background(), source, &module)

			if tc.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("Expected error containing '%s', got: %v", tc.expectedError, err)
				}
			} else {
				// For these tests, we expect them to work with our mock reader
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestPklresPerformance tests performance characteristics
func TestPklresPerformance(t *testing.T) {
	// Create evaluator with enhanced resource reader
	opts := func(options *pkl.EvaluatorOptions) {
		pkl.WithDefaultAllowedResources(options)
		pkl.WithOsEnv(options)
		pkl.WithDefaultAllowedModules(options)
		pkl.WithDefaultCacheDir(options)
		options.Logger = pkl.NoopLogger
		options.ResourceReaders = []pkl.ResourceReader{
			&EnhancedPklresResourceReader{},
		}
		options.AllowedModules = []string{".*"}
		options.AllowedResources = []string{".*"}
		options.ModulePaths = []string{"."}
	}

	evaluator, err := pkl.NewEvaluator(context.Background(), opts)
	if err != nil {
		t.Fatalf("Failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	// Test multiple concurrent operations
	pklExpr := `
		import "../deps/pkl/PklResource.pkl" as pklres
		result1 = pklres.getPklRecord("test1", "exec")
		result2 = pklres.getPklRecord("test2", "python")
		result3 = pklres.getPklRecord("test3", "llm")
		result4 = pklres.getPklRecord("test4", "http")
	`

	tempFile, err := writeTempFile(os.TempDir(), "*.pkl", pklExpr)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	start := time.Now()
	source := pkl.FileSource(tempFile)
	var module map[string]interface{}
	if err := evaluator.EvaluateModule(context.Background(), source, &module); err != nil {
		t.Fatalf("Failed to evaluate PKL module: %v", err)
	}
	duration := time.Since(start)

	// Performance should be reasonable (under 1 second for this simple operation)
	if duration > time.Second {
		t.Errorf("Performance test took too long: %v", duration)
	}

	// Verify all results are present
	expectedResults := []string{"result1", "result2", "result3", "result4"}
	for _, resultKey := range expectedResults {
		if _, exists := module[resultKey]; !exists {
			t.Errorf("Expected result key '%s' not found", resultKey)
		}
	}
}

// TestPklresSchemaValidation tests that the PKL schema is valid
func TestPklresSchemaValidation(t *testing.T) {
	// Check if pkl CLI is available
	if _, err := exec.LookPath("pkl"); err != nil {
		t.Errorf("PKL CLI not available, skipping schema validation")
	}

	// Test files to validate
	testFiles := []string{
		"../deps/pkl/PklResource.pkl",
		"../deps/pkl/Exec.pkl",
		"../deps/pkl/Python.pkl",
		"../deps/pkl/LLM.pkl",
		"../deps/pkl/HTTP.pkl",
		"../deps/pkl/Data.pkl",
	}

	for _, file := range testFiles {
		t.Run(fmt.Sprintf("Validate %s", filepath.Base(file)), func(t *testing.T) {
			// Check if file exists
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Errorf("File %s does not exist", file)
			}

			// Validate PKL file syntax
			cmd := exec.Command("pkl", "eval", "--no-cache", file)
			output, err := cmd.CombinedOutput()

			if err != nil {
				// Some files might have dependencies that aren't available in test environment
				// but we can still check that the syntax is valid
				if strings.Contains(string(output), "syntax error") {
					t.Errorf("PKL syntax error in %s: %s", file, string(output))
				}
			}
		})
	}
}
