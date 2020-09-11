package ai

import (
	"math"
	"myFirstProject/game/world/ent"

	"github.com/faiface/pixel"
)

// Target ...
type Target interface {
	ent.Existence
	ent.Velocity
	ent.Finite
	ent.Competence
}

// Shooter ...
type Shooter interface {
	ent.Speed
	Target
}

// Predict takes shooter and target and returns point where shooter should
// shoot to in order to hit moving target
func Predict(shooter Shooter, target Target) (pixel.Vec, bool) {
	return ProceduralPredict(shooter.GetPos(), target.GetPos(), target.GetVel().Sub(shooter.GetVel()), shooter.GetSpeed())
}

// ProceduralPredict does the same thing as Predict it juts takes a raw arguments
func ProceduralPredict(shooter, target, targetVelocity pixel.Vec, bulletSpeed float64) (pixel.Vec, bool) {
	d := target.Sub(shooter)

	a := targetVelocity.X*targetVelocity.X + targetVelocity.Y*targetVelocity.Y - bulletSpeed*bulletSpeed
	b := 2 * (d.X*targetVelocity.X + d.Y*targetVelocity.Y)
	c := d.X*d.X + d.Y*d.Y

	cof := b*b - 4*a*c

	if cof < 0 {
		return pixel.Vec{}, false
	}

	dis := math.Sqrt(cof)
	a *= 2
	t1, t2 := (-b+dis)/a, (-b-dis)/a

	if t1 >= 0 && (t1 < t2 || t2 < 0) {
		return target.Add(targetVelocity.Scaled(t1)), true
	}

	return target.Add(targetVelocity.Scaled(t2)), true
}
