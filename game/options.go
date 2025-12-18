package game

const (
	GAME_MODE_1P = 1
	GAME_MODE_2P = 2
)

const (
	GHOST_AI_OFF = 0
	GHOST_AI_ON  = 1
)

const (
	DIFFICULTY_EASY   = 0
	DIFFICULTY_NORMAL = 1
	DIFFICULTY_HARD   = 2
)

// Game options
type GameOptions struct {
	GameMode   int
	Difficulty int
	FrameRate  int // fps
	MaxGhosts  int // 1, 2, 3 or 4
	GhostAi    int // 0=off, 1=on
	Lives      int // 3, 5 or 10
}

func DefaultGameOptions() GameOptions {
	return GameOptions{
		GameMode:   GAME_MODE_1P,
		Difficulty: DIFFICULTY_EASY,
		FrameRate:  60,
		MaxGhosts:  4,
		GhostAi:    GHOST_AI_ON,
		Lives:      3,
	}
}
