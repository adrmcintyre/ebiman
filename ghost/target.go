package ghost

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/pacman"
)

func (g *Actor) UpdateTarget(pm *pacman.Actor, blinky *Actor) {
	switch g.Mode {
	case MODE_RETURNING:
		g.TargetPos = g.HomePos
	case MODE_PLAYING:
		switch g.SubMode {
		case SUBMODE_SCATTER:
			g.TargetPos = g.ScatterPos
		case SUBMODE_CHASE:
			g.TargetPos = g.GetChaseTarget(pm, blinky)
		}
	}
}

func (g *Actor) GetChaseTarget(pm *pacman.Actor, blinky *Actor) geom.Position {
	switch g.Id {
	case PINKY:
		targetPos := pm.Pos.Add(pm.Dir.Scale(4 * 8))
		if pm.Dir.IsUp() {
			targetPos.X -= 4 * 8
		}
		return targetPos
	case INKY:
		return pm.Pos.Add(pm.Dir.Scale(4 * 8)).Add(pm.Pos.Sub(blinky.Pos))
	case CLYDE:
		if g.Pos.TileDistSq(pm.Pos) < 64 {
			return g.ScatterPos
		}
	}

	return pm.Pos
}
