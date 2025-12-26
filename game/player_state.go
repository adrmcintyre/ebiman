package game

import (
	"github.com/adrmcintyre/ebiman/option"
)

// SavePlayerState saves the state of the specified player.
func (g *Game) SavePlayerState(i int) {
	p := &g.SavedPlayer[i]
	ls := &g.LevelState

	p.State = ls.State
	for i, gh := range g.Ghosts {
		p.DotLimits[i] = gh.DotLimit
	}
}

// LoadPlayerState reinitialises the active player
// state from the specified saved state.
func (g *Game) LoadPlayerState(i int) {
	p := &g.SavedPlayer[i]
	ls := &g.LevelState

	ls.State = p.State
	for i, gh := range g.Ghosts {
		gh.DotLimit = p.DotLimits[i]
	}
}

// LoadNextPlayerState saves the state of the current player,
// and loads in the next player with lives remaining.
// If no such player is found, returns false.
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
