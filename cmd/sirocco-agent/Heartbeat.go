package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type HeartbeatCommand struct {
	Cmd       string `json:"cmd"`
	WorkerID  string `json:"worker_id"`
	Timestamp int64  `json:"timestamp"`
}

type HeartbeatResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	ServerTime int64  `json:"server_time"`
}

// Heartbeat is a lightweight liveness signal sent by the worker
func Heartbeat(w http.ResponseWriter, r *http.Request) {
	var hb HeartbeatCommand

	if err := json.NewDecoder(r.Body).Decode(&hb); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if hb.Cmd != "HEARTBEAT" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if hb.WorkerID == "" {
		http.Error(w, "worker_id cannot be empty", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()

	// Basic sanity check (avoid broken clocks)
	if hb.Timestamp > now+30 || hb.Timestamp < now-300 {
		http.Error(w, "invalid timestamp (clock drift too large)", http.StatusBadRequest)
		return
	}

	// 🔥 In real Sirocco system this would:
	// - update worker last_seen in cluster state
	// - mark node as alive in switch registry
	// - reset failure detection timers
	// - feed autoscaler health graph

	sendHeartbeatJSON(w, HeartbeatResponse{
		Success:    true,
		Message:    "heartbeat received",
		ServerTime: now,
	})
}

func sendHeartbeatJSON(w http.ResponseWriter, data HeartbeatResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}