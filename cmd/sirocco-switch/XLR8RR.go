

package main

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"time"
)


type ShardNode struct {
	ShardID int
	Role    string // "primary" or "replica"
	Host    string
	Port    int
	Status  string // "ONLINE", "OFFLINE"
	Weight  int
}

type ShardMap struct {
	Shards    map[int]*ShardInfo
	LoadedAt  time.Time
}


type ShardInfo struct {
	ShardID   int
	HashStart uint64
	HashEnd   uint64
	Primary   *ShardNode
	Replicas  []ShardNode
	UpdatedAt time.Time

	// 🔥 NEW: round robin pointer
	RRIndex int
}





// XLR8 routes a query to the correct shard and node
// based on hash-based sharding + round-robin node selection.
func XLR8(key any, sm *ShardMap, write bool) (RouteResult, error) {

	// -----------------------------
	// Validate shard map existence
	// -----------------------------
	if sm == nil || len(sm.Shards) == 0 {
		return RouteResult{}, fmt.Errorf("empty shard map")
	}

	// -----------------------------
	// 1. Hash the key (deterministic routing input)
	// -----------------------------
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%v", key)))
	hash := uint64(h.Sum32())

	// -----------------------------
	// 2. Find the shard matching the hash range
	// -----------------------------
	var target *ShardInfo

	for _, shard := range sm.Shards {
		if hash >= shard.HashStart && hash <= shard.HashEnd {
			target = shard
			break
		}
	}

	// -----------------------------
	// 3. Fallback if no range matched
	// (prevents routing failure on misconfigured ranges)
	// -----------------------------
	if target == nil {
		for _, shard := range sm.Shards {
			target = shard
			break
		}
	}

	if target == nil {
		return RouteResult{}, fmt.Errorf("no shard found")
	}

	// -----------------------------
	// 4. Select best node inside shard
	// (uses round robin for replicas)
	// -----------------------------
	node, role := pickNodeRR(target, write)
	if node == nil {
		return RouteResult{}, fmt.Errorf("no available node in shard %d", target.ShardID)
	}

	// -----------------------------
	// 5. Return routing decision
	// -----------------------------
	return RouteResult{
		ShardID: target.ShardID,
		Host:    node.Host,
		Port:    node.Port,
		Role:    role,
		Reason:  "XLR8 range + round-robin routing",
	}, nil
}




// pickNodeRR decides which node inside a shard should handle the query.
// - WRITE: prefers primary, falls back to replicas (RR)
// - READ: prefers replicas (RR), falls back to primary
func pickNodeRR(s *ShardInfo, write bool) (*ShardNode, string) {

	// -----------------------------
	// WRITE PATH
	// -----------------------------
	if write {
		// Primary is always preferred for writes
		if s.Primary != nil && s.Primary.Status == "ONLINE" {
			return s.Primary, "primary"
		}

		// If primary is down, use replicas (round robin)
		return pickReplicaRR(s)
	}

	// -----------------------------
	// READ PATH
	// -----------------------------

	// Try replicas first using round robin
	if node := pickReplicaRR(s); node != nil {
		return node, "replica"
	}

	// Fallback to primary if no replicas available
	if s.Primary != nil && s.Primary.Status == "ONLINE" {
		return s.Primary, "primary"
	}

	return nil, ""
}



// pickReplicaRR selects an ONLINE replica using round robin.
// It keeps state in shard.RRIndex to ensure fairness across requests.
func pickReplicaRR(s *ShardInfo) *ShardNode {

	n := len(s.Replicas)

	// No replicas available
	if n == 0 {
		return nil
	}

	// Starting point for round robin cycle
	start := s.RRIndex

	// Try all replicas starting from RR position
	for i := 0; i < n; i++ {

		// Circular index calculation
		idx := (start + i) % n

		// Only select ONLINE nodes
		if s.Replicas[idx].Status == "ONLINE" {

			// Move RR pointer forward for next request
			s.RRIndex = (idx + 1) % n

			return &s.Replicas[idx]
		}
	}

	// No available replicas found
	return nil
}