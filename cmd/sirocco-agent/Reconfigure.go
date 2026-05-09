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

// Reconfigure updates the worker role dynamically at runtime
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