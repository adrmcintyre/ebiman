package video

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/tile"
)

// FlashPills causes the pills to flash by periodically
// setting their palettes to black.
func (v *Video) FlashPills() {
	const flashFrames = 10

	v.flashCycle += 1
	if v.flashCycle > flashFrames {
		v.flashCycle = 0
		v.flashOff = !v.flashOff

		pal := color.PalBlack
		if v.flashOff {
			pal = color.PalMaze
		}
		for _, pos := range geom.PowerPills {
			index := tileIndex(pos.TileXY())
			if v.TileRam[index] == tile.Power {
				v.palRam[index] = pal
			}
		}
	}
}
