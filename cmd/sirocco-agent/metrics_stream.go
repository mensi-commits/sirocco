package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type MetricsStreamCommand struct {
	Cmd string `json:"cmd"`
}

type MetricsStreamData struct {
	QPS           float64 `json:"qps"`            // queries per second
	LatencyMs     float64 `json:"latency_ms"`     // average latency
	DiskUsage     float64 `json:"disk_usage"`     // percentage (0–1)
	CacheHitRatio float64 `json:"cache_hit_ratio"`// percentage (0–1)
	ActiveConns   int     `json:"active_conns"`
}

type MetricsStreamResponse struct {
	Success     bool              `json:"success"`
	Message     string            `json:"message"`
	WorkerID    string            `json:"worker_id"`
	Timestamp   int64             `json:"timestamp"`
	Metrics     MetricsStreamData `json:"metrics"`
}

// MetricsStream continuously reports detailed performance metrics
// from a worker node to the Switch (control plane).
//
// It is used for:
//   - load balancing decisions (XLR8 routing)
//   - autoscaling triggers
//   - detecting hotspots or overloaded shards
//   - performance optimization across replicas
//
// Metrics include:
//   - QPS (queries per second)
//   - average latency
//   - disk usage
//   - cache hit ratio
//   - active connections
//
// In a production system, this would typically be:
//   - streamed periodically (not single-shot HTTP)
//   - pushed to a monitoring system (Prometheus / custom registry)
//   - aggregated by the Switch for global cluster decisions
func MetricsStream(w http.ResponseWriter, r *http.Request) {
	var cmd MetricsStreamCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "METRICS_STREAM" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	// 🧠 In a real system, these would be collected from runtime:
	metrics := MetricsStreamData{
		QPS:           120.5,   // placeholder
		LatencyMs:     4.8,     // placeholder
		DiskUsage:     0.62,    // placeholder
		CacheHitRatio: 0.91,    // placeholder
		ActiveConns:   87,      // placeholder
	}

	now := time.Now().Unix()

	sendMetricsJSON(w, MetricsStreamResponse{
		Success:   true,
		Message:   "metrics snapshot delivered",
		WorkerID:  "worker-unknown", // normally injected from registration context
		Timestamp: now,
		Metrics:   metrics,
	})
}

func sendMetricsJSON(w http.ResponseWriter, data MetricsStreamResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}