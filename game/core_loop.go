package game

import (
	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/data"
)

// ResetGame resets game state as if power up has just occurred.
func (g *Game) ResetGame() {
	v := g.Video
	v.ClearTiles()
	v.ClearPalette()
	v.ColorMaze(false)
	v.Write1Up()
	v.WriteHighScore(0)
	v.WriteScoreAt(1, 1, 0)
}

// ShowOptionsScreen returns a continuation that will
// display the options screen and then start a new game.
func (g *Game) ShowOptionsScreen() Return {
	return withCoro(
		(*Game).AnimOptionsScreen,
		(*Game).BeginNewGame,
	)
}

// BeginNewGame sets up for a new game, and returns a continuation
// that will show the READY animation and then start the core game loop.
func (g *Game) BeginNewGame() Return {
	g.LevelConfig.Init(0, g.Options.Difficulty)
	g.LevelState.Init(0)
	g.LevelState.LevelStart()
	g.LevelState.SetLives(g.Options.Lives)
	g.LevelState.ClearScores()
	g.LevelState.PillState.Reset()

	g.Audio.UnMute()
	g.Audio.PlaySong(audio.SongStartup)

	g.ScheduleDelay(3000)

	return withCoro(
		(*Game).AnimReady,
		(*Game).EnterNewGameLoop,
	)
}

// EnterNewGameLoop sets the core game loop running.
func (g *Game) EnterNewGameLoop() Return {
	// sync each player's saved state to be the same
	g.SavePlayerState(0)
	if g.Options.NumPlayers() > 1 {
		g.SavePlayerState(1)
	}

	g.RunningGame = true
	//ugh
	g.LevelState.FrameCounter = 0
	g.LevelState.UpdateCounter = 0

	return thenContinue
}

// UpdateState is the state update routine for the core game loop.
//
// It returns a continuation describing whether to run any animations,
// or simply continue as normal. On first entry it returns a continuation
// to display the options screen and then begin a new game.
func (g *Game) UpdateState() Return {
	g.LevelState.UpdateCounter += 1

	if !g.RunningGame {
		return g.ShowOptionsScreen()
	}

	g.GhostsTunnel()

	ghostsPulsed := g.GhostsPulse()
	pacmanPulsed := g.PacmanPulse()

	demoMode := g.LevelState.DemoMode

	if !demoMode {
		g.CheckGhostsLeaveHome()

		revert := g.ManagePanicStations()
		g.CheckGhostsRevert(revert)
		g.PacmanRevert(revert)
		g.CheckGhostsSwitchTactics(revert)

		g.GhostsSteer(ghostsPulsed)
		g.CheckGhostsReturned()
		g.PacmanSteer(pacmanPulsed)

		dotsEaten := g.LevelState.DotsEaten
		if dotsEaten > 228 {
			g.Audio.PlayBackgroundEffect(audio.Background5)
		} else if dotsEaten > 212 {
			g.Audio.PlayBackgroundEffect(audio.Background4)
		} else if dotsEaten > 180 {
			g.Audio.PlayBackgroundEffect(audio.Background3)
		} else if dotsEaten > 116 {
			g.Audio.PlayBackgroundEffect(audio.Background2)
		} else {
			g.Audio.PlayBackgroundEffect(audio.Background1)
		}

	}

	g.GhostsMove(ghostsPulsed)
	g.PacmanMove(pacmanPulsed)

	g.CheckTimeoutBonus()
	g.CheckTimeoutBonusScore()

	alive := !g.PacmanCollide()

	if demoMode {
		return thenContinue
	}

	if alive {
		return g.SurviveStep1()
	}

	return withCoro(
		(*Game).AnimPacmanDie,
		(*Game).DieStep1,
	)
}

// DieStep1 is invoked when pacman has just died (post-animation).
// Here we determine if the player can continue with any remaining
// lives, or if it is time to run the GAME OVER animation.
func (g *Game) DieStep1() Return {
	ls := &g.LevelState

	ls.PillState.Save(g.Video)

	// death of pacman triggers global dot counter
	ls.PacmanDiedThisLevel = true
	ls.DecrementLives()

	if ls.Lives > 0 {
		return g.DieStep2()
	}

	return withCoro(
		(*Game).AnimGameOver,
		(*Game).DieStep2,
	)
}

// DieStep2 determines if another player can still continue
// playing after the previous player died, or if it's time
// to return to the splash screen.
func (g *Game) DieStep2() Return {
	if !g.LoadNextPlayerState() {
		g.GameState = GameStateSplashStart
		return thenStop
	}
	return g.DieStep3()
}

// DieStep3 initialises play for the next player (the same player
// if they have lives remaining and it's a single player game),
// then schedules the READY animation.
func (g *Game) DieStep3() Return {
	ls := &g.LevelState

	g.LevelConfig.Init(ls.LevelNumber, g.Options.Difficulty)
	g.LevelState.Init(ls.LevelNumber)

	// TODO refactor this spaghetti
	{
		saved := &g.SavedPlayer[g.PlayerNumber]
		// these get clobbered by LevelInit...
		ls.PacmanDiedThisLevel = saved.PacmanDiedThisLevel
		ls.DotsSinceDeathCounter = saved.DotsSinceDeathCounter
		ls.DotsRemaining = saved.DotsRemaining
		ls.DotsEaten = saved.DotsEaten
	}

	g.LevelState.LevelStart()
	g.PacmanResetIdleTimer()

	return withCoro(
		(*Game).AnimReady,
		(*Game).SurviveStep1,
	)
}

// SurviveStep1 continues play when pacman is still alive.
// If no dots remain, the end of level animation is scheduled
// followed by starting the next level, otherwise play simply
// continues.
func (g *Game) SurviveStep1() Return {
	if g.LevelState.DotsRemaining > 0 {
		return thenContinue
	}

	return withCoro(
		(*Game).AnimEndLevel,
		(*Game).BeginNewLevel,
	)
}

// BeginNewLevel gets the next level ready, then schedules
// the READY animation, after which play continues.
func (g *Game) BeginNewLevel() Return {
	ls := &g.LevelState
	ls.LevelNumber += 1

	ls.PillState.Reset() // mark all pills as uneaten

	// level config may be different between players (due to differing level number)
	g.LevelConfig.Init(ls.LevelNumber, g.Options.Difficulty)
	g.LevelState.Init(ls.LevelNumber)
	g.LevelState.LevelStart()

	return withCoro(
		(*Game).AnimReady,
		(*Game).SurviveStep3,
	)
}

// SurviveStep3 is invoked after the READY animation for
// a new level has run.
func (g *Game) SurviveStep3() Return {
	g.PacmanResetIdleTimer()

	return thenStop
}

// ManagePanicStations manages the flashing of panicked ghosts,
// Returns true if they were panicked but aren't any more.
func (g *Game) ManagePanicStations() bool {
	ls := &g.LevelState

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
		g.Audio.StopBackgroundEffect(audio.EnergiserEaten)
	}

	return revert
}
