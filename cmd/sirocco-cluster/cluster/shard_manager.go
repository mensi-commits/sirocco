package cluster

import (
	"sirocco-cluster/agent"
)

func CreateShard(shardID string) error {
	node := PickNode()

	AddShard(Shard{
		ID:     shardID,
		NodeID: node.ID,
		Port:   3306,
		Status: "deploying",
	})

	payload := map[string]any{
		"action": "create_shard",
		"data": map[string]any{
			"shard_id": shardID,
			"port":     3306,
		},
	}

	return agent.SendCommand(node.IP, payload)
}