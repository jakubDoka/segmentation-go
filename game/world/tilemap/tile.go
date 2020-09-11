package tilemap

import (
	"fmt"
	"libs/pathfinder"
	"myFirstProject/game/content/segments/nodes"
	"myFirstProject/game/world/ent"

	"github.com/faiface/pixel"
)

// layer enum
const (
	FLOOR = iota
	OVERLAY
	BLOCK
)

// Tile is component of tle map. Is used to render and get
// cost for pathfinding
type Tile struct {
	Layers []*Block
	pixel.Rect
	Drops    nodes.Item
	InRange  map[ent.Strength]uint16
	Occupied bool
	X, Y     int
}

// GetCost returns travel cost of tile
func (t *Tile) GetCost(team uint16) int {
	cost := 0
	if t.Layers[BLOCK] != nil {
		return pathfinder.Infinity
	}
	for k, v := range t.InRange {
		if v != team {
			cost += k.GetDamage()
		}
	}
	return cost + 1 //+ t.Layers[FLOOR].Cost
}

// Draw draws a tile to target
func (t *Tile) Draw(target pixel.Target) {
	pos := t.Center()
	for _, l := range t.Layers {
		if l == nil {
			continue
		}
		l.Draw(target, &pos)
	}
}

/* ==Debug== */
func (t *Tile) DebugPos() pixel.Vec {
	return t.Center()
}

func (t *Tile) DebugString() string {
	g := World.Paths[0].Graphs[0]
	g.Lock()
	g.D.Lock()
	defer g.Unlock()
	defer g.D.Unlock()
	return fmt.Sprintln(g.Mapped[t.Y][t.X])
}
