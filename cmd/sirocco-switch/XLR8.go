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

type ShardInfo struct {
	ShardID    int
	HashStart  uint64
	HashEnd    uint64
	Primary    *ShardNode
	Replicas   []ShardNode
	UpdatedAt  time.Time
}

type ShardMap struct {
	Shards    map[int]*ShardInfo
	LoadedAt  time.Time
}


// -----------------------------
// ROUTE RESULT
// -----------------------------
type RouteResult struct {
	ShardID int
	Host    string
	Port    int
	Role    string
	Reason  string
}

// -----------------------------
// XLR8 ROUTER (ShardMap BASED)
// -----------------------------
func XLR8(key any, sm *ShardMap, write bool) (RouteResult, error) {

	if sm == nil || len(sm.Shards) == 0 {
		return RouteResult{}, fmt.Errorf("empty shard map")
	}

	// 1. Hash key
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%v", key)))
	hash := uint64(h.Sum32())

	// 2. Find correct shard by HASH RANGE
	var target *ShardInfo

	for _, shard := range sm.Shards {
		if hash >= shard.HashStart && hash <= shard.HashEnd {
			target = shard
			break
		}
	}

	// fallback (if no range match)
	if target == nil {
		for _, shard := range sm.Shards {
			target = shard
			break
		}
	}

	if target == nil {
		return RouteResult{}, fmt.Errorf("no shard found")
	}

	// 3. Pick best node
	node, role := pickNode(target, write)

	if node == nil {
		return RouteResult{}, fmt.Errorf("no available node in shard %d", target.ShardID)
	}

	return RouteResult{
		ShardID: target.ShardID,
		Host:    node.Host,
		Port:    node.Port,
		Role:    role,
		Reason:  "XLR8 range-based routing",
	}, nil
}

// -----------------------------
// NODE SELECTION
// -----------------------------
func pickNode(s *ShardInfo, write bool) (*ShardNode, string) {

	// WRITE → primary first
	if write {
		if s.Primary != nil && s.Primary.Status == "ONLINE" {
			return s.Primary, "primary"
		}

		for i := range s.Replicas {
			if s.Replicas[i].Status == "ONLINE" {
				return &s.Replicas[i], "replica"
			}
		}

		return nil, ""
	}

	// READ → replicas first
	for i := range s.Replicas {
		if s.Replicas[i].Status == "ONLINE" {
			return &s.Replicas[i], "replica"
		}
	}

	// fallback to primary
	if s.Primary != nil && s.Primary.Status == "ONLINE" {
		return s.Primary, "primary"
	}

	return nil, ""
}

// -----------------------------
// OPTIONAL: helper if needed elsewhere
// -----------------------------
func hashToUint64(v any) uint64 {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	return uint64(h.Sum32())
}

// optional debug helper
func debugShardRange(s *ShardInfo) string {
	return strconv.Itoa(s.ShardID) +
		" [" + fmt.Sprintf("%d", s.HashStart) +
		"-" + fmt.Sprintf("%d", s.HashEnd) + "]"
}