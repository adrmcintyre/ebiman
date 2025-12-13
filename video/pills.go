package video

import (
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
)

var pillCoords = [4]struct{ X, Y int }{
	{1, 6},
	{26, 6},
	{1, 26},
	{26, 26},
}

func (v *Video) FlashPills() {
	const FLASH_FRAMES = 10

	v.FlashCycle += 1
	if v.FlashCycle > FLASH_FRAMES {
		v.FlashCycle = 0
		v.FlashOff = !v.FlashOff

		pal := palette.BLACK
		if v.FlashOff {
			pal = palette.MAZE
		}
		for _, coords := range pillCoords {
			index := tileIndex(coords.X, coords.Y)
			if v.TileRam[index] == tile.POWER {
				v.PalRam[index] = pal
			}
		}
	}
}
