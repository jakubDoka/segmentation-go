package worms

import (
	"myFirstProject/game"
	"myFirstProject/game/content/segments"
	"myFirstProject/game/world/graphics"
	"myFirstProject/game/world/tilemap"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// InputState is for translating client input to server
type InputState struct {
	Shoot   bool
	Pointer pixel.Vec
}

// Shooting returns whether worm is shooting. Its for input handling
func (i *InputState) Shooting() bool {
	return i.Shoot
}

// GetPointer returns where is client that is controling worm pointing with his mouse
func (i *InputState) GetPointer() pixel.Vec {
	return i.Pointer
}

// SegmentInput handles players interaction with segments
func (e *Entity) SegmentInput() {
	selectPressed := game.Win.JustPressed(pixelgl.MouseButtonRight)
	usePressed := game.Win.JustPressed(pixelgl.MouseButtonLeft)
	transMode := game.Win.Pressed(pixelgl.KeyLeftShift)

	prev := game.Selected
	prevT := game.TransSelected

	selectedT := prevT != -1
	selected := prev != -1

	var howered *segments.Entity

	for i, s := range e.Segments {
		if !s.Determination.GetRect().Contains(graphics.Mouse) {
			continue
		}

		howered = s
		if selectPressed {
			if transMode {
				if s.Deployed {
					game.TransSelected = i
					game.Selected = -1
				} else {
					//TODO alert
				}
			} else {
				game.Selected = i
				game.TransSelected = -1
			}
			break
		}

		if !usePressed {
			break
		}

		if selectedT && s.Deployed {
			if i == game.TransSelected {
				if len(s.Node.Cons) != 0 {
					s.Node.CutAllCons()
				} else {
					s.ConAll()
				}
			} else {
				seg := e.Segments[game.TransSelected]
				if seg.Node.Cons[howered.Node] {
					seg.Node.CutConnection(howered.Node)
				} else if seg.Node.Recs[howered.Node] {
					// TODO Warm this is not possible
				} else {
					seg.Node.Connect(howered.Node)
				}
			}
		} else if selected {
			if i == game.Selected {
				if s.Independent() {
					s.Pickup()
				}
				e.Move(game.Selected, -1)
			} else if !s.Independent() {
				e.Segments[game.Selected].Pickup()
				e.Move(game.Selected, i)
			}
		}
		break
	}

	if game.Selected != -1 {
		seg := e.Segments[game.Selected]
		tile := tilemap.World.GetTileByPos(graphics.Mouse)
		canPlace := (seg.PlaceFilter == nil || seg.CanPlace(tile)) && !tile.Occupied
		if howered == nil && prev == game.Selected {
			if !seg.GetRect().Contains(graphics.Mouse) {
				if usePressed {
					if canPlace {
						if seg.Deployed {
							tilemap.World.OnTileChange(seg.Module, true, false, e.Team, 0)
						}
						e.Place(game.Selected)
						seg.Goto(tile.Center())
					} else {
						//TODO warm player
					}

				} else if selectPressed {
					game.Selected = -1
				}
			}
		}
		e.DrawMoveSelection(seg, howered, canPlace)
	} else if game.TransSelected != -1 {
		seg := e.Segments[game.TransSelected]
		if howered == nil && selectPressed {
			game.TransSelected = -1
		}
		e.DrawTransSelection(seg, howered)
	}
}

// DrawTransSelection draws ui selection of segments for connecting transport paths
func (e *Entity) DrawTransSelection(seg, howered *segments.Entity) {
	rect := seg.Determination.GetRect()
	pos := seg.Determination.Vec
	game.Im.Push(rect.Min, rect.Max)
	game.Im.Rectangle(5)
	for c := range seg.Node.Cons {
		game.Im.Push(pos, c.Vec)
		game.Im.Line(10)
		game.Im.Push(c.Vec)
		game.Im.Circle(20, 10)
	}
	if howered != nil && howered.Deployed && howered != seg {
		game.Im.Push(howered.Determination.Vec)
		game.Im.Circle(40, 10)
	}
}

// DrawMoveSelection draws ui selection of segments for changing order and deploying
func (e *Entity) DrawMoveSelection(seg, howered *segments.Entity, canPlace bool) {
	pos := seg.Determination.Vec
	game.Im.SetColorMask(SelectionColor)
	game.Im.Push(pos)
	radius := seg.Determination.GetSpriteSize() / 2
	game.Im.Circle(radius, 3)
	var dest pixel.Vec
	if howered != nil {
		if howered != seg {
			dest = howered.Back
			if game.Selected != 0 && howered == e.Segments[game.Selected-1] {
				return
			}
		} else {
			dest = e.Segments[0].Vec
			if game.Selected == 0 {
				return
			}
		}
	} else {
		tile := tilemap.World.GetTileByPos(graphics.Mouse)
		if !canPlace {
			v := tile.Vertices()
			game.Im.Push(v[0], v[2])
			game.Im.Line(10)
			game.Im.Push(v[1], v[3])
			game.Im.Line(10)
		}
		game.Im.Push(tile.Min, tile.Max)
		game.Im.Rectangle(5)
		dest = tile.Center()
	}
	game.Im.Push(dest)
	game.Im.Circle(5, 0)
	game.Im.Push(pos.Add(pos.To(dest).Unit().Scaled(radius)))
	game.Im.Push(dest)
	game.Im.Line(3)
}

// Input handles input related to movement of worm and shooting
func (e *Entity) Input() {
	e.InputState.Pointer = graphics.Mouse
	e.InputState.Shoot = game.Win.Pressed(pixelgl.MouseButtonLeft) && game.ShouldShoot()
	if game.IsNetworking {
		e.WriteInput()
	}
}

// MovementInput handles movement of worm
func (e *Entity) MovementInput() {
	if game.Win.Pressed(pixelgl.KeyA) {
		e.Rot += e.Steer * game.Delta
	} else if game.Win.Pressed(pixelgl.KeyD) {
		e.Rot -= e.Steer * game.Delta
	}

	if game.Win.Pressed(pixelgl.KeyW) {
		e.Charge()
	}
}
