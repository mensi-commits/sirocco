package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type WriteCommand struct {
	Cmd      string `json:"cmd"`
	SQL      string `json:"sql"`
	ShardID  int    `json:"shard_id"`
	TxID     string `json:"tx_id"`
}

type WriteResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ShardID      int    `json:"shard_id"`
	TxID         string `json:"tx_id"`
	AffectedRows int64  `json:"affected_rows,omitempty"`
	Error        string `json:"error,omitempty"`
}

// ExecuteWrite handles INSERT / UPDATE / DELETE inside a transaction
func ExecuteWrite(w http.ResponseWriter, r *http.Request) {
	var cmd WriteCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "EXECUTE_WRITE" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.SQL == "" {
		http.Error(w, "sql cannot be empty", http.StatusBadRequest)
		return
	}

	// Start transaction (durability boundary)
	tx, err := db.Begin()
	if err != nil {
		sendWriteJSON(w, WriteResponse{
			Success: false,
			Message: "failed to start transaction",
			ShardID: cmd.ShardID,
			TxID:    cmd.TxID,
			Error:   err.Error(),
		})
		return
	}

	// Execute write inside transaction
	res, err := tx.Exec(cmd.SQL)
	if err != nil {
		tx.Rollback()
		sendWriteJSON(w, WriteResponse{
			Success: false,
			Message: "write execution failed",
			ShardID: cmd.ShardID,
			TxID:    cmd.TxID,
			Error:   err.Error(),
		})
		return
	}

	affected, _ := res.RowsAffected()

	// Commit transaction (durability guarantee)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		sendWriteJSON(w, WriteResponse{
			Success: false,
			Message: "transaction commit failed",
			ShardID: cmd.ShardID,
			TxID:    cmd.TxID,
			Error:   err.Error(),
		})
		return
	}

	// Optional: replication hook (handled by cluster later)
	// replicateToFollowers(cmd.ShardID, cmd.SQL)

	sendWriteJSON(w, WriteResponse{
		Success:      true,
		Message:      "write executed successfully",
		ShardID:      cmd.ShardID,
		TxID:         cmd.TxID,
		AffectedRows: affected,
	})
}

// helper
func sendWriteJSON(w http.ResponseWriter, data WriteResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// shared DB connection (this is the shard database)
var db *sql.DB