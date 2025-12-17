package geom

type Position struct {
	X, Y int
}

func TilePos(x, y int) Position {
	return Position{x * 8, y * 8}
}

func (p Position) TileX() int {
	return p.X / 8
}

func (p Position) TileY() int {
	return p.Y / 8
}

func (p Position) TileXY() (int, int) {
	return p.X / 8, p.Y / 8
}

func (p Position) TileEq(q Position) bool {
	pX, pY := p.TileXY()
	qX, qY := q.TileXY()
	return pX == qX && pY == qY
}

func (p Position) TileDistSq(q Position) int {
	dx := p.TileX() - q.TileX()
	dy := p.TileY() - q.TileY()
	return dx*dx + dy*dy
}

func (p Position) Add(d Delta) Position {
	return Position{p.X + d.X, p.Y + d.Y}
}

func (p Position) Sub(q Position) Delta {
	return Delta{p.X - q.X, p.Y - q.Y}
}

func (p Position) WrapTunnel() Position {
	// wrap left<->right to account for tunnel
	const wrapWidth = 8 * 28
	return Position{(p.X + wrapWidth) % wrapWidth, p.Y}
}
