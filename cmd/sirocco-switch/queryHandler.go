
package main


import (

	"log"
	"encoding/json"
    "net/http"

)


type QueryRequest struct {
    SQL string `json:"sql"`
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 🔥 LOG RAW QUERY (debugging)
	log.Printf("Incoming SQL: %s\n", req.SQL)

	info, err := ParseSQL(req.SQL)
	if err != nil {
		log.Printf("Parse error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 🔥 LOG PARSED RESULT
	log.Printf("Parsed -> Type: %s | Table: %s\n", info.Type, info.Table)

	response := map[string]any{
		"status": "parsed",
		"sql":    req.SQL,
		"type":   info.Type,
		"table":  info.Table,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
