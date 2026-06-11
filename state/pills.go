package state

import (
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/video"
)

// Pills describes the state of pills and power ups.
type Pills struct {
	PillTiles  [240]video.Tile // true for each uneaten pill
	PowerPills [4]bool         // true for each uneaten power pill
	NetCharge  int             // total charge of all pills
}

// Reset restores the state of each pill and power up to uneaten.
func (ps *Pills) Reset() {
	for i := range ps.PillTiles {
		ps.PillTiles[i] = video.TilePill
	}
	for i := range ps.PowerPills {
		ps.PowerPills[i] = true
	}
	ps.cacheNetCharge()
}

// Save retrieves the current state of each pill and power up
// from the screen's tile data.
func (ps *Pills) Save(v *video.Video) {
	tileIndex := 0
	for i := range ps.PillTiles {
		tileIndex += int(data.PillData[i])
		t := v.GetTileAtIndex(tileIndex)
		if t.IsPill() {
			t = video.TilePill
		}
		ps.PillTiles[i] = t
	}

	for i, pos := range geom.PowerPills {
		t := v.GetTile(pos.TileXY())
		ps.PowerPills[i] = t == video.TilePower || t == video.TilePowerSmall
	}

	ps.cacheNetCharge()
}

// cacheNetCharge recalculated the NetCharge property.
func (ps *Pills) cacheNetCharge() {
	ps.NetCharge = 0
	for _, t := range ps.PillTiles {
		ps.NetCharge += t.Charge()
	}
}

// Draw places tiles representing the state of each pill and power up.
func (ps *Pills) Draw(v *video.Video) {

	tileIndex := 0
	for i, t := range ps.PillTiles {
		tileIndex += int(data.PillData[i])
		v.SetTileAtIndex(tileIndex, t)
	}

	for i, bit := range ps.PowerPills {
		x, y := geom.PowerPills[i].TileXY()
		if bit {
			v.SetTile(x, y, video.TilePower)
		} else {
			v.SetTile(x, y, video.TileSpace)
		}
	}
}
