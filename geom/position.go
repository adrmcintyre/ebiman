package geom

// A Position represents the position of a pixel on the simulated display.
type Position struct {
	X, Y int
}

// TilePos constructs a new position corresponding to the given tile location.
func TilePos(x, y int) Position {
	return Position{x * 8, y * 8}
}

// TileX returns the x-component of the tile-coordinate for this position.
func (p Position) TileX() int {
	return p.X / 8
}

// TileX returns the y-component of the tile-coordinate for this position.
func (p Position) TileY() int {
	return p.Y / 8
}

// TileXY returns the x- and y- components of the tile-coordinates for this position
func (p Position) TileXY() (int, int) {
	return p.X / 8, p.Y / 8
}

// TileEq returns true if p and q are positions within co-incident tiles.
func (p Position) TileEq(q Position) bool {
	pX, pY := p.TileXY()
	qX, qY := q.TileXY()
	return pX == qX && pY == qY
}

// TileDistSq returns the square of the distance between p and q when considered as tiles.
func (p Position) TileDistSq(q Position) int {
	dx := p.TileX() - q.TileX()
	dy := p.TileY() - q.TileY()
	return dx*dx + dy*dy
}

// Add returns p offset by d as a Position.
func (p Position) Add(d Delta) Position {
	return Position{p.X + d.X, p.Y + d.Y}
}

// Sub returns q subtracted from p as a Delta.
func (p Position) Sub(q Position) Delta {
	return Delta{p.X - q.X, p.Y - q.Y}
}

// WrapTunnel returns a new position with the possibility invalid x-coordinate
// of p fixed by wrapping left<->right.
func (p Position) WrapTunnel() Position {
	// wrap left<->right to account for tunnel
	const wrapWidth = 8 * 28
	return Position{(p.X + wrapWidth) % wrapWidth, p.Y}
}
