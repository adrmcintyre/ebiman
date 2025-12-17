package palette

type colorByte byte

// Format: bbgg:grrr
//
//	7654 3210
var colorData = [32]colorByte{
	0x00, //  0 - 00:000:000 black
	0x07, //  1 - 00:000:111 red (blinky)
	0x66, //  2 - 01:100:110 brown (apple stalk)
	0xef, //  3 - 11:101:111 pink (pinky)
	0x00, //  4 - unused
	0xf8, //  5 - 11:111:000 cyan (inky)
	0xea, //  6 - 11:101:010 mid-blue-cyan "steel" (key, bell)
	0x6f, //  7 - 01:101:111 orange (clyde)
	0x00, //  8 - unused
	0x3f, //  9 - 00:111:111 yellow (pacman)
	0x00, // 10 - unused
	0xc9, // 11 - 11:001:001 blue (scared ghost)
	0x38, // 12 - 00:111:000 green (leaf)
	0xaa, // 13 - 10:101:010 dark-cyan (pineapple wood)
	0xaf, // 14 - 10:101:111 pill
	0xf6, // 15 - 11:110:110 white

	// entries 16..31 not used
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}
