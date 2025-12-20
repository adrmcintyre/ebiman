package level

import (
	"github.com/adrmcintyre/poweraid/audio"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/option"
	"github.com/adrmcintyre/poweraid/player"
	"github.com/adrmcintyre/poweraid/video"
)

// Current level variables
type State struct {
	FrameCounter      int  // how many frames have elapsed since game start  TODO - should only be references by animations
	UpdateCounter     int  // how many game updates since game start         TODO - should only be referenced by game logic
	BlueTimeout       int  // ghosts stop being blue once update_counter passes this time
	WhiteBlueTimeout  int  // next frame at which to flash white/blue
	IsWhite           bool // are ghosts currently white?
	IsFlashing        bool // are ghosts currently flashing?
	IdleAfter         int  // pacman is considered idle when frame counter exceeds this value  TODO - should be in updates
	BonusTimeout      int  // update at which bonus is due to vanish, if non-zero
	BonusScoreTimeout int  // update wt which bonus score is due to vanish

	GhostsEaten int // ghosts eaten since last power dot

	player.State // current PlayerState

	// Game variables
	Score1    int // player1 - total points scored
	Score2    int // player2 - total points scored
	HighScore int // highest score since power-on (TODO - store this in flash memory???)
	DemoMode  bool
}

func DefaultState() State {
	return State{
		DemoMode: true,
	}
}

func (s *State) Init(i int) {
	s.LevelNumber = i
	s.PacmanDiedThisLevel = false
	s.DotsSinceDeathCounter = 0
	s.DotsRemaining = data.DOTS_COUNT
	s.DotsEaten = 0
}

// Start of level, or pacman lost a life and is restarting
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

func (s *State) SetLives(lives int) {
	s.Lives = lives
}

func (s *State) DecrementLives() {
	s.SetLives(s.Lives - 1)
}

func (s *State) AwardExtraLife() {
	s.SetLives(s.Lives + 1)
}

// -------------------------------------------------------------------------
// Scoring
// -------------------------------------------------------------------------

func (s *State) WriteScores(v *video.Video, gameMode int) {
	v.WriteHighScore(s.HighScore)
	v.WriteScoreAt(1, 1, s.Score1)
	if gameMode == option.GAME_MODE_2P {
		v.WriteScoreAt(20, 1, s.Score2)
	}
}

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

func (s *State) ClearScores() {
	s.Score1 = 0
	s.Score2 = 0
}

func (s *State) IncrementScore(playerNumber int, delta int) {
	oldScore := s.Score1
	if playerNumber == 1 {
		oldScore = s.Score2
	}
	newScore := oldScore + delta

	// pac man very generously awards one and only one extra life!
	if oldScore < data.EXTRA_LIFE_SCORE && newScore >= data.EXTRA_LIFE_SCORE {
		s.AwardExtraLife()
		// TODO nasty having this dependency
		audio.PlayEffect1(audio.Effect1_ExtraLife)

	}

	s.SetScore(playerNumber, newScore)
}
