package memcacheha

import "sync"

// NodeList represents a list of memcache servers configured/discovered by this client.
type NodeList struct {
	lock  sync.RWMutex // on nodes
	nodes map[string]*Node
}

// NewNodeList returns a new, empty NodeList
func NewNodeList() *NodeList {
	return &NodeList{
		lock:  sync.RWMutex{},
		nodes: map[string]*Node{},
	}
}

// GetHealthyNodes returns a map of config endpoints to Nodes where the node IsHealthy is true
func (nodeList *NodeList) GetHealthyNodes() map[string]*Node {
	out := map[string]*Node{}
	nodeList.lock.RLock()
	for _, node := range nodeList.nodes {
		if node.IsHealthy() {
			out[node.Endpoint] = node
		}
	}
	nodeList.lock.RUnlock()
	return out
}

// GetHealthyNodeCount returns the count of Nodes where the node IsHealthy is true
func (nodeList *NodeList) GetHealthyNodeCount() int {
	healthy := 0
	nodeList.lock.RLock()
	for _, node := range nodeList.nodes {
		if node.IsHealthy() {
			healthy++
		}
	}
	nodeList.lock.RUnlock()
	return healthy
}

// Exists returns true if a node for the given endpoint exists
func (nodeList *NodeList) Exists(nodeAddr string) bool {
	nodeList.lock.RLock()
	_, found := nodeList.nodes[nodeAddr]
	nodeList.lock.RUnlock()
	return found
}

// Add the given node to this list
func (nodeList *NodeList) Add(node *Node) {
	nodeList.lock.Lock()
	nodeList.nodes[node.Endpoint] = node
	nodeList.lock.Unlock()
}

func (nodeList *NodeList) SetNodes(nodeAddrs map[string]bool) []string {
	removedAddrs := []string{}
	nodeList.lock.Lock()
	for nodeAddr := range nodeList.nodes {
		if !nodeAddrs[nodeAddr] {
			delete(nodeList.nodes, nodeAddr)
			removedAddrs = append(removedAddrs, nodeAddr)
		}
	}
	nodeList.lock.Unlock()
	return removedAddrs
}

func (nodeList *NodeList) HealthCheck() error {
	nodeList.lock.RLock()
	defer nodeList.lock.RUnlock()
	for _, node := range nodeList.nodes {
		if _, err := node.HealthCheck(); err != nil {
			return err
		}
	}
	return nil
}
