package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/mensi/siroccodb/internal/client"
	"github.com/mensi/siroccodb/internal/protocol"
)

type Worker struct {
	id      string
	addr    string
	role    string
	cluster string
	db      *sql.DB
}

func main() {
	id := flag.String("id", "worker-1", "worker id")
	addr := flag.String("addr", ":8091", "worker listen address")
	cluster := flag.String("cluster", "http://localhost:8081", "cluster base url")
	mysql := flag.String("mysql", "root:root@tcp(127.0.0.1:3306)/sirocco", "mysql dsn")
	role := flag.String("role", "primary", "worker role")
	flag.Parse()

	log.Println("[BOOT] Starting worker...")
	log.Printf("[BOOT] MySQL DSN: %s", *mysql)

	db, err := sql.Open("mysql", *mysql)
	if err != nil {
		log.Fatal("[DB] connection error:", err)
	}

	// test connection immediately
	if err := db.Ping(); err != nil {
		log.Fatal("[DB] ping failed:", err)
	}
	log.Println("[DB] Connected successfully")

	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	w := &Worker{
		id:      *id,
		addr:    *addr,
		role:    *role,
		cluster: *cluster,
		db:      db,
	}

	log.Printf("[WORKER] ID=%s ROLE=%s ADDR=%s", w.id, w.role, w.addr)

	if err := w.register(); err != nil {
		log.Printf("[REGISTER] warning: %v", err)
	} else {
		log.Println("[REGISTER] success")
	}

	go w.heartbeatLoop()

	mux := http.NewServeMux()

	// ================= HEALTH =================
	mux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /health called")
		writeJSON(rw, 200, map[string]any{
			"ok":     true,
			"worker": w.id,
			"role":   w.role,
		})
	})

	// ================= EXECUTE =================
	mux.HandleFunc("/execute", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /execute called")

		if r.Method != http.MethodPost {
			log.Println("[HTTP] invalid method:", r.Method)
			http.Error(rw, "method not allowed", 405)
			return
		}

		var req protocol.ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("[EXECUTE] decode error:", err)
			http.Error(rw, err.Error(), 400)
			return
		}

		log.Printf("[EXECUTE] operation=%s table=%s key=%s", req.Operation, req.Table, req.Key)

		resp := w.execute(req)

		log.Printf("[EXECUTE RESULT] ok=%v affected=%d count=%d",
			resp.OK, resp.Affected, resp.Count)

		writeJSON(rw, 200, resp)
	})

	// ================= STATE =================
	mux.HandleFunc("/state", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("[HTTP] /state called")

		rows, err := w.db.Query("SELECT user_id, data FROM users")
		if err != nil {
			log.Println("[STATE] query error:", err)
			http.Error(rw, err.Error(), 500)
			return
		}
		defer rows.Close()

		data := map[string]any{}

		for rows.Next() {
			var id string
			var jsonData string
			rows.Scan(&id, &jsonData)
			data[id] = jsonData
		}

		log.Printf("[STATE] returned %d rows", len(data))

		writeJSON(rw, 200, map[string]any{
			"worker": w.id,
			"data":   data,
		})
	})

	log.Printf("[SERVER] Worker running on %s", w.addr)
	log.Fatal(http.ListenAndServe(w.addr, mux))
}

// ================= MYSQL EXECUTION =================

func (w *Worker) execute(req protocol.ExecuteRequest) protocol.ExecuteResponse {

	log.Printf("[DB EXEC] %s on table=%s key=%s", req.Operation, req.Table, req.Key)

	switch req.Operation {

	case "insert":
		jsonData, _ := json.Marshal(req.Columns)

		log.Printf("[DB INSERT] key=%s data=%s", req.Key, string(jsonData))

		_, err := w.db.Exec(`
			INSERT INTO users (user_id, data)
			VALUES (?, ?)
			ON DUPLICATE KEY UPDATE data = ?
		`, req.Key, jsonData, jsonData)

		if err != nil {
			log.Println("[DB INSERT ERROR]", err)
			return protocol.ExecuteResponse{OK: false, Message: err.Error()}
		}

		log.Println("[DB INSERT] success")
		return protocol.ExecuteResponse{OK: true, Affected: 1, Message: "inserted"}

	case "select":
		log.Printf("[DB SELECT] key=%s", req.Key)

		row := w.db.QueryRow("SELECT data FROM users WHERE user_id = ?", req.Key)

		var data string
		err := row.Scan(&data)
		if err != nil {
			log.Println("[DB SELECT] not found:", req.Key)
			return protocol.ExecuteResponse{OK: true, Count: 0}
		}

		log.Println("[DB SELECT] found:", req.Key)

		return protocol.ExecuteResponse{
			OK:    true,
			Count: 1,
			Row:   map[string]string{"data": data},
		}

	case "delete":
		log.Printf("[DB DELETE] key=%s", req.Key)

		res, err := w.db.Exec("DELETE FROM users WHERE user_id = ?", req.Key)
		if err != nil {
			log.Println("[DB DELETE ERROR]", err)
			return protocol.ExecuteResponse{OK: false, Message: err.Error()}
		}

		affected, _ := res.RowsAffected()
		log.Printf("[DB DELETE] affected=%d", affected)

		return protocol.ExecuteResponse{OK: true, Affected: int(affected)}

	case "count":
		log.Println("[DB COUNT] called")

		var c int
		w.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&c)

		log.Printf("[DB COUNT] result=%d", c)

		return protocol.ExecuteResponse{OK: true, Count: c}

	default:
		log.Println("[DB EXEC] unsupported operation:", req.Operation)
		return protocol.ExecuteResponse{OK: false, Message: "unsupported"}
	}
}

// ================= CLUSTER =================

func (w *Worker) register() error {
	log.Println("[CLUSTER] registering worker...")

	req := protocol.WorkerRegistrationRequest{
		ID:      w.id,
		Address: w.addr,
		Role:    w.role,
	}

	return client.JSON(http.MethodPost, w.cluster+"/workers/register", req, nil)
}

func (w *Worker) heartbeatLoop() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		log.Println("[CLUSTER] heartbeat sent")

		_ = client.JSON(http.MethodPost, w.cluster+"/workers/heartbeat",
			protocol.HeartbeatRequest{
				ID:      w.id,
				Healthy: true,
				Load:    0,
			}, nil)
	}
}

// ================= HELPERS =================

func writeJSON(wr http.ResponseWriter, status int, v any) {
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(status)
	json.NewEncoder(wr).Encode(v)
}
