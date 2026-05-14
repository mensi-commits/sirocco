package main

import (
	"fmt"

	"encoding/json"
    "net/http"

	"github.com/xwb1989/sqlparser"

	"sirocco-switch/ParseSQL"
	"sirocco-switch/queryHandler"
)

// -------------------------
// Example usage
// -------------------------
func main() {

	http.HandleFunc("/query", queryHandler)
	fmt.Println("Sirocco switch running on :8080")
    http.ListenAndServe("0.0.0.0:8080", nil)

}