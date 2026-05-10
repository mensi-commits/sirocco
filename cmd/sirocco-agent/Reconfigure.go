package main

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type ReconfigureCommand struct {
	Cmd     string `json:"cmd"`
	Role    string `json:"role"`     // "primary", "replica", "readonly"
	ShardID int    `json:"shard_id"`
}

type ReconfigureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Role    string `json:"role"`
	ShardID int    `json:"shard_id"`
}

// global worker state (dynamic role)
var workerRole atomic.Value

// Reconfigure dynamically updates the role of a worker at runtime
// without requiring a restart.
//
// It is used by the Switch (control plane) to change how a worker
// participates in the Sirocco cluster.
//
// Supported roles:
//   - primary   : handles write operations and authoritative data
//   - replica   : receives replicated data and serves read traffic
//   - readonly  : serves read-only queries without accepting writes
//
// Responsibilities:
//   - Validate the requested role
//   - Atomically update the worker's runtime role
//   - Allow the Switch to immediately adjust routing decisions
//
// This mechanism enables live topology changes such as:
//   - failover (replica → primary)
//   - load redistribution
//   - read/write role separation
//
// Note:
// Role changes take effect immediately and may influence query routing,
// replication behavior, and write acceptance policies.
func Reconfigure(w http.ResponseWriter, r *http.Request) {
	var cmd ReconfigureCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "RECONFIGURE" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.Role == "" {
		http.Error(w, "role cannot be empty", http.StatusBadRequest)
		return
	}

	// validate allowed roles
	switch cmd.Role {
	case "primary", "replica", "readonly":
		// valid
	default:
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	// update runtime role (hot swap)
	workerRole.Store(cmd.Role)

	// 🧠 In real Sirocco system:
	// - switch updates routing rules (XLR8)
	// - replica may start applying WAL stream
	// - primary may enable write mode
	// - readonly blocks writes

	sendReconfigureJSON(w, ReconfigureResponse{
		Success: true,
		Message: "worker reconfigured successfully",
		Role:    cmd.Role,
		ShardID: cmd.ShardID,
	})
}

// helper
func GetWorkerRole() string {
	if v := workerRole.Load(); v != nil {
		return v.(string)
	}
	return "unknown"
}

func sendReconfigureJSON(w http.ResponseWriter, data ReconfigureResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}