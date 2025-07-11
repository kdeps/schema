package test

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

// MonitoringSystem provides real-time monitoring capabilities
type MonitoringSystem struct {
	metrics   *MetricsCollector
	dashboard *DashboardServer
	alerts    *AlertManager
	config    *MonitoringConfig
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
}

// MonitoringConfig configures the monitoring system
type MonitoringConfig struct {
	DashboardPort   int             `json:"dashboard_port"`
	MetricsInterval time.Duration   `json:"metrics_interval"`
	AlertThresholds AlertThresholds `json:"alert_thresholds"`
	EnableWebUI     bool            `json:"enable_web_ui"`
	EnableAlerts    bool            `json:"enable_alerts"`
	RetentionPeriod time.Duration   `json:"retention_period"`
}

// MetricsCollector collects and stores metrics
type MetricsCollector struct {
	metrics sync.Map
	config  *MonitoringConfig
}

// MetricSeries represents a time series of metrics
type MetricSeries struct {
	Name       string        `json:"name"`
	DataPoints []MetricPoint `json:"data_points"`
	LastValue  float64       `json:"last_value"`
	Trend      string        `json:"trend"`
	mu         sync.Mutex    // protects DataPoints, LastValue, Trend
}

// MetricPoint represents a single metric measurement
type MetricPoint struct {
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
}

// DashboardServer serves the monitoring dashboard
type DashboardServer struct {
	server *http.Server
	config *MonitoringConfig
}

// AlertManager manages alerts and notifications
type AlertManager struct {
	alerts []*Alert
	mu     sync.RWMutex
	config *MonitoringConfig
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

// NewMonitoringSystem creates a new monitoring system
func NewMonitoringSystem(config *MonitoringConfig) *MonitoringSystem {
	if config == nil {
		config = &MonitoringConfig{
			DashboardPort:   8080,
			MetricsInterval: 30 * time.Second,
			EnableWebUI:     true,
			EnableAlerts:    true,
			RetentionPeriod: 24 * time.Hour,
			AlertThresholds: AlertThresholds{
				FailureRatePercent: 0.1,  // 10%
				TestDurationMs:     5000, // 5 seconds in ms
				MemoryUsagePercent: 0.8,  // 80%
				CPUUsagePercent:    0.9,  // 90%
			},
		}
	}

	// Ensure minimum interval to prevent panic
	if config.MetricsInterval <= 0 {
		config.MetricsInterval = 30 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitoring := &MonitoringSystem{
		metrics: &MetricsCollector{
			metrics: sync.Map{},
			config:  config,
		},
		alerts: &AlertManager{
			alerts: make([]*Alert, 0),
			config: config,
		},
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}

	// Start metrics collection
	go monitoring.startMetricsCollection()

	// Start dashboard if enabled
	if config.EnableWebUI {
		monitoring.dashboard = monitoring.createDashboard()
		go monitoring.startDashboard()
	}

	// Start alert monitoring if enabled
	if config.EnableAlerts {
		go monitoring.startAlertMonitoring()
	}

	return monitoring
}

// RecordMetric records a new metric
func (ms *MonitoringSystem) RecordMetric(name string, value float64, tags map[string]string) {
	seriesIface, _ := ms.metrics.metrics.Load(name)
	if seriesIface == nil {
		series := &MetricSeries{
			Name:       name,
			DataPoints: make([]MetricPoint, 0),
		}
		ms.metrics.metrics.Store(name, series)
	}

	point := MetricPoint{
		Timestamp: time.Now(),
		Value:     value,
		Tags:      tags,
	}

	seriesIface, _ = ms.metrics.metrics.Load(name)
	series := seriesIface.(*MetricSeries)

	series.mu.Lock()
	defer series.mu.Unlock()

	series.DataPoints = append(series.DataPoints, point)
	series.LastValue = value

	// Calculate trend
	if len(series.DataPoints) >= 2 {
		recent := series.DataPoints[len(series.DataPoints)-1].Value
		previous := series.DataPoints[len(series.DataPoints)-2].Value
		if recent > previous {
			series.Trend = "increasing"
		} else if recent < previous {
			series.Trend = "decreasing"
		} else {
			series.Trend = "stable"
		}
	}

	// Cleanup old data points
	ms.cleanupOldMetrics()
}

// GetMetrics returns current metrics
func (ms *MonitoringSystem) GetMetrics() map[string]*MetricSeries {
	metrics := make(map[string]*MetricSeries)
	ms.metrics.metrics.Range(func(key, value interface{}) bool {
		series := value.(*MetricSeries)
		series.mu.Lock()
		metrics[key.(string)] = &MetricSeries{
			Name:       series.Name,
			DataPoints: append([]MetricPoint(nil), series.DataPoints...),
			LastValue:  series.LastValue,
			Trend:      series.Trend,
		}
		series.mu.Unlock()
		return true
	})

	return metrics
}

// CreateAlert creates a new alert
func (ms *MonitoringSystem) CreateAlert(alertType, severity, message string) {
	ms.alerts.mu.Lock()
	defer ms.alerts.mu.Unlock()

	alert := &Alert{
		ID:           fmt.Sprintf("alert_%d", time.Now().Unix()),
		Type:         alertType,
		Severity:     severity,
		Message:      message,
		Timestamp:    time.Now(),
		Acknowledged: false,
		Resolved:     false,
	}

	ms.alerts.alerts = append(ms.alerts.alerts, alert)
}

// GetAlerts returns current alerts
func (ms *MonitoringSystem) GetAlerts() []*Alert {
	ms.alerts.mu.RLock()
	defer ms.alerts.mu.RUnlock()

	alerts := make([]*Alert, len(ms.alerts.alerts))
	copy(alerts, ms.alerts.alerts)
	return alerts
}

// AcknowledgeAlert acknowledges an alert
func (ms *MonitoringSystem) AcknowledgeAlert(alertID string) {
	ms.alerts.mu.Lock()
	defer ms.alerts.mu.Unlock()

	for _, alert := range ms.alerts.alerts {
		if alert.ID == alertID {
			alert.Acknowledged = true
			break
		}
	}
}

// startMetricsCollection starts the metrics collection loop
func (ms *MonitoringSystem) startMetricsCollection() {
	ticker := time.NewTicker(ms.config.MetricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ms.ctx.Done():
			return
		case <-ticker.C:
			ms.collectSystemMetrics()
		}
	}
}

// collectSystemMetrics collects system-level metrics
func (ms *MonitoringSystem) collectSystemMetrics() {
	// Collect memory usage
	memUsage := ms.getMemoryUsage()
	ms.RecordMetric("memory_usage", memUsage, map[string]string{"type": "system"})

	// Collect CPU usage
	cpuUsage := ms.getCPUUsage()
	ms.RecordMetric("cpu_usage", cpuUsage, map[string]string{"type": "system"})

	// Collect test metrics
	testMetrics := ms.getTestMetrics()
	for name, value := range testMetrics {
		ms.RecordMetric(name, value, map[string]string{"type": "test"})
	}
}

// startAlertMonitoring starts the alert monitoring loop
func (ms *MonitoringSystem) startAlertMonitoring() {
	ticker := time.NewTicker(ms.config.MetricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ms.ctx.Done():
			return
		case <-ticker.C:
			ms.checkAlertThresholds()
		}
	}
}

// checkAlertThresholds checks if any metrics exceed alert thresholds
func (ms *MonitoringSystem) checkAlertThresholds() {
	metrics := ms.GetMetrics()

	// Check memory usage
	if memSeries, exists := metrics["memory_usage"]; exists {
		if memSeries.LastValue > ms.config.AlertThresholds.MemoryUsagePercent {
			ms.CreateAlert("memory", "warning",
				fmt.Sprintf("Memory usage is high: %.2f%%", memSeries.LastValue*100))
		}
	}

	// Check CPU usage
	if cpuSeries, exists := metrics["cpu_usage"]; exists {
		if cpuSeries.LastValue > ms.config.AlertThresholds.CPUUsagePercent {
			ms.CreateAlert("cpu", "warning",
				fmt.Sprintf("CPU usage is high: %.2f%%", cpuSeries.LastValue*100))
		}
	}

	// Check test failure rate
	if failureSeries, exists := metrics["test_failure_rate"]; exists {
		if failureSeries.LastValue > ms.config.AlertThresholds.FailureRatePercent {
			ms.CreateAlert("test_failure", "critical",
				fmt.Sprintf("Test failure rate is high: %.2f%%", failureSeries.LastValue*100))
		}
	}
}

// createDashboard creates the dashboard server
func (ms *MonitoringSystem) createDashboard() *DashboardServer {
	mux := http.NewServeMux()

	// Dashboard routes
	mux.HandleFunc("/", ms.dashboardHandler)
	mux.HandleFunc("/api/metrics", ms.metricsAPIHandler)
	mux.HandleFunc("/api/alerts", ms.alertsAPIHandler)
	mux.HandleFunc("/api/alerts/acknowledge", ms.acknowledgeAlertHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", ms.config.DashboardPort),
		Handler: mux,
	}

	return &DashboardServer{
		server: server,
		config: ms.config,
	}
}

// startDashboard starts the dashboard server
func (ms *MonitoringSystem) startDashboard() {
	fmt.Printf("Starting monitoring dashboard on port %d\n", ms.config.DashboardPort)
	if err := ms.dashboard.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Dashboard server error: %v\n", err)
	}
}

// dashboardHandler serves the main dashboard page
func (ms *MonitoringSystem) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Monitoring Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .metric { background: #f5f5f5; padding: 15px; margin: 10px; border-radius: 5px; }
        .alert { background: #fff3cd; padding: 10px; margin: 10px; border-radius: 3px; }
        .critical { background: #f8d7da; }
        .warning { background: #fff3cd; }
        .info { background: #d1ecf1; }
    </style>
    <script>
        function refreshMetrics() {
            fetch('/api/metrics')
                .then(response => response.json())
                .then(data => {
                    const container = document.getElementById('metrics');
                    container.innerHTML = '';
                    Object.keys(data).forEach(key => {
                        const metric = data[key];
                        container.innerHTML += '<div class="metric">' +
                            '<h3>' + metric.name + '</h3>' +
                            '<p>Current Value: ' + metric.last_value + '</p>' +
                            '<p>Trend: ' + metric.trend + '</p>' +
                            '</div>';
                    });
                });
        }
        
        function refreshAlerts() {
            fetch('/api/alerts')
                .then(response => response.json())
                .then(data => {
                    const container = document.getElementById('alerts');
                    container.innerHTML = '';
                    data.forEach(alert => {
                        container.innerHTML += '<div class="alert ' + alert.severity + '">' +
                            '<h4>' + alert.type + ' (' + alert.severity + ')</h4>' +
                            '<p>' + alert.message + '</p>' +
                            '<p>Time: ' + alert.timestamp + '</p>' +
                            '</div>';
                    });
                });
        }
        
        // Refresh every 30 seconds
        setInterval(refreshMetrics, 30000);
        setInterval(refreshAlerts, 30000);
        
        // Initial load
        refreshMetrics();
        refreshAlerts();
    </script>
</head>
<body>
    <h1>Test Monitoring Dashboard</h1>
    
    <h2>Metrics</h2>
    <div id="metrics"></div>
    
    <h2>Alerts</h2>
    <div id="alerts"></div>
</body>
</html>`

	tmplParsed, _ := template.New("dashboard").Parse(tmpl)
	tmplParsed.Execute(w, nil)
}

// metricsAPIHandler serves metrics data as JSON
func (ms *MonitoringSystem) metricsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms.GetMetrics())
}

// alertsAPIHandler serves alerts data as JSON
func (ms *MonitoringSystem) alertsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms.GetAlerts())
}

// acknowledgeAlertHandler handles alert acknowledgment
func (ms *MonitoringSystem) acknowledgeAlertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	alertID := r.URL.Query().Get("id")
	if alertID == "" {
		http.Error(w, "Alert ID required", http.StatusBadRequest)
		return
	}

	ms.AcknowledgeAlert(alertID)
	w.WriteHeader(http.StatusOK)
}

// cleanupOldMetrics removes old metric data points
func (ms *MonitoringSystem) cleanupOldMetrics() {
	cutoff := time.Now().Add(-ms.config.RetentionPeriod)

	ms.metrics.metrics.Range(func(key, value interface{}) bool {
		series := value.(*MetricSeries)
		var validPoints []MetricPoint
		for _, point := range series.DataPoints {
			if point.Timestamp.After(cutoff) {
				validPoints = append(validPoints, point)
			}
		}
		series.DataPoints = validPoints
		return true
	})
}

// Helper methods for system metrics (simplified implementations)
func (ms *MonitoringSystem) getMemoryUsage() float64 {
	// Simplified memory usage calculation
	// In a real implementation, this would use runtime.ReadMemStats()
	return 0.5 // 50% for demo
}

func (ms *MonitoringSystem) getCPUUsage() float64 {
	// Simplified CPU usage calculation
	// In a real implementation, this would use system calls
	return 0.3 // 30% for demo
}

func (ms *MonitoringSystem) getTestMetrics() map[string]float64 {
	// Simplified test metrics
	// In a real implementation, this would collect from test results
	return map[string]float64{
		"test_failure_rate": 0.05, // 5%
		"test_duration":     2.5,  // 2.5 seconds
		"test_count":        100,  // 100 tests
	}
}

// Close shuts down the monitoring system
func (ms *MonitoringSystem) Close() {
	ms.cancel()
	if ms.dashboard != nil {
		ms.dashboard.server.Shutdown(context.Background())
	}
}
