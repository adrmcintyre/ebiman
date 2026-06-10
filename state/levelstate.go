package state

// LevelState describes the dynamic state of an in-play level.
type LevelState struct {
	FrameCounter      int  // frames elapsed since game start - TODO should only be references by animations
	UpdateCounter     int  // updates elapsed since game start - TODO should only be referenced by game logic
	BlueTimeout       int  // ghosts stop being blue once UpdateCounter exceeds this
	WhiteBlueTimeout  int  // next frame at which to flash white/blue
	IsWhite           bool // are ghosts currently white?
	IsFlashing        bool // are ghosts currently flashing?
	IdleAfter         int  // pacman is idle when FrameCounter exceeds this - TODO should be in updates
	BonusTimeout      int  // bonus vanishes when UpdateCounter exceeds this; visible if non-zero
	BonusScoreTimeout int  // bonus score vanishes when UpdateCounter exceeds this; visible if non-zero
	GhostsEaten       int  // ghosts eaten since last power dot
}

// DefaultLevelState returns a State uninitialised, except for DemoMode being enabled.
func DefaultLevelState() LevelState {
	return LevelState{}
}

// LevelStart initialises the state ready for the current level
// to start (or restart if pacman lost a life and we are continuing).
func (s *LevelState) LevelStart() {
	s.GhostsEaten = 0
	s.FrameCounter = 0
	s.UpdateCounter = 0
	s.BlueTimeout = 0
	s.WhiteBlueTimeout = 0
	s.IsFlashing = false
	s.BonusTimeout = 0
	s.BonusScoreTimeout = 0
}
