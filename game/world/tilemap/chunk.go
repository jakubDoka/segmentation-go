package tilemap

import (
	"myFirstProject/game"

	"github.com/faiface/pixel"
)

// Chunk is for effective map rendering
type Chunk struct {
	Tiles []*Tile
	*pixel.Batch
	pixel.Rect
}

// DrawTiles Draws all tiles to chunks Batch
func (c *Chunk) DrawTiles() {
	c.Clear()
	for _, t := range c.Tiles {
		t.Draw(c)
	}
}

// Draw draws the Chunk
func (c *Chunk) Draw() {
	c.Batch.Draw(game.Env)
}
