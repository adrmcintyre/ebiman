package video

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/palette"
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

	v.FlashCycle += 1
	if v.FlashCycle > FLASH_FRAMES {
		v.FlashCycle = 0
		v.FlashOff = !v.FlashOff

		pal := palette.BLACK
		if v.FlashOff {
			pal = palette.MAZE
		}
		for _, pos := range powerPillPos {
			index := TileIndex(pos.TileXY())
			if v.TileRam[index] == tile.POWER {
				v.PalRam[index] = pal
			}
		}
	}
}
