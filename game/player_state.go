package game

import (
	"github.com/adrmcintyre/poweraid/option"
)

// TODO finish this
func (g *Game) SavePlayerState(i int) {
	p := &g.SavedPlayer[i]
	ls := &g.LevelState

	p.State = ls.State
	for i, gh := range g.Ghosts {
		p.DotLimits[i] = gh.DotLimit
	}
}

// TODO finish this
func (g *Game) LoadPlayerState(i int) {
	p := &g.SavedPlayer[i]
	ls := &g.LevelState

	ls.State = p.State
	for i, gh := range g.Ghosts {
		gh.DotLimit = p.DotLimits[i]
	}
}

func (g *Game) LoadNextPlayerState() bool {
	g.SavePlayerState(g.PlayerNumber)

	n := 1
	if g.Options.GameMode == option.GAME_MODE_2P {
		n = 2
	}
	for range n {
		g.PlayerNumber = (g.PlayerNumber + 1) % n
		if g.SavedPlayer[g.PlayerNumber].Lives > 0 {
			g.LoadPlayerState(g.PlayerNumber)
			return true
		}
	}
	return false
}
