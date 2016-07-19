package memcacheha

type NodeList struct {
  nodes map[string]*Node
}

func NewNodeList() *NodeList {
  return &NodeList{
    nodes: map[string]*Node{},
  }
}

func (me *NodeList) GetHealthyNodes() map[string]*Node {
  out := map[string]*Node{}
  for _, node := range me.nodes {
    if node.IsHealthy { out[node.Endpoint] = node }
  }
  return out
}

func (me *NodeList) Exists(nodeAddr string) bool {
  _, found := me.nodes[nodeAddr]
  return found
}

func (me *NodeList) Add(node *Node) {
  me.nodes[node.Endpoint] = node
}