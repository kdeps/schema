//go:build ignore
// +build ignore

// This file intentionally uses package main for CLI usage.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type AnalyticsHistory struct {
	GeneratedAt time.Time `json:"generated_at"`
	Summary     struct {
		PassRate    float64       `json:"pass_rate"`
		AvgDuration time.Duration `json:"avg_duration"`
		TotalTests  int           `json:"total_tests"`
		FailedTests int           `json:"failed_tests"`
		FlakyTests  int           `json:"flaky_tests"`
	} `json:"summary"`
}

func main() {
	dir := "history"
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read history directory: %v\n", err)
		os.Exit(1)
	}

	var histories []AnalyticsHistory
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, file.Name())
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		var hist AnalyticsHistory
		if err := json.NewDecoder(f).Decode(&hist); err == nil {
			histories = append(histories, hist)
		}
		f.Close()
	}

	if len(histories) == 0 {
		fmt.Println("No analytics history found.")
		return
	}

	sort.Slice(histories, func(i, j int) bool {
		return histories[i].GeneratedAt.Before(histories[j].GeneratedAt)
	})

	fmt.Println("\n📈 Historical Analytics Trends:")
	fmt.Printf("%-20s %-10s %-10s %-10s %-10s\n", "Date", "PassRate", "AvgDur(s)", "Total", "Failed")
	for _, h := range histories {
		fmt.Printf("%-20s %8.2f%% %10.2f %10d %10d\n",
			h.GeneratedAt.Format("2006-01-02 15:04"),
			h.Summary.PassRate,
			h.Summary.AvgDuration.Seconds(),
			h.Summary.TotalTests,
			h.Summary.FailedTests)
	}

	// Optionally, print trend direction
	if len(histories) > 1 {
		first := histories[0]
		last := histories[len(histories)-1]
		fmt.Println("\nTrend Summary:")
		fmt.Printf("Pass Rate:   %.2f%% → %.2f%%\n", first.Summary.PassRate, last.Summary.PassRate)
		fmt.Printf("Avg Duration: %.2fs → %.2fs\n", first.Summary.AvgDuration.Seconds(), last.Summary.AvgDuration.Seconds())
		fmt.Printf("Failed:      %d → %d\n", first.Summary.FailedTests, last.Summary.FailedTests)
	}
}
