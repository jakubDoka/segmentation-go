package turrets

import (
	"libs/mathm/angles"
	"libs/mathm/floats"
	"libs/mathm/vectors"
	"libs/netm"
	"libs/threads"
	"math"
	"myFirstProject/game"
	"myFirstProject/game/content/bullets"
	"myFirstProject/game/content/segments"
	"myFirstProject/game/content/segments/nodes"
	"myFirstProject/game/network"
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/ent/ai"
	"myFirstProject/game/world/ent/ai/targets"

	"github.com/faiface/pixel"
)

const unnoticeableAngle = math.Pi / 20

// Type is for defining kinds of Turrets that then can be produced by the type
type Type struct {
	*bullets.Type
	pixel.Rect
	Speed, Inaccuracy, VelInacuracy, ReloadSpeed, Hull, Offset, Pushback, FixSpeed,
	DeployBoost float64
	MaxAmmo, AmmoPerShot int
	ID                   uint8
	Ammo                 nodes.Item
	Shoot                func(e *Entity, dir float64)
}

// New returns new turret entity ready to use
func (t *Type) New(target *segments.Entity) {
	turret := &Entity{
		Determination: ent.Determination{
			Visibility: ent.Visibility{Sprite: *pixel.NewSprite(nil, t.Rect)},
		},
		Type:    t,
		Entity:  target,
		Storage: &target.Node.Storage,
		Side:    1,
		Boost:   1,
	}

	turret.Stored[t.Ammo] = t.MaxAmmo

	target.Node.Mode = nodes.Consumer
	target.Module = turret
}

// GetSight returns range of turret
func (t *Type) GetSight() float64 {
	return t.Range
}

// GetSpeed returns optimal bullet speed
func (t *Type) GetSpeed() float64 {
	return t.Type.Speed + t.VelInacuracy/2
}

// Entity is Turret it stores all data about its existence
type Entity struct {
	*Type
	*segments.Entity
	*nodes.Storage
	netm.Interpolator
	ent.Determination
	Target      ai.Target
	Unloaded    bool
	Side, Boost float64
}

func (e *Entity) Write(b *netm.Buffer) {
	b.PutUint8(ID)
	b.PutUint8(e.Type.ID)
}

// WriteUpdate pushes state to buffer
func (e *Entity) WriteUpdate(b *netm.Buffer) {
	b.PutFloat64(e.Rot)
}

// ReadUpdate updates state of segment from buffer
func (e *Entity) ReadUpdate(b *netm.Buffer) {
	e.Interpolator.Update(pixel.ZV, b.Float64())
}

// Transform applays boost to turret
func (e *Entity) Transform(deploy bool) {
	if deploy {
		e.Boost = e.DeployBoost
	} else {
		e.Boost = 1
	}
}

// Draw draws turret
func (e *Entity) Draw(t pixel.Target) {
	e.Determination.Draw(t)
	e.Determination.Offset = vectors.MoveTwards(e.Determination.Offset, pixel.ZV, e.FixSpeed*game.Delta)
}

// Turn turn turret to position
func (e *Entity) Turn(pos pixel.Vec) {
	e.Rot = angles.TurnTwards(e.Rot, e.Determination.To(pos).Angle(), e.Speed*e.Boost*game.Delta)
}

// Reload reloads turret
func (e *Entity) Reload() {
	e.Unloaded = true
	threads.Queue.PostDeferred(e.ReloadSpeed/e.Boost, func() { e.Unloaded = false })
}

// Input handles turret input from player
func (e *Entity) Input() {

	e.Turn(e.GetPointer())

	if e.CanShoot() && e.Shooting() {
		e.Shoot(e.Rot)
	}
}

// CanShoot returns whether turret can shoot
func (e *Entity) CanShoot() bool {
	return !e.Unloaded && e.Stored[e.Ammo] >= e.AmmoPerShot
}

// Shoot makes turret shoot in given direction
func (e *Entity) Shoot(dir float64) {
	e.Type.Shoot(e, dir)
	e.Use(e.Ammo, e.AmmoPerShot)
	e.Reload()
}

// FireBullet is local call
func (e *Entity) FireBullet(dir float64) {
	ve := e.VelInacuracy / e.Boost
	in := e.Inaccuracy / e.Boost
	e.FireBulletRemote(
		e.Entity.GetVel().Add(pixel.V(floats.RandRange(0, ve), 0).Rotated(dir)),
		dir+floats.RandRange(-in, in),
	)
}

// FireBulletRemote makes turret fire a bullet in direction
func (e *Entity) FireBulletRemote(vel pixel.Vec, dir float64) {
	e.Side = -e.Side

	e.Type.Type.New(
		e.Determination.Vec.Add(pixel.V(e.Hull, e.Type.Offset*e.Side).Rotated(e.Rot)).Add(e.Determination.Offset),
		vel,
		dir,
		e.GetTeam(),
	)
	e.Determination.Offset = pixel.V(-e.Pushback, 0).Rotated(dir)
	if game.IsServer {
		b := netm.Buffer{}
		b.PutUint8(2)
		b.PutUint16(e.GetID())
		b.PutUint16(e.Entity.ID)
		b.PutFloat64(vel.X)
		b.PutFloat64(vel.Y)
		b.PutFloat64(dir)
		network.Clients.TCPAppendExcept(b, e.GetIP())
	}
}

// Update ...
func (e *Entity) Update(pos pixel.Vec, twerk float64) {
	e.Determination.Vec = pos
	e.Rot += twerk

	if e.Transforming {
		return
	}

	if !e.Deployed && (e.IsControlled() || e.IsClient()) {
		e.Input()
		return
	}

	if game.IsNetworking && !e.IsControlled() {
		_, e.Rot = e.Interpolate(game.Alfa)
		return
	}

	if !targets.IsValidTarget(e.Target, e.Vec, e.Range) {
		e.Target = ai.Scanner.GetTarget(e.Range, e.Vec, e.GetTeam())
		e.Entity.Target = e.Target
		if e.Target == nil {
			return
		}
	}

	dist, ok := ai.Predict(e, e.Target)
	if !ok {
		return
	}

	e.Turn(dist)

	dir := e.To(dist).Angle()
	if angles.Between(e.Rot, dir) < unnoticeableAngle && e.CanShoot() {
		e.Shoot(dir)
	}
}
