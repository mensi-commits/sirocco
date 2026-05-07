package cluster

type Shard struct {
	ID     string
	NodeID string
	Port   int
	Status string
}

func AddShard(s Shard) {
	State.mu.Lock()
	defer State.mu.Unlock()

	State.Shards[s.ID] = &s
}