package gui

import (
	"myFirstProject/game"
	"myFirstProject/game/world/graphics"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ButtonText struct {
	Text  string
	Scale pixel.Vec
}

// Button si gui button. It is component of Layer
type Button struct {
	Txt                                         ButtonText
	Idle, Pressed, Howered, Disable, prev, rect pixel.Rect
	Pos, Scale                                  pixel.Vec
	Disabled, Selected                          bool
	Call                                        func()
}

func (b *Button) textMatrix(mut pixel.Matrix) pixel.Matrix {
	return mut.Moved(b.Pos.Scaled(-1)).ScaledXY(pixel.ZV, pixel.V(1/b.Txt.Scale.X, 1/b.Txt.Scale.Y))
}

// NewButton creates new button
func NewButton(pos, scl pixel.Vec, idle, hovered, pressed, Disable pixel.Rect) *Button {
	w, h := idle.W()/2*scl.X, idle.H()/2*scl.Y
	button := &Button{
		Pos:     pos,
		Scale:   scl,
		Idle:    idle,
		Disable: idle,
		Pressed: idle,
		Howered: idle,
		rect:    pixel.R(pos.X-w, pos.Y-h, pos.X+w, pos.Y+h),
	}
	if hovered != pixel.ZR {
		button.Howered = hovered
	}
	if pressed != pixel.ZR {
		button.Pressed = pressed
	}
	if Disable != pixel.ZR {
		button.Disable = Disable
	}

	return button
}

// Update updates button
func (b *Button) Update() pixel.Rect {
	if b.Disabled {
		return b.Disable
	}

	if !b.rect.Contains(game.Win.MousePosition()) {
		b.Selected = false
		return b.Idle
	}

	if game.Win.JustPressed(pixelgl.MouseButtonLeft) || b.Selected {
		b.Selected = true
		if game.Win.JustReleased(pixelgl.MouseButtonLeft) {
			b.Selected = false
			b.Call()
		}
		return b.Pressed
	}

	return b.Howered
}

// GetPos returns position of button in game world so it is always at the same spot
// from view of camera
func (b *Button) GetPos() pixel.Vec {
	return graphics.Matrix.Project(b.Pos)
}
