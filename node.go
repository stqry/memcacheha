package memcacheha

import (
	"github.com/bradfitz/gomemcache/memcache"

	"crypto/rand"
	"fmt"
	"time"
)

// Node represents a single Memcache server.
type Node struct {
	Endpoint string
	Log      Logger

	IsHealthy       bool
	LastHealthCheck time.Time

	client *memcache.Client
}

// NewNode returns a new Node with the given Logger and endpoint (host:port)
func NewNode(log Logger, endpoint string, timeout time.Duration) *Node {
	node := &Node{
		Endpoint:        endpoint,
		Log:             newScopedLogger("Node "+endpoint, log),
		IsHealthy:       false,
		LastHealthCheck: time.Now().Add(-1 * HEALTHCHECK_PERIOD),
		client:          memcache.New(endpoint),
	}
	node.client.Timeout = timeout
	return node
}

// Add an item to the memcache server represented by this node and send the response to the given channel
func (node *Node) Add(item *Item, finishChan chan (*NodeResponse)) {
	go func() {
		if item.Expiration != nil && !item.Expiration.After(time.Now()) {
			if finishChan != nil {
				finishChan <- NewNodeResponse(node, nil, nil)
			}
			return
		}
		if item.Expiration != nil {
			node.Log.Debug("ADD %s Expire %s", item.Key, *item.Expiration)
		} else {
			node.Log.Debug("ADD %s", item.Key)
		}
		err := node.client.Add(item.AsMemcacheItem())
		if finishChan != nil {
			finishChan <- node.getNodeResponse(nil, err)
		}
	}()
}

// Set an item in the memcache server represented by this node and send the response to the given channel
func (node *Node) Set(item *Item, finishChan chan (*NodeResponse)) {
	go func() {
		if item.Expiration != nil && !item.Expiration.After(time.Now()) {
			if finishChan != nil {
				finishChan <- NewNodeResponse(node, nil, nil)
			}
			return
		}
		if item.Expiration != nil {
			node.Log.Debug("SET %s Expire %s", item.Key, *item.Expiration)
		} else {
			node.Log.Debug("SET %s", item.Key)
		}
		err := node.client.Set(item.AsMemcacheItem())
		if finishChan != nil {
			finishChan <- node.getNodeResponse(nil, err)
		}
	}()
}

// Get an item with the given key from the memcache server represented by this node and send the response to the given channel
func (node *Node) Get(key string, finishChan chan (*NodeResponse)) {
	go func() {
		node.Log.Debug("GET %s", key)
		item, err := node.client.Get(key)
		if finishChan != nil {
			finishChan <- node.getNodeResponse(item, err)
		}
	}()
}

// Delete an item with the given key from the memcache server represented by this node and send the response to the given channel
func (node *Node) Delete(key string, finishChan chan (*NodeResponse)) {
	go func() {
		node.Log.Debug("DELETE %s", key)
		err := node.client.Delete(key)
		if finishChan != nil {
			finishChan <- node.getNodeResponse(nil, err)
		}
	}()
}

// Touch an item with the given key, updating its expiry.
func (node *Node) Touch(key string, seconds int32, finishChan chan (*NodeResponse)) {
	go func() {
		node.Log.Debug("TOUCH %s", key)
		err := node.client.Touch(key, seconds)
		if finishChan != nil {
			finishChan <- node.getNodeResponse(nil, err)
		}
	}()
}

// HealthCheck performs a healthcheck on the memcache server represented by this node, update IsHealthy, and return it
func (node *Node) HealthCheck() (bool, error) {
	// Read a Random key, expect ErrCacheMiss
	x := make([]byte, 32)
	_, err := rand.Read(x)
	if err != nil {
		return false, err
	}
	_, err = node.client.Get(fmt.Sprintf("%02x", x))
	if err != nil && err != memcache.ErrCacheMiss {
		return false, err
	}
	node.getNodeResponse(nil, err)
	return node.IsHealthy, nil
}

func (node *Node) getNodeResponse(item *memcache.Item, err error) *NodeResponse {
	var haitem *Item
	node.LastHealthCheck = time.Now()
	if err != nil &&
		err != memcache.ErrCacheMiss &&
		err != memcache.ErrCASConflict &&
		err != memcache.ErrNotStored &&
		err != memcache.ErrNoStats &&
		err != memcache.ErrMalformedKey {
		node.markUnhealthy(err)
	} else {
		node.markHealthy()
		if item != nil {
			haitem, err = NewItemFromMemcacheItem(item)
		}
	}
	return NewNodeResponse(node, haitem, err)
}

func (node *Node) markHealthy() {
	if !node.IsHealthy {
		node.Log.Info("Healthy")
	}
	node.IsHealthy = true
}
func (node *Node) markUnhealthy(err error) {
	if node.IsHealthy {
		node.Log.Warn("Unhealthy (%s)", err)
	}
	node.IsHealthy = false
}
