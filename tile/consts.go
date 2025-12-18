package tile

type Tile byte

const (
	HEX_BASE Tile = 0x00 // 0x00..0x0f = 0..F
	PILL     Tile = 0x10
	// 0x11 - invisible pill?
	POWER_SMALL Tile = 0x12
	// 0x13 - invisible power-small?
	POWER Tile = 0x14
	// 0x15 - invisible power?
	// 0x16-0x1f - unused
	PACMAN_BASE Tile = 0x20 // 0x20..0x23
	// 0x24 - unused
	POINT      Tile = 0x25
	QUOTES     Tile = 0x26
	QUOTES2    Tile = 0x27 // duplicate
	NAMCO      Tile = 0x28 // 0x28..0x2F
	DIGIT_BASE Tile = 0x30 // 0x30..0x39 = 0..9
	SLASH      Tile = 0x3A
	MINUS      Tile = 0x3B
	// 0x3C..0x3F square angles - unused?
	SPACE      Tile = 0x40
	ALPHA_BASE Tile = 0x41 // 0x41..0x5A = A..Z
	EXCLAM     Tile = 0x5B
	COPYRIGHT  Tile = 0x5C
	PTS        Tile = 0x5D // 0x5D..0x5F
	// 0x6E..0x7F - unused

	// Bonus points tiles
	BONUS_SCORE_MIN Tile = 0x80
	BONUS_SCORE_MAX Tile = 0x8F

	// 0x80 - unused
	// tiles for 100-700 bonus points
	SCORE_100 Tile = 0x81
	SCORE_300 Tile = 0x82
	SCORE_500 Tile = 0x83
	SCORE_700 Tile = 0x84
	SCORE_X00 Tile = 0x85

	// tiles for 1000-5000 bonus points
	SCORE_1000   Tile = 0x86
	SCORE_2000_1 Tile = 0x87
	SCORE_2000_2 Tile = 0x88
	SCORE_3000_1 Tile = 0x89
	SCORE_3000_2 Tile = 0x8A
	SCORE_5000_1 Tile = 0x8B
	SCORE_5000_2 Tile = 0x8C
	SCORE_X000_1 Tile = 0x8D
	SCORE_X000_2 Tile = 0x8E
	// 0x8f - unused

	// Bonus tiles
	CHERRY_BASE     Tile = 0x90
	STRAWBERRY_BASE Tile = 0x94
	ORANGE_BASE     Tile = 0x98
	BELL_BASE       Tile = 0x9C
	APPLE_BASE      Tile = 0xA0
	PINEAPPLE_BASE  Tile = 0xA4
	GALAXIAN_BASE   Tile = 0xA8
	KEY_BASE        Tile = 0xAC

	GHOST_BASE Tile = 0xB0 // 0xB0..0xB5

	SPACE_BASE Tile = 0xB8 // 0xB8..0xBB - not used by pacman - we're stealing this for convenience

	// 0xb6..0xcd - unused

	MAZE_MIN Tile = 0xCE
	MAZE_MAX Tile = 0xFF
	// 0xce..0xff - maze parts
	GATE_LEFT  Tile = 0xCE
	GATE_RIGHT Tile = 0xCF
	HOME_LEFT  Tile = 0xFC
	HOME_RIGHT Tile = 0xFD
)
