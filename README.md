# memcacheha

[![](https://godoc.org/github.com/apitalent/memcacheha?status.svg)](https://godoc.org/github.com/apitalent/memcacheha)

memcacheha wraps [gomemcache](https://github.com/bradfitz/gomemcache) to provide HA (highly available) functionality with lazy client-side synchronization.

# How is this different from gomemcache multi-node?

[gomemcache](https://github.com/bradfitz/gomemcache) performs client-side sharding (distributes keys across multiple memcache nodes), whereas memcacheha 
is designed to write to all memcache nodes and synchronise nodes with missing data during reads. This is useful 
in situations where memcache availability and consistency is not negligible, i.e. when memcache is being used as a 
session store.

## Operation

MemcacheHA operates as a Client nanoservice, maintaining a pool of connections to all configured or discovered memcache
nodes.

The client checks the health of all configured nodes periodically (every 5 seconds, HEALTHCHECK_PERIOD)

Writes are mirrored to all nodes concurrently, and consistency is achieved by not returning until all writes 
have acknowledged or timed out. Reads are performed from at least n/2 nodes where n is the total number of currently
healthy nodes - if at least one node returns data, items are (transparently) written to nodes with missing data. 

## Autodiscovery

Nodes are discovered through [NodeSource](./node_source.go)s - currently, two are available:

* [StaticNodeSource](./static_node_source.go) - Allows nodes to be configured statically (e.g. from a config file or ENV)
* [ElastiCacheNodeSource](./elasticache_node_source.go) - Retreives nodes from an AWS ElastiCache cluster

Multiple sources can be used, passed to `New` in [Client](./client.go). All sources will be queried once every 10 seconds (GET_NODES_PERIOD).

## Example

```golang
 	// You can use any type that implements the Logger interface.
	logger = log.NewConsoleLogger("debug")

	// Configure an AWS ElasticCache Source
	source := memcacheha.NewElastiCacheNodeSource(logger, "ap-southeast-2", "myMemcacheCluster")  

	// Get a new client
	client := memcacheha.New(logger, source)

	// Start the nanoservice
	client.Start()

	// ...use client as if you were talking to one memcache via gomemcache...

	// Stop the nanoservice
	client.Stop()
```

## Detail

### Failover condition assumptions

* Only one node will be lost at once
* A node (re)joining the cluster will be empty.

### Unconditional Writing

* Items will be concurrently written to all healthy nodes. The write will not return until:
	* All nodes have been written to and responded, or timed out

### Conditional Writing

* Items will be concurrently written to all healthy nodes. The write will not return until:
	* All nodes have been written to and responded, or timed out
* If any node responds with conditional write fail:
	* The value will be re-read from that node and unconditionally written to all healthy nodes
	* The call will return with conditional write fail only after all nodes have responded or timed out on the second write

### Reading

* If no healthy nodes are available, the client will return an error.
* Ceil(n/2) random nodes of _n_ healthy nodes are selected for reads.
* When all nodes return a cache miss, the response is a cache miss.
* If any node returns a hit and any other node(s) return a miss, the value will be written to the missing nodes

### Deleting

* Keys will be concurrently deleted from all healthy nodes.
* **CAVEAT:** If a node drops from the cluster, misses a DELETE, and then rejoins the cluster maintaining its old data, the next GET will synchronise the data to all nodes again. This behaviour can be mitigated by always setting expiry timeouts on keys.

### Health checks

* Health checks occur on all nodes periodically, and also as part of any node operation
* A node health check will pass if:
	* The node responds to a GET for a random string with a cache miss within a timeout (100ms)
* A node health check will fail if:
	* The node fails to respond to any operation within a timeout (100ms)
	* The node responds with a Server Error

## Contributing

Contributions are welcome. Please follow the [code of conduct](./code_of_conduct.md).

## License

Please see the [license](./LICENSE) file.

## Credits

Inspired by and uses [gomemcache](https://github.com/bradfitz/gomemcache) by [Brad Fitzpatrick](https://github.com/bradfitz).