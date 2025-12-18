package game

import "github.com/adrmcintyre/poweraid/data"

func (g *Game) EatPill() {
	g.LevelState.IncrementScore(g.PlayerNumber, data.DOT_SCORE)
	g.CountDot()
	g.Pacman.StallTimer = data.DOT_STALL
}

func (g *Game) EatPower() {
	g.LevelState.IncrementScore(g.PlayerNumber, data.POWER_SCORE)
	g.CountDot()
	g.Pacman.StallTimer = data.POWER_STALL
	g.Pacman.Pcm = g.LevelConfig.Speeds.PacmanBlue

	g.LevelState.BlueTimeout = g.LevelState.UpdateCounter + g.LevelConfig.BlueTime
	g.LevelState.WhiteBlueTimeout = g.LevelState.BlueTimeout - g.LevelConfig.WhiteBlueCount*data.WHITE_BLUE_PERIOD
	g.LevelState.IsFlashing = false
	g.LevelState.IsWhite = false
	g.LevelState.GhostsEaten = 0

	// If some ghost is already scared, don't scare additional ghosts
	alreadyScared := false
	for j := range 4 {
		if g.Ghosts[j].SubMode == SUBMODE_SCARED {
			alreadyScared = true
			break
		}
	}

	if !alreadyScared {
		for j := range 4 {
			ghost := &g.Ghosts[j]
			if ghost.Mode == MODE_PLAYING || ghost.Mode == MODE_HOME {
				ghost.SetSubMode(SUBMODE_SCARED)
				ghost.Pcm = g.LevelConfig.Speeds.GhostBlue
			}
		}
	}
}

func (g *Game) CountDot() {
	g.LevelState.DotsRemaining -= 1
	g.LevelState.DotsEaten += 1

	switch g.LevelState.DotsEaten {
	case data.FIRST_BONUS_DOTS, data.SECOND_BONUS_DOTS:
		g.DropBonus()
	}

	g.PacmanResetIdleTimer()

	if g.LevelState.PacmanDiedThisLevel {
		g.LevelState.DotsSinceDeathCounter += 1
	} else {
		for j := 1; j < 4; j++ {
			ghost := &g.Ghosts[j]
			if ghost.Mode == MODE_HOME {
				ghost.DotsAtHomeCounter += 1
				break
			}
		}
	}
}
