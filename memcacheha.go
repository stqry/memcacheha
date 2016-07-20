// memcacheha wraps github.com/bradfitz/gomemcache/memcache to provide HA (highly available) functionality with lazy client-side synchronization.
package memcacheha

import(
  "github.com/apitalent/memcacheha/log"
  "github.com/bradfitz/gomemcache/memcache"
  "time"
)

const VERSION = "0.1.0"

var(
  GET_NODES_PERIOD time.Duration = time.Duration(10 * time.Second)
  HEALTHCHECK_PERIOD time.Duration = time.Duration(5 * time.Second)
)

// The MemcacheHA type represents the cluster client.
type MemcacheHA struct {
  Nodes *NodeList
  Sources []NodeSource
  Log log.Logger

  shutdownChan chan(int)
  running bool
}

// Return a new MemcacheHA with the specified logger and NodeSources
func NewMemcacheHA(logger log.Logger, sources ...NodeSource) *MemcacheHA {
  i := &MemcacheHA{
    Nodes: NewNodeList(),
    Sources: sources,
    Log: logger,
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

  // Concurrently write to all healthy nodes
  for _, node := range nodes { node.Add(item, statusChan) }

  // True if any node returns ErrNotStored
  doSync := false
  // These are the nodes that don't contain the value
  var nodesToSync []*Node 

  // Handle responses
  go func(){
    defer func(){ r := recover(); if r!=nil { finishChan <- ErrUnknown } }()

    // Get response from all nodes
    for ; nodeCount > 0; nodeCount-- {
      response := <- statusChan
      if response.Error == memcache.ErrNotStored { doSync = true }
      if response.Error == nil { nodesToSync = append(nodesToSync, response.Node) }
      // We ignore other errors
    }

    // Where there any ErrNotStored?
    if doSync {
      if len(nodesToSync)>0 {
        me.Log.Info("Add: Synchronising %d nodes", len(nodesToSync))
        // Re-read the original
        item, err := me.Get(item.Key)
        if err!=nil {
          // Write to all sync nodes unconditionally
          for _, node := range nodesToSync { node.Set(item, nil) }
        } 
      }

      finishChan <- memcache.ErrNotStored
      return
    } 

    // If this happened, writes to all nodes failed
    if me.Nodes.GetHealthyNodeCount() == 0 {
      finishChan <- ErrNoHealthyNodes
      return
    }

    // All good
    finishChan <- nil
  }()

  // Return result
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
    // Panic handler
    defer func(){ r := recover(); if r!=nil { finishChan <- ErrUnknown } }()

    for ; nodeCount > 0; nodeCount-- {
      // We actually don't care about errors, Node handles them.
      <- statusChan
    }

    // If this happened, writes to all nodes failed
    if me.Nodes.GetHealthyNodeCount() == 0 {
      finishChan <- ErrNoHealthyNodes
      return
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

  // Reduce to less than 3 nodes - thanks to golang, these 2 nodes will be random
  for k, _ := range nodes {
    if len(nodes) < 3 { break }
    delete(nodes, k)
  }

  finishChan := make(chan(*NodeResponse))
  statusChan := make(chan(*NodeResponse), nodeCount)

  // Concurrently read from nodes
  for _, node := range nodes { node.Get(key, statusChan) }

  // These are the nodes to sync to if we get some ErrCacheMiss from requests
  var nodesToSync []*Node 

  // Handle responses
  go func(){
    // Panic handler
    defer func(){ r := recover(); if r!=nil { finishChan <- NewNodeResponse(nil, nil, ErrUnknown) } }()

    // Placeholder for result
    var item *memcache.Item

    // Get response from all nodes
    for ; nodeCount > 0; nodeCount-- {
      response := <- statusChan
      if response.Error == memcache.ErrCacheMiss { nodesToSync = append(nodesToSync, response.Node) }
      if response.Error == nil && response.Item != nil { item = response.Item }
    }

    // Did we find an item from any node?
    if item!=nil {
      if len(nodesToSync)>0 {
        me.Log.Info("Get: Synchronising %d nodes", len(nodesToSync))
        // Resync by writing to missing nodes 
        for _, node := range nodesToSync { node.Set(item, nil) }
      }

      // Return Item
      finishChan <- NewNodeResponse(nil, item, nil)
      return
    }

    // Something nasty happened
    finishChan <- nil
  }()

  // Wait for aggregate response
  res := <- finishChan
  // return result
  if res!=nil { return res.Item, res.Error }

  // End case for something nasty
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

  // If any node returns ErrCacheMiss return this instead.
  var errToReturn error = nil

  // Handle responses
  go func(){
    // Panic handler
    defer func(){ r := recover(); if r!=nil { finishChan <- ErrUnknown } }()

    for ; nodeCount > 0; nodeCount-- {
      response := <- statusChan
      if response.Error == memcache.ErrCacheMiss { errToReturn = memcache.ErrCacheMiss }
    }

    // If this happened, writes to all nodes failed
    if me.Nodes.GetHealthyNodeCount() == 0 {
      finishChan <- ErrNoHealthyNodes
      return
    }

    finishChan <- errToReturn
  }()

  return <- finishChan
}

// Start the MemcacheHA client. This should be called before any operations are called.
func (me *MemcacheHA) Start() error {
  if me.running != false { return ErrAlreadyRunning } 
  go me.runloop()
  return nil
}

func (me *MemcacheHA) runloop() {
  me.Log.Info("Running")
  timerChannel := time.After(time.Duration(time.Second))
  lastGetNodes := time.Time{}
  lastHealthCheck := time.Time{}
  me.running = true

  for {
    select {
    case <- timerChannel:
      now := time.Now()

      if lastGetNodes.Add(GET_NODES_PERIOD).Before(now) {
        me.GetNodes()
        lastGetNodes = time.Now()
      }

      if lastHealthCheck.Add(HEALTHCHECK_PERIOD).Before(now) {
        me.HealthCheck()
        lastHealthCheck = time.Now()
      }

      timerChannel = time.After(time.Duration(time.Second/10))

    case <- me.shutdownChan: 
      me.running = false
      me.Log.Info("Stopped")
      me.shutdownChan <- 2
      return
    }
  }

}

// Update the list of nodes in the client from the configured sources.
func (me *MemcacheHA) GetNodes() {
  for _, source := range me.Sources {
    nodes, err := source.GetNodes()
    if err != nil { me.Log.Error("GetNodes: Source Error: %s", err); return }

    for _, nodeAddr := range nodes {
      if !me.Nodes.Exists(nodeAddr) { 
        me.Log.Info("GetNodes: New Node %s", nodeAddr)
        node := NewNode(me.Log, nodeAddr); me.Nodes.Add(node); node.HealthCheck()
      }
    }
  }
}

// Perform a healthcheck on all nodes.
func (me *MemcacheHA) HealthCheck() {
  for _, node := range me.Nodes.Nodes {
    node.HealthCheck()
  }
}

// Stop the MemcacheHA client.
func (me *MemcacheHA) Stop() error {
  if me.running != true { return ErrAlreadyRunning }
  me.shutdownChan <- 1
  <- me.shutdownChan
  return nil
}