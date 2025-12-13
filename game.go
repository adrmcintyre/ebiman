package main

import (
	"github.com/adrmcintyre/poweraid/data"
)

func (g *Game) ResetGame() {
	v := &g.Video
	ls := &g.LevelState

	g.RunningGame = false
	g.PlayerNumber = 0

	v.ClearTiles()                        // zero out tile mem
	v.ClearPalette()                      // zero out palette mem
	ls.ClearScores(v, g.Options.GameMode) // reset score + write to top status tiles
	ls.BonusState.ClearBonuses(v)         // reset bonuses + write to bottom status tiles
	ls.SetLives(v, 0)                     // set lives to 0 + write to bottom status tiles
	g.WritePlayerUp(v)
	ls.WriteHighscore(v)

	v.ColorMaze(0) // set maze palette + top status

	g.LevelInit(0)  // init level state
	ls.LevelStart() // reset any state relating to a new life

	g.Action = ActionSplash
}

func (g *Game) RunGame() Return {
	g.Action = ActionRun
	return ThenContinue
}

func (g *Game) StartGame() Return {
	v := &g.Video
	ls := &g.LevelState

	ls.DemoMode = false

	v.ClearTiles() // zero out splash screen cruft

	return WithAnim(
		(*Game).AnimStartButtonScreen,
		(*Game).StartGameStep2,
	)
}

func (g *Game) StartGameStep2() Return {
	v := &g.Video
	ls := &g.LevelState
	// set starting lives and update bottom status
	ls.SetLives(v, g.Options.Lives)
	ls.ClearScores(v, g.Options.GameMode)

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

	v.ClearTiles()                       // zero out tiles
	v.ClearPalette()                     // zero out palettes
	g.LevelState.DotState.ResetPellets() // mark all pills as uneaten
	v.ColorMaze(0)                       // set maze + top status palettes

	if g.PlayerNumber == 0 {
		v.Write1Up()
	} else {
		v.Write2Up()
	}

	ls.WriteHighscore(v)
	ls.WriteScores(v, g.Options.GameMode)

	v.DecodeTiles()              // draw out the maze
	ls.DotState.DecodePellets(v) // populate with pills
	g.LevelInit(level)
	ls.LevelStart()
	g.GhostsStart()                             // reset ghosts to starting positions
	g.Pacman.Start(g.LevelConfig.Speeds.Pacman) // reset pacman to starting position
	g.BonusActor.BonusStart()
	g.HideBonusScore()
	g.HideBonus()
	g.LevelState.BonusState.WriteBonuses(&g.Video)
	g.Video.WriteLives(g.LevelState.Lives)
}

func (g *Game) UpdateState() Return {
	g.LevelState.UpdateCounter += 1

	if !g.RunningGame {
		return g.StartGame()
	}

	var ghostPulsed [4]bool
	for j := range 4 {
		g.Ghosts[j].GhostTunnel(g.LevelConfig.Speeds.Tunnel)
		ghostPulsed[j] = g.GhostPulse(j)
	}

	pacmanPulsed := g.Pacman.Pulse()
	if pacmanPulsed {
		// TODO not clear if he should stall for a specified number of frames, updates, or pulses
		// let's go with pulses for now
		if g.Pacman.StallTimer > 0 {
			g.Pacman.StallTimer -= 1
			pacmanPulsed = false
		}
	}

	ls := &g.LevelState

	if !ls.DemoMode {
		g.GhostsLeave()

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
			ls.GhostsEaten == g.Options.MaxGhosts)

		if revert {
			ls.BlueTimeout = 0
			ls.WhiteBlueTimeout = 0
		}

		g.GhostsRevert(revert)
		g.PacmanRevert(revert)

		for j := range 4 {
			if ghostPulsed[j] {
				g.Ghosts[j].SteerGhost(&g.Video, &g.Pacman, &g.Ghosts[BLINKY], &g.LevelConfig.Speeds, g.Options.GhostAi)
			}
		}

		if pacmanPulsed {
			inDir := GetJoystickDirection()
			g.Pacman.SteerPacman(&g.Video, inDir)
		}
	}

	for j := range 4 {
		if ghostPulsed[j] {
			g.Ghosts[j].MoveGhost()
		}
	}

	if pacmanPulsed {
		g.Pacman.MovePacman(&g.Video)
	}

	g.TimeoutBonus()
	g.TimeoutBonusScore()

	if ls.DemoMode {
		g.CollidePacman()
		return ThenContinue
	}

	if g.CollidePacman() {
		return WithAnim(
			(*Game).AnimPacmanDie,
			(*Game).DieStep1,
		)
	}

	return g.SurviveStep1()
}

func (g *Game) DieStep1() Return {
	g.LevelState.DotState.SavePellets(&g.Video)

	// death of pacman triggers global dot counter
	ls := &g.LevelState
	ls.GlobalDotCounterEnabled = true
	ls.DecrementLives(&g.Video)

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

	ls.DotState.DecodePellets(&g.Video)
	g.LevelInit(ls.LevelNumber)

	// TODO refactor this spaghetti
	{
		saved := &g.SavedPlayer[g.PlayerNumber]
		// these get clobbered by LevelInit...
		ls.GlobalDotCounterEnabled = saved.GlobalDotCounterEnabled
		ls.GlobalDotCounter = saved.GlobalDotCounter
		ls.DotsRemaining = saved.DotsRemaining
		ls.DotsEaten = saved.DotsEaten
	}

	ls.LevelStart()
	g.GhostsStart()
	g.Pacman.Start(g.LevelConfig.Speeds.Pacman)
	g.LevelState.BonusState.WriteBonuses(&g.Video)
	g.Video.WriteLives(g.LevelState.Lives)
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
