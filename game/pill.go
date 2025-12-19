package game

import (
	"github.com/adrmcintyre/poweraid/bonus"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/ghost"
)

func (g *Game) EatPill() {
	g.LevelState.IncrementScore(g.PlayerNumber, data.DOT_SCORE)
	g.CountPill()
	g.Pacman.StallTimer = data.DOT_STALL
}

func (g *Game) EatPower() {
	g.LevelState.IncrementScore(g.PlayerNumber, data.POWER_SCORE)
	g.CountPill()
	g.Pacman.StallTimer = data.POWER_STALL
	g.Pacman.Pcm = g.LevelConfig.Speeds.PacmanBlue

	g.LevelState.BlueTimeout = g.LevelState.UpdateCounter + g.LevelConfig.BlueTime
	g.LevelState.WhiteBlueTimeout = g.LevelState.BlueTimeout - g.LevelConfig.WhiteBlueCount*data.WHITE_BLUE_PERIOD
	g.LevelState.IsFlashing = false
	g.LevelState.IsWhite = false
	g.LevelState.GhostsEaten = 0

	// If some ghost is already scared, don't scare additional ghosts
	alreadyScared := false
	for _, gh := range g.Ghosts {
		if gh.SubMode == ghost.SUBMODE_SCARED {
			alreadyScared = true
			break
		}
	}

	if !alreadyScared {
		for _, gh := range g.Ghosts {
			if gh.Mode == ghost.MODE_PLAYING || gh.Mode == ghost.MODE_HOME {
				gh.SetSubMode(ghost.SUBMODE_SCARED)
				gh.Pcm = g.LevelConfig.Speeds.GhostBlue
			}
		}
	}
}

func (g *Game) CountPill() {
	g.LevelState.DotsRemaining -= 1
	g.LevelState.DotsEaten += 1

	switch g.LevelState.DotsEaten {
	case bonus.FIRST_BONUS_DOTS, bonus.SECOND_BONUS_DOTS:
		g.DropBonus()
	}

	g.PacmanResetIdleTimer()

	if g.LevelState.PacmanDiedThisLevel {
		g.LevelState.DotsSinceDeathCounter += 1
	} else {
		for _, gh := range g.Ghosts {
			if gh.Id != ghost.BLINKY && gh.Mode == ghost.MODE_HOME {
				gh.DotsAtHomeCounter += 1
				break
			}
		}
	}
}
