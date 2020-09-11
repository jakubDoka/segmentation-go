package bullets

import (
	"libs/graphics/g2d/textures"

	"github.com/faiface/pixel"
)

// Batch is batch where all bullets are drawn
var Batch *pixel.Batch

// Sheet contains bullet textures
var Sheet textures.PixSheet

// All bullet types
var (
	Standard Type // Standard is most basic bullet
)

// StandardBehavior si just going forward
func StandardBehavior(e *Entity) {
	e.Move()
	e.DealDamage()
}

// Load loads all bullet types
func Load(sheet textures.PixSheet) {
	Batch = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Picture)
	Sheet = sheet
	Standard = Init(Type{
		Speed:          700,
		Livetime:       1,
		Size:           10,
		TrailFrequency: 1.0 / 60.0,
		Damage:         3,
		Rect:           Sheet.Regs[0],
		Behavior:       StandardBehavior,
		//Trail:          effects.DuoTrail,
		//StartExplosion: effects.DuoShoot,
		//EndExplosion:   effects.DuoExp,
	})
}
