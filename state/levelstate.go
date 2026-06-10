package state

// LevelState describes the dynamic state of an in-play level.
type LevelState struct {
	FrameCounter        int  // frames elapsed since game start - TODO should only be references by animations
	UpdateCounter       int  // updates elapsed since game start - TODO should only be referenced by game logic
	GhostsScaredTimeout int  // ghosts stop being blue once UpdateCounter exceeds this
	GhostsFlashTimeout  int  // next frame at which to flash white/blue
	GhostsAreWhite      bool // are ghosts currently white?
	GhostsAreFlashing   bool // are ghosts currently flashing?
	IdleAfter           int  // pacman is idle when FrameCounter exceeds this - TODO should be in updates
	BonusTimeout        int  // bonus vanishes when UpdateCounter exceeds this; visible if non-zero
	BonusScoreTimeout   int  // bonus score vanishes when UpdateCounter exceeds this; visible if non-zero
	GhostsEaten         int  // ghosts eaten since last power dot
}

// NewLevelState returns a new state ready for the current level to start
// (or restart if pacman lost a life and we are continuing).
func NewLevelState() *LevelState {
	return &LevelState{}
}
