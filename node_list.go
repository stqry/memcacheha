package memcacheha

// NodeList represents a list of memcache servers configured/discovered by this client.
type NodeList struct {
	Nodes map[string]*Node
}

// NewNodeList returns a new, empty NodeList
func NewNodeList() *NodeList {
	return &NodeList{
		Nodes: map[string]*Node{},
	}
}

// GetHealthyNodes returns a map of config endpoints to Nodes where the node IsHealthy is true
func (nodeList *NodeList) GetHealthyNodes() map[string]*Node {
	out := map[string]*Node{}
	for _, node := range nodeList.Nodes {
		if node.IsHealthy {
			out[node.Endpoint] = node
		}
	}
	return out
}

// GetHealthyNodeCount returns the count of Nodes where the node IsHealthy is true
func (nodeList *NodeList) GetHealthyNodeCount() int {
	healthy := 0
	for _, node := range nodeList.Nodes {
		if node.IsHealthy {
			healthy++
		}
	}
	return healthy
}

// Exists returns true if a node for the given endpoint exists
func (nodeList *NodeList) Exists(nodeAddr string) bool {
	_, found := nodeList.Nodes[nodeAddr]
	return found
}

// Add the given node to this list
func (nodeList *NodeList) Add(node *Node) {
	nodeList.Nodes[node.Endpoint] = node
}
