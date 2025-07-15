package test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// ImprovementTracker tracks and manages code improvements
type ImprovementTracker struct {
	improvements map[string]*Improvement
	mu           sync.RWMutex
	config       *ImprovementConfig
}

// Improvement represents a code improvement
type Improvement struct {
	ID          string                 `json:"id"`
	Category    string                 `json:"category"`
	Priority    int                    `json:"priority"` // 1-5, 5 being highest
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"` // "pending", "in_progress", "completed", "deferred"
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// ImprovementConfig configures the improvement tracker
type ImprovementConfig struct {
	AutoTrackPerformance bool          `json:"auto_track_performance"`
	AutoTrackSecurity    bool          `json:"auto_track_security"`
	AutoTrackQuality     bool          `json:"auto_track_quality"`
	ReviewInterval       time.Duration `json:"review_interval"`
	MaxImprovements      int           `json:"max_improvements"`
}

// NewImprovementTracker creates a new improvement tracker
func NewImprovementTracker(config *ImprovementConfig) *ImprovementTracker {
	if config == nil {
		config = &ImprovementConfig{
			AutoTrackPerformance: true,
			AutoTrackSecurity:    true,
			AutoTrackQuality:     true,
			ReviewInterval:       24 * time.Hour,
			MaxImprovements:      100,
		}
	}

	return &ImprovementTracker{
		improvements: make(map[string]*Improvement),
		config:       config,
	}
}

// AddImprovement adds a new improvement
func (it *ImprovementTracker) AddImprovement(improvement *Improvement) error {
	it.mu.Lock()
	defer it.mu.Unlock()

	if len(it.improvements) >= it.config.MaxImprovements {
		return fmt.Errorf("maximum number of improvements reached")
	}

	if improvement.ID == "" {
		improvement.ID = fmt.Sprintf("imp_%d", time.Now().UnixNano())
	}

	if improvement.CreatedAt.IsZero() {
		improvement.CreatedAt = time.Now()
	}

	improvement.UpdatedAt = time.Now()

	it.improvements[improvement.ID] = improvement
	return nil
}

// GetImprovement retrieves an improvement by ID
func (it *ImprovementTracker) GetImprovement(id string) (*Improvement, bool) {
	it.mu.RLock()
	defer it.mu.RUnlock()

	improvement, exists := it.improvements[id]
	return improvement, exists
}

// ListImprovements returns all improvements, optionally filtered
func (it *ImprovementTracker) ListImprovements(filters ...ImprovementFilter) []*Improvement {
	it.mu.RLock()
	defer it.mu.RUnlock()

	var improvements []*Improvement
	for _, improvement := range it.improvements {
		if it.matchesFilters(improvement, filters...) {
			improvements = append(improvements, improvement)
		}
	}

	return improvements
}

// UpdateImprovement updates an existing improvement
func (it *ImprovementTracker) UpdateImprovement(id string, updates map[string]interface{}) error {
	it.mu.Lock()
	defer it.mu.Unlock()

	improvement, exists := it.improvements[id]
	if !exists {
		return fmt.Errorf("improvement not found: %s", id)
	}

	// Apply updates
	for key, value := range updates {
		switch key {
		case "status":
			if status, ok := value.(string); ok {
				improvement.Status = status
				if status == "completed" && improvement.CompletedAt == nil {
					now := time.Now()
					improvement.CompletedAt = &now
				}
			}
		case "priority":
			if priority, ok := value.(int); ok {
				improvement.Priority = priority
			}
		case "description":
			if desc, ok := value.(string); ok {
				improvement.Description = desc
			}
		case "tags":
			if tags, ok := value.([]string); ok {
				improvement.Tags = tags
			}
		}
	}

	improvement.UpdatedAt = time.Now()
	return nil
}

// DeleteImprovement removes an improvement
func (it *ImprovementTracker) DeleteImprovement(id string) error {
	it.mu.Lock()
	defer it.mu.Unlock()

	if _, exists := it.improvements[id]; !exists {
		return fmt.Errorf("improvement not found: %s", id)
	}

	delete(it.improvements, id)
	return nil
}

// ImprovementFilter filters improvements
type ImprovementFilter func(*Improvement) bool

// FilterByCategory filters by category
func FilterByCategory(category string) ImprovementFilter {
	return func(imp *Improvement) bool {
		return imp.Category == category
	}
}

// FilterByStatus filters by status
func FilterByStatus(status string) ImprovementFilter {
	return func(imp *Improvement) bool {
		return imp.Status == status
	}
}

// FilterByPriority filters by priority
func FilterByPriority(minPriority int) ImprovementFilter {
	return func(imp *Improvement) bool {
		return imp.Priority >= minPriority
	}
}

// FilterByTag filters by tag
func FilterByTag(tag string) ImprovementFilter {
	return func(imp *Improvement) bool {
		for _, t := range imp.Tags {
			if t == tag {
				return true
			}
		}
		return false
	}
}

// matchesFilters checks if an improvement matches all filters
func (it *ImprovementTracker) matchesFilters(improvement *Improvement, filters ...ImprovementFilter) bool {
	for _, filter := range filters {
		if !filter(improvement) {
			return false
		}
	}
	return true
}

// GenerateImprovementReport generates a comprehensive improvement report
func (it *ImprovementTracker) GenerateImprovementReport() *ImprovementReport {
	it.mu.RLock()
	defer it.mu.RUnlock()

	report := &ImprovementReport{
		GeneratedAt: time.Now(),
		Summary:     make(map[string]int),
		ByCategory:  make(map[string][]*Improvement),
		ByPriority:  make(map[int][]*Improvement),
		ByStatus:    make(map[string][]*Improvement),
	}

	for _, improvement := range it.improvements {
		// Update summary
		report.TotalImprovements++
		report.Summary[improvement.Category]++

		// Group by category
		report.ByCategory[improvement.Category] = append(report.ByCategory[improvement.Category], improvement)

		// Group by priority
		report.ByPriority[improvement.Priority] = append(report.ByPriority[improvement.Priority], improvement)

		// Group by status
		report.ByStatus[improvement.Status] = append(report.ByStatus[improvement.Status], improvement)

		// Track completion rate
		if improvement.Status == "completed" {
			report.CompletedImprovements++
		}
	}

	if report.TotalImprovements > 0 {
		report.CompletionRate = float64(report.CompletedImprovements) / float64(report.TotalImprovements) * 100
	}

	return report
}

// ImprovementReport represents a comprehensive improvement report
type ImprovementReport struct {
	GeneratedAt           time.Time                 `json:"generated_at"`
	TotalImprovements     int                       `json:"total_improvements"`
	CompletedImprovements int                       `json:"completed_improvements"`
	CompletionRate        float64                   `json:"completion_rate"`
	Summary               map[string]int            `json:"summary"`
	ByCategory            map[string][]*Improvement `json:"by_category"`
	ByPriority            map[int][]*Improvement    `json:"by_priority"`
	ByStatus              map[string][]*Improvement `json:"by_status"`
}

// AutoTrackPerformance automatically tracks performance-related improvements
func (it *ImprovementTracker) AutoTrackPerformance(ctx context.Context) {
	if !it.config.AutoTrackPerformance {
		return
	}

	// Track common performance improvements
	performanceImprovements := []*Improvement{
		{
			Category:    "performance",
			Priority:    4,
			Title:       "Implement Connection Pooling",
			Description: "Add connection pooling for database and HTTP connections to reduce connection overhead",
			Status:      "pending",
			Tags:        []string{"database", "http", "connection-pooling"},
		},
		{
			Category:    "performance",
			Priority:    3,
			Title:       "Add Response Caching",
			Description: "Implement intelligent caching for frequently accessed data and API responses",
			Status:      "pending",
			Tags:        []string{"caching", "api", "response-time"},
		},
		{
			Category:    "performance",
			Priority:    3,
			Title:       "Optimize Memory Usage",
			Description: "Review and optimize memory allocation patterns, especially in hot paths",
			Status:      "pending",
			Tags:        []string{"memory", "optimization", "gc"},
		},
	}

	for _, imp := range performanceImprovements {
		it.AddImprovement(imp)
	}
}

// AutoTrackSecurity automatically tracks security-related improvements
func (it *ImprovementTracker) AutoTrackSecurity(ctx context.Context) {
	if !it.config.AutoTrackSecurity {
		return
	}

	// Track common security improvements
	securityImprovements := []*Improvement{
		{
			Category:    "security",
			Priority:    5,
			Title:       "Add Input Validation",
			Description: "Implement comprehensive input validation for all user inputs and API endpoints",
			Status:      "pending",
			Tags:        []string{"validation", "security", "api"},
		},
		{
			Category:    "security",
			Priority:    4,
			Title:       "Implement Rate Limiting",
			Description: "Add rate limiting to prevent abuse and DoS attacks",
			Status:      "pending",
			Tags:        []string{"rate-limiting", "security", "dos-protection"},
		},
		{
			Category:    "security",
			Priority:    4,
			Title:       "Add Request Logging",
			Description: "Implement comprehensive request logging for security monitoring and debugging",
			Status:      "pending",
			Tags:        []string{"logging", "security", "monitoring"},
		},
	}

	for _, imp := range securityImprovements {
		it.AddImprovement(imp)
	}
}

// AutoTrackQuality automatically tracks code quality improvements
func (it *ImprovementTracker) AutoTrackQuality(ctx context.Context) {
	if !it.config.AutoTrackQuality {
		return
	}

	// Track common code quality improvements
	qualityImprovements := []*Improvement{
		{
			Category:    "quality",
			Priority:    3,
			Title:       "Add Comprehensive Error Handling",
			Description: "Implement proper error handling and user-friendly error messages throughout the codebase",
			Status:      "pending",
			Tags:        []string{"error-handling", "user-experience", "robustness"},
		},
		{
			Category:    "quality",
			Priority:    3,
			Title:       "Improve Code Documentation",
			Description: "Add comprehensive documentation for all public APIs and complex functions",
			Status:      "pending",
			Tags:        []string{"documentation", "api", "maintainability"},
		},
		{
			Category:    "quality",
			Priority:    2,
			Title:       "Add Code Metrics",
			Description: "Implement code quality metrics and automated quality checks in CI/CD",
			Status:      "pending",
			Tags:        []string{"metrics", "ci-cd", "quality-gates"},
		},
	}

	for _, imp := range qualityImprovements {
		it.AddImprovement(imp)
	}
}

// TestImprovementTracker tests the improvement tracking functionality
func TestImprovementTracker(t *testing.T) {
	config := &ImprovementConfig{
		AutoTrackPerformance: true,
		AutoTrackSecurity:    true,
		AutoTrackQuality:     true,
		ReviewInterval:       1 * time.Hour,
		MaxImprovements:      50,
	}

	tracker := NewImprovementTracker(config)

	// Test adding improvements
	t.Run("AddImprovements", func(t *testing.T) {
		improvement := &Improvement{
			Category:    "performance",
			Priority:    4,
			Title:       "Test Improvement",
			Description: "A test improvement for validation",
			Status:      "pending",
			Tags:        []string{"test", "validation"},
		}

		err := tracker.AddImprovement(improvement)
		if err != nil {
			t.Errorf("Failed to add improvement: %v", err)
		}

		if improvement.ID == "" {
			t.Error("Expected improvement ID to be set")
		}

		if improvement.CreatedAt.IsZero() {
			t.Error("Expected creation time to be set")
		}
	})

	// Test retrieving improvements
	t.Run("GetImprovements", func(t *testing.T) {
		improvement := &Improvement{
			Category:    "security",
			Priority:    5,
			Title:       "Security Test",
			Description: "A security test improvement",
			Status:      "pending",
		}

		tracker.AddImprovement(improvement)

		retrieved, exists := tracker.GetImprovement(improvement.ID)
		if !exists {
			t.Error("Expected improvement to exist")
		}

		if retrieved.Title != improvement.Title {
			t.Errorf("Expected title %s, got %s", improvement.Title, retrieved.Title)
		}
	})

	// Test filtering
	t.Run("FilterImprovements", func(t *testing.T) {
		// Add some test improvements
		improvements := []*Improvement{
			{Category: "performance", Priority: 3, Title: "Perf 1", Status: "pending"},
			{Category: "security", Priority: 4, Title: "Sec 1", Status: "completed"},
			{Category: "quality", Priority: 2, Title: "Qual 1", Status: "pending"},
		}

		for _, imp := range improvements {
			tracker.AddImprovement(imp)
		}

		// Test category filter
		perfImprovements := tracker.ListImprovements(FilterByCategory("performance"))
		if len(perfImprovements) == 0 {
			t.Error("Expected performance improvements")
		}

		// Test status filter
		completedImprovements := tracker.ListImprovements(FilterByStatus("completed"))
		if len(completedImprovements) == 0 {
			t.Error("Expected completed improvements")
		}

		// Test priority filter
		highPriorityImprovements := tracker.ListImprovements(FilterByPriority(4))
		if len(highPriorityImprovements) == 0 {
			t.Error("Expected high priority improvements")
		}
	})

	// Test report generation
	t.Run("GenerateReport", func(t *testing.T) {
		report := tracker.GenerateImprovementReport()

		if report.TotalImprovements == 0 {
			t.Error("Expected total improvements to be greater than 0")
		}

		if report.GeneratedAt.IsZero() {
			t.Error("Expected report generation time to be set")
		}

		if len(report.Summary) == 0 {
			t.Error("Expected summary to contain data")
		}
	})

	// Test auto-tracking
	t.Run("AutoTracking", func(t *testing.T) {
		ctx := context.Background()

		tracker.AutoTrackPerformance(ctx)
		tracker.AutoTrackSecurity(ctx)
		tracker.AutoTrackQuality(ctx)

		allImprovements := tracker.ListImprovements()
		if len(allImprovements) == 0 {
			t.Error("Expected auto-tracked improvements")
		}

		// Verify categories are present
		categories := make(map[string]bool)
		for _, imp := range allImprovements {
			categories[imp.Category] = true
		}

		expectedCategories := []string{"performance", "security", "quality"}
		for _, category := range expectedCategories {
			if !categories[category] {
				t.Errorf("Expected category %s to be present", category)
			}
		}
	})

	t.Log("✅ Improvement tracker functionality working correctly")
}
