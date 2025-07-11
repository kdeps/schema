package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDistributedTestSystem(t *testing.T) {
	t.Run("NewDistributedTestSystem", func(t *testing.T) {
		config := &DistributedConfig{
			CoordinatorPort:   8081,
			WorkerPort:        8082,
			HeartbeatInterval: 30 * time.Second,
			WorkerTimeout:     2 * time.Minute,
			MaxWorkers:        10,
			LoadBalancing:     "round_robin",
		}

		system := NewDistributedTestSystem(config)
		defer system.Close()

		if system == nil {
			t.Fatal("Expected distributed test system to be created")
		}

		if system.config != config {
			t.Error("Expected config to match provided config")
		}

		if system.coordinator == nil {
			t.Error("Expected coordinator to be created")
		}
	})

	t.Run("CreateWorker", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		worker := system.CreateWorker("worker1", "localhost:8082")

		if worker == nil {
			t.Fatal("Expected worker to be created")
		}

		if worker.id != "worker1" {
			t.Errorf("Expected worker ID to be 'worker1', got %s", worker.id)
		}

		if worker.address != "localhost:8082" {
			t.Errorf("Expected worker address to be 'localhost:8082', got %s", worker.address)
		}

		// Check that worker was added to system
		system.mu.RLock()
		_, exists := system.workers["worker1"]
		system.mu.RUnlock()

		if !exists {
			t.Error("Expected worker to be added to system")
		}
	})

	t.Run("SubmitTask", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		parameters := map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		}

		taskID := system.SubmitTask("test_function", parameters, 1)

		if taskID == "" {
			t.Fatal("Expected task ID to be returned")
		}

		// Check that task was added to coordinator
		system.coordinator.mu.RLock()
		taskFound := false
		for _, task := range system.coordinator.tasks {
			if task.ID == taskID {
				taskFound = true
				if task.TestName != "test_function" {
					t.Errorf("Expected test name to be 'test_function', got %s", task.TestName)
				}
				if task.Priority != 1 {
					t.Errorf("Expected priority to be 1, got %d", task.Priority)
				}
				if task.Status != "pending" {
					t.Errorf("Expected status to be 'pending', got %s", task.Status)
				}
				break
			}
		}
		system.coordinator.mu.RUnlock()

		if !taskFound {
			t.Error("Expected task to be added to coordinator")
		}
	})

	t.Run("GetTaskResult", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		// Submit a task
		taskID := system.SubmitTask("test_function", nil, 1)

		// Try to get result (should not exist yet)
		result, exists := system.GetTaskResult(taskID)
		if exists {
			t.Error("Expected task result to not exist yet")
		}

		// Manually add a result
		system.coordinator.mu.Lock()
		system.coordinator.results[taskID] = &DistributedTestResult{
			TaskID:      taskID,
			Status:      "pass",
			Duration:    1 * time.Second,
			Output:      "Test passed",
			CompletedAt: time.Now(),
		}
		system.coordinator.mu.Unlock()

		// Now get the result
		result, exists = system.GetTaskResult(taskID)
		if !exists {
			t.Error("Expected task result to exist")
		}

		if result.Status != "pass" {
			t.Errorf("Expected result status to be 'pass', got %s", result.Status)
		}
	})

	t.Run("GetSystemStatus", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		// Create some workers
		system.CreateWorker("worker1", "localhost:8082")
		system.CreateWorker("worker2", "localhost:8083")

		// Manually register workers with coordinator (since CreateWorker doesn't do this automatically)
		system.coordinator.mu.Lock()
		system.coordinator.workers["worker1"] = &WorkerInfo{
			ID:            "worker1",
			Address:       "localhost:8082",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          0,
		}
		system.coordinator.workers["worker2"] = &WorkerInfo{
			ID:            "worker2",
			Address:       "localhost:8083",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          0,
		}
		system.coordinator.mu.Unlock()

		// Submit some tasks
		system.SubmitTask("test1", nil, 1)
		system.SubmitTask("test2", nil, 2)

		status := system.GetSystemStatus()

		// Check workers section
		workers, exists := status["workers"].(map[string]interface{})
		if !exists {
			t.Fatal("Expected workers section in status")
		}

		if workers["total"].(int) != 2 {
			t.Errorf("Expected 2 total workers, got %d", workers["total"].(int))
		}

		// Check tasks section
		tasks, exists := status["tasks"].(map[string]interface{})
		if !exists {
			t.Fatal("Expected tasks section in status")
		}

		if tasks["total"].(int) != 2 {
			t.Errorf("Expected 2 total tasks, got %d", tasks["total"].(int))
		}

		if tasks["pending"].(int) != 2 {
			t.Errorf("Expected 2 pending tasks, got %d", tasks["pending"].(int))
		}
	})
}

func TestLoadBalancing(t *testing.T) {
	t.Run("RoundRobinLoadBalancing", func(t *testing.T) {
		system := NewDistributedTestSystem(&DistributedConfig{
			CoordinatorPort:   8081,
			WorkerPort:        8082,
			HeartbeatInterval: 30 * time.Second,
			WorkerTimeout:     2 * time.Minute,
			MaxWorkers:        10,
			LoadBalancing:     "round_robin",
		})
		defer system.Close()

		// Create workers
		system.CreateWorker("worker1", "localhost:8082")
		system.CreateWorker("worker2", "localhost:8083")

		// Manually register workers with coordinator
		system.coordinator.mu.Lock()
		system.coordinator.workers["worker1"] = &WorkerInfo{
			ID:            "worker1",
			Address:       "localhost:8082",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          0,
		}
		system.coordinator.workers["worker2"] = &WorkerInfo{
			ID:            "worker2",
			Address:       "localhost:8083",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          0,
		}
		system.coordinator.mu.Unlock()

		// Submit tasks and check assignment
		task1 := &TestTask{
			ID:        "task1",
			TestName:  "test1",
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		task2 := &TestTask{
			ID:        "task2",
			TestName:  "test2",
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		system.assignTask(task1)
		system.assignTask(task2)

		// Check that tasks were assigned to different workers
		if task1.AssignedTo == task2.AssignedTo {
			t.Error("Expected tasks to be assigned to different workers in round-robin")
		}
	})

	t.Run("LeastLoadedLoadBalancing", func(t *testing.T) {
		system := NewDistributedTestSystem(&DistributedConfig{
			CoordinatorPort:   8081,
			WorkerPort:        8082,
			HeartbeatInterval: 30 * time.Second,
			WorkerTimeout:     2 * time.Minute,
			MaxWorkers:        10,
			LoadBalancing:     "least_loaded",
		})
		defer system.Close()

		// Create workers with different loads
		system.coordinator.mu.Lock()
		system.coordinator.workers["worker1"] = &WorkerInfo{
			ID:            "worker1",
			Address:       "localhost:8082",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          5, // Higher load
		}
		system.coordinator.workers["worker2"] = &WorkerInfo{
			ID:            "worker2",
			Address:       "localhost:8083",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          1, // Lower load
		}
		system.coordinator.mu.Unlock()

		task := &TestTask{
			ID:        "task1",
			TestName:  "test1",
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		system.assignTask(task)

		// Should be assigned to worker2 (least loaded)
		if task.AssignedTo != "worker2" {
			t.Errorf("Expected task to be assigned to worker2 (least loaded), got %s", task.AssignedTo)
		}
	})
}

func TestCoordinatorHTTPHandlers(t *testing.T) {
	t.Run("RegisterWorkerHandler", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		// Create proper worker info
		workerInfo := WorkerInfo{
			ID:            "test_worker",
			Address:       "localhost:8082",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          0,
		}

		// Send proper JSON data
		jsonData, _ := json.Marshal(workerInfo)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		system.coordinator.registerWorkerHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check that worker was registered
		system.coordinator.mu.RLock()
		_, exists := system.coordinator.workers["test_worker"]
		system.coordinator.mu.RUnlock()

		if !exists {
			t.Error("Expected worker to be registered")
		}
	})

	t.Run("HeartbeatHandler", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		// Register a worker first
		system.coordinator.mu.Lock()
		system.coordinator.workers["test_worker"] = &WorkerInfo{
			ID:            "test_worker",
			Address:       "localhost:8082",
			Status:        "idle",
			LastHeartbeat: time.Now().Add(-1 * time.Hour), // Old heartbeat
		}
		system.coordinator.mu.Unlock()

		// Create proper heartbeat data
		heartbeat := map[string]interface{}{
			"worker_id": "test_worker",
			"status":    "busy",
			"load":      3.0,
		}

		// Send proper JSON data
		jsonData, _ := json.Marshal(heartbeat)
		req := httptest.NewRequest("POST", "/heartbeat", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		system.coordinator.heartbeatHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check that heartbeat was updated
		system.coordinator.mu.RLock()
		worker := system.coordinator.workers["test_worker"]
		system.coordinator.mu.RUnlock()

		if worker.Status != "busy" {
			t.Errorf("Expected worker status to be 'busy', got %s", worker.Status)
		}

		if worker.Load != 3 {
			t.Errorf("Expected worker load to be 3, got %d", worker.Load)
		}
	})

	t.Run("CompleteTaskHandler", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		// Submit a task
		taskID := system.SubmitTask("test_function", nil, 1)

		// Create proper request body with JSON data
		result := DistributedTestResult{
			TaskID:      taskID,
			Status:      "pass",
			Duration:    1 * time.Second,
			Output:      "Test passed",
			CompletedAt: time.Now(),
		}

		jsonData, _ := json.Marshal(result)
		req := httptest.NewRequest("POST", "/task/complete", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		system.coordinator.completeTaskHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check that result was stored
		storedResult, exists := system.GetTaskResult(taskID)
		if !exists {
			t.Error("Expected task result to be stored")
		}

		if storedResult.Status != "pass" {
			t.Errorf("Expected result status to be 'pass', got %s", storedResult.Status)
		}
	})

	t.Run("StatusHandler", func(t *testing.T) {
		system := NewDistributedTestSystem(nil)
		defer system.Close()

		// Register some workers
		system.coordinator.mu.Lock()
		system.coordinator.workers["worker1"] = &WorkerInfo{
			ID:            "worker1",
			Address:       "localhost:8082",
			Status:        "idle",
			LastHeartbeat: time.Now(),
		}
		system.coordinator.mu.Unlock()

		req := httptest.NewRequest("GET", "/status", nil)
		w := httptest.NewRecorder()

		system.coordinator.statusHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected JSON content type, got %s", contentType)
		}

		// Parse response
		var workers map[string]*WorkerInfo
		if err := json.NewDecoder(w.Body).Decode(&workers); err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}

		if len(workers) != 1 {
			t.Errorf("Expected 1 worker in response, got %d", len(workers))
		}
	})
}

func TestWorkerManagement(t *testing.T) {
	t.Run("WorkerTimeout", func(t *testing.T) {
		config := &DistributedConfig{
			WorkerTimeout: 100 * time.Millisecond, // Short timeout for testing
		}
		system := NewDistributedTestSystem(config)
		defer system.Close()

		// Register a worker with old heartbeat
		system.coordinator.mu.Lock()
		system.coordinator.workers["old_worker"] = &WorkerInfo{
			ID:            "old_worker",
			Address:       "localhost:8082",
			Status:        "busy",
			LastHeartbeat: time.Now().Add(-200 * time.Millisecond), // Very old
		}
		system.coordinator.mu.Unlock()

		// Wait for cleanup
		time.Sleep(150 * time.Millisecond)

		// Trigger cleanup by calling monitorWorkers
		system.cleanupOfflineWorkers()

		// Check that worker was marked as offline
		system.coordinator.mu.RLock()
		worker := system.coordinator.workers["old_worker"]
		system.coordinator.mu.RUnlock()

		if worker.Status != "offline" {
			t.Errorf("Expected worker to be marked as offline, got %s", worker.Status)
		}
	})

	t.Run("TaskReassignment", func(t *testing.T) {
		config := &DistributedConfig{
			WorkerTimeout: 100 * time.Millisecond,
		}
		system := NewDistributedTestSystem(config)
		defer system.Close()

		// Create a task assigned to a worker
		task := &TestTask{
			ID:         "task1",
			TestName:   "test1",
			Status:     "running",
			AssignedTo: "worker1",
			CreatedAt:  time.Now(),
		}

		system.coordinator.mu.Lock()
		system.coordinator.tasks = append(system.coordinator.tasks, task)
		system.coordinator.workers["worker1"] = &WorkerInfo{
			ID:            "worker1",
			Address:       "localhost:8082",
			Status:        "busy",
			LastHeartbeat: time.Now().Add(-200 * time.Millisecond), // Old heartbeat
		}
		system.coordinator.mu.Unlock()

		// Trigger cleanup
		system.cleanupOfflineWorkers()

		// Check that task was reassigned
		if task.Status != "pending" {
			t.Errorf("Expected task status to be 'pending' after reassignment, got %s", task.Status)
		}

		if task.AssignedTo != "" {
			t.Errorf("Expected task to be unassigned after worker went offline, got %s", task.AssignedTo)
		}
	})
}

func TestDistributedSystemIntegration(t *testing.T) {
	t.Run("EndToEndWorkflow", func(t *testing.T) {
		system := NewDistributedTestSystem(&DistributedConfig{
			CoordinatorPort:   8081,
			WorkerPort:        8082,
			HeartbeatInterval: 30 * time.Second,
			WorkerTimeout:     2 * time.Minute,
			MaxWorkers:        10,
			LoadBalancing:     "round_robin",
		})
		defer system.Close()

		// Create workers (not used in this test but would be in real implementation)
		system.CreateWorker("worker1", "localhost:8082")
		system.CreateWorker("worker2", "localhost:8083")

		// Register workers with coordinator
		system.coordinator.mu.Lock()
		system.coordinator.workers["worker1"] = &WorkerInfo{
			ID:            "worker1",
			Address:       "localhost:8082",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          0,
		}
		system.coordinator.workers["worker2"] = &WorkerInfo{
			ID:            "worker2",
			Address:       "localhost:8083",
			Status:        "idle",
			LastHeartbeat: time.Now(),
			Load:          0,
		}
		system.coordinator.mu.Unlock()

		// Submit multiple tasks
		taskIDs := make([]string, 5)
		for i := 0; i < 5; i++ {
			taskIDs[i] = system.SubmitTask(fmt.Sprintf("test_%d", i), map[string]interface{}{
				"iteration": i,
			}, i+1)
		}

		// Check system status
		status := system.GetSystemStatus()
		tasks := status["tasks"].(map[string]interface{})
		if tasks["total"].(int) != 5 {
			t.Errorf("Expected 5 total tasks, got %d", tasks["total"].(int))
		}

		// Simulate task completion
		for _, taskID := range taskIDs {
			result := &DistributedTestResult{
				TaskID:      taskID,
				Status:      "pass",
				Duration:    time.Duration(taskID[len(taskID)-1]-'0') * time.Second,
				Output:      fmt.Sprintf("Task %s completed", taskID),
				CompletedAt: time.Now(),
			}

			system.coordinator.mu.Lock()
			system.coordinator.results[taskID] = result
			system.coordinator.mu.Unlock()
		}

		// Verify all results
		for _, taskID := range taskIDs {
			result, exists := system.GetTaskResult(taskID)
			if !exists {
				t.Errorf("Expected result for task %s", taskID)
			}
			if result.Status != "pass" {
				t.Errorf("Expected task %s to have 'pass' status, got %s", taskID, result.Status)
			}
		}
	})
}

// Benchmark tests for distributed system
func BenchmarkDistributedSystem(b *testing.B) {
	system := NewDistributedTestSystem(nil)
	defer system.Close()

	b.Run("SubmitTask", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			system.SubmitTask(fmt.Sprintf("benchmark_test_%d", i), map[string]interface{}{
				"benchmark": true,
				"iteration": i,
			}, i%5)
		}
	})

	b.Run("GetSystemStatus", func(b *testing.B) {
		// Pre-populate with some data
		for i := 0; i < 100; i++ {
			system.SubmitTask(fmt.Sprintf("status_test_%d", i), nil, 1)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			system.GetSystemStatus()
		}
	})

	b.Run("CreateWorker", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			system.CreateWorker(fmt.Sprintf("benchmark_worker_%d", i), fmt.Sprintf("localhost:%d", 8082+i))
		}
	})
}
