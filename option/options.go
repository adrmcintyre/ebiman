package option

// Constants for configuring the overall game play in the options screen.
// All constants are untyped ints to simplify the menu code.

const (
	ModeClassic1P  int = 1 // single player
	ModeClassic2P  int = 2 // two players (turn based)
	ModeElectric1P int = 3 // single player with charges
	ModeElectric2P int = 4 // two players with charges
)

const (
	GhostAiOff int = 0 // ghosts hunt at random
	GhostAiOn  int = 1 // ghosts can hunt actively
)

const (
	DifficulyEasy   int = 0 // easier than normal
	DifficultyMedum int = 1 // "classic" mode
	DifficultyHard  int = 2 // stupidly hard
)

// Options describes the selected game play options.
type Options struct {
	Mode       int // Mode* number of players and style
	Difficulty int // Difficulty* how hard is the game
	FrameRate  int // (unused) fps
	MaxGhosts  int // (unused) 1, 2, 3 or 4
	GhostAi    int // GhostAi*
	Lives      int // 3, 5 or 10
}

// DefaultOptions returns a sensible default set of game play options.
func DefaultOptions() Options {
	return Options{
		Mode:       ModeElectric1P,
		Difficulty: DifficultyMedum,
		FrameRate:  60,
		MaxGhosts:  4,
		GhostAi:    GhostAiOn,
		Lives:      5,
	}
}

// NumPlayers returns the number of player selected.
func (o *Options) NumPlayers() int {
	switch o.Mode {
	case ModeClassic1P, ModeElectric1P:
		return 1
	case ModeClassic2P, ModeElectric2P:
		return 2
	default:
		return 1
	}
}

// IsElectric returns true if an "electric" game mode is selected.
func (o *Options) IsElectric() bool {
	switch o.Mode {
	case ModeElectric1P, ModeElectric2P:
		return true
	default:
		return false
	}
}
