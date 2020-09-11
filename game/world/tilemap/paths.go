package tilemap

import (
	"libs/mathm/ints"
	"libs/mathm/points"
	"libs/pathfinder"
)

// TeamPath is collection of pathgraphs with common costs
type TeamPath struct {
	Graphs []*pathfinder.PathGraph
	Costs  [][]int
}

// Paths stores PathGraphs
type Paths map[uint16]*TeamPath

// AddTeam adds new team to Paths
func (p Paths) AddTeam(m *Tilemap, team uint16) {
	tp := &TeamPath{
		Costs: ints.BoundArr2D(pathfinder.Infinity, m.GetCosts(team)),
	}
	content := make([]*pathfinder.PathGraph, len(TargetTypes))
	for i := range content {
		content[i] = pathfinder.New(m.Point, tp.Costs)
		content[i].Map(m.Get(team, i))
	}
	tp.Graphs = content
	p[team] = tp
}

// GetStep returns next step to go from given point
func (p Paths) GetStep(team uint16, kind int, pos points.Point) points.Point {
	return p[team].Graphs[kind].GetStep(pos)
}

// Update of paths
func (p Paths) Update(m *Tilemap, team uint16, places []points.Point, remove bool) {
	for k, t := range p {
		if k != team {
			if remove {
				for _, p := range places {
					t.Costs[p.R+1][p.C+1] -= m.GetTile(p).GetCost(k)
				}
			} else {
				for _, p := range places {
					t.Costs[p.R+1][p.C+1] += m.GetTile(p).GetCost(k)
				}
			}

			for i, p := range t.Graphs {
				p.Lock()
				p.D.Lock()
				p.Costs = t.Costs
				p.Unlock()
				p.D.Unlock()
				go p.Map(m.Get(k, i))
			}

		}
	}
}

// Targets is datastructure that stores all targets
type Targets map[uint16][][]points.Point

// AddTeam adds new team to targets
func (p Targets) AddTeam(team uint16) {
	p[team] = make([][]points.Point, len(TargetTypes))
}

// Add adds new target to targets
func (p Targets) Add(team uint16, kind int, target points.Point) {
	p[team][kind] = append(p[team][kind], target)

}

// Remove removed target from targets
func (p Targets) Remove(team uint16, kind int, target points.Point) {
	arr := p[team][kind]
	for i, t := range arr {
		if t.Eq(target) {
			p[team][kind] = append(arr[:i], arr[i+1:]...)
			return
		}
	}
}

// Get gets all targets for given team and TargetType
func (p Targets) Get(team uint16, kind int) []points.Point {
	targets := []points.Point{}
	for k, v := range p {
		if k != team {
			targets = append(targets, v[kind]...)
		}
	}
	return targets
}
