package game

import "github.com/adrmcintyre/poweraid/data"

func (g *Game) GhostsStart() {
	for i := range 4 {
		g.Ghosts[i].Start(
			g.LevelConfig.Speeds.Ghost,
			g.Options.MaxGhosts,
			&g.LevelConfig.DotLimits,
		)
	}
}

func (g *Game) GhostsLeaveHome() {

	// blinky always leaves immediately
	blinky := g.Ghosts[BLINKY]
	if blinky.Mode == MODE_HOME {
		blinky.SetLeaveState()
	}

	// check remaining ghosts - only one may leave
	for id := PINKY; id < GhostId(g.Options.MaxGhosts); id++ {
		ghost := &g.Ghosts[id]
		if ghost.Mode == MODE_HOME {
			leave := false
			// A ghost will leave if pacman has been idle for too long
			if g.IsPacmanIdle() {
				g.PacmanResetIdleTimer()
				leave = true
			} else if g.LevelState.PacmanDiedThisLevel {
				if g.LevelState.DotsSinceDeathCounter == ghost.AllDotLimit {
					if id == CLYDE {
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

func (g *Game) GhostsRevert(revert bool) {
	for j := range 4 {
		ghost := &g.Ghosts[j]
		if revert && ghost.SubMode == SUBMODE_SCARED {
			ghost.Pcm = g.LevelConfig.Speeds.Ghost
		}
	}
}

func (g *Game) GhostsSteer(pulsed [4]bool) {
	v := &g.Video
	pacman := &g.Pacman
	blinky := &g.Ghosts[BLINKY]
	speeds := &g.LevelConfig.Speeds
	ai := g.Options.GhostAi == GHOST_AI_ON

	for j := range 4 {
		if pulsed[j] {
			g.Ghosts[j].Steer(v, pacman, blinky, speeds, ai)
		}
	}
}

func (g *Game) GhostsPulse() (pulsed [4]bool) {
	for j := range 4 {
		g.Ghosts[j].Tunnel(g.LevelConfig.Speeds.Tunnel)
		pulsed[j] = g.GhostPulse(j)
	}
	return pulsed
}

func (g *Game) GhostPulse(i int) bool {
	ghost := &g.Ghosts[i]

	pcm := &ghost.Pcm

	isBlinky := ghost.Id == BLINKY
	isHunting := ghost.Mode == MODE_PLAYING && ghost.SubMode != SUBMODE_SCARED
	isClydeOut := g.Ghosts[CLYDE].Mode != MODE_HOME

	if ghost.TunnelPcm != 0 {
		pcm = &ghost.TunnelPcm
	} else if isBlinky && isHunting && isClydeOut {
		if g.LevelState.DotsRemaining <= g.LevelConfig.ElroyPills2 {
			pcm = &g.LevelConfig.Speeds.Elroy2
		} else if g.LevelState.DotsRemaining <= g.LevelConfig.ElroyPills1 {
			pcm = &g.LevelConfig.Speeds.Elroy1
		}
	}

	return pcm.Pulse()
}

func (g *Game) GhostsMove(pulsed [4]bool) {
	for j := range 4 {
		if pulsed[j] {
			g.Ghosts[j].Move()
		}
	}
}

func (g *Game) GhostConsume(ghost *GhostActor) {
	ghostScore := &data.GhostScore[g.LevelState.GhostsEaten]
	g.LevelState.IncrementScore(g.PlayerNumber, ghostScore.Score)

	ghost.ScoreLook = ghostScore.Look
	ghost.Mode = MODE_RETURNING
	ghost.Pcm = data.PCM_MAX

	g.Pacman.Visible = false

	g.ScheduleDelay(data.DISPLAY_GHOST_SCORE_MS)
	g.AddTask(TASK_GHOST_RETURN, int(ghost.Id))
}

func (g *Game) GhostReturn(id int) {
	ghost := &g.Ghosts[id]
	ghost.ScoreLook = 0

	g.Pacman.Visible = true

	g.LevelState.GhostsEaten += 1
}

func (g *Game) DrawGhosts() {
	for j := range 4 {
		g.Ghosts[j].DrawGhost(&g.Video, g.LevelState.IsWhite, g.LevelState.FrameCounter&8 > 0)
	}
}
