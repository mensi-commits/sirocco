Creating shards is **not a worker responsibility**. In a well-designed system like Sirocco, shard creation is a **control-plane operation**, because it changes global topology and routing.

---

# 🧠 Who creates shards?

## ⚡ 1. Cluster Autoscaler (decision maker)

Detects:

- high load
- full capacity
- need for scaling

Then decides:

> “We need a new shard”

---

## 🧩 2. Shard Manager (actual creator)

This is the component that **physically creates the shard**.

It:

- allocates a new shard ID
- assigns it to a worker
- defines key range or hash space
- updates metadata

---

## 🔀 3. Metadata / Routing Store (source of truth)

After creation, it updates:

- shard → node mapping
- routing rules
- partition ranges

Without this update, XLR8 would not know the shard exists.

---

## 🖥️ 4. Worker (just executes setup)

Workers only **receive a command to initialize a shard**, they do not decide to create it.

Example command:

```json id="c1"
{
  "cmd": "CREATE_SHARD",
  "shard_id": 7,
  "strategy": "hash",
  "range": {
    "start": 7000,
    "end": 7999
  }
}
```

Worker actions:

- initialize empty database schema
- prepare storage
- optionally restore from snapshot
- start accepting queries

---

# ⚡ Full shard creation flow

## 1. Load detected

```text id="s1"
Shard 3 at 95% CPU
```

## 2. Autoscaler decides

```text id="s2"
Split or create Shard 8
```

## 3. Shard Manager executes

- creates shard definition
- assigns worker node

## 4. Worker initializes shard

- creates DB instance
- loads schema

## 5. Rebalancer moves data

- migrates part of old shard → new shard

## 6. Metadata update

- routing tables updated

## 7. XLR8 adapts automatically

- no code change required

---

# 🧠 Important design principle

> Workers never “create shards”. They only **host shards**.

---

# ⚡ What worker _can_ do related to shards

Workers can only:

### ✅ Initialize shard

- `CREATE_SHARD`

### ✅ Load shard data

- `LOAD_SHARD`

### ✅ Migrate shard data

- `MIGRATE_DATA`

### ❌ NOT allowed:

- deciding when to create shards
- choosing shard boundaries
- updating routing tables

---

# 🚀 Simple mental model

| Component     | Responsibility  |
| ------------- | --------------- |
| Autoscaler    | decides scaling |
| Shard Manager | creates shards  |
| Worker        | hosts shard     |
| Rebalancer    | moves data      |
| XLR8          | routes queries  |

---

# 💥 One-line answer

👉 Shards are created by the **Shard Manager (triggered by the Autoscaler)**, and workers only **initialize and host them**, they never decide to create them.

---

If you want, I can next design the **exact CREATE_SHARD algorithm (how ranges are split + how data migration happens safely without downtime)**—that’s the hardest part of the system.
