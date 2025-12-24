package pill

import (
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/tile"
	"github.com/adrmcintyre/ebiman/video"
)

// State describes the state of pills and power ups.
type State struct {
	PillBits   [240]bool // true for each uneaten pill
	PowerPills [4]bool   // true for each uneaten power pill
}

// Reset restores the state of each pill and power up to uneaten.
func (ds *State) Reset() {
	for i := range ds.PillBits {
		ds.PillBits[i] = true
	}
	for i := range ds.PowerPills {
		ds.PowerPills[i] = true
	}
}

// Save retrieves the current state of each pill and power up
// from the screen's tile data.
func (ds *State) Save(v *video.Video) {
	pillIndex := 0
	tileIndex := 0

	// FIXME peeking directly into TileRam - not very nice
	for i := range ds.PillBits {
		tileIndex += int(pillData[pillIndex])
		pillIndex += 1
		ds.PillBits[i] = v.TileRam[tileIndex] == tile.PILL
	}

	for i, pos := range geom.POWER_PILLS {
		t := v.GetTile(pos.TileXY())
		ds.PowerPills[i] = t == tile.POWER || t == tile.POWER_SMALL
	}
}

// Draw places tiles representing the state of each pill and power up.
func (ds *State) Draw(v *video.Video) {
	pillIndex := 0
	tileIndex := 0

	// FIXME poking directly into TileRam - not very nice
	for _, bit := range ds.PillBits {
		tileIndex += int(pillData[pillIndex])
		pillIndex++
		if bit {
			v.TileRam[tileIndex] = tile.PILL
		} else {
			v.TileRam[tileIndex] = tile.SPACE
		}
	}

	for i, bit := range ds.PowerPills {
		x, y := geom.POWER_PILLS[i].TileXY()
		if bit {
			v.SetTile(x, y, tile.POWER)
		} else {
			v.SetTile(x, y, tile.SPACE)

		}
	}
}
