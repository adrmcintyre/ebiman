package ghost

import (
	"github.com/adrmcintyre/poweraid/geom"
)

// UpdateTarget ensures the ghost seeks the correct target
// based on its current mode and submode.
func (g *Actor) UpdateTarget() {
	switch g.Mode {
	case MODE_RETURNING:
		g.TargetPos = g.HomePos
	case MODE_PLAYING:
		switch g.SubMode {
		case SUBMODE_SCATTER:
			g.TargetPos = g.ScatterPos
		case SUBMODE_CHASE:
			g.TargetPos = g.GetChaseTarget()
		}
	}
}

// GetChaseTarget returns the screen position to target
// if chase behaviour is active.
func (g *Actor) GetChaseTarget() geom.Position {
	pm := g.Pacman
	switch g.Id {
	case PINKY:
		targetPos := pm.Pos.Add(pm.Dir.ScaleUp(4 * 8))
		if pm.Dir.IsUp() {
			targetPos.X -= 4 * 8
		}
		return targetPos
	case INKY:
		return pm.Pos.Add(pm.Dir.ScaleUp(4 * 8)).Add(pm.Pos.Sub(g.Blinky.Pos))
	case CLYDE:
		if g.Pos.TileDistSq(pm.Pos) < 64 {
			return g.ScatterPos
		}
	}

	return pm.Pos
}
