package option

const (
	GAME_MODE_1P int = 1
	GAME_MODE_2P int = 2
)

const (
	GHOST_AI_OFF int = 0
	GHOST_AI_ON  int = 1
)

const (
	DIFFICULTY_EASY   int = 0
	DIFFICULTY_MEDIUM int = 1
	DIFFICULTY_HARD   int = 2
)

// Game options
type Options struct {
	GameMode   int
	Difficulty int
	FrameRate  int // fps
	MaxGhosts  int // 1, 2, 3 or 4
	GhostAi    int // 0=off, 1=on
	Lives      int // 3, 5 or 10
}

func MakeOptions() Options {
	return Options{
		GameMode:   GAME_MODE_1P,
		Difficulty: DIFFICULTY_EASY,
		FrameRate:  60,
		MaxGhosts:  4,
		GhostAi:    GHOST_AI_ON,
		Lives:      3,
	}
}
