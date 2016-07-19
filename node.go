package memcacheha

import(
  "github.com/bradfitz/gomemcache/memcache"
  
  "crypto/rand"
  "fmt"
  "time"
)

type Node struct {
  Endpoint string

  IsHealthy bool
  LastHealthCheck time.Time

  client *memcache.Client
}

func NewNode(endpoint string) *Node {
  return &Node{
    Endpoint: endpoint,
    IsHealthy: false,
    LastHealthCheck: time.Now().Add(-1 * HEALTHCHECK_PERIOD),
    client: memcache.New(endpoint),
  }
} 

func (me *Node) Add(item *memcache.Item, finishChan chan(*NodeResponse)) {
  go func(){
    err := me.client.Add(item)
    if finishChan!=nil { finishChan <- me.getNodeResponse(nil, err) }
  }()
}

func (me *Node) Set(item *memcache.Item, finishChan chan(*NodeResponse)) {
  go func(){
    err := me.client.Set(item)
    if finishChan!=nil { finishChan <- me.getNodeResponse(nil, err) }
  }()
}

func (me *Node) Get(key string, finishChan chan(*NodeResponse)) {
  go func(){
    item, err := me.client.Get(key)
    if finishChan!=nil { finishChan <- me.getNodeResponse(item, err) }
  }()
}

func (me *Node) Delete(key string, finishChan chan(*NodeResponse)) {
  go func(){
    err := me.client.Delete(key)
    if finishChan!=nil { finishChan <- me.getNodeResponse(nil, err) }
  }()
}

func (me *Node) HealthCheck() {
  // Read a Random key, expect ErrCacheMiss
  x := make([]byte,32)
  rand.Read(x)
  _, err := me.client.Get(fmt.Sprintf("%02x", x))
  me.getNodeResponse(nil, err)
}

func (me *Node) getNodeResponse(item *memcache.Item, err error) *NodeResponse {
  me.LastHealthCheck = time.Now()
  if err == nil { me.markHealthy() } else { me.markUnhealthy() }
  return NewNodeResponse(me, item, err)
}

func (me *Node) markHealthy() { me.IsHealthy = true }
func (me *Node) markUnhealthy() { me.IsHealthy = false }