package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// LayoutStyle specifies a layout for touch controls
type LayoutStyle int

// Enumerate each layout style
const (
	LayoutCircles LayoutStyle = iota
	LayoutRectsLUDR
	LayoutRectsLRUD
	LayoutTriangles
	LayoutWedges
)

// MakeTouchLayout creates a touch layout in the area bounded by l, t, w, h
func MakeTouchLayout(style LayoutStyle, l, t, w, h int) *TouchLayout {
	r := l + w
	b := t + h
	layout := NewTouchLayout()

	// whole play area is "space"
	layout.BindRegion(&Rect{0, 0, w, t}, ebiten.KeySpace, false)

	switch style {
	case LayoutCircles:
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
		layout.BindRegion(&Circle{mx - ox, my, radius}, ebiten.KeyLeft, true)
		layout.BindRegion(&Circle{mx + ox, my, radius}, ebiten.KeyRight, true)
		layout.BindRegion(&Circle{mx, my - radius, radius}, ebiten.KeyUp, true)
		layout.BindRegion(&Circle{mx, my + radius, radius}, ebiten.KeyDown, true)

	case LayoutRectsLUDR:
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
		layout.BindRegion(&Rect{l, t, x1, b}, ebiten.KeyLeft, true)
		layout.BindRegion(&Rect{x1, t, x2, my}, ebiten.KeyUp, true)
		layout.BindRegion(&Rect{x1, my, x2, b}, ebiten.KeyDown, true)
		layout.BindRegion(&Rect{x2, t, r, b}, ebiten.KeyRight, true)

	case LayoutRectsLRUD:
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
		layout.BindRegion(&Rect{l, t, x1, b}, ebiten.KeyLeft, true)
		layout.BindRegion(&Rect{x1, t, x2, b}, ebiten.KeyRight, true)
		layout.BindRegion(&Rect{x2, t, r, my}, ebiten.KeyUp, true)
		layout.BindRegion(&Rect{x2, my, r, b}, ebiten.KeyDown, true)

	case LayoutTriangles:
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
		layout.BindRegion(&Triangle{l, t, l, b, mx, my}, ebiten.KeyLeft, true)
		layout.BindRegion(&Triangle{r, t, r, b, mx, my}, ebiten.KeyRight, true)
		layout.BindRegion(&Triangle{l, t, r, t, mx, my}, ebiten.KeyUp, true)
		layout.BindRegion(&Triangle{l, b, r, b, mx, my}, ebiten.KeyDown, true)

	case LayoutWedges:
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
		layout.BindRegion(&Quad{l, t, l, b, x1, y2, x1, y1}, ebiten.KeyLeft, true)
		layout.BindRegion(&Quad{r, t, r, b, x2, y2, x2, y1}, ebiten.KeyRight, true)
		layout.BindRegion(&Quad{l, t, r, t, x2, y1, x1, y1}, ebiten.KeyUp, true)
		layout.BindRegion(&Quad{l, b, r, b, x2, y2, x1, y2}, ebiten.KeyDown, true)
	}

	return layout
}
