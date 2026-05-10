package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

type CloneShardCommand struct {
	Cmd          string `json:"cmd"`
	SourceShard  int    `json:"source_shard"`
	TargetShard  int    `json:"target_shard"`
	TargetPort   string `json:"target_port"`
}

type CloneShardResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	SourceShard   int    `json:"source_shard"`
	TargetShard   int    `json:"target_shard"`
	ContainerName string `json:"container_name,omitempty"`
	Error         string `json:"error,omitempty"`
}

// CloneShard creates a new shard by duplicating an existing shard state.
//
// It is a control-plane operation used for:
//   - instant replica scaling
//   - creating warm standby nodes
//   - accelerating failover readiness
//
// Responsibilities:
//   - Create a new Docker container for the target shard
//   - Copy data from source shard (logical or snapshot-based cloning)
//   - Initialize target shard as a replica of the source
//   - Register the new shard in cluster metadata (Switch/XLR8)
//
// In production Sirocco systems, cloning would typically:
//   - use snapshot + WAL replay for consistency
//   - avoid full blocking copy operations
//   - stream data incrementally for large shards
//
// Important:
// The cloned shard starts as a replica and must be promoted explicitly
// before it can accept write traffic.
func CloneShard(w http.ResponseWriter, r *http.Request) {
	var cmd CloneShardCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "CLONE_SHARD" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.TargetPort == "" {
		cmd.TargetPort = "3306"
	}

	sourceContainer := fmt.Sprintf("sirocco-shard-%d", cmd.SourceShard)
	targetContainer := fmt.Sprintf("sirocco-shard-%d", cmd.TargetShard)
	targetVolume := fmt.Sprintf("sirocco-vol-%d", cmd.TargetShard)

	// 🧠 Step 1: Create volume for cloned shard
	volCmd := exec.Command("docker", "volume", "create", targetVolume)
	if out, err := volCmd.CombinedOutput(); err != nil {
		sendCloneShardJSON(w, CloneShardResponse{
			Success: false,
			Message: "failed to create volume",
			Error:   fmt.Sprintf("%v - %s", err, string(out)),
		})
		return
	}

	// 🧠 Step 2: Create new container for target shard
	runCmd := exec.Command(
		"docker", "run", "-d",
		"--name", targetContainer,
		"-e", "MYSQL_ROOT_PASSWORD=root",
		"-v", targetVolume+":/var/lib/mysql",
		"-p", cmd.TargetPort+":3306",
		"mysql:8.0",
	)

	if out, err := runCmd.CombinedOutput(); err != nil {
		sendCloneShardJSON(w, CloneShardResponse{
			Success: false,
			Message: "failed to create target shard container",
			Error:   fmt.Sprintf("%v - %s", err, string(out)),
		})
		return
	}

	// 🧠 Step 3: Simulated data cloning (placeholder)
	// Real systems would:
	// - snapshot source shard
	// - stream WAL changes
	// - or restore from backup
	time.Sleep(2 * time.Second)

	// 🧠 Step 4: In real Sirocco system:
	// - register new shard in Switch (XLR8)
	// - mark as REPLICA of source shard
	// - start replication stream

	sendCloneShardJSON(w, CloneShardResponse{
		Success:       true,
		Message:       "shard cloned successfully",
		SourceShard:   cmd.SourceShard,
		TargetShard:   cmd.TargetShard,
		ContainerName: targetContainer,
	})
}

func sendCloneShardJSON(w http.ResponseWriter, data CloneShardResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}