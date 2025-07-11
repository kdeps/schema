package test

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// TestFilter provides advanced test filtering capabilities
type TestFilter struct {
	criteria *FilterCriteria
	matcher  *FilterMatcher
}

// FilterCriteria defines filtering criteria
type FilterCriteria struct {
	// Basic filters
	NamePattern       string   `json:"name_pattern"`
	NameRegex         string   `json:"name_regex"`
	Tags              []string `json:"tags"`
	ExcludeTags       []string `json:"exclude_tags"`
	Categories        []string `json:"categories"`
	ExcludeCategories []string `json:"exclude_categories"`

	// Status filters
	Status        []string `json:"status"` // "pass", "fail", "skip", "flaky"
	ExcludeStatus []string `json:"exclude_status"`

	// Time-based filters
	LastRunAfter  *time.Time     `json:"last_run_after"`
	LastRunBefore *time.Time     `json:"last_run_before"`
	DurationMin   *time.Duration `json:"duration_min"`
	DurationMax   *time.Duration `json:"duration_max"`

	// Performance filters
	PerformanceThreshold float64 `json:"performance_threshold"`
	MemoryThreshold      uint64  `json:"memory_threshold"`

	// Metadata filters
	Metadata      map[string]interface{} `json:"metadata"`
	MetadataRegex map[string]string      `json:"metadata_regex"`

	// Dependency filters
	HasDependencies     bool     `json:"has_dependencies"`
	Dependencies        []string `json:"dependencies"`
	ExcludeDependencies []string `json:"exclude_dependencies"`

	// Resource filters
	Resources        []string `json:"resources"`
	ExcludeResources []string `json:"exclude_resources"`

	// Priority filters
	PriorityMin int `json:"priority_min"`
	PriorityMax int `json:"priority_max"`

	// Flaky test filters
	FlakyOnly      bool    `json:"flaky_only"`
	StableOnly     bool    `json:"stable_only"`
	FlakyThreshold float64 `json:"flaky_threshold"` // percentage of failures to consider flaky

	// Custom filters
	CustomFilters []CustomFilter `json:"custom_filters"`
}

// CustomFilter allows custom filtering logic
type CustomFilter struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Function    func(*TestInfo) bool   `json:"-"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// TestInfo contains information about a test for filtering
type TestInfo struct {
	Name         string                 `json:"name"`
	Tags         []string               `json:"tags"`
	Category     string                 `json:"category"`
	Status       string                 `json:"status"`
	Duration     time.Duration          `json:"duration"`
	LastRun      *time.Time             `json:"last_run"`
	Metadata     map[string]interface{} `json:"metadata"`
	Dependencies []string               `json:"dependencies"`
	Resources    []string               `json:"resources"`
	Priority     int                    `json:"priority"`
	FlakyScore   float64                `json:"flaky_score"`
	Performance  *PerformanceMetrics    `json:"performance"`
}

// PerformanceMetrics contains performance data
type PerformanceMetrics struct {
	AvgDuration  time.Duration `json:"avg_duration"`
	MaxDuration  time.Duration `json:"max_duration"`
	MinDuration  time.Duration `json:"min_duration"`
	MemoryUsage  uint64        `json:"memory_usage"`
	CPUUsage     float64       `json:"cpu_usage"`
	OpsPerSecond float64       `json:"ops_per_second"`
}

// FilterMatcher handles the actual filtering logic
type FilterMatcher struct {
	nameRegex     *regexp.Regexp
	metadataRegex map[string]*regexp.Regexp
}

// NewTestFilter creates a new test filter
func NewTestFilter() *TestFilter {
	return &TestFilter{
		criteria: &FilterCriteria{
			Metadata:      make(map[string]interface{}),
			MetadataRegex: make(map[string]string),
			CustomFilters: make([]CustomFilter, 0),
		},
		matcher: &FilterMatcher{
			metadataRegex: make(map[string]*regexp.Regexp),
		},
	}
}

// SetNamePattern sets a simple name pattern filter
func (tf *TestFilter) SetNamePattern(pattern string) *TestFilter {
	tf.criteria.NamePattern = pattern
	return tf
}

// SetNameRegex sets a regex name filter
func (tf *TestFilter) SetNameRegex(regex string) *TestFilter {
	tf.criteria.NameRegex = regex
	if regex != "" {
		tf.matcher.nameRegex = regexp.MustCompile(regex)
	}
	return tf
}

// AddTags adds tags to include
func (tf *TestFilter) AddTags(tags ...string) *TestFilter {
	tf.criteria.Tags = append(tf.criteria.Tags, tags...)
	return tf
}

// ExcludeTags adds tags to exclude
func (tf *TestFilter) ExcludeTags(tags ...string) *TestFilter {
	tf.criteria.ExcludeTags = append(tf.criteria.ExcludeTags, tags...)
	return tf
}

// AddCategories adds categories to include
func (tf *TestFilter) AddCategories(categories ...string) *TestFilter {
	tf.criteria.Categories = append(tf.criteria.Categories, categories...)
	return tf
}

// ExcludeCategories adds categories to exclude
func (tf *TestFilter) ExcludeCategories(categories ...string) *TestFilter {
	tf.criteria.ExcludeCategories = append(tf.criteria.ExcludeCategories, categories...)
	return tf
}

// SetStatus sets status filters
func (tf *TestFilter) SetStatus(status ...string) *TestFilter {
	tf.criteria.Status = status
	return tf
}

// ExcludeStatus sets status filters to exclude
func (tf *TestFilter) ExcludeStatus(status ...string) *TestFilter {
	tf.criteria.ExcludeStatus = status
	return tf
}

// SetTimeRange sets time-based filters
func (tf *TestFilter) SetTimeRange(after, before *time.Time) *TestFilter {
	tf.criteria.LastRunAfter = after
	tf.criteria.LastRunBefore = before
	return tf
}

// SetDurationRange sets duration-based filters
func (tf *TestFilter) SetDurationRange(min, max *time.Duration) *TestFilter {
	tf.criteria.DurationMin = min
	tf.criteria.DurationMax = max
	return tf
}

// SetPerformanceThreshold sets performance threshold
func (tf *TestFilter) SetPerformanceThreshold(threshold float64) *TestFilter {
	tf.criteria.PerformanceThreshold = threshold
	return tf
}

// SetMemoryThreshold sets memory threshold
func (tf *TestFilter) SetMemoryThreshold(threshold uint64) *TestFilter {
	tf.criteria.MemoryThreshold = threshold
	return tf
}

// AddMetadata adds metadata filter
func (tf *TestFilter) AddMetadata(key string, value interface{}) *TestFilter {
	tf.criteria.Metadata[key] = value
	return tf
}

// AddMetadataRegex adds metadata regex filter
func (tf *TestFilter) AddMetadataRegex(key, regex string) *TestFilter {
	tf.criteria.MetadataRegex[key] = regex
	if regex != "" {
		tf.matcher.metadataRegex[key] = regexp.MustCompile(regex)
	}
	return tf
}

// SetDependencyFilter sets dependency filters
func (tf *TestFilter) SetDependencyFilter(hasDeps bool, deps, excludeDeps []string) *TestFilter {
	tf.criteria.HasDependencies = hasDeps
	tf.criteria.Dependencies = deps
	tf.criteria.ExcludeDependencies = excludeDeps
	return tf
}

// SetResourceFilter sets resource filters
func (tf *TestFilter) SetResourceFilter(resources, excludeResources []string) *TestFilter {
	tf.criteria.Resources = resources
	tf.criteria.ExcludeResources = excludeResources
	return tf
}

// SetPriorityRange sets priority range
func (tf *TestFilter) SetPriorityRange(min, max int) *TestFilter {
	tf.criteria.PriorityMin = min
	tf.criteria.PriorityMax = max
	return tf
}

// SetFlakyFilter sets flaky test filters
func (tf *TestFilter) SetFlakyFilter(flakyOnly, stableOnly bool, threshold float64) *TestFilter {
	tf.criteria.FlakyOnly = flakyOnly
	tf.criteria.StableOnly = stableOnly
	tf.criteria.FlakyThreshold = threshold
	return tf
}

// AddCustomFilter adds a custom filter
func (tf *TestFilter) AddCustomFilter(name, description string, fn func(*TestInfo) bool, params map[string]interface{}) *TestFilter {
	tf.criteria.CustomFilters = append(tf.criteria.CustomFilters, CustomFilter{
		Name:        name,
		Description: description,
		Function:    fn,
		Parameters:  params,
	})
	return tf
}

// Matches checks if a test matches the filter criteria
func (tf *TestFilter) Matches(test *TestInfo) bool {
	// Name pattern filter
	if tf.criteria.NamePattern != "" {
		if !strings.Contains(strings.ToLower(test.Name), strings.ToLower(tf.criteria.NamePattern)) {
			return false
		}
	}

	// Name regex filter
	if tf.matcher.nameRegex != nil {
		if !tf.matcher.nameRegex.MatchString(test.Name) {
			return false
		}
	}

	// Tags filter
	if len(tf.criteria.Tags) > 0 {
		if !tf.hasAnyTag(test.Tags, tf.criteria.Tags) {
			return false
		}
	}

	// Exclude tags filter
	if len(tf.criteria.ExcludeTags) > 0 {
		if tf.hasAnyTag(test.Tags, tf.criteria.ExcludeTags) {
			return false
		}
	}

	// Categories filter
	if len(tf.criteria.Categories) > 0 {
		if !tf.hasAnyCategory(test.Category, tf.criteria.Categories) {
			return false
		}
	}

	// Exclude categories filter
	if len(tf.criteria.ExcludeCategories) > 0 {
		if tf.hasAnyCategory(test.Category, tf.criteria.ExcludeCategories) {
			return false
		}
	}

	// Status filter
	if len(tf.criteria.Status) > 0 {
		if !tf.hasAnyStatus(test.Status, tf.criteria.Status) {
			return false
		}
	}

	// Exclude status filter
	if len(tf.criteria.ExcludeStatus) > 0 {
		if tf.hasAnyStatus(test.Status, tf.criteria.ExcludeStatus) {
			return false
		}
	}

	// Time-based filters
	if tf.criteria.LastRunAfter != nil && test.LastRun != nil {
		if test.LastRun.Before(*tf.criteria.LastRunAfter) {
			return false
		}
	}

	if tf.criteria.LastRunBefore != nil && test.LastRun != nil {
		if test.LastRun.After(*tf.criteria.LastRunBefore) {
			return false
		}
	}

	// Duration filters
	if tf.criteria.DurationMin != nil {
		if test.Duration < *tf.criteria.DurationMin {
			return false
		}
	}

	if tf.criteria.DurationMax != nil {
		if test.Duration > *tf.criteria.DurationMax {
			return false
		}
	}

	// Performance filters
	if tf.criteria.PerformanceThreshold > 0 && test.Performance != nil {
		if test.Performance.OpsPerSecond < tf.criteria.PerformanceThreshold {
			return false
		}
	}

	if tf.criteria.MemoryThreshold > 0 && test.Performance != nil {
		if test.Performance.MemoryUsage > tf.criteria.MemoryThreshold {
			return false
		}
	}

	// Metadata filters
	for key, value := range tf.criteria.Metadata {
		if test.Metadata == nil {
			return false
		}
		if testValue, exists := test.Metadata[key]; !exists || testValue != value {
			return false
		}
	}

	// Metadata regex filters
	for key, regex := range tf.matcher.metadataRegex {
		if test.Metadata == nil {
			return false
		}
		if testValue, exists := test.Metadata[key]; !exists {
			return false
		} else {
			valueStr := fmt.Sprintf("%v", testValue)
			if !regex.MatchString(valueStr) {
				return false
			}
		}
	}

	// Dependency filters
	if tf.criteria.HasDependencies {
		if len(test.Dependencies) == 0 {
			return false
		}
	}

	if len(tf.criteria.Dependencies) > 0 {
		if !tf.hasAnyDependency(test.Dependencies, tf.criteria.Dependencies) {
			return false
		}
	}

	if len(tf.criteria.ExcludeDependencies) > 0 {
		if tf.hasAnyDependency(test.Dependencies, tf.criteria.ExcludeDependencies) {
			return false
		}
	}

	// Resource filters
	if len(tf.criteria.Resources) > 0 {
		if !tf.hasAnyResource(test.Resources, tf.criteria.Resources) {
			return false
		}
	}

	if len(tf.criteria.ExcludeResources) > 0 {
		if tf.hasAnyResource(test.Resources, tf.criteria.ExcludeResources) {
			return false
		}
	}

	// Priority filters
	if tf.criteria.PriorityMin > 0 {
		if test.Priority < tf.criteria.PriorityMin {
			return false
		}
	}

	if tf.criteria.PriorityMax > 0 {
		if test.Priority > tf.criteria.PriorityMax {
			return false
		}
	}

	// Flaky filters
	if tf.criteria.FlakyOnly {
		if test.FlakyScore < tf.criteria.FlakyThreshold {
			return false
		}
	}

	if tf.criteria.StableOnly {
		if test.FlakyScore >= tf.criteria.FlakyThreshold {
			return false
		}
	}

	// Custom filters
	for _, customFilter := range tf.criteria.CustomFilters {
		if customFilter.Function != nil && !customFilter.Function(test) {
			return false
		}
	}

	return true
}

// FilterTests filters a list of tests based on criteria
func (tf *TestFilter) FilterTests(tests []*TestInfo) []*TestInfo {
	var filtered []*TestInfo
	for _, test := range tests {
		if tf.Matches(test) {
			filtered = append(filtered, test)
		}
	}
	return filtered
}

// GetCriteria returns the current filter criteria
func (tf *TestFilter) GetCriteria() *FilterCriteria {
	return tf.criteria
}

// Clone creates a copy of the filter
func (tf *TestFilter) Clone() *TestFilter {
	clone := NewTestFilter()

	// Copy basic criteria
	clone.criteria.NamePattern = tf.criteria.NamePattern
	clone.criteria.NameRegex = tf.criteria.NameRegex
	clone.criteria.Tags = append([]string{}, tf.criteria.Tags...)
	clone.criteria.ExcludeTags = append([]string{}, tf.criteria.ExcludeTags...)
	clone.criteria.Categories = append([]string{}, tf.criteria.Categories...)
	clone.criteria.ExcludeCategories = append([]string{}, tf.criteria.ExcludeCategories...)
	clone.criteria.Status = append([]string{}, tf.criteria.Status...)
	clone.criteria.ExcludeStatus = append([]string{}, tf.criteria.ExcludeStatus...)

	// Copy time filters
	if tf.criteria.LastRunAfter != nil {
		after := *tf.criteria.LastRunAfter
		clone.criteria.LastRunAfter = &after
	}
	if tf.criteria.LastRunBefore != nil {
		before := *tf.criteria.LastRunBefore
		clone.criteria.LastRunBefore = &before
	}

	// Copy duration filters
	if tf.criteria.DurationMin != nil {
		min := *tf.criteria.DurationMin
		clone.criteria.DurationMin = &min
	}
	if tf.criteria.DurationMax != nil {
		max := *tf.criteria.DurationMax
		clone.criteria.DurationMax = &max
	}

	// Copy performance filters
	clone.criteria.PerformanceThreshold = tf.criteria.PerformanceThreshold
	clone.criteria.MemoryThreshold = tf.criteria.MemoryThreshold

	// Copy metadata filters
	for k, v := range tf.criteria.Metadata {
		clone.criteria.Metadata[k] = v
	}
	for k, v := range tf.criteria.MetadataRegex {
		clone.criteria.MetadataRegex[k] = v
	}

	// Copy dependency filters
	clone.criteria.HasDependencies = tf.criteria.HasDependencies
	clone.criteria.Dependencies = append([]string{}, tf.criteria.Dependencies...)
	clone.criteria.ExcludeDependencies = append([]string{}, tf.criteria.ExcludeDependencies...)

	// Copy resource filters
	clone.criteria.Resources = append([]string{}, tf.criteria.Resources...)
	clone.criteria.ExcludeResources = append([]string{}, tf.criteria.ExcludeResources...)

	// Copy priority filters
	clone.criteria.PriorityMin = tf.criteria.PriorityMin
	clone.criteria.PriorityMax = tf.criteria.PriorityMax

	// Copy flaky filters
	clone.criteria.FlakyOnly = tf.criteria.FlakyOnly
	clone.criteria.StableOnly = tf.criteria.StableOnly
	clone.criteria.FlakyThreshold = tf.criteria.FlakyThreshold

	// Copy custom filters
	clone.criteria.CustomFilters = append([]CustomFilter{}, tf.criteria.CustomFilters...)

	// Recompile regex patterns
	if tf.criteria.NameRegex != "" {
		clone.SetNameRegex(tf.criteria.NameRegex)
	}
	for key, regex := range tf.criteria.MetadataRegex {
		clone.AddMetadataRegex(key, regex)
	}

	return clone
}

// Helper methods
func (tf *TestFilter) hasAnyTag(testTags, filterTags []string) bool {
	for _, filterTag := range filterTags {
		for _, testTag := range testTags {
			if testTag == filterTag {
				return true
			}
		}
	}
	return false
}

func (tf *TestFilter) hasAnyCategory(testCategory string, filterCategories []string) bool {
	for _, filterCategory := range filterCategories {
		if testCategory == filterCategory {
			return true
		}
	}
	return false
}

func (tf *TestFilter) hasAnyStatus(testStatus string, filterStatuses []string) bool {
	for _, filterStatus := range filterStatuses {
		if testStatus == filterStatus {
			return true
		}
	}
	return false
}

func (tf *TestFilter) hasAnyDependency(testDeps, filterDeps []string) bool {
	for _, filterDep := range filterDeps {
		for _, testDep := range testDeps {
			if testDep == filterDep {
				return true
			}
		}
	}
	return false
}

func (tf *TestFilter) hasAnyResource(testResources, filterResources []string) bool {
	for _, filterResource := range filterResources {
		for _, testResource := range testResources {
			if testResource == filterResource {
				return true
			}
		}
	}
	return false
}

// Predefined filters
func NewQuickTestsFilter() *TestFilter {
	return NewTestFilter().
		SetDurationRange(nil, &[]time.Duration{5 * time.Second}[0]).
		AddTags("quick", "fast")
}

func NewSlowTestsFilter() *TestFilter {
	return NewTestFilter().
		SetDurationRange(&[]time.Duration{5 * time.Second}[0], nil).
		AddTags("slow", "integration")
}

func NewFlakyTestsFilter() *TestFilter {
	return NewTestFilter().
		SetFlakyFilter(true, false, 0.1).
		AddTags("flaky")
}

func NewStableTestsFilter() *TestFilter {
	return NewTestFilter().
		SetFlakyFilter(false, true, 0.1).
		AddTags("stable")
}

func NewUnitTestsFilter() *TestFilter {
	return NewTestFilter().
		AddCategories("unit").
		ExcludeCategories("integration", "e2e")
}

func NewIntegrationTestsFilter() *TestFilter {
	return NewTestFilter().
		AddCategories("integration").
		ExcludeCategories("unit", "e2e")
}

func NewE2ETestsFilter() *TestFilter {
	return NewTestFilter().
		AddCategories("e2e").
		ExcludeCategories("unit", "integration")
}

func NewPerformanceTestsFilter() *TestFilter {
	return NewTestFilter().
		AddTags("performance", "benchmark").
		SetPerformanceThreshold(1000) // 1000 ops/sec minimum
}

func NewResourceIntensiveTestsFilter() *TestFilter {
	return NewTestFilter().
		AddTags("resource-intensive").
		SetMemoryThreshold(100 * 1024 * 1024) // 100MB threshold
}
