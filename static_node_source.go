package memcacheha

// StaticNodeSource represents a static list of nodes
type StaticNodeSource []string

// NewStaticNodeSource returns a new StaticNodeSource with the given endpoints
func NewStaticNodeSource(nodes ...string) *StaticNodeSource {
	staticNodeSource := StaticNodeSource(nodes)
	return &staticNodeSource
}

// GetNodes implements NodeSource, return a slice of configured endpoints
func (staticNodeSource *StaticNodeSource) GetNodes() ([]string, error) {
	return *staticNodeSource, nil
}
