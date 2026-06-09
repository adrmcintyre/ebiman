package state

import (
	"github.com/adrmcintyre/ebiman/pill"
)

// PlayerState represents the current player state.
type PlayerState struct {
	LevelNumber           int         // current level (0-based)
	Score                 int         // points scored
	Lives                 int         // lives remaining
	PillState             pill.State  // state of each pill and power pill
	DotsEaten             int         // number of dots eaten this level
	DotsRemaining         int         // number of dots left this level
	PacmanDiedThisLevel   bool        // true if pacman died in this level
	DotsSinceDeathCounter int         // dots consumed since pacman died
	BonusStatus           BonusStatus // bonuses awarded
}

// SavedPlayerState represents a saved player state.
// It is used for retaing the state of a player while the other player is active.
type SavedPlayerState struct {
	PlayerState
	DotLimits [4]int // personal dot limits per ghost
}
