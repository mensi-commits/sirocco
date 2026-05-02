package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mensi/siroccodb/internal/protocol"
)

/* =========================
   CONFIG
========================= */

var UI_ENDPOINT = "http://192.168.1.101:5001/api/events"

var httpClient = &http.Client{
	Timeout: 2 * time.Second,
}

/* =========================
   EVENT SYSTEM
========================= */

type Event struct {
	Type    string      `json:"type"`
	Time    time.Time   `json:"time"`
	Payload interface{} `json:"payload"`
}

// sendEvent sends an event to the UI endpoint asynchronously.
// It logs the event and any errors that occur during sending.
// The event is sent as a JSON payload with the following structure:
// {
//   "type": "event_type",
//   "time": "timestamp",
//   "payload": { ... }
// }
// Events types include: "worker_registered", "heartbeat", "route", etc.

func sendEvent(eventType string, payload any) {
	body, err := json.Marshal(Event{
		Type:    eventType,
		Time:    time.Now().UTC(),
		Payload: payload,
	})

	if err != nil {
		log.Println("[EVENT ERROR] marshal failed:", err)
		return
	}

	log.Println("[EVENT] sending:", eventType)

	go func() {
		req, _ := http.NewRequest("POST", UI_ENDPOINT, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Println("[EVENT ERROR] send failed:", err)
			return
		}
		defer resp.Body.Close()

		log.Println("[EVENT] sent:", eventType, "status:", resp.StatusCode)
	}()
}

/* =========================
   WORKER
========================= */

type Worker struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Role      string    `json:"role"`
	Healthy   bool      `json:"healthy"`
	Load      int       `json:"load"`
	UpdatedAt time.Time `json:"updated_at"`
}

/* =========================
   CLUSTER STATE
========================= */

type Cluster struct {
	mu      sync.RWMutex
	version int64
	workers map[string]Worker
}

func NewCluster() *Cluster {
	return &Cluster{
		workers: make(map[string]Worker),
	}
}

func (c *Cluster) bump() int64 {
	c.version++
	log.Println("[CLUSTER] version bumped:", c.version)
	return c.version
}

/* =========================
   REGISTER
========================= */

func (c *Cluster) register(w Worker) protocol.WorkerRegistrationResponse {
	log.Println("[REGISTER] worker:", w.ID, w.Address)

	c.mu.Lock()
	defer c.mu.Unlock()

	if w.Role == "" {
		w.Role = "primary"
	}

	w.Healthy = true
	w.UpdatedAt = time.Now().UTC()

	c.workers[w.ID] = w

	log.Println("[REGISTER] stored worker:", w.ID)

	sendEvent("worker_registered", w)

	return protocol.WorkerRegistrationResponse{
		OK:         true,
		Version:    c.bump(),
		WorkerID:   w.ID,
		Registered: w.UpdatedAt,
	}
}

/* =========================
   HEARTBEAT
========================= */

func (c *Cluster) heartbeat(req protocol.HeartbeatRequest) protocol.HeartbeatResponse {
	log.Println("[HEARTBEAT] from:", req.ID, "load:", req.Load)

	c.mu.Lock()
	defer c.mu.Unlock()

	w, ok := c.workers[req.ID]
	if !ok {
		log.Println("[HEARTBEAT] unknown worker:", req.ID)
		return protocol.HeartbeatResponse{OK: false, Version: c.bump()}
	}

	w.Healthy = req.Healthy
	w.Load = req.Load
	w.UpdatedAt = time.Now().UTC()

	c.workers[req.ID] = w

	sendEvent("heartbeat", w)

	return protocol.HeartbeatResponse{
		OK:      true,
		Version: c.bump(),
	}
}

/* =========================
   WORKERS
========================= */

func (c *Cluster) allWorkers() []Worker {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]Worker, 0, len(c.workers))
	for _, w := range c.workers {
		out = append(out, w)
	}

	log.Println("[WORKERS] total workers:", len(out))

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})

	return out
}

func (c *Cluster) activeWorkers() []Worker {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := []Worker{}
	for _, w := range c.workers {
		if w.Healthy {
			out = append(out, w)
		}
	}

	log.Println("[WORKERS] active workers:", len(out))

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})

	return out
}

/* =========================
   ROUTING
========================= */

func (c *Cluster) route(table, key, mode string) (protocol.RouteResponse, error) {
	log.Println("[ROUTE] request table:", table, "key:", key, "mode:", mode)

	workers := c.activeWorkers()
	if len(workers) == 0 {
		log.Println("[ROUTE ERROR] no healthy workers")
		return protocol.RouteResponse{}, fmt.Errorf("no healthy workers")
	}

	idx := 0
	if key != "" {
		h := fnv.New32a()
		_, _ = h.Write([]byte(strings.ToLower(table) + ":" + key))
		idx = int(h.Sum32() % uint32(len(workers)))
	}

	w := workers[idx]

	log.Println("[ROUTE] selected worker:", w.ID, "addr:", w.Address)

	c.mu.Lock()
	v := c.bump()
	c.mu.Unlock()

	resp := protocol.RouteResponse{
		OK:         true,
		Version:    v,
		Table:      table,
		Key:        key,
		Shard:      fmt.Sprintf("%s_shard_%d", table, idx),
		WorkerID:   w.ID,
		Address:    w.Address,
		Role:       w.Role,
		ResolvedAt: time.Now().UTC(),
		Cacheable:  true,
	}

	sendEvent("route", resp)

	return resp, nil
}

/* =========================
   JSON
========================= */

func writeJSON(w http.ResponseWriter, status int, v any) {
	log.Println("[HTTP] response status:", status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

/* =========================
   MAIN
========================= */

func main() {
	addr := flag.String("addr", "0.0.0.0:8081", "cluster address")
	flag.Parse()

	cluster := NewCluster()

	mux := http.NewServeMux()

	// HEALTH CHECK ENDPOINT FOR CLUSTER STATUS.
	// Receives GET /health from any source, responds with cluster status and version.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /health")
		writeJSON(w, 200, map[string]any{
			"ok":      true,
			"version": cluster.version,
		})
	})

	// WORKERS ENDPOINT TO LIST ALL REGISTERED WORKERS.
	// Receives GET /workers, responds with list of all workers and their status.
	mux.HandleFunc("/workers", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /workers")
		writeJSON(w, 200, map[string]any{
			"ok":      true,
			"version": cluster.version,
			"workers": cluster.allWorkers(),
		})
	})








	
// Events types include: "worker_registered", "heartbeat", "route", etc.

	// ROUTE ENDPOINT FOR UI TO GET ROUTING INFO.
	mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /route")

		table := r.URL.Query().Get("table")
		key := r.URL.Query().Get("key")
		mode := r.URL.Query().Get("mode")

		resp, err := cluster.route(table, key, mode)
		if err != nil {
			log.Println("[ROUTE ERROR]", err)
			http.Error(w, err.Error(), 500)
			return
		}

		writeJSON(w, 200, resp)
	})


	// HEARTBEAT ENDPOINT FOR WORKERS.
	mux.HandleFunc("/workers/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /workers/heartbeat")

		var req protocol.HeartbeatRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		writeJSON(w, 200, cluster.heartbeat(req))
	})

	mux.HandleFunc("/workers/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /workers/register")

		var req protocol.WorkerRegistrationRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		resp := cluster.register(Worker{
			ID:      req.ID,
			Address: req.Address,
			Role:    req.Role,
		})

		writeJSON(w, 200, resp)
	})


	log.Printf("Cluster running on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
