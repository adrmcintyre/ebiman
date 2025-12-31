package color

import "image/color"

// A colorByte is a bb:ggg:rrr colour triplet.
type colorByte byte

// Colour palettes contain colour indexes corresponding to entries in this table.
// We use octal notation below as it nicely breaks out each colour channel.
var colorData = [32]colorByte{
	0_000, //  0 - black
	0_007, //  1 - red (blinky)
	0_146, //  2 - brown (apple stalk)
	0_357, //  3 - pink (pinky)
	0_000, //
	0_370, //  5 - cyan (inky)
	0_352, //  6 - mid-blue-cyan "steel" (key, bell)
	0_157, //  7 - orange (clyde)
	0_000, //
	0_077, //  9 - yellow (pacman)
	0_000, //
	0_311, // 11 - blue (scared ghost)
	0_070, // 12 - green (leaf)
	0_252, // 13 - dark-cyan (pineapple wood)
	0_257, // 14 - pill
	0_366, // 15 - white

	// entries 16..31 not used
	0000, 0000, 0000, 0000, 0000, 0000, 0000, 0000,
	0000, 0000, 0000, 0000, 0000, 0000, 0000, 0000,
}

// Channel maps each possible value in a 2-bpp bitmap to a single colour channel.
var Channel = []color.Color{
	color.RGBA{},                       // transparent
	color.RGBA{0x00, 0xff, 0x00, 0xff}, // green
	color.RGBA{0xff, 0x00, 0x00, 0xff}, // red
	color.RGBA{0x00, 0x00, 0xff, 0xff}, // blue
}
