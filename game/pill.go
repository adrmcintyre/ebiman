package game

import (
	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/bonus"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/ghost"
	"github.com/adrmcintyre/ebiman/tile"
)

// EatPill is called when pacman has gone over a pill.
func (g *Game) EatPill(t tile.Tile) {
	charge := t.Charge()
	g.LevelState.PillState.NetCharge -= charge
	switch charge {
	case 0:
		g.IncrementScore(data.DOT_SCORE)
	case -1, 1:
		g.IncrementScore(data.DOT_SCORE_CHARGE1)
	case -2, 2:
		g.IncrementScore(data.DOT_SCORE_CHARGE2)
	}

	g.CountPill()
	g.Pacman.StallTimer = data.DOT_STALL

	if g.LevelState.DotsEaten&1 == 0 {
		g.Audio.PlayPacmanEffect(audio.DotEatenEven)
	} else {
		g.Audio.PlayPacmanEffect(audio.DotEatenOdd)
	}
}

// EatPower is called when pacman has eaten a power pill.
func (g *Game) EatPower() {
	g.IncrementScore(data.POWER_SCORE)
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
	g.Audio.PlayBackgroundEffect(audio.EnergiserEaten)
}

// CountPill is called whenever pacman consumes a pill or power up.
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
