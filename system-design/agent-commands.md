A worker in your Sirocco-style system is basically a **dumb executor**: it does not decide anything, it only receives **commands from the switch/control plane** and executes them on its local shard DB.

So the key is to define a **clean, minimal command set**.

---

# ⚙️ Core Worker Command Set

## 1. 📥 ExecuteQuery (main command)

This is the most important one.

```json id="c1"
{
  "cmd": "EXECUTE_QUERY",
  "sql": "SELECT * FROM users WHERE id=55",
  "shard_id": 3,
  "read_only": true
}
```

👉 Worker:

- sends SQL to local DB
- returns result

---

## 2. ✍️ ExecuteWrite

For INSERT / UPDATE / DELETE

```json id="c2"
{
  "cmd": "EXECUTE_WRITE",
  "sql": "UPDATE users SET name='Ali' WHERE id=55",
  "shard_id": 3,
  "tx_id": "abc123"
}
```

👉 Worker:

- executes inside transaction
- ensures durability
- may replicate if primary

---

## 3. 🧠 HealthReport

Worker reports status to cluster:

```json id="c3"
{
  "cmd": "HEALTH_REPORT",
  "cpu": 0.72,
  "memory": 0.65,
  "active_connections": 120,
  "status": "healthy"
}
```

---

## 4. 💓 Heartbeat

Lightweight liveness signal:

```json id="c4"
{
  "cmd": "HEARTBEAT",
  "worker_id": "w-12",
  "timestamp": 1710000000
}
```

---

## 5. 📦 LoadShard

Used during scaling / migration:

```json id="c5"
{
  "cmd": "LOAD_SHARD",
  "shard_id": 4,
  "data_source": "s3://backup/shard4.dump"
}
```

👉 Worker:

- restores shard data
- prepares DB instance

---

## 6. 🔄 MigrateData

Used during rebalancing:

```json id="c6"
{
  "cmd": "MIGRATE_DATA",
  "from_shard": 2,
  "to_shard": 5,
  "range": {
    "start": 0,
    "end": 5000
  }
}
```

👉 Worker:

- exports subset of data
- sends to target worker

---

## 7. 🧾 StreamReplication (optional but advanced)

```json id="c7"
{
  "cmd": "STREAM_REPLICA",
  "source": "primary",
  "mode": "async"
}
```

👉 Used for:

- keeping replicas in sync
- near real-time replication

---

## 8. ⛔ Shutdown / Drain

Used for scaling down or maintenance:

```json id="c8"
{
  "cmd": "DRAIN"
}
```

👉 Worker:

- stops accepting new queries
- finishes current ones
- prepares for shutdown or reassignment

---

## 9. 🔁 Reconfigure

Update worker role dynamically:

```json id="c9"
{
  "cmd": "RECONFIGURE",
  "role": "replica",
  "shard_id": 3
}
```

---

# 🧠 Mental model (important)

Workers only do 3 things:

### 🟢 Compute

- execute SQL

### 🟡 Maintain

- heartbeat
- health reporting

### 🔵 Participate in cluster changes

- load shards
- migrate data
- drain/shutdown

---

# ⚡ Minimal “real production” set

If you want the cleanest version:

```text id="m1"
EXECUTE_QUERY
EXECUTE_WRITE
HEARTBEAT
HEALTH_REPORT
DRAIN
MIGRATE_DATA
```

That’s enough for a full sharded DB system.

---

# 🚀 Key design principle

> The worker should NEVER decide anything.

It only:

- receives commands
- executes locally
- reports back

Everything else belongs to:
👉 Switch / Controller / Autoscaler

---

If you want next step, I can design:

- the **Worker internal architecture (Go service design)**
- or the **protocol between Switch ↔ Worker (gRPC schema)**
- or even a **full Sirocco cluster message flow diagram**

Just tell me.
