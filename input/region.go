package input

import (
	image_color "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// define the stroke for rendering regions
const (
	antialias   = true
	strokeWidth = 1.0
)

// outlineColor defines the colour of a region's outline
var outlineColor = image_color.Gray{0x80}

// A Region is a shape which can report if a point is in its interior.
// Used for detecting touches.
type Region interface {
	Contains(x, y int) bool
	Draw(img *ebiten.Image)
	Centre() (int, int)
}

// A Rect is a rectangular Region.
type Rect struct {
	X0, Y0 int
	X1, Y1 int
}

var _ Region = (*Rect)(nil)

// Contains returns true if (x,y) lies inside the region.
func (rect Rect) Contains(x int, y int) bool {
	hIn := rect.X0 <= x && x <= rect.X1
	vIn := rect.Y0 <= y && y <= rect.Y1
	return hIn && vIn
}

// Draw paints the region onto img.
func (rect Rect) Draw(img *ebiten.Image) {
	var (
		x0 = float32(rect.X0)
		y0 = float32(rect.Y0)
		x1 = float32(rect.X1)
		y1 = float32(rect.Y1)
	)
	vector.StrokeRect(img, x0, y0, x1-x0, y1-y0, strokeWidth, outlineColor, antialias)
}

// Centre returns the region's centre point.
func (rect Rect) Centre() (int, int) {
	return (rect.X0 + rect.X1) / 2, (rect.Y0 + rect.Y1) / 2
}

// A Circle is a circular Region.
type Circle struct {
	X, Y, R int
}

var _ Region = (*Circle)(nil)

// Contains returns true if (x,y) lies inside the region.
func (c Circle) Contains(x int, y int) bool {
	dx, dy := c.X-x, c.Y-y
	return dx*dx+dy*dy <= c.R*c.R
}

// Centre returns the region's centre point.
func (c Circle) Centre() (int, int) {
	return c.X, c.Y
}

// Draw paints the region onto img.
func (c Circle) Draw(img *ebiten.Image) {
	vector.StrokeCircle(img, float32(c.X), float32(c.Y), float32(c.R), strokeWidth, outlineColor, antialias)
}

// A Triangle is a triangular region.
type Triangle struct {
	X1, Y1 int
	X2, Y2 int
	X3, Y3 int
}

// sign is a determines which side of line (x1,y1)-(x2,y2) the point (px,py) lies.
func sign(px, py, x1, y1, x2, y2 int) int {
	return (px-x2)*(y1-y2) - (x1-x2)*(py-y2)
}

// Contains returns true if (x,y) lies inside the region.
func (t Triangle) Contains(x, y int) bool {
	d12 := sign(x, y, t.X1, t.Y1, t.X2, t.Y2)
	d23 := sign(x, y, t.X2, t.Y2, t.X3, t.Y3)
	d31 := sign(x, y, t.X3, t.Y3, t.X1, t.Y1)

	hasNeg := (d12 < 0) || (d23 < 0) || (d31 < 0)
	hasPos := (d12 > 0) || (d23 > 0) || (d31 > 0)

	return !(hasNeg && hasPos)
}

// Centre returns the region's centre point.
func (t Triangle) Centre() (int, int) {
	return (t.X1 + t.X2 + t.X3) / 3, (t.Y1 + t.Y2 + t.Y3) / 3
}

// Draw paints the region onto img.
func (t Triangle) Draw(img *ebiten.Image) {
	var (
		x1 = float32(t.X1)
		y1 = float32(t.Y1)
		x2 = float32(t.X2)
		y2 = float32(t.Y2)
		x3 = float32(t.X3)
		y3 = float32(t.Y3)
	)
	vector.StrokeLine(img, x1, y1, x2, y2, strokeWidth, outlineColor, antialias)
	vector.StrokeLine(img, x2, y2, x3, y3, strokeWidth, outlineColor, antialias)
	vector.StrokeLine(img, x3, y3, x1, y1, strokeWidth, outlineColor, antialias)
}

// A Quad is a convex quadrilateral.
type Quad struct {
	X1, Y1 int
	X2, Y2 int
	X3, Y3 int
	X4, Y4 int
}

// Contains returns true if (x,y) lies inside the region.
func (q Quad) Contains(x, y int) bool {
	// Divide the quad into two triangles and test each
	// for containment.
	//
	//    x1     x2
	// y1 .-------. y2
	//     \ p  .' \
	//      \ .'    \
	//    y4 '-------' y3
	//       x4    x3
	d12 := sign(x, y, q.X1, q.Y1, q.X2, q.Y2)
	d23 := sign(x, y, q.X2, q.Y2, q.X3, q.Y3)
	d31 := sign(x, y, q.X3, q.Y3, q.X1, q.Y1)
	d34 := sign(x, y, q.X3, q.Y3, q.X4, q.Y4)
	d41 := sign(x, y, q.X4, q.Y4, q.X1, q.Y1)
	d13 := -d31

	hasNeg1 := (d12 < 0) || (d23 < 0) || (d31 < 0)
	hasPos1 := (d12 > 0) || (d23 > 0) || (d31 > 0)
	hasNeg2 := (d34 < 0) || (d41 < 0) || (d13 < 0)
	hasPos2 := (d34 > 0) || (d41 > 0) || (d13 > 0)

	return !(hasNeg1 && hasPos1 && hasNeg2 && hasPos2)
}

// Centre returns the region's centre point.
func (q Quad) Centre() (int, int) {
	return (q.X1 + q.X2 + q.X3 + q.X4) / 4, (q.Y1 + q.Y2 + q.Y3 + q.Y4) / 4
}

// Draw paints the region onto img.
func (q Quad) Draw(img *ebiten.Image) {
	var (
		x1 = float32(q.X1)
		y1 = float32(q.Y1)
		x2 = float32(q.X2)
		y2 = float32(q.Y2)
		x3 = float32(q.X3)
		y3 = float32(q.Y3)
		x4 = float32(q.X4)
		y4 = float32(q.Y4)
	)
	vector.StrokeLine(img, x1, y1, x2, y2, strokeWidth, outlineColor, antialias)
	vector.StrokeLine(img, x2, y2, x3, y3, strokeWidth, outlineColor, antialias)
	vector.StrokeLine(img, x3, y3, x4, y4, strokeWidth, outlineColor, antialias)
	vector.StrokeLine(img, x4, y4, x1, y1, strokeWidth, outlineColor, antialias)
}

// An invisibleRegion adapts a region to make it invisible.
type invisibleRegion struct {
	Region
}

var _ Region = (*invisibleRegion)(nil)

// Draw is a no-op for invisible regions.
func (ir invisibleRegion) Draw(img *ebiten.Image) {
}

// Centre returns (0,0) for invisible regions.
func (it invisibleRegion) Centre() (int, int) {
	return 0, 0
}
