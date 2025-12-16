package main

type PlayerState struct {
	PacmanDiedThisLevel   bool // true if pacman has died during the current level
	DotsSinceDeathCounter int
	DotsRemaining         int        // number of dots left in this level
	DotsEaten             int        // number of dots eaten in this level
	BonusState            BonusState // bonuses awarded
	Score                 int        // total points scored
	Lives                 int        // number of lives remaining
	LevelNumber           int        // current level (0-based)
	DotState              DotState   // state of each pill and power pill
}

type SavedPlayerState struct {
	PlayerState
	DotLimits [4]int // personal dot limits per ghost
}

// TODO finish this
func (g *Game) SavePlayerState(i int) {
	p := &g.SavedPlayer[i]
	ls := &g.LevelState

	p.PlayerState = ls.PlayerState
	for i := range 4 {
		p.DotLimits[i] = g.Ghosts[i].DotLimit
	}
}

// TODO finish this
func (g *Game) LoadPlayerState(i int) {
	p := &g.SavedPlayer[i]
	ls := &g.LevelState

	ls.PlayerState = p.PlayerState
	for i := range 4 {
		g.Ghosts[i].DotLimit = p.DotLimits[i]
	}
}

func (g *Game) LoadNextPlayerState() bool {
	g.SavePlayerState(g.PlayerNumber)

	n := 1
	if g.Options.GameMode == GAME_MODE_2P {
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
