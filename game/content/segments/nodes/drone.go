package nodes

import (
	"myFirstProject/game"
	"myFirstProject/game/world/ent"

	"github.com/faiface/pixel"
)

// DType can define kinds of drones
type DType struct {
	Cap, MaxHealth     int
	Speed, Size, Sight float64
	pixel.Rect
}

// DNew returns new drone ready to use
func (t *DType) DNew(pos pixel.Vec, team uint16, pack *Package, source *Entity) *DEntity {
	drone := &DEntity{

		Determination: ent.Determination{
			Transform: ent.Transform{
				APos: ent.APos{Vec: pos},
			},
			Visibility: ent.Visibility{Sprite: *pixel.NewSprite(nil, t.Rect)},
		},

		DType:  t,
		ATeam:  ent.ATeam{Team: team},
		P:      pack,
		Dest:   pack.Path[pack.Progress],
		Source: source,
	}

	return drone
}

// DEntity stores information about drone
type DEntity struct {
	*DType
	ent.Determination
	Dead, Done bool
	ent.ATeam
	P            *Package
	Dest, Source *Entity
}

// Update ...
func (e *DEntity) Update() {
	vec := e.To(e.Dest.Vec)
	if vec.Len() < e.Dest.PickupRange {
		if e.Done {
			e.Dead = true
			return
		}
		e.Dest.Receive(e.P)
		e.Done = true
		e.Dest = e.Source
		return
	}

	vel := vec.Unit().Scaled(e.Speed * game.Delta)
	e.Vec = e.Vec.Add(vel)
}
