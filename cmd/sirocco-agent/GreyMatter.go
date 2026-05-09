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

// GrayMatter loads/restores a shard from a backup source (scaling/migration)
func GrayMatter(w http.ResponseWriter, r *http.Request) {
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

	// NOTE: For now we only support local dump files for security reasons.
	// Example supported input: "/backups/shard4.sql"
	// If you want S3 support, you must download it first using AWS CLI or SDK.
	if len(cmd.DataSource) > 5 && cmd.DataSource[:5] == "s3://" {
		sendJSON(w, LoadShardResponse{
			Success: false,
			Message: "S3 sources not supported yet (download dump locally first)",
			ShardID: cmd.ShardID,
			Error:   "unsupported data_source",
		})
		return
	}

	// Create shard database if not exists
	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS shard_%d", cmd.ShardID)
	if _, err := db.Exec(createDB); err != nil {
		sendJSON(w, LoadShardResponse{
			Success: false,
			Message: "failed to create shard database",
			ShardID: cmd.ShardID,
			Error:   err.Error(),
		})
		return
	}

	// Restore dump into shard DB using mysql CLI (fast + reliable)
	// mysql -u root -pPASS shard_4 < dump.sql
	restoreCmd := exec.Command(
		"bash", "-c",
		fmt.Sprintf("mysql -u root -p%s shard_%d < %s", mysqlRootPassword, cmd.ShardID, cmd.DataSource),
	)

	out, err := restoreCmd.CombinedOutput()
	if err != nil {
		sendJSON(w, LoadShardResponse{
			Success: false,
			Message: "failed to restore shard dump",
			ShardID: cmd.ShardID,
			Error:   fmt.Sprintf("%v - %s", err, string(out)),
		})
		return
	}

	sendJSON(w, LoadShardResponse{
		Success: true,
		Message: "shard loaded successfully",
		ShardID: cmd.ShardID,
	})
}