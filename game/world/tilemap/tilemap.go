package tilemap

import (
	"libs/mathm/ints"
	"libs/mathm/points"
	"myFirstProject/game"
	"myFirstProject/game/world/ent"
	"myFirstProject/game/world/graphics"

	"github.com/faiface/pixel"
)

// TargetType ...
type TargetType int

// TargetType enum
const (
	Core TargetType = iota
	Defence
	Producer
	Transporter
)

// TargetTypes are values of TargetType enum
var TargetTypes []TargetType = []TargetType{
	Core,
	Defence,
	Producer,
	Transporter,
}

// Batch is batch where tile map is drawn
var Batch pixel.Batch

// Sheet contains map textures
var Sheet pixel.Picture

// World is instance of current game world
var World *Tilemap = New(1000, 1000, 10, 64)

// Tilemap offers efficient rendering of tiles and pathfinding
type Tilemap struct {
	pixel.Rect
	points.Point
	Tiles               [][]*Tile
	Chunks              [][]*Chunk
	TexRegions          []pixel.Rect
	Sheet               pixel.Picture
	TileSize, ChunkSize pixel.Vec
	Targets
	Paths
}

// DemageSource is used for mapping tile costs
type DemageSource interface {
	ent.Strength
	ent.Averness
	ent.Existence
}

// New creates instance of tilemap
func New(w, h, resolution int, tileSize float64) *Tilemap {
	tilemap := &Tilemap{
		TileSize: pixel.V(tileSize, tileSize),
		Point:    points.New(w, h),
		Paths:    map[uint16]*TeamPath{},
		Targets:  map[uint16][][]points.Point{},
	}
	tilemap.GenerateTiles(w, h)
	tilemap.GenerateChunks(resolution)

	return tilemap
}

// Draw draws visible part of tilemap
func (t *Tilemap) Draw() {
	minX, maxX := int(graphics.View.Min.X/t.ChunkSize.X), int(graphics.View.Max.X/t.ChunkSize.X)
	minY, maxY := int(graphics.View.Min.Y/t.ChunkSize.Y), int(graphics.View.Max.Y/t.ChunkSize.Y)
	for y := minY; y < maxY; y++ {
		if y == t.R {
			break
		}
		if y < 0 {
			continue
		}
		for x := minX; x < maxX; x++ {
			if x == t.C {
				break
			}
			if x < 0 {
				continue
			}

			t.Chunks[y][x].Draw()
		}
	}
}

// GetStep returns next position to go to
func (t *Tilemap) GetStep(team uint16, kind int, pos pixel.Vec) pixel.Vec {
	return t.CenterOf(t.Paths.GetStep(team, kind, t.GetMapPoint(pos)))
}

// AddTeam adds team to world
func (t *Tilemap) AddTeam(team uint16) {
	t.Targets.AddTeam(team)
	t.Paths.AddTeam(t, team)
}

// OnTileChange handles addition or removal of demage source
func (t *Tilemap) OnTileChange(ent DemageSource, remove, update bool, team uint16, kind int) {
	pos := ent.GetPos()
	//Updateing targets
	point := t.GetMapPoint(pos)
	tile := t.GetTile(point)
	if remove {
		t.Remove(team, kind, point)
		tile.Occupied = false
	} else {
		t.Targets.Add(team, kind, point)
		tile.Occupied = true
	}

	// updating pathfinder
	places := t.GetPointsInRange(pos, ent.GetSight()+400)
	for _, p := range places {
		tile := t.GetTile(p)
		if remove {
			delete(tile.InRange, ent)
		} else {
			tile.InRange[ent] = team
		}

	}
	if !update {
		return
	}
	t.Paths.Update(t, team, places, remove)
}

// DrawTiles is good for small tile maps
func (t *Tilemap) DrawTiles(target pixel.Target) {
	for _, row := range t.Tiles {
		for _, ti := range row {
			ti.Draw(game.Env)
		}
	}
}

// GenerateChunks does not take to account the reminder, don't use 3 x 3 if map has 4 x 4 tiles
// use 2 x 2. Pregenerate chunks before you use DrawChunks method. Resolution bigger then
// map size is redundant.
func (t *Tilemap) GenerateChunks(resolution int) {
	ChunkSizeX, ChunkSizeY := t.C/resolution, t.R/resolution
	ChunkWidth, ChunkHeight := float64(ChunkSizeX)*t.TileSize.X, float64(ChunkSizeY)*t.TileSize.Y
	t.ChunkSize = pixel.V(ChunkWidth, ChunkHeight)

	for y := range t.Chunks {
		t.Chunks[y] = make([]*Chunk, resolution)
		for x := range t.Chunks[y] {
			orx, ory := float64(x)*ChunkWidth, float64(y)*ChunkHeight
			chunk := Chunk{
				Batch: pixel.NewBatch(&pixel.TrianglesData{}, t.Sheet),
				Rect:  pixel.R(orx, ory, orx+ChunkWidth, ory+ChunkHeight),
				Tiles: make([]*Tile, 0, ChunkSizeX*ChunkSizeY),
			}

			beginX, beginY := ChunkSizeX*x, ChunkSizeY*y
			endX, endY := beginX+ChunkSizeX, beginY+ChunkSizeY
			for r := beginY; r < endY; r++ {
				for c := beginX; c < endX; c++ {
					chunk.Tiles = append(chunk.Tiles, t.Tiles[r][c])
				}
			}
			t.Chunks[y][x] = &chunk
		}
	}
}

// GenerateTiles generates tile map
func (t *Tilemap) GenerateTiles(w, h int) {
	t.Rect = pixel.R(0, 0, float64(w)*t.TileSize.X, float64(h)*t.TileSize.Y)
	t.Tiles = make([][]*Tile, h)
	for y := range t.Tiles {
		row := make([]*Tile, w)
		for x := range row {
			orx, ory := float64(x)*t.TileSize.X, float64(y)*t.TileSize.Y
			row[x] = &Tile{
				Rect:    pixel.R(orx, ory, orx+t.TileSize.X, ory+t.TileSize.Y),
				Layers:  make([]*Block, 3),
				InRange: map[ent.Strength]uint16{},
				X:       x,
				Y:       y,
			}
		}
		t.Tiles[y] = row
	}
}

// GetMapPoint converts vector to tilemap coordinates
func (t *Tilemap) GetMapPoint(pos pixel.Vec) points.Point {
	return points.New(
		ints.Clamp(int(pos.X/t.TileSize.X), 0, t.C-1),
		ints.Clamp(int(pos.Y/t.TileSize.Y), 0, t.R-1),
	)
}

// GetTile gets tile of tile map
func (t *Tilemap) GetTile(pos points.Point) *Tile {
	if pos.C < 0 || pos.R < 0 || pos.C >= t.C || pos.R >= t.R {
		return nil
	}
	return t.Tiles[pos.R][pos.C]
}

// GetTileByPos returns the tile that contains given vector
func (t *Tilemap) GetTileByPos(pos pixel.Vec) *Tile {
	return t.GetTile(t.GetMapPoint(pos))
}

// SetTile sets tile overlay, floor and block
func (t *Tilemap) SetTile(x, y int, floor, overlay, block *Block) {
	tile := t.Tiles[y][x]
	tile.Layers[FLOOR] = floor
	tile.Layers[OVERLAY] = overlay
	tile.Layers[BLOCK] = block
}

// CenterOf returns real position of point
func (t *Tilemap) CenterOf(pos points.Point) pixel.Vec {
	return pixel.V(
		float64(pos.C)*t.TileSize.X+t.TileSize.X/2,
		float64(pos.R)*t.TileSize.Y+t.TileSize.Y/2,
	)
}

// GetPointsInRange returns all tiles that are in given circle
func (t *Tilemap) GetPointsInRange(pos pixel.Vec, r float64) []points.Point {
	minX, minY, maxX, maxY := t.getTileRect(pos.Add(pixel.V(-r, -r)), pos.Add(pixel.V(r, r)))

	res := make([]points.Point, 0, (maxX-minX)*(maxY-minY))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := points.New(x, y)
			if t.CenterOf(p).To(pos).Len() > r {
				continue
			}

			res = append(res, p)

		}
	}

	return res
}

// GetIntersectingTiles gets all tiles that itersect rectangle
func (t *Tilemap) GetIntersectingTiles(rect pixel.Rect) []*Tile {
	minX, minY, maxX, maxY := t.getTileRect(rect.Min, rect.Max)

	res := make([]*Tile, 0, (maxX-minX)*(maxY-minY))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			res = append(res, t.Tiles[y][x])
		}
	}
	return res
}

func (t *Tilemap) getTileRect(min, max pixel.Vec) (int, int, int, int) {
	minp := t.GetMapPoint(min)
	maxp := t.GetMapPoint(max)

	return ints.Max(0, minp.C), ints.Max(0, minp.R), ints.Min(t.C, maxp.C), ints.Min(t.R, maxp.R)
}

// GetNeighbors gets all tiles that touches given tilemap point
func (t *Tilemap) GetNeighbors(pos points.Point) []*Tile {
	return []*Tile{
		t.GetTile(pos.Add(points.D4[0])),
		t.GetTile(pos.Add(points.D4[1])),
		t.GetTile(pos.Add(points.D4[2])),
		t.GetTile(pos.Add(points.D4[3])),
	}
}

// GetCosts retruns all tiles costs
func (t *Tilemap) GetCosts(team uint16) [][]int {
	costs := make([][]int, len(t.Tiles))
	for y, row := range t.Tiles {
		costs[y] = make([]int, len(t.Tiles[0]))
		for x, ti := range row {
			costs[y][x] = ti.GetCost(team)
		}
	}
	return costs
}

// ForEach is action for each tile
func (t *Tilemap) ForEach(con func(*Tile)) {
	for _, row := range t.Tiles {
		for _, t := range row {
			con(t)
		}
	}
}
