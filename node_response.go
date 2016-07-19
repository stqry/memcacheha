package memcacheha

import(
  "github.com/bradfitz/gomemcache/memcache"
)

type NodeResponse struct {
  Node *Node
  Item *memcache.Item
  Error error
}

func NewNodeResponse(node *Node, item *memcache.Item, err error) *NodeResponse {
  return &NodeResponse{
    Node: node,
    Item: item,
    Error: err,
  }
}