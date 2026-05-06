package main

import (
	"fmt"

	"encoding/json"
    "net/http"

	"github.com/xwb1989/sqlparser"
)

// QueryInfo holds extracted routing metadata
type QueryInfo struct {
	Type  string
	Table string
}


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

	info, err := ParseSQL(req.SQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"status": "parsed",
		"sql":    req.SQL,
		"type":   info.Type,
		"table":  info.Table,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// ParseSQL analyzes a SQL query and extracts routing info
func ParseSQL(query string) (*QueryInfo, error) {

	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("invalid SQL: %v", err)
	}

	info := &QueryInfo{
		Type:  "UNKNOWN",
		Table: "",
	}

	switch v := stmt.(type) {

	// SELECT
	case *sqlparser.Select:
		info.Type = "SELECT"
		info.Table = extractTableFromSelect(v.From)

	// INSERT
	case *sqlparser.Insert:
		info.Type = "INSERT"
		info.Table = v.Table.Name.String()

	// UPDATE
	case *sqlparser.Update:
		info.Type = "UPDATE"
		info.Table = extractTableGeneric(v.TableExprs)

	// DELETE
	case *sqlparser.Delete:
		info.Type = "DELETE"
		info.Table = extractTableGeneric(v.TableExprs)

	default:
		info.Type = "UNKNOWN"
	}

	return info, nil
}

// Extract table from SELECT
func extractTableFromSelect(from sqlparser.TableExprs) string {
	for _, expr := range from {
		if aliased, ok := expr.(*sqlparser.AliasedTableExpr); ok {
			if tbl, ok := aliased.Expr.(sqlparser.TableName); ok {
				return tbl.Name.String()
			}
		}
	}
	return ""
}

// Extract table from UPDATE/DELETE
func extractTableGeneric(exprs sqlparser.TableExprs) string {
	for _, expr := range exprs {
		if aliased, ok := expr.(*sqlparser.AliasedTableExpr); ok {
			if tbl, ok := aliased.Expr.(sqlparser.TableName); ok {
				return tbl.Name.String()
			}
		}
	}
	return ""
}

// -------------------------
// Example usage
// -------------------------
func main() {

	http.HandleFunc("/query", queryHandler)
	fmt.Println("Sirocco switch running on :8080")
    http.ListenAndServe(":8080", nil)

}