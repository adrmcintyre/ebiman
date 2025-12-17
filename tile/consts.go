package tile

const (
	HEX_BASE byte = 0x00 // 0x00..0x0f = 0..F
	PILL     byte = 0x10
	// 0x11 - invisible pill?
	POWER_SMALL byte = 0x12
	// 0x13 - invisible power-small?
	POWER byte = 0x14
	// 0x15 - invisible power?
	// 0x16-0x1f - unused
	PACMAN_BASE byte = 0x20 // 0x20..0x23
	// 0x24 - unused
	POINT      byte = 0x25
	QUOTES     byte = 0x26
	QUOTES2    byte = 0x27 // duplicate
	NAMCO      byte = 0x28 // 0x28..0x2F
	DIGIT_BASE byte = 0x30 // 0x30..0x39 = 0..9
	SLASH      byte = 0x3A
	MINUS      byte = 0x3B
	// 0x3C..0x3F square angles - unused?
	SPACE      byte = 0x40
	ALPHA_BASE byte = 0x41 // 0x41..0x5A = A..Z
	EXCLAM     byte = 0x5B
	COPYRIGHT  byte = 0x5C
	PTS        byte = 0x5D // 0x5D..0x5F
	// 0x6E..0x7F - unused

	// Bonus points tiles
	BONUS_SCORE_MIN byte = 0x80
	BONUS_SCORE_MAX byte = 0x8F

	// 0x80 - unused
	// tiles for 100-700 bonus points
	SCORE_100 byte = 0x81
	SCORE_300 byte = 0x82
	SCORE_500 byte = 0x83
	SCORE_700 byte = 0x84
	SCORE_X00 byte = 0x85

	// tiles for 1000-5000 bonus points
	SCORE_1000   byte = 0x86
	SCORE_2000_1 byte = 0x87
	SCORE_2000_2 byte = 0x88
	SCORE_3000_1 byte = 0x89
	SCORE_3000_2 byte = 0x8A
	SCORE_5000_1 byte = 0x8B
	SCORE_5000_2 byte = 0x8C
	SCORE_X000_1 byte = 0x8D
	SCORE_X000_2 byte = 0x8E
	// 0x8f - unused

	// Bonus tiles
	CHERRY_BASE     byte = 0x90
	STRAWBERRY_BASE byte = 0x94
	ORANGE_BASE     byte = 0x98
	BELL_BASE       byte = 0x9C
	APPLE_BASE      byte = 0xA0
	PINEAPPLE_BASE  byte = 0xA4
	GALAXIAN_BASE   byte = 0xA8
	KEY_BASE        byte = 0xAC

	GHOST_BASE byte = 0xB0 // 0xB0..0xB5

	SPACE_BASE byte = 0xB8 // 0xB8..0xBB - not used by pacman - we're stealing this for convenience

	// 0xb6..0xcd - unused

	MAZE_MIN byte = 0xCE
	MAZE_MAX byte = 0xFF
	// 0xce..0xff - maze parts
	GATE_LEFT  byte = 0xCE
	GATE_RIGHT byte = 0xCF
	HOME_LEFT  byte = 0xFC
	HOME_RIGHT byte = 0xFD
)
