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

	// tunnel position in tile co-ordinates
	TunnelTileY     = 17
	TunnelTileLeft  = 5
	TunnelTileRight = 22

	// tunnel limits in pixels
	TunnelLeft  = 4
	TunnelRight = 220
	TunnelWidth = TunnelRight - TunnelLeft

	// highest and lowest valid pixel positions for actors
	MazeTop    = 12
	MazeBottom = 260

	// SafeTiles are tiles the ghosts will turn up to chase pac man
	SafeTiles = []Position{
		TilePos(12, 13),
		TilePos(15, 13),
		TilePos(12, 25),
		TilePos(15, 25),
	}

	// PowerPills lists the location of each power pill.
	PowerPills = []Position{
		TilePos(1, 6),
		TilePos(26, 6),
		TilePos(1, 26),
		TilePos(26, 26),
	}
)
