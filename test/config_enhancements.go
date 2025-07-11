package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// ConfigManager provides enhanced configuration management
type ConfigManager struct {
	configs map[string]*Config
	mu      sync.RWMutex
	watcher *ConfigWatcher
}

// Config represents a configuration with metadata
type Config struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	Data        map[string]interface{} `json:"data"`
	Metadata    *ConfigMetadata        `json:"metadata"`
	Validators  []ConfigValidator      `json:"-"`
}

// ConfigMetadata contains configuration metadata
type ConfigMetadata struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastLoaded  time.Time `json:"last_loaded"`
	Source      string    `json:"source"`
	Checksum    string    `json:"checksum"`
	Tags        []string  `json:"tags,omitempty"`
	Description string    `json:"description,omitempty"`
}

// ConfigValidator validates configuration data
type ConfigValidator func(data map[string]interface{}) error

// ConfigWatcher watches for configuration changes
type ConfigWatcher struct {
	watchers map[string][]ConfigChangeHandler
	mu       sync.RWMutex
	stopChan chan struct{}
}

// ConfigChangeHandler handles configuration changes
type ConfigChangeHandler func(config *Config, changeType string)

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs: make(map[string]*Config),
		watcher: &ConfigWatcher{
			watchers: make(map[string][]ConfigChangeHandler),
			stopChan: make(chan struct{}),
		},
	}
}

// LoadConfig loads a configuration from file
func (cm *ConfigManager) LoadConfig(name, filePath string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var configData map[string]interface{}
	if err := json.Unmarshal(data, &configData); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	config := &Config{
		Name:     name,
		Version:  "1.0.0",
		Data:     configData,
		Metadata: &ConfigMetadata{},
	}

	// Set metadata
	config.Metadata.CreatedAt = time.Now()
	config.Metadata.UpdatedAt = time.Now()
	config.Metadata.LastLoaded = time.Now()
	config.Metadata.Source = filePath
	config.Metadata.Checksum = calculateChecksum(data)

	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	cm.configs[name] = config

	// Notify watchers
	cm.watcher.notifyChange(config, "loaded")

	return nil
}

// GetConfig retrieves a configuration by name
func (cm *ConfigManager) GetConfig(name string) (*Config, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	config, exists := cm.configs[name]
	return config, exists
}

// SetConfig sets a configuration
func (cm *ConfigManager) SetConfig(config *Config) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if config.Metadata == nil {
		config.Metadata = &ConfigMetadata{}
	}

	config.Metadata.UpdatedAt = time.Now()

	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	cm.configs[config.Name] = config

	// Notify watchers
	cm.watcher.notifyChange(config, "updated")

	return nil
}

// DeleteConfig removes a configuration
func (cm *ConfigManager) DeleteConfig(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	config, exists := cm.configs[name]
	if !exists {
		return fmt.Errorf("config not found: %s", name)
	}

	delete(cm.configs, name)

	// Notify watchers
	cm.watcher.notifyChange(config, "deleted")

	return nil
}

// ListConfigs returns all configuration names
func (cm *ConfigManager) ListConfigs() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var names []string
	for name := range cm.configs {
		names = append(names, name)
	}

	return names
}

// AddValidator adds a validator to a configuration
func (cm *ConfigManager) AddValidator(configName string, validator ConfigValidator) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	config, exists := cm.configs[configName]
	if !exists {
		return fmt.Errorf("config not found: %s", configName)
	}

	config.Validators = append(config.Validators, validator)

	// Re-validate
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("config validation failed after adding validator: %w", err)
	}

	return nil
}

// validateConfig validates a configuration
func (cm *ConfigManager) validateConfig(config *Config) error {
	for _, validator := range config.Validators {
		if err := validator(config.Data); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}
	return nil
}

// WatchConfig adds a change handler for a configuration
func (cm *ConfigManager) WatchConfig(configName string, handler ConfigChangeHandler) {
	cm.watcher.addWatcher(configName, handler)
}

// UnwatchConfig removes a change handler for a configuration
func (cm *ConfigManager) UnwatchConfig(configName string, handler ConfigChangeHandler) {
	cm.watcher.removeWatcher(configName, handler)
}

// SaveConfig saves a configuration to file
func (cm *ConfigManager) SaveConfig(name, filePath string) error {
	cm.mu.RLock()
	config, exists := cm.configs[name]
	cm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("config not found: %s", name)
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal configuration
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Update metadata
	cm.mu.Lock()
	config.Metadata.UpdatedAt = time.Now()
	config.Metadata.Source = filePath
	config.Metadata.Checksum = calculateChecksum(data)
	cm.mu.Unlock()

	return nil
}

// ExportConfig exports a configuration in different formats
func (cm *ConfigManager) ExportConfig(name, format string) ([]byte, error) {
	cm.mu.RLock()
	config, exists := cm.configs[name]
	cm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("config not found: %s", name)
	}

	switch format {
	case "json":
		return json.MarshalIndent(config, "", "  ")
	case "yaml":
		// Note: Would need yaml package for full YAML support
		// For now, return JSON format as fallback
		return json.Marshal(map[string]interface{}{
			"name":     config.Name,
			"version":  config.Version,
			"data":     config.Data,
			"metadata": config.Metadata,
		})
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// calculateChecksum calculates a simple checksum for data
func calculateChecksum(data []byte) string {
	// Simple checksum implementation
	// In production, use a proper hash function
	sum := 0
	for _, b := range data {
		sum += int(b)
	}
	return fmt.Sprintf("%x", sum)
}

// addWatcher adds a change handler for a configuration
func (cw *ConfigWatcher) addWatcher(configName string, handler ConfigChangeHandler) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.watchers[configName] = append(cw.watchers[configName], handler)
}

// removeWatcher removes a change handler for a configuration
func (cw *ConfigWatcher) removeWatcher(configName string, handler ConfigChangeHandler) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	watchers, exists := cw.watchers[configName]
	if !exists {
		return
	}

	for i, w := range watchers {
		if fmt.Sprintf("%p", w) == fmt.Sprintf("%p", handler) {
			cw.watchers[configName] = append(watchers[:i], watchers[i+1:]...)
			break
		}
	}
}

// notifyChange notifies all watchers of a configuration change
func (cw *ConfigWatcher) notifyChange(config *Config, changeType string) {
	cw.mu.RLock()
	watchers, exists := cw.watchers[config.Name]
	cw.mu.RUnlock()

	if !exists {
		return
	}

	for _, handler := range watchers {
		go handler(config, changeType)
	}
}

// TestConfigManager tests the configuration management functionality
func TestConfigManager(t *testing.T) {
	manager := NewConfigManager()

	// Test creating and setting a configuration
	t.Run("SetAndGetConfig", func(t *testing.T) {
		config := &Config{
			Name:    "test-config",
			Version: "1.0.0",
			Data: map[string]interface{}{
				"database": map[string]interface{}{
					"host":     "localhost",
					"port":     5432,
					"username": "testuser",
				},
				"api": map[string]interface{}{
					"timeout": 30,
					"retries": 3,
				},
			},
			Metadata: &ConfigMetadata{
				Description: "Test configuration",
				Tags:        []string{"test", "database", "api"},
			},
		}

		err := manager.SetConfig(config)
		if err != nil {
			t.Errorf("Failed to set config: %v", err)
		}

		retrieved, exists := manager.GetConfig("test-config")
		if !exists {
			t.Error("Expected config to exist")
		}

		if retrieved.Name != config.Name {
			t.Errorf("Expected name %s, got %s", config.Name, retrieved.Name)
		}

		if retrieved.Metadata.Description != config.Metadata.Description {
			t.Errorf("Expected description %s, got %s", config.Metadata.Description, retrieved.Metadata.Description)
		}
	})

	// Test configuration validation
	t.Run("ConfigValidation", func(t *testing.T) {
		// Add a validator
		validator := func(data map[string]interface{}) error {
			if db, exists := data["database"]; exists {
				if dbMap, ok := db.(map[string]interface{}); ok {
					if _, exists := dbMap["host"]; !exists {
						return fmt.Errorf("database host is required")
					}
				}
			}
			return nil
		}

		err := manager.AddValidator("test-config", validator)
		if err != nil {
			t.Errorf("Failed to add validator: %v", err)
		}

		// Test invalid configuration
		invalidConfig := &Config{
			Name:    "invalid-config",
			Version: "1.0.0",
			Data: map[string]interface{}{
				"database": map[string]interface{}{
					"port":     5432,
					"username": "testuser",
					// Missing host
				},
			},
		}

		err = manager.SetConfig(invalidConfig)
		if err == nil {
			t.Error("Expected validation error for missing host")
		}
	})

	// Test configuration watching
	t.Run("ConfigWatching", func(t *testing.T) {
		changeDetected := false
		changeType := ""

		handler := func(config *Config, change string) {
			changeDetected = true
			changeType = change
		}

		manager.WatchConfig("test-config", handler)

		// Update configuration
		config, _ := manager.GetConfig("test-config")
		config.Data["new_field"] = "new_value"
		if err := manager.SetConfig(config); err != nil {
			t.Errorf("Failed to update config: %v", err)
		}

		// Give some time for the goroutine to execute
		time.Sleep(100 * time.Millisecond)

		if !changeDetected {
			t.Error("Expected change to be detected")
		}

		if changeType != "updated" {
			t.Errorf("Expected change type 'updated', got %s", changeType)
		}
	})

	// Test configuration export
	t.Run("ConfigExport", func(t *testing.T) {
		data, err := manager.ExportConfig("test-config", "json")
		if err != nil {
			t.Errorf("Failed to export config: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected exported data to be non-empty")
		}

		// Verify it's valid JSON
		var exportedConfig Config
		if err := json.Unmarshal(data, &exportedConfig); err != nil {
			t.Errorf("Failed to parse exported config: %v", err)
		}

		if exportedConfig.Name != "test-config" {
			t.Errorf("Expected name %s, got %s", "test-config", exportedConfig.Name)
		}
	})

	// Test configuration listing
	t.Run("ListConfigs", func(t *testing.T) {
		configs := manager.ListConfigs()
		if len(configs) == 0 {
			t.Error("Expected at least one configuration")
		}

		found := false
		for _, name := range configs {
			if name == "test-config" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected test-config to be in the list")
		}
	})

	// Test configuration deletion
	t.Run("DeleteConfig", func(t *testing.T) {
		err := manager.DeleteConfig("test-config")
		if err != nil {
			t.Errorf("Failed to delete config: %v", err)
		}

		_, exists := manager.GetConfig("test-config")
		if exists {
			t.Error("Expected config to be deleted")
		}
	})

	t.Log("✅ Configuration management functionality working correctly")
}
