package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type FailoverNoticeCommand struct {
	Cmd      string `json:"cmd"`
	WorkerID string `json:"worker_id"`
	ShardID  int    `json:"shard_id"`
	Reason   string `json:"reason"` // e.g. "high_latency", "db_failure", "disk_full"
}

type FailoverNoticeResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	WorkerID  string `json:"worker_id"`
	ShardID   int    `json:"shard_id"`
	NotifiedAt int64  `json:"notified_at"`
}

// FailoverNotice is sent by a worker when it detects internal instability
// that may affect query correctness or availability.
//
// It acts as a proactive failure signal to the Switch (XLR8 control plane)
// so the system can immediately:
//   - stop routing traffic to the affected worker
//   - trigger replica promotion (failover)
//   - redistribute shards if necessary
//
// This is different from Heartbeat:
//   - Heartbeat = "I am alive"
//   - FailoverNotice = "I am unhealthy and should be avoided"
//
// Typical triggers:
//   - high latency spikes
//   - database connection failures
//   - disk or memory exhaustion
//   - replication lag overload
func FailoverNotice(w http.ResponseWriter, r *http.Request) {
	var cmd FailoverNoticeCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "FAILOVER_NOTICE" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.WorkerID == "" {
		http.Error(w, "worker_id cannot be empty", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()

	// 🧠 In a real Sirocco system:
	// - Switch immediately marks worker as DEGRADED
	// - routing engine (XLR8) removes it from active pool
	// - autoscaler may trigger replication or failover
	// - metrics system logs incident

	sendFailoverNoticeJSON(w, FailoverNoticeResponse{
		Success:   true,
		Message:   "failover notice received",
		WorkerID:  cmd.WorkerID,
		ShardID:   cmd.ShardID,
		NotifiedAt: now,
	})
}

func sendFailoverNoticeJSON(w http.ResponseWriter, data FailoverNoticeResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}