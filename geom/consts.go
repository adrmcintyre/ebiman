package geom

// These constants define some important locations in the maze.
var (
	PACMAN_START  = Position{108, 208}
	BLINKY_START  = Position{108, HOME_EXITED_Y}
	HOME_CENTRE   = Position{108, 136}
	HOME_TOP      = HOME_CENTRE.Y - 4
	HOME_BOTTOM   = HOME_CENTRE.Y + 4
	HOME_EXITED_Y = 112
	BONUS_POS     = Position{HOME_CENTRE.X, 160}

	BLINKY_HOME = HOME_CENTRE
	INKY_HOME   = HOME_CENTRE.Add(Delta{-16, 0})
	PINKY_HOME  = HOME_CENTRE
	CLYDE_HOME  = HOME_CENTRE.Add(Delta{16, 0})

	BLINKY_SCATTER = TilePos(25, 0)
	INKY_SCATTER   = TilePos(25, 36)
	PINKY_SCATTER  = TilePos(2, 2)
	CLYDE_SCATTER  = TilePos(0, 36)
)

// POWER_PILLS lists the location of each power pill
var POWER_PILLS = []Position{
	TilePos(1, 6),
	TilePos(26, 6),
	TilePos(1, 26),
	TilePos(26, 26),
}
