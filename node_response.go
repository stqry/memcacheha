package memcacheha

import (
	"github.com/bradfitz/gomemcache/memcache"
)

// NodeResponse represents a reply from a node
type NodeResponse struct {
	Node  *Node
	Item  *memcache.Item
	Error error
}

// NewNodeResponse returns a new NodeResponse with the specified Node, Item and Error
func NewNodeResponse(node *Node, item *memcache.Item, err error) *NodeResponse {
	return &NodeResponse{
		Node:  node,
		Item:  item,
		Error: err,
	}
}
