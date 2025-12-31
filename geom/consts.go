package geom

// These constants define some important locations in the maze.
var (
	HomeCentre  = Position{108, 136}
	HomeTop     = HomeCentre.Y - 4
	HomeBottom  = HomeCentre.Y + 4
	HomeExitedY = 112

	PacmanStart = Position{108, 208}
	BlinkyStart = Position{108, HomeExitedY}
	BonusPos    = Position{HomeCentre.X, 160}

	BlinkyHome = HomeCentre
	InkyHome   = HomeCentre.Add(Delta{-16, 0})
	PinkyHome  = HomeCentre
	ClydeHome  = HomeCentre.Add(Delta{16, 0})

	BlinkyScatter = TilePos(25, 0)
	InkyScatter   = TilePos(25, 36)
	PinkyScatter  = TilePos(2, 2)
	ClydeScatter  = TilePos(0, 36)
)

// PowerPills lists the location of each power pill.
var PowerPills = []Position{
	TilePos(1, 6),
	TilePos(26, 6),
	TilePos(1, 26),
	TilePos(26, 26),
}
