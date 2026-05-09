package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type LoadShardCommand struct {
	Cmd        string `json:"cmd"`
	ShardID    int    `json:"shard_id"`
	DataSource string `json:"data_source"`
}

type LoadShardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ShardID int    `json:"shard_id"`
	Error   string `json:"error,omitempty"`
}

// LoadShard restores a shard database from a backup source
func LoadShard(w http.ResponseWriter, r *http.Request) {
	var cmd LoadShardCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "LOAD_SHARD" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.DataSource == "" {
		http.Error(w, "data_source cannot be empty", http.StatusBadRequest)
		return
	}

	shardDB := fmt.Sprintf("shard_%d", cmd.ShardID)

	// Step 1: Create shard database if not exists
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", shardDB))
	if err != nil {
		sendLoadShardJSON(w, LoadShardResponse{
			Success: false,
			Message: "failed to create shard database",
			ShardID: cmd.ShardID,
			Error:   err.Error(),
		})
		return
	}

	// Step 2: Restore data from dump file
	// NOTE: expects local file (downloaded from S3 beforehand if needed)
	restoreCmd := exec.Command(
		"bash", "-c",
		fmt.Sprintf("mysql -u root -p$MYSQL_ROOT_PASSWORD %s < %s", shardDB, cmd.DataSource),
	)

	output, err := restoreCmd.CombinedOutput()
	if err != nil {
		sendLoadShardJSON(w, LoadShardResponse{
			Success: false,
			Message: "failed to restore shard from dump",
			ShardID: cmd.ShardID,
			Error:   fmt.Sprintf("%v - %s", err, string(output)),
		})
		return
	}

	// Step 3: Mark shard ready (in real system: register in cluster metadata)
	sendLoadShardJSON(w, LoadShardResponse{
		Success: true,
		Message: "shard loaded successfully",
		ShardID: cmd.ShardID,
	})
}

// helper
func sendLoadShardJSON(w http.ResponseWriter, data LoadShardResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// shared DB connection
var db *sql.DB