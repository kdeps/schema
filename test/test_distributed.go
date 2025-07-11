package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// DistributedTestSystem provides distributed test execution capabilities
type DistributedTestSystem struct {
	coordinator *TestCoordinator
	workers     map[string]*TestWorker
	config      *DistributedConfig
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// DistributedConfig configures distributed testing
type DistributedConfig struct {
	CoordinatorPort   int           `json:"coordinator_port"`
	WorkerPort        int           `json:"worker_port"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	WorkerTimeout     time.Duration `json:"worker_timeout"`
	MaxWorkers        int           `json:"max_workers"`
	LoadBalancing     string        `json:"load_balancing"` // "round_robin", "least_loaded", "random"
}

// TestCoordinator manages distributed test execution
type TestCoordinator struct {
	workers map[string]*WorkerInfo
	tasks   []*TestTask
	results map[string]*DistributedTestResult
	config  *DistributedConfig
	mu      sync.RWMutex
	server  *http.Server
}

// WorkerInfo represents information about a test worker
type WorkerInfo struct {
	ID            string    `json:"id"`
	Address       string    `json:"address"`
	Status        string    `json:"status"` // "idle", "busy", "offline"
	LastHeartbeat time.Time `json:"last_heartbeat"`
	Load          int       `json:"load"` // number of active tasks
	Capabilities  []string  `json:"capabilities"`
}

// TestTask represents a test task to be executed
type TestTask struct {
	ID          string                 `json:"id"`
	TestName    string                 `json:"test_name"`
	Parameters  map[string]interface{} `json:"parameters"`
	Priority    int                    `json:"priority"`
	AssignedTo  string                 `json:"assigned_to"`
	Status      string                 `json:"status"` // "pending", "running", "completed", "failed"
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	Result      *DistributedTestResult `json:"result"`
}

// TestWorker represents a test execution worker
type TestWorker struct {
	id          string
	address     string
	coordinator string
	config      *DistributedConfig
	ctx         context.Context
	cancel      context.CancelFunc
}

// DistributedTestResult represents the result of a distributed test execution
type DistributedTestResult struct {
	TaskID      string                 `json:"task_id"`
	Status      string                 `json:"status"` // "pass", "fail", "skip", "error"
	Duration    time.Duration          `json:"duration"`
	Output      string                 `json:"output"`
	Error       string                 `json:"error"`
	Metrics     map[string]interface{} `json:"metrics"`
	CompletedAt time.Time              `json:"completed_at"`
}

// NewDistributedTestSystem creates a new distributed test system
func NewDistributedTestSystem(config *DistributedConfig) *DistributedTestSystem {
	if config == nil {
		config = &DistributedConfig{
			CoordinatorPort:   8081,
			WorkerPort:        8082,
			HeartbeatInterval: 30 * time.Second,
			WorkerTimeout:     2 * time.Minute,
			MaxWorkers:        10,
			LoadBalancing:     "round_robin",
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	system := &DistributedTestSystem{
		workers: make(map[string]*TestWorker),
		config:  config,
		ctx:     ctx,
		cancel:  cancel,
	}

	system.coordinator = system.createCoordinator()

	return system
}

// createCoordinator creates the test coordinator
func (dts *DistributedTestSystem) createCoordinator() *TestCoordinator {
	coordinator := &TestCoordinator{
		workers: make(map[string]*WorkerInfo),
		tasks:   make([]*TestTask, 0),
		results: make(map[string]*DistributedTestResult),
		config:  dts.config,
	}

	// Create HTTP server for coordinator
	mux := http.NewServeMux()
	mux.HandleFunc("/register", coordinator.registerWorkerHandler)
	mux.HandleFunc("/heartbeat", coordinator.heartbeatHandler)
	mux.HandleFunc("/task/assign", coordinator.assignTaskHandler)
	mux.HandleFunc("/task/complete", coordinator.completeTaskHandler)
	mux.HandleFunc("/status", coordinator.statusHandler)

	coordinator.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", dts.config.CoordinatorPort),
		Handler: mux,
	}

	return coordinator
}

// StartCoordinator starts the coordinator server
func (dts *DistributedTestSystem) StartCoordinator() error {
	fmt.Printf("Starting test coordinator on port %d\n", dts.config.CoordinatorPort)

	// Start worker monitoring
	go dts.monitorWorkers()

	return dts.coordinator.server.ListenAndServe()
}

// CreateWorker creates a new test worker
func (dts *DistributedTestSystem) CreateWorker(id, address string) *TestWorker {
	worker := &TestWorker{
		id:          id,
		address:     address,
		coordinator: fmt.Sprintf("http://localhost:%d", dts.config.CoordinatorPort),
		config:      dts.config,
	}

	ctx, cancel := context.WithCancel(dts.ctx)
	worker.ctx = ctx
	worker.cancel = cancel

	dts.mu.Lock()
	dts.workers[id] = worker
	dts.mu.Unlock()

	// Start worker
	go worker.start()

	return worker
}

// SubmitTask submits a test task for execution
func (dts *DistributedTestSystem) SubmitTask(testName string, parameters map[string]interface{}, priority int) string {
	task := &TestTask{
		ID:         fmt.Sprintf("task_%d", time.Now().UnixNano()),
		TestName:   testName,
		Parameters: parameters,
		Priority:   priority,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	dts.coordinator.mu.Lock()
	dts.coordinator.tasks = append(dts.coordinator.tasks, task)
	dts.coordinator.mu.Unlock()

	// Try to assign task immediately
	go dts.assignTask(task)

	return task.ID
}

// GetTaskResult gets the result of a task
func (dts *DistributedTestSystem) GetTaskResult(taskID string) (*DistributedTestResult, bool) {
	dts.coordinator.mu.RLock()
	defer dts.coordinator.mu.RUnlock()

	result, exists := dts.coordinator.results[taskID]
	return result, exists
}

// GetSystemStatus returns the current status of the distributed system
func (dts *DistributedTestSystem) GetSystemStatus() map[string]interface{} {
	dts.coordinator.mu.RLock()
	defer dts.coordinator.mu.RUnlock()

	status := map[string]interface{}{
		"workers": map[string]interface{}{
			"total":  len(dts.coordinator.workers),
			"online": 0,
			"busy":   0,
			"idle":   0,
		},
		"tasks": map[string]interface{}{
			"total":     len(dts.coordinator.tasks),
			"pending":   0,
			"running":   0,
			"completed": 0,
			"failed":    0,
		},
	}

	// Count worker statuses
	for _, worker := range dts.coordinator.workers {
		if time.Since(worker.LastHeartbeat) < dts.config.WorkerTimeout {
			status["workers"].(map[string]interface{})["online"] = status["workers"].(map[string]interface{})["online"].(int) + 1
			if worker.Status == "busy" {
				status["workers"].(map[string]interface{})["busy"] = status["workers"].(map[string]interface{})["busy"].(int) + 1
			} else {
				status["workers"].(map[string]interface{})["idle"] = status["workers"].(map[string]interface{})["idle"].(int) + 1
			}
		}
	}

	// Count task statuses
	for _, task := range dts.coordinator.tasks {
		switch task.Status {
		case "pending":
			status["tasks"].(map[string]interface{})["pending"] = status["tasks"].(map[string]interface{})["pending"].(int) + 1
		case "running":
			status["tasks"].(map[string]interface{})["running"] = status["tasks"].(map[string]interface{})["running"].(int) + 1
		case "completed":
			status["tasks"].(map[string]interface{})["completed"] = status["tasks"].(map[string]interface{})["completed"].(int) + 1
		case "failed":
			status["tasks"].(map[string]interface{})["failed"] = status["tasks"].(map[string]interface{})["failed"].(int) + 1
		}
	}

	return status
}

// assignTask assigns a task to an available worker
func (dts *DistributedTestSystem) assignTask(task *TestTask) {
	dts.coordinator.mu.Lock()
	defer dts.coordinator.mu.Unlock()

	// Find available worker based on load balancing strategy
	var selectedWorker *WorkerInfo
	switch dts.config.LoadBalancing {
	case "round_robin":
		selectedWorker = dts.selectWorkerRoundRobin()
	case "least_loaded":
		selectedWorker = dts.selectWorkerLeastLoaded()
	case "random":
		selectedWorker = dts.selectWorkerRandom()
	default:
		selectedWorker = dts.selectWorkerRoundRobin()
	}

	if selectedWorker != nil {
		task.AssignedTo = selectedWorker.ID
		task.Status = "running"
		now := time.Now()
		task.StartedAt = &now
		selectedWorker.Status = "busy"
		selectedWorker.Load++
	}
}

// selectWorkerRoundRobin selects worker using round-robin strategy
func (dts *DistributedTestSystem) selectWorkerRoundRobin() *WorkerInfo {
	// Simplified round-robin implementation
	for _, worker := range dts.coordinator.workers {
		if worker.Status == "idle" && time.Since(worker.LastHeartbeat) < dts.config.WorkerTimeout {
			return worker
		}
	}
	return nil
}

// selectWorkerLeastLoaded selects the least loaded worker
func (dts *DistributedTestSystem) selectWorkerLeastLoaded() *WorkerInfo {
	var selectedWorker *WorkerInfo
	minLoad := int(^uint(0) >> 1) // Max int

	for _, worker := range dts.coordinator.workers {
		if worker.Status == "idle" && time.Since(worker.LastHeartbeat) < dts.config.WorkerTimeout {
			if worker.Load < minLoad {
				minLoad = worker.Load
				selectedWorker = worker
			}
		}
	}
	return selectedWorker
}

// selectWorkerRandom selects a random available worker
func (dts *DistributedTestSystem) selectWorkerRandom() *WorkerInfo {
	// Simplified random selection
	return dts.selectWorkerRoundRobin()
}

// monitorWorkers monitors worker health and cleans up offline workers
func (dts *DistributedTestSystem) monitorWorkers() {
	ticker := time.NewTicker(dts.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-dts.ctx.Done():
			return
		case <-ticker.C:
			dts.cleanupOfflineWorkers()
		}
	}
}

// cleanupOfflineWorkers removes workers that haven't sent heartbeats
func (dts *DistributedTestSystem) cleanupOfflineWorkers() {
	dts.coordinator.mu.Lock()
	defer dts.coordinator.mu.Unlock()

	for id, worker := range dts.coordinator.workers {
		if time.Since(worker.LastHeartbeat) > dts.config.WorkerTimeout {
			// Mark worker as offline
			worker.Status = "offline"

			// Reassign tasks from offline worker
			for _, task := range dts.coordinator.tasks {
				if task.AssignedTo == id && task.Status == "running" {
					task.Status = "pending"
					task.AssignedTo = ""
					task.StartedAt = nil
				}
			}
		}
	}
}

// start starts the worker
func (tw *TestWorker) start() {
	// Start heartbeat loop
	go tw.sendHeartbeat()

	// Start task execution loop
	go tw.executeTasks()
}

// sendHeartbeat sends periodic heartbeats to coordinator
func (tw *TestWorker) sendHeartbeat() {
	ticker := time.NewTicker(tw.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tw.ctx.Done():
			return
		case <-ticker.C:
			tw.sendHeartbeatRequest()
		}
	}
}

// sendHeartbeatRequest sends a heartbeat to the coordinator
func (tw *TestWorker) sendHeartbeatRequest() {
	// In a real implementation, this would send the actual JSON data
	// For now, we'll just make the HTTP call without the body
	http.Post(tw.coordinator+"/heartbeat", "application/json", nil)
}

// executeTasks executes assigned tasks
func (tw *TestWorker) executeTasks() {
	// Simplified task execution
	// In a real implementation, this would poll for tasks and execute them
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-tw.ctx.Done():
			return
		case <-ticker.C:
			// Poll for tasks (simplified)
		}
	}
}

// HTTP handlers for coordinator
func (tc *TestCoordinator) registerWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var workerInfo WorkerInfo
	if err := json.NewDecoder(r.Body).Decode(&workerInfo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tc.mu.Lock()
	tc.workers[workerInfo.ID] = &workerInfo
	tc.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (tc *TestCoordinator) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var heartbeat map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&heartbeat); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workerID := heartbeat["worker_id"].(string)
	tc.mu.Lock()
	if worker, exists := tc.workers[workerID]; exists {
		worker.LastHeartbeat = time.Now()
		worker.Status = heartbeat["status"].(string)
		if load, ok := heartbeat["load"].(float64); ok {
			worker.Load = int(load)
		}
	}
	tc.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (tc *TestCoordinator) assignTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Simplified task assignment
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "no_tasks"})
}

func (tc *TestCoordinator) completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var result DistributedTestResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tc.mu.Lock()
	tc.results[result.TaskID] = &result

	// Update task status
	for _, task := range tc.tasks {
		if task.ID == result.TaskID {
			task.Status = result.Status
			now := time.Now()
			task.CompletedAt = &now
			task.Result = &result
			break
		}
	}
	tc.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (tc *TestCoordinator) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tc.workers)
}

// Close shuts down the distributed test system
func (dts *DistributedTestSystem) Close() {
	dts.cancel()
	if dts.coordinator.server != nil {
		dts.coordinator.server.Shutdown(context.Background())
	}
}
