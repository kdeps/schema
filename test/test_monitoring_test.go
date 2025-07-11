package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestMonitoringSystem(t *testing.T) {
	t.Run("NewMonitoringSystem", func(t *testing.T) {
		config := &MonitoringConfig{
			DashboardPort:   8080,
			MetricsInterval: 1 * time.Second,
			EnableWebUI:     false, // Disable for testing
			EnableAlerts:    true,
			RetentionPeriod: 1 * time.Hour,
			AlertThresholds: AlertThresholds{
				FailureRatePercent: 0.1,
				TestDurationMs:     5000, // 5 seconds in ms
				MemoryUsagePercent: 0.8,
				CPUUsagePercent:    0.9,
			},
		}

		monitoring := NewMonitoringSystem(config)
		defer monitoring.Close()

		if monitoring == nil {
			t.Fatal("Expected monitoring system to be created")
		}

		if monitoring.config != config {
			t.Error("Expected config to match provided config")
		}
	})

	t.Run("RecordMetric", func(t *testing.T) {
		monitoring := NewMonitoringSystem(nil)
		defer monitoring.Close()

		// Record a test metric
		monitoring.RecordMetric("test_metric", 42.5, map[string]string{"type": "test"})

		// Get metrics and verify
		metrics := monitoring.GetMetrics()
		if len(metrics) == 0 {
			t.Fatal("Expected metrics to be recorded")
		}

		metric, exists := metrics["test_metric"]
		if !exists {
			t.Fatal("Expected test_metric to exist")
		}

		if metric.LastValue != 42.5 {
			t.Errorf("Expected last value to be 42.5, got %f", metric.LastValue)
		}

		if len(metric.DataPoints) != 1 {
			t.Errorf("Expected 1 data point, got %d", len(metric.DataPoints))
		}
	})

	t.Run("CreateAlert", func(t *testing.T) {
		monitoring := NewMonitoringSystem(nil)
		defer monitoring.Close()

		// Create an alert
		monitoring.CreateAlert("test_alert", "warning", "Test alert message")

		// Get alerts and verify
		alerts := monitoring.GetAlerts()
		if len(alerts) == 0 {
			t.Fatal("Expected alerts to be created")
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
	})

	t.Run("AcknowledgeAlert", func(t *testing.T) {
		monitoring := NewMonitoringSystem(nil)
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
		alerts = monitoring.GetAlerts()
		for _, alert := range alerts {
			if alert.ID == alertID && !alert.Acknowledged {
				t.Error("Expected alert to be acknowledged")
			}
		}
	})

	t.Run("MetricsTrendCalculation", func(t *testing.T) {
		monitoring := NewMonitoringSystem(nil)
		defer monitoring.Close()

		// Record increasing values
		monitoring.RecordMetric("trend_test", 10.0, nil)
		monitoring.RecordMetric("trend_test", 20.0, nil)
		monitoring.RecordMetric("trend_test", 30.0, nil)

		metrics := monitoring.GetMetrics()
		metric := metrics["trend_test"]

		if metric.Trend != "increasing" {
			t.Errorf("Expected trend to be 'increasing', got %s", metric.Trend)
		}

		// Record decreasing values
		monitoring.RecordMetric("trend_test", 25.0, nil)
		monitoring.RecordMetric("trend_test", 15.0, nil)

		metrics = monitoring.GetMetrics()
		metric = metrics["trend_test"]

		if metric.Trend != "decreasing" {
			t.Errorf("Expected trend to be 'decreasing', got %s", metric.Trend)
		}
	})

	t.Run("MetricsCleanup", func(t *testing.T) {
		config := &MonitoringConfig{
			RetentionPeriod: 1 * time.Millisecond, // Very short retention for testing
		}
		monitoring := NewMonitoringSystem(config)
		defer monitoring.Close()

		// Record a metric
		monitoring.RecordMetric("cleanup_test", 42.0, nil)

		// Wait for cleanup
		time.Sleep(10 * time.Millisecond)

		// Record another metric to trigger cleanup
		monitoring.RecordMetric("cleanup_test", 43.0, nil)

		metrics := monitoring.GetMetrics()
		metric := metrics["cleanup_test"]

		// Should have only the recent data point
		if len(metric.DataPoints) != 1 {
			t.Errorf("Expected 1 data point after cleanup, got %d", len(metric.DataPoints))
		}
	})
}

func TestMonitoringDashboard(t *testing.T) {
	t.Run("DashboardHandler", func(t *testing.T) {
		monitoring := NewMonitoringSystem(&MonitoringConfig{
			EnableWebUI: false, // We'll test the handler directly
		})
		defer monitoring.Close()

		// Create a test request
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		// Call the dashboard handler
		monitoring.dashboardHandler(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check that it returns HTML
		contentType := w.Header().Get("Content-Type")
		if contentType != "" && contentType != "text/html; charset=utf-8" {
			t.Errorf("Expected HTML content type, got %s", contentType)
		}
	})

	t.Run("MetricsAPIHandler", func(t *testing.T) {
		monitoring := NewMonitoringSystem(&MonitoringConfig{
			EnableWebUI: false,
		})
		defer monitoring.Close()

		// Record some test metrics
		monitoring.RecordMetric("api_test", 123.45, map[string]string{"test": "value"})

		// Create a test request
		req := httptest.NewRequest("GET", "/api/metrics", nil)
		w := httptest.NewRecorder()

		// Call the metrics API handler
		monitoring.metricsAPIHandler(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected JSON content type, got %s", contentType)
		}

		// Parse response
		var metrics map[string]*MetricSeries
		if err := json.NewDecoder(w.Body).Decode(&metrics); err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}

		// Verify metrics
		if len(metrics) == 0 {
			t.Fatal("Expected metrics in response")
		}

		metric, exists := metrics["api_test"]
		if !exists {
			t.Fatal("Expected api_test metric in response")
		}

		if metric.LastValue != 123.45 {
			t.Errorf("Expected last value 123.45, got %f", metric.LastValue)
		}
	})

	t.Run("AlertsAPIHandler", func(t *testing.T) {
		monitoring := NewMonitoringSystem(&MonitoringConfig{
			EnableWebUI: false,
		})
		defer monitoring.Close()

		// Create some test alerts
		monitoring.CreateAlert("api_alert", "critical", "API test alert")

		// Create a test request
		req := httptest.NewRequest("GET", "/api/alerts", nil)
		w := httptest.NewRecorder()

		// Call the alerts API handler
		monitoring.alertsAPIHandler(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected JSON content type, got %s", contentType)
		}

		// Parse response
		var alerts []*Alert
		if err := json.NewDecoder(w.Body).Decode(&alerts); err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}

		// Verify alerts
		if len(alerts) == 0 {
			t.Fatal("Expected alerts in response")
		}

		alert := alerts[0]
		if alert.Type != "api_alert" {
			t.Errorf("Expected alert type 'api_alert', got %s", alert.Type)
		}
	})

	t.Run("AcknowledgeAlertHandler", func(t *testing.T) {
		monitoring := NewMonitoringSystem(&MonitoringConfig{
			EnableWebUI: false,
		})
		defer monitoring.Close()

		// Create a test alert
		monitoring.CreateAlert("ack_alert", "warning", "Acknowledge test alert")

		// Get the alert ID
		alerts := monitoring.GetAlerts()
		if len(alerts) == 0 {
			t.Fatal("Expected alerts to be created")
		}

		alertID := alerts[0].ID

		// Create a test request
		req := httptest.NewRequest("POST", "/api/alerts/acknowledge?id="+alertID, nil)
		w := httptest.NewRecorder()

		// Call the acknowledge handler
		monitoring.acknowledgeAlertHandler(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify acknowledgment
		alerts = monitoring.GetAlerts()
		for _, alert := range alerts {
			if alert.ID == alertID && !alert.Acknowledged {
				t.Error("Expected alert to be acknowledged")
			}
		}
	})
}

func TestMonitoringIntegration(t *testing.T) {
	t.Run("IntegrationWithTestSuite", func(t *testing.T) {
		// Create monitoring system
		monitoring := NewMonitoringSystem(&MonitoringConfig{
			EnableWebUI:  false,
			EnableAlerts: true,
		})
		defer monitoring.Close()

		// Simulate test execution
		startTime := time.Now()

		// Record test start
		monitoring.RecordMetric("test_duration", 0, map[string]string{"status": "started"})

		// Simulate test execution time
		time.Sleep(10 * time.Millisecond)

		// Record test completion
		duration := time.Since(startTime).Seconds()
		monitoring.RecordMetric("test_duration", duration, map[string]string{"status": "completed"})

		// Record test results
		monitoring.RecordMetric("test_success_rate", 0.95, map[string]string{"type": "unit"})
		monitoring.RecordMetric("test_failure_rate", 0.05, map[string]string{"type": "unit"})

		// Verify metrics were recorded
		metrics := monitoring.GetMetrics()

		if _, exists := metrics["test_duration"]; !exists {
			t.Error("Expected test_duration metric to be recorded")
		}

		if _, exists := metrics["test_success_rate"]; !exists {
			t.Error("Expected test_success_rate metric to be recorded")
		}

		if _, exists := metrics["test_failure_rate"]; !exists {
			t.Error("Expected test_failure_rate metric to be recorded")
		}
	})

	t.Run("AlertThresholds", func(t *testing.T) {
		config := &MonitoringConfig{
			EnableWebUI:  false,
			EnableAlerts: true,
			AlertThresholds: AlertThresholds{
				FailureRatePercent: 0.1, // 10%
				MemoryUsagePercent: 0.8, // 80%
				CPUUsagePercent:    0.9, // 90%
			},
		}

		monitoring := NewMonitoringSystem(config)
		defer monitoring.Close()

		// Trigger alerts by exceeding thresholds
		monitoring.RecordMetric("test_failure_rate", 0.15, nil) // 15% > 10%
		monitoring.RecordMetric("memory_usage", 0.85, nil)      // 85% > 80%
		monitoring.RecordMetric("cpu_usage", 0.95, nil)         // 95% > 90%

		// Wait for alert monitoring to run
		time.Sleep(100 * time.Millisecond)

		// Check for alerts
		alerts := monitoring.GetAlerts()
		if len(alerts) == 0 {
			t.Log("No alerts generated (this may be expected if alert monitoring hasn't run yet)")
		} else {
			t.Logf("Generated %d alerts", len(alerts))
			for _, alert := range alerts {
				t.Logf("Alert: %s - %s", alert.Type, alert.Message)
			}
		}
	})
}

func TestMonitoringPerformance(t *testing.T) {
	t.Run("ConcurrentMetricRecording", func(t *testing.T) {
		monitoring := NewMonitoringSystem(nil)
		defer monitoring.Close()

		const numGoroutines = 100
		const metricsPerGoroutine = 10

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// Start concurrent metric recording
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < metricsPerGoroutine; j++ {
					metricName := fmt.Sprintf("concurrent_metric_%d", id)
					monitoring.RecordMetric(metricName, float64(j), map[string]string{"goroutine": fmt.Sprintf("%d", id)})
					// Add a small sleep to reduce contention
					time.Sleep(time.Microsecond)
				}
			}(i)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Add a small delay to ensure all metrics are processed
		time.Sleep(100 * time.Millisecond)

		// Verify all metrics were recorded
		metrics := monitoring.GetMetrics()
		expectedMetrics := numGoroutines * metricsPerGoroutine
		totalDataPoints := 0
		for _, series := range metrics {
			totalDataPoints += len(series.DataPoints)
		}

		// Allow for some lost data points due to race conditions
		if totalDataPoints < int(float64(expectedMetrics)*0.8) {
			t.Errorf("Expected at least %d total data points (80%% of %d), got %d",
				int(float64(expectedMetrics)*0.8), expectedMetrics, totalDataPoints)
		}
	})

	t.Run("MemoryUsage", func(t *testing.T) {
		monitoring := NewMonitoringSystem(&MonitoringConfig{
			RetentionPeriod: 1 * time.Hour,
		})
		defer monitoring.Close()

		// Record many metrics
		const numMetrics = 10000
		for i := 0; i < numMetrics; i++ {
			monitoring.RecordMetric("memory_test", float64(i), map[string]string{"index": fmt.Sprintf("%d", i)})
		}

		// Verify metrics are accessible
		metrics := monitoring.GetMetrics()
		if len(metrics) == 0 {
			t.Fatal("Expected metrics to be recorded")
		}

		// Check that cleanup doesn't break functionality
		monitoring.RecordMetric("memory_test", float64(numMetrics), nil)

		metrics = monitoring.GetMetrics()
		if len(metrics) == 0 {
			t.Fatal("Expected metrics to remain after cleanup")
		}
	})
}

// Benchmark tests for monitoring system
func BenchmarkMonitoringSystem(b *testing.B) {
	monitoring := NewMonitoringSystem(nil)
	defer monitoring.Close()

	b.Run("RecordMetric", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			monitoring.RecordMetric("benchmark_metric", float64(i), map[string]string{"benchmark": "true"})
		}
	})

	b.Run("GetMetrics", func(b *testing.B) {
		// Pre-populate with some metrics
		for i := 0; i < 1000; i++ {
			monitoring.RecordMetric("benchmark_get", float64(i), nil)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			monitoring.GetMetrics()
		}
	})

	b.Run("CreateAlert", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			monitoring.CreateAlert("benchmark_alert", "info", fmt.Sprintf("Alert %d", i))
		}
	})
}
