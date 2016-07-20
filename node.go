package memcacheha

import(
  "github.com/apitalent/memcacheha/log"

  "github.com/bradfitz/gomemcache/memcache"
  
  "crypto/rand"
  "fmt"
  "time"
)

// Node represents a single Memcache server.
type Node struct {
  Endpoint string
  Log log.Logger

  IsHealthy bool
  LastHealthCheck time.Time

  client *memcache.Client
}

// Return a new Node with the given Logger and endpoint (host:port)
func NewNode(logger log.Logger, endpoint string) *Node {
  return &Node{
    Endpoint: endpoint,
    Log: log.NewScopedLogger("Node "+endpoint, logger),
    IsHealthy: false,
    LastHealthCheck: time.Now().Add(-1 * HEALTHCHECK_PERIOD),
    client: memcache.New(endpoint),
  }
} 

// Add an item to the memcache server represented by this node and send the response to the given channel
func (me *Node) Add(item *memcache.Item, finishChan chan(*NodeResponse)) {
  go func(){
    me.Log.Debug("ADD %s", item.Key)
    err := me.client.Add(item)
    if finishChan!=nil { finishChan <- me.getNodeResponse(nil, err) }
  }()
}

// Set an item in the memcache server represented by this node and send the response to the given channel
func (me *Node) Set(item *memcache.Item, finishChan chan(*NodeResponse)) {
  go func(){
    me.Log.Debug("SET %s", item.Key)
    err := me.client.Set(item)
    if finishChan!=nil { finishChan <- me.getNodeResponse(nil, err) }
  }()
}

// Get an item with the given key from the memcache server represented by this node and send the response to the given channel
func (me *Node) Get(key string, finishChan chan(*NodeResponse)) {
  go func(){
    me.Log.Debug("GET %s", key)
    item, err := me.client.Get(key)
    if finishChan!=nil { finishChan <- me.getNodeResponse(item, err) }
  }()
}

// Delete an item with the given key from the memcache server represented by this node and send the response to the given channel
func (me *Node) Delete(key string, finishChan chan(*NodeResponse)) {
  go func(){
    me.Log.Debug("DELETE %s", key)
    err := me.client.Delete(key)
    if finishChan!=nil { finishChan <- me.getNodeResponse(nil, err) }
  }()
}

func (me *Node) Touch(key string, seconds int32, finishChan chan(*NodeResponse)) {
  go func(){
    me.Log.Debug("TOUCH %s", key)
    err := me.client.Touch(key, seconds)
    if finishChan!=nil { finishChan <- me.getNodeResponse(nil, err) }
  }()
}

// Perform a healthcheck on the memcache server represented by this node, update IsHealthy, and return it
func (me *Node) HealthCheck() bool {
  // Read a Random key, expect ErrCacheMiss
  x := make([]byte,32)
  rand.Read(x)
  _, err := me.client.Get(fmt.Sprintf("%02x", x))
  me.getNodeResponse(nil, err)
  return me.IsHealthy
}

func (me *Node) getNodeResponse(item *memcache.Item, err error) *NodeResponse {
  me.LastHealthCheck = time.Now()
  if err != nil && err != memcache.ErrCacheMiss && err != memcache.ErrCASConflict && err != memcache.ErrNotStored && err != memcache.ErrNoStats && err != memcache.ErrMalformedKey { 
    me.markUnhealthy(err) 
  } else { 
    me.markHealthy(err) 
  }
  return NewNodeResponse(me, item, err)
}

func (me *Node) markHealthy(err error) { 
  if !me.IsHealthy { me.Log.Info("Healthy") }
  me.IsHealthy = true
}
func (me *Node) markUnhealthy(err error) {
  if me.IsHealthy { me.Log.Warn("Unhealthy (%s)", err) }
  me.IsHealthy = false
}