package memcacheha

// NodeList represents a list of memcache servers configured/discovered by this client.
type NodeList struct {
	Nodes map[string]*Node
}

// Return a new, empty NodeList
func NewNodeList() *NodeList {
	return &NodeList{
		Nodes: map[string]*Node{},
	}
}

// Return a map of config endpoints to Nodes where the node IsHealthy is true
func (me *NodeList) GetHealthyNodes() map[string]*Node {
	out := map[string]*Node{}
	for _, node := range me.Nodes {
		if node.IsHealthy {
			out[node.Endpoint] = node
		}
	}
	return out
}

// Return the count of Nodes where the node IsHealthy is true
func (me *NodeList) GetHealthyNodeCount() int {
	healthy := 0
	for _, node := range me.Nodes {
		if node.IsHealthy {
			healthy++
		}
	}
	return healthy
}

// Return true if a node for the given endpoint exists
func (me *NodeList) Exists(nodeAddr string) bool {
	_, found := me.Nodes[nodeAddr]
	return found
}

// Add the given node to this list
func (me *NodeList) Add(node *Node) {
	me.Nodes[node.Endpoint] = node
}
