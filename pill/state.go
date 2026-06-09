package pill

import (
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/video"
)

// State describes the state of pills and power ups.
type State struct {
	PillTiles  [240]video.Tile // true for each uneaten pill
	PowerPills [4]bool         // true for each uneaten power pill
	NetCharge  int             // total charge of all pills
}

// Reset restores the state of each pill and power up to uneaten.
func (ds *State) Reset() {
	for i := range ds.PillTiles {
		ds.PillTiles[i] = video.TilePill
	}
	for i := range ds.PowerPills {
		ds.PowerPills[i] = true
	}
	ds.cacheNetCharge()
}

// Save retrieves the current state of each pill and power up
// from the screen's tile data.
func (ds *State) Save(v *video.Video) {
	tileIndex := 0
	for i := range ds.PillTiles {
		tileIndex += int(pillData[i])
		t := v.GetTileAtIndex(tileIndex)
		if t.IsPill() {
			t = video.TilePill
		}
		ds.PillTiles[i] = t
	}

	for i, pos := range geom.PowerPills {
		t := v.GetTile(pos.TileXY())
		ds.PowerPills[i] = t == video.TilePower || t == video.TilePowerSmall
	}

	ds.cacheNetCharge()
}

// cacheNetCharge recalculated the NetCharge property.
func (ds *State) cacheNetCharge() {
	ds.NetCharge = 0
	for _, t := range ds.PillTiles {
		ds.NetCharge += t.Charge()
	}
}

// Draw places tiles representing the state of each pill and power up.
func (ds *State) Draw(v *video.Video) {

	tileIndex := 0
	for i, t := range ds.PillTiles {
		tileIndex += int(pillData[i])
		v.SetTileAtIndex(tileIndex, t)
	}

	for i, bit := range ds.PowerPills {
		x, y := geom.PowerPills[i].TileXY()
		if bit {
			v.SetTile(x, y, video.TilePower)
		} else {
			v.SetTile(x, y, video.TileSpace)
		}
	}
}
