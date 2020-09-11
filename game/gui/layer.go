package gui

import (
	"libs/graphics/g2d/particles"
	"libs/textbox"
	"myFirstProject/game"
	"myFirstProject/game/world/graphics"

	"github.com/faiface/pixel"
)

var zr = struct {
	Position  pixel.Vec
	Color     pixel.RGBA
	Picture   pixel.Vec
	Intensity float64
}{Color: pixel.Alpha(1), Intensity: 1}

// Layer can hold buttons or layers
type Layer struct {
	txtDrawer *textbox.TextDrawer
	Trig      pixel.TrianglesData
	pixel.Drawer
	TexBoxes []*textbox.Textbox
	Buttons  []*Button
	Hidden   bool
	pixel.Matrix
}

// New is layer constructor
func New(pic pixel.Picture, txtDrawer *textbox.TextDrawer) *Layer {
	layer := Layer{txtDrawer: txtDrawer}
	layer.Picture = pic
	layer.Triangles = &layer.Trig
	return &layer
}

// Draw draws layer
func (l *Layer) Draw(t pixel.Target, m pixel.Matrix) {
	if l.Hidden {
		return
	}

	l.Matrix = m

	for i, c := range l.Buttons {
		reg := c.Update()

		c.prev = reg

		var (
			w, h    = reg.W() / 2, reg.H() / 2
			bounds  = reg.Vertices()
			corners = [4]pixel.Vec{
				pixel.V(w, -h),
				pixel.V(w, h),
				pixel.V(-w, h),
				pixel.V(-w, -h),
			}
		)

		i *= 6
		for j, v := range particles.RectTrigPattern {
			j = j + i
			l.Trig[j].Position = m.Unproject(c.Pos.Add(corners[v].ScaledXY(c.Scale)))
			l.Trig[j].Picture = bounds[v]
		}

		l.Dirty()
	}

	l.Drawer.Draw(t)

	for _, b := range l.Buttons {
		l.txtDrawer.Write(b.Txt.Text, b.textMatrix(m), pixel.ZV)
	}

	l.drawText(m)
}

func (l *Layer) drawText(m pixel.Matrix) {
	for _, t := range l.TexBoxes {
		t.Input(game.Win, graphics.Mouse, game.Delta)
		t.Draw(game.Win, m)
	}
}

// Add adds button to layer
func (l *Layer) Add(button *Button) {
	l.Buttons = append(l.Buttons, button)
	for i := 0; i < 6; i++ {
		l.Trig = append(l.Trig, zr)
	}
}

// AddText adds textbox to layer
func (l *Layer) AddText(t *textbox.Textbox) {
	l.TexBoxes = append(l.TexBoxes, t)
}
