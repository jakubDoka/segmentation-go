package targets

import (
	"myFirstProject/game/world/ent/ai"

	"github.com/faiface/pixel"
)

func IsValidTarget(t ai.Target, pos pixel.Vec, rad float64) bool {
	return t != nil && !t.IsDead() && t.GetPos().To(pos).Len() < rad

}
