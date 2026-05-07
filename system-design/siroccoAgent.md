The **agent** is the software running on every VPS that makes the VPS a **managed node** in the Sirocco cluster.

Think of it like:

- **Controller = brain**
- **Switch = traffic router**
- **Agent = hands/executor**

Without the agent, your controller cannot safely create shards, start containers, monitor MySQL, or perform failover.

---

# What the Agent does (real role)

## 1. Node Registration (Join the cluster)

When the agent starts, it connects to the controller and says:

> "I am VPS-12, my IP is X, I have 8 CPU, 16GB RAM, 200GB disk."

So the controller adds it into the cluster topology.

---

## 2. Deploy and Manage Shard Containers

The controller sends commands like:

- "create shard_5 container on port 3310"
- "start shard_5"
- "stop shard_5"
- "restart shard_5"

The agent executes those commands locally using Docker.

So **agent is responsible for MySQL containers lifecycle**.

---

## 3. Persistent Storage Handling

Agent ensures shard data is stored in:

`/var/lib/sirocco/shards/shard_5`

So containers can be deleted/recreated without losing data.

---

## 4. Configure Replication

If controller says:

> "make shard_5 a replica of shard_5 primary"

Agent will execute MySQL replication setup:

- connect to primary
- run `CHANGE MASTER TO`
- run `START REPLICA`
- verify replication status

So replication is automated.

---

## 5. Health Monitoring + Metrics Reporting

Agent continuously monitors:

- MySQL alive or down
- replication lag
- disk usage
- CPU load
- memory usage
- container status

Then sends periodic heartbeats to controller:

```json
{
  "node": "vps-12",
  "shard": "shard_5",
  "mysql_status": "healthy",
  "lag_seconds": 2,
  "disk_used": "40GB"
}
```

---

## 6. Failover Execution

When primary dies, controller decides:

> "promote replica shard_5 on VPS-12"

Agent executes locally:

- stop replication
- disable read_only
- restart MySQL if needed

Then controller updates the switch routing.

---

# How Agent interacts with the cluster

## Interaction Flow (clean design)

### 1. Agent → Controller (heartbeat)

Agent sends every 5 seconds:

- node health
- shard health
- metrics

### 2. Controller → Agent (commands)

Controller sends commands like:

- deploy shard container
- configure replication
- run backup
- restore shard
- promote replica

### 3. Switch ← Controller (routing updates)

Controller updates the routing table of the switch:

- shard locations
- primary/replica endpoints

So the switch always routes queries to the correct place.

---

# Typical full workflow example

### Admin adds VPS

1. Admin installs `sirocco-agent` on VPS.
2. Agent calls controller `/register`.

### Controller uses it

3. Controller says: "create shard_8 replica here".
4. Agent deploys mysql container.
5. Agent configures replication.
6. Agent reports status = ready.

### Cluster updated

7. Controller updates metadata DB.
8. Switch routing is updated.

---

# Agent Responsibilities Summary (simple)

The agent is responsible for:

✅ container lifecycle (start/stop/create/remove)
✅ volume management (persistent storage)
✅ mysql setup (users, configs, ports)
✅ replication setup
✅ health checks and metrics
✅ failover execution
✅ backup and restore (optional)

---

# Why this is the correct architecture

Because the controller should never directly run system commands on VPS.
The agent is the secure local daemon that does it.

---

If you want, I can now design the **full communication protocol** between controller and agent (endpoints + JSON format) like a real product.
