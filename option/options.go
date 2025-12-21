package option

// Constants for configuring the overall game play in the options screen.
// All constants are untyped ints to simplify the menu code.

const (
	GAME_MODE_1P int = 1 // single player
	GAME_MODE_2P int = 2 // two players (turn based)
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
	GameMode   int // GAME_MODE_* number of players
	Difficulty int // DIFFICULTY_* how hard is the game
	FrameRate  int // (unused) fps
	MaxGhosts  int // (unused) 1, 2, 3 or 4
	GhostAi    int // GHOST_AI_*
	Lives      int // 3, 5 or 10
}

// DefaultOptions returns a sensible default set of game play options.
func DefaultOptions() Options {
	return Options{
		GameMode:   GAME_MODE_1P,
		Difficulty: DIFFICULTY_MEDIUM,
		FrameRate:  60,
		MaxGhosts:  4,
		GhostAi:    GHOST_AI_ON,
		Lives:      3,
	}
}
