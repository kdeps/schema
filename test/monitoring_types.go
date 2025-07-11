package test

import "time"

// AlertThresholds defines when alerts should be triggered
// Shared between performance and monitoring systems
// Only define here, import elsewhere
type AlertThresholds struct {
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	MemoryUsagePercent float64 `json:"memory_usage_percent"`
	TestDurationMs     int64   `json:"test_duration_ms"`
	FailureRatePercent float64 `json:"failure_rate_percent"`
}

// PerformanceAlert represents a performance alert
// Shared between performance and monitoring systems
type PerformanceAlert struct {
	ID           string    `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	Severity     string    `json:"severity"` // "low", "medium", "high", "critical"
	Category     string    `json:"category"`
	Message      string    `json:"message"`
	Metric       string    `json:"metric"`
	Value        float64   `json:"value"`
	Threshold    float64   `json:"threshold"`
	Acknowledged bool      `json:"acknowledged"`
}

// MemoryMetrics tracks memory usage
type MemoryMetrics struct {
	Alloc      uint64  `json:"alloc"`
	TotalAlloc uint64  `json:"total_alloc"`
	Sys        uint64  `json:"sys"`
	NumGC      uint32  `json:"num_gc"`
	HeapUsage  float64 `json:"heap_usage_percent"`
}

// ResourceUsageMetrics tracks resource consumption
type ResourceUsageMetrics struct {
	OpenFiles    int     `json:"open_files"`
	NetworkIO    float64 `json:"network_io_mb"`
	DiskIO       float64 `json:"disk_io_mb"`
	CacheHitRate float64 `json:"cache_hit_rate"`
}
