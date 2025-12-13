package main

func (g *GhostActor) UpdateTarget(pacman *PacmanActor, blinky *GhostActor) {
	if g.Mode == MODE_RETURNING {
		g.TargetPos = Position{g.HomePos.X / 8, g.HomePos.Y / 8}
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
		switch g.Id {
		case BLINKY:
			g.TargetPos = Position{pm.Pos.X / 8, pm.Pos.Y / 8}
		case PINKY:
			g.TargetPos = Position{
				pm.Pos.X/8 + 4*pm.Vel.Vx,
				pm.Pos.Y/8 + 4*pm.Vel.Vy,
			}
			if pm.Vel.Vy < 0 {
				g.TargetPos.X -= 4
			}
		case INKY:
			g.TargetPos = Position{
				2*(pm.Pos.X/8+2*pm.Vel.Vx) - blinky.Motion.Pos.X/8,
				2*(pm.Pos.Y/8+2*pm.Vel.Vy) - blinky.Motion.Pos.Y/8,
			}
		case CLYDE:
			dx := g.Motion.Pos.X/8 - pm.Pos.X/8
			dy := g.Motion.Pos.Y/8 - pm.Pos.Y/8
			if d2 := dx*dx + dy*dy; d2 < 64 {
				g.TargetPos = g.ScatterPos
			} else {
				g.TargetPos = Position{pm.Pos.X, pm.Pos.Y}
			}
		}
	}
}
