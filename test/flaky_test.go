package test

import (
	"os"
	"strconv"
	"testing"
	"time"
)

// FlakyTestConfig holds configuration for flaky test detection
type FlakyTestConfig struct {
	MaxRetries        int
	SlowTestThreshold time.Duration
}

// DefaultFlakyTestConfig returns the default config, reading from env if set
func DefaultFlakyTestConfig() FlakyTestConfig {
	maxRetries := 2
	if v := os.Getenv("FLAKY_MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxRetries = n
		}
	}
	return FlakyTestConfig{
		MaxRetries:        maxRetries,
		SlowTestThreshold: 2 * time.Second,
	}
}

// RunFlakyTest wraps a test function with flaky detection and retry logic
func RunFlakyTest(t *testing.T, name string, testFunc func(t *testing.T)) {
	cfg := DefaultFlakyTestConfig()
	retries := 0
	wasFlaky := false
	failed := false
	var duration time.Duration

	for {
		start := time.Now()
		innerT := &testing.T{}
		testFunc(innerT)
		duration = time.Since(start)

		if !innerT.Failed() {
			failed = false
			break
		}
		retries++
		failed = true
		if retries > cfg.MaxRetries {
			break
		}
		wasFlaky = true
	}

	// Output formatting
	status := "PASS"
	if failed {
		status = "FAIL"
	} else if wasFlaky {
		status = "FLAKY"
	}
	PrintTestResult(name, status, duration, "")

	if failed {
		t.Fail()
	}
}
