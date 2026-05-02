package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mensi/siroccodb/internal/client"
	"github.com/mensi/siroccodb/internal/protocol"
	"github.com/mensi/siroccodb/internal/sqlparse"
)

/* =========================
   CACHE
========================= */

type cacheEntry struct {
	Route     protocol.RouteResponse
	ExpiresAt time.Time
}

type Switch struct {
    clusterURL string
    cacheTTL   time.Duration

    mu    sync.RWMutex
    cache map[string]cacheEntry

    idMu   sync.Mutex
    nextID int64
}

func NewSwitch(clusterURL string, cacheTTL time.Duration) *Switch {
    log.Printf("[BOOT] Switch cluster=%s cacheTTL=%s", clusterURL, cacheTTL)

    return &Switch{
        clusterURL: strings.TrimRight(clusterURL, "/"),
        cacheTTL:   cacheTTL,
        cache:      make(map[string]cacheEntry),
        nextID:     0,
    }
}

func (s *Switch) generateID() string {
    s.idMu.Lock()
    defer s.idMu.Unlock()

    s.nextID++
    return fmt.Sprintf("%d", s.nextID)
}

func cacheKey(table, key, mode string) string {
	if key == "" {
		key = "*"
	}
	return strings.ToLower(table) + "|" + key + "|" + mode
}

/* =========================
   ROUTE CACHE
========================= */

func (s *Switch) getRoute(table, key, mode string) (protocol.RouteResponse, error) {
	k := cacheKey(table, key, mode)

	s.mu.RLock()
	if entry, ok := s.cache[k]; ok && time.Now().Before(entry.ExpiresAt) {
		log.Printf("[CACHE HIT] %s → %s", k, entry.Route.WorkerID)
		s.mu.RUnlock()
		return entry.Route, nil
	}
	s.mu.RUnlock()

	log.Printf("[CACHE MISS] %s", k)

	url := fmt.Sprintf("%s/route?table=%s&key=%s&mode=%s",
		s.clusterURL, table, key, mode)

	var route protocol.RouteResponse
	if err := client.GetJSON(url, &route); err != nil {
		return protocol.RouteResponse{}, err
	}

	s.mu.Lock()
	s.cache[k] = cacheEntry{
		Route:     route,
		ExpiresAt: time.Now().Add(s.cacheTTL),
	}
	s.mu.Unlock()

	log.Printf("[ROUTE OK] worker=%s addr=%s shard=%s",
		route.WorkerID, route.Address, route.Shard)

	return route, nil
}

func (s *Switch) clearCache() {
	log.Println("[CACHE] cleared")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = make(map[string]cacheEntry)
}

/* =========================
   MAIN
========================= */

func main() {
	addr := flag.String("addr", ":8080", "switch address")
	cluster := flag.String("cluster", "http://localhost:8081", "cluster url")
	cacheTTL := flag.Duration("cache-ttl", 10*time.Second, "cache ttl")
	flag.Parse()

	log.Printf("[BOOT] starting switch on %s", *addr)

	sw := NewSwitch(*cluster, *cacheTTL)

	mux := http.NewServeMux()

	/* ---------- HEALTH ---------- */
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{
			"ok":      true,
			"cluster": sw.clusterURL,
			"cache":   len(sw.cache),
		})
	})

	/* ---------- CACHE CLEAR ---------- */
	mux.HandleFunc("/cache/flush", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}
		sw.clearCache()
		writeJSON(w, 200, map[string]any{"ok": true})
	})

	/* ---------- QUERY ---------- */
	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}

		var req protocol.QueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		log.Printf("[QUERY] %s", req.SQL)

		resp, err := sw.handleQuery(req.SQL)
		if err != nil {
			log.Printf("[QUERY ERROR] %v", err)
			http.Error(w, err.Error(), 400)
			return
		}

		writeJSON(w, 200, resp)
	})

	log.Printf("[READY] listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

/* =========================
   QUERY ENGINE
========================= */


// executed when request to /query is received. It parses the SQL, 
// determines the operation, gets routing info, and forwards the request to the appropriate worker.
func (s *Switch) handleQuery(sql string) (protocol.QueryResponse, error) {
	log.Printf("[PARSE] %s", sql)

	q, err := sqlparse.Parse(sql)

	if err != nil {
		return protocol.QueryResponse{}, err
	}

	b, _ := json.MarshalIndent(q, "", "  ")
	log.Println("[PARSE RESULT]\n", string(b))

	log.Printf("[OP] %s table=%s key=%s", q.Operation, q.Table, q.KeyValue)

	switch q.Operation {

	/* ================= INSERT ================= */
	case sqlparse.OpInsert:
    // Switch ALWAYS generates its own shard key (ID)
    q.KeyValue = s.generateID()
    log.Printf("[ID GEN] generated id=%s for table=%s", q.KeyValue, q.Table)

    route, err := s.getRoute(q.Table, q.KeyValue, "write")
    if err != nil {
        return protocol.QueryResponse{}, err
    }

    resp, err := postExecute(route.Address, protocol.ExecuteRequest{
        Operation: "insert",
        Table:     q.Table,
        Key:       q.KeyValue,
        Columns:   q.Columns,
    })
    if err != nil {
        return protocol.QueryResponse{}, err
    }

    return protocol.QueryResponse{
        OK:       resp.OK,
        Message:  resp.Message,
        Affected: resp.Affected,
        Route:    &route,
    }, nil

	/* ================= SELECT (FIXED) ================= */
	case sqlparse.OpSelect:
		route, err := s.getRoute(q.Table, q.KeyValue, "read")
		if err != nil {
			return protocol.QueryResponse{}, err
		}

		resp, err := postExecute(route.Address, protocol.ExecuteRequest{
			Operation: "select",
			Table:     q.Table,
			Key:       q.KeyValue,
		})
		if err != nil {
			return protocol.QueryResponse{}, err
		}

		return protocol.QueryResponse{
			OK:    resp.OK,
			Row:   resp.Row,
			Rows:  resp.Rows,
			Count: resp.Count,
			Route: &route,
		}, nil

	default:
		log.Printf("[UNSUPPORTED OP] %s", q.Operation)
		return protocol.QueryResponse{}, fmt.Errorf("unsupported operation: %s", q.Operation)
	}
}

/* =========================
   EXECUTE
========================= */
func normalizeURL(addr string) string {
	addr = strings.TrimSpace(addr)

	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return addr
	}

	// :8091 → localhost
	if strings.HasPrefix(addr, ":") {
		return "http://127.0.0.1" + addr
	}

	// 0.0.0.0:8091 → localhost:8091
	if strings.HasPrefix(addr, "0.0.0.0") {
		return "http://127.0.0.1" + strings.TrimPrefix(addr, "0.0.0.0")
	}

	return "http://" + addr
}
func postExecute(addr string, req protocol.ExecuteRequest) (protocol.ExecuteResponse, error) {
	normalized := normalizeURL(addr)
	url := normalized + "/execute"

	log.Printf("[EXECUTE] POST %s op=%s", url, req.Operation)

	var resp protocol.ExecuteResponse
	if err := client.JSON(http.MethodPost, url, req, &resp); err != nil {
		log.Printf("[EXECUTE ERROR] %v", err)
		return protocol.ExecuteResponse{}, err
	}

	return resp, nil
}
/* =========================
   UTIL
========================= */

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
