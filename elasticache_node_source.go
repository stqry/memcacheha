package memcacheha

import(
  "github.com/aws/aws-sdk-go/service/elasticache"

  "errors"
  "fmt"
)

const(
  ELASTICACHE_CLUSTER_TYPE_MEMCACHE = "memcache"
)

type ElastiCacheNodeSource struct {
  CacheClusterId string

  client *elasticache.ElastiCache
}

func NewElastiCacheNodeSource(awsRegion string, cacheClusterId string) *ElastiCacheNodeSource {
  inst := &ElastiCacheNodeSource{
    CacheClusterId: cacheClusterId,
    client: elasticache.New(nil),
  }
  return inst
}

func (me *ElastiCacheNodeSource) GetNodes() ([]string, error) {
  // Create input struct
  x := true
  input := &elasticache.DescribeCacheClustersInput{
    CacheClusterId: &me.CacheClusterId,
    ShowCacheNodeInfo: &x,
  }

  // Get the AWS cache cluster
  output, err := me.client.DescribeCacheClusters(input)
  if err !=nil { return nil, ErrElastiCacheNotAvailable }

  // Set up output
  var out []string

  // Check that there is only one cluster, and that it is a memcache cluster
  if len(output.CacheClusters)!=1 { return nil, ErrElastiCacheMultipleClusters }
  cluster := output.CacheClusters[0]
  if *cluster.CacheNodeType != ELASTICACHE_CLUSTER_TYPE_MEMCACHE { return nil, ErrElastiCacheNotMemcache }

  // Iterate nodes, get addresses
  for _, node := range cluster.CacheNodes {
    ep := node.Endpoint
    out = append(out, fmt.Sprintf("%s:%d", ep.Address, ep.Port))
  }

  return out, nil
}

var(
  ErrElastiCacheNotAvailable = errors.New("memcacheha: elasticache: not available")
  ErrElastiCacheMultipleClusters = errors.New("memcacheha: elasticache: DescribeCacheClusters returned more than one cluster")
  ErrElastiCacheNotMemcache = errors.New("memcacheha: elasticache: Not a memcache cluster")
)
