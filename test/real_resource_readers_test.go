package test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/kdeps/pkg/agent"
	"github.com/kdeps/kdeps/pkg/logging"
	"github.com/kdeps/kdeps/pkg/pklres"
	"github.com/spf13/afero"
)

// TestRealAgentReader tests the actual agent resource reader from kdeps
func TestRealAgentReader(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "agent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test agents directory structure
	agentsDir := filepath.Join(tempDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("Failed to create agents directory: %v", err)
	}

	// Create a test agent with workflow.pkl
	testAgentDir := filepath.Join(agentsDir, "test-agent", "1.0.0")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("Failed to create test agent directory: %v", err)
	}

	// Create a simple workflow.pkl file
	workflowContent := `
@ModuleInfo { minPklVersion = "0.28.2" }

module test.Workflow

ActionID = "test-action"
Description = "Test action for agent reader testing"
`
	if err := os.WriteFile(filepath.Join(testAgentDir, "workflow.pkl"), []byte(workflowContent), 0644); err != nil {
		t.Fatalf("Failed to write workflow.pkl: %v", err)
	}

	// Create logger
	logger := logging.NewTestLogger()

	// Initialize the real agent reader
	fs := afero.NewOsFs()
	agentReader, err := agent.InitializeAgent(fs, tempDir, "test-agent", "1.0.0", logger)
	if err != nil {
		t.Fatalf("Failed to initialize agent reader: %v", err)
	}
	defer agentReader.Close()

	testCases := []struct {
		name     string
		uri      string
		expected string
		contains bool
	}{
		{
			name:     "Resolve local action ID",
			uri:      "agent:/test-action",
			expected: "@test-agent/test-action:1.0.0",
			contains: true,
		},
		{
			name:     "Resolve with query parameters",
			uri:      "agent:/test-action?agent=test-agent&version=1.0.0",
			expected: "@test-agent/test-action:1.0.0",
			contains: true,
		},
		{
			name:     "List installed agents",
			uri:      "agent:/?op=list-installed",
			expected: "test-agent",
			contains: true,
		},
		{
			name:     "List agent resources",
			uri:      "agent:/test-agent?op=list&agent=test-agent&version=1.0.0",
			expected: "copy_all_resources",
			contains: true,
		},
		{
			name:     "Resolve canonical agent ID",
			uri:      "agent:/@test-agent:1.0.0",
			expected: "@test-agent:1.0.0",
			contains: true,
		},
		{
			name:     "Resolve canonical action ID",
			uri:      "agent:/@test-agent/test-action:1.0.0",
			expected: "@test-agent/test-action:1.0.0",
			contains: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedURL, err := url.Parse(tc.uri)
			if err != nil {
				t.Fatalf("Failed to parse URI: %v", err)
			}

			result, err := agentReader.Read(*parsedURL)
			if err != nil {
				t.Fatalf("Failed to read from agent reader: %v", err)
			}

			resultStr := string(result)
			if tc.contains {
				if !strings.Contains(resultStr, tc.expected) {
					t.Errorf("Expected result to contain '%s', got: %s", tc.expected, resultStr)
				}
			} else {
				if resultStr != tc.expected {
					t.Errorf("Expected '%s', got: %s", tc.expected, resultStr)
				}
			}
		})
	}
}

// TestRealPklresReader tests the actual pklres resource reader from kdeps
func TestRealPklresReader(t *testing.T) {
	// Create a temporary database file
	tempDB, err := os.CreateTemp("", "pklres-test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp database: %v", err)
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	// Initialize the real pklres reader
	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		t.Fatalf("Failed to initialize pklres reader: %v", err)
	}
	// Note: pklresReader doesn't have a Close method, database will be cleaned up when temp file is removed

	testCases := []struct {
		name     string
		uri      string
		expected string
		contains bool
	}{
		{
			name:     "Set record",
			uri:      "pklres:/test-id?op=set&type=exec&key=command&value=echo hello",
			expected: "echo hello",
			contains: false,
		},
		{
			name:     "Get record",
			uri:      "pklres:/test-id?op=get&type=exec&key=command",
			expected: "echo hello",
			contains: false,
		},
		{
			name:     "Get non-existent record",
			uri:      "pklres:/nonexistent?op=get&type=exec&key=command",
			expected: "",
			contains: false,
		},
		{
			name:     "Set record without key",
			uri:      "pklres:/test-id2?op=set&type=python&value=print('hello')",
			expected: "print('hello')",
			contains: false,
		},
		{
			name:     "Get record without key",
			uri:      "pklres:/test-id2?op=get&type=python",
			expected: "print('hello')",
			contains: false,
		},
		{
			name:     "List records",
			uri:      "pklres:/?op=list&type=exec",
			expected: "test-id",
			contains: true,
		},
		{
			name:     "Delete specific key",
			uri:      "pklres:/test-id?op=delete&type=exec&key=command",
			expected: "Deleted 1 record(s)",
			contains: true,
		},
		{
			name:     "Clear records by type",
			uri:      "pklres:/?op=clear&type=python",
			expected: "Cleared 1 records",
			contains: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedURL, err := url.Parse(tc.uri)
			if err != nil {
				t.Fatalf("Failed to parse URI: %v", err)
			}

			result, err := pklresReader.Read(*parsedURL)
			if err != nil {
				t.Fatalf("Failed to read from pklres reader: %v", err)
			}

			resultStr := string(result)
			if tc.contains {
				if !strings.Contains(resultStr, tc.expected) {
					t.Errorf("Expected result to contain '%s', got: %s", tc.expected, resultStr)
				}
			} else {
				if resultStr != tc.expected {
					t.Errorf("Expected '%s', got: %s", tc.expected, resultStr)
				}
			}
		})
	}
}

// TestRealResourceReadersIntegration tests both readers together by calling actual PKL functions
func TestRealResourceReadersIntegration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test agents directory structure
	agentsDir := filepath.Join(tempDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("Failed to create agents directory: %v", err)
	}

	// Create a test agent with workflow.pkl
	testAgentDir := filepath.Join(agentsDir, "test-agent", "1.0.0")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("Failed to create test agent directory: %v", err)
	}

	// Create a simple workflow.pkl file
	workflowContent := `
@ModuleInfo { minPklVersion = "0.28.2" }

module test.Workflow

ActionID = "test-action"
Description = "Test action for integration testing"
`
	if err := os.WriteFile(filepath.Join(testAgentDir, "workflow.pkl"), []byte(workflowContent), 0644); err != nil {
		t.Fatalf("Failed to write workflow.pkl: %v", err)
	}

	// Set up deps/pkl directory with required PKL schema files
	depsDir := filepath.Join(tempDir, "deps", "pkl")
	if err := os.MkdirAll(depsDir, 0755); err != nil {
		t.Fatalf("Failed to create deps/pkl directory: %v", err)
	}
	// Copy Agent.pkl and PklResource.pkl
	schemaRoot := filepath.Clean(filepath.Join("..", "deps", "pkl"))
	requiredFiles := []string{"Agent.pkl", "PklResource.pkl"}
	for _, fname := range requiredFiles {
		src := filepath.Join(schemaRoot, fname)
		dst := filepath.Join(depsDir, fname)
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", dst, err)
		}
	}

	// Create temporary database file
	tempDB, err := os.CreateTemp("", "pklres-integration-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp database: %v", err)
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	// Initialize real readers
	logger := logging.NewTestLogger()
	fs := afero.NewOsFs()

	agentReader, err := agent.InitializeAgent(fs, tempDir, "test-agent", "1.0.0", logger)
	if err != nil {
		t.Fatalf("Failed to initialize agent reader: %v", err)
	}
	defer agentReader.Close()

	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		t.Fatalf("Failed to initialize pklres reader: %v", err)
	}
	// Note: pklresReader doesn't have a Close method, database will be cleaned up when temp file is removed

	// Create evaluator with real readers
	evaluator, err := NewTestEvaluator(agentReader, pklresReader)
	if err != nil {
		t.Fatalf("Failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	// Test PKL functions that call the resource readers
	testCases := []struct {
		name     string
		pklExpr  string
		expected string
		contains bool
	}{
		{
			name: "Test Agent.resolveActionID function",
			pklExpr: `
				import "deps/pkl/Agent.pkl" as agent
				result = agent.resolveActionID("test-action")
			`,
			expected: "@test-agent/test-action:1.0.0",
			contains: true,
		},
		{
			name: "Test PklResource.setPklValue and getPklValue functions",
			pklExpr: `
				import "deps/pkl/PklResource.pkl" as pklres
				setResult = pklres.setPklValue("test-id", "exec", "command", "echo hello")
				getResult = pklres.getPklValue("test-id", "exec", "command")
				result = getResult
			`,
			expected: "echo hello",
			contains: false,
		},
		{
			name: "Test PklResource.getPklRecord function",
			pklExpr: `
				import "deps/pkl/PklResource.pkl" as pklres
				setResult = pklres.setPklValue("test-id2", "python", "script", "print('world')")
				getResult = pklres.getPklRecord("test-id2", "python")
				result = getResult
			`,
			expected: "print('world')",
			contains: false,
		},
		{
			name: "Test Agent.resolveActionID with null input",
			pklExpr: `
				import "deps/pkl/Agent.pkl" as agent
				result = agent.resolveActionID(null)
			`,
			expected: "",
			contains: false,
		},
		{
			name: "Test Agent.resolveActionID with empty string",
			pklExpr: `
				import "deps/pkl/Agent.pkl" as agent
				result = agent.resolveActionID("")
			`,
			expected: "",
			contains: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write PKL code directly into tempDir
			pklFile := filepath.Join(tempDir, fmt.Sprintf("test_%s.pkl", strings.ReplaceAll(tc.name, " ", "_")))
			if err := os.WriteFile(pklFile, []byte(tc.pklExpr), 0644); err != nil {
				t.Fatalf("Failed to write PKL test file: %v", err)
			}
			defer os.Remove(pklFile)

			// Change working directory to tempDir for PKL import resolution
			origWD, _ := os.Getwd()
			if err := os.Chdir(tempDir); err != nil {
				t.Fatalf("Failed to chdir to tempDir: %v", err)
			}
			defer os.Chdir(origWD)

			source := pkl.FileSource(pklFile)
			var module map[string]interface{}
			if err := evaluator.EvaluateModule(context.Background(), source, &module); err != nil {
				t.Fatalf("Failed to evaluate PKL module: %v", err)
			}

			resultStr := fmt.Sprintf("%v", module["result"])
			if tc.contains {
				if !strings.Contains(resultStr, tc.expected) {
					t.Errorf("Expected result to contain '%s', got: %s", tc.expected, resultStr)
				}
			} else {
				if resultStr != tc.expected {
					t.Errorf("Expected '%s', got: %s", tc.expected, resultStr)
				}
			}
		})
	}
}

// TestAgentReaderDatabaseOperations tests database operations of the agent reader
func TestAgentReaderDatabaseOperations(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "agent-db-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test agents directory structure
	agentsDir := filepath.Join(tempDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("Failed to create agents directory: %v", err)
	}

	// Create multiple test agents
	testAgents := []struct {
		name    string
		version string
		actions []string
	}{
		{
			name:    "agent1",
			version: "1.0.0",
			actions: []string{"action1", "action2"},
		},
		{
			name:    "agent1",
			version: "2.0.0",
			actions: []string{"action1", "action3"},
		},
		{
			name:    "agent2",
			version: "1.0.0",
			actions: []string{"action1"},
		},
	}

	for _, agent := range testAgents {
		agentDir := filepath.Join(agentsDir, agent.name, agent.version)
		if err := os.MkdirAll(agentDir, 0755); err != nil {
			t.Fatalf("Failed to create agent directory: %v", err)
		}

		// Create workflow.pkl with actions
		workflowContent := fmt.Sprintf(`
@ModuleInfo { minPklVersion = "0.28.2" }

module test.Workflow

%s
`, func() string {
			var actions []string
			for _, action := range agent.actions {
				actions = append(actions, fmt.Sprintf(`ActionID = "%s"`, action))
			}
			return strings.Join(actions, "\n")
		}())

		if err := os.WriteFile(filepath.Join(agentDir, "workflow.pkl"), []byte(workflowContent), 0644); err != nil {
			t.Fatalf("Failed to write workflow.pkl: %v", err)
		}
	}

	// Initialize agent reader
	logger := logging.NewTestLogger()
	fs := afero.NewOsFs()

	agentReader, err := agent.InitializeAgent(fs, tempDir, "agent1", "1.0.0", logger)
	if err != nil {
		t.Fatalf("Failed to initialize agent reader: %v", err)
	}
	defer agentReader.Close()

	// Test latest version resolution
	testCases := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "Resolve agent1 without version (should get latest)",
			uri:      "agent:/@agent1",
			expected: "@agent1:2.0.0",
		},
		{
			name:     "Resolve agent1 action without version (should get latest)",
			uri:      "agent:/@agent1/action1",
			expected: "@agent1/action1:2.0.0",
		},
		{
			name:     "Resolve specific version",
			uri:      "agent:/@agent1:1.0.0",
			expected: "@agent1:1.0.0",
		},
		{
			name:     "Resolve specific version action",
			uri:      "agent:/@agent1/action1:1.0.0",
			expected: "@agent1/action1:1.0.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedURL, err := url.Parse(tc.uri)
			if err != nil {
				t.Fatalf("Failed to parse URI: %v", err)
			}

			result, err := agentReader.Read(*parsedURL)
			if err != nil {
				t.Fatalf("Failed to read from agent reader: %v", err)
			}

			resultStr := string(result)
			if resultStr != tc.expected {
				t.Errorf("Expected '%s', got: %s", tc.expected, resultStr)
			}
		})
	}

	// Test listing installed agents
	t.Run("List installed agents", func(t *testing.T) {
		parsedURL, err := url.Parse("agent:/?op=list-installed")
		if err != nil {
			t.Fatalf("Failed to parse URI: %v", err)
		}

		result, err := agentReader.Read(*parsedURL)
		if err != nil {
			t.Fatalf("Failed to read from agent reader: %v", err)
		}

		var agents []agent.AgentInfo
		if err := json.Unmarshal(result, &agents); err != nil {
			t.Fatalf("Failed to unmarshal agents: %v", err)
		}

		expectedCount := 3 // agent1:1.0.0, agent1:2.0.0, agent2:1.0.0
		if len(agents) != expectedCount {
			t.Errorf("Expected %d agents, got %d", expectedCount, len(agents))
		}

		// Check that we have the expected agents
		agentMap := make(map[string]bool)
		for _, a := range agents {
			agentMap[fmt.Sprintf("%s:%s", a.Name, a.Version)] = true
		}

		expectedAgents := []string{"agent1:1.0.0", "agent1:2.0.0", "agent2:1.0.0"}
		for _, expected := range expectedAgents {
			if !agentMap[expected] {
				t.Errorf("Expected agent %s not found", expected)
			}
		}
	})
}

// TestPklresReaderConcurrency tests concurrent operations on the pklres reader
func TestPklresReaderConcurrency(t *testing.T) {
	// Create a temporary database file
	tempDB, err := os.CreateTemp("", "pklres-concurrency-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp database: %v", err)
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	// Initialize the real pklres reader
	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		t.Fatalf("Failed to initialize pklres reader: %v", err)
	}
	// Note: pklresReader doesn't have a Close method, database will be cleaned up when temp file is removed

	// Test concurrent set operations
	t.Run("Concurrent set operations", func(t *testing.T) {
		const numGoroutines = 10
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				uri := fmt.Sprintf("pklres:/test-id-%d?op=set&type=exec&key=command&value=echo hello-%d", id, id)
				parsedURL, err := url.Parse(uri)
				if err != nil {
					t.Errorf("Failed to parse URI: %v", err)
					return
				}

				result, err := pklresReader.Read(*parsedURL)
				if err != nil {
					t.Errorf("Failed to set record: %v", err)
					return
				}

				expected := fmt.Sprintf("echo hello-%d", id)
				if string(result) != expected {
					t.Errorf("Expected '%s', got '%s'", expected, string(result))
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})

	// Test concurrent get operations
	t.Run("Concurrent get operations", func(t *testing.T) {
		const numGoroutines = 10
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				uri := fmt.Sprintf("pklres:/test-id-%d?op=get&type=exec&key=command", id)
				parsedURL, err := url.Parse(uri)
				if err != nil {
					t.Errorf("Failed to parse URI: %v", err)
					return
				}

				result, err := pklresReader.Read(*parsedURL)
				if err != nil {
					t.Errorf("Failed to get record: %v", err)
					return
				}

				expected := fmt.Sprintf("echo hello-%d", id)
				if string(result) != expected {
					t.Errorf("Expected '%s', got '%s'", expected, string(result))
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}

// TestResourceReaderErrorHandling tests error handling in both readers
func TestResourceReaderErrorHandling(t *testing.T) {
	// Test agent reader error handling
	t.Run("Agent reader error handling", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir, err := os.MkdirTemp("", "agent-error-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		logger := logging.NewTestLogger()
		fs := afero.NewOsFs()

		agentReader, err := agent.InitializeAgent(fs, tempDir, "", "", logger)
		if err != nil {
			t.Fatalf("Failed to initialize agent reader: %v", err)
		}
		defer agentReader.Close()

		// Test invalid URIs
		invalidURIs := []string{
			"agent:/",               // No action ID
			"agent:/@",              // Empty agent ID
			"agent:/@/action",       // Empty agent ID with action
			"agent:/invalid-format", // Invalid format without @
		}

		for _, uri := range invalidURIs {
			t.Run(fmt.Sprintf("Invalid URI: %s", uri), func(t *testing.T) {
				parsedURL, err := url.Parse(uri)
				if err != nil {
					t.Fatalf("Failed to parse URI: %v", err)
				}

				_, err = agentReader.Read(*parsedURL)
				if err == nil {
					t.Errorf("Expected error for invalid URI: %s", uri)
				}
			})
		}
	})

	// Test pklres reader error handling
	t.Run("Pklres reader error handling", func(t *testing.T) {
		// Create a temporary database file
		tempDB, err := os.CreateTemp("", "pklres-error-*.db")
		if err != nil {
			t.Fatalf("Failed to create temp database: %v", err)
		}
		defer os.Remove(tempDB.Name())
		tempDB.Close()

		pklresReader, err := pklres.InitializePklResource(tempDB.Name())
		if err != nil {
			t.Fatalf("Failed to initialize pklres reader: %v", err)
		}
		// Note: pklresReader doesn't have a Close method, database will be cleaned up when temp file is removed

		// Test invalid operations
		invalidURIs := []string{
			"pklres:/?op=invalid",               // Invalid operation
			"pklres:/?op=set",                   // Missing required parameters
			"pklres:/?op=set&type=exec",         // Missing id and value
			"pklres:/?op=set&id=test&type=exec", // Missing value
			"pklres:/?op=get",                   // Missing required parameters
			"pklres:/?op=get&type=exec",         // Missing id
			"pklres:/?op=delete",                // Missing required parameters
			"pklres:/?op=clear",                 // Missing type
			"pklres:/?op=list",                  // Missing type
		}

		for _, uri := range invalidURIs {
			t.Run(fmt.Sprintf("Invalid URI: %s", uri), func(t *testing.T) {
				parsedURL, err := url.Parse(uri)
				if err != nil {
					t.Fatalf("Failed to parse URI: %v", err)
				}

				_, err = pklresReader.Read(*parsedURL)
				if err == nil {
					t.Errorf("Expected error for invalid URI: %s", uri)
				}
			})
		}
	})
}

func TestSessionReader_Base64AndPlain(t *testing.T) {
	base64Val := base64.StdEncoding.EncodeToString([]byte("Hello Session"))
	plainVal := "Hello Session"
	// Simulate getRecord logic
	decode := func(val string) string {
		decoded, err := base64.StdEncoding.DecodeString(val)
		if err == nil {
			return string(decoded)
		}
		return val
	}
	if got := decode(base64Val); got != "Hello Session" {
		t.Errorf("Session base64 decode failed: got %q", got)
	}
	if got := decode(plainVal); got != plainVal {
		t.Errorf("Session plain decode failed: got %q", got)
	}
}

func TestToolReader_Base64AndPlain(t *testing.T) {
	base64Val := base64.StdEncoding.EncodeToString([]byte("Hello Tool"))
	plainVal := "Hello Tool"
	decode := func(val string) string {
		decoded, err := base64.StdEncoding.DecodeString(val)
		if err == nil {
			return string(decoded)
		}
		return val
	}
	if got := decode(base64Val); got != "Hello Tool" {
		t.Errorf("Tool base64 decode failed: got %q", got)
	}
	if got := decode(plainVal); got != plainVal {
		t.Errorf("Tool plain decode failed: got %q", got)
	}
}

func TestMemoryReader_Base64AndPlain(t *testing.T) {
	base64Val := base64.StdEncoding.EncodeToString([]byte("Hello Memory"))
	plainVal := "Hello Memory"
	decode := func(val string) string {
		decoded, err := base64.StdEncoding.DecodeString(val)
		if err == nil {
			return string(decoded)
		}
		return val
	}
	if got := decode(base64Val); got != "Hello Memory" {
		t.Errorf("Memory base64 decode failed: got %q", got)
	}
	if got := decode(plainVal); got != plainVal {
		t.Errorf("Memory plain decode failed: got %q", got)
	}
}

func TestItemReader_Base64AndPlain(t *testing.T) {
	base64Val := base64.StdEncoding.EncodeToString([]byte("Hello Item"))
	plainVal := "Hello Item"
	decode := func(val string) string {
		decoded, err := base64.StdEncoding.DecodeString(val)
		if err == nil {
			return string(decoded)
		}
		return val
	}
	if got := decode(base64Val); got != "Hello Item" {
		t.Errorf("Item base64 decode failed: got %q", got)
	}
	if got := decode(plainVal); got != plainVal {
		t.Errorf("Item plain decode failed: got %q", got)
	}
}
