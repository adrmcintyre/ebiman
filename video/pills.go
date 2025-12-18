package video

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/tile"
)

func (v *Video) FlashPills() {
	const FLASH_FRAMES = 10

	v.flashCycle += 1
	if v.flashCycle > FLASH_FRAMES {
		v.flashCycle = 0
		v.flashOff = !v.flashOff

		pal := color.PAL_BLACK
		if v.flashOff {
			pal = color.PAL_MAZE
		}
		for _, pos := range geom.POWER_PILLS {
			index := tileIndex(pos.TileXY())
			if v.TileRam[index] == tile.POWER {
				v.palRam[index] = pal
			}
		}
	}
}
