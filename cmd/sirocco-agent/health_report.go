package main

import (
	"encoding/json"
	"net/http"
)

type HealthReportCommand struct {
	Cmd               string  `json:"cmd"`
	CPU               float64 `json:"cpu"`
	Memory            float64 `json:"memory"`
	ActiveConnections int     `json:"active_connections"`
	Status            string  `json:"status"`
}

type HealthReportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HealthReport collects and returns the current runtime health state
// of the worker to the cluster control plane.
//
// It is used by the Switch layer and monitoring system to assess
// node health and make routing / scaling decisions.
//
// Reported metrics typically include:
//   - CPU usage
//   - Memory usage
//   - Active database connections
//   - Overall worker status (healthy, degraded, unhealthy)
//
// This function does not modify system state. It is purely observational
// and is called periodically as part of cluster monitoring.
func HealthReport(w http.ResponseWriter, r *http.Request) {
	var report HealthReportCommand

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if report.Cmd != "HEALTH_REPORT" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if report.CPU < 0 || report.CPU > 1 {
		http.Error(w, "invalid cpu value", http.StatusBadRequest)
		return
	}

	if report.Memory < 0 || report.Memory > 1 {
		http.Error(w, "invalid memory value", http.StatusBadRequest)
		return
	}

	if report.Status == "" {
		report.Status = "unknown"
	}

	// In a real Sirocco cluster:
	// - this would update switch/cluster state
	// - feed autoscaler decisions
	// - update node health registry

	sendHealthJSON(w, HealthReportResponse{
		Success: true,
		Message: "health report received",
	})
}

func sendHealthJSON(w http.ResponseWriter, data HealthReportResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}	