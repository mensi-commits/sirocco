package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type WorkerRequest struct {
	SQL   string `json:"sql"`
	Type  string `json:"type"`
	Table string `json:"table"`
	Shard string `json:"shard"` 
}

// map shard → DB connection
var shardDBs = map[string]*sql.DB{}

func initShardDBs() {
	// Example shards
	shards := map[string]string{
		"shard_1": "root:pass@tcp(127.0.0.1:3306)/db1",
		"shard_2": "root:pass@tcp(127.0.0.1:3306)/db2",
	}

	for name, dsn := range shards {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("failed to connect shard %s: %v", name, err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("cannot ping shard %s: %v", name, err)
		}

		shardDBs[name] = db
		log.Printf("Connected to %s", name)
	}
}

// -----------------------
// EXECUTION HANDLER
// -----------------------
func executeHandler(w http.ResponseWriter, r *http.Request) {
	var req WorkerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	db, ok := shardDBs[req.Shard]
	if !ok {
		http.Error(w, "unknown shard", http.StatusBadRequest)
		return
	}

	log.Printf("Executing on shard=%s | SQL=%s", req.Shard, req.SQL)

	rows, err := db.Query(req.SQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()

	results := []map[string]any{}

	for rows.Next() {
		cols := make([]any, len(columns))
		ptrs := make([]any, len(columns))

		for i := range cols {
			ptrs[i] = &cols[i]
		}

		rows.Scan(ptrs...)

		row := map[string]any{}
		for i, col := range columns {
			row[col] = cols[i]
		}

		results = append(results, row)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"data":   results,
	})
}

// -----------------------
func main() {
	initShardDBs()

	http.HandleFunc("/execute", executeHandler)

	fmt.Println("Worker running on :9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}