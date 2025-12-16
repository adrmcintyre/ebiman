package data

import (
	"github.com/adrmcintyre/poweraid/sprite"
)

const (
	FPS = 60 // frames-per-sec
	UPS = 30 // updates-per-sec
)

// The actors advance on each update according to a pulse train.
// The more bits set, the more updates will cause it to advance,
// and so the faster the actor moves.
type PCM uint32

// The bits in each pulse train are spread out as evenly as possible to reduce
// jerkiness. Each constant is named as a percentage of pacman's maximum speed.
const (
	PCM_5   PCM = 0x20000000 // 1 pulse every 32 updates
	PCM_10  PCM = 0x20002000 // 2 etc
	PCM_15  PCM = 0x20040100 // 3
	PCM_20  PCM = 0x20202020 // 4
	PCM_25  PCM = 0x20810408 // 5
	PCM_30  PCM = 0x20842084 // 6
	PCM_35  PCM = 0x22110884 // 7
	PCM_40  PCM = 0x22222222 // 8
	PCM_45  PCM = 0x24489122 // 9
	PCM_50  PCM = 0x24922492 // 10
	PCM_55  PCM = 0x24924925 // 11
	PCM_60  PCM = 0x25252525 // 12
	PCM_65  PCM = 0x25A4925A // 13
	PCM_70  PCM = 0x259A259A // 14
	PCM_75  PCM = 0x2AAA5555 // 15
	PCM_80  PCM = 0x55555555 // 16
	PCM_85  PCM = 0x6AAAD555 // 17
	PCM_90  PCM = 0x6AD56AD5 // 18
	PCM_95  PCM = 0x5AD6B5AD // 19
	PCM_100 PCM = 0x6D6D6D6D // 20 - pacman's maximum speed
	PCM_105 PCM = 0x6DB6DB6D // 21
	PCM_110 PCM = 0x6DBB6DBB // 22
	PCM_MAX PCM = 0xFFFFFFFF // eyes return at full pelt
)

func (pcm *PCM) Pulse() bool {
	msb := *pcm >> 31
	*pcm = (*pcm << 1) | msb
	return msb != 0
}

type LevelEntry struct {
	SpeedIndex    int
	unused        int
	DotLimitIndex int
	ElroyIndex    int
	BlueIndex     int
	IdleIndex     int
}

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

type Speeds struct {
	Pacman     PCM
	PacmanBlue PCM
	Elroy2     PCM
	Elroy1     PCM
	Ghost      PCM
	GhostBlue  PCM
	Tunnel     PCM
}

type SwitchTacticsEntry [7]int

type SpeedDataEntry struct {
	Easy          Speeds
	Medium        Speeds
	Hard          Speeds
	SwitchTactics SwitchTacticsEntry
}

var SpeedData = [4]SpeedDataEntry{
	// Indexes are offset by +3

	// 3 - level 1
	{
		Speeds{PCM_90, PCM_95, PCM_80, PCM_75, PCM_70, PCM_45, PCM_40},
		Speeds{PCM_80, PCM_90, PCM_85, PCM_80, PCM_75, PCM_50, PCM_40},
		Speeds{PCM_80, PCM_90, PCM_90, PCM_85, PCM_80, PCM_55, PCM_45},
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
		Speeds{PCM_95, PCM_100, PCM_90, PCM_85, PCM_80, PCM_50, PCM_40},
		Speeds{PCM_90, PCM_95, PCM_95, PCM_90, PCM_85, PCM_55, PCM_45},
		Speeds{PCM_90, PCM_95, PCM_100, PCM_95, PCM_90, PCM_60, PCM_50},
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
		Speeds{PCM_105, PCM_105, PCM_100, PCM_95, PCM_90, PCM_55, PCM_45},
		Speeds{PCM_100, PCM_100, PCM_105, PCM_100, PCM_95, PCM_60, PCM_50},
		Speeds{PCM_100, PCM_100, PCM_110, PCM_105, PCM_100, PCM_65, PCM_55},
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
		Speeds{PCM_95, 1, PCM_105, PCM_95, PCM_90, 0, PCM_45},
		Speeds{PCM_90, 0, PCM_105, PCM_100, PCM_95, 0, PCM_50},
		Speeds{PCM_90, 0, PCM_110, PCM_105, PCM_100, 0, PCM_55},
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

// when ghost.dot_counter reaches this value, ghost goes out of home
type DotLimitEntry struct {
	Pinky int
	Inky  int
	Clyde int
}

var DotLimit = [4]DotLimitEntry{
	{20, 30, 70}, // 0 - this entry appears to be unused
	{0, 30, 60},  // 1
	{0, 0, 50},   // 2
	{0, 0, 0},    // 3
}

// remaining number of pills when first difficulty flag is set
type ElroyEntry struct {
	Pills1 int // when first difficulty flag set (cruise elroy 1)
	Pills2 int // when second difficulty flag  set (cruise elroy 2)
}

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

// Time the ghosts stay blue when pacman eats a big pill
type BlueControlEntry struct {
	BlueTime       int // total time to remain blue (including flashes)
	WhiteBlueCount int // number of white/blue flashes
}

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

// Number of units before a ghost goes out of home when pacman is idle.
// Pacman Dossier claims these are 4, 4, 3, whereas comments in mspacman.asm claim
// they are 2, 2, 1.5, so it's not clear if these are measured in FPS or UPS.
var IdleLimit = [3]int{
	4 * FPS, // 0
	4 * FPS, // 1
	3 * FPS, // 2
}

const (
	EXTRA_LIFE_SCORE = 10000
	DOT_SCORE        = 10 // value of a dot
	POWER_SCORE      = 50 // value of a power pill
)

const (
	DOT_STALL   = 1 // how long pacman stalls after eating a dot
	POWER_STALL = 4 // how long pacman stalls after eating a power pill
)

const DISPLAY_GHOST_SCORE_MS = 1000 // how long to display ghost's score

const WHITE_BLUE_PERIOD = 14 // number of updates between white and blue

// points awarded for consecutive ghosts
type GhostScoreEntry struct {
	Score  int  // points to award
	Sprite byte // sprite to display
}

var GhostScore = [4]GhostScoreEntry{
	{200, sprite.SCORE_200},
	{400, sprite.SCORE_400},
	{800, sprite.SCORE_800},
	{1600, sprite.SCORE_1600},
}
