package color

type colorByte byte

// Format: bbgg:grrr
// Note - entries are in octal
var colorData = [32]colorByte{
	0000, //  0 - 00:000:000 black
	0007, //  1 - 00:000:111 red (blinky)
	0146, //  2 - 01:100:110 brown (apple stalk)
	0357, //  3 - 11:101:111 pink (pinky)
	0000, //  4 - unused
	0370, //  5 - 11:111:000 cyan (inky)
	0352, //  6 - 11:101:010 mid-blue-cyan "steel" (key, bell)
	0157, //  7 - 01:101:111 orange (clyde)
	0000, //  8 - unused
	0077, //  9 - 00:111:111 yellow (pacman)
	0000, // 10 - unused
	0311, // 11 - 11:001:001 blue (scared ghost)
	0070, // 12 - 00:111:000 green (leaf)
	0252, // 13 - 10:101:010 dark-cyan (pineapple wood)
	0257, // 14 - 10:101:111 pill
	0366, // 15 - 11:110:110 white

	// entries 16..31 not used
	0000, 0000, 0000, 0000, 0000, 0000, 0000, 0000,
	0000, 0000, 0000, 0000, 0000, 0000, 0000, 0000,
}
