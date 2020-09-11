package segments

import (
	"libs/graphics/g2d/textures"
	"myFirstProject/game/content/segments/nodes"

	"github.com/faiface/pixel"
)

// Batches and textures
var (
	Turrets *pixel.Batch
	Worms   *pixel.Batch
	Sheet   textures.PixSheet
)

// all segment types
var (
	Standard, StandardEnd Type
)

// Types is Type map for networking purposes
var Types = map[uint8]*Type{}

// Load loads all segments
func Load(sheet textures.PixSheet) {
	Turrets = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Picture)
	Worms = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Picture)
	Sheet = sheet
	Standard = Type{
		ID:           0,
		Type:         &nodes.Standard,
		Base:         64,
		Size:         20,
		Speed:        300,
		MaxHealth:    100,
		DeployTime:   3,
		SegmentPiece: Sheet.Regs[2],
		BasePiece:    Sheet.Regs[0],
		CocPiece:     Sheet.Regs[3],
	}
	Types[0] = &Standard
	StandardEnd = Type{
		ID:           1,
		Type:         &nodes.Standard,
		Base:         64,
		Speed:        500,
		SegmentPiece: Sheet.Regs[4],
		BasePiece:    Sheet.Regs[0],
		CocPiece:     Sheet.Regs[3],
	}
	Types[1] = &StandardEnd
}
