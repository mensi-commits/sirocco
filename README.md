Here is a clean, professional `README.md` without emojis, ready for real project use.

---

# SiroccoDB

A lightweight distributed database system built with Go, featuring a SQL-like query engine, a routing switch layer, a cluster coordinator, and worker nodes backed by MySQL.

---

# Architecture

```
Client / Flask App
        ↓
   Switch (8080)
        ↓
   Cluster (8081)
        ↓
 Workers (8091, 8092, ...)
        ↓
     MySQL
```

---

# Features

- SQL-like query support (INSERT, SELECT, UPDATE, DELETE, COUNT)
- Key-based sharding using `user_id`
- Switch layer with routing cache
- Cluster-based worker routing
- Worker nodes executing SQL against MySQL
- JSON-based internal communication protocol

---

# Requirements

## System

- Go 1.20 or higher
- MySQL or MariaDB 10+
- Linux, WSL, or Windows

## Go dependencies

Run:

```bash
go mod tidy
```

---

# Database Setup

## 1. Start MySQL

```bash
sudo systemctl start mariadb
```

---

## 2. Create database and user

```sql
CREATE DATABASE sirocco;

CREATE USER 'sirocco'@'%' IDENTIFIED BY 'sirocco';
GRANT ALL PRIVILEGES ON sirocco.* TO 'sirocco'@'%';
FLUSH PRIVILEGES;
```

---

## 3. Create table

```sql
USE sirocco;

CREATE TABLE users (
    user_id VARCHAR(64) PRIMARY KEY,
    data JSON
);
```

---

# Running the System

## 1. Start Cluster

```bash
make run-cluster
```

or

```bash
go run ./cmd/cluster -addr 0.0.0.0:8081
```

---

## 2. Start Switch

```bash
make run-switch
```

or

```bash
go run ./cmd/switch -addr 0.0.0.0:8080 -cluster http://127.0.0.1:8081
```

---

## 3. Start Workers

### Worker 1

```bash
make run-worker1
```

or

```bash
go run ./cmd/worker \
  -id worker-1 \
  -addr 0.0.0.0:8091 \
  -cluster http://127.0.0.1:8081 \
  -mysql "sirocco:sirocco@tcp(127.0.0.1:3306)/sirocco"
```

---

### Worker 2

```bash
make run-worker2
```

or

```bash
go run ./cmd/worker \
  -id worker-2 \
  -addr 0.0.0.0:8092 \
  -cluster http://127.0.0.1:8081 \
  -mysql "sirocco:sirocco@tcp(127.0.0.1:3306)/sirocco"
```

---

## 4. Run Flask Dashboard (optional)

```bash
python app.py
```

Access:

```
http://127.0.0.1:5000
```

---

# Testing API

## Insert

```bash
curl -X POST http://127.0.0.1:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql":"INSERT INTO users (user_id, name, email) VALUES (1, \"john\", \"john@mail.com\")"}'
```

---

## Select

```bash
curl -X POST http://127.0.0.1:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql":"SELECT * FROM users WHERE user_id = 1"}'
```

---

## Count

```sql
SELECT COUNT(*) FROM users;
```

---

# Health Checks

## Switch

```bash
curl http://127.0.0.1:8080/health
```

## Cluster

```bash
curl http://127.0.0.1:8081/health
```

---

# Debugging

## Check MariaDB status

```bash
sudo systemctl status mariadb
```

## Check worker connectivity

```bash
curl http://127.0.0.1:8091/health
```

---

# System Design

- Switch handles query parsing and routing
- Cluster decides which worker handles a key
- Workers execute queries on MySQL
- MySQL stores the actual data

---

# Known Issues

- Workers must support all operations defined in switch routing
- MySQL user must have correct privileges for remote access
- Cluster and switch URLs must match network environment

---

# Future Improvements

- Replication between workers
- Automatic shard rebalancing
- Query planner optimization
- Real-time dashboard with WebSockets
- Fault tolerance and failover system

---

If you want next step, I can build:

- production-grade cluster dashboard UI
- live worker monitoring system
- Grafana-style metrics
- automatic failover and replication layer

Here’s your customized version for your repo **`mensi-commits/sirocco`**:

---

## License

This project is open source and available under the [MIT License](LICENSE).

---

 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=mensi-commits/sirocco&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=mensi-commits/sirocco&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=mensi-commits/sirocco&type=date&legend=top-left" />
 </picture>

---

## 📊 Alternative Style View

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=mensi-commits/sirocco&style=landscape1&theme=dark" />
  <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=mensi-commits/sirocco&style=landscape1" />
  <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=mensi-commits/sirocco&style=landscape1" />
</picture>

---

Built with &#9749; by [mensi](https://mensi-commits.github.io/)
