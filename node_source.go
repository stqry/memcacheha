package memcacheha

type NodeSource interface {
  GetNodes() ([]string, error)
}