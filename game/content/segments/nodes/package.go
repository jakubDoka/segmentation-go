package nodes

// Package is way of transfering resources.
// It also stores the path to destinations.
type Package struct {
	Path []*Entity
	Storage
	Progress int
}

// Fail lets know the destination that delivery failed
func (p *Package) Fail() {
	for k, v := range p.Stored {
		p.Requested[k] += v
	}
}

// NewPackage returns new Package ready to use
func NewPackage(path []*Entity) *Package {
	return &Package{
		Path: path,
		Storage: Storage{
			Requested: path[len(path)-1].Requested,
			Stored:    map[Item]int{},
		},
	}
}
