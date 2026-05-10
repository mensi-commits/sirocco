package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type RegisterWorkerCommand struct {
	Cmd        string `json:"cmd"`
	WorkerID   string `json:"worker_id"`
	Capacity   int    `json:"capacity"`   // max connections / load capacity
	ShardCount int    `json:"shard_count"`
	Role       string `json:"role"`       // worker role: replica / primary / readonly
}

type RegisterWorkerResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	WorkerID    string `json:"worker_id"`
	RegisteredAt int64  `json:"registered_at"`
}

// RegisterWorker registers a worker node inside the Sirocco cluster.
//
// It is the initial handshake between a worker and the Switch (control plane).
// Without registration, a worker is not eligible to receive queries or shard assignments.
//
// Responsibilities:
//   - Validate worker identity and metadata
//   - Register worker in the cluster registry (Switch-side state)
//   - Assign role, capacity, and initial shard allocation metadata
//   - Mark worker as active and available for routing
//
// In a production Sirocco system:
//   - Switch stores worker in a distributed registry (e.g. etcd / raft log)
//   - worker becomes part of routing decisions (XLR8)
//   - load balancer begins assigning shards to it
//
// This function must be called before any other worker operations
// such as ExecuteRead, ExecuteWrite, or Heartbeat.
func RegisterWorker(w http.ResponseWriter, r *http.Request) {
	var cmd RegisterWorkerCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "REGISTER_WORKER" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.WorkerID == "" {
		http.Error(w, "worker_id cannot be empty", http.StatusBadRequest)
		return
	}

	if cmd.Role == "" {
		cmd.Role = "replica"
	}

	now := time.Now().Unix()

	// 🧠 In real Sirocco system:
	// - persist worker in Switch registry
	// - assign worker metadata (capacity, shards, role)
	// - mark worker as ACTIVE
	// - include in XLR8 routing pool

	sendRegisterWorkerJSON(w, RegisterWorkerResponse{
		Success:      true,
		Message:      "worker registered successfully",
		WorkerID:     cmd.WorkerID,
		RegisteredAt: now,
	})
}

func sendRegisterWorkerJSON(w http.ResponseWriter, data RegisterWorkerResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}