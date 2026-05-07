package cluster

import "sync"

var State *ClusterState

type ClusterState struct {
	Nodes  map[string]*Node
	Shards map[string]*Shard
	mu     sync.RWMutex
}

func Init() {
	State = &ClusterState{
		Nodes:  make(map[string]*Node),
		Shards: make(map[string]*Shard),
	}
}