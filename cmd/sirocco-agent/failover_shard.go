package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type FailoverShardCommand struct {
	Cmd          string `json:"cmd"`
	ShardID      int    `json:"shard_id"`
	NewPrimaryID string `json:"new_primary_id"` // worker/container id
	ReplicaURL   string `json:"replica_url"`    // target worker to promote
}

type FailoverShardResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	ShardID       int    `json:"shard_id"`
	OldPrimary    string `json:"old_primary,omitempty"`
	NewPrimary    string `json:"new_primary,omitempty"`
	Error         string `json:"error,omitempty"`
}

// FailoverShard performs a shard leadership transition by promoting a replica
// to become the new primary node.
//
// It is used by the cluster control plane when the current primary shard
// becomes unavailable, degraded, or unresponsive.
//
// Responsibilities:
//   - Validate failover request parameters
//   - Trigger role reconfiguration on the selected replica worker
//   - Promote the replica to primary status
//   - Prepare the system for routing updates (handled by Switch/XLR8)
//
// In a complete Sirocco system, this operation is part of the high-availability
// mechanism and is typically triggered automatically by health checks or
// failure detection systems.
//
// Important:
// This function does NOT update routing tables directly. The Switch (XLR8)
// is responsible for propagating the new primary shard information across
// the cluster.
func FailoverShard(w http.ResponseWriter, r *http.Request) {
	var cmd FailoverShardCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "FAILOVER_SHARD" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.ReplicaURL == "" {
		http.Error(w, "replica_url cannot be empty", http.StatusBadRequest)
		return
	}

	if cmd.NewPrimaryID == "" {
		http.Error(w, "new_primary_id cannot be empty", http.StatusBadRequest)
		return
	}

	// 🧠 Step 1: Promote replica container (conceptual switch)
	// In real systems: update replication role + stop slave mode
	promoteCmd := exec.Command(
		"curl", "-X", "POST",
		cmd.ReplicaURL+"/reconfigure",
		"-d", fmt.Sprintf(`{"cmd":"RECONFIGURE","role":"primary","shard_id":%d}`, cmd.ShardID),
	)

	output, err := promoteCmd.CombinedOutput()
	if err != nil {
		sendFailoverJSON(w, FailoverShardResponse{
			Success: false,
			Message: "failed to promote replica",
			ShardID: cmd.ShardID,
			Error:   fmt.Sprintf("%v - %s", err, string(output)),
		})
		return
	}

	// 🧠 Step 2: (Optional real system step)
	// - stop writes on old primary
	// - freeze replication streams
	// - redirect switch routing table

	// 🧠 Step 3: Update cluster metadata (pseudo)
	// switch.UpdateShardPrimary(cmd.ShardID, cmd.NewPrimaryID)

	sendFailoverJSON(w, FailoverShardResponse{
		Success:    true,
		Message:    "failover completed successfully",
		ShardID:    cmd.ShardID,
		NewPrimary: cmd.NewPrimaryID,
	})
}

func sendFailoverJSON(w http.ResponseWriter, data FailoverShardResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}