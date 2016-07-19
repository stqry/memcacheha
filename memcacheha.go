package memcacheha

import(
  "github.com/bradfitz/gomemcache/memcache"

  "time"
)

const VERSION = "0.0.1"

var(
  GET_NODES_PERIOD time.Duration = time.Duration(60 * time.Second)
  HEALTHCHECK_PERIOD time.Duration = time.Duration(10 * time.Second)
)

type MemcacheHA struct {
  Nodes *NodeList
  Sources []NodeSource

  shutdownChan chan(int)
  running bool
}

func NewMemcacheHA(sources ...NodeSource) *MemcacheHA {
  i := &MemcacheHA{
    Nodes: NewNodeList(),
    Sources: sources,
    shutdownChan: make(chan(int)),
    running: false,
  }
  return i
}

// Add writes the given item, if no value already exists for its key. ErrNotStored is returned if that condition is not met.
func (me *MemcacheHA) Add(item *memcache.Item) error {
  // Get all nodes that are marked healthy
  nodes := me.Nodes.GetHealthyNodes()
  nodeCount := len(nodes)

  // Bug out early if no nodes
  if nodeCount == 0 { return ErrNoHealthyNodes }

  finishChan := make(chan(error))
  statusChan := make(chan(*NodeResponse), nodeCount)

  // Concurrently write to all nodes
  for _, node := range nodes { node.Add(item, statusChan) }

  // Handle responses
  go func(){
    defer func(){ r := recover(); if r != nil { finishChan <- ErrUnknown } }()

    for ; nodeCount > 0; nodeCount-- {
      response := <- statusChan
      if response.Error != nil { finishChan <- ErrUnknown; return }
    }

    finishChan <- nil
  }()

  // Wait for final response and return
  return <- finishChan
}

// Set writes the given item, unconditionally.
func (me *MemcacheHA) Set(item *memcache.Item) error {
  // Get all nodes that are marked healthy
  nodes := me.Nodes.GetHealthyNodes()
  nodeCount := len(nodes)

  // Bug out early if no nodes
  if nodeCount == 0 { return ErrNoHealthyNodes }

  finishChan := make(chan(error))
  statusChan := make(chan(*NodeResponse), nodeCount)

  // Concurrently write to all nodes
  for _, node := range nodes { node.Set(item, statusChan) }

  // Handle responses
  go func(){
    defer func(){ r := recover(); if r!=nil { finishChan <- ErrUnknown } }()

    for ; nodeCount > 0; nodeCount-- {
      response := <- statusChan
      if response.Error != nil { finishChan <- ErrUnknown; return }
    }

    finishChan <- nil
  }()

  // Wait for final response and return
  return <- finishChan
}

// Get gets the item for the given key. ErrCacheMiss is returned for a memcache cache miss. The key must be at most 250 bytes in length.
func (me *MemcacheHA) Get(key string) (*memcache.Item, error) {
  // Get all nodes that are marked healthy
  nodes := me.Nodes.GetHealthyNodes()
  nodeCount := len(nodes)

  // Bug out early if no nodes
  if nodeCount == 0 { return nil, ErrNoHealthyNodes }

  // Reduce to less than 3 nodes - thanks to golang, these will be random
  for k, _ := range nodes {
    if len(nodes) < 3 { break }
    delete(nodes, k)
  }

  finishChan := make(chan(*NodeResponse))
  statusChan := make(chan(*NodeResponse), nodeCount)

  // Concurrently read to all nodes
  for _, node := range nodes { node.Get(key, statusChan) }

  var nodesToSync []*Node 

  // Handle responses
  go func(){
    defer func(){ r := recover(); if r!=nil { finishChan <- NewNodeResponse(nil, nil, ErrUnknown) } }()

    var item *memcache.Item

    for ; nodeCount > 0; nodeCount-- {
      response := <- statusChan
      if response.Error == memcache.ErrCacheMiss { nodesToSync = append(nodesToSync, response.Node) }
      if response.Error == nil && response.Item != nil { item = response.Item }
    }

    // Resync by writing to missing nodes
    if item!=nil {
      for _, node := range nodesToSync {
        node.Set(item, nil)
      }
    }

    finishChan <- nil
  }()

  res := <- finishChan
  if res!=nil { return res.Item, res.Error }

  return nil, ErrUnknown
}

// Delete deletes the item with the provided key. The error ErrCacheMiss is returned if the item didn't already exist in the cache.
func (me *MemcacheHA) Delete(key string) error {
  // Get all nodes that are marked healthy
  nodes := me.Nodes.GetHealthyNodes()
  nodeCount := len(nodes)

  // Bug out early if no nodes
  if len(nodes) == 0 { return ErrNoHealthyNodes }

  finishChan := make(chan(error))
  statusChan := make(chan(*NodeResponse), nodeCount)

  // Concurrently delete from all nodes
  for _, node := range nodes { node.Delete(key, statusChan) }

  // Handle responses
  go func(){
    defer func(){ r := recover(); if r!=nil { finishChan <- ErrUnknown } }()

    for ; nodeCount > 0; nodeCount-- {
      response := <- statusChan
      if response.Error != nil { finishChan <- ErrUnknown; return }
    }

    finishChan <- nil
  }()

  return <- finishChan
}

func (me *MemcacheHA) Start() error {
  if me.running != false { return ErrAlreadyRunning } 
  go me.runloop()
  return nil
}

func (me *MemcacheHA) runloop() {
  timerChannel := time.After(time.Duration(time.Second))

  lastGetNodes := time.Now().Add(-1*GET_NODES_PERIOD)

  for {
    select {
    case <- timerChannel:
      now := time.Now()

      if lastGetNodes.Add(GET_NODES_PERIOD).Before(now) {
        me.GetNodes()
        lastGetNodes = time.Now()
      }

      timerChannel = time.After(time.Duration(time.Second))

    case <- me.shutdownChan: 
      me.running = false
      me.shutdownChan <- 2
      return
    }
  }
}

func (me *MemcacheHA) GetNodes() {
  for _, source := range me.Sources {

    nodes, _ := source.GetNodes()

    for _, nodeAddr := range nodes {
      if !me.Nodes.Exists(nodeAddr) { node := NewNode(nodeAddr); me.Nodes.Add(node); node.HealthCheck() }
    }
  }
}

func (me *MemcacheHA) Stop() error {
  if me.running != true { return ErrAlreadyRunning }
  me.shutdownChan <- 1
  <- me.shutdownChan
  return nil
}