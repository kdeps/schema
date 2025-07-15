package test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/kdeps/kdeps/pkg/pklres"
)

// BenchmarkPKLFileEvaluation benchmarks PKL file evaluation performance
func BenchmarkPKLFileEvaluation(b *testing.B) {
	evaluator, err := NewTestEvaluator(&AgentResourceReader{}, &PklresResourceReader{})
	if err != nil {
		b.Fatalf("Failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	// Benchmark different PKL file types
	benchmarks := []struct {
		name     string
		fileName string
	}{
		{"Simple_PKL_File", "test_pklres_integration.pkl"},
		{"Exec_Tests", "exec_tests_pass.pkl"},
		{"Python_Tests", "python_tests_pass.pkl"},
		{"LLM_Tests", "llm_tests_pass.pkl"},
		{"HTTP_Tests", "http_tests_pass.pkl"},
		{"Data_Tests", "data_tests_pass.pkl"},
		{"Pklres_Tests", "pklres_tests_pass.pkl"},
		{"All_Tests", "all_tests_pass.pkl"},
		{"Test_Summary", "test_summary.pkl"},
	}

	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			BenchmarkPKLEvaluation(b, evaluator, bench.fileName)
		})
	}
}

// BenchmarkResourceReaderOperations benchmarks resource reader performance
func BenchmarkResourceReaderOperations(b *testing.B) {
	// Benchmark agent resource reader
	b.Run("Agent_Reader_Resolve", func(b *testing.B) {
		reader := &AgentResourceReader{}
		uri, _ := url.Parse("agent:/test-action")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to read agent resource: %v", err)
			}
		}
	})

	// Benchmark pklres resource reader
	b.Run("Pklres_Reader_Get", func(b *testing.B) {
		reader := &PklresResourceReader{}
		uri, _ := url.Parse("pklres:/test-id?type=exec&key=command&op=get")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to read pklres resource: %v", err)
			}
		}
	})

	// Benchmark pklres resource reader with set operation
	b.Run("Pklres_Reader_Set", func(b *testing.B) {
		reader := &PklresResourceReader{}
		uri, _ := url.Parse("pklres:/test-id?type=exec&key=command&op=set&value=echo%20hello")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to set pklres resource: %v", err)
			}
		}
	})
}

// BenchmarkRealPklresReader benchmarks the real pklres reader
func BenchmarkRealPklresReader(b *testing.B) {
	// Create a temporary database for benchmarking
	tempDB, err := os.CreateTemp("", "pklres-bench-*.db")
	if err != nil {
		b.Fatalf("Failed to create temp DB: %v", err)
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		b.Fatalf("Failed to initialize pklres reader: %v", err)
	}

	// Benchmark set operations
	b.Run("Real_Pklres_Set", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			uri, _ := url.Parse(fmt.Sprintf("pklres:/bench-id-%d?op=set&type=exec&key=command&value=echo%%20hello%%20%d", i, i))
			_, err := pklresReader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to set pklres resource: %v", err)
			}
		}
	})

	// Benchmark get operations
	b.Run("Real_Pklres_Get", func(b *testing.B) {
		// Pre-populate some data
		for i := 0; i < 100; i++ {
			uri, _ := url.Parse(fmt.Sprintf("pklres:/bench-id-%d?op=set&type=exec&key=command&value=echo%%20hello%%20%d", i, i))
			_, err := pklresReader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to set pklres resource: %v", err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			uri, _ := url.Parse(fmt.Sprintf("pklres:/bench-id-%d?op=get&type=exec&key=command", i%100))
			_, err := pklresReader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to get pklres resource: %v", err)
			}
		}
	})

	// Benchmark list operations
	b.Run("Real_Pklres_List", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			uri, _ := url.Parse("pklres:/?op=list&type=exec")
			_, err := pklresReader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to list pklres resources: %v", err)
			}
		}
	})
}

// BenchmarkConcurrentOperations benchmarks concurrent resource operations
func BenchmarkConcurrentOperations(b *testing.B) {
	tempDB, err := os.CreateTemp("", "pklres-concurrent-bench-*.db")
	if err != nil {
		b.Fatalf("Failed to create temp DB: %v", err)
	}
	defer os.Remove(tempDB.Name())
	tempDB.Close()

	pklresReader, err := pklres.InitializePklResource(tempDB.Name())
	if err != nil {
		b.Fatalf("Failed to initialize pklres reader: %v", err)
	}

	// Benchmark concurrent set operations
	b.Run("Concurrent_Set", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				uri, _ := url.Parse(fmt.Sprintf("pklres:/concurrent-id-%d?op=set&type=exec&key=command&value=echo%%20hello%%20%d", i, i))
				_, err := pklresReader.Read(*uri)
				if err != nil {
					b.Fatalf("Failed to set pklres resource: %v", err)
				}
				i++
			}
		})
	})

	// Benchmark concurrent get operations
	b.Run("Concurrent_Get", func(b *testing.B) {
		// Pre-populate data
		for i := 0; i < 1000; i++ {
			uri, _ := url.Parse(fmt.Sprintf("pklres:/concurrent-get-id-%d?op=set&type=exec&key=command&value=echo%%20hello%%20%d", i, i))
			_, err := pklresReader.Read(*uri)
			if err != nil {
				b.Fatalf("Failed to set pklres resource: %v", err)
			}
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				uri, _ := url.Parse(fmt.Sprintf("pklres:/concurrent-get-id-%d?op=get&type=exec&key=command", i%1000))
				_, err := pklresReader.Read(*uri)
				if err != nil {
					b.Fatalf("Failed to get pklres resource: %v", err)
				}
				i++
			}
		})
	})
}

// BenchmarkPKLEvaluatorCreation benchmarks evaluator creation performance
func BenchmarkPKLEvaluatorCreation(b *testing.B) {
	b.Run("Mock_Readers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			evaluator, err := NewTestEvaluator(&AgentResourceReader{}, &PklresResourceReader{})
			if err != nil {
				b.Fatalf("Failed to create evaluator: %v", err)
			}
			evaluator.Close()
		}
	})

	b.Run("Real_Readers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tempDB, err := os.CreateTemp("", "pklres-bench-*.db")
			if err != nil {
				b.Fatalf("Failed to create temp DB: %v", err)
			}
			tempDB.Close()
			defer os.Remove(tempDB.Name())

			pklresReader, err := pklres.InitializePklResource(tempDB.Name())
			if err != nil {
				b.Fatalf("Failed to initialize pklres reader: %v", err)
			}

			evaluator, err := NewTestEvaluator(&AgentResourceReader{}, pklresReader)
			if err != nil {
				b.Fatalf("Failed to create evaluator: %v", err)
			}
			evaluator.Close()
		}
	})
}
