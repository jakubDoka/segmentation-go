package graphics

import (
	"libs/mathm"
	"myFirstProject/game"
	"myFirstProject/game/world/ent"

	"github.com/faiface/pixel"
)

// Cam is a current game camera
var Cam Camera

// View is current viewport rect
var View pixel.Rect

// ParticleView is viewport for particles, its bigger then View for estetics
var ParticleView pixel.Rect

// Mouse is current mouse position
var Mouse pixel.Vec

// Matrix is current camera matrix
var Matrix pixel.Matrix

// Camera is game camera
type Camera struct {
	Scale pixel.Vec
	pixel.Vec
	MinZ, MaxZ, Sensitivity float64
	ent.Existence
}

// OnInput handles input
func (c *Camera) OnInput() {
	c.Scale = c.Scale.Scaled(1 + game.Win.MouseScroll().Y*c.Sensitivity)
	c.Scale.X = mathm.Clampf(c.Scale.X, c.MinZ, c.MaxZ)
	c.Scale.Y = mathm.Clampf(c.Scale.Y, c.MinZ, c.MaxZ)
}

// GetMatrix returns matrix representation of camera
func (c *Camera) GetMatrix() pixel.Matrix {
	return pixel.IM.ScaledXY(c.Vec, c.Scale).Moved(game.Win.Bounds().Center().Sub(c.Vec))
}

// GetRect viewport rect
func (c *Camera) GetRect() pixel.Rect {
	r := game.Win.Bounds()
	hw, hh := 1/c.Scale.X*r.W()/2, 1/c.Scale.Y*r.H()/2
	return pixel.R(c.X-hw, c.Y-hh, c.X+hw, c.Y+hh)
}

// GetGlobalMousePosition returns mouse position in global coordinates
func (c *Camera) GetGlobalMousePosition() pixel.Vec {
	return Matrix.Unproject(game.Win.MousePosition())
}

// Update ...
func (c *Camera) Update() {
	c.OnInput()
	View = c.GetRect()
	vec := View.Min.To(View.Max).Scaled(.3)
	ParticleView = pixel.Rect{
		Min: View.Min.Sub(vec),
		Max: View.Max.Add(vec),
	}
	Matrix = c.GetMatrix()
	game.Win.SetMatrix(Matrix)
	Mouse = c.GetGlobalMousePosition()
	if c.Existence == nil {
		return
	}
	c.Vec = pixel.Lerp(c.Vec, c.GetPos(), game.Delta*10)

}
