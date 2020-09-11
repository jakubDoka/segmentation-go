package bullets

import (
	"libs/collizions"
	"libs/threads"
	"myFirstProject/game"
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/ent/ai"
	"myFirstProject/game/world/graphics"

	"github.com/faiface/pixel"
)

// Type si struct that can be used to define kinds of bullets
type Type struct {
	pixel.Rect
	Speed, Livetime, Size, Range, TrailFrequency float64
	HitAir, HitGround                            bool
	Damage                                       int
	Behavior                                     func(e *Entity)
	StartExplosion, Trail, EndExplosion          func(pos, dir pixel.Vec)
}

// Init initializes Bullet type
func Init(t Type) Type {
	t.Range = t.Livetime * t.Speed
	return t
}

// GetDamage returns bullet damage
func (t *Type) GetDamage() int {
	return t.Damage
}

// GetSpeed returns speed of bullet
func (t *Type) GetSpeed() float64 {
	return t.Speed
}

// New creates new bullet
func (t *Type) New(start, vel pixel.Vec, rot float64, team uint16) *Entity {
	resVel := vel.Add(pixel.V(t.Speed, 0).Rotated(rot))

	return t.Custom(start, resVel.Angle(), 1, resVel.Len()/t.Speed, team)
}

// Custom creates bullet with custom speed and livetime
func (t *Type) Custom(start pixel.Vec, rot, liveMul, speedMul float64, team uint16) *Entity {
	bullet := &Entity{
		Determination: ent.Determination{
			Transform: ent.Transform{
				APos: ent.APos{Vec: start},
				ARot: ent.ARot{Rot: rot},
			},
			Visibility: ent.Visibility{Sprite: *pixel.NewSprite(nil, t.Rect)},
		},
		ATeam:    ent.ATeam{Team: team},
		SpeedMul: speedMul,
		Type:     t,
	}
	if t.StartExplosion != nil {
		t.StartExplosion(start, pixel.Unit(rot))
	}

	bullet.TrailCycle()
	threads.Queue.PostDeferred(t.Livetime*liveMul, func() { bullet.Dead = true })
	Q.Add(bullet)
	return bullet
}

// Entity is bullet that si controlled by BulletTypes Behavior
type Entity struct {
	*Type
	ent.Determination
	ent.Mortality
	ent.ATeam
	Dead, NotTrail bool
	SpeedMul       float64
}

// GetRect is interface method gor collision detection
func (e *Entity) GetRect() pixel.Rect {
	return pixel.R(e.X-e.Size, e.Y-e.Size, e.X+e.Size, e.Y+e.Size)
}

// GetColliding returns list of objects that collide with the bullet
func (e *Entity) GetColliding() []collizions.Shape {
	return ai.Scanner.GetCollidingEnemy(e.GetRect(), e.Team)
}

// DealDamage deals damage to all coliding shapes
func (e *Entity) DealDamage() {
	for _, c := range e.GetColliding() {
		val, ok := c.(ent.Hitable)
		if ok {
			e.Dead = true
			if game.IsNetworking {
				return
			}
			val.Demage(e.Damage)

		}
	}
}

// Move moves a bullet
func (e *Entity) Move() {
	if !e.NotTrail && e.Trail != nil {
		e.TrailCycle()
		e.Trail(e.Vec, pixel.Unit(e.Rot))
	}
	e.Vec = e.Add(pixel.V(e.Speed*e.SpeedMul, 0).Rotated(e.Rot).Scaled(game.Delta))
}

// TrailCycle executes trail cycle
func (e *Entity) TrailCycle() {
	e.NotTrail = true
	threads.Queue.PostDeferred(e.TrailFrequency, func() { e.NotTrail = false })
}

// Q for worms
var Q Queue

// Queue updates all worms and deletes dead ones
type Queue struct {
	Entities []*Entity
}

// Add adds entity to queue
func (q *Queue) Add(e *Entity) {
	q.Entities = append(q.Entities, e)
}

// Update ...
func (q *Queue) Update() {
	i := 0
	for _, e := range q.Entities {
		e.Behavior(e)

		if e.Dead {
			if e.EndExplosion != nil {
				e.EndExplosion(e.Vec, pixel.Unit(e.Rot))
			}
			continue
		}

		if e.Determination.GetRect().Intersects(graphics.View) {
			e.Draw(Batch)
		}

		q.Entities[i] = e
		i++
	}

	if i == len(q.Entities) {
		return
	}

	for j := i; j < len(q.Entities); j++ {
		q.Entities[j] = nil
	}

	q.Entities = q.Entities[:i]
}
