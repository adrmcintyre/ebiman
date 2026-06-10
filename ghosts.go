package main

import (
	"github.com/adrmcintyre/ebiman/actor"
	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/data"
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

// CheckGhostsLeaveHome releases any ghost that meets its
// condition for leaving home.
func (g *Game) CheckGhostsLeaveHome() {
	for _, gh := range g.Ghosts {
		g.CheckGhostLeaveHome(gh)
	}
}

// CheckGhostLeaveHome releases the specified ghost if its
// conditions for leaving home are met.
func (g *Game) CheckGhostLeaveHome(gh *actor.Ghost) {
	switch {
	case int(gh.Id) >= g.Options.MaxGhosts:
		return

	case gh.Mode != actor.GhostModeHome:
		return

	// blinky never hangs around
	case gh.Id == actor.Blinky:
		gh.SetLeaveState()

	// a ghost will leave if pacman has been idle for too long
	case g.IsPacmanIdle():
		g.PacmanResetIdleTimer()
		gh.SetLeaveState()

	// if pacman has died, refer to the global dot counter
	case g.Player.PacmanDiedThisLevel && g.Player.DotsSinceDeathCounter == gh.AllDotLimit:
		if gh.Id == actor.Clyde {
			g.Player.PacmanDiedThisLevel = false
			g.Player.DotsSinceDeathCounter = 0
		}
		gh.SetLeaveState()

	// otherwise use the ghost's personal dot counter
	case gh.DotsAtHomeCounter >= gh.DotLimit:
		gh.SetLeaveState()
	}
}

// CheckGhostsSwitchTactics switches unscared ghosts
// between their scatter and chase behaviours.
func (g *Game) CheckGhostsSwitchTactics() {
	subMode := actor.GhostSubModeScattering

	tactics := []actor.GhostSubMode{
		actor.GhostSubModeChasing,
		actor.GhostSubModeScattering,
	}
	for i, frame := range g.LevelConfig.SwitchTactics {
		if g.LevelState.FrameCounter >= frame {
			subMode = tactics[i%len(tactics)]
			break
		}
	}

	for _, gh := range g.Ghosts {
		if gh.SubMode != actor.GhostSubModeScared {
			gh.SetSubMode(subMode)
		}
	}
}

// GhostsScare puts all ghosts into the scared state.
func (g *Game) GhostsScare() {
	ls := g.LevelState
	ls.GhostsScaredTimeout = ls.UpdateCounter + g.LevelConfig.ScaredTime
	ls.GhostsFlashTimeout = ls.GhostsScaredTimeout - g.LevelConfig.WhiteBlueCount*data.WhiteBluePeriod
	ls.GhostsAreFlashing = false
	ls.GhostsAreWhite = false
	ls.GhostsEaten = 0

	for _, gh := range g.Ghosts {
		gh.Scare(g.LevelConfig.Speeds.GhostBlue)
	}
}

// GhostsUnscare sets each scared ghost back to normal.
func (g *Game) GhostsUnscare() {
	ls := g.LevelState
	ls.GhostsScaredTimeout = 0
	ls.GhostsFlashTimeout = 0

	for _, gh := range g.Ghosts {
		if gh.SubMode == actor.GhostSubModeScared {
			gh.SubMode = actor.GhostSubModeChasing
			gh.Pcm = g.LevelConfig.Speeds.Ghost
		}
	}
}

// GhostsSteer manages the navigation of the maze for
// each ghost currently on the move (have just pulsed).
func (g *Game) GhostsSteer(pulsed [4]bool) {
	v := g.Video
	speeds := &g.LevelConfig.Speeds
	ai := g.Options.GhostAi == GhostAiOn

	for i, gh := range g.Ghosts {
		if pulsed[i] {
			gh.Steer(v, speeds, ai)
		}
	}
}

// UpdateGhostsReturnAudio cancels the special "returning" audio
// when there are no more returning ghosts.
func (g *Game) UpdateGhostReturnAudio() {
	for _, gh := range g.Ghosts {
		if gh.Mode == actor.GhostModeReturning {
			return
		}
	}
	g.Audio.StopBackgroundEffect(audio.EyesReturning)
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
// ghost, and reports if it pulsed (i.e. is due for a movement update).
func (g *Game) GhostPulse(gh *actor.Ghost) bool {
	if gh.TunnelPcm != 0 {
		return gh.TunnelPcm.Pulse()
	}

	isBlinky := gh.Id == actor.Blinky
	isHunting := gh.Mode == actor.GhostModePlaying && gh.SubMode != actor.GhostSubModeScared
	isClydeOut := g.Ghosts[actor.Clyde].Mode != actor.GhostModeHome

	if isBlinky && isHunting && isClydeOut {
		if g.Player.DotsRemaining <= g.LevelConfig.ElroyPills2 {
			return g.LevelConfig.Speeds.Elroy2.Pulse()
		}
		if g.Player.DotsRemaining <= g.LevelConfig.ElroyPills1 {
			return g.LevelConfig.Speeds.Elroy1.Pulse()
		}
	}

	return gh.Pcm.Pulse()
}

// GhostsMove advances the position of each ghost that just had a pulse.
func (g *Game) GhostsMove(pulsed [4]bool) {
	for i, gh := range g.Ghosts {
		if pulsed[i] {
			isNewTile := gh.Move()
			if !g.DemoMode && g.Options.IsElectric() {
				if isNewTile {
					g.Player.Pills.NetCharge += gh.CheckModifyCharge(g.Video, g.LevelState.FrameCounter, g.LevelConfig.Electric)
				}
			}
		}
	}
}

// PacmanEatsGhost is triggered when pacman collides with a vulnerable ghost.
// Pacman vanishes, and the ghost's score value is displayed, during a brief
// pause, and the ghost schedule to be put into "eyes returning" mode.
func (g *Game) PacmanEatsGhost(gh *actor.Ghost) {
	ghostScore := &data.GhostScore[g.LevelState.GhostsEaten]
	g.IncrementScore(ghostScore.Score)

	gh.SetEaten(ghostScore.Look)
	g.Pacman.Visible = false

	g.ScheduleDelay(data.DisplayGhostScoreMs)
	g.AddTask(TaskGhostReturn, int(gh.Id))
	g.Audio.PlayBackgroundEffect(audio.EyesReturning)
}

// GhostReturn is invoked when the pause to see the ghost's score value has
// expired. The score is hidden and pacman reappears.
func (g *Game) GhostReturn(id int) {
	gh := g.Ghosts[id]
	gh.HideScore()

	g.Pacman.Visible = true

	g.LevelState.GhostsEaten += 1
}

func (g *Game) NotifyGhostsPillEaten() {
	for _, gh := range g.Ghosts {
		if gh.Id == actor.Blinky {
			continue
		}
		if gh.Mode == actor.GhostModeHome {
			gh.DotsAtHomeCounter += 1
			break
		}
	}
}

// DrawGhosts schedules the ghosts to be rendered as sprites in the next frame.
func (g *Game) DrawGhosts() {
	for _, gh := range g.Ghosts {
		gh.Draw(g.Video, g.LevelState.GhostsAreWhite, g.LevelState.FrameCounter&8 > 0)
	}
}
