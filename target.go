package main

import "github.com/adrmcintyre/poweraid/geom"

func (g *GhostActor) UpdateTarget(pm *PacmanActor, blinky *GhostActor) {
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

func (g *GhostActor) GetChaseTarget(pm *PacmanActor, blinky *GhostActor) geom.Position {
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
