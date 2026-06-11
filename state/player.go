package state

import "github.com/adrmcintyre/ebiman/data"

// Player represents the current player state.
type Player struct {
	LevelNumber           int         // current level (0-based)
	Score                 int         // points scored
	Lives                 int         // lives remaining
	Pills                 Pills       // state of each pill and power pill
	DotsEaten             int         // number of dots eaten this level
	DotsRemaining         int         // number of dots left this level
	PacmanDiedThisLevel   bool        // true if pacman died in this level
	DotsSinceDeathCounter int         // dots consumed since pacman died
	BonusStatus           BonusStatus // bonuses awarded
}

// New returns a new player.
func NewPlayer() *Player {
	return &Player{}
}

// StartLevel initialises the state ready for the given level.
func (p *Player) StartLevel(levelNumber int) {
	p.LevelNumber = levelNumber
	p.PacmanDiedThisLevel = false
	p.DotsSinceDeathCounter = 0
	p.DotsRemaining = data.DotsCount
	p.DotsEaten = 0
	p.Pills.Reset()
}

// SetLives sets the current number of lives.
func (p *Player) SetLives(lives int) {
	p.Lives = lives
}

// DecrementLives removes a life.
func (p *Player) DecrementLives() {
	p.Lives -= 1
}

// AwardExtraLife adds a new life.
func (p *Player) AwardExtraLife() {
	p.Lives += 1
}
