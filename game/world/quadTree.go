package world

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Collidable interface {
	GetRect() pixel.Rect
	IsDead() bool
}

type Common struct {
	Depth int
	Level int
	Cap   int //max amount of objects per quadrant, if there is more quadrant splits
}

type Quadtree struct {
	pixel.Rect
	tl, tr, bl, br, pr *Quadtree
	Shapes             []Collidable
	Common
	splitted bool
}

// Creates new quad tree reference.
// bounds - defines position of quad tree and its size. If shapes goes out of bounds they
// will not be assigned to quadrants and the tree will be ineffective.
// depth - resolution of quad tree. It lavais splits in half so if bounds size is 100 x 100
// and depth is 2 smallest quadrants will be 25 x 25. Making resolution too high is redundant
// if shapes cannot fit into smallest quadrants.
// cap - sets maximal capacity of quadrant before it splits to 4 smaller. Making can too big is
// inefficient. optimal value can be 10 but its allways better to test what works the best.
func NewQuadTree(bounds pixel.Rect, depth, cap int) *Quadtree {
	return &Quadtree{
		Rect: bounds,
		Common: Common{
			Depth: depth,
			Cap:   cap,
		},
	}
}

// generates subquadrants, always check if quadrant is not already splitted
func (q *Quadtree) split() {
	q.splitted = true
	newCommon := q.Common
	newCommon.Level++
	halfH := q.H() / 2
	halfW := q.W() / 2
	center := q.Center()
	q.tl = &Quadtree{
		Rect: pixel.Rect{
			Min: pixel.V(q.Min.X, q.Min.Y+halfH),
			Max: pixel.V(q.Max.X-halfW, q.Max.Y),
		},
		pr:     q,
		Common: newCommon,
	}
	q.tr = &Quadtree{
		Rect: pixel.Rect{
			Min: center,
			Max: q.Max,
		},
		pr:     q,
		Common: newCommon,
	}
	q.bl = &Quadtree{
		Rect: pixel.Rect{
			Min: q.Min,
			Max: center,
		},
		pr:     q,
		Common: newCommon,
	}
	q.br = &Quadtree{
		Rect: pixel.Rect{
			Min: pixel.V(q.Min.X+halfW, q.Min.Y),
			Max: pixel.V(q.Max.X, q.Min.Y+halfH),
		},
		pr:     q,
		Common: newCommon,
	}
}

func (q *Quadtree) fits(rect pixel.Rect) bool {
	return rect.Max.X > q.Min.X && rect.Max.X < q.Max.X && rect.Min.Y > q.Min.Y && rect.Max.Y < q.Max.Y
}

// finds out in witch subquadrant the shape belongs to. Shape has to overlap only with one quadrant,
// otherwise it returns nil
func (q *Quadtree) getSub(rect pixel.Rect) *Quadtree {
	vertical := q.Min.X + q.W()/2
	horizontal := q.Min.Y + q.H()/2

	if !q.fits(rect) {
		return nil
	}

	left := rect.Max.X <= vertical
	right := rect.Min.X >= vertical
	if rect.Min.Y >= horizontal {
		if left {
			return q.tl
		} else if right {
			return q.tr
		}
	} else if rect.Max.Y <= horizontal {
		if left {
			return q.bl
		} else if right {
			return q.br
		}
	}
	return nil
}

// Adds the shape to quad tree and asians it to correct quadrant.
// Proper way is adding all shapes first and then detecting collisions.
// For struct to implement Collidable interface it has to implement
// GetRect() *pixel.Rect. GetRect function also slightly affects performance.
func (q *Quadtree) Insert(collidable Collidable) {
	rect := collidable.GetRect()

	if q.splitted {
		fitting := q.getSub(rect)
		if fitting != nil {
			fitting.Insert(collidable)
			return
		}
		q.Shapes = append(q.Shapes, collidable)
		return
	}
	q.Shapes = append(q.Shapes, collidable)
	if q.Cap <= len(q.Shapes) && q.Level != q.Depth {
		q.split()
		new := []Collidable{}
		for _, s := range q.Shapes {
			fitting := q.getSub(s.GetRect())
			if fitting != nil {
				fitting.Insert(s)
			} else {
				new = append(new, s)
			}
		}
		q.Shapes = new
	}
}

// Update reassigns shapes to correct quadrants if needed
func (q *Quadtree) Update() {
	for _, c := range q.update() {
		q.Insert(c)
	}
}

func (q *Quadtree) update() []Collidable {
	moved := []Collidable{}
	i := 0
	for _, s := range q.Shapes {
		if s.IsDead() {
			continue
		}
		if !q.fits(s.GetRect()) {
			moved = append(moved, s)
			continue
		}

		q.Shapes[i] = s
		i++
	}

	for j := i; j < len(q.Shapes); j++ {
		q.Shapes[j] = nil
	}

	q.Shapes = q.Shapes[:i]

	if !q.splitted {
		return moved
	}

	moved = append(moved, q.tl.update()...)
	moved = append(moved, q.tr.update()...)
	moved = append(moved, q.bl.update()...)
	moved = append(moved, q.br.update()...)
	return moved
}

func (q *Quadtree) Remove(c Collidable) {
	for i, s := range q.Shapes {
		if s == c {
			end := len(q.Shapes) - 1
			q.Shapes[i] = q.Shapes[end]
			q.Shapes[end] = nil
			q.Shapes = q.Shapes[:end]
			return
		}
	}

	if !q.splitted {
		return
	}
	q.tr.Remove(c)
	q.tl.Remove(c)
	q.bl.Remove(c)
	q.br.Remove(c)
	/*rect := c.GetRect()
	var current *Quadtree
	sub := q
	for sub != nil {
		current = sub
		sub = current.getSub(rect)
	}
	for current.pr != nil {

		for i, s := range current.Shapes {
			if s == c {
				end := len(current.Shapes) - 1
				current.Shapes[i] = current.Shapes[end]
				current.Shapes[end] = nil
				current.Shapes = current.Shapes[:end]
				return nil
			}
		}
		current = current.pr
	}
	return errors.New("hell no")*/
}

// returns all coliding collidables, if rect belongs to object that is already
// inserted in tree it returns is as well
func (q *Quadtree) GetColliding(rect pixel.Rect) []Collidable {
	colliding := []Collidable{}

	for _, c := range q.Shapes {
		if c.GetRect().Intersects(rect) {
			colliding = append(colliding, c)
		}
	}

	if q.splitted {
		if q.tl.Intersects(rect) {
			colliding = append(colliding, q.tl.GetColliding(rect)...)
		}
		if q.tr.Intersects(rect) {
			colliding = append(colliding, q.tr.GetColliding(rect)...)
		}
		if q.bl.Intersects(rect) {
			colliding = append(colliding, q.bl.GetColliding(rect)...)
		}
		if q.br.Intersects(rect) {
			colliding = append(colliding, q.br.GetColliding(rect)...)
		}
	}

	return colliding
}

func (q *Quadtree) GetShapeCount() int {
	res := len(q.Shapes)
	if q.splitted {
		res += q.tl.GetShapeCount() + q.tr.GetShapeCount() + q.bl.GetShapeCount() + q.br.GetShapeCount()
	}
	return res
}

func (q *Quadtree) GetSmallestQuad(rect pixel.Rect) *Quadtree {
	current := q
	for {
		sub := current.getSub(rect)
		if sub == nil {
			break
		}
		current = sub
	}
	return current
}

// Resets the tree, use this every frame before inserting all shapes
// other wise you will run out of memory eventually and tree will not even work properly
func (q *Quadtree) Clear() {
	q.Shapes = []Collidable{}
	q.tl, q.tr, q.bl, q.br = nil, nil, nil, nil
	q.splitted = false
}

// visualizes state of quadtree
func (q *Quadtree) Draw(id *imdraw.IMDraw, thickness float64) {
	id.Push(q.Min)
	id.Push(q.Max)
	id.Rectangle(thickness)
	if !q.splitted {
		return
	}
	q.tl.Draw(id, thickness)
	q.tr.Draw(id, thickness)
	q.bl.Draw(id, thickness)
	q.br.Draw(id, thickness)
}
