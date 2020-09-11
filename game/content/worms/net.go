package worms

import (
	"libs/netm"
	"myFirstProject/game"
	"myFirstProject/game/content/segments"
	"myFirstProject/game/network"
	"myFirstProject/game/world/graphics"
	"time"

	"github.com/faiface/pixel"
)

// ReadUpdate reads worm update from buffer
func (e *Entity) ReadUpdate(b *netm.Buffer) {
	e.Front.ReadUpdate(b)
	length := b.Uint16()
	for i := uint16(0); i < length; i++ {
		id := b.Uint16()
		seg, ok := e.Indexed[id]
		if ok {
			seg.ReadUpdate(b)
		} else {
			segments.BlancReadUpdate(b)
		}
	}
	e.End.ReadUpdate(b)
}

// WriteInput input sends new inputstate to server
func (e *Entity) WriteInput() {
	b := netm.Buffer{}
	b.PutUint8(1)
	b.PutUint16(e.ID)
	b.PutBool(e.InputState.Shoot)
	b.PutFloat64(e.Pointer.X)
	b.PutFloat64(e.Pointer.Y)
	network.Server.UDP.Append(b)
}

// ReadInput id server callback for updating inputstate
func (q *Queue) ReadInput(b *netm.Buffer) {
	worm, ok := q.Indexed[b.Uint16()]
	if ok {
		worm.InputState.Read(b)
	}

}

// Read reads input from buffer
func (i *InputState) Read(b *netm.Buffer) {
	i.Shoot = b.Bool()
	i.Pointer.X = b.Float64()
	i.Pointer.Y = b.Float64()
}

// WriteUpdate writes state of worm to buffer
func (e *Entity) WriteUpdate(b *netm.Buffer) {
	e.Front.WriteUpdate(b)
	b.PutUint16(uint16(len(e.Segments)))
	for _, s := range e.Segments {
		b.PutUint16(s.ID)
		s.WriteUpdate(b)
	}
	e.End.WriteUpdate(b)
}

// Write marshals the worm
func (e *Entity) Write(b *netm.Buffer) {
	b.PutUint8(e.Type.ID)
	b.PutUint16(e.ID)
	b.PutUint16(e.Team)
	b.PutUint16(uint16(len(e.Segments)))
	for _, s := range e.Segments {
		s.Write(b)
	}
}

// Read unmarshal the worm
func Read(b *netm.Buffer) *Entity {
	ent := Types[b.Uint8()].New(b.Uint16(), b.Uint16())
	length := b.Uint16()
	for i := uint16(0); i < length; i++ {
		seg := segments.Read(b)
		ent.Add(seg)
	}
	if b.Failed {
		panic("f")
	}
	return ent
}

// PopSeg sends package declaring death of segment to clients
func (e *Entity) PopSeg(segID uint16) {
	b := netm.Buffer{}
	b.PutUint8(3)
	b.PutUint16(e.ID)
	b.PutUint16(segID)
	network.Clients.TCP.Append(b)
}

// PopSeg deletes segment from worms.
func (q *Queue) PopSeg(b *netm.Buffer) {
	id := b.Uint16()
	q.Indexed[id].DeleteSegment(b.Uint16())
}

// MoveSeg sends package declaring order change of segments
func (e *Entity) MoveSeg(from, to int) {
	b := netm.Buffer{}
	b.PutUint8(5)
	b.PutUint16(e.ID)
	b.PutUint16(uint16(from))
	b.PutUint16(uint16(to))
	network.Clients.TCP.Append(b)
}

// MoveSegment is client callback that calls move on worm
func (q *Queue) MoveSegment(b *netm.Buffer) {
	q.Indexed[b.Uint16()].Move(int(b.Uint16()), int(b.Uint16()))
}

// SyncVel synchronisms worms velocity
func (e *Entity) SyncVel() {
	b := netm.Buffer{}
	b.PutUint8(2)
	b.PutUint16(e.ID)
	b.PutFloat64(e.Vel.X)
	b.PutFloat64(e.Vel.Y)
	network.Server.UDP.Append(b)
}

// SetVel is remote callback triggered by client that sets velocity of worm
func (q *Queue) SetVel(b *netm.Buffer) {
	w, ok := q.Indexed[b.Uint16()]
	if ok {
		w.Vel = pixel.V(b.Float64(), b.Float64())
	}
}

// ReadUpdate ...
func (q *Queue) ReadUpdate(b *netm.Buffer) {

	for _, b := range b.LossySplit() {
		w, ok := q.Indexed[b.Uint16()]
		if ok {
			w.ReadUpdate(b)
		}
	}
	game.UpdateSpacing = time.Since(game.LastUpdate).Seconds()
	game.LastUpdate = time.Now()
}

// DelWorm declares worms death to clients
func (q *Queue) DelWorm(id uint16) {
	b := netm.Buffer{}
	b.PutUint8(4)
	b.PutUint16(id)
	network.Clients.TCP.Append(b)
}

// Pop deletes worm from collection
func (q *Queue) Pop(b *netm.Buffer) {
	seg := q.Indexed[b.Uint16()]
	delete(q.Indexed, seg.ID)
	q.Entities = append(q.Entities[:seg.Idx], q.Entities[seg.Idx+1:]...)
	q.ReIndex()
}

// Goto is Server callback that places a segment on client request
func (q *Queue) Goto(b *netm.Buffer) {
	// TODO Make safety checks
	q.Indexed[b.Uint16()].Indexed[b.Uint16()].Goto(pixel.V(b.Float64(), b.Float64()))
}

// Deploy is remote callback of segment deploy
func (q *Queue) Deploy(b *netm.Buffer) {
	q.Indexed[b.Uint16()].Indexed[b.Uint16()].Deploy()
}

// Undeploy is remote callback of segment undeploy
func (q *Queue) Undeploy(b *netm.Buffer) {
	q.Indexed[b.Uint16()].Indexed[b.Uint16()].Undeploy()
}

// TurretShoot makes turret shoot remotely
func (q *Queue) TurretShoot(b *netm.Buffer) {
	q.Indexed[b.Uint16()].Indexed[b.Uint16()].Module.FireBulletRemote(
		pixel.V(b.Float64(), b.Float64()),
		b.Float64(),
	)
}

// Write translates Queue to bites
func (q *Queue) Write(b *netm.Buffer) {
	b.PutUint8(1)
	b.PutUint16(uint16(len(q.Entities)))
	for _, w := range q.Entities {
		w.Write(b)
	}
}

// WriteUpdate ...
func (q *Queue) WriteUpdate(b *netm.Buffer) {
	b.PutUint8(1)
	for _, w := range q.Entities {
		wb := netm.Buffer{}
		wb.PutUint16(w.ID)
		w.WriteUpdate(&wb)
		b.Append(wb)
	}
}

func (q *Queue) Read(b *netm.Buffer) {
	q.Indexed = map[uint16]*Entity{}
	q.Entities = []*Entity{}
	length := b.Uint16()
	for i := uint16(0); i < length; i++ {
		Read(b)
	}
}

// AddWorm sends data about new worm to clients
func (q *Queue) AddWorm(e *Entity) {
	b := netm.Buffer{}
	b.PutUint8(8)
	e.Write(&b)
	network.Clients.TCP.Append(b)
}

// AddWormC is client callback for AddWorm
func (q *Queue) AddWormC(b *netm.Buffer) {
	Read(b)
}

// BindToWorm gives player control over the worm
func (q *Queue) BindToWorm(b *netm.Buffer) {
	worm := q.Indexed[b.Uint16()]
	worm.Controlled = true
	graphics.Cam.Existence = worm
	game.PlayerTeam = worm.Team
}
