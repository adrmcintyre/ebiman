package video

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/tile"
)

var powerPillPos = [4]geom.Position{
	geom.TilePos(1, 6),
	geom.TilePos(26, 6),
	geom.TilePos(1, 26),
	geom.TilePos(26, 26),
}

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
		for _, pos := range powerPillPos {
			index := TileIndex(pos.TileXY())
			if v.TileRam[index] == tile.POWER {
				v.PalRam[index] = pal
			}
		}
	}
}
