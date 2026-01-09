package game

import (
	"math"

	"github.com/adrmcintyre/ebiman/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// LayoutStyle specifies a layout for touch controls
type LayoutStyle int

// Enumerate each layout style
const (
	layoutCircles LayoutStyle = iota
	layoutRectsLUDR
	layoutRectsLRUD
	layoutTriangles
	layoutWedges
)

// MakeTouchLayout creates a touch layout in the area bounded by l, t, w, h
func MakeTouchLayout(style LayoutStyle, l, t, w, h int) *input.TouchLayout {
	r := l + w
	b := t + h
	layout := input.NewTouchLayout()

	// whole play area is "space"
	layout.BindRegion(&input.Rect{0, 0, w, t}, ebiten.KeySpace, false)

	switch style {
	case layoutCircles:
		// l       mx        r
		// +-----------------+ t
		// |       ___       |
		// |  ___ ( U ) ___  |
		// | ( L ) --- ( R ) | my
		// | `--- ( D ) ---' |
		// |      `---'      |
		// +-----------------+ b
		radius := 25 // TODO figure this out dynamically
		mx := (l + r) / 2
		my := (t + b) / 2
		ox := int(float64(radius) * math.Sqrt(3))
		layout.BindRegion(&input.Circle{mx - ox, my, radius}, ebiten.KeyLeft, true)
		layout.BindRegion(&input.Circle{mx + ox, my, radius}, ebiten.KeyRight, true)
		layout.BindRegion(&input.Circle{mx, my - radius, radius}, ebiten.KeyUp, true)
		layout.BindRegion(&input.Circle{mx, my + radius, radius}, ebiten.KeyDown, true)

	case layoutRectsLUDR:
		// l    x1      x2   r
		// +----+-------+----+ t
		// |    |       |    |
		// |    |   U   |    |
		// | L  +-------+  R | my
		// |    |   D   |    |
		// |    |       |    |
		// +----+-------+----+ b
		my := (t + b) / 2
		x1 := l + w/4
		x2 := r - w/4
		layout.BindRegion(&input.Rect{l, t, x1, b}, ebiten.KeyLeft, true)
		layout.BindRegion(&input.Rect{x1, t, x2, my}, ebiten.KeyUp, true)
		layout.BindRegion(&input.Rect{x1, my, x2, b}, ebiten.KeyDown, true)
		layout.BindRegion(&input.Rect{x2, t, r, b}, ebiten.KeyRight, true)

	case layoutRectsLRUD:
		// l    x1   x2      r
		// +----+----+-------+ t
		// |    |    |       |
		// |    |    |   U   |
		// | L  | R  +-------+ my
		// |    |    |   D   |
		// |    |    |       |
		// +----+----+-------+ b
		my := (t + b) / 2
		x1 := l + w/4
		x2 := l + w/2
		layout.BindRegion(&input.Rect{l, t, x1, b}, ebiten.KeyLeft, true)
		layout.BindRegion(&input.Rect{x1, t, x2, b}, ebiten.KeyRight, true)
		layout.BindRegion(&input.Rect{x2, t, r, my}, ebiten.KeyUp, true)
		layout.BindRegion(&input.Rect{x2, my, r, b}, ebiten.KeyDown, true)

	case layoutTriangles:
		// l        mx       r
		// +-----------------+ t
		// | `.     U      . |
		// |    `.     . '   |
		// |  L    `x'    R  | my
		// |     . ' ' .     |
		// | . '    D    ' . |
		// +-----------------+ b
		mx := (l + r) / 2
		my := (t + b) / 2
		layout.BindRegion(&input.Triangle{l, t, l, b, mx, my}, ebiten.KeyLeft, true)
		layout.BindRegion(&input.Triangle{r, t, r, b, mx, my}, ebiten.KeyRight, true)
		layout.BindRegion(&input.Triangle{l, t, r, t, mx, my}, ebiten.KeyUp, true)
		layout.BindRegion(&input.Triangle{l, b, r, b, mx, my}, ebiten.KeyDown, true)

	case layoutWedges:
		// l     x1   x2    r
		// +----------------+ t
		// | `.    U      . |
		// |    `+----+ '   | y1
		// |  L  |dead| R   |
		// |     +----+     | y2
		// | . '   D    ' . |
		// +----------------+ b
		x1, x2 := l+w*3/8, l+w*5/8
		y1, y2 := t+h*3/8, t+h*5/8
		layout.BindRegion(&input.Quad{l, t, l, b, x1, y2, x1, y1}, ebiten.KeyLeft, true)
		layout.BindRegion(&input.Quad{r, t, r, b, x2, y2, x2, y1}, ebiten.KeyRight, true)
		layout.BindRegion(&input.Quad{l, t, r, t, x2, y1, x1, y1}, ebiten.KeyUp, true)
		layout.BindRegion(&input.Quad{l, b, r, b, x2, y2, x1, y2}, ebiten.KeyDown, true)
	}

	return layout
}
