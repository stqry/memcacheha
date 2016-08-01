package memcacheha

// NodeResponse represents a reply from a node
type NodeResponse struct {
	Node  *Node
	Item  *Item
	Error error
  NewValue uint64
}

// NewNodeResponse returns a new NodeResponse with the specified Node, Item and Error
func NewNodeResponse(node *Node, item *Item, err error, newValue uint64) *NodeResponse {
	return &NodeResponse{
		Node:  node,
		Item:  item,
		Error: err,
    NewValue: newValue,
	}
}
