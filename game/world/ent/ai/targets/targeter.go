package targets

/*import (
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/ent/ai"

	"github.com/faiface/pixel"
)

// TargetType is enum type for withdrawing groups of targets from targeter
type TargetType int

// TargetType enum
const (
	UNIT TargetType = iota
	CORE
	SEGMENT
)

// TargetTypes is slice of all targetTypes
var TargetTypes []TargetType = []TargetType{
	UNIT,
	CORE,
	SEGMENT,
}

// Finder is an entry to general purpose targeter
var Finder *Targeter

// Targeter finds target for entity and manages targets
type Targeter map[int][][][]ent.Existence

// AddTeam adds new team to Targeter so you can insert objects freely
func (t *Targeter) AddTeam(id int) {
	tv := *t
	tv[id] = make([][][]ent.Existence, len(ai.Layers))
	for i := range tv[id] {
		tv[id][i] = make([][]ent.Existence, len(TargetTypes))
	}
}

// GetAirTarget finds a best target for air unit from a given position
func (t *Targeter) GetAirTarget(pos pixel.Vec, team int, tp ...TargetType) ent.Existence {
	var best ent.Existence
	for k, v := range *t {
		if k != team {
			for _, l := range v {
				for _, i := range tp {
					for _, e := range l[i] {
						if best == nil || e.GetPos().To(pos).Len() < best.GetPos().To(pos).Len() {
							best = e
						}
					}
				}
			}
		}
	}
	return best
}
*/
