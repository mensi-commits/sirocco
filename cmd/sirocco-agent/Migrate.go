package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type RangeSpec struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type MigrateDataCommand struct {
	Cmd          string    `json:"cmd"`
	FromShard    int       `json:"from_shard"`
	ToShard      int       `json:"to_shard"`
	TargetWorker string    `json:"target_worker"` // URL of target worker
	Table        string    `json:"table"`
	Range        RangeSpec `json:"range"`
}

type MigrateDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type MigratedRow map[string]interface{}

type ImportPayload struct {
	Cmd     string        `json:"cmd"`
	ShardID int           `json:"shard_id"`
	Table   string        `json:"table"`
	Rows    []MigratedRow `json:"rows"`
}

// Upgrade migrates a range of rows from one shard to another worker
func Upgrade(w http.ResponseWriter, r *http.Request) {
	var cmd MigrateDataCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "MIGRATE_DATA" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.TargetWorker == "" || cmd.Table == "" {
		http.Error(w, "target_worker and table are required", http.StatusBadRequest)
		return
	}

	if cmd.Range.Start > cmd.Range.End {
		http.Error(w, "invalid range", http.StatusBadRequest)
		return
	}

	sourceDB := fmt.Sprintf("shard_%d", cmd.FromShard)

	query := fmt.Sprintf(
		"SELECT * FROM %s.%s WHERE id >= %d AND id <= %d",
		sourceDB,
		cmd.Table,
		cmd.Range.Start,
		cmd.Range.End,
	)

	rows, err := db.Query(query)
	if err != nil {
		sendJSON(w, MigrateDataResponse{
			Success: false,
			Message: "failed to export rows",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		sendJSON(w, MigrateDataResponse{
			Success: false,
			Message: "failed to read columns",
			Error:   err.Error(),
		})
		return
	}

	exported := []MigratedRow{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))

		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			sendJSON(w, MigrateDataResponse{
				Success: false,
				Message: "failed to scan row",
				Error:   err.Error(),
			})
			return
		}

		rowMap := MigratedRow{}
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		exported = append(exported, rowMap)
	}

	// Send rows to target worker
	payload := ImportPayload{
		Cmd:     "IMPORT_DATA",
		ShardID: cmd.ToShard,
		Table:   cmd.Table,
		Rows:    exported,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(cmd.TargetWorker+"/import", "application/json", bytes.NewBuffer(body))
	if err != nil {
		sendJSON(w, MigrateDataResponse{
			Success: false,
			Message: "failed to send rows to target worker",
			Error:   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		sendJSON(w, MigrateDataResponse{
			Success: false,
			Message: "target worker rejected import",
			Error:   resp.Status,
		})
		return
	}

	// OPTIONAL: delete migrated rows from source shard after successful transfer
	deleteQuery := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE id >= %d AND id <= %d",
		sourceDB,
		cmd.Table,
		cmd.Range.Start,
		cmd.Range.End,
	)

	if _, err := db.Exec(deleteQuery); err != nil {
		sendJSON(w, MigrateDataResponse{
			Success: false,
			Message: "migration succeeded but failed to delete source rows",
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, MigrateDataResponse{
		Success: true,
		Message: fmt.Sprintf("migrated %d rows successfully", len(exported)),
	})
}

// --------------------------
// TARGET SIDE IMPORT HANDLER
// --------------------------

type ImportDataCommand struct {
	Cmd     string        `json:"cmd"`
	ShardID int           `json:"shard_id"`
	Table   string        `json:"table"`
	Rows    []MigratedRow `json:"rows"`
}

type ImportDataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ImportData receives migrated rows and inserts them into local shard
func ImportData(w http.ResponseWriter, r *http.Request) {
	var cmd ImportDataCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "IMPORT_DATA" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.Table == "" {
		http.Error(w, "table is required", http.StatusBadRequest)
		return
	}

	targetDB := fmt.Sprintf("shard_%d", cmd.ShardID)

	tx, err := db.Begin()
	if err != nil {
		sendJSON(w, ImportDataResponse{
			Success: false,
			Message: "failed to begin transaction",
			Error:   err.Error(),
		})
		return
	}

	for _, row := range cmd.Rows {
		// Build dynamic INSERT statement
		cols := []string{}
		vals := []interface{}{}
		placeholders := []string{}

		for k, v := range row {
			cols = append(cols, k)
			vals = append(vals, v)
			placeholders = append(placeholders, "?")
		}

		insertSQL := fmt.Sprintf(
			"INSERT INTO %s.%s (%s) VALUES (%s)",
			targetDB,
			cmd.Table,
			join(cols, ","),
			join(placeholders, ","),
		)

		_, err := tx.Exec(insertSQL, vals...)
		if err != nil {
			tx.Rollback()
			sendJSON(w, ImportDataResponse{
				Success: false,
				Message: "failed to insert row",
				Error:   err.Error(),
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		sendJSON(w, ImportDataResponse{
			Success: false,
			Message: "failed to commit import",
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, ImportDataResponse{
		Success: true,
		Message: fmt.Sprintf("imported %d rows successfully", len(cmd.Rows)),
	})
}

// --------------------------
// SMALL UTILITY
// --------------------------
func join(arr []string, sep string) string {
	out := ""
	for i, s := range arr {
		if i > 0 {
			out += sep
		}
		out += s
	}
	return out
}