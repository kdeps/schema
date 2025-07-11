package test

import (
	"sync"
	"testing"
	"time"
)

func TestMonitoringSystem(t *testing.T) {
	t.Run("NewMonitoringSystem", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		if monitoring == nil {
			t.Fatal("Expected monitoring system to be created")
		}

		alerts := monitoring.GetAlerts()
		if len(alerts) != 0 {
			t.Error("Expected no alerts initially")
		}
	})

	t.Run("CreateAlert", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create an alert
		monitoring.CreateAlert("test_alert", "warning", "Test alert message")

		// Get alerts and verify
		alerts := monitoring.GetAlerts()
		if len(alerts) != 1 {
			t.Fatal("Expected 1 alert to be created")
		}

		alert := alerts[0]
		if alert.Type != "test_alert" {
			t.Errorf("Expected alert type to be 'test_alert', got %s", alert.Type)
		}

		if alert.Severity != "warning" {
			t.Errorf("Expected alert severity to be 'warning', got %s", alert.Severity)
		}

		if alert.Message != "Test alert message" {
			t.Errorf("Expected alert message to be 'Test alert message', got %s", alert.Message)
		}

		if alert.Acknowledged {
			t.Error("Expected alert to not be acknowledged initially")
		}

		if alert.Resolved {
			t.Error("Expected alert to not be resolved initially")
		}
	})

	t.Run("AcknowledgeAlert", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create an alert
		monitoring.CreateAlert("test_alert", "warning", "Test alert message")

		// Get the alert ID
		alerts := monitoring.GetAlerts()
		if len(alerts) == 0 {
			t.Fatal("Expected alerts to be created")
		}

		alertID := alerts[0].ID

		// Acknowledge the alert
		monitoring.AcknowledgeAlert(alertID)

		// Verify acknowledgment
		alert := monitoring.GetAlertByID(alertID)
		if alert == nil {
			t.Fatal("Expected alert to be found")
		}

		if !alert.Acknowledged {
			t.Error("Expected alert to be acknowledged")
		}
	})

	t.Run("ResolveAlert", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create an alert
		monitoring.CreateAlert("test_alert", "warning", "Test alert message")

		// Get the alert ID
		alerts := monitoring.GetAlerts()
		if len(alerts) == 0 {
			t.Fatal("Expected alerts to be created")
		}

		alertID := alerts[0].ID

		// Resolve the alert
		monitoring.ResolveAlert(alertID)

		// Verify resolution
		alert := monitoring.GetAlertByID(alertID)
		if alert == nil {
			t.Fatal("Expected alert to be found")
		}

		if !alert.Resolved {
			t.Error("Expected alert to be resolved")
		}
	})

	t.Run("GetAlertByID", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create an alert
		monitoring.CreateAlert("test_alert", "warning", "Test alert message")

		// Get the alert ID
		alerts := monitoring.GetAlerts()
		if len(alerts) == 0 {
			t.Fatal("Expected alerts to be created")
		}

		alertID := alerts[0].ID

		// Get alert by ID
		alert := monitoring.GetAlertByID(alertID)
		if alert == nil {
			t.Fatal("Expected alert to be found")
		}

		if alert.ID != alertID {
			t.Errorf("Expected alert ID to be %s, got %s", alertID, alert.ID)
		}

		// Test non-existent alert
		nonExistentAlert := monitoring.GetAlertByID("non_existent")
		if nonExistentAlert != nil {
			t.Error("Expected non-existent alert to return nil")
		}
	})

	t.Run("GetAlertsByType", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create multiple alerts of different types
		monitoring.CreateAlert("type_a", "warning", "Alert A1")
		monitoring.CreateAlert("type_b", "error", "Alert B1")
		monitoring.CreateAlert("type_a", "info", "Alert A2")
		monitoring.CreateAlert("type_c", "critical", "Alert C1")

		// Get alerts by type
		typeAAlerts := monitoring.GetAlertsByType("type_a")
		if len(typeAAlerts) != 2 {
			t.Errorf("Expected 2 type_a alerts, got %d", len(typeAAlerts))
		}

		typeBAlerts := monitoring.GetAlertsByType("type_b")
		if len(typeBAlerts) != 1 {
			t.Errorf("Expected 1 type_b alert, got %d", len(typeBAlerts))
		}

		// Verify all type_a alerts have correct type
		for _, alert := range typeAAlerts {
			if alert.Type != "type_a" {
				t.Errorf("Expected alert type to be 'type_a', got %s", alert.Type)
			}
		}
	})

	t.Run("GetAlertsBySeverity", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create multiple alerts of different severities
		monitoring.CreateAlert("test", "warning", "Warning alert")
		monitoring.CreateAlert("test", "error", "Error alert")
		monitoring.CreateAlert("test", "warning", "Another warning")
		monitoring.CreateAlert("test", "critical", "Critical alert")

		// Get alerts by severity
		warningAlerts := monitoring.GetAlertsBySeverity("warning")
		if len(warningAlerts) != 2 {
			t.Errorf("Expected 2 warning alerts, got %d", len(warningAlerts))
		}

		errorAlerts := monitoring.GetAlertsBySeverity("error")
		if len(errorAlerts) != 1 {
			t.Errorf("Expected 1 error alert, got %d", len(errorAlerts))
		}

		// Verify all warning alerts have correct severity
		for _, alert := range warningAlerts {
			if alert.Severity != "warning" {
				t.Errorf("Expected alert severity to be 'warning', got %s", alert.Severity)
			}
		}
	})

	t.Run("ClearAlerts", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create multiple alerts
		monitoring.CreateAlert("test1", "warning", "Alert 1")
		monitoring.CreateAlert("test2", "error", "Alert 2")
		monitoring.CreateAlert("test3", "info", "Alert 3")

		// Verify alerts exist
		alerts := monitoring.GetAlerts()
		if len(alerts) != 3 {
			t.Errorf("Expected 3 alerts, got %d", len(alerts))
		}

		// Clear all alerts
		monitoring.ClearAlerts()

		// Verify alerts are cleared
		alerts = monitoring.GetAlerts()
		if len(alerts) != 0 {
			t.Errorf("Expected 0 alerts after clearing, got %d", len(alerts))
		}
	})
}

func TestMonitoringIntegration(t *testing.T) {
	t.Run("IntegrationWithTestSuite", func(t *testing.T) {
		// Create monitoring system
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create test suite
		testSuite := NewTestSuite()

		// Simulate test execution
		testCases := []struct {
			name   string
			status string
			error  string
		}{
			{"Test1", "PASS", ""},
			{"Test2", "FAIL", "Assertion failed"},
			{"Test3", "SKIP", "Skipped due to dependency"},
			{"Test4", "PASS", ""},
		}

		for _, tc := range testCases {
			// Record test result
			result := TestResult{
				Name:     tc.name,
				Status:   tc.status,
				Duration: 1500 * time.Millisecond, // 1.5 seconds
				Error:    nil,
				Message:  tc.error,
			}

			// Add to test suite metrics
			testSuite.metrics.mu.Lock()
			testSuite.metrics.TestResults[tc.name] = result
			testSuite.metrics.TotalTests++
			switch tc.status {
			case "PASS":
				testSuite.metrics.PassedTests++
			case "FAIL":
				testSuite.metrics.FailedTests++
			case "SKIP":
				testSuite.metrics.SkippedTests++
			}
			testSuite.metrics.mu.Unlock()

			// Create alerts for failures
			if tc.status == "FAIL" {
				monitoring.CreateAlert("test_failure", "error",
					"Test failed: "+tc.name+" - "+tc.error)
			}
		}

		// Verify test suite
		metrics := testSuite.GetMetrics()
		if metrics.TotalTests != 4 {
			t.Errorf("Expected 4 total tests, got %d", metrics.TotalTests)
		}

		if metrics.PassedTests != 2 {
			t.Errorf("Expected 2 passed tests, got %d", metrics.PassedTests)
		}

		if metrics.FailedTests != 1 {
			t.Errorf("Expected 1 failed test, got %d", metrics.FailedTests)
		}

		if metrics.SkippedTests != 1 {
			t.Errorf("Expected 1 skipped test, got %d", metrics.SkippedTests)
		}

		// Verify alerts
		alerts := monitoring.GetAlerts()
		if len(alerts) != 1 {
			t.Errorf("Expected 1 alert for test failure, got %d", len(alerts))
		}

		failureAlerts := monitoring.GetAlertsByType("test_failure")
		if len(failureAlerts) != 1 {
			t.Errorf("Expected 1 test_failure alert, got %d", len(failureAlerts))
		}

		alert := failureAlerts[0]
		if alert.Severity != "error" {
			t.Errorf("Expected alert severity to be 'error', got %s", alert.Severity)
		}
	})

	t.Run("ConcurrentAlertCreation", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		const numGoroutines = 10
		const alertsPerGoroutine = 5

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// Start concurrent alert creation
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < alertsPerGoroutine; j++ {
					alertType := "concurrent_alert"
					severity := "info"
					message := "Concurrent alert from goroutine " + string(rune(id+'0'))
					monitoring.CreateAlert(alertType, severity, message)
				}
			}(i)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Verify all alerts were created
		alerts := monitoring.GetAlerts()
		expectedTotal := numGoroutines * alertsPerGoroutine
		if len(alerts) != expectedTotal {
			t.Errorf("Expected %d total alerts, got %d", expectedTotal, len(alerts))
		}

		// Verify all alerts have correct type
		concurrentAlerts := monitoring.GetAlertsByType("concurrent_alert")
		if len(concurrentAlerts) != expectedTotal {
			t.Errorf("Expected %d concurrent_alert alerts, got %d", expectedTotal, len(concurrentAlerts))
		}
	})

	t.Run("AlertLifecycle", func(t *testing.T) {
		monitoring := NewMonitoringSystem()
		defer monitoring.Close()

		// Create an alert
		monitoring.CreateAlert("lifecycle_test", "warning", "Test alert for lifecycle")

		// Get the alert
		alerts := monitoring.GetAlerts()
		if len(alerts) == 0 {
			t.Fatal("Expected alert to be created")
		}

		alertID := alerts[0].ID
		alert := monitoring.GetAlertByID(alertID)

		// Verify initial state
		if alert.Acknowledged {
			t.Error("Expected alert to not be acknowledged initially")
		}

		if alert.Resolved {
			t.Error("Expected alert to not be resolved initially")
		}

		// Acknowledge the alert
		monitoring.AcknowledgeAlert(alertID)
		alert = monitoring.GetAlertByID(alertID)

		if !alert.Acknowledged {
			t.Error("Expected alert to be acknowledged")
		}

		if alert.Resolved {
			t.Error("Expected alert to not be resolved yet")
		}

		// Resolve the alert
		monitoring.ResolveAlert(alertID)
		alert = monitoring.GetAlertByID(alertID)

		if !alert.Acknowledged {
			t.Error("Expected alert to remain acknowledged")
		}

		if !alert.Resolved {
			t.Error("Expected alert to be resolved")
		}
	})
}
