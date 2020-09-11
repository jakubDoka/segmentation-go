package worms

import (
	"libs/mathm/angles"
	"libs/mathm/uints"
	"myFirstProject/game"
	"myFirstProject/game/content/segments"
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/ent/ai"
	"myFirstProject/game/world/graphics"
	"myFirstProject/game/world/tilemap"

	"github.com/faiface/pixel"
)

// FRICTION si worlds friction
var FRICTION float64 = .02 // TODO move friction to world

// SelectionColor is color of uis selection when interactin with worms segments
var SelectionColor pixel.RGBA = pixel.RGBA{R: 1, G: 1, B: 1, A: .5}

// Type is used to define kinds of worms
type Type struct {
	Ends                          *segments.Type
	Acceleration, Steer, MaxSpeed float64
	MaxSegments                   int
	ID                            uint8
	Objective                     tilemap.TargetType
}

// New creates new worm entity ready to use
func (w *Type) New(id, team uint16) *Entity {
	worm := &Entity{
		Type:    w,
		ATeam:   ent.ATeam{Team: team},
		Front:   w.Ends.New(pixel.ZV),
		End:     w.Ends.New(pixel.ZV),
		IDStock: uints.New(3000),
		Indexed: map[uint16]*segments.Entity{},
	}

	worm.Front.Determination.Front = true
	worm.End.Determination.End = true

	Q.Add(worm, id)
	return worm
}

// Entity stores information about worm
type Entity struct {
	*Type
	Sight         float64
	Front, End    *segments.Entity
	Segments      []*segments.Entity
	Indexed       map[uint16]*segments.Entity
	Target        ai.Target
	Apart, Client bool
	ID            uint16
	Idx           int
	IP            string
	InputState
	uints.IDStock
	ent.ATeam
	ent.Movement
	ent.ARot
	ent.Control
}

// GetIP returns ip addres of player controling the worm
func (e *Entity) GetIP() string {
	return e.IP
}

// GetID returns worms id
func (e *Entity) GetID() uint16 {
	return e.ID
}

// Add adds new segment to worm
func (e *Entity) Add(s *segments.Entity) {
	e.Segments = append(e.Segments, s)

	if game.IsServer {
		s.ID = e.Take()
	}
	e.Indexed[s.ID] = s
	s.Parent = e
	ai.Scanner.Insert(s)
	e.ReIndex()
}

// GetPos returns position of worm
func (e *Entity) GetPos() pixel.Vec {
	if len(e.Segments) == 0 {
		return pixel.ZV
	}
	return e.Front.Determination.Vec
}

// Charge speed up the worm
func (e *Entity) Charge() {
	e.Vel = e.Vel.Add(pixel.V(e.Acceleration*game.Delta, 0).Rotated(e.Rot))
}

// Turn turns worm to position
func (e *Entity) Turn(pos pixel.Vec) {
	e.Rot = angles.TurnTwards(e.Rot, e.GetPos().To(pos).Angle(), e.Steer*game.Delta)
}

func draw(s *segments.Entity) {
	if s.Determination.GetRect().Intersects(graphics.View) {
		s.Draw()
	}
}

// IsClient returns whether worm is client controlled
func (e *Entity) IsClient() bool {
	return e.Client
}

func (e *Entity) moveSegment(p, s *segments.Entity, vel pixel.Vec) pixel.Vec {
	pos := s.Determination.Vec
	vec := s.Vec.To(p.Back)
	if vec.Len() > e.MaxSpeed*game.Delta*2 {
		e.Apart = s.Speed < e.MaxSpeed
		return s.Pull(vec.Unit().Scaled(s.Speed * game.Delta))
	}
	vel = s.Pull(vel)
	s.Vel = pos.To(s.Determination.Vec).Scaled(1 / game.Delta)
	return vel.Add(s.Pull(s.Vec.To(p.Back)))
}

// DeleteSegment removes segment as network callback
func (e *Entity) DeleteSegment(id uint16) {
	seg := e.Indexed[id]
	seg.Kill()
	e.Segments = append(e.Segments[:seg.Idx], e.Segments[seg.Idx+1:]...)
	e.ReIndex()
}

// Move moves segment on index to different t index
func (e *Entity) Move(from, to int) {
	if from > to {
		to++
	}
	if e.Controlled {
		game.Selected = to
	}
	moved := e.Segments[from]
	e.Segments = append(e.Segments[:from], e.Segments[from+1:]...)
	temp := make([]*segments.Entity, to)
	copy(temp, e.Segments[:to])
	e.Segments = append(append(temp, moved), e.Segments[to:]...)
	if game.IsServer {
		e.MoveSeg(from, to)
	}
	e.ReIndex()
}

// LastMoving returns last moving segment
func (e *Entity) LastMoving() int {
	for i, s := range e.Segments {
		if s.Independent() {
			if i == 0 {
				return 0
			}
			return i - 1
		}
	}
	return len(e.Segments) - 1
}

// PushBack puts segment from index to ent position
func (e *Entity) PushBack(idx int) {
	e.Move(idx, e.LastMoving())
}

// Place places the turret
func (e *Entity) Place(idx int) {
	e.Move(idx, len(e.Segments)-1)
}

// ReIndex updates indexes of segments
func (e *Entity) ReIndex() {
	for i, s := range e.Segments {
		s.Idx = uint16(i)
	}
}

// Q for worms
var Q = Queue{Indexed: map[uint16]*Entity{}, IDStock: uints.New(3000)}

// Queue updates all worms and deletes dead ones
type Queue struct {
	Indexed  map[uint16]*Entity
	Entities []*Entity
	uints.IDStock
}

// ReIndex updates indexes of worms
func (q *Queue) ReIndex() {
	for i, w := range q.Entities {
		w.Idx = i
	}
}

// Add adds entity to queue
func (q *Queue) Add(e *Entity, id uint16) {
	e.ID = id
	e.Idx = len(q.Entities)
	q.Indexed[id] = e
	q.Entities = append(q.Entities, e)
}

// Update ...
func (q *Queue) Update() {
	i := 0
	for _, e := range q.Entities {
		if !game.IsNetworking && len(e.Segments) == 0 {
			if game.IsServer {
				q.DelWorm(e.ID)
				q.Give(e.ID)
			}
			q.ReIndex()
			delete(q.Indexed, e.ID)
			continue
		}

		e.Update()

		q.Entities[i] = e
		i++
	}

	if i != len(q.Entities) {
		for j := i; j < len(q.Entities); j++ {
			q.Entities[j] = nil
		}

		q.Entities = q.Entities[:i]
	}
}
