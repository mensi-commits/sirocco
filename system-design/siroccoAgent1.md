Good question — because right now your agent is already at the point where a single `main.go` will become unmaintainable fast.

In Go, you don’t “class-split” like in Java; you split by **packages + responsibilities**, not just files.

---

# 🧱 Clean way to structure your Sirocco Agent

## 📁 Recommended folder structure

```
sirocco-agent/
│
├── main.go
│
├── config/
│   └── config.go
│
├── api/
│   ├── server.go
│   ├── routes.go
│   └── middleware.go
│
├── node/
│   ├── register.go
│   └── heartbeat.go
│
├── shard/
│   ├── create.go
│   ├── lifecycle.go
│   └── delete.go
│
├── docker/
│   └── client.go
│
├── types/
│   └── models.go
│
└── utils/
    ├── network.go
    └── system.go
```

This gives you **clear separation of responsibilities**, like a real distributed system agent.

---

# 🧠 What goes where

## 1. `main.go` (entry point only)

👉 Only bootstraps everything

```go
package main

import (
	"sirocco-agent/api"
	"sirocco-agent/config"
)

func main() {
	cfg := config.Load()

	server := api.NewServer(cfg)
	server.Start()
}
```

---

## 2. `config/config.go`

```go
package config

import "os"

type Config struct {
	Port string
	Token string
	ControllerURL string
	DataDir string
}

func Load() Config {
	return Config{
		Port: os.Getenv("PORT"),
		Token: os.Getenv("TOKEN"),
		ControllerURL: os.Getenv("CONTROLLER_URL"),
		DataDir: "/var/lib/sirocco",
	}
}
```

---

## 3. `api/server.go`

```go
package api

import (
	"net/http"
	"sirocco-agent/config"
)

type Server struct {
	cfg config.Config
	mux *http.ServeMux
}

func NewServer(cfg config.Config) *Server {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) Start() {
	http.ListenAndServe(":"+s.cfg.Port, s.mux)
}
```

---

## 4. `api/routes.go`

```go
package api

import (
	"net/http"
	"sirocco-agent/node"
	"sirocco-agent/shard"
)

func (s *Server) routes() {

	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	s.mux.HandleFunc("/command", s.auth(shard.HandleCommand(s.cfg)))

	s.mux.HandleFunc("/node/register", node.RegisterHandler(s.cfg))
}
```

---

## 5. `node/heartbeat.go`

```go
package node

import (
	"time"
	"sirocco-agent/config"
)

func StartHeartbeat(cfg config.Config) {
	for {
		// send heartbeat to controller
		time.Sleep(5 * time.Second)
	}
}
```

---

## 6. `shard/create.go`

```go
package shard

import (
	"fmt"
	"os/exec"
	"sirocco-agent/config"
)

func CreateShard(cfg config.Config, name string, port int, password string) error {

	dataPath := fmt.Sprintf("%s/%s", cfg.DataDir, name)

	cmd := exec.Command("docker", "run", "-d",
		"--name", name,
		"-e", "MYSQL_ROOT_PASSWORD="+password,
		"-p", fmt.Sprintf("%d:3306", port),
		"-v", fmt.Sprintf("%s:/var/lib/mysql", dataPath),
		"mysql:8.0",
	)

	return cmd.Run()
}
```

---

## 7. `shard/lifecycle.go`

```go
package shard

import "os/exec"

func Start(name string) error {
	return exec.Command("docker", "start", name).Run()
}

func Stop(name string) error {
	return exec.Command("docker", "stop", name).Run()
}
```

---

## 8. `shard/delete.go`

```go
package shard

import (
	"os"
	"os/exec"
	"sirocco-agent/config"
)

func Delete(cfg config.Config, name string) error {
	exec.Command("docker", "rm", "-f", name).Run()
	return os.RemoveAll(cfg.DataDir + "/" + name)
}
```

---

## 9. `types/models.go`

```go
package types

type Command struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

type ShardConfig struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}
```

---

## 10. `docker/client.go` (optional abstraction)

```go
package docker

import "os/exec"

func Run(args ...string) error {
	return exec.Command("docker", args...).Run()
}
```

---

# 🧠 Why this structure is correct

## ✔ Separation of concerns

- API layer → HTTP only
- Node layer → registration + heartbeat
- Shard layer → MySQL lifecycle
- Docker layer → system execution
- Types → shared models

---

## ✔ Scales like real systems

This is how systems like:

- Kubernetes (kubelet)
- Vitess
- etcd agents

are structured.

---

## ✔ Easy to extend

Later you can add:

```
replication/
monitoring/
backup/
security/
```

without touching core logic.

---

# 🚀 Key idea to remember

> In Go, architecture is not about classes — it's about packages that represent responsibilities.

---

If you want next step, I can help you:

- design the **controller-agent protocol (like Kubernetes API)**
- or build **auto-sharding logic inside controller**
- or implement **replication manager per shard**
