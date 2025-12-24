package pill

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

// State describes the state of pills and power ups.
type State struct {
	PillTiles  [240]tile.Tile // true for each uneaten pill
	PowerPills [4]bool        // true for each uneaten power pill
}

// Reset restores the state of each pill and power up to uneaten.
func (ds *State) Reset() {
	for i := range ds.PillTiles {
		ds.PillTiles[i] = tile.PILL
		//if i&4 == 0 {
		//	ds.PillTiles[i] = tile.PILL_PLUS
		//} else {
		//	ds.PillTiles[i] = tile.PILL_MINUS
		//}
	}
	for i := range ds.PowerPills {
		ds.PowerPills[i] = true
	}
}

// Save retrieves the current state of each pill and power up
// from the screen's tile data.
func (ds *State) Save(v *video.Video) {

	// FIXME peeking directly into TileRam - not very nice
	tileIndex := 0
	for i := range ds.PillTiles {
		tileIndex += int(pillData[i])
		ds.PillTiles[i] = v.TileRam[tileIndex]
	}

	for i, pos := range geom.POWER_PILLS {
		t := v.GetTile(pos.TileXY())
		ds.PowerPills[i] = t == tile.POWER || t == tile.POWER_SMALL
	}
}

// Draw places tiles representing the state of each pill and power up.
func (ds *State) Draw(v *video.Video) {

	// FIXME poking directly into TileRam - not very nice
	tileIndex := 0
	for i, t := range ds.PillTiles {
		tileIndex += int(pillData[i])
		v.TileRam[tileIndex] = t
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
