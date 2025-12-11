package main

func (g *GhostActor) UpdateTarget(pacman *PacmanActor, blinky *GhostActor) {
	if g.Mode == MODE_RETURNING {
		g.TargetX = g.HomeX / 8
		g.TargetY = g.HomeY / 8
		return
	}

	if g.Mode != MODE_PLAYING {
		return
	}

	if g.SubMode == SUBMODE_SCATTER {
		switch g.Id {
		case BLINKY, PINKY, INKY, CLYDE:
			g.TargetX = g.ScatterX
			g.TargetY = g.ScatterY
		}
	} else if g.SubMode == SUBMODE_CHASE {
		pm := pacman.Motion
		switch g.Id {
		case BLINKY:
			g.TargetX = pm.X / 8
			g.TargetY = pm.Y / 8
		case PINKY:
			g.TargetX = pm.X/8 + 4*pm.Vx
			g.TargetY = pm.Y/8 + 4*pm.Vy
			if pm.Vy < 0 {
				g.TargetX -= 4
			}
		case INKY:
			g.TargetX = 2*(pm.X/8+2*pm.Vx) - blinky.Motion.X/8
			g.TargetY = 2*(pm.Y/8+2*pm.Vy) - blinky.Motion.Y/8
		case CLYDE:
			dx := g.Motion.X/8 - pm.X/8
			dy := g.Motion.Y/8 - pm.Y/8
			if d2 := dx*dx + dy*dy; d2 < 64 {
				g.TargetX = g.ScatterX
				g.TargetY = g.ScatterY
			} else {
				g.TargetX = pm.X / 8
				g.TargetY = pm.Y / 8
			}
		}
	}
}
