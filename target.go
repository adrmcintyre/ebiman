package main

import "github.com/adrmcintyre/poweraid/video"

func (g *GhostActor) UpdateTarget(pacman *PacmanActor, blinky *GhostActor) {
	if g.Mode == MODE_RETURNING {
		g.TargetPos = g.HomePos.ToTilePos()
		return
	}

	if g.Mode != MODE_PLAYING {
		return
	}

	if g.SubMode == SUBMODE_SCATTER {
		switch g.Id {
		case BLINKY, PINKY, INKY, CLYDE:
			g.TargetPos = g.ScatterPos
		}
	} else if g.SubMode == SUBMODE_CHASE {
		pm := pacman.Motion
		pmPos := pm.Pos.ToTilePos()

		switch g.Id {
		case BLINKY:
			g.TargetPos = pmPos
		case PINKY:
			g.TargetPos = video.TilePos{
				pmPos.X + 4*pm.Vel.Vx,
				pmPos.Y + 4*pm.Vel.Vy,
			}
			if pm.Vel.Vy < 0 {
				g.TargetPos.X -= 4
			}
		case INKY:
			blinkyPos := blinky.Motion.Pos.ToTilePos()
			g.TargetPos = video.TilePos{
				2*(pmPos.X+2*pm.Vel.Vx) - blinkyPos.X,
				2*(pmPos.Y+2*pm.Vel.Vy) - blinkyPos.Y,
			}
		case CLYDE:
			ghostPos := g.Motion.Pos.ToTilePos()
			if ghostPos.DistSq(pmPos) < 64 {
				g.TargetPos = g.ScatterPos
			} else {
				g.TargetPos = pmPos
			}
		}
	}
}
