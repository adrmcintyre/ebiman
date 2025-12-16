package main

const (
	GAME_MODE_1P = 1
	GAME_MODE_2P = 2
)

// Game options
type GameOptions struct {
	GameMode   int
	Difficulty int  // difficulty level: 0=easy, 1=normal, 2=hard
	FrameRate  int  // fps
	MaxGhosts  int  // 1, 2, 3 or 4
	GhostAi    bool // 0=off, 1=on
	Lives      int  // 3, 5 or 10
}

func DefaultGameOptions() GameOptions {
	return GameOptions{
		GameMode:   GAME_MODE_1P,
		Difficulty: 0,
		FrameRate:  60,
		MaxGhosts:  4,
		GhostAi:    true,
		Lives:      3,
	}
}
