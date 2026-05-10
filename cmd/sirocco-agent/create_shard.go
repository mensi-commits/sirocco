package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type CreateShardCommand struct {
	Cmd     string `json:"cmd"`
	ShardID int    `json:"shard_id"`
	Port    string `json:"port"`
}

type CreateShardResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	ShardID    int    `json:"shard_id"`
	Container  string `json:"container"`
	Volume     string `json:"volume"`
	DSN        string `json:"dsn"`
	Error      string `json:"error,omitempty"`
}

// CreateShard provisions a new database shard as an isolated Docker container
// with persistent storage attached.
//
// It is a control-plane operation executed by the Shard Manager and triggered
// by the Autoscaler when the system requires horizontal scaling due to load.
//
// Responsibilities:
//   - Create a persistent Docker volume for shard data
//   - Launch a dedicated MySQL container for the shard
//   - Bind the container to a network port for worker access
//   - Generate a DSN (connection string) used by the Switch (XLR8) for routing
//
// Each shard is fully isolated and represents an independent database instance.
// The persistence layer ensures data survives container restarts or failures.
//
// Important:
// This function must NOT be exposed to clients or workers directly.
// It is strictly part of the cluster control plane and must only be invoked
// by internal orchestration components (Autoscaler / Shard Manager).
//
// After creation, the shard must be registered in the cluster metadata store
// so that the Switch can route queries correctly.
func CreateShard(w http.ResponseWriter, r *http.Request) {
	var cmd CreateShardCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "CREATE_SHARD" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	containerName := fmt.Sprintf("sirocco-shard-%d", cmd.ShardID)
	volumeName := fmt.Sprintf("sirocco-vol-%d", cmd.ShardID)

	if cmd.Port == "" {
		cmd.Port = "3306"
	}

	// Step 1: Create Docker volume (persistent storage)
	volCmd := exec.Command("docker", "volume", "create", volumeName)
	if out, err := volCmd.CombinedOutput(); err != nil {
		sendCreateShard(w, CreateShardResponse{
			Success: false,
			Message: "failed to create volume",
			Error:   fmt.Sprintf("%v - %s", err, string(out)),
		})
		return
	}

	// Step 2: Create MySQL container for shard
	runCmd := exec.Command(
		"docker", "run", "-d",
		"--name", containerName,
		"-e", "MYSQL_ROOT_PASSWORD=root",
		"-v", volumeName+":/var/lib/mysql",
		"-p", cmd.Port+":3306",
		"mysql:8.0",
	)

	output, err := runCmd.CombinedOutput()
	if err != nil {
		sendCreateShard(w, CreateShardResponse{
			Success: false,
			Message: "failed to create shard container",
			Error:   fmt.Sprintf("%v - %s", err, string(output)),
		})
		return
	}

	// Step 3: Build DSN (for workers / switch routing)
	dsn := fmt.Sprintf("root:root@tcp(localhost:%s)/shard_%d", cmd.Port, cmd.ShardID)

	sendCreateShard(w, CreateShardResponse{
		Success:   true,
		Message:   "shard created successfully",
		ShardID:   cmd.ShardID,
		Container: containerName,
		Volume:    volumeName,
		DSN:       dsn,
	})
}

func sendCreateShard(w http.ResponseWriter, data CreateShardResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}