package api

import (
	"encoding/json"
	"net/http"
	"sirocco-cluster/cluster"
	"sirocco-cluster/agent"
	"time"
)

type RegisterRequest struct {
	ID string `json:"id"`
	IP string `json:"ip"`
}

func registerNode(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	cluster.AddNode(cluster.Node{
		ID:       req.ID,
		IP:       req.IP,
		LastSeen: time.Now(),
		Status:   "online",
	})

	w.Write([]byte(`{"status":"registered"}`))
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	cluster.UpdateHeartbeat(req.ID)

	w.Write([]byte(`{"status":"ok"}`))
}