package ghost

import (
	"github.com/adrmcintyre/ebiman/geom"
)

// UpdateTarget ensures the ghost seeks the correct target
// based on its current mode and submode.
func (g *Actor) UpdateTarget() {
	switch g.Mode {
	case ModeReturning:
		g.TargetPos = g.HomePos
	case ModePlaying:
		switch g.SubMode {
		case SubModeScatter:
			g.TargetPos = g.ScatterPos
		case SubModeChase:
			g.TargetPos = g.GetChaseTarget()
		}
	}
}

// GetChaseTarget returns the screen position to target
// if chase behaviour is active.
func (g *Actor) GetChaseTarget() geom.Position {
	pm := g.Pacman
	switch g.Id {
	case Pinky:
		targetPos := pm.Pos.Add(pm.Dir.ScaleUp(4 * 8))
		if pm.Dir.IsUp() {
			targetPos.X -= 4 * 8
		}
		return targetPos
	case Inky:
		return pm.Pos.Add(pm.Dir.ScaleUp(4 * 8)).Add(pm.Pos.Sub(g.Blinky.Pos))
	case Clyde:
		if g.Pos.TileDistSq(pm.Pos) < 64 {
			return g.ScatterPos
		}
	}

	return pm.Pos
}
