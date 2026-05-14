package cluster

import (
	"context"
	"database/sql"
	"fmt"
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

// LoadClusterMetadata loads shard topology from metadata database.
// It builds a shard map for the switch/router.
func LoadClusterMetadata(ctx context.Context, db *sql.DB) (*ShardMap, error) {
	// Example schema assumptions:
	// shards table:
	//   shard_id INT PK
	//   hash_start BIGINT
	//   hash_end BIGINT
	//   updated_at TIMESTAMP
	//
	// shard_nodes table:
	//   shard_id INT
	//   role VARCHAR (primary/replica)
	//   host VARCHAR
	//   port INT
	//   status VARCHAR
	//   weight INT

	query := `
		SELECT 
			s.shard_id,
			s.hash_start,
			s.hash_end,
			s.updated_at,
			n.role,
			n.host,
			n.port,
			n.status,
			n.weight
		FROM shards s
		JOIN shard_nodes n ON n.shard_id = s.shard_id
		WHERE n.status = 'ONLINE'
		ORDER BY s.shard_id;
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query shard metadata failed: %w", err)
	}
	defer rows.Close()

	shardMap := &ShardMap{
		Shards:   make(map[int]*ShardInfo),
		LoadedAt: time.Now(),
	}

	for rows.Next() {
		var shardID int
		var hashStart uint64
		var hashEnd uint64
		var updatedAt time.Time

		var role string
		var host string
		var port int
		var status string
		var weight int

		err := rows.Scan(
			&shardID,
			&hashStart,
			&hashEnd,
			&updatedAt,
			&role,
			&host,
			&port,
			&status,
			&weight,
		)
		if err != nil {
			return nil, fmt.Errorf("scan shard metadata failed: %w", err)
		}

		if _, ok := shardMap.Shards[shardID]; !ok {
			shardMap.Shards[shardID] = &ShardInfo{
				ShardID:   shardID,
				HashStart: hashStart,
				HashEnd:   hashEnd,
				UpdatedAt: updatedAt,
				Replicas:  []ShardNode{},
			}
		}

		node := ShardNode{
			ShardID: shardID,
			Role:    role,
			Host:    host,
			Port:    port,
			Status:  status,
			Weight:  weight,
		}

		// Assign primary vs replicas
		if role == "primary" {
			shardMap.Shards[shardID].Primary = &node
		} else {
			shardMap.Shards[shardID].Replicas = append(shardMap.Shards[shardID].Replicas, node)
		}

		// keep most recent updated_at
		if updatedAt.After(shardMap.Shards[shardID].UpdatedAt) {
			shardMap.Shards[shardID].UpdatedAt = updatedAt
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return shardMap, nil
}