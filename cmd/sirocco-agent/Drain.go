package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type DrainCommand struct {
	Cmd string `json:"cmd"`
}

type DrainResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	State   string `json:"state"`
}

// global worker state
var acceptingRequests int32 = 1 // 1 = active, 0 = draining

var inFlight sync.WaitGroup

// Drain stops the worker from accepting new queries and waits for active ones
func Drain(w http.ResponseWriter, r *http.Request) {
	var cmd DrainCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "DRAIN" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	// ⛔ Stop accepting new requests
	atomic.StoreInt32(&acceptingRequests, 0)

	// 🧠 In real Sirocco:
	// - switch stops routing queries here
	// - worker is marked "draining" in cluster metadata
	// - autoscaler avoids assigning new load

	// ⏳ Wait for in-flight queries to finish
	done := make(chan struct{})

	go func() {
		inFlight.Wait()
		close(done)
	}()

	select {
	case <-done:
		// all queries finished
	case <-time.After(30 * time.Second):
		// safety timeout (force drain)
	}

	sendDrainJSON(w, DrainResponse{
		Success: true,
		Message: "worker drained successfully",
		State:   "drained",
	})
}

// helper to check before executing any query
func AllowRequest() bool {
	return atomic.LoadInt32(&acceptingRequests) == 1
}

func sendDrainJSON(w http.ResponseWriter, data DrainResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}