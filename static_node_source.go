package memcacheha

// StaticNodeSource represents a static list of nodes
type StaticNodeSource []string

// NewStaticNodeSource returns a new StaticNodeSource with the given endpoints
func NewStaticNodeSource(nodes ...string) *StaticNodeSource {
	x := StaticNodeSource(nodes)
	return &x
}

// GetNodes implements NodeSource, return a slice of configured endpoints
func (me *StaticNodeSource) GetNodes() ([]string, error) {
	return *me, nil
}
