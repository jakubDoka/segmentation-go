package segments

import (
	"fmt"
	"libs/graphics/animations"
	"libs/mathm/angles"
	"libs/mathm/floats"
	"libs/mathm/vectors"
	"libs/netm"
	"math"
	"myFirstProject/game"
	"myFirstProject/game/content/segments/nodes"
	"myFirstProject/game/network"
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/ent/ai"
	"myFirstProject/game/world/tilemap"

	"github.com/faiface/pixel"
)

type ModuleType interface {
	New(*Entity)
}

type Mod interface {
	Get(uint8) ModuleType
}

var Modules []Mod

// Module is what segment is carrying
type Module interface {
	ent.Existence
	ent.Averness
	ent.Drawable
	ent.Strength
	Transform(bool)
	Update(pixel.Vec, float64)
	Write(b *netm.Buffer)
	WriteUpdate(b *netm.Buffer)
	ReadUpdate(b *netm.Buffer)
	FireBulletRemote(vel pixel.Vec, dir float64)
}

// PlaceFilter indicates whether segment can or cannot be places of concrete tile
type PlaceFilter interface {
	CanPlace(*tilemap.Tile) bool
}

// Type is for defining kinds of segments
type Type struct {
	*nodes.Type
	Base, Size, Speed, DeployTime     float64
	MaxHealth                         int
	ID                                uint8
	CocPiece, BasePiece, SegmentPiece pixel.Rect
}

// New returns new segment ready to use
func (s *Type) New(pos pixel.Vec) *Entity {

	segment := &Entity{
		Determination: Determination{
			Determination: ent.Determination{
				Visibility: ent.Visibility{Sprite: *pixel.NewSprite(nil, s.SegmentPiece)},
			},

			Stand: Stand{
				Sprite: *pixel.NewSprite(nil, s.BasePiece),
			},
			Con: Con{
				Sprite:   *pixel.NewSprite(nil, s.CocPiece),
				Progress: 1,
				Offset:   pixel.V(-s.Base/2, 0),
			},
		},
		Health: s.MaxHealth,
		Type:   s,
		Vec:    pixel.V(s.Base/2, 0),
		Back:   pixel.V(-s.Base/2, 0),
		Node:   s.Type.New(),
	}
	segment.Node.Competence = segment
	segment.Move(pos)
	return segment
}

// Stand is a turret stand tha apears when turret is deployed
type Stand struct {
	pixel.Sprite
	Progress float64
}

// Draw draws the stand
func (s *Stand) Draw(t pixel.Target, pos pixel.Vec) {
	if s.Progress != 0 {
		s.Sprite.Draw(t, pixel.IM.Scaled(pixel.ZV, s.Progress).Moved(pos))
	}
}

// Con is the connection point that appears when segment is moving
type Con struct {
	pixel.Sprite
	Progress float64
	Offset   pixel.Vec
}

// Draw draws the stand
func (c *Con) Draw(t pixel.Target, pos pixel.Vec, rot float64) {
	if c.Progress != 0 {
		c.Sprite.Draw(t, pixel.IM.Rotated(pixel.ZV, rot).Moved(pos.Add(c.Offset.Scaled(c.Progress).Rotated(rot))))
	}
}

// Transform is special interface of segment because it cannot use net Transform
type Transform interface {
	ent.Existence
	ent.Rot
}

// Determination is special struct for segment because it cannot use ent.Determination
type Determination struct {
	ent.Determination
	ent.Mortality
	netm.Interpolator
	Con
	Stand
	End, Front bool
}

// Interpolate is called only on client side and works only here
// it produces smooth movement eventhough update is done only 10 times a second
func (d *Determination) Interpolate() {
	d.Vec, d.Rot = d.Interpolator.Interpolate(game.Alfa)
}

// GetOffset returns distance betwen segment on server and segment on client side
func (d *Determination) GetOffset() float64 {
	return d.Vec.To(d.Tar).Len()
}

// Update updates the interpolator state
func (d *Determination) Update(pos pixel.Vec, rot float64) {
	d.Interpolator.Update(pos, rot)
}

// Draw draws segment
func (d *Determination) Draw(t pixel.Target) {
	if !d.End {
		d.Con.Draw(t, d.Vec, d.Rot)
	}

	if d.End {
		d.Rot += math.Pi
	}

	d.Stand.Draw(t, d.Vec)
	d.Visibility.Sprite.Draw(t, pixel.IM.Rotated(pixel.ZV, d.Rot).Moved(d.Vec))
}

// GetRect returns visibility rectangle
func (d *Determination) GetRect() pixel.Rect {
	return floats.Cube(d.Vec, d.GetSpriteSize()/2)
}

// Parent is there just to make it more clear what is worm pointer
type Parent interface {
	ent.Competence
	ent.Origin
	ent.SelfControl
	GetPointer() pixel.Vec
	Shooting() bool
	IsClient() bool
	GetIP() string
}

// Entity holds all information about the segment
type Entity struct {
	*Type
	Determination
	PlaceFilter
	ent.Volume
	ent.Movement
	Parent
	pixel.Vec
	Target                       ai.Target
	Health                       int
	Module                       Module
	Node                         *nodes.Entity
	Back                         pixel.Vec
	Deployed, Busy, Transforming bool
	Dest                         pixel.Vec
	Twerk                        float64
	Idx, ID                      uint16
}

func (e *Entity) Write(b *netm.Buffer) {
	b.PutUint8(e.Type.ID)
	b.PutUint16(e.ID)
	b.PutBool(e.Deployed)
	if e.Module != nil {
		e.Module.Write(b)
	}
}

// WriteUpdate pushes state to buffer
func (e *Entity) WriteUpdate(b *netm.Buffer) {
	b.PutFloat64(e.Determination.Vec.X)
	b.PutFloat64(e.Determination.Vec.Y)
	b.PutFloat64(e.GetRot())
	if e.Module != nil {
		e.Module.WriteUpdate(b)
	}
}

// Read returns new entity from buffer
func Read(b *netm.Buffer) *Entity {
	ent := Types[b.Uint8()].New(pixel.ZV)
	ent.ID = b.Uint16()
	if b.Bool() {
		ent.DeployState()
	}
	Modules[b.Uint8()].Get(b.Uint8()).New(ent)
	return ent
}

// ReadUpdate updates state of segment from buffer
func (e *Entity) ReadUpdate(b *netm.Buffer) {
	e.Determination.Update(pixel.V(b.Float64(), b.Float64()), b.Float64())
	if e.Module != nil {
		e.Module.ReadUpdate(b)
	}
}

// BlancReadUpdate When segment is missing witch may happend because update is driven by udp blank
// read is used to skip sequence
func BlancReadUpdate(b *netm.Buffer) {
	b.Float64()
	b.Float64()
	b.Float64()
	b.Float64()
}

// Correct corrects segments position
func (e *Entity) Correct() {
	e.SetPos(e.Tar)
	e.SetRot(e.TarR)
}

// GetRect returns hitbox of segment
func (e *Entity) GetRect() pixel.Rect {
	return floats.Cube(e.Determination.Vec, e.Size)
}

// Kill removes the segment from the game
func (e *Entity) Kill() {
	if e.Deployed {
		tilemap.World.OnTileChange(e.Module, true, true, e.GetTeam(), 0)
	}
}

// Demage deals demage to segment and kills it if needed
func (e *Entity) Demage(amount int) {
	e.Health -= amount
	if e.Health < 1 {
		e.Dead = true
	}
}

// GetPosition returns center of segment
func (e *Entity) GetPosition() pixel.Vec {
	return e.GetDirection().Scaled(0.5).Add(e.Back)
}

// GetDirection returns vector from back to front of the segment
func (e *Entity) GetDirection() pixel.Vec {
	return e.Sub(e.Back)
}

// GetRot returns rotation of segment
func (e *Entity) GetRot() float64 {
	return e.GetDirection().Angle()
}

// Independent returns wheather segment should not be moved by worm
func (e *Entity) Independent() bool {
	return e.Deployed || e.Busy
}

// Pull pulls the segment by its front and return how match back moved
func (e *Entity) Pull(vel pixel.Vec) pixel.Vec {
	dir := e.GetDirection()
	rot := dir.Angle()
	vAng := vel.Angle()
	extra := pixel.ZV

	if vel.Len() > e.Base {
		step := pixel.V(e.Base, 0).Rotated(vAng)
		e.Back = e.Vec
		e.Vec = e.Vec.Add(step)
		vel = vel.Sub(step)
		extra = dir
		dir = step
		rot = vAng
	}

	projected := pixel.V(vel.Len(), 0).Rotated(angles.Between(rot, vAng))
	backMoveLen := e.Base - math.Sqrt(e.Base*e.Base-projected.Y*projected.Y) + projected.X
	backMoveVel := pixel.V(backMoveLen, 0).Rotated(rot)

	e.Vec = e.Add(vel)
	e.Back = e.Back.Add(backMoveVel)

	e.Determination.Vec = e.GetPosition()
	e.Determination.Rot = e.GetRot()

	return backMoveVel.Add(extra)
}

// DeployState sets state of segment to deployed
func (e *Entity) DeployState() {
	e.Stand.Progress = 1
	e.Con.Progress = 0
	e.Deployed = true
	e.Transforming = false
	if game.IsServer {
		tilemap.World.OnTileChange(e.Module, true, false, e.GetTeam(), 0)
	}
}

// Deploy triggers deploy animation of segment
func (e *Entity) Deploy() {
	if game.IsServer {
		b := netm.Buffer{}
		b.PutUint8(6)
		e.putPath(&b)
		network.Clients.TCP.Append(b)
	}
	e.Module.Transform(true)
	e.Transforming = true
	e.Busy = false
	animations.Queue.Add(func(delta float64) bool {
		dif := delta / e.DeployTime
		e.Stand.Progress += dif
		e.Con.Progress -= dif
		if e.Con.Progress < 0 {
			e.DeployState()
			return true
		}
		return false
	})
}

// Undeploy triggers undeploy animation of segment
func (e *Entity) Undeploy() {
	if game.IsServer {
		b := netm.Buffer{}
		b.PutUint8(7)
		e.putPath(&b)
		network.Clients.TCPAppendExcept(b, e.GetIP())
	} else if game.IsNetworking && e.IsControlled() {
		b := netm.Buffer{}
		b.PutUint8(2)
		e.putPath(&b)
		network.Server.TCP.Append(b)
	}
	e.Module.Transform(false)
	e.Transforming = true
	animations.Queue.Add(func(delta float64) bool {
		dif := delta / e.DeployTime
		e.Stand.Progress -= dif
		e.Con.Progress += dif
		if e.Stand.Progress < 0 {
			e.Stand.Progress = 0
			e.Con.Progress = 1
			e.Transforming = false
			e.Deployed = false
			tilemap.World.OnTileChange(e.Module, true, true, e.GetTeam(), 0)
			return true
		}
		return false
	})
}

// Goto sends segment to given position
func (e *Entity) Goto(pos pixel.Vec) {
	if game.IsNetworking && e.IsControlled() {
		b := netm.Buffer{}
		b.PutUint8(1)
		e.putPath(&b)
		b.PutFloat64(pos.X)
		b.PutFloat64(pos.Y)
		network.Server.TCP.Append(b)
	}
	if e.Deployed {
		e.Pickup()
	}
	e.Dest = pos
	e.Busy = true
}

// putPath puts if info to find segment on other side
func (e *Entity) putPath(b *netm.Buffer) {
	b.PutUint16(e.GetID())
	b.PutUint16(e.ID)
}

// Interupt stops segment from going to destination
func (e *Entity) Interupt() {
	e.Busy = false
}

// Pickup conects segment back to worm
func (e *Entity) Pickup() {
	e.Node.CutAll()
	e.Undeploy()
}

// MoveToDest moves segment to destination
func (e *Entity) MoveToDest() {
	step := e.Speed * game.Delta
	if e.Vec.To(e.Dest).Len() < e.Base {
		e.Move(vectors.VelTwards(e.Determination.Vec, e.Dest, step))
		e.Rotate(angles.VelTwards(e.Determination.Rot, 0, step/100))

		if e.Determination.Vec == e.Dest && e.GetRot() == 0 {
			e.Vel = pixel.ZV
			e.Determination.Vec = e.GetPosition()
			e.Node.Vec = e.Determination.Vec
			e.Deploy()
		}
		return
	}

	e.Pull(e.Determination.Vec.To(e.Dest).Unit().Scaled(step))
}

//Move moves base by vector
func (e *Entity) Move(vel pixel.Vec) {
	e.Vec = e.Add(vel)
	e.Back = e.Back.Add(vel)
	e.Determination.Vec = e.Determination.Vec.Add(vel)
}

//RotateAround rotates base around given point
func (e *Entity) RotateAround(around pixel.Vec, angle float64) {
	e.Vec = e.Sub(around).Rotated(angle).Add(around)
	e.Back = e.Back.Sub(around).Rotated(angle).Add(around)
	e.Determination.Rot = e.GetRot()
}

//Rotate rotates segment around center
func (e *Entity) Rotate(angle float64) {
	e.RotateAround(e.Determination.Vec, angle)
}

// SetPos is used only by network callbacks, it does not keep rotation
func (e *Entity) SetPos(pos pixel.Vec) {
	e.Vec = pixel.V(e.Base/2, 0)
	e.Back = pixel.V(-e.Base/2, 0)
	e.Determination.Vec = pixel.ZV
	e.Move(pos)

}

// SetRot sets rotation of the segment
func (e *Entity) SetRot(angle float64) {
	e.Vec = e.Determination.Vec.Add(pixel.V(e.Base/2, 0).Rotated(angle))
	e.Back = e.Determination.Vec.Add(pixel.V(-e.Base/2, 0).Rotated(angle))
	e.Determination.Rot = angle
}

// ConAll connects to all nodes in range
func (e *Entity) ConAll() {
	for _, n := range ai.Scanner.GetBudsInRange(e.Node.Sight, e.Vec, e.GetTeam()) {
		val, ok := n.(*Entity)
		if ok && val.Deployed && !val.Node.Recs[val.Node] {
			e.Node.Connect(val.Node)
		}
	}
}

// Update updates all modules of segment
func (e *Entity) Update(twerk float64) {

	e.Module.Update(e.Determination.Vec, twerk)
	if e.Deployed {
		e.Node.Update()
	}
}

// Draw draw ewerithing that needs to be drawn
func (e *Entity) Draw() {
	if e.Deployed {
		e.Determination.Draw(Turrets)
		if e.End || e.Front {
			return
		}
		e.Module.Draw(Turrets)
	} else {
		e.Determination.Draw(Worms)
		if e.End || e.Front {
			return
		}
		e.Module.Draw(Worms)
	}

}

/* ===Debug=== */

// DebugPos ...
func (e *Entity) DebugPos() pixel.Vec {
	return e.Determination.Vec.Add(pixel.V(0, 50))
}

// DebugString ...
func (e *Entity) DebugString() string {
	return fmt.Sprint(e.Node.Requested, e.Node.Stored)
}
