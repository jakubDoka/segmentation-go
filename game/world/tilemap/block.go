package tilemap

import "github.com/faiface/pixel"

// Block is a tile layer content. For rendering and getting cost
type Block struct {
	pixel.Sprite
	Cost int
}

// Draw draws the block
func (b *Block) Draw(target pixel.Target, pos *pixel.Vec) {
	b.Sprite.Draw(target, pixel.Matrix{1, 0, 0, 1, pos.X, pos.Y})
}
