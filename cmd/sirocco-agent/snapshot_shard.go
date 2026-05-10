package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

type SnapshotShardCommand struct {
	Cmd     string `json:"cmd"`
	ShardID int    `json:"shard_id"`
}

type SnapshotShardResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ShardID     int    `json:"shard_id"`
	SnapshotURI string `json:"snapshot_uri,omitempty"`
	Error       string `json:"error,omitempty"`
}

// SnapshotShard creates a consistent backup of a shard database.
//
// It is a control-plane operation used for:
//   - backups and disaster recovery
//   - shard migration preparation
//   - restoring failed or corrupted shards
//
// Responsibilities:
//   - Freeze or ensure consistency of the shard (conceptually)
//   - Export database contents using mysqldump
//   - Store snapshot in a durable location (e.g. file system or object storage)
//   - Return a snapshot reference (URI)
//
// In a production Sirocco system, this would:
//   - integrate with WAL (Write-Ahead Log) for consistency
//   - store snapshots in S3 / distributed storage
//   - support incremental + full backups
//
// Important:
// Snapshots are read-only and used as recovery or migration checkpoints.
func SnapshotShard(w http.ResponseWriter, r *http.Request) {
	var cmd SnapshotShardCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "SNAPSHOT_SHARD" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	containerName := fmt.Sprintf("sirocco-shard-%d", cmd.ShardID)

	// Snapshot file path (local placeholder for real object storage)
	timestamp := time.Now().Unix()
	snapshotFile := fmt.Sprintf("/tmp/shard_%d_snapshot_%d.sql", cmd.ShardID, timestamp)

	// 🧠 Step 1: Create database dump using mysqldump
	dumpCmd := exec.Command(
		"docker", "exec", containerName,
		"sh", "-c",
		fmt.Sprintf("mysqldump -u root -proot --all-databases > %s", snapshotFile),
	)

	output, err := dumpCmd.CombinedOutput()
	if err != nil {
		sendSnapshotShardJSON(w, SnapshotShardResponse{
			Success: false,
			Message: "failed to create snapshot",
			ShardID: cmd.ShardID,
			Error:   fmt.Sprintf("%v - %s", err, string(output)),
		})
		return
	}

	// 🧠 Step 2: In real systems:
	// - upload snapshotFile to S3 / distributed storage
	// - register snapshot in metadata store
	// - tag with shard version + WAL offset

	snapshotURI := fmt.Sprintf("local://%s", snapshotFile)

	sendSnapshotShardJSON(w, SnapshotShardResponse{
		Success:     true,
		Message:     "snapshot created successfully",
		ShardID:     cmd.ShardID,
		SnapshotURI: snapshotURI,
	})
}

func sendSnapshotShardJSON(w http.ResponseWriter, data SnapshotShardResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}