### 🔥 Workflow: Query Navigation in Sirocco DB (Client → Cluster → Shard)

This is the typical path of a SQL query inside a Sirocco cluster.

---

## 1) Client sends query

The application (client) sends a SQL query like:

```sql
SELECT * FROM users WHERE id = 55;
```

It connects **not directly to MySQL shards**, but to the **Sirocco entry point**:

✅ Sirocco Router / Switch (single endpoint)

Example:

```
client → sirocco-switch:3307
```

---

## 2) Switch receives the SQL query

The **Switch** is the front door of the cluster.

It accepts:

- SQL query
- client connection
- authentication info

---

## 3) SQL Parsing & Query Analysis

The switch parses the query:

- detects query type: SELECT / INSERT / UPDATE / DELETE
- extracts table name: `users`
- extracts key condition: `id = 55`

This is done by a **SQL Parser module**.

---

## 4) Shard key detection

The switch checks if the query contains a shard key.

Example shard key:

- `users.id`

So it sees:

```
id = 55
```

That means it can route correctly.

If shard key is missing → it may broadcast or reject (depends on design).

---

## 5) Shard mapping lookup (Metadata)

The switch consults cluster metadata:

Metadata contains:

- shard ranges
- shard IDs
- which worker owns which shard
- which DB instance is primary/replica

Example metadata rule:

```
Shard 0: id 0 - 9999
Shard 1: id 10000 - 19999
```

So `id=55` → belongs to **Shard 0**.

---

## 6) Node selection

Now the switch decides where to send the query:

- If query is WRITE → go to shard PRIMARY
- If query is READ → can go to a REPLICA (load balancing)

Example:

```
Shard 0 primary = worker-2
Shard 0 replicas = worker-5, worker-6
```

SELECT → choose worker-5 (replica)

---

## 7) Forward query to worker

Switch forwards query + shard info:

```json
{
  "sql": "SELECT * FROM users WHERE id=55",
  "shard_id": 0,
  "target_role": "replica"
}
```

So the worker knows **exactly which database to use**.

---

## 8) Worker executes on shard database

The worker is connected to the actual MySQL instance of that shard:

Example:

```
worker-5 → mysql-shard0-replica:3306
```

Then it executes the query.

---

## 9) Worker returns results to switch

The worker returns:

- rows (for SELECT)
- affected rows (for UPDATE/DELETE)
- inserted ID (for INSERT)
- errors if any

---

## 10) Switch returns response to client

Finally:

```
worker → switch → client
```

The client receives the result like it was a single normal MySQL database.

---

# ⚡ Full Flow Summary (Simple Line)

**Client → Switch (router) → SQL parser → shard lookup → worker selection → shard DB → results → back**

---

# Special Cases

### Case A: Query has no shard key

Example:

```sql
SELECT * FROM users WHERE email='x@gmail.com';
```

Switch cannot know shard directly.

Possible behaviors:

1. ❌ reject (best for performance)
2. 🔁 broadcast query to all shards and merge results (expensive)
3. use secondary index service (advanced)

---

### Case B: JOIN across tables

Example:

```sql
SELECT * FROM users u JOIN orders o ON u.id=o.user_id;
```

If tables are in different shards, the switch must:

- block it
- or execute distributed join (complex + slow)

Most sharded systems **avoid cross-shard joins**.

---

### Case C: INSERT query

```sql
INSERT INTO users(id,name) VALUES(55,'Ali');
```

Same routing steps happen:

- parse
- shard key found
- shard selected
- send to primary shard node

---

# The Brain of Query Navigation

✅ Switch/Router is the **decision maker**
✅ Metadata store is the **map**
✅ Workers are the **executors**
✅ Shards are the **real MySQL databases**

---

If you want, I can draw this as an **architecture diagram** (Client → Switch → Workers → Shards + Metadata).
