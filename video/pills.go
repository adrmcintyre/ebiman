package video

import (
	"github.com/adrmcintyre/ebiman/geom"
)

// FlashPills causes the pills to flash by periodically
// setting their palettes to black.
func (v *Video) FlashPills() {
	const flashFrames = 10

	v.flashCycle += 1
	if v.flashCycle > flashFrames {
		v.flashCycle = 0
		v.flashOff = !v.flashOff

		pal := PalBlack
		if v.flashOff {
			pal = PalMaze
		}
		for _, pos := range geom.PowerPills {
			index := tileIndex(pos.TileXY())
			if v.tileRam[index] == TilePower {
				v.palRam[index] = pal
			}
		}
	}
}
