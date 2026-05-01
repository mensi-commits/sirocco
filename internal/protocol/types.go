package protocol

import "time"

type WorkerRegistrationRequest struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Role    string `json:"role"`
}

type WorkerRegistrationResponse struct {
	OK         bool      `json:"ok"`
	Version    int64     `json:"version"`
	WorkerID   string    `json:"worker_id"`
	Registered time.Time `json:"registered_at"`
}

type HeartbeatRequest struct {
	ID       string `json:"id"`
	Healthy  bool   `json:"healthy"`
	Load     int    `json:"load"`
	Message  string `json:"message,omitempty"`
}

type HeartbeatResponse struct {
	OK      bool  `json:"ok"`
	Version int64 `json:"version"`
}

type RouteRequest struct {
	Table string `json:"table"`
	Key   string `json:"key"`
	Mode  string `json:"mode"` // read | write | count
}

type RouteResponse struct {
	OK         bool      `json:"ok"`
	Version    int64     `json:"version"`
	Table      string    `json:"table"`
	Key        string    `json:"key"`
	Shard      string    `json:"shard"`
	WorkerID   string    `json:"worker_id"`
	Address    string    `json:"address"`
	Role       string    `json:"role"`
	ResolvedAt  time.Time `json:"resolved_at"`
	Cacheable  bool      `json:"cacheable"`
}

type ExecuteRequest struct {
	Operation string            `json:"operation"` // insert | select | update | delete | count
	Table     string            `json:"table"`
	Key       string            `json:"key,omitempty"`
	Columns   map[string]string `json:"columns,omitempty"`
	Updates   map[string]string `json:"updates,omitempty"`
}

type ExecuteResponse struct {
	OK        bool              `json:"ok"`
	Message   string            `json:"message,omitempty"`
	Rows      []map[string]string `json:"rows,omitempty"`
	Row       map[string]string `json:"row,omitempty"`
	Affected  int               `json:"affected,omitempty"`
	Count     int               `json:"count,omitempty"`
	WorkerID  string            `json:"worker_id,omitempty"`
	Table     string            `json:"table,omitempty"`
	Operation string            `json:"operation,omitempty"`
}

type QueryRequest struct {
	SQL string `json:"sql"`
}

type QueryResponse struct {
	OK       bool              `json:"ok"`
	Message  string            `json:"message,omitempty"`
	Rows     []map[string]string `json:"rows,omitempty"`
	Row      map[string]string `json:"row,omitempty"`
	Count    int               `json:"count,omitempty"`
	Affected int               `json:"affected,omitempty"`
	Route    *RouteResponse    `json:"route,omitempty"`
}
