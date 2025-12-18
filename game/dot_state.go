package game

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

type DotState struct {
	PillBits   [30]byte     // bitmap of uneaten pills
	PowerPills [4]tile.Tile // tile at each power pill location
}

func (ds *DotState) ResetPellets() {
	for i := range ds.PillBits {
		ds.PillBits[i] = 0xff
	}
	for i := range ds.PowerPills {
		ds.PowerPills[i] = tile.POWER
	}
}

func (ds *DotState) DrawPellets(v *video.Video) {
	pillIndex := 0
	tileIndex := 0

	// FIXME poking directly into TileRam - not very nice
	for _, bits := range ds.PillBits {
		for mask := byte(0x80); mask > 0; mask >>= 1 {
			tileIndex += int(data.Pill[pillIndex])
			pillIndex++
			if bits&mask != 0 {
				v.TileRam[tileIndex] = tile.PILL
			} else {
				v.TileRam[tileIndex] = tile.SPACE
			}
		}
	}
	// TODO derive tile X,Y coords instead
	v.TileRam[3*32+4] = ds.PowerPills[0]
	v.TileRam[3*32+24] = ds.PowerPills[1]
	v.TileRam[28*32+4] = ds.PowerPills[2]
	v.TileRam[28*32+24] = ds.PowerPills[3]
}

func (ds *DotState) SavePellets(v *video.Video) {
	pillIndex := 0
	tileIndex := 0

	// FIXME peeking directly into TileRam - not very nice
	for i := range ds.PillBits {
		bits := byte(0)
		for mask := byte(0x80); mask != 0; mask >>= 1 {
			tileIndex += int(data.Pill[pillIndex])
			pillIndex += 1
			if v.TileRam[tileIndex] == tile.PILL {
				bits |= mask
			}
		}
		ds.PillBits[i] = bits
	}

	// TODO derive tile X,Y coords instead
	ds.PowerPills[0] = v.TileRam[3*32+4]
	ds.PowerPills[1] = v.TileRam[3*32+24]
	ds.PowerPills[2] = v.TileRam[28*32+4]
	ds.PowerPills[3] = v.TileRam[28*32+24]
}
