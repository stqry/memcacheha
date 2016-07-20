package memcacheha

type StaticNodeSource []string

func NewStaticNodeSource(nodes ...string) *StaticNodeSource {
  x := StaticNodeSource(nodes)
  return &x
}

func (me *StaticNodeSource) GetNodes() ([]string, error) {
  return *me, nil
}