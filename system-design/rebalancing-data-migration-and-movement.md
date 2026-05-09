Good—this is exactly the situation where a distributed system stops being “just sharding” and becomes a **dynamic resource allocation system**.

If **all shards are full but there are still idle agents available**, the system _does not immediately split shards_. Instead, it goes through a smarter sequence:

---

# 🧠 Core idea

You don’t add shards first.
You first **expand capacity of existing shards using idle agents**.

So the system tries to **convert idle agents → active shard capacity**.

---

# ⚡ Step 1: Elastic expansion (use free agents)

The **Cluster Autoscaler / Shard Manager** will:

- detect high shard saturation (CPU, IO, latency)
- detect available idle agents
- assign those agents as:

### ➜ new replicas for existing shards

Example:

```text id="e1"
Shard 3:
- Primary (full)
- Replica A (full)
→ add Replica B (new agent)
```

👉 This immediately spreads read load + sometimes write buffering pressure.

---

# 🔀 Step 2: Load redistribution (horizontal balancing)

Once new agents join:

The system updates routing:

- XLR8 starts sending **read queries to new replicas**
- hot shards are partially offloaded
- traffic gets redistributed dynamically

👉 This is “scale-out without changing shard count”

---

# 🧩 Step 3: Micro-partitioning (internal shard splitting without full rebalance)

If pressure is still high:

Instead of splitting whole shards immediately, system may:

- create **virtual sub-shards**
- assign them to new agents

Example:

```text id="e2"
Shard 3 → becomes:
  - 3.1 (agent A)
  - 3.2 (agent B)
  - 3.3 (agent C)
```

👉 This is faster than full shard migration

---

# 🧠 Step 4: Promote agents into new shard holders

If load continues increasing:

Now system moves to real scaling:

- create **new shards**
- assign idle agents as owners
- migrate part of data from overloaded shards

👉 This is the real horizontal scaling step

---

# 🔁 Step 5: Rebalancing (data movement phase)

Now the Rebalancer kicks in:

- moves key ranges or hash partitions
- ensures consistency
- avoids downtime (dual write or replication sync)

---

# 🚨 Step 6: Temporary fallback protections (while scaling happens)

During all this:

- XLR8 may temporarily:
  - prefer replicas over primaries
  - route to less-loaded shards even if slightly suboptimal

- system may apply:
  - throttling per client
  - queueing writes

---

# 🧠 Key insight (very important)

The system prefers this order:

### 1. Use existing shards more efficiently

→ add replicas using idle agents

### 2. Increase parallelism inside shards

→ virtual partitions / sub-shards

### 3. Only then increase shard count

→ full split + rebalancing

---

# ⚡ Simple mental model

When load grows and agents are free:

```text id="e3"
Idle Agents
    ↓
Become replicas first
    ↓
Redistribute traffic
    ↓
If still overloaded → split shards
```

---

# 🧭 Who does what?

| Component     | Role                           |
| ------------- | ------------------------------ |
| Monitoring    | detects overload + idle agents |
| Autoscaler    | decides expansion strategy     |
| Shard Manager | assigns agents to shards       |
| Rebalancer    | moves data                     |
| XLR8          | just uses updated routing      |

---

# 💥 One-line answer

👉 If shards are full but agents are available, the system first **expands existing shards by adding replicas and redistributing load**, and only then proceeds to **true shard splitting if needed**.

---

If you want, I can show you the **exact decision algorithm (like a flowchart or Go pseudocode) that decides: replica vs split vs throttle**—that’s the real “brain” of systems like Sirocco.
