package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/input"
)

func (g *Game) ResetGame() {
	v := &g.Video
	ls := &g.LevelState

	v.ClearTiles()   // zero out tile mem
	v.ClearPalette() // zero out palette mem
	v.ColorMaze()    // set maze palette + top status

	ls.ClearScores()             // reset score
	ls.SetLives(0)               // set lives to 0 + write to bottom status tiles
	ls.BonusState.ClearBonuses() // reset bonuses + write to bottom status tiles

	g.LevelInit(0)  // init level state
	ls.LevelStart() // reset any state relating to a new life

	g.RunningGame = false
	g.PlayerNumber = 0
	g.Action = ActionSplash
}

func (g *Game) RunGame() Return {
	g.Action = ActionRun
	return ThenContinue
}

func (g *Game) StartGame() Return {
	g.Video.ClearTiles() // zero out splash screen cruft
	g.LevelState.DemoMode = false

	return WithAnim(
		(*Game).AnimStartButtonScreen,
		(*Game).StartGameStep2,
	)
}

func (g *Game) StartGameStep2() Return {
	ls := &g.LevelState
	// set starting lives
	ls.SetLives(g.Options.Lives)
	ls.ClearScores()

	g.BeginLevel(0)

	return WithAnim(
		(*Game).AnimReady,
		(*Game).StartGameStep3,
	)
}

func (g *Game) StartGameStep3() Return {
	// sync each player's saved state to be the same
	g.SavePlayerState(0)
	if g.Options.GameMode == GAME_MODE_2P {
		g.SavePlayerState(1)
	}

	g.RunningGame = true

	return ThenContinue
}

func (g *Game) BeginLevel(level int) {
	v := &g.Video
	ls := &g.LevelState

	v.ClearTiles()   // zero out tiles
	v.ClearPalette() // zero out palettes
	v.ColorMaze()    // set maze + top status palettes

	if g.PlayerNumber == 0 {
		v.Write1Up()
	} else {
		v.Write2Up()
	}

	v.DecodeTiles()            // draw out the maze
	ls.DotState.ResetPellets() // mark all pills as uneaten
	ls.DotState.DrawPellets(v) // populate with pills

	g.LevelInit(level)
	ls.LevelStart()
	g.GhostsStart()                             // reset ghosts to starting positions
	g.Pacman.Start(g.LevelConfig.Speeds.Pacman) // reset pacman to starting position
	g.BonusActor.BonusStart()
	g.HideBonusScore()
	g.HideBonus()
}

func (g *Game) UpdateState() Return {
	g.LevelState.UpdateCounter += 1

	if !g.RunningGame {
		return g.StartGame()
	}

	ghostsPulsed := g.GhostsPulse()
	pacmanPulsed := g.PacmanPulse()

	demoMode := g.LevelState.DemoMode

	if !demoMode {
		g.LeaveHome()

		g.PanicStations()

		g.GhostsSteer(ghostsPulsed)
		g.PacmanSteer(pacmanPulsed)
	}

	g.GhostsMove(ghostsPulsed)
	g.PacmanMove(pacmanPulsed)

	g.TimeoutBonus()
	g.TimeoutBonusScore()

	dead := g.PacmanCollide()

	if demoMode {
		return ThenContinue
	}

	if dead {
		return WithAnim(
			(*Game).AnimPacmanDie,
			(*Game).DieStep1,
		)
	}

	return g.SurviveStep1()
}

func (g *Game) DieStep1() Return {
	ls := &g.LevelState

	ls.DotState.SavePellets(&g.Video)

	// death of pacman triggers global dot counter
	ls.PacmanDiedThisLevel = true
	ls.DecrementLives()

	if ls.Lives == 0 {
		return WithAnim(
			(*Game).AnimGameOver,
			(*Game).DieStep2,
		)
	}
	return g.DieStep2()
}

func (g *Game) DieStep2() Return {
	if !g.LoadNextPlayerState() {
		g.ResetGame()
		return ThenStop
	}
	return g.DieStep3()
}

func (g *Game) DieStep3() Return {
	ls := &g.LevelState

	ls.DotState.DrawPellets(&g.Video)
	g.LevelInit(ls.LevelNumber)

	// TODO refactor this spaghetti
	{
		saved := &g.SavedPlayer[g.PlayerNumber]
		// these get clobbered by LevelInit...
		ls.PacmanDiedThisLevel = saved.PacmanDiedThisLevel
		ls.DotsSinceDeathCounter = saved.DotsSinceDeathCounter
		ls.DotsRemaining = saved.DotsRemaining
		ls.DotsEaten = saved.DotsEaten
	}

	ls.LevelStart()
	g.GhostsStart()
	g.PacmanStart()
	return WithAnim(
		(*Game).AnimReady,
		(*Game).SurviveStep1,
	)
}

func (g *Game) SurviveStep1() Return {
	ls := &g.LevelState

	if ls.DotsRemaining == 0 {
		return WithAnim(
			(*Game).AnimEndLevel,
			(*Game).SurviveStep2,
		)
	}

	return ThenContinue
}

func (g *Game) SurviveStep2() Return {
	ls := &g.LevelState
	ls.LevelNumber += 1
	g.BeginLevel(ls.LevelNumber)

	return WithAnim(
		(*Game).AnimReady,
		(*Game).SurviveStep3,
	)
}

func (g *Game) SurviveStep3() Return {
	return ThenStop
}

func (g *Game) PacmanStart() {
	g.Pacman.Start(g.LevelConfig.Speeds.Pacman)
}

func (g *Game) PacmanPulse() bool {
	pulsed := g.Pacman.Pulse()
	if pulsed {
		// TODO not clear if he should stall for a specified number of frames, updates, or pulses
		// let's go with pulses for now
		if g.Pacman.StallTimer > 0 {
			g.Pacman.StallTimer -= 1
			return false
		}
	}
	return pulsed
}

func (g *Game) PacmanSteer(pulsed bool) {
	if pulsed {
		inDir := input.GetJoystickDirection()
		g.Pacman.Steer(&g.Video, inDir)
	}
}

func (g *Game) PacmanMove(pulsed bool) {
	if pulsed {
		g.Pacman.MovePacman(&g.Video)
	}
}

func (g *Game) GhostsStart() {
	for i := range 4 {
		g.Ghosts[i].Start(
			g.LevelConfig.Speeds.Ghost,
			g.Options.MaxGhosts,
			&g.LevelConfig.DotLimits,
		)
	}
}

func (g *Game) GhostsPulse() (pulsed [4]bool) {
	for j := range 4 {
		g.Ghosts[j].Tunnel(g.LevelConfig.Speeds.Tunnel)
		pulsed[j] = g.Pulse(j)
	}
	return pulsed
}

func (g *Game) GhostsSteer(pulsed [4]bool) {
	v := &g.Video
	pacman := &g.Pacman
	blinky := &g.Ghosts[BLINKY]
	speeds := &g.LevelConfig.Speeds
	ai := g.Options.GhostAi

	for j := range 4 {
		if pulsed[j] {
			g.Ghosts[j].Steer(v, pacman, blinky, speeds, ai)
		}
	}
}

func (g *Game) GhostsMove(pulsed [4]bool) {
	for j := range 4 {
		if pulsed[j] {
			g.Ghosts[j].MoveGhost()
		}
	}
}

func (g *Game) PanicStations() {
	ls := &g.LevelState
	maxGhosts := g.Options.MaxGhosts

	if ls.WhiteBlueTimeout != 0 && ls.UpdateCounter >= ls.WhiteBlueTimeout {
		ls.IsFlashing = true
		ls.IsWhite = !ls.IsWhite
		ls.WhiteBlueTimeout += data.WHITE_BLUE_PERIOD
	}
	if ls.BlueTimeout != 0 && ls.UpdateCounter < ls.BlueTimeout {
		// TODO - will need to clear the effect while a ghost is being eaten
		// if blocking delays are removed in the future - see EatGhost().
	}

	revert := ls.BlueTimeout != 0 && (ls.UpdateCounter >= ls.BlueTimeout ||
		ls.GhostsEaten == maxGhosts)

	if revert {
		ls.BlueTimeout = 0
		ls.WhiteBlueTimeout = 0
	}

	g.GhostsRevert(revert)
	g.GhostsSwitchTactics(revert)
	g.PacmanRevert(revert)
}
