package worms

import (
	"libs/mathm/angles"
	"myFirstProject/game"
	"myFirstProject/game/world/tilemap"

	"github.com/faiface/pixel"
)

// NetworkUpdate handles optimized updating of worm in case of server connection
func (e *Entity) NetworkUpdate() {
	e.Front.Interpolate()
	draw(e.Front)
	for _, s := range e.Segments {
		s.Interpolate()
		s.Module.Update(s.Determination.Vec, 0)
		draw(s)
	}
	e.End.Interpolate()
	draw(e.End)
}

// Correct corrects worm
func (e *Entity) Correct() {
	e.Front.Correct()
	for _, s := range e.Segments {
		s.Correct()
	}
	e.End.Correct()
}

// AIUpdate si logical control of worm
func (e *Entity) AIUpdate() {
	pos := e.GetPos()
	var dest pixel.Vec

	// Follow and orbit target if there is any
	if e.Target != nil {
		tPos := e.Target.GetPos()
		vec := tPos.To(pos)
		if vec.Len() > e.Sight*.7 { // .7 is kinda optimal...
			dest = tPos
		} else {
			dest = vec.Rotated(.2).Add(tPos)
		}
		e.Target = nil // target is figured out from turret targets so target will be up to date
	} else { //otherwise follow path
		dest = tilemap.World.GetStep(e.Team, int(e.Objective), pos)
	}

	// worm has some segments disconnected so he has to wait for them to reconnect
	if e.Apart {
		e.Apart = false
		return
	}

	// now change velocity and rotation accordingly
	e.Turn(dest)
	e.Charge()
}

// Update ...
func (e *Entity) Update() {
	if e.Controlled {
		e.SegmentInput()
		e.Input()
		if game.IsNetworking {
			if e.Front.Determination.GetOffset() > 100 {
				e.Correct()
			}
			e.SyncVel()
		}
		e.MovementInput()
	} else if game.IsNetworking {
		e.NetworkUpdate()
		return
	} else if !e.Client {
		e.AIUpdate()
	}

	e.Vel = e.Vel.Sub(e.Vel.Scaled(FRICTION))
	if e.Vel.Len() > e.MaxSpeed {
		e.Vel = e.Vel.Unit().Scaled(e.MaxSpeed)
	}

	draw(e.Front)
	vel := e.Front.Pull(e.Vel.Scaled(game.Delta))

	c := 0
	p := e.Front

	for i, s := range e.Segments {
		draw(s)

		originalRot := s.GetRot()
		// special cases
		if s.Dead {
			if game.IsServer {
				e.Give(s.ID)
				e.PopSeg(s.ID)
				delete(e.Indexed, s.ID)
			}
			if game.Selected == i {
				game.Selected = -1
			}
			if game.TransSelected == i {
				game.TransSelected = -1
			}
			continue
		} else if s.Busy {
			s.MoveToDest()
		} else if !s.Deployed && !s.Transforming {
			vel = e.moveSegment(p, s, vel)
			p = s

			if s.Target != nil {
				e.Target = s.Target
			}
		}

		s.Update(angles.To(originalRot, s.GetRot()))

		e.Segments[c] = s
		c++
	}
	draw(e.End)
	e.moveSegment(p, e.End, vel)

	if c != len(e.Segments) {
		for j := c; j < len(e.Segments); j++ {
			e.Segments[j] = nil
		}

		e.Segments = e.Segments[:c]
		e.ReIndex()
	}

}
