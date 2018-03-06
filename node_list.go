package memcacheha

import "sync"

// NodeList represents a list of memcache servers configured/discovered by this client.
type NodeList struct {
	Nodes map[string]*Node
	lock  sync.RWMutex
}

// NewNodeList returns a new, empty NodeList
func NewNodeList() *NodeList {
	return &NodeList{
		Nodes: map[string]*Node{},
		lock:  sync.RWMutex{},
	}
}

// GetHealthyNodes returns a map of config endpoints to Nodes where the node IsHealthy is true
func (nodeList *NodeList) GetHealthyNodes() map[string]*Node {
	out := map[string]*Node{}
	nodeList.lock.RLock()
	for _, node := range nodeList.Nodes {
		node.lock.RLock()
		if node.IsHealthy {
			out[node.Endpoint] = node
		}
		node.lock.RUnlock()
	}
	nodeList.lock.RUnlock()
	return out
}

// GetHealthyNodeCount returns the count of Nodes where the node IsHealthy is true
func (nodeList *NodeList) GetHealthyNodeCount() int {
	healthy := 0
	nodeList.lock.RLock()
	for _, node := range nodeList.Nodes {
		node.lock.RLock()
		if node.IsHealthy {
			healthy++
		}
		node.lock.RUnlock()
	}
	nodeList.lock.RUnlock()
	return healthy
}

// Exists returns true if a node for the given endpoint exists
func (nodeList *NodeList) Exists(nodeAddr string) bool {
	nodeList.lock.RLock()
	_, found := nodeList.Nodes[nodeAddr]
	nodeList.lock.RUnlock()
	return found
}

// Add the given node to this list
func (nodeList *NodeList) Add(node *Node) {
	nodeList.lock.Lock()
	nodeList.Nodes[node.Endpoint] = node
	nodeList.lock.Unlock()
}
