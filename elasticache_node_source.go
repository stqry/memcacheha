package memcacheha

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"

	"errors"
	"fmt"
)

const (
	// ELASTICACHE_ENGINE_MEMCACHE is the AWS Engine type for a memcached cluster
	ELASTICACHE_ENGINE_MEMCACHE = "memcached"
)

// ElastiCacheNodeSource represents a source of nodes from an AWS ElastiCache cluster
type ElastiCacheNodeSource struct {
	AWSRegion      string
	CacheClusterId string
	Log            Logger
}

// NewElastiCacheNodeSource returns a new ElastiCacheNodeSource with the given logger, AWS region, and cache cluster ID
func NewElastiCacheNodeSource(log Logger, awsRegion string, cacheClusterId string) *ElastiCacheNodeSource {
	inst := &ElastiCacheNodeSource{
		AWSRegion:      awsRegion,
		CacheClusterId: cacheClusterId,
		Log:            newScopedLogger("ElastiCache Source", log),
	}
	return inst
}

// GetNodes implements NodeSource, querying the AWS API to get the nodes in the configured CacheClusterId
func (elastiCacheNodeSource *ElastiCacheNodeSource) GetNodes() ([]string, error) {
	// AWS Session / Client
	sess := session.New(&aws.Config{Region: aws.String(elastiCacheNodeSource.AWSRegion)})
	client := elasticache.New(sess)

	// Create input struct
	x := true
	input := &elasticache.DescribeCacheClustersInput{
		CacheClusterId:    &elastiCacheNodeSource.CacheClusterId,
		ShowCacheNodeInfo: &x,
	}

	// Get the AWS cache cluster
	output, err := client.DescribeCacheClusters(input)
	if err != nil {
		return nil, err
	}

	// Set up output
	var out []string

	// Check that there is only one cluster, and that it is a memcache cluster
	if len(output.CacheClusters) > 1 {
		return nil, ErrElastiCacheMultipleClusters
	}
	cluster := output.CacheClusters[0]
	if *cluster.Engine != ELASTICACHE_ENGINE_MEMCACHE {
		return nil, fmt.Errorf("Not a memcache cluster, type %s", *cluster.Engine)
	}

	// Iterate nodes, get addresses
	for _, node := range cluster.CacheNodes {
		if node != nil {
			ep := node.Endpoint
			if ep != nil {
				out = append(out, fmt.Sprintf("%s:%d", *ep.Address, *ep.Port))
			}
		}
	}

	return out, nil
}

var (
	// ErrElastiCacheMultipleClusters is an error meaning that the AWS discovery call returned more than one cluster
	ErrElastiCacheMultipleClusters = errors.New("DescribeCacheClusters returned more than one cluster")

	// ErrElastiCacheNotMemcache is an error meaning that the AWS discovery call returned a cluster that is not a memcached cluster
	ErrElastiCacheNotMemcache = errors.New("Not a memcache cluster")
)
