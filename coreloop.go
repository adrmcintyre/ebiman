package main

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
	playerNumber := 0
	levelNumber := 0

	g.LevelConfig.Init(levelNumber, g.Options.Difficulty)
	g.LevelState.LevelStart()
	g.ClearScores()

	for i := range g.Options.NumPlayers() {
		player := &g.Players[i]
		player.SetLives(g.Options.Lives)
		player.Init(levelNumber)
		player.Pills.Reset()
	}
	g.PlayerNumber = playerNumber
	g.Player = &g.Players[g.PlayerNumber]

	g.Audio.UnMute()
	if g.Options.IsElectric() {
		g.Audio.PlaySong(audio.SongAlternateStartup)
	} else {
		g.Audio.PlaySong(audio.SongStartup)
	}

	g.ScheduleDelay(3000)

	return withCoro(
		(*Game).AnimReady,
		(*Game).EnterNewGameLoop,
	)
}

// EnterNewGameLoop sets the core game loop running.
func (g *Game) EnterNewGameLoop() Return {
	g.LevelState.FrameCounter = 0
	g.LevelState.UpdateCounter = 0
	g.RunningGame = true

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

	if !g.DemoMode {
		g.CheckGhostsLeaveHome()

		revert := g.ManagePanicStations()
		g.CheckGhostsRevert(revert)
		g.PacmanRevert(revert)
		g.CheckGhostsSwitchTactics(revert)

		g.GhostsSteer(ghostsPulsed)
		g.CheckGhostsReturned()
		g.PacmanSteer(pacmanPulsed)
		g.UpdateSirenAudio()
	}

	g.GhostsMove(ghostsPulsed)
	g.PacmanMove(pacmanPulsed)

	g.CheckTimeoutBonus()
	g.CheckTimeoutBonusScore()

	dead := g.PacmanCollide() || g.ElectricOverload()

	if g.DemoMode {
		return thenContinue
	}

	if dead {
		return withCoro(
			(*Game).AnimPacmanDie,
			(*Game).DieStep1,
		)
	}

	return g.SurviveStep1()

}

func (g *Game) UpdateSirenAudio() {
	var effect audio.BackgroundEffectId
	eaten := g.Player.DotsEaten
	switch {
	case eaten <= 116:
		effect = audio.Background1
	case eaten <= 180:
		effect = audio.Background2
	case eaten <= 212:
		effect = audio.Background3
	case eaten <= 228:
		effect = audio.Background4
	default:
		effect = audio.Background5
	}
	g.Audio.PlayBackgroundEffect(effect)
}

// DieStep1 is invoked when pacman has just died (post-animation).
// Here we determine if the player can continue with any remaining
// lives, or if it is time to run the GAME OVER animation.
func (g *Game) DieStep1() Return {
	player := g.Player
	player.Pills.Save(g.Video)

	// death of pacman triggers global dot counter
	player.PacmanDiedThisLevel = true
	player.DecrementLives()

	if player.Lives > 0 {
		return g.DieStep2()
	}

	g.RegisterScore()

	return withCoro(
		(*Game).AnimGameOver,
		(*Game).DieStep2,
	)
}

// DieStep2 determines if another player can still continue
// playing after the previous player died, or if it's time
// to return to the splash screen.
func (g *Game) DieStep2() Return {
	if !g.NextPlayer() {
		g.GameState = GameStateSplashStart
		return thenStop
	}
	return g.DieStep3()
}

// DieStep3 initialises play for the next player (the same player
// if they have lives remaining and it's a single player game),
// then schedules the READY animation.
func (g *Game) DieStep3() Return {
	g.LevelConfig.Init(g.Player.LevelNumber, g.Options.Difficulty)
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
	if g.Player.DotsRemaining > 0 {
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
	player := g.Player
	player.LevelNumber += 1
	player.Pills.Reset() // mark all pills as uneaten
	player.Init(player.LevelNumber)

	// level config may be different between players (due to differing level number)
	g.LevelConfig.Init(player.LevelNumber, g.Options.Difficulty)
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
		ls.WhiteBlueTimeout += data.WhiteBluePeriod
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
