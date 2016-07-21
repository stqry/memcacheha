# memcacheha

[![](https://godoc.org/github.com/apitalent/memcacheha?status.svg)](https://godoc.org/github.com/apitalent/memcacheha)

memcacheha wraps `github.com/bradfitz/gomemcache/memcache` to provide HA (highly available) functionality with lazy client-side synchronization.

# How is this different from gomemcache multi-server?

gomemcache performs client-side sharding (distributes keys across multiple memcache servers), whereas memcacheha 
is designed to write to all memcache servers and synchronise servers with missing data during reads. This is useful 
in situations where memcache availability and consistency is not negligible, i.e. when memcache is being used as a 
session store.

## Operation

memcache operates as a Client nanoservice, maintaining a pool of connections to all configured or discovered memcache
servers. 

The client checks the health of all configured servers periodically (every 5 seconds, HEALTHCHECK_PERIOD)

Writes are mirrored to all servers concurrently, and consistency is achieved by not returning until all writes 
have acknowledged or timed out. Reads are performed from at least n/2 servers where n is the total number of currently
healthy nodes - if at least one server returns data, items are (transparently) written to servers with missing data. 

## Autodiscovery

Nodes are discovered through [NodeSource](./node_source.go)s - currently, two are available:

* [StaticNodeSource](./static_node_source.go) - Allows nodes to be configured statically (e.g. from a config file or ENV)
* [ElastiCacheNodeSource](./elasticache_node_source) - Retreives nodes from an AWS ElastiCache cluster

Multiple sources can be used, passed to `New` in [Client](./client.go). All sources will be queried once every 10 seconds (GET_NODES_PERIOD).

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
	* The call will return with conditional write fail
	* The value will be re-read from that node and unconditionally written to all healthy nodes
		* The call will not return until all nodes have responded or timed out on the second write

### Reading

* If no healthy nodes are available, the client will return an error.
* Two random healthy nodes are selected for reads.
* When both nodes return a cache miss, the response is a cache miss.
* If one node returns a cache miss and another node returns a hit, the value will be written to the missing node
* If both nodes timeout or error, they are marked unhealthy the read is restarted.

### Deleting

* Keys will be concurrently deleted from 

### Health checks

* Health checks occur on all nodes periodically, and also as part of any node operation
* A node health check will pass if:
	* The node responds to a GET for a random string with a cache miss within a timeout (100ms)
* A node health check will fail if:
	* The node fails to respond to any operation within a timeout (100ms)
	* The node responds with a Server Error

