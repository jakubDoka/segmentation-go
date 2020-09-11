package ent

import (
	"github.com/faiface/pixel"
)

type Visibility struct {
	pixel.Sprite
	Layer int
}

func (v *Visibility) GetSpriteSize() float64 {
	rect := v.Frame()
	return rect.Min.To(rect.Max).Len()
}

type Control struct {
	Controlled bool
}

func (c *Control) Set(state bool) {
	c.Controlled = state
}

func (c *Control) IsControlled() bool {
	return c.Controlled
}

type APos struct {
	pixel.Vec
}

func (e *APos) GetPos() pixel.Vec {
	return e.Vec
}

func (e *APos) SetPos(pos pixel.Vec) {
	e.Vec = pos
}

type ARot struct {
	Rot float64
}

func (r *ARot) SetRot(rot float64) {
	r.Rot += rot
}

type Transform struct {
	APos
	ARot
}

func (e *Transform) GetMatrix() pixel.Matrix {
	return pixel.IM.Rotated(pixel.ZV, e.Rot).Moved(e.Vec)
}

type Determination struct {
	Visibility
	Transform
	Offset pixel.Vec
}

func (d *Determination) GetRect() pixel.Rect {
	size := d.GetSpriteSize() / 2
	pos := d.Add(d.Offset)
	return pixel.R(pos.X-size, pos.Y-size, pos.X+size, pos.Y+size)
}

func (d *Determination) Draw(t pixel.Target) {
	d.Sprite.Draw(t, d.GetMatrix().Moved(d.Offset))
}

type Mortality struct {
	Dead bool
}

func (m *Mortality) Kill() {
	m.Dead = true
}

func (m *Mortality) IsDead() bool {
	return m.Dead
}

type Movement struct {
	Vel pixel.Vec
}

func (s *Movement) GetVel() pixel.Vec {
	return s.Vel
}

type Solidity struct {
	Transform
	Volume
}

func (s *Solidity) GetRect() pixel.Rect {
	size := s.GetSize()
	return pixel.R(s.X-size, s.Y-size, s.X+size, s.Y+size)
}

type ATeam struct {
	Team uint16
}

func (t *ATeam) GetTeam() uint16 {
	return t.Team
}

type Harmfullness struct {
	Damage int
}

func (d *Harmfullness) GetDamage() int {
	return d.Damage
}
