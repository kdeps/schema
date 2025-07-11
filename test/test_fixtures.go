package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestFixture represents a test fixture with data and metadata
type TestFixture struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Data        map[string]interface{} `json:"data"`
	Created     time.Time              `json:"created"`
	Version     string                 `json:"version"`
}

// FixtureManager manages test fixtures and data generation
type FixtureManager struct {
	fixtures map[string]*TestFixture
	baseDir  string
}

// NewFixtureManager creates a new fixture manager
func NewFixtureManager(baseDir string) *FixtureManager {
	return &FixtureManager{
		fixtures: make(map[string]*TestFixture),
		baseDir:  baseDir,
	}
}

// LoadFixtures loads all fixtures from the fixtures directory
func (fm *FixtureManager) LoadFixtures() error {
	fixturesDir := filepath.Join(fm.baseDir, "fixtures")
	if err := os.MkdirAll(fixturesDir, 0755); err != nil {
		return fmt.Errorf("failed to create fixtures directory: %v", err)
	}

	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		return fmt.Errorf("failed to read fixtures directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fixturePath := filepath.Join(fixturesDir, entry.Name())
		data, err := os.ReadFile(fixturePath)
		if err != nil {
			return fmt.Errorf("failed to read fixture %s: %v", entry.Name(), err)
		}

		var fixture TestFixture
		if err := json.Unmarshal(data, &fixture); err != nil {
			return fmt.Errorf("failed to unmarshal fixture %s: %v", entry.Name(), err)
		}

		fm.fixtures[fixture.Name] = &fixture
	}

	return nil
}

// GetFixture retrieves a fixture by name
func (fm *FixtureManager) GetFixture(name string) (*TestFixture, error) {
	fixture, exists := fm.fixtures[name]
	if !exists {
		return nil, fmt.Errorf("fixture '%s' not found", name)
	}
	return fixture, nil
}

// CreateFixture creates a new fixture
func (fm *FixtureManager) CreateFixture(name, description, category string, data map[string]interface{}) error {
	fixture := &TestFixture{
		Name:        name,
		Description: description,
		Category:    category,
		Data:        data,
		Created:     time.Now(),
		Version:     "1.0.0",
	}

	fm.fixtures[name] = fixture
	return fm.SaveFixture(fixture)
}

// SaveFixture saves a fixture to disk
func (fm *FixtureManager) SaveFixture(fixture *TestFixture) error {
	fixturesDir := filepath.Join(fm.baseDir, "fixtures")
	if err := os.MkdirAll(fixturesDir, 0755); err != nil {
		return fmt.Errorf("failed to create fixtures directory: %v", err)
	}

	data, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal fixture: %v", err)
	}

	fixturePath := filepath.Join(fixturesDir, fixture.Name+".json")
	if err := os.WriteFile(fixturePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write fixture: %v", err)
	}

	return nil
}

// ListFixtures returns all available fixtures
func (fm *FixtureManager) ListFixtures() []*TestFixture {
	var fixtures []*TestFixture
	for _, fixture := range fm.fixtures {
		fixtures = append(fixtures, fixture)
	}
	return fixtures
}

// TestDataGenerator generates test data for various scenarios
type TestDataGenerator struct {
	fixtureManager *FixtureManager
}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator(baseDir string) *TestDataGenerator {
	return &TestDataGenerator{
		fixtureManager: NewFixtureManager(baseDir),
	}
}

// GenerateExecResource generates test data for exec resources
func (tdg *TestDataGenerator) GenerateExecResource() map[string]interface{} {
	return map[string]interface{}{
		"command": "echo 'Hello, World!'",
		"args":    []string{"-n", "test"},
		"env": map[string]string{
			"TEST_ENV": "test_value",
			"DEBUG":    "true",
		},
		"timeout":     30,
		"working_dir": "/tmp",
	}
}

// GeneratePythonResource generates test data for Python resources
func (tdg *TestDataGenerator) GeneratePythonResource() map[string]interface{} {
	return map[string]interface{}{
		"code": `
import sys
import json

def main():
    result = {"message": "Hello from Python", "version": sys.version}
    print(json.dumps(result))

if __name__ == "__main__":
    main()
`,
		"requirements":   []string{"requests", "pandas"},
		"python_version": "3.9",
		"timeout":        60,
	}
}

// GenerateLLMResource generates test data for LLM resources
func (tdg *TestDataGenerator) GenerateLLMResource() map[string]interface{} {
	return map[string]interface{}{
		"model":       "gpt-4",
		"prompt":      "Explain the concept of test-driven development in one sentence.",
		"max_tokens":  100,
		"temperature": 0.7,
		"api_key":     "test_api_key",
		"endpoint":    "https://api.openai.com/v1/chat/completions",
	}
}

// GenerateHTTPResource generates test data for HTTP resources
func (tdg *TestDataGenerator) GenerateHTTPResource() map[string]interface{} {
	return map[string]interface{}{
		"url":    "https://api.example.com/data",
		"method": "GET",
		"headers": map[string]string{
			"Authorization": "Bearer test_token",
			"Content-Type":  "application/json",
		},
		"body":    `{"query": "test"}`,
		"timeout": 30,
		"retries": 3,
	}
}

// GenerateDataResource generates test data for data resources
func (tdg *TestDataGenerator) GenerateDataResource() map[string]interface{} {
	return map[string]interface{}{
		"format": "json",
		"data": map[string]interface{}{
			"users": []map[string]interface{}{
				{"id": 1, "name": "Alice", "email": "alice@example.com"},
				{"id": 2, "name": "Bob", "email": "bob@example.com"},
			},
			"metadata": map[string]interface{}{
				"total":   2,
				"created": time.Now().Format(time.RFC3339),
			},
		},
		"schema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"users": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"id":    map[string]interface{}{"type": "integer"},
							"name":  map[string]interface{}{"type": "string"},
							"email": map[string]interface{}{"type": "string"},
						},
					},
				},
			},
		},
	}
}

// GenerateAgentResource generates test data for agent resources
func (tdg *TestDataGenerator) GenerateAgentResource() map[string]interface{} {
	return map[string]interface{}{
		"agent_id": "test-agent-001",
		"action":   "process_data",
		"version":  "1.0.0",
		"parameters": map[string]interface{}{
			"input_file":    "data.json",
			"output_format": "csv",
			"batch_size":    100,
		},
		"config": map[string]interface{}{
			"timeout":  300,
			"retries":  3,
			"parallel": true,
		},
	}
}

// GeneratePklresResource generates test data for pklres resources
func (tdg *TestDataGenerator) GeneratePklresResource() map[string]interface{} {
	return map[string]interface{}{
		"id":    "test-resource-001",
		"type":  "exec",
		"key":   "command",
		"value": "echo 'Hello from pklres'",
		"metadata": map[string]interface{}{
			"created": time.Now().Format(time.RFC3339),
			"version": "1.0.0",
			"tags":    []string{"test", "example"},
		},
	}
}

// GenerateTestScenario generates a complete test scenario
func (tdg *TestDataGenerator) GenerateTestScenario(scenarioType string) map[string]interface{} {
	switch scenarioType {
	case "basic":
		return map[string]interface{}{
			"exec":   tdg.GenerateExecResource(),
			"python": tdg.GeneratePythonResource(),
			"llm":    tdg.GenerateLLMResource(),
		}
	case "advanced":
		return map[string]interface{}{
			"exec":   tdg.GenerateExecResource(),
			"python": tdg.GeneratePythonResource(),
			"llm":    tdg.GenerateLLMResource(),
			"http":   tdg.GenerateHTTPResource(),
			"data":   tdg.GenerateDataResource(),
			"agent":  tdg.GenerateAgentResource(),
			"pklres": tdg.GeneratePklresResource(),
		}
	case "performance":
		return map[string]interface{}{
			"exec": map[string]interface{}{
				"command": "stress-ng --cpu 1 --timeout 10s",
				"timeout": 15,
			},
			"http": map[string]interface{}{
				"url":     "https://httpbin.org/delay/1",
				"method":  "GET",
				"timeout": 5,
			},
		}
	case "error":
		return map[string]interface{}{
			"exec": map[string]interface{}{
				"command": "nonexistent-command",
				"timeout": 5,
			},
			"http": map[string]interface{}{
				"url":     "https://invalid-url-that-does-not-exist.com",
				"method":  "GET",
				"timeout": 5,
			},
		}
	default:
		return tdg.GenerateTestScenario("basic")
	}
}

// CreateDefaultFixtures creates default test fixtures
func (tdg *TestDataGenerator) CreateDefaultFixtures() error {
	fixtures := []struct {
		name        string
		description string
		category    string
		data        map[string]interface{}
	}{
		{
			name:        "basic_exec",
			description: "Basic exec resource for testing",
			category:    "exec",
			data:        tdg.GenerateExecResource(),
		},
		{
			name:        "basic_python",
			description: "Basic Python resource for testing",
			category:    "python",
			data:        tdg.GeneratePythonResource(),
		},
		{
			name:        "basic_llm",
			description: "Basic LLM resource for testing",
			category:    "llm",
			data:        tdg.GenerateLLMResource(),
		},
		{
			name:        "basic_http",
			description: "Basic HTTP resource for testing",
			category:    "http",
			data:        tdg.GenerateHTTPResource(),
		},
		{
			name:        "basic_data",
			description: "Basic data resource for testing",
			category:    "data",
			data:        tdg.GenerateDataResource(),
		},
		{
			name:        "basic_agent",
			description: "Basic agent resource for testing",
			category:    "agent",
			data:        tdg.GenerateAgentResource(),
		},
		{
			name:        "basic_pklres",
			description: "Basic pklres resource for testing",
			category:    "pklres",
			data:        tdg.GeneratePklresResource(),
		},
		{
			name:        "scenario_basic",
			description: "Basic test scenario with common resources",
			category:    "scenario",
			data:        tdg.GenerateTestScenario("basic"),
		},
		{
			name:        "scenario_advanced",
			description: "Advanced test scenario with all resources",
			category:    "scenario",
			data:        tdg.GenerateTestScenario("advanced"),
		},
		{
			name:        "scenario_performance",
			description: "Performance test scenario",
			category:    "scenario",
			data:        tdg.GenerateTestScenario("performance"),
		},
		{
			name:        "scenario_error",
			description: "Error test scenario for negative testing",
			category:    "scenario",
			data:        tdg.GenerateTestScenario("error"),
		},
	}

	for _, fixture := range fixtures {
		if err := tdg.fixtureManager.CreateFixture(fixture.name, fixture.description, fixture.category, fixture.data); err != nil {
			return fmt.Errorf("failed to create fixture %s: %v", fixture.name, err)
		}
	}

	return nil
}

// GetTestData retrieves test data for a specific resource type
func (tdg *TestDataGenerator) GetTestData(resourceType string) (map[string]interface{}, error) {
	switch resourceType {
	case "exec":
		return tdg.GenerateExecResource(), nil
	case "python":
		return tdg.GeneratePythonResource(), nil
	case "llm":
		return tdg.GenerateLLMResource(), nil
	case "http":
		return tdg.GenerateHTTPResource(), nil
	case "data":
		return tdg.GenerateDataResource(), nil
	case "agent":
		return tdg.GenerateAgentResource(), nil
	case "pklres":
		return tdg.GeneratePklresResource(), nil
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

// TestScenarioBuilder builds complex test scenarios
type TestScenarioBuilder struct {
	generator *TestDataGenerator
	scenario  map[string]interface{}
}

// NewTestScenarioBuilder creates a new test scenario builder
func NewTestScenarioBuilder(generator *TestDataGenerator) *TestScenarioBuilder {
	return &TestScenarioBuilder{
		generator: generator,
		scenario:  make(map[string]interface{}),
	}
}

// AddResource adds a resource to the scenario
func (tsb *TestScenarioBuilder) AddResource(resourceType string) *TestScenarioBuilder {
	data, err := tsb.generator.GetTestData(resourceType)
	if err == nil {
		tsb.scenario[resourceType] = data
	}
	return tsb
}

// AddCustomResource adds a custom resource to the scenario
func (tsb *TestScenarioBuilder) AddCustomResource(name string, data map[string]interface{}) *TestScenarioBuilder {
	tsb.scenario[name] = data
	return tsb
}

// SetMetadata sets metadata for the scenario
func (tsb *TestScenarioBuilder) SetMetadata(metadata map[string]interface{}) *TestScenarioBuilder {
	tsb.scenario["metadata"] = metadata
	return tsb
}

// Build returns the built scenario
func (tsb *TestScenarioBuilder) Build() map[string]interface{} {
	return tsb.scenario
}

// SaveScenario saves the scenario as a fixture
func (tsb *TestScenarioBuilder) SaveScenario(name, description string) error {
	return tsb.generator.fixtureManager.CreateFixture(name, description, "custom", tsb.scenario)
}
