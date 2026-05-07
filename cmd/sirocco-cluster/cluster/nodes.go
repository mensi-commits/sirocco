package cluster

import "time"

type Node struct {
	ID        string
	IP        string
	LastSeen  time.Time
	CPU       int
	Status    string
}

func AddNode(n Node) {
	State.mu.Lock()
	defer State.mu.Unlock()

	State.Nodes[n.ID] = &n
}

func UpdateHeartbeat(id string) {
	State.mu.Lock()
	defer State.mu.Unlock()

	if node, ok := State.Nodes[id]; ok {
		node.LastSeen = time.Now()
		node.Status = "online"
	}
}