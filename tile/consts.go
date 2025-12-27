package tile

// A Tile identifies a specific tile bitmap.
type Tile byte

// The identifiers for each tile.
const (
	// base of digits 0..9, A..F - usable for in-display hex debug info
	HEX_BASE Tile = 0x00 // 0x00..0x0F

	// maze dots
	PILL        Tile = 0x10 // a pill dot
	POWER_SMALL Tile = 0x12 // small power-up (unused)
	POWER       Tile = 0x14 // standard power-up
	POINT       Tile = 0x25 // '.'

	// punctuation
	QUOTES    Tile = 0x26 // '"'
	QUOTES2   Tile = 0x27 // '"' (dup)
	SLASH     Tile = 0x3A // '/'
	MINUS     Tile = 0x3B // '-'
	SPACE     Tile = 0x40 // ' '
	EXCLAM    Tile = 0x5B // '!'
	COPYRIGHT Tile = 0x5C // (c)
	PTS       Tile = 0x5D // 0x5D..0x5F

	// 4 consecutive tiles for a pacman life in lower status area
	PACMAN_BASE Tile = 0x20 // 0x20..0x23

	// 8 consecutive tiles forming "NAMCO"
	NAMCO Tile = 0x28 // 0x28..0x2F

	// base of digits 0..9 (coincides with ASCII values)
	DIGIT_BASE Tile = 0x30 // 0x30..0x39

	// base of uppercase alpha A..Z (coincides with ASCII values)
	ALPHA_BASE Tile = 0x41 // 0x41..0x5A

	// tiles between SCORE_MIN and SCORE_MAX are used for
	// displaying bonus and ghost scores, and need to be
	// combined in various ways to form the correct text
	SCORE_MIN Tile = 0x80
	SCORE_MAX Tile = 0x8F

	// 0x80 - unused
	SCORE_100 Tile = 0x81 // leading text of 100
	SCORE_300 Tile = 0x82 // leading text of 300
	SCORE_500 Tile = 0x83 // leading text of 500
	SCORE_700 Tile = 0x84 // leading text of 700
	SCORE_X00 Tile = 0x85 // trailing text of a hundreds score

	// leading text of 1000
	SCORE_1000 Tile = 0x86

	// leading text of 2000
	SCORE_2000_1 Tile = 0x87
	SCORE_2000_2 Tile = 0x88

	// leading text of 3000
	SCORE_3000_1 Tile = 0x89
	SCORE_3000_2 Tile = 0x8A

	// leading text of 5000
	SCORE_5000_1 Tile = 0x8B
	SCORE_5000_2 Tile = 0x8C

	// trailing text of a thousands score
	SCORE_X000_1 Tile = 0x8D
	SCORE_X000_2 Tile = 0x8E

	// groups of 4 consecutive tiles forming bonuses for display in lower status area
	CHERRY_BASE     Tile = 0x90 // 0x90..0x93
	STRAWBERRY_BASE Tile = 0x94 // 0x94..0x97
	ORANGE_BASE     Tile = 0x98 // 0x98..0x9B
	BELL_BASE       Tile = 0x9C // 0x9C..0x9F
	APPLE_BASE      Tile = 0xA0 // 0xA0..0xA3
	PINEAPPLE_BASE  Tile = 0xA4 // 0xA4..0xA7
	GALAXIAN_BASE   Tile = 0xA8 // 0xA8..0xAB
	KEY_BASE        Tile = 0xAC // 0xAC..0xAF

	// 6 consecutive tiles for a ghost, part of a cut-scene animation
	GHOST_BASE Tile = 0xB0 // 0xB0..0xB5

	// a series of 4 consecutive blank tiles, for erasing items in lower status area
	SPACE_BASE Tile = 0xB8 // 0xB8..0xBB - not used by pacman

	// tiles between MAZE_MIN and MAZE_MAX are used to draw the maze
	MAZE_MIN   Tile = 0xCE
	GATE_LEFT  Tile = 0xCE // left half of the gate to the ghost's home
	GATE_RIGHT Tile = 0xCF // right half of the gate to the ghost's home
	HOME_LEFT  Tile = 0xFC // left part of interior of ghost's home
	HOME_RIGHT Tile = 0xFD // right part of interior of ghost's home
	MAZE_MAX   Tile = 0xFF

	PILL_MINUS  Tile = 0x6E
	PILL_PLUS   Tile = 0x6F
	PILL_MINUS2 Tile = 0x70
	PILL_PLUS2  Tile = 0x71

	// Following tiles are unused:
	// 0x11       - PILL_INVISIBLE?
	// 0x13       - POWER_SMALL_INVISIBLE?
	// 0x15       - POWER_INVISIBLE?
	// 0x16-0x1F  - unused
	// 0x24       - unused
	// 0x3C..0x3F - square angles
	// 0x6E..0x7F - unused
	// 0x8F       - unused
	// 0xB6..0xB7 - unused
	// 0xBC..0xCD - unused
)
