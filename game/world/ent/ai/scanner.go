package ai

import (
	"libs/collizions"
	"libs/mathm/floats"
	"myFirstProject/game/world/ent"

	"github.com/faiface/pixel"
)

// Scanner stores radar for general use
var Scanner Radar = Radar{collizions.New(64000, 64000, 128)}

// Radar is Group of quadtree that sorts collidables to teams and layers
type Radar struct {
	*collizions.Graph
}

// GetTarget gets the closest target in the range from position, team and layers are filters
func (r Radar) GetTarget(sight float64, pos pixel.Vec, team uint16) Target {
	found := r.GetColliding(floats.Cube(pos, sight))

	if len(found) == 0 {
		return nil
	}

	candidates := []Target{}
	for _, c := range found {
		t, ok := c.(Target)
		if ok && t.GetTeam() != team {
			candidates = append(candidates, t)
		}
	}

	count := len(candidates)
	if count == 0 {
		return nil
	}

	best := candidates[0]
	for i := 1; i < count; i++ {
		if best.GetPos().To(pos).Len() > candidates[i].GetPos().To(pos).Len() {
			best = candidates[i]
		}
	}

	if pos.To(best.GetPos()).Len() > sight {
		return nil
	}

	return best
}

// GetBudsInRange gets the closest target in the range from position, team and layers are filters
func (r Radar) GetBudsInRange(sight float64, pos pixel.Vec, team uint16) []collizions.Shape {
	found := r.GetColliding(pixel.R(pos.X-sight, pos.Y-sight, pos.X+sight, pos.Y+sight))

	count := len(found)
	if count == 0 {
		return nil
	}

	i := 0
	for _, c := range found {
		if pos.To(c.GetRect().Center()).Len() > sight {
			continue
		}

		co, ok := c.(ent.Competence)
		if !ok || co.GetTeam() != team {
			continue
		}

		found[i] = c
		i++
	}

	for j := i; j < count; j++ {
		found[j] = nil
	}

	return found[:i]
}

// GetCollidingEnemy feeds the slice with all enemy colliding shapes with giver rect
func (r Radar) GetCollidingEnemy(rect pixel.Rect, team uint16) (shapes []collizions.Shape) {
	shapes = r.GetColliding(rect)

	i := 0
	for _, c := range shapes {
		co, ok := c.(ent.Competence)
		if !ok || co.GetTeam() == team {
			continue
		}

		shapes[i] = c
		i++
	}

	for j := i; j < len(shapes); j++ {
		shapes[j] = nil
	}

	return shapes[:i]
}
