package turrets

import (
	"myFirstProject/game/content/bullets"
	"myFirstProject/game/content/segments"

	"github.com/faiface/pixel"
)

// all turret types
var (
	Standard Type
	ID       uint8
)

// Module is container of types
type Module []*Type

// Get implements modules.Module interface
func (m Module) Get(i uint8) segments.ModuleType {
	return m[i]
}

// Types is Type array for networking purposes
var Types = Module{}

// StandardShoot is most basic shooting pattern: pew reload pew
func StandardShoot(e *Entity, dir float64) {
	e.FireBullet(dir)
}

// Load loads all Turret types
func Load(reg []pixel.Rect) {
	Standard = Type{
		ID:           0,
		Rect:         reg[1],
		Speed:        3,
		Type:         &bullets.Standard,
		Hull:         20,
		Offset:       16,
		ReloadSpeed:  .7,
		Inaccuracy:   .05,
		VelInacuracy: 100,
		MaxAmmo:      20,
		AmmoPerShot:  0,
		Pushback:     4,
		FixSpeed:     50,
		DeployBoost:  2,
		Shoot:        StandardShoot,
	}
	Types = append(Types, &Standard)

	segments.Modules = append(segments.Modules, Types)
}
