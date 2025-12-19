package pill

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

type State struct {
	PillBits   [240]bool    // bitmap of uneaten pills
	PowerPills [4]tile.Tile // tile at each power pill location
}

func (ds *State) Reset() {
	for i := range ds.PillBits {
		ds.PillBits[i] = true
	}
	for i := range ds.PowerPills {
		ds.PowerPills[i] = tile.POWER
	}
}

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
		ds.PowerPills[i] = v.GetTile(pos.TileXY())
	}
}

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

	for i, pos := range geom.POWER_PILLS {
		x, y := pos.TileXY()
		v.SetTile(x, y, ds.PowerPills[i])
	}
}
