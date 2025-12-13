package main

func (g *GhostActor) GhostSetLeaveState() {
	g.Mode = MODE_LEAVING
}

func (g *Game) GhostsLeave() {
	blinky := g.Ghosts[BLINKY]
	if blinky.Mode == MODE_HOME {
		blinky.GhostSetLeaveState()
	}

	// check remaining ghosts - only one may leave
	for i := 1; i < g.Options.MaxGhosts; i++ {
		ghost := &g.Ghosts[i]
		if ghost.Mode == MODE_HOME {
			leave := false
			if g.IsPacmanIdle() {
				g.PacmanIsActive()
				leave = true
			} else if g.LevelState.GlobalDotCounterEnabled {
				if g.LevelState.GlobalDotCounter == ghost.GlobalDotLimit {
					if ghost.Id == CLYDE {
						g.LevelState.GlobalDotCounterEnabled = false
						g.LevelState.GlobalDotCounter = 0
					}
					leave = true
				}
			} else {
				leave = ghost.DotCounter >= ghost.DotLimit
			}

			if leave {
				ghost.GhostSetLeaveState()
				break
			}
		}
	}
}

func (g *GhostActor) GhostSetSubmode(subMode SubMode) {
	// Ghosts are forced to reverse direction by the system anytime the mode
	// changes from: chase-to-scatter, chase-to-frightened, scatter-to-chase,
	// and scatter-to-frightened.
	// Ghosts do not reverse direction when changing back from frightened to
	// chase or scatter modes.
	switch g.SubMode {
	case subMode:
		return

	case SUBMODE_CHASE:
		if subMode == SUBMODE_SCARED || subMode == SUBMODE_SCATTER {
			g.ReversePending = true
		}

	case SUBMODE_SCATTER:
		if subMode == SUBMODE_SCARED || subMode == SUBMODE_CHASE {
			g.ReversePending = true
		}
	}
	g.SubMode = subMode
}

func (g *Game) GhostsRevert(timeout bool) {
	for j := range 4 {
		ghost := &g.Ghosts[j]
		if ghost.SubMode == SUBMODE_SCARED {
			if !timeout {
				continue
			}
			ghost.Motion.Pcm = g.LevelConfig.Speeds.Ghost
		}

		for i := 6; i >= 0; i-- {
			if g.LevelState.FrameCounter >= g.LevelConfig.ScatterChase[i] {
				if i&1 == 0 {
					ghost.GhostSetSubmode(SUBMODE_CHASE)
				} else {
					ghost.GhostSetSubmode(SUBMODE_SCATTER)
				}
				break
			}
		}
	}
}

func (g *Game) PacmanRevert(timeout bool) {
	if timeout {
		g.Pacman.Motion.Pcm = g.LevelConfig.Speeds.Pacman
	}
}
