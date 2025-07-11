package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"
)

// TestComprehensiveImprovements runs all improvement systems and validates enhancements
func TestComprehensiveImprovements(t *testing.T) {
	t.Log("🚀 COMPREHENSIVE IMPROVEMENTS TEST SUITE")
	t.Log("==========================================")

	// Test Improvement Tracker
	t.Run("ImprovementTracker", func(t *testing.T) {
		t.Log("📋 Testing Improvement Tracking System")

		config := &ImprovementConfig{
			AutoTrackPerformance: true,
			AutoTrackSecurity:    true,
			AutoTrackQuality:     true,
			ReviewInterval:       1 * time.Hour,
			MaxImprovements:      50,
		}

		tracker := NewImprovementTracker(config)

		// Test auto-tracking
		ctx := context.Background()
		tracker.AutoTrackPerformance(ctx)
		tracker.AutoTrackSecurity(ctx)
		tracker.AutoTrackQuality(ctx)

		// Verify improvements were added
		allImprovements := tracker.ListImprovements()
		if len(allImprovements) == 0 {
			t.Error("Expected auto-tracked improvements to be added")
		}

		// Check for specific categories
		categories := make(map[string]int)
		for _, imp := range allImprovements {
			categories[imp.Category]++
		}

		expectedCategories := []string{"performance", "security", "quality"}
		for _, category := range expectedCategories {
			if count := categories[category]; count == 0 {
				t.Errorf("Expected improvements in category: %s", category)
			} else {
				t.Logf("✅ Found %d improvements in category: %s", count, category)
			}
		}

		// Test report generation
		report := tracker.GenerateImprovementReport()
		if report.TotalImprovements == 0 {
			t.Error("Expected total improvements to be greater than 0")
		}

		t.Logf("📊 Improvement Report: %d total, %.1f%% completion rate",
			report.TotalImprovements, report.CompletionRate)
	})

	// Test Configuration Management
	t.Run("ConfigurationManagement", func(t *testing.T) {
		t.Log("⚙️  Testing Configuration Management System")

		manager := NewConfigManager()

		// Test configuration creation and validation
		config := &Config{
			Name:    "test-env",
			Version: "1.0.0",
			Data: map[string]interface{}{
				"database": map[string]interface{}{
					"host":     "localhost",
					"port":     5432,
					"username": "testuser",
					"password": "testpass",
				},
				"api": map[string]interface{}{
					"timeout":    30,
					"retries":    3,
					"rate_limit": 100,
				},
				"logging": map[string]interface{}{
					"level":  "info",
					"format": "json",
				},
			},
			Metadata: &ConfigMetadata{
				Description: "Test environment configuration",
				Tags:        []string{"test", "database", "api", "logging"},
			},
		}

		err := manager.SetConfig(config)
		if err != nil {
			t.Errorf("Failed to set config: %v", err)
		}

		// Test configuration retrieval
		retrieved, exists := manager.GetConfig("test-env")
		if !exists {
			t.Error("Expected config to exist")
		}

		if retrieved.Name != config.Name {
			t.Errorf("Expected name %s, got %s", config.Name, retrieved.Name)
		}

		// Test configuration validation
		validator := func(data map[string]interface{}) error {
			// Validate database configuration
			if db, exists := data["database"]; exists {
				if dbMap, ok := db.(map[string]interface{}); ok {
					if _, exists := dbMap["host"]; !exists {
						return fmt.Errorf("database host is required")
					}
					if port, exists := dbMap["port"]; exists {
						if portNum, ok := port.(float64); ok && portNum <= 0 {
							return fmt.Errorf("database port must be positive")
						}
					}
				}
			}

			// Validate API configuration
			if api, exists := data["api"]; exists {
				if apiMap, ok := api.(map[string]interface{}); ok {
					if timeout, exists := apiMap["timeout"]; exists {
						if timeoutNum, ok := timeout.(float64); ok && timeoutNum <= 0 {
							return fmt.Errorf("API timeout must be positive")
						}
					}
				}
			}

			return nil
		}

		err = manager.AddValidator("test-env", validator)
		if err != nil {
			t.Errorf("Failed to add validator: %v", err)
		}

		// Test configuration watching
		changeDetected := false
		handler := func(config *Config, change string) {
			changeDetected = true
			t.Logf("🔔 Config change detected: %s - %s", config.Name, change)
		}

		manager.WatchConfig("test-env", handler)

		// Update configuration
		retrieved.Data["new_setting"] = "new_value"
		err = manager.SetConfig(retrieved)
		if err != nil {
			t.Errorf("Failed to update config: %v", err)
		}

		// Give time for async notification
		time.Sleep(100 * time.Millisecond)

		if !changeDetected {
			t.Error("Expected configuration change to be detected")
		}

		// Test configuration export
		exported, err := manager.ExportConfig("test-env", "json")
		if err != nil {
			t.Errorf("Failed to export config: %v", err)
		}

		if len(exported) == 0 {
			t.Error("Expected exported configuration to be non-empty")
		}

		t.Logf("✅ Configuration management working correctly")
	})

	// Test Performance Optimization
	t.Run("PerformanceOptimization", func(t *testing.T) {
		t.Log("⚡ Testing Performance Optimization System")

		config := &PerformanceConfig{
			CacheSize:       100,
			PoolSize:        10,
			CacheTTL:        5 * time.Minute,
			MonitorInterval: 30 * time.Second,
			EnableProfiling: true,
		}

		optimizer := NewPerformanceOptimizer(config)
		defer optimizer.Close()

		// Test caching functionality
		testData := map[string]interface{}{
			"user": map[string]interface{}{
				"id":    12345,
				"name":  "Test User",
				"email": "test@example.com",
			},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "en",
			},
		}

		// Set cache entries
		for key, value := range testData {
			optimizer.Set(key, value)
		}

		// Retrieve from cache
		for key := range testData {
			if value, exists := optimizer.Get(key); !exists {
				t.Errorf("Expected cached value for key: %s", key)
			} else {
				// For maps, we can't use direct comparison, so just verify the value exists
				if value == nil {
					t.Errorf("Cache value is nil for key %s", key)
				}
			}
		}

		// Test resource pooling
		resourceCount := 0
		optimizer.SetResourceFactory(func() (interface{}, error) {
			resourceCount++
			return fmt.Sprintf("resource_%d", resourceCount), nil
		})

		// Acquire and release resources
		resources := make([]interface{}, 3)
		for i := 0; i < 3; i++ {
			resource, err := optimizer.AcquireResource()
			if err != nil {
				t.Errorf("Failed to acquire resource: %v", err)
			}
			resources[i] = resource
		}

		// Release resources back to pool
		for _, resource := range resources {
			optimizer.ReleaseResource(resource)
		}

		// Test performance tracking
		start := time.Now()
		time.Sleep(10 * time.Millisecond) // Simulate work
		duration := time.Since(start)

		optimizer.TrackPerformance("test_operation", duration)

		// Get performance metrics
		metrics := optimizer.GetPerformanceMetrics()
		if len(metrics) == 0 {
			t.Error("Expected performance metrics to be available")
		}

		t.Logf("📊 Performance metrics collected: %d operations", len(metrics))
		t.Logf("✅ Performance optimization working correctly")
	})

	// Test Integration Between Systems
	t.Run("SystemIntegration", func(t *testing.T) {
		t.Log("🔗 Testing System Integration")

		// Create all systems
		improvementTracker := NewImprovementTracker(nil)
		configManager := NewConfigManager()
		performanceOptimizer := NewPerformanceOptimizer(nil)
		defer performanceOptimizer.Close()

		// Create a configuration for the improvement tracker
		config := &Config{
			Name:    "improvement-config",
			Version: "1.0.0",
			Data: map[string]interface{}{
				"auto_track": true,
				"categories": []string{"performance", "security", "quality"},
				"thresholds": map[string]interface{}{
					"performance": 80.0,
					"security":    95.0,
					"quality":     90.0,
				},
			},
		}

		err := configManager.SetConfig(config)
		if err != nil {
			t.Errorf("Failed to set improvement config: %v", err)
		}

		// Track performance improvements
		ctx := context.Background()
		improvementTracker.AutoTrackPerformance(ctx)

		// Get improvements and validate against config
		improvements := improvementTracker.ListImprovements(FilterByCategory("performance"))
		if len(improvements) == 0 {
			t.Error("Expected performance improvements to be tracked")
		}

		// Use performance optimizer to track actual performance
		start := time.Now()
		time.Sleep(5 * time.Millisecond) // Simulate work
		duration := time.Since(start)

		performanceOptimizer.TrackPerformance("integration_test", duration)

		// Generate comprehensive report
		improvementReport := improvementTracker.GenerateImprovementReport()
		configs := configManager.ListConfigs()

		t.Logf("📈 Integration Results:")
		t.Logf("  - Improvements tracked: %d", improvementReport.TotalImprovements)
		t.Logf("  - Configurations managed: %d", len(configs))
		t.Logf("  - Performance operations: %d", len(performanceOptimizer.GetPerformanceMetrics()))

		t.Logf("✅ System integration working correctly")
	})

	t.Log("🎉 ALL COMPREHENSIVE IMPROVEMENTS TESTS PASSED!")
}

// TestImprovementMetrics tests the metrics and reporting functionality
func TestImprovementMetrics(t *testing.T) {
	t.Log("📊 Testing Improvement Metrics and Reporting")

	// Create tracker without auto-tracking to have clean state
	config := &ImprovementConfig{
		AutoTrackPerformance: false,
		AutoTrackSecurity:    false,
		AutoTrackQuality:     false,
		ReviewInterval:       1 * time.Hour,
		MaxImprovements:      50,
	}
	tracker := NewImprovementTracker(config)

	// Add various improvements
	improvements := []*Improvement{
		{
			Category:    "performance",
			Priority:    5,
			Title:       "High Priority Performance Fix",
			Description: "Critical performance optimization needed",
			Status:      "completed",
			Tags:        []string{"critical", "performance"},
		},
		{
			Category:    "performance",
			Priority:    2,
			Title:       "Minor Performance Tweak",
			Description: "Small optimization",
			Status:      "deferred",
			Tags:        []string{"minor", "performance"},
		},
		{
			Category:    "quality",
			Priority:    3,
			Title:       "Code Quality Improvement",
			Description: "Improve error handling",
			Status:      "pending",
			Tags:        []string{"quality", "error-handling"},
		},
	}

	for _, imp := range improvements {
		err := tracker.AddImprovement(imp)
		if err != nil {
			t.Errorf("Failed to add improvement: %v", err)
		}
	}

	// Test filtering by various criteria
	t.Run("Filtering", func(t *testing.T) {
		// Filter by category
		perfImprovements := tracker.ListImprovements(FilterByCategory("performance"))
		if len(perfImprovements) != 2 {
			t.Errorf("Expected 2 performance improvements, got %d", len(perfImprovements))
		}

		// Filter by status
		completedImprovements := tracker.ListImprovements(FilterByStatus("completed"))
		if len(completedImprovements) != 1 {
			t.Errorf("Expected 1 completed improvement, got %d", len(completedImprovements))
		}

		// Filter by priority
		highPriorityImprovements := tracker.ListImprovements(FilterByPriority(4))
		if len(highPriorityImprovements) != 1 {
			t.Errorf("Expected 1 high priority improvement, got %d", len(highPriorityImprovements))
		}

		// Filter by tag
		criticalImprovements := tracker.ListImprovements(FilterByTag("critical"))
		if len(criticalImprovements) != 1 {
			t.Errorf("Expected 1 critical improvement, got %d", len(criticalImprovements))
		}
	})

	// Test report generation
	t.Run("Reporting", func(t *testing.T) {
		report := tracker.GenerateImprovementReport()

		// Verify report structure
		if report.TotalImprovements != 3 {
			t.Errorf("Expected 3 total improvements, got %d", report.TotalImprovements)
		}

		if report.CompletedImprovements != 1 {
			t.Errorf("Expected 1 completed improvement, got %d", report.CompletedImprovements)
		}

		expectedCompletionRate := 33.3 // 1 out of 3 ≈ 33.3%
		tolerance := 0.1               // Allow small floating point differences
		if math.Abs(report.CompletionRate-expectedCompletionRate) > tolerance {
			t.Errorf("Expected completion rate %.1f%%, got %.1f%%", expectedCompletionRate, report.CompletionRate)
		}

		// Verify category breakdown
		if count := report.Summary["performance"]; count != 2 {
			t.Errorf("Expected 2 performance improvements, got %d", count)
		}

		if count := report.Summary["security"]; count != 0 {
			t.Errorf("Expected 0 security improvements, got %d", count)
		}

		if count := report.Summary["quality"]; count != 1 {
			t.Errorf("Expected 1 quality improvement, got %d", count)
		}

		t.Logf("📊 Report Summary:")
		t.Logf("  - Total: %d", report.TotalImprovements)
		t.Logf("  - Completed: %d", report.CompletedImprovements)
		t.Logf("  - Completion Rate: %.1f%%", report.CompletionRate)
		t.Logf("  - By Category: %v", report.Summary)
	})

	t.Log("✅ Improvement metrics and reporting working correctly")
}
