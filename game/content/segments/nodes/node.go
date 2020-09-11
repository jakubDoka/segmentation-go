package nodes

import (
	"myFirstProject/game/content/bullets"
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/graphics"
	"sync"

	"github.com/faiface/pixel"
)

// Mode is node behavior mode
type Mode int

// Mode enum
const (
	Consumer Mode = iota
	Producer
	Transporter
	All
)

// Type can define kinds of Nodes
type Type struct {
	*DType
	PickupRange float64
}

// Entity stores information about node
type Entity struct {
	*Type
	sync.Mutex
	ent.Competence
	Storage
	pixel.Vec
	Cons, Recs map[*Entity]bool
	packs      []*Package
	Drones     []*DEntity

	Mode Mode
	cost float64
}

// New creates new node ready to use
func (t *Type) New() *Entity {
	new := &Entity{
		Cons:    map[*Entity]bool{},
		Recs:    map[*Entity]bool{},
		Storage: *NewStorage(),
		Type:    t,
		cost:    -1,
	}

	return new
}

func getPathLength(path []*Entity) float64 {
	if len(path) == 0 {
		return 0
	}
	current := path[0]
	var res float64 = 0
	for i := 1; i < len(path); i++ {
		res += path[i].Sub(current.Vec).Len()
	}
	return res
}

func extract(cons map[*Entity]bool) []*Entity {
	c := make([]*Entity, 0, len(cons))
	for k := range cons {
		c = append(c, k)
	}
	return c
}

// Query finds all possible package requesters
func (e *Entity) Query() map[*Entity]*Package {
	e.cost = 0
	var level int
	var pathLength float64
	paths := map[*Entity]*Entity{}
	packs := map[*Entity]*Package{}
	toQuery := extract(e.Cons)
	seen := []*Entity{}
	for len(toQuery) != 0 {
		new := []*Entity{}
		level++
		for _, c := range toQuery {
			for o := range c.Cons {
				nextCost := c.cost + o.To(c.Vec).Len()
				if o.cost == -1 {
					new = append(new, o)
					o.cost = nextCost
					seen = append(seen, o)
					paths[o] = c
				} else if o.cost > nextCost {
					new = append(new, o)
					o.cost = nextCost
					paths[o] = c
				}
			}
			if len(c.Requested) > 0 {
				var path = make([]*Entity, level)
				current := c
				for i := level - 1; i >= 0; i-- {
					path[i] = current
					current = paths[current]
				}
				if val, ok := packs[c]; ok {
					if getPathLength(val.Path) > pathLength {
						packs[c] = NewPackage(path)
					}
				} else {
					packs[c] = NewPackage(path)
				}
			}
		}
		toQuery = new
	}
	e.cost = -1
	for _, s := range seen {
		s.cost = -1
	}

	return packs
}

func getNeeded(packages map[*Entity]*Package) map[Item]bool {
	needed := map[Item]bool{}
	for _, p := range packages {
		for k := range p.Requested {
			needed[k] = true
		}
	}
	return needed
}

func getValues(packages map[*Entity]*Package) []*Package {
	values := make([]*Package, 0, len(packages))
	for _, v := range packages {
		values = append(values, v)
	}
	return values
}

func (e *Entity) resolve(packs []*Package) {
	for _, k := range Items {
		stored, ok := e.Stored[k]
		if !ok {
			continue
		}
		for _, p := range packs {
			val, ok := p.Requested[k]
			if !ok {
				continue
			}
			if val > stored {
				p.Stored[k] = stored
				p.Requested[k] -= stored
				delete(e.Stored, k)
				break
			} else {
				p.Stored[k] = p.Requested[k]
				e.Stored[k] -= p.Requested[k]
				delete(p.Requested, k)
			}
		}
	}

	for _, v := range packs {
		if !v.IsEmpty() {
			e.packs = append(e.packs, v)
		}
	}
}

// Distribute Distributes all Packages
func (e *Entity) Distribute() {
	for _, v := range e.packs {
		e.Drones = append(e.Drones, e.DNew(e.Vec, e.GetTeam(), v, e))
		v.Progress++
	}
	e.packs = []*Package{}
}

// Receive receives the package
func (e *Entity) Receive(p *Package) {
	if p.Progress == len(p.Path) {
		e.Merge(&p.Storage)
		return
	}
	e.packs = append(e.packs, p)
}

// Connect connects to given node
func (e *Entity) Connect(o *Entity) {
	//e.Lock()
	e.Cons[o] = true
	o.Recs[e] = true
	//e.Unlock()
}

// CutConnection cuts connection with given node
func (e *Entity) CutConnection(o *Entity) {
	//e.Lock()
	delete(e.Cons, o)
	delete(o.Recs, e)
	//e.Unlock()
}

// Update ...
func (e *Entity) Update() {
	if len(e.packs) != 0 {
		e.Distribute()
	}
	e.UpdateDrones()
	switch e.Mode {
	case Consumer:
		return
	case Transporter:
		return
	}
	packs := e.Query()
	e.resolve(getValues(packs))
}

// UpdateDrones updates state of drones
func (e *Entity) UpdateDrones() {
	i := 0
	for _, d := range e.Drones {
		d.Update()
		if d.Determination.GetRect().Intersects(graphics.View) {
			d.Draw(bullets.Batch)
		}
		if d.Dead {
			continue
		}

		e.Drones[i] = d
		i++

	}

	for j := i; j < len(e.Drones); j++ {
		e.Drones[j] = nil
	}

	e.Drones = e.Drones[:i]
}

// CutAllCons cuts all connections
func (e *Entity) CutAllCons() {
	for k := range e.Cons {
		e.CutConnection(k)
	}
}

// CutAllRecs deleted all senders
func (e *Entity) CutAllRecs() {
	for o := range e.Recs {
		o.CutConnection(e)
	}
}

// CutAll deletes both receivers and senders
func (e *Entity) CutAll() {
	e.CutAllRecs()
	e.CutAllCons()
}
