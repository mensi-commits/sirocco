package main

import (
	"encoding/json"
	"net/http"
)

type ImportDataCommand struct {
	Cmd     string `json:"cmd"`
	Table   string `json:"table"`
	ShardID int    `json:"shard_id"`
	Rows    []Row  `json:"rows"`
}

type ImportDataResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ShardID     int    `json:"shard_id"`
	InsertedRows int   `json:"inserted_rows"`
	Error       string `json:"error,omitempty"`
}

// Row represents a generic database row used for migration/import
type Row map[string]interface{}

// ImportData inserts a batch of rows into the local shard database.
//
// It is a core ingestion endpoint used by:
//   - MigrateData (live shard rebalancing)
//   - LoadShard (initial restoration from backup)
//
// Responsibilities:
//   - Receive structured data from another worker or control plane
//   - Insert rows into the correct shard table
//   - Ensure data is applied to the local database instance
//
// In a production Sirocco system, this function would:
//   - use prepared statements or bulk inserts for performance
//   - validate schema consistency before insertion
//   - optionally support idempotency (avoid duplicate imports)
//
// Important:
// This endpoint should NEVER be exposed publicly.
// It must only be called internally by trusted cluster components.
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
		http.Error(w, "table cannot be empty", http.StatusBadRequest)
		return
	}

	if len(cmd.Rows) == 0 {
		sendImportJSON(w, ImportDataResponse{
			Success: false,
			Message: "no rows to import",
			ShardID: cmd.ShardID,
		})
		return
	}

	inserted := 0

	// 🧠 In a real Sirocco system:
	// - use bulk INSERT / COPY for performance
	// - wrap in transaction for consistency
	// - map rows → SQL schema safely

	for range cmd.Rows {
		// simulate insertion
		inserted++
	}

	sendImportJSON(w, ImportDataResponse{
		Success:      true,
		Message:      "data imported successfully",
		ShardID:      cmd.ShardID,
		InsertedRows: inserted,
	})
}

func sendImportJSON(w http.ResponseWriter, data ImportDataResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}