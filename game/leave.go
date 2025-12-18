package game

func (g *Game) GhostsLeaveHome() {
	// blinky always leaves immediately
	blinky := g.Ghosts[BLINKY]
	if blinky.Mode == MODE_HOME {
		blinky.SetLeaveState()
	}

	// check remaining ghosts - only one may leave
	for i := 1; i < g.Options.MaxGhosts; i++ {
		ghost := &g.Ghosts[i]
		if ghost.Mode == MODE_HOME {
			leave := false
			// A ghost will leave if pacman has been idle for too long
			if g.IsPacmanIdle() {
				g.PacmanResetIdleTimer()
				leave = true
			} else if g.LevelState.PacmanDiedThisLevel {
				if g.LevelState.DotsSinceDeathCounter == ghost.AllDotLimit {
					if i == CLYDE {
						g.LevelState.PacmanDiedThisLevel = false
						g.LevelState.DotsSinceDeathCounter = 0
					}
					leave = true
				}
			} else {
				leave = ghost.DotsAtHomeCounter >= ghost.DotLimit
			}

			if leave {
				ghost.SetLeaveState()
				break
			}
		}
	}
}

func (g *Game) GhostsRevert(revert bool) {
	for j := range 4 {
		ghost := &g.Ghosts[j]
		if revert && ghost.SubMode == SUBMODE_SCARED {
			ghost.Pcm = g.LevelConfig.Speeds.Ghost
		}
	}
}

func (g *Game) GhostsSwitchTactics(revert bool) {
	subModes := []SubMode{
		SUBMODE_CHASE,
		SUBMODE_SCATTER,
	}
	for i, frame := range g.LevelConfig.SwitchTactics {
		if g.LevelState.FrameCounter >= frame {
			for j := range 4 {
				ghost := &g.Ghosts[j]
				if revert || ghost.SubMode != SUBMODE_SCARED {
					ghost.SetSubMode(subModes[i%len(subModes)])
				}
			}
			break
		}
	}
}

func (g *Game) PacmanRevert(revert bool) {
	if revert {
		g.Pacman.Pcm = g.LevelConfig.Speeds.Pacman
	}
}
