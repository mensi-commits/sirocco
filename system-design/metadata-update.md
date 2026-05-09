Good decision splitting it out. Now you’re hitting the **real distributed-systems problem**: _keeping routing metadata consistent between Switch and Cluster._

There are a few standard ways to do it, and the best systems usually combine **push + pull + versioning**.

---

# 🧠 First: what “metadata” means here

This is what must stay in sync:

- shard → node mapping
- key ranges / hash rings
- replica locations
- node health state
- load info (optional for routing)

Think of it as:

> 📍 “the map XLR8 uses to decide where queries go”

---

# 🚀 Option 1 (Best): Central Metadata Store (source of truth)

### Architecture:

```text id="m1"
Cluster Manager → Metadata Store ← Switch
```

Use something like:

- etcd
- Consul
- PostgreSQL (simple version)
- Redis (fast but less strict)

---

## How it works

### 1. Cluster updates metadata

When:

- shard is created
- rebalanced
- node fails

It writes:

```json id="c1"
{
  "shard_id": 3,
  "nodes": ["A", "B"],
  "version": 42
}
```

---

### 2. Switch reads metadata

Two ways:

#### Pull model:

- Switch periodically fetches updates

```text id="p1"
GET /metadata/latest
```

#### Cache + refresh:

- keeps local copy
- refresh every N seconds or on version change

---

### 3. Versioning (VERY IMPORTANT)

Each update has:

```text id="v1"
metadata_version = 42 → 43 → 44
```

Switch only applies newer versions.

---

# ⚡ Option 2: Push-based updates (faster)

### Architecture:

```text id="m2"
Cluster → Switch (stream / websocket / gRPC stream)
```

Cluster pushes:

```json id="c2"
{
  "type": "METADATA_UPDATE",
  "version": 43,
  "changes": [...]
}
```

---

## Pros:

- near real-time routing updates
- no polling delay

## Cons:

- more complex
- needs reconnect handling

---

# 🔥 Option 3: Hybrid (BEST PRACTICE)

Real systems use this:

### 1. Push updates (fast path)

Cluster streams changes to Switch

### 2. Pull fallback (recovery path)

Switch periodically syncs full snapshot

---

# 🧠 How XLR8 should use it

Inside your switch:

```text id="x1"
local_cache ← metadata_store
```

Flow:

1. XLR8 receives query
2. reads local metadata cache
3. routes instantly
4. if version outdated → refresh async

---

# 🧩 Important design rules

## ❌ Don’t:

- query cluster for every request (too slow)
- let workers update metadata directly
- rely on single switch memory without sync

## ✅ Do:

- make metadata **event-driven**
- use **versioned state**
- keep switch **stateless + cache-based**

---

# ⚡ Simple mental model

```text id="m3"
Cluster = writes truth
Metadata Store = holds truth
Switch = cached reader of truth
```

---

# 🏆 Recommended design for Sirocco

Go with:

### ✔ Central metadata store (etcd or DB)

### ✔ Switch local cache

### ✔ Push updates + periodic sync

---

# 💥 Final answer

👉 The cluster and switch should share metadata through a **central versioned metadata store**, with **push updates for speed and pull sync for safety**.

---

If you want next step, I can design:

- the **exact metadata schema for Sirocco**
- or a **Go implementation of the metadata sync system (push + cache + versioning)**
- or even the **full cluster state machine (how shards evolve over time)**
