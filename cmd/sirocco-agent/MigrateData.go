package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type MigrateDataCommand struct {
	Cmd        string `json:"cmd"`
	FromShard  int    `json:"from_shard"`
	ToShard    int    `json:"to_shard"`
	TargetURL  string `json:"target_url"`
	Table      string `json:"table"`
	Range      Range  `json:"range"`
}

type MigrateDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type Row map[string]interface{}

type ImportCommand struct {
	Cmd     string `json:"cmd"`
	Table   string `json:"table"`
	ShardID int    `json:"shard_id"`
	Rows    []Row  `json:"rows"`
}

// MigrateData exports a subset of data from one shard and sends it to another worker
func MigrateData(w http.ResponseWriter, r *http.Request) {
	var cmd MigrateDataCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "MIGRATE_DATA" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.TargetURL == "" || cmd.Table == "" {
		http.Error(w, "target_url and table required", http.StatusBadRequest)
		return
	}

	sourceDB := fmt.Sprintf("shard_%d", cmd.FromShard)

	query := fmt.Sprintf(
		"SELECT * FROM %s.%s WHERE id BETWEEN %d AND %d",
		sourceDB,
		cmd.Table,
		cmd.Range.Start,
		cmd.Range.End,
	)

	rows, err := db.Query(query)
	if err != nil {
		sendMigrateJSON(w, MigrateDataResponse{
			Success: false,
			Message: "failed to query source shard",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		sendMigrateJSON(w, MigrateDataResponse{
			Success: false,
			Message: "failed to read columns",
			Error:   err.Error(),
		})
		return
	}

	exported := []Row{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))

		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			sendMigrateJSON(w, MigrateDataResponse{
				Success: false,
				Message: "row scan failed",
				Error:   err.Error(),
			})
			return
		}

		row := Row{}
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		exported = append(exported, row)
	}

	// Send data to target worker
	payload := ImportCommand{
		Cmd:     "IMPORT_DATA",
		Table:   cmd.Table,
		ShardID: cmd.ToShard,
		Rows:    exported,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		cmd.TargetURL+"/import",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		sendMigrateJSON(w, MigrateDataResponse{
			Success: false,
			Message: "failed to send data to target worker",
			Error:   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	sendMigrateJSON(w, MigrateDataResponse{
		Success: true,
		Message: fmt.Sprintf("migrated %d rows successfully", len(exported)),
	})
}

// helper
func sendMigrateJSON(w http.ResponseWriter, data MigrateDataResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// shared DB connection
var db *sql.DB