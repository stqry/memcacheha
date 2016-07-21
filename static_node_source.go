package memcacheha

// StaticNodeSource represents a static list of nodes
type StaticNodeSource []string

// Return a new StaticNodeSource with the given endpoints
func NewStaticNodeSource(nodes ...string) *StaticNodeSource {
	x := StaticNodeSource(nodes)
	return &x
}

// Implement NodeSource, return a slice of configured endpoints
func (me *StaticNodeSource) GetNodes() ([]string, error) {
	return *me, nil
}
