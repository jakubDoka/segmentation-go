package effects

import (
	"libs/graphics/g2d/particles"
	"math"

	"github.com/faiface/pixel"
)

var (
	Trig   = NewSymmetricTriangle(1)
	Circle = NewSymmetricPolygon(1, 1, 1, 12)
	Bubble = NewSymmetricPolygon(1, .7, 0, 12)
	hoop   = NewSymmetricHoop(1, 12)
)

func NewSymmetricPolygon(size, in, out, sides float64) particles.Polygon {
	points := []pixel.Vec{}
	for i := 0.0; i < sides; i++ {
		points = append(points, pixel.V(size, 0).Rotated(i*2*math.Pi/sides))
	}
	return particles.Polygon{
		Points:       points,
		InnerOpacity: in,
		OuterOpacity: out,
		Cube:         particles.Cube{Radius: size},
	}
}

func NewSymmetricHoop(size, sides float64) particles.Hoop {
	points := []pixel.Vec{}
	for i := 0.0; i < sides; i++ {
		points = append(points, pixel.V(size, 0).Rotated(i*2*math.Pi/sides))
	}
	return particles.Hoop{
		Points: points,
		Cube:   particles.Cube{Radius: size},
	}
}

func NewSymmetricTriangle(size float64) particles.Triangle {
	points := [3]pixel.Vec{}
	for i := 0.0; i < 3; i++ {
		points[int(i)] = pixel.V(size, 0).Rotated(i * 2 * math.Pi / 3)
	}
	return particles.Triangle{
		Points: points,
		Cube:   particles.Cube{Radius: size},
	}
}
