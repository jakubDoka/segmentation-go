package effects

import "github.com/faiface/pixel"

type PointGen struct{}

func (p PointGen) Pos() pixel.Vec {
	return pixel.ZV
}
