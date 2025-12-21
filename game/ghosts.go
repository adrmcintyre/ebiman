package game

import (
	"github.com/adrmcintyre/poweraid/audio"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/ghost"
	"github.com/adrmcintyre/poweraid/option"
)

// GhostsStart initialises all of the ghosts at level start / restart.
func (g *Game) GhostsStart() {
	for _, gh := range g.Ghosts {
		gh.Start(
			g.LevelConfig.Speeds.Ghost,
			g.Options.MaxGhosts,
			&g.LevelConfig.DotLimits,
		)
	}
}

// CheckGhostsLeaveHome releases ghosts that meet their
// condition for leaving their home.
func (g *Game) CheckGhostsLeaveHome() {

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

// CheckGhostsSwitchTactics switches the ghosts between
// their scatter and chase behaviours when the necessary
// game triggers are met.
func (g *Game) CheckGhostsSwitchTactics(revert bool) {
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

// CheckGhostsRevert sets each ghost's speed back to normal
// when they are no longer scared.
func (g *Game) CheckGhostsRevert(revert bool) {
	for _, gh := range g.Ghosts {
		if revert && gh.SubMode == ghost.SUBMODE_SCARED {
			gh.Pcm = g.LevelConfig.Speeds.Ghost
		}
	}
}

// GhostsSteer manages the navigation of the maze for
// each ghost currently on the move (have just pulsed).
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

// CheckGhostsReturned cancels the special "returning" audio
// when there are no more returning ghosts.
func (g *Game) CheckGhostsReturned() {
	numReturning := 0
	for _, gh := range g.Ghosts {
		if gh.Mode == ghost.MODE_RETURNING {
			numReturning += 1
		}
	}
	if numReturning == 0 {
		g.Audio.StopBackgroundEffect(audio.EyesReturning)
	}
}

func (g *Game) GhostsTunnel() {
	for _, gh := range g.Ghosts {
		gh.CheckTunnelSpeed(g.LevelConfig.Speeds.Tunnel)
	}
}

// GhostsPulse advances each ghost's pulse train, and reports
// those which pulsed (i.e. are due a movement update).
func (g *Game) GhostsPulse() (pulsed [4]bool) {
	for i, gh := range g.Ghosts {
		pulsed[i] = g.GhostPulse(gh)
	}
	return pulsed
}

// GhostPulse advances the appropriate pulse train for a specific
// ghost, and reports if is pulsed (i.e. is due for a movement update).
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

// GhostsMove advances the position of each ghost that just had a pulse.
func (g *Game) GhostsMove(pulsed [4]bool) {
	for i, gh := range g.Ghosts {
		if pulsed[i] {
			gh.Move()
		}
	}
}

// PacmanEatsGhost is triggered when pacman collides with a vulnerable ghost.
// Pacman vanishes, and the ghost's score value is displayed, during a brief
// pause, and the ghost schedule to be put into "eyes returning" mode.
func (g *Game) PacmanEatsGhost(gh *ghost.Actor) {
	ghostScore := &data.GhostScore[g.LevelState.GhostsEaten]
	g.IncrementScore(ghostScore.Score)

	gh.ScoreLook = ghostScore.Look
	gh.Mode = ghost.MODE_RETURNING
	gh.Pcm = data.PCM_MAX

	g.Pacman.Visible = false

	g.ScheduleDelay(data.DISPLAY_GHOST_SCORE_MS)
	g.AddTask(TASK_GHOST_RETURN, int(gh.Id))
	g.Audio.PlayBackgroundEffect(audio.EyesReturning)
}

// GhostReturn is invoked when the pause to see the ghost's score value has
// expired. The score is hidden and pacman reappears.
func (g *Game) GhostReturn(id int) {
	gh := g.Ghosts[id]
	gh.ScoreLook = 0

	g.Pacman.Visible = true

	g.LevelState.GhostsEaten += 1
}

// DrawGhosts schedules the ghosts to be rendered as sprites in the next frame.
func (g *Game) DrawGhosts() {
	for _, gh := range g.Ghosts {
		gh.Draw(&g.Video, g.LevelState.IsWhite, g.LevelState.FrameCounter&8 > 0)
	}
}
