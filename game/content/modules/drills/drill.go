package drills

import (
	"libs/threads"
	"math"
	"myFirstProject/game"
	"myFirstProject/game/content/segments"
	"myFirstProject/game/content/segments/nodes"
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/tilemap"

	"github.com/faiface/pixel"
)

// Type is to exprex generam drill stats
type Type struct {
	MiningSpeed float64
	CanMine     map[nodes.Item]bool
	Capacity    int
	pixel.Rect
}

// GetDamage is just to satisfy interface
func (t *Type) GetDamage() int {
	return 0
}

// GetSight is just to satisfy interface
func (t *Type) GetSight() float64 {
	return 0
}

// Set is just to satisfy interface
func (t *Type) Set(bool) {}

// New creates new drill and attatches it to segment
/*func (t *Type) New(target *segments.Entity) {
	drill := Entity{
		Determination: ent.Determination{
			Visibility: ent.Visibility{Sprite: *pixel.NewSprite(segments.Sheet, t.Rect)},
		},
		Type:    t,
		Entity:  target,
		Storage: &target.Node.Storage,
	}

	target.Node.Mode = nodes.Producer
	target.Module = &drill
}*/

// Entity stores unique information about concrete drill
type Entity struct {
	*Type
	*segments.Entity
	*nodes.Storage
	ent.Determination
	Mining bool
}

// Draw draws drill
func (e *Entity) Draw(t pixel.Target) {
	e.Determination.Draw(t)
}

// GetMinedItem returns item drill is currently mining
func (e *Entity) GetMinedItem() nodes.Item {
	return tilemap.World.GetTileByPos(e.Vec).Drops
}

// CanMine returns whether drill can mine
func (e *Entity) CanMine() bool {
	item := e.GetMinedItem()
	return !e.Mining && e.Deployed && e.Stored[item] < e.Capacity
}

// Mine executes mining cycle
func (e *Entity) Mine() {
	e.Mining = true
	threads.Queue.PostDeferred(e.MiningSpeed, func() {
		e.Mining = false
		e.Stored[e.GetMinedItem()]++
	})
}

// Update ...
func (e *Entity) Update(pos pixel.Vec, twerk float64) {
	e.Determination.Vec = pos
	e.Rot += twerk

	if e.CanMine() {
		e.Mine()
	}

	if e.Mining {
		e.Rot += math.Phi * 3 * game.Delta
	}
}
