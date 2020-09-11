package content

import (
	"image/color"
	"image/draw"
	"libs/graphics/g2d/colors"
	"libs/graphics/g2d/textures"
	"math/rand"
	"myFirstProject/game"
	"myFirstProject/game/content/bullets"
	"myFirstProject/game/content/effects"
	"myFirstProject/game/content/modules/turrets"
	"myFirstProject/game/content/segments"
	"myFirstProject/game/content/segments/nodes"
	"myFirstProject/game/content/worms"
	"time"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

var Mixer = .1
var UiSheet textures.PixSheet
var Rand *rand.Rand

func drawAndClear(b *pixel.Batch, t pixel.Target) {
	b.Draw(t)
	b.Clear()
}

// Draw draws everithing
func Draw() {
	effects.DrawToBatch()
	drawAndClear(segments.Turrets, game.Win)
	drawAndClear(segments.Worms, game.Win)
	drawAndClear(bullets.Batch, game.Win)
	effects.Draw()
	game.Im.Draw(game.Win)
	game.Im.Clear()
}

func shadeIMG(img draw.Image) {
	textures.Shade(img, img.Bounds(), textures.Circular, false, 4, 100, 20, colornames.Black)
}

// RandomColor returns random color
func RandomColor() pixel.RGBA {
	return pixel.RGB(Rand.Float64(), Rand.Float64(), Rand.Float64())
}

func mix(img draw.Image) {
	b, t := RandomColor(), Rand.Float64()
	textures.ForEachFull(img, func(col color.Color) color.Color {
		return colors.Mix(col, b, t)
	})
}

// Load loads ewerithing
func Load() {
	Rand = rand.New(rand.NewSource(time.Now().Unix()))
	parts, _ := textures.New("assets/textures/parts_1.png", 64, 64)
	bases, _ := textures.New("assets/textures/bases_1.png", 72, 72)
	_bullets, _ := textures.New("assets/textures/bullets.png", 64, 64)
	buttons, _ := textures.New("assets/textures/ui.png", 134, 134)

	//Shading
	shadeIMG(parts)
	shadeIMG(bases)
	shadeIMG(buttons)

	//mixing
	pressedButtons := buttons.Copy()
	mix(pressedButtons.Image)

	//Merging
	_worms := textures.Merge(parts, bases)
	mix(_worms.Image)
	ui := textures.Merge(buttons, pressedButtons)
	ui.Crop()
	UiSheet = *textures.FromHSheet(ui)

	//Converting
	pixWorms := textures.FromHSheet(_worms)

	effects.Load()
	bullets.Load(*textures.FromSheet(_bullets))
	turrets.Load(pixWorms.Regs)
	nodes.Load(pixWorms.Regs)
	segments.Load(*pixWorms)
	worms.Load()
}
