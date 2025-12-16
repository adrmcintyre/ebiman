package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

type BonusState struct {
	BonusIndicator [7]int // last 7 bonuses awarded, most recent first
	BonusCount     int    // how many bonuses awarded so far
}

type DotState struct {
	PillBits   [30]byte // bitmap of uneaten pills
	PowerPills [4]byte  // tile at each power pill location
}

// Current level variables
type LevelState struct {
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

	PlayerState // current PlayerState

	// Game variables
	Score1    int // player1 - total points scored
	Score2    int // player2 - total points scored
	HighScore int // highest score since power-on (TODO - store this in flash memory???)
	DemoMode  bool
}

func DefaultLevelState() LevelState {
	return LevelState{
		DemoMode: true,
	}
}

func (g *Game) PacmanResetIdleTimer() {
	g.LevelState.IdleAfter = g.LevelState.FrameCounter + g.LevelConfig.IdleLimit
}

func (g *Game) IsPacmanIdle() bool {
	return g.LevelState.FrameCounter >= g.LevelState.IdleAfter

}

func (g *Game) LevelInit(levelNumber int) {
	g.LevelConfig.Init(levelNumber, g.Options.Difficulty)
	g.LevelState.Init(levelNumber)
}

func (cfg *LevelConfig) Init(levelNumber int, difficulty int) {
	levelIndex := min(levelNumber, len(data.Level)-1)
	level := data.Level[levelIndex]

	speeds := data.SpeedData[level.SpeedIndex-3]

	switch difficulty {
	case 0:
		cfg.Speeds = speeds.Easy
	case 1:
		cfg.Speeds = speeds.Medium
	case 2:
		cfg.Speeds = speeds.Hard
	}

	cfg.SwitchTactics = speeds.SwitchTactics
	cfg.DotLimits = data.DotLimit[level.DotLimitIndex]
	elroy := data.Elroy[level.ElroyIndex]
	cfg.ElroyPills1 = elroy.Pills1
	cfg.ElroyPills2 = elroy.Pills2

	blueControl := data.BlueControl[level.BlueIndex]
	cfg.BlueTime = blueControl.BlueTime
	cfg.WhiteBlueCount = blueControl.WhiteBlueCount

	if difficulty == 0 {
		cfg.BlueTime *= 2
	}

	cfg.IdleLimit = data.IdleLimit[level.IdleIndex]
	cfg.BonusType = data.BonusType[levelIndex]
	cfg.BonusInfo = data.BonusInfo[cfg.BonusType]
}

func (st *LevelState) Init(i int) {
	st.LevelNumber = i
	st.PacmanDiedThisLevel = false
	st.DotsSinceDeathCounter = 0
	st.DotsRemaining = data.DOTS_COUNT
	st.DotsEaten = 0
}

func (g *Game) LevelStart() {
	g.LevelState.LevelStart()
	g.PacmanResetIdleTimer()
}

// Start of level, or pacman lost a life and is restarting
func (ls *LevelState) LevelStart() {
	ls.GhostsEaten = 0
	ls.FrameCounter = 0
	ls.UpdateCounter = 0
	ls.BlueTimeout = 0
	ls.WhiteBlueTimeout = 0
	ls.IsFlashing = false
	ls.BonusTimeout = 0
	ls.BonusScoreTimeout = 0
}

func (g *Game) WritePlayerUp(v *video.Video) {
	v.WritePlayerUp(g.PlayerNumber)
}

func (g *Game) ClearPlayerUp(v *video.Video) {
	v.ClearPlayerUp(g.PlayerNumber)
}

// -------------------------------------------------------------------------
// Lives
// -------------------------------------------------------------------------

func (ls *LevelState) SetLives(lives int) {
	ls.Lives = lives
}

func (ls *LevelState) DecrementLives() {
	ls.SetLives(ls.Lives - 1)
}

func (ls *LevelState) AwardExtraLife() {
	ls.SetLives(ls.Lives + 1)
}

// -------------------------------------------------------------------------
// Scoring
// -------------------------------------------------------------------------

func (ls *LevelState) WriteScores(v *video.Video, gameMode int) {
	v.WriteHighScore(ls.HighScore)
	v.WriteScoreAt(1, 1, ls.Score1)
	if gameMode == GAME_MODE_2P {
		v.WriteScoreAt(20, 1, ls.Score2)
	}
}

func (ls *LevelState) SetScore(playerNumber int, score int) {
	if score > ls.HighScore {
		ls.HighScore = score
	}

	if playerNumber == 0 {
		ls.Score1 = score
	} else {
		ls.Score2 = score
	}
}

func (ls *LevelState) ClearScores() {
	ls.Score1 = 0
	ls.Score2 = 0
}

func (ls *LevelState) IncrementScore(playerNumber int, delta int) {
	oldScore := ls.Score1
	if playerNumber == 1 {
		oldScore = ls.Score2
	}
	newScore := oldScore + delta

	// pac man very generously awards one and only one extra life!
	if oldScore < data.EXTRA_LIFE_SCORE && newScore >= data.EXTRA_LIFE_SCORE {
		ls.AwardExtraLife()
	}

	ls.SetScore(playerNumber, newScore)
}

// -------------------------------------------------------------------------
// Bonuses (fruit)
// -------------------------------------------------------------------------

func (bs *BonusState) WriteBonuses(v *video.Video) {
	tileBase := tile.SPACE_BASE
	pal := palette.BLACK
	j := 0
	for i := range 7 {
		if i+bs.BonusCount >= 7 {
			info := &data.BonusInfo[bs.BonusIndicator[j]]
			j++
			tileBase = info.BaseTile
			pal = info.Pal
		}
		v.SetStatusQuad(12+i*2, tileBase, pal)
	}
}

func (bs *BonusState) ClearBonuses() {
	bs.BonusCount = 0
}

func (bs *BonusState) AddBonus(bonus int) {
	bs.BonusCount += 1
	copy(bs.BonusIndicator[1:], bs.BonusIndicator[:6])
	bs.BonusIndicator[0] = bonus
}
