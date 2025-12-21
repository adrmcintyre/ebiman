package player

import (
	"github.com/adrmcintyre/poweraid/bonus"
	"github.com/adrmcintyre/poweraid/pill"
)

// State represents the current player state.
type State struct {
	LevelNumber           int          // current level (0-based)
	Score                 int          // points scored
	Lives                 int          // lives remaining
	PillState             pill.State   // state of each pill and power pill
	DotsEaten             int          // number of dots eaten this level
	DotsRemaining         int          // number of dots left this level
	PacmanDiedThisLevel   bool         // true if pacman died in this level
	DotsSinceDeathCounter int          // dots consumed since pacman died
	BonusStatus           bonus.Status // bonuses awarded
}

// SavedState represents a saved player state.
// It is used for retaing the state of a player while the other player is active.
type SavedState struct {
	State
	DotLimits [4]int // personal dot limits per ghost
}
