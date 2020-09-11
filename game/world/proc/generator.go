package proc

import (
	"math/rand"
)

type Generator struct {
	W, H, D int
	Obj     map[int]*GenObj
}

type GenObj struct {
	Percentage, Layer, Priority int
}

type Result [][][]int

func (r *Result) ForEach(con func(int, int)) {
	for y, row := range *r {
		for x := range row {
			con(x, y)
		}
	}
}

func (r *Result) CountSurrounding(x, y, z, id int) int {
	ra := *r
	minX, maxX := x-1, x+2
	minY, maxY := y-1, y+2
	w, h := len(ra[0]), len(ra)
	res := 0
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			if x < 0 || y < 0 || x == w || y == h {
				res += 1
				continue
			}
			if ra[y][x][z] == id {
				res += 1
			}
		}
	}
	return res
}

func (g *Generator) Generate(seed, iters int) Result {
	rd := rand.New(rand.NewSource(int64(seed)))
	res := g.CreateMap()
	w, h := len(res[0])-1, len(res)-1
	for id, o := range g.Obj {
		res.ForEach(func(x, y int) {
			val := res[y][x][o.Layer]
			other, ok := g.Obj[val]

			if (!ok || other.Priority < o.Priority) && ((x == 0 || y == 0 || x == w || y == h) || rd.Intn(100) < o.Percentage) {
				res[y][x][o.Layer] = id
			}
		})
		for i := 0; i < iters; i++ {
			res.ForEach(func(x, y int) {
				if res.CountSurrounding(x, y, o.Layer, id) > 4 {
					res[y][x][o.Layer] = id
				} else {
					res[y][x][o.Layer] = 0
				}
			})
		}
	}
	return res
}

func (g *Generator) CreateMap() Result {
	res := make([][][]int, g.H)
	for y := range res {
		row := make([][]int, g.W)
		for x := range row {
			row[x] = make([]int, g.D)
		}
		res[y] = row
	}
	return res
}
