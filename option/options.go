package option

// Constants for configuring the overall game play in the options screen.
// All constants are untyped ints to simplify the menu code.

const (
	MODE_CLASSIC_1P  int = 1 // single player
	MODE_CLASSIC_2P  int = 2 // two players (turn based)
	MODE_ELECTRIC_1P int = 3 // single player with charges
	MODE_ELECTRIC_2P int = 4 // two players with charges
)

const (
	GHOST_AI_OFF int = 0 // ghosts hunt at random
	GHOST_AI_ON  int = 1 // ghosts can hunt actively
)

const (
	DIFFICULTY_EASY   int = 0 // easier than normal
	DIFFICULTY_MEDIUM int = 1 // "classic" mode
	DIFFICULTY_HARD   int = 2 // stupidly hard
)

// Options describes the selected game play options.
type Options struct {
	Mode       int // MODE_* number of players and style
	Difficulty int // DIFFICULTY_* how hard is the game
	FrameRate  int // (unused) fps
	MaxGhosts  int // (unused) 1, 2, 3 or 4
	GhostAi    int // GHOST_AI_*
	Lives      int // 3, 5 or 10
}

// DefaultOptions returns a sensible default set of game play options.
func DefaultOptions() Options {
	return Options{
		Mode:       MODE_ELECTRIC_1P,
		Difficulty: DIFFICULTY_MEDIUM,
		FrameRate:  60,
		MaxGhosts:  4,
		GhostAi:    GHOST_AI_ON,
		Lives:      5,
	}
}

// NumPlayers returns the number of player selected.
func (o *Options) NumPlayers() int {
	switch o.Mode {
	case MODE_CLASSIC_1P, MODE_ELECTRIC_1P:
		return 1
	case MODE_CLASSIC_2P, MODE_ELECTRIC_2P:
		return 2
	default:
		return 1
	}
}

// IsElectric returns true if an "electric" game mode is selected.
func (o *Options) IsElectric() bool {
	switch o.Mode {
	case MODE_ELECTRIC_1P, MODE_ELECTRIC_2P:
		return true
	default:
		return false
	}
}
