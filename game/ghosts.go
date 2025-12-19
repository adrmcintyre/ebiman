package game

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/ghost"
	"github.com/adrmcintyre/poweraid/option"
)

func (g *Game) GhostsStart() {
	for _, gh := range g.Ghosts {
		gh.Start(
			g.LevelConfig.Speeds.Ghost,
			g.Options.MaxGhosts,
			&g.LevelConfig.DotLimits,
		)
	}
}

func (g *Game) GhostsLeaveHome() {

	// blinky always leaves immediately
	blinky := g.Ghosts[ghost.BLINKY]
	if blinky.Mode == ghost.MODE_HOME {
		blinky.SetLeaveState()
	}

	// check remaining ghosts - only one may leave
	for _, gh := range g.Ghosts {
		if gh.Id == ghost.BLINKY {
			continue
		}
		if gh.Mode == ghost.MODE_HOME {
			leave := false
			// A ghost will leave if pacman has been idle for too long
			if g.IsPacmanIdle() {
				g.PacmanResetIdleTimer()
				leave = true
			} else if g.LevelState.PacmanDiedThisLevel {
				if g.LevelState.DotsSinceDeathCounter == gh.AllDotLimit {
					if gh.Id == ghost.CLYDE {
						g.LevelState.PacmanDiedThisLevel = false
						g.LevelState.DotsSinceDeathCounter = 0
					}
					leave = true
				}
			} else {
				leave = gh.DotsAtHomeCounter >= gh.DotLimit
			}

			if leave {
				gh.SetLeaveState()
				break
			}
		}
	}
}

func (g *Game) GhostsSwitchTactics(revert bool) {
	subModes := []ghost.SubMode{
		ghost.SUBMODE_CHASE,
		ghost.SUBMODE_SCATTER,
	}
	for i, frame := range g.LevelConfig.SwitchTactics {
		if g.LevelState.FrameCounter >= frame {
			for _, gh := range g.Ghosts {
				if revert || gh.SubMode != ghost.SUBMODE_SCARED {
					gh.SetSubMode(subModes[i%len(subModes)])
				}
			}
			break
		}
	}
}

func (g *Game) GhostsRevert(revert bool) {
	for _, gh := range g.Ghosts {
		if revert && gh.SubMode == ghost.SUBMODE_SCARED {
			gh.Pcm = g.LevelConfig.Speeds.Ghost
		}
	}
}

func (g *Game) GhostsSteer(pulsed [4]bool) {
	v := &g.Video
	speeds := &g.LevelConfig.Speeds
	ai := g.Options.GhostAi == option.GHOST_AI_ON

	for i, gh := range g.Ghosts {
		if pulsed[i] {
			gh.Steer(v, speeds, ai)
		}
	}
}

func (g *Game) GhostsPulse() (pulsed [4]bool) {
	for i, gh := range g.Ghosts {
		gh.Tunnel(g.LevelConfig.Speeds.Tunnel)
		pulsed[i] = g.GhostPulse(gh)
	}
	return pulsed
}

func (g *Game) GhostPulse(gh *ghost.Actor) bool {
	pcm := &gh.Pcm

	isBlinky := gh.Id == ghost.BLINKY
	isHunting := gh.Mode == ghost.MODE_PLAYING && gh.SubMode != ghost.SUBMODE_SCARED
	isClydeOut := g.Ghosts[ghost.CLYDE].Mode != ghost.MODE_HOME

	if gh.TunnelPcm != 0 {
		pcm = &gh.TunnelPcm
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
	for i, gh := range g.Ghosts {
		if pulsed[i] {
			gh.Move()
		}
	}
}

func (g *Game) GhostConsume(gh *ghost.Actor) {
	ghostScore := &data.GhostScore[g.LevelState.GhostsEaten]
	g.LevelState.IncrementScore(g.PlayerNumber, ghostScore.Score)

	gh.ScoreLook = ghostScore.Look
	gh.Mode = ghost.MODE_RETURNING
	gh.Pcm = data.PCM_MAX

	g.Pacman.Visible = false

	g.ScheduleDelay(data.DISPLAY_GHOST_SCORE_MS)
	g.AddTask(TASK_GHOST_RETURN, int(gh.Id))
}

func (g *Game) GhostReturn(id int) {
	gh := g.Ghosts[id]
	gh.ScoreLook = 0

	g.Pacman.Visible = true

	g.LevelState.GhostsEaten += 1
}

func (g *Game) DrawGhosts() {
	for _, gh := range g.Ghosts {
		gh.Draw(&g.Video, g.LevelState.IsWhite, g.LevelState.FrameCounter&8 > 0)
	}
}
