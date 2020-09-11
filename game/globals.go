package game

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// IsServer indicate whether the game instance is server
var IsServer bool

// Halted if someon connects whole process have to stop to ensure tha sync is accurate
var Halted int32

// IsNetworking indicate whether game is connected to network
var IsNetworking bool

// LastUpdate is time when last udp update from server wos received
var LastUpdate time.Time

// UpdateSpacing is time between remote updates
var UpdateSpacing float64

// Alfa is progress of remote updates used for interpolation
var Alfa float64

// Delta is frame time
var Delta float64 = 1 / 60

// Bullets renders bullets
var Bullets *pixel.Batch

// Worms renders all worm moving segments
var Worms *pixel.Batch

// Turrets renders all turrets
var Turrets *pixel.Batch

// Env renders environment
var Env *pixel.Batch

// Textures is spritesheet with all textures
var Textures pixel.Picture

// GuiTextures is spritesheet that contains all gui textures
var GuiTextures pixel.Picture

// Win is game window
var Win *pixelgl.Window

// PlayerTeam is team where local player is
var PlayerTeam uint16

// PlayerID is is id of worm player is controling
var PlayerID uint16

// Selected determinate witch segment has player selected
var Selected = -1

// TransSelected determinate with segment is selected for transport config
var TransSelected = -1

// Im is global IMDraw instance
var Im *imdraw.IMDraw = imdraw.New(nil)

// ShouldShoot returns whether turrets should react on players input
func ShouldShoot() bool {
	return Selected == -1 && TransSelected == -1
}
