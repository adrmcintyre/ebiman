package player

import (
	"github.com/adrmcintyre/poweraid/bonus"
	"github.com/adrmcintyre/poweraid/pill"
)

type State struct {
	LevelNumber           int        // current level (0-based)
	Score                 int        // total points scored
	Lives                 int        // number of lives remaining
	PillState             pill.State // state of each pill and power pill
	DotsEaten             int        // number of dots eaten in this level
	DotsRemaining         int        // number of dots left in this level
	PacmanDiedThisLevel   bool       // true if pacman has died during the current level
	DotsSinceDeathCounter int
	BonusState            bonus.State // bonuses awarded
}

type SavedState struct {
	State
	DotLimits [4]int // personal dot limits per ghost
}
