package main

import (
	"github.com/adrmcintyre/ebiman/audio"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/video"
)

// EatPill is called when pacman has gone over a pill.
func (g *Game) EatPill(t video.Tile) {
	charge := t.Charge()
	g.Player.Pills.NetCharge -= charge
	switch charge {
	case 0:
		g.IncrementScore(data.DotScore)
	case -1, 1:
		g.IncrementScore(data.DotScoreCharge1)
	case -2, 2:
		g.IncrementScore(data.DotScoreCharge2)
	}

	g.CountPill()
	g.Pacman.StallTimer = data.DotStall

	if g.Player.DotsEaten&1 == 0 {
		g.Audio.PlayPacmanEffect(audio.DotEatenEven)
	} else {
		g.Audio.PlayPacmanEffect(audio.DotEatenOdd)
	}
}

// EatPower is called when pacman has eaten a power pill.
func (g *Game) EatPower() {
	g.IncrementScore(data.PowerScore)
	g.CountPill()
	g.Pacman.StallTimer = data.PowerStall
	g.Pacman.Pcm = g.LevelConfig.Speeds.PacmanBlue

	g.GhostsScare()
	g.Audio.PlayBackgroundEffect(audio.EnergiserEaten)
}

// CountPill is called whenever pacman consumes a pill or power up.
func (g *Game) CountPill() {
	g.Player.DotsRemaining -= 1
	g.Player.DotsEaten += 1

	switch g.Player.DotsEaten {
	case data.FirstBonusDots, data.SecondBonusDots:
		g.DropBonus()
	}

	g.PacmanResetIdleTimer()

	if g.Player.PacmanDiedThisLevel {
		g.Player.DotsSinceDeathCounter += 1
	} else {
		g.NotifyGhostsPillEaten()
	}
}
