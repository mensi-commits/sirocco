package sirocco

import (
"hash/fnv"
"sort"
)

type QueryInfo struct {
Operation string
Table string
ShardKey string
ShardValue any
}

type Node struct {
ID string
Load float64 // 0.0 → 1.0
Alive bool
IsPrim bool
}

type Shard struct {
ID int
Nodes []Node
}

type ClusterState struct {
Shards []Shard
}

// Final routing result
type ShardRoute struct {
ShardID int
NodeID string
Reason string
}

// -----------------------------
// XLR8 ROUTING ENGINE
// -----------------------------
func XLR8(query string, info QueryInfo, cluster ClusterState) ShardRoute {

    // 1. If shard key is missing → fallback (broadcast-safe strategy)
    if info.ShardKey == "" || info.ShardValue == nil {
    	return fallbackRoute(cluster, "missing shard key → fallback routing")
    }

    // 2. Determine shard using consistent hashing
    shardID := hashShard(info.ShardValue, len(cluster.Shards))
    shard := cluster.Shards[shardID]

    // 3. Filter alive nodes
    aliveNodes := filterAlive(shard.Nodes)
    if len(aliveNodes) == 0 {
    	return fallbackRoute(cluster, "no alive nodes in shard")
    }

    // 4. Decide read/write strategy
    isWrite := info.Operation == "INSERT" ||
    	info.Operation == "UPDATE" ||
    	info.Operation == "DELETE"

    // 5. Select best node
    var selected Node

    if isWrite {
    	// Writes → prefer primary
    	selected = pickPrimary(aliveNodes)
    	if selected.ID == "" {
    		selected = pickLeastLoaded(aliveNodes)
    	}
    } else {
    	// Reads → load-balanced replica selection
    	selected = pickLeastLoaded(aliveNodes)
    }

    return ShardRoute{
    	ShardID: shard.ID,
    	NodeID:  selected.ID,
    	Reason:  "XLR8 routed using hash + health + load strategy",
    }

}

// -----------------------------
// HASH SHARDING (CONSISTENT STYLE)
// -----------------------------
func hashShard(value any, shardCount int) int {
h := fnv.New32a()
h.Write([]byte(toString(value)))
return int(h.Sum32()) % shardCount
}

// -----------------------------
// NODE SELECTION HELPERS
// -----------------------------
func filterAlive(nodes []Node) []Node {
out := []Node{}
for \_, n := range nodes {
if n.Alive {
out = append(out, n)
}
}
return out
}

func pickPrimary(nodes []Node) Node {
for \_, n := range nodes {
if n.IsPrim && n.Alive {
return n
}
}
return Node{}
}

func pickLeastLoaded(nodes []Node) Node {
sort.Slice(nodes, func(i, j int) bool {
return nodes[i].Load < nodes[j].Load
})
return nodes[0]
}

// -----------------------------
// FALLBACK STRATEGY
// -----------------------------
func fallbackRoute(cluster ClusterState, reason string) ShardRoute {
// naive safe fallback: pick first healthy node in cluster
for _, s := range cluster.Shards {
for _, n := range s.Nodes {
if n.Alive {
return ShardRoute{
ShardID: s.ID,
NodeID: n.ID,
Reason: reason,
}
}
}
}
return ShardRoute{ShardID: -1, NodeID: "", Reason: "cluster unavailable"}
}

// -----------------------------
// UTILS
// -----------------------------
func toString(v any) string {
return fmt.Sprintf("%v", v)
}
