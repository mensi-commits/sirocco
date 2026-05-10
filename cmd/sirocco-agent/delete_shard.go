package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type DeleteShardCommand struct {
	Cmd     string `json:"cmd"`
	ShardID int    `json:"shard_id"`
}

type DeleteShardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ShardID int    `json:"shard_id"`
	Error   string `json:"error,omitempty"`
}

// DeleteShard permanently removes a shard from the system.
//
// It is a control-plane operation executed by the Shard Manager.
//
// Responsibilities:
//   - Stop and remove the Docker container hosting the shard
//   - Remove the associated persistent storage volume
//   - Clean up shard metadata from the cluster registry
//
// This operation is used during:
//   - downscaling (reducing cluster size)
//   - full shard rebalancing
//   - decommissioning unhealthy nodes
//
// Important:
// This is a destructive operation. Once executed, all data in the shard
// will be lost unless previously migrated or backed up.
func DeleteShard(w http.ResponseWriter, r *http.Request) {
	var cmd DeleteShardCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "DELETE_SHARD" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	containerName := fmt.Sprintf("sirocco-shard-%d", cmd.ShardID)
	volumeName := fmt.Sprintf("sirocco-vol-%d", cmd.ShardID)

	// 🧠 Step 1: Stop and remove Docker container
	stopCmd := exec.Command("docker", "rm", "-f", containerName)
	stopOut, err := stopCmd.CombinedOutput()
	if err != nil {
		sendDeleteShardJSON(w, DeleteShardResponse{
			Success: false,
			Message: "failed to remove container",
			ShardID: cmd.ShardID,
			Error:   fmt.Sprintf("%v - %s", err, string(stopOut)),
		})
		return
	}

	// 🧠 Step 2: Remove persistent volume
	volCmd := exec.Command("docker", "volume", "rm", volumeName)
	volOut, err := volCmd.CombinedOutput()
	if err != nil {
		sendDeleteShardJSON(w, DeleteShardResponse{
			Success: false,
			Message: "failed to remove volume",
			ShardID: cmd.ShardID,
			Error:   fmt.Sprintf("%v - %s", err, string(volOut)),
		})
		return
	}

	// 🧠 Step 3: Remove metadata (pseudo step)
	// switch.DeleteShardMetadata(cmd.ShardID)

	sendDeleteShardJSON(w, DeleteShardResponse{
		Success: true,
		Message: "shard deleted successfully",
		ShardID: cmd.ShardID,
	})
}

func sendDeleteShardJSON(w http.ResponseWriter, data DeleteShardResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}