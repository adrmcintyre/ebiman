package tile

// A Tile identifies a specific tile bitmap.
type Tile byte

// The identifiers for each tile.
const (
	// base of digits 0..9, A..F - usable for in-display hex debug info
	HexBase Tile = 0x00 // 0x00..0x0F

	// maze dots
	Pill       Tile = 0x10 // a pill dot
	PowerSmall Tile = 0x12 // small power-up (unused)
	Power      Tile = 0x14 // standard power-up
	Point      Tile = 0x25 // '.'

	// punctuation
	Quotes    Tile = 0x26 // '"'
	Quotes2   Tile = 0x27 // '"' (dup)
	Slash     Tile = 0x3A // '/'
	Minus     Tile = 0x3B // '-'
	Space     Tile = 0x40 // ' '
	Exclam    Tile = 0x5B // '!'
	Copyright Tile = 0x5C // (c)
	Pts       Tile = 0x5D // 0x5D..0x5F

	// 4 consecutive tiles for a pacman life in lower status area
	PacmanBase Tile = 0x20 // 0x20..0x23

	// 8 consecutive tiles forming "Namco"
	Namco Tile = 0x28 // 0x28..0x2F

	// base of digits 0..9 (coincides with ASCII values)
	DigitBase Tile = 0x30 // 0x30..0x39

	// base of uppercase alpha A..Z (coincides with ASCII values)
	AlphaBase Tile = 0x41 // 0x41..0x5A

	// tiles between ScoreMin and ScoreMax are used for
	// displaying bonus and ghost scores, and need to be
	// combined in various ways to form the correct text
	ScoreMin Tile = 0x80
	ScoreMax Tile = 0x8F

	// 0x80 - unused
	Score100 Tile = 0x81 // leading text of 100
	Score300 Tile = 0x82 // leading text of 300
	Score500 Tile = 0x83 // leading text of 500
	Score700 Tile = 0x84 // leading text of 700
	ScoreX00 Tile = 0x85 // trailing text of a hundreds score

	// leading text of 1000
	Score1000 Tile = 0x86

	// leading text of 2000
	Score2000_1 Tile = 0x87
	Score2000_2 Tile = 0x88

	// leading text of 3000
	Score3000_1 Tile = 0x89
	Score3000_2 Tile = 0x8A

	// leading text of 5000
	Score5000_1 Tile = 0x8B
	Score5000_2 Tile = 0x8C

	// trailing text of a thousands score
	ScoreX000_1 Tile = 0x8D
	ScoreX000_2 Tile = 0x8E

	// groups of 4 consecutive tiles forming bonuses for display in lower status area
	CherryBase     Tile = 0x90 // 0x90..0x93
	StrawberryBase Tile = 0x94 // 0x94..0x97
	OrangeBase     Tile = 0x98 // 0x98..0x9B
	BellBase       Tile = 0x9C // 0x9C..0x9F
	AppleBase      Tile = 0xA0 // 0xA0..0xA3
	PineappleBase  Tile = 0xA4 // 0xA4..0xA7
	GalaxianBase   Tile = 0xA8 // 0xA8..0xAB
	KeyBase        Tile = 0xAC // 0xAC..0xAF

	// 6 consecutive tiles for a ghost, part of a cut-scene animation
	GhostBase Tile = 0xB0 // 0xB0..0xB5

	// a series of 4 consecutive blank tiles, for erasing items in lower status area
	SpaceBase Tile = 0xB8 // 0xB8..0xBB - not used by pacman

	// tiles between MazeMin and MazeMax are used to draw the maze
	MazeMin   Tile = 0xCE
	GateLeft  Tile = 0xCE // left half of the gate to the ghost's home
	GateRight Tile = 0xCF // right half of the gate to the ghost's home
	HomeLeft  Tile = 0xFC // left part of interior of ghost's home
	HomeRight Tile = 0xFD // right part of interior of ghost's home
	MazeMax   Tile = 0xFF

	PillMinus  Tile = 0x6E
	PillPlus   Tile = 0x6F
	PillMinus2 Tile = 0x70
	PillPlus2  Tile = 0x71

	// Following tiles are unused:
	// 0x11       - PillInvisible?
	// 0x13       - PowerSmallInvisible?
	// 0x15       - PowerInvisible?
	// 0x16-0x1F  - unused
	// 0x24       - unused
	// 0x3C..0x3F - square angles
	// 0x6E..0x7F - unused
	// 0x8F       - unused
	// 0xB6..0xB7 - unused
	// 0xBC..0xCD - unused
)
