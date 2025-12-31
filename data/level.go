package data

// LevelEntry describes key attributes of a level.
type LevelEntry struct {
	SpeedIndex    int // an index into the SpeedData array
	unused        int //
	DotLimitIndex int // an index into the DotLimit array
	ElroyIndex    int // an index into the Elroy array
	BlueIndex     int // an index into the BlueControl array
	IdleIndex     int // an index into the IdleLimit array
}

// Level provides the LevelEntry data for each level in the game.
// Things become increasingly frantic as the levels progress,
// until level 21 when subsequent levels repeat level 21.
var Level = [21]LevelEntry{
	{3, 1, 1, 0, 2, 0}, // level 1
	{4, 1, 2, 1, 3, 0}, // level 2
	{4, 1, 3, 2, 4, 1}, // level 3
	{4, 2, 3, 2, 5, 1}, // level 4
	{5, 0, 3, 2, 6, 2}, // level 5
	{5, 1, 3, 3, 3, 2}, // level 6
	{5, 2, 3, 3, 6, 2}, // level 7
	{5, 2, 3, 3, 6, 2}, // level 8
	{5, 0, 3, 4, 7, 2}, // level 9
	{5, 1, 3, 4, 3, 2}, // level 10
	{5, 2, 3, 4, 6, 2}, // level 11
	{5, 2, 3, 5, 7, 2}, // level 12
	{5, 0, 3, 5, 7, 2}, // level 13
	{5, 2, 3, 5, 5, 2}, // level 14
	{5, 1, 3, 6, 7, 2}, // level 15
	{5, 2, 3, 6, 7, 2}, // level 16
	{5, 2, 3, 6, 8, 2}, // level 17
	{5, 2, 3, 6, 7, 2}, // level 18
	{5, 2, 3, 7, 8, 2}, // level 19
	{5, 2, 3, 7, 8, 2}, // level 20
	{6, 2, 3, 7, 8, 2}, // level 21+
}

// Speeds describes the speed of movement of each element in the game.
type Speeds struct {
	Pacman     PCM // pacman's normal speed
	PacmanBlue PCM // pacman's speed when powered up
	Elroy2     PCM // blinky's speed when "cruise elroy" mode triggered for the second time
	Elroy1     PCM // blinky's speed when "cruise elroy" mode triggered for the first time
	Ghost      PCM // the ghosts' normal speed
	GhostBlue  PCM // the ghosts' speed when frightened
	Tunnel     PCM // the ghosts' speed when navigating the tunnel
}

// A SwitchTacticsEntry describes when the ghosts switch between
// their chase and scatter behaviours (in reverse order of frame count).
type SwitchTacticsEntry [7]int

// A SpeedDataEntry describes the speeds and tactics-switching behaviour
// of a level for various difficulty settings.
type SpeedDataEntry struct {
	Easy          Speeds             // speeds to use in "easy" mode
	Medium        Speeds             // speeds to use in "medium" mode
	Hard          Speeds             // speeds to use in "hard" mode
	SwitchTactics SwitchTacticsEntry // when ghosts switch tactics
}

// SpeedData defines the speed and tactics-switching behaviour for
// groups of levels.
var SpeedData = [4]SpeedDataEntry{
	// Indexes are offset by +3

	// 3 - level 1
	{
		Speeds{PCM90, PCM95, PCM80, PCM75, PCM70, PCM45, PCM40},
		Speeds{PCM80, PCM90, PCM85, PCM80, PCM75, PCM50, PCM40},
		Speeds{PCM80, PCM90, PCM90, PCM85, PCM80, PCM55, PCM45},
		SwitchTacticsEntry{
			84 * FPS,
			79 * FPS,
			59 * FPS,
			54 * FPS,
			34 * FPS,
			27 * FPS,
			7 * FPS,
		},
	},

	// 4 - levels 2-4
	{
		Speeds{PCM95, PCM100, PCM90, PCM85, PCM80, PCM50, PCM40},
		Speeds{PCM90, PCM95, PCM95, PCM90, PCM85, PCM55, PCM45},
		Speeds{PCM90, PCM95, PCM100, PCM95, PCM90, PCM60, PCM50},
		SwitchTacticsEntry{
			0xFFFF,
			0xFFFE,
			59 * FPS,
			54 * FPS,
			34 * FPS,
			27 * FPS,
			7 * FPS,
		},
	},

	// 5 - levels 5-20
	{
		Speeds{PCM105, PCM105, PCM100, PCM95, PCM90, PCM55, PCM45},
		Speeds{PCM100, PCM100, PCM105, PCM100, PCM95, PCM60, PCM50},
		Speeds{PCM100, PCM100, PCM110, PCM105, PCM100, PCM65, PCM55},
		SwitchTacticsEntry{
			0xFFFF,
			0xFFFE,
			55 * FPS,
			50 * FPS,
			30 * FPS,
			25 * FPS,
			5 * FPS,
		},
	},

	// 6 - levels 21+
	{
		// Energizers have no effect on these levels, so pacman_blue and
		// ghost_blue are unused - we arbitrarily set them to 0.
		Speeds{PCM95, 1, PCM105, PCM95, PCM90, 0, PCM45},
		Speeds{PCM90, 0, PCM105, PCM100, PCM95, 0, PCM50},
		Speeds{PCM90, 0, PCM110, PCM105, PCM100, 0, PCM55},
		SwitchTacticsEntry{
			0xFFFF,
			0xFFFE,
			55 * FPS,
			50 * FPS,
			30 * FPS,
			25 * FPS,
			5 * FPS,
		},
	},
}

// A DotLimitEntry describes how many dots pacman must consume
// while a ghost is at home before it should be released.
type DotLimitEntry struct {
	Pinky int
	Inky  int
	Clyde int
}

// DotLimit defines a DotLimitEntry for each LevelEntry.DotLimitIndex
var DotLimit = [4]DotLimitEntry{
	{20, 30, 70}, // 0 - this entry appears to be unused
	{0, 30, 60},  // 1
	{0, 0, 50},   // 2
	{0, 0, 0},    // 3
}

// An ElroyEntry describes when blinky's "cruise elroy" mode is
// triggered for the first and second time, after the given
// number of pills remain.
type ElroyEntry struct {
	Pills1 int // when cruise elroy 1 is triggered
	Pills2 int // when cruise elroy 2 is triggered
}

// Elroy defines an ElroyEntry for each LevelEntry.ElroyIndex
var Elroy = [9]ElroyEntry{
	{20, 10},  // 0
	{30, 15},  // 1
	{40, 20},  // 2
	{50, 25},  // 3
	{60, 30},  // 4
	{80, 40},  // 5
	{100, 50}, // 6
	{120, 60}, // 7
	{140, 70}, // 8
}

// A BlueControlEntry describes how long the ghosts stay blue after
// pacman eats a power up.
type BlueControlEntry struct {
	BlueTime       int // total time to remain blue (including flashes)
	WhiteBlueCount int // number of white/blue flashes
}

// BlueControl defines a BlueControlEntry for each LevelEntry.BlueIndex
var BlueControl = [9]BlueControlEntry{
	{8 * 4 * UPS, 9}, // 0  (not used)
	{7 * 4 * UPS, 9}, // 1  (not used)
	{6 * 4 * UPS, 9}, // 2
	{5 * 4 * UPS, 9}, // 3
	{4 * 4 * UPS, 9}, // 4
	{3 * 4 * UPS, 9}, // 5
	{2 * 4 * UPS, 9}, // 6
	{1 * 4 * UPS, 5}, // 7
	{1, 0},           // 8
}

// IdleLimit defines when is ghost is released due to pacman being idle
// (not eating) for a given number of frames. There is an entry for
// each LevelEntry.IdleLimitIndex
//
// Note - the Pacman Dossier claims these are 4, 4, 3, whereas comments in
// mspacman.asm claim they are 2, 2, 1.5, so it's not clear if these are
// measured in FPS or UPS.
var IdleLimit = [3]int{
	4 * FPS, // 0
	4 * FPS, // 1
	3 * FPS, // 2
}
