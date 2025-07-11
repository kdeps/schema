package test

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PerformanceOptimizer provides advanced performance optimization
type PerformanceOptimizer struct {
	cache   *ResourceCache
	pool    *PerformanceResourcePool
	monitor *PerformanceMonitor
	config  *PerformanceConfig
	mu      sync.RWMutex
}

// PerformanceConfig configures performance optimization
type PerformanceConfig struct {
	CacheSize       int           `json:"cache_size"`
	PoolSize        int           `json:"pool_size"`
	CacheTTL        time.Duration `json:"cache_ttl"`
	MonitorInterval time.Duration `json:"monitor_interval"`
	EnableProfiling bool          `json:"enable_profiling"`
}

// ResourceCache provides intelligent caching for expensive operations
type ResourceCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
	size  int
	ttl   time.Duration
}

// CacheItem represents a cached resource
type CacheItem struct {
	Value       interface{}
	ExpiresAt   time.Time
	AccessCount int
	LastAccess  time.Time
}

// PerformanceResourcePool manages reusable resources
type PerformanceResourcePool struct {
	resources chan interface{}
	factory   ResourceFactory
	mu        sync.Mutex
}

// ResourceFactory creates new resources
type ResourceFactory func() (interface{}, error)

// PerformanceMonitor tracks real-time performance metrics
type PerformanceMonitor struct {
	metrics     map[string]*PerformanceMetric
	alerts      []*PerformanceAlert
	config      *PerformanceMonitorConfig
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	alertChan   chan *PerformanceAlert
	metricsChan chan *SystemMetrics
}

// PerformanceMonitorConfig configures the performance monitor
type PerformanceMonitorConfig struct {
	Enabled          bool             `json:"enabled"`
	MonitorInterval  time.Duration    `json:"monitor_interval"`
	AlertThresholds  *AlertThresholds `json:"alert_thresholds"`
	MetricsRetention time.Duration    `json:"metrics_retention"`
	EnableProfiling  bool             `json:"enable_profiling"`
	ExportMetrics    bool             `json:"export_metrics"`
	ExportPath       string           `json:"export_path"`
}

// SystemMetrics tracks system-level performance indicators
type SystemMetrics struct {
	Timestamp      time.Time                     `json:"timestamp"`
	CPUUsage       float64                       `json:"cpu_usage"`
	MemoryUsage    *MemoryMetrics                `json:"memory_usage"`
	GoroutineCount int                           `json:"goroutine_count"`
	TestMetrics    map[string]*PerformanceMetric `json:"test_metrics"`
	ResourceUsage  *ResourceUsageMetrics         `json:"resource_usage"`
}

// PerformanceMetric tracks performance data
type PerformanceMetric struct {
	Name            string
	Count           int64
	TotalDuration   time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
	AverageDuration time.Duration
	LastUpdated     time.Time
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer(config *PerformanceConfig) *PerformanceOptimizer {
	if config == nil {
		config = &PerformanceConfig{
			CacheSize:       100,
			PoolSize:        10,
			CacheTTL:        5 * time.Minute,
			MonitorInterval: 30 * time.Second,
			EnableProfiling: false,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	optimizer := &PerformanceOptimizer{
		cache: &ResourceCache{
			items: make(map[string]*CacheItem),
			size:  config.CacheSize,
			ttl:   config.CacheTTL,
		},
		pool: &PerformanceResourcePool{
			resources: make(chan interface{}, config.PoolSize),
		},
		monitor: &PerformanceMonitor{
			metrics: make(map[string]*PerformanceMetric),
			ctx:     ctx,
			cancel:  cancel,
		},
		config: config,
	}

	// Start monitoring
	go optimizer.monitor.startMonitoring(config.MonitorInterval)

	return optimizer
}

// Get retrieves a value from cache
func (po *PerformanceOptimizer) Get(key string) (interface{}, bool) {
	po.cache.mu.RLock()
	defer po.cache.mu.RUnlock()

	item, exists := po.cache.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.ExpiresAt) {
		po.cache.mu.RUnlock()
		po.cache.mu.Lock()
		delete(po.cache.items, key)
		po.cache.mu.Unlock()
		po.cache.mu.RLock()
		return nil, false
	}

	// Update access stats
	item.AccessCount++
	item.LastAccess = time.Now()

	return item.Value, true
}

// Set stores a value in cache
func (po *PerformanceOptimizer) Set(key string, value interface{}) {
	po.cache.mu.Lock()
	defer po.cache.mu.Unlock()

	// Check cache size
	if len(po.cache.items) >= po.cache.size {
		po.evictLRU()
	}

	po.cache.items[key] = &CacheItem{
		Value:       value,
		ExpiresAt:   time.Now().Add(po.cache.ttl),
		AccessCount: 1,
		LastAccess:  time.Now(),
	}
}

// evictLRU removes least recently used items
func (po *PerformanceOptimizer) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range po.cache.items {
		if oldestKey == "" || item.LastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.LastAccess
		}
	}

	if oldestKey != "" {
		delete(po.cache.items, oldestKey)
	}
}

// AcquireResource gets a resource from the pool
func (po *PerformanceOptimizer) AcquireResource() (interface{}, error) {
	select {
	case resource := <-po.pool.resources:
		return resource, nil
	default:
		// Pool is empty, create new resource
		if po.pool.factory != nil {
			return po.pool.factory()
		}
		return nil, fmt.Errorf("no resources available and no factory configured")
	}
}

// ReleaseResource returns a resource to the pool
func (po *PerformanceOptimizer) ReleaseResource(resource interface{}) {
	select {
	case po.pool.resources <- resource:
		// Successfully returned to pool
	default:
		// Pool is full, discard resource
	}
}

// SetResourceFactory configures the resource factory
func (po *PerformanceOptimizer) SetResourceFactory(factory ResourceFactory) {
	po.pool.mu.Lock()
	defer po.pool.mu.Unlock()
	po.pool.factory = factory
}

// TrackPerformance records performance metrics
func (po *PerformanceOptimizer) TrackPerformance(name string, duration time.Duration) {
	po.monitor.mu.Lock()
	defer po.monitor.mu.Unlock()

	metric, exists := po.monitor.metrics[name]
	if !exists {
		metric = &PerformanceMetric{
			Name: name,
		}
		po.monitor.metrics[name] = metric
	}

	metric.Count++
	metric.TotalDuration += duration
	metric.LastUpdated = time.Now()

	// Update min/max
	if metric.Count == 1 || duration < metric.MinDuration {
		metric.MinDuration = duration
	}
	if duration > metric.MaxDuration {
		metric.MaxDuration = duration
	}

	// Calculate average
	metric.AverageDuration = metric.TotalDuration / time.Duration(metric.Count)
}

// GetPerformanceMetrics returns current performance metrics
func (po *PerformanceOptimizer) GetPerformanceMetrics() map[string]*PerformanceMetric {
	po.monitor.mu.RLock()
	defer po.monitor.mu.RUnlock()

	metrics := make(map[string]*PerformanceMetric)
	for name, metric := range po.monitor.metrics {
		metrics[name] = &PerformanceMetric{
			Name:            metric.Name,
			Count:           metric.Count,
			TotalDuration:   metric.TotalDuration,
			MinDuration:     metric.MinDuration,
			MaxDuration:     metric.MaxDuration,
			AverageDuration: metric.AverageDuration,
			LastUpdated:     metric.LastUpdated,
		}
	}

	return metrics
}

// startMonitoring starts the performance monitoring loop
func (pm *PerformanceMonitor) startMonitoring(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-pm.ctx.Done():
			return
		case <-ticker.C:
			pm.reportMetrics()
		}
	}
}

// reportMetrics reports current performance metrics
func (pm *PerformanceMonitor) reportMetrics() {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if len(pm.metrics) == 0 {
		return
	}

	fmt.Printf("\n=== Performance Metrics Report ===\n")
	fmt.Printf("Generated: %s\n\n", time.Now().Format(time.RFC3339))

	for name, metric := range pm.metrics {
		fmt.Printf("Metric: %s\n", name)
		fmt.Printf("  Count: %d\n", metric.Count)
		fmt.Printf("  Total Duration: %v\n", metric.TotalDuration)
		fmt.Printf("  Average Duration: %v\n", metric.AverageDuration)
		fmt.Printf("  Min Duration: %v\n", metric.MinDuration)
		fmt.Printf("  Max Duration: %v\n", metric.MaxDuration)
		fmt.Printf("  Last Updated: %s\n", metric.LastUpdated.Format(time.RFC3339))
		fmt.Println()
	}
}

// Close cleans up resources
func (po *PerformanceOptimizer) Close() {
	po.monitor.cancel()
}
