package level

import (
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/player"
	"github.com/adrmcintyre/ebiman/video"
)

// State describes the dynamic state of an in-play level.
type State struct {
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

	player.State // current player's State

	// Game variables
	Score1    int  // player1 - total points scored
	Score2    int  // player2 - total points scored
	HighScore int  // highest score since power-on
	DemoMode  bool // when in demo mode certain behaviour is turned off
}

// DefaultState returns a State uninitialised, except for DemoMode being enabled.
func DefaultState() State {
	return State{
		DemoMode: true,
	}
}

// Init initialises the state ready for the given level.
func (s *State) Init(levelNumber int) {
	s.LevelNumber = levelNumber
	s.PacmanDiedThisLevel = false
	s.DotsSinceDeathCounter = 0
	s.DotsRemaining = data.DOTS_COUNT
	s.DotsEaten = 0
}

// LevelStart initialises the state ready for the current level
// to start (or restart if pacman lost a life and we are continuing).
func (s *State) LevelStart() {
	s.GhostsEaten = 0
	s.FrameCounter = 0
	s.UpdateCounter = 0
	s.BlueTimeout = 0
	s.WhiteBlueTimeout = 0
	s.IsFlashing = false
	s.BonusTimeout = 0
	s.BonusScoreTimeout = 0
}

// -------------------------------------------------------------------------
// Lives
// -------------------------------------------------------------------------

// SetLives sets the current number of lives.
func (s *State) SetLives(lives int) {
	s.Lives = lives
}

// DecrementLives removes a life.
func (s *State) DecrementLives() {
	s.SetLives(s.Lives - 1)
}

// AwardExtraLife adds a new life.
func (s *State) AwardExtraLife() {
	s.SetLives(s.Lives + 1)
}

// -------------------------------------------------------------------------
// Scoring
// -------------------------------------------------------------------------

// WriteScores writes the tiles for displaying the current
// high-score and player(s) score(s) into the top status area.
func (s *State) WriteScores(v *video.Video, gameMode int) {
	v.WriteHighScore(s.HighScore)
	v.WriteScoreAt(1, 1, s.Score1)
	v.WriteChargeAt(20, 1, s.PillState.NetCharge)
	//	if gameMode == option.GAME_MODE_2P {
	//		v.WriteScoreAt(20, 1, s.Score2)
	//	}
}

// SetScore records the specified player's latest score,
// updating the high-score if appropriate.
func (s *State) SetScore(playerNumber int, score int) {
	if score > s.HighScore {
		s.HighScore = score
	}

	if playerNumber == 0 {
		s.Score1 = score
	} else {
		s.Score2 = score
	}
}

// ClearScores resets both players' scores, leaving the high-score intact.
func (s *State) ClearScores() {
	s.Score1 = 0
	s.Score2 = 0
}
