package worms

import (
	"math"
	"myFirstProject/game/content/segments"
)

// all worm types
var (
	Standard Type
)

// Types is Type map for networking purposes
var Types = map[uint8]*Type{}

// Load loads all worms
func Load() {
	Standard = Type{
		ID:           0,
		Ends:         &segments.StandardEnd,
		Acceleration: 500,
		Steer:        math.Phi,
		MaxSpeed:     500,
	}
	Types[0] = &Standard
}
