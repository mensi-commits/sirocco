package main

import (
	"fmt"


    "net/http"

	
)

// -------------------------
// Example usage
// -------------------------
func main() {

	http.HandleFunc("/query", queryHandler)
	fmt.Println("Sirocco switch running on :8080")
    http.ListenAndServe("0.0.0.0:8080", nil)

}