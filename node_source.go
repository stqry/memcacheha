package memcacheha

// NodeSource is an interface defining the GetNodes function. All node sources must implement NodeSource.
type NodeSource interface {
  GetNodes() ([]string, error)
}