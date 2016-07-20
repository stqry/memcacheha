package memcacheha

type NodeList struct {
  Nodes map[string]*Node
}

func NewNodeList() *NodeList {
  return &NodeList{
    Nodes: map[string]*Node{},
  }
}

func (me *NodeList) GetHealthyNodes() map[string]*Node {
  out := map[string]*Node{}
  for _, node := range me.Nodes {
    if node.IsHealthy { out[node.Endpoint] = node }
  }
  return out
}

func (me *NodeList) GetHealthyNodeCount() int {
  healthy := 0
  for _, node := range me.Nodes {
    if node.IsHealthy { healthy++ }
  }
  return healthy
}

func (me *NodeList) Exists(nodeAddr string) bool {
  _, found := me.Nodes[nodeAddr]
  return found
}

func (me *NodeList) Add(node *Node) {
  me.Nodes[node.Endpoint] = node
}