package main

import "github.com/adrmcintyre/poweraid/video"

func (g *GhostActor) UpdateTarget(pm *PacmanActor, blinky *GhostActor) {
	switch g.Mode {
	case MODE_RETURNING:
		g.TargetPos = g.HomePos.ToTilePos()
	case MODE_PLAYING:
		switch g.SubMode {
		case SUBMODE_SCATTER:
			g.TargetPos = g.ScatterPos
		case SUBMODE_CHASE:
			g.TargetPos = g.GetChaseTarget(pm, blinky)
		}
	}
}

func (g *GhostActor) GetChaseTarget(pm *PacmanActor, blinky *GhostActor) video.TilePos {
	pmPos := pm.Pos.ToTilePos()

	switch g.Id {
	case PINKY:
		targetPos := video.TilePos{
			pmPos.X + 4*pm.Vel.Vx,
			pmPos.Y + 4*pm.Vel.Vy,
		}
		if pm.Vel.Vy < 0 {
			targetPos.X -= 4
		}
		return targetPos
	case INKY:
		blinkyPos := blinky.Pos.ToTilePos()
		return video.TilePos{
			2*(pmPos.X+2*pm.Vel.Vx) - blinkyPos.X,
			2*(pmPos.Y+2*pm.Vel.Vy) - blinkyPos.Y,
		}
	case CLYDE:
		if g.Pos.ToTilePos().DistSq(pmPos) < 64 {
			return g.ScatterPos
		}
	}

	return pmPos
}
