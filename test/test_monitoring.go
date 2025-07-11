package test

import (
	"context"
	"sync"
	"time"
)

// MonitoringSystem provides basic monitoring capabilities for integration tests
type MonitoringSystem struct {
	alerts []*Alert
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// Alert represents a monitoring alert
type Alert struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
	Acknowledged bool      `json:"acknowledged"`
	Resolved     bool      `json:"resolved"`
}

// NewMonitoringSystem creates a new simplified monitoring system
func NewMonitoringSystem() *MonitoringSystem {
	ctx, cancel := context.WithCancel(context.Background())

	return &MonitoringSystem{
		alerts: make([]*Alert, 0),
		ctx:    ctx,
		cancel: cancel,
	}
}

// CreateAlert creates a new alert
func (ms *MonitoringSystem) CreateAlert(alertType, severity, message string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	alert := &Alert{
		ID:           generateAlertID(),
		Type:         alertType,
		Severity:     severity,
		Message:      message,
		Timestamp:    time.Now(),
		Acknowledged: false,
		Resolved:     false,
	}

	ms.alerts = append(ms.alerts, alert)
}

// GetAlerts returns current alerts
func (ms *MonitoringSystem) GetAlerts() []*Alert {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	alerts := make([]*Alert, len(ms.alerts))
	copy(alerts, ms.alerts)
	return alerts
}

// AcknowledgeAlert acknowledges an alert
func (ms *MonitoringSystem) AcknowledgeAlert(alertID string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, alert := range ms.alerts {
		if alert.ID == alertID {
			alert.Acknowledged = true
			break
		}
	}
}

// ResolveAlert resolves an alert
func (ms *MonitoringSystem) ResolveAlert(alertID string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, alert := range ms.alerts {
		if alert.ID == alertID {
			alert.Resolved = true
			break
		}
	}
}

// GetAlertByID returns a specific alert by ID
func (ms *MonitoringSystem) GetAlertByID(alertID string) *Alert {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for _, alert := range ms.alerts {
		if alert.ID == alertID {
			return alert
		}
	}
	return nil
}

// GetAlertsByType returns alerts filtered by type
func (ms *MonitoringSystem) GetAlertsByType(alertType string) []*Alert {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var filtered []*Alert
	for _, alert := range ms.alerts {
		if alert.Type == alertType {
			filtered = append(filtered, alert)
		}
	}
	return filtered
}

// GetAlertsBySeverity returns alerts filtered by severity
func (ms *MonitoringSystem) GetAlertsBySeverity(severity string) []*Alert {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var filtered []*Alert
	for _, alert := range ms.alerts {
		if alert.Severity == severity {
			filtered = append(filtered, alert)
		}
	}
	return filtered
}

// ClearAlerts removes all alerts
func (ms *MonitoringSystem) ClearAlerts() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.alerts = make([]*Alert, 0)
}

// Close shuts down the monitoring system
func (ms *MonitoringSystem) Close() {
	ms.cancel()
}

// Helper function to generate alert IDs
func generateAlertID() string {
	return "alert_" + time.Now().Format("20060102150405")
}
