package effects

import (
	"myFirstProject/game/world/graphics"

	"github.com/faiface/pixel"
)

func DuoShoot(pos, dir pixel.Vec) {
	if graphics.ParticleView.Contains(pos) {
		DuoShootExp.Explode(5, pos, dir)
	}
}

func DuoTrail(pos, dir pixel.Vec) {
	if graphics.ParticleView.Contains(pos) {
		DuoBulletTrail.Explode(1, pos, dir)
	}
}

func DuoExp(pos, dir pixel.Vec) {
	if graphics.ParticleView.Contains(pos) {
		DuoBulletExpWave.Explode(5, pos, dir)
		//DuoBulletExpSharp.Explode(6, pos, dir)
	}
}
