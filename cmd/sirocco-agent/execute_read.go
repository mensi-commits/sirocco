package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type WorkerCommand struct {
	Cmd      string `json:"cmd"`
	SQL      string `json:"sql"`
	ShardID  int    `json:"shard_id"`
	ReadOnly bool   `json:"read_only"`
}

type WorkerResponse struct {
	Success      bool                     `json:"success"`
	Message      string                   `json:"message"`
	ShardID      int                      `json:"shard_id"`
	Rows         []map[string]interface{} `json:"rows,omitempty"`
	AffectedRows int64                    `json:"affected_rows,omitempty"`
	Error        string                   `json:"error,omitempty"`
}

// ExecuteRead handles read-only SQL queries (SELECT statements)
// executed on the local shard database.
//
// It is invoked by the Switch layer (XLR8 router) when a read request
// is routed to this worker.
//
// Responsibilities:
//   - Execute SELECT queries against the assigned shard
//   - Format database rows into structured JSON
//   - Return results to the cluster in a consistent response format
//
// Note:
// This function is intended for read operations only. Write operations
// should be handled by ExecuteWrite to maintain shard consistency.

// {
//   "cmd": "EXECUTE_QUERY",
//   "sql": "SELECT * FROM users WHERE id=55",
//   "shard_id": 3,
//   "read_only": true
// }
func ExecuteRead(w http.ResponseWriter, r *http.Request) {
	var cmd WorkerCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "EXECUTE_QUERY" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.SQL == "" {
		http.Error(w, "sql cannot be empty", http.StatusBadRequest)
		return
	}

	// Try SELECT first
	rows, err := db.Query(cmd.SQL)
	if err == nil {
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			sendJSON(w, WorkerResponse{
				Success: false,
				Message: "failed to read columns",
				ShardID: cmd.ShardID,
				Error:   err.Error(),
			})
			return
		}

		results := []map[string]interface{}{}

		for rows.Next() {
			values := make([]interface{}, len(cols))
			ptrs := make([]interface{}, len(cols))

			for i := range values {
				ptrs[i] = &values[i]
			}

			if err := rows.Scan(ptrs...); err != nil {
				sendJSON(w, WorkerResponse{
					Success: false,
					Message: "failed to scan row",
					ShardID: cmd.ShardID,
					Error:   err.Error(),
				})
				return
			}

			row := map[string]interface{}{}
			for i, col := range cols {
				val := values[i]
				if b, ok := val.([]byte); ok {
					row[col] = string(b)
				} else {
					row[col] = val
				}
			}

			results = append(results, row)
		}

		sendJSON(w, WorkerResponse{
			Success: true,
			Message: "query executed successfully",
			ShardID: cmd.ShardID,
			Rows:    results,
		})
		return
	}

	// Fallback: INSERT / UPDATE / DELETE (should ideally not be here in read path)
	res, execErr := db.Exec(cmd.SQL)
	if execErr != nil {
		sendJSON(w, WorkerResponse{
			Success: false,
			Message: "query failed",
			ShardID: cmd.ShardID,
			Error:   execErr.Error(),
		})
		return
	}

	affected, _ := res.RowsAffected()

	sendJSON(w, WorkerResponse{
		Success:      true,
		Message:      "query executed successfully",
		ShardID:      cmd.ShardID,
		AffectedRows: affected,
	})
}

func sendJSON(w http.ResponseWriter, data WorkerResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

var db *sql.DB