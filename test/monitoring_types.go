package test

import "time"

// IntegrationTestConfig holds configuration for integration tests
type IntegrationTestConfig struct {
	Timeout    time.Duration `json:"timeout"`
	Retries    int           `json:"retries"`
	Parallel   bool          `json:"parallel"`
	MaxWorkers int           `json:"max_workers"`
	LogLevel   string        `json:"log_level"`
}
