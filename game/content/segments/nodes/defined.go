package nodes

import "github.com/faiface/pixel"

// all node types
var (
	Standard Type // most basic Node
)

// all drone types
var (
	StandardD DType // most basic drone
)

// Load loads all nodes and drones
func Load(reg []pixel.Rect) {
	loadDrones(reg)

	Standard = Type{
		PickupRange: 20,
		DType:       &StandardD,
	}
}

func loadDrones(reg []pixel.Rect) {
	StandardD = DType{
		Sight: 1000,
		Speed: 100,
		Rect:  reg[0],
	}
}
