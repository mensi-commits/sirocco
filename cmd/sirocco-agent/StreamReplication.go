package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type StreamReplicationCommand struct {
	Cmd     string `json:"cmd"`
	Source  string `json:"source"` // e.g. "primary"
	Mode    string `json:"mode"`   // "async" or "sync"
}

type StreamReplicationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Mode    string `json:"mode"`
}

// StreamReplication starts a replication stream from primary to replicas
func StreamReplication(w http.ResponseWriter, r *http.Request) {
	var cmd StreamReplicationCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "STREAM_REPLICA" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.Source == "" {
		http.Error(w, "source cannot be empty", http.StatusBadRequest)
		return
	}

	if cmd.Mode == "" {
		cmd.Mode = "async"
	}

	// 🧠 In a real Sirocco system:
	// - this would attach to WAL (write-ahead log)
	// - stream changes to replica workers
	// - maintain offsets / LSN
	// - handle retry + replay

	go func() {
		// Simulated replication loop (replace with WAL stream later)
		for {
			// Example: pull new changes from primary log
			// and push to replicas

			time.Sleep(2 * time.Second)
		}
	}()

	sendReplicationJSON(w, StreamReplicationResponse{
		Success: true,
		Message: "replication stream started",
		Mode:    cmd.Mode,
	})
}

func sendReplicationJSON(w http.ResponseWriter, data StreamReplicationResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}