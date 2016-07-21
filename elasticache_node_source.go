package memcacheha

import(
  "github.com/apitalent/memcacheha/log"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/elasticache"

  "errors"
  "fmt"
)

const(
  ELASTICACHE_CLUSTER_TYPE_MEMCACHE = "memcache"
)

// ElastiCacheNodeSource represents a source of nodes from an AWS ElastiCache cluster
type ElastiCacheNodeSource struct {
  AWSRegion string
  CacheClusterId string
  Log log.Logger
}

// Return a new ElastiCacheNodeSource with the given logger, AWS region, and cache cluster ID
func NewElastiCacheNodeSource(logger log.Logger, awsRegion string, cacheClusterId string) *ElastiCacheNodeSource {
  inst := &ElastiCacheNodeSource{
    AWSRegion: awsRegion,
    CacheClusterId: cacheClusterId,
    Log: log.NewScopedLogger("ElastiCache Source",logger),
  }
  return inst
}

// Implement NodeSource, query the AWS API and get the nodes associated with the configured CacheClusterId
func (me *ElastiCacheNodeSource) GetNodes() ([]string, error) {
  // AWS Session / Client
  sess := session.New(&aws.Config{Region: aws.String(me.AWSRegion)})
  client := elasticache.New(sess)

  // Create input struct
  x := true
  input := &elasticache.DescribeCacheClustersInput{
    CacheClusterId: &me.CacheClusterId,
    ShowCacheNodeInfo: &x,
  }

  // Get the AWS cache cluster
  output, err := client.DescribeCacheClusters(input)
  if err !=nil { return nil, err }

  // Set up output
  var out []string

  // Check that there is only one cluster, and that it is a memcache cluster
  if len(output.CacheClusters)!=1 { return nil, ErrElastiCacheMultipleClusters }
  cluster := output.CacheClusters[0]
  if *cluster.CacheNodeType != ELASTICACHE_CLUSTER_TYPE_MEMCACHE { return nil, fmt.Errorf("Not a memcache cluster, type %s", *cluster.CacheNodeType) }

  // Iterate nodes, get addresses
  for _, node := range cluster.CacheNodes {
    ep := node.Endpoint
    out = append(out, fmt.Sprintf("%s:%d", ep.Address, ep.Port))
  }

  return out, nil
}

var(
  ErrElastiCacheMultipleClusters = errors.New("DescribeCacheClusters returned more than one cluster")
  ErrElastiCacheNotMemcache = errors.New("Not a memcache cluster")
)
