package cluster

func PickNode() *Node {
	State.mu.RLock()
	defer State.mu.RUnlock()

	var selected *Node

	for _, n := range State.Nodes {
		if selected == nil || n.LastSeen.After(selected.LastSeen) {
			selected = n
		}
	}

	return selected
}