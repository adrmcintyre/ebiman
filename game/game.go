package game

import (
	"github.com/adrmcintyre/poweraid/audio"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/option"
)

func (g *Game) ResetGame() {
	v := &g.Video
	v.ClearTiles()
	v.ClearPalette()
	v.ColorMaze()
	v.Write1Up()
	v.WriteHighScore(0)
	v.WriteScoreAt(1, 1, 0)
}

func (g *Game) ShowOptionsScreen() Return {
	return WithAnim(
		(*Game).AnimOptionsScreen,
		(*Game).BeginNewGame,
	)
}

func (g *Game) BeginNewGame() Return {
	g.LevelConfig.Init(0, g.Options.Difficulty)
	g.LevelState.Init(0)
	g.LevelState.LevelStart()
	g.LevelState.SetLives(g.Options.Lives)
	g.LevelState.ClearScores()
	g.LevelState.PillState.Reset()

	audio.PlaySong(audio.SongStartup)
	return WithAnim(
		(*Game).AnimReady,
		(*Game).EnterNewGameLoop,
	)
}

func (g *Game) EnterNewGameLoop() Return {
	// sync each player's saved state to be the same
	g.SavePlayerState(0)
	if g.Options.GameMode == option.GAME_MODE_2P {
		g.SavePlayerState(1)
	}

	g.RunningGame = true
	//ugh
	g.LevelState.FrameCounter = 0
	g.LevelState.UpdateCounter = 0

	return ThenContinue
}

func (g *Game) UpdateState() Return {
	g.LevelState.UpdateCounter += 1

	if !g.RunningGame {
		return g.ShowOptionsScreen()
	}

	ghostsPulsed := g.GhostsPulse()
	pacmanPulsed := g.PacmanPulse()

	demoMode := g.LevelState.DemoMode

	if !demoMode {
		g.GhostsLeaveHome()

		g.PanicStations()

		g.GhostsSteer(ghostsPulsed)
		g.PacmanSteer(pacmanPulsed)

		dotsEaten := g.LevelState.DotsEaten
		if dotsEaten > 228 {
			audio.PlayBackgroundEffect(audio.Background5)
		} else if dotsEaten > 212 {
			audio.PlayBackgroundEffect(audio.Background4)
		} else if dotsEaten > 180 {
			audio.PlayBackgroundEffect(audio.Background3)
		} else if dotsEaten > 116 {
			audio.PlayBackgroundEffect(audio.Background2)
		} else {
			audio.PlayBackgroundEffect(audio.Background1)
		}

	}

	g.GhostsMove(ghostsPulsed)
	g.PacmanMove(pacmanPulsed)

	g.TimeoutBonus()
	g.TimeoutBonusScore()

	alive := !g.PacmanCollide()

	if demoMode {
		return ThenContinue
	}

	if alive {
		return g.SurviveStep1()
	}

	return WithAnim(
		(*Game).AnimPacmanDie,
		(*Game).DieStep1,
	)
}

func (g *Game) DieStep1() Return {
	ls := &g.LevelState

	ls.PillState.Save(&g.Video)

	// death of pacman triggers global dot counter
	ls.PacmanDiedThisLevel = true
	ls.DecrementLives()

	if ls.Lives > 0 {
		return g.DieStep2()
	}

	return WithAnim(
		(*Game).AnimGameOver,
		(*Game).DieStep2,
	)
}

func (g *Game) DieStep2() Return {
	if !g.LoadNextPlayerState() {
		g.Action = ActionSplash
		return ThenStop
	}
	return g.DieStep3()
}

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

	return WithAnim(
		(*Game).AnimReady,
		(*Game).SurviveStep1,
	)
}

func (g *Game) SurviveStep1() Return {
	if g.LevelState.DotsRemaining > 0 {
		return ThenContinue
	}

	return WithAnim(
		(*Game).AnimEndLevel,
		(*Game).BeginNewLevel,
	)
}

func (g *Game) BeginNewLevel() Return {
	ls := &g.LevelState
	ls.LevelNumber += 1

	ls.PillState.Reset() // mark all pills as uneaten

	// level config may be different between players (due to differing level number)
	g.LevelConfig.Init(ls.LevelNumber, g.Options.Difficulty)
	g.LevelState.Init(ls.LevelNumber)
	g.LevelState.LevelStart()

	return WithAnim(
		(*Game).AnimReady,
		(*Game).SurviveStep3,
	)
}

func (g *Game) SurviveStep3() Return {
	g.PacmanResetIdleTimer()

	return ThenStop
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
		audio.StopBackgroundEffect(audio.EnergiserEaten)
	}

	g.GhostsRevert(revert)
	g.GhostsSwitchTactics(revert)
	g.PacmanRevert(revert)
}
