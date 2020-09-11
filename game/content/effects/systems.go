package effects

import (
	"libs/graphics/g2d/particles"
	"libs/graphics/g2d/particles/curves"
	"math"
	"myFirstProject/game"
	"sync"

	"github.com/faiface/pixel"
)

// Batch where all effects are drawn
var Batch *pixel.Batch

var mut sync.Mutex

var systems []*particles.System

// all effect types
var (
	DuoBulletTrail    *particles.System
	DuoBulletExpWave  *particles.System
	DuoBulletExpSharp *particles.System
	DuoShootExp       *particles.System
)

// Load loads all effects
func Load() {
	Batch = pixel.NewBatch(&pixel.TrianglesData{}, nil)

	DuoShootExp = particles.New(PointGen{},
		particles.Behavior{
			Scale:              particles.Prop(12, curves.BumpBig),
			Velocity:           0,
			VelocityRandomness: 50,
			Livetime:           .3,
			LivetimeRandomness: .1,
			Rotation:           math.Pi / 4,
			SpreadAngle:        math.Pi / 1.5,
			Color: curves.NewColor(
				pixel.RGBA{R: 1, G: .8, B: .7, A: .8},
				pixel.RGBA{R: 0.08, G: .5, B: .54, A: .8},
				curves.New(pixel.ZV, pixel.V(0, 1), pixel.V(2, 1), pixel.V(-1, 0)),
			),
		},
		Bubble.Cube,
	)
	systems = append(systems, DuoShootExp)

	DuoBulletExpWave = particles.New(PointGen{},
		particles.Behavior{
			Scale:              particles.Prop(40, curves.BumpBig),
			ScaleRandomness:    20,
			Livetime:           .3,
			LivetimeRandomness: -.2,
			Rotation:           math.Pi / 4,
			SpreadAngle:        math.Pi,
			Color:              DuoShootExp.Color,
		},
		hoop,
	)
	systems = append(systems, DuoBulletExpWave)

	DuoBulletExpSharp = particles.New(PointGen{},
		particles.Behavior{
			Scale:              particles.Prop(6, curves.None),
			Velocity:           100,
			VelocityRandomness: 100,
			ScaleRandomness:    1,
			Livetime:           .2,
			LivetimeRandomness: .1,
			Twerk:              10,
			Rotation:           math.Pi / 4,
			SpreadAngle:        math.Pi,
			Color: curves.NewColor(
				pixel.RGBA{R: 1, G: .8, B: .7, A: .8},
				pixel.RGBA{R: 1, G: .8, B: .7, A: .8},
				curves.None,
			),
		},
		Trig,
	)
	systems = append(systems, DuoBulletExpSharp)

	DuoBulletTrail = particles.New(PointGen{},
		particles.Behavior{
			Scale:    particles.Prop(7, curves.LinearDecreasing),
			Livetime: .5,
			Rotation: math.Pi / 4,
			Color:    DuoShootExp.Color,
		},
		Bubble.Cube,
	)
	systems = append(systems, DuoBulletTrail)
}

// Update updates all systems
func Update(delta float64) {
	for _, sys := range systems {
		go sys.Update(delta, pixel.ZR)
	}
}

// DrawToBatch draws all particles to batch on thread
func DrawToBatch() {
	mut.Lock()
	go func() {
		for _, sys := range systems {
			sys.Draw(Batch)
		}
		mut.Unlock()
	}()
}

// Draw draws the batch to window
func Draw() {
	mut.Lock()
	Batch.Draw(game.Win)
	Batch.Clear()
	mut.Unlock()
}
