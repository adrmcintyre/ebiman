package video

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
	hOffset int = 8
	vOffset int = 8
)

type Video struct {
	TileRam     [1024]byte               // tiles
	PalRam      [1024]byte               // per-tile colour palettes
	CursorX     int                      // current cursor position for adding tiles
	CursorY     int                      // current cursor position for adding tiles
	Sprites     [MAX_SPRITES]SpriteState // attributes of each sprite
	SpriteCount int                      // how many sprites are active
	FlashCycle  int                      // control flashing of dots
	FlashOff    bool                     // """
}

func (v *Video) ColorMaze() {
	v.FlashMaze(false)
	for x := 11; x <= 16; x++ {
		v.ColorTile(x, 14, palette.PAL_26)
		v.ColorTile(x, 26, palette.PAL_26)
	}
	for y := 16; y <= 18; y++ {
		for x := 23; x <= 27; x++ {
			v.ColorTile(x, y, palette.TUNNEL)
		}
		for x := 0; x <= 4; x++ {
			v.ColorTile(x, y, palette.TUNNEL)
		}
	}
	v.ColorTile(13, 15, palette.GATE)
	v.ColorTile(14, 15, palette.GATE)
}

func (v *Video) FlashMaze(flash bool) {
	pal := palette.MAZE
	if flash {
		pal = palette.MAZE_FLASH
	}
	for y := 2; y <= 33; y++ {
		for x := range 28 {
			v.ColorTile(x, y, pal)
		}
	}
	for y := range 2 {
		for x := range 32 {
			v.ColorTile(x, y, palette.SCORE)
		}
	}
}

func (v *Video) ClearTiles() {
	for i := range 1024 {
		v.TileRam[i] = tile.SPACE
	}
}

func (v *Video) ClearPalette() {
	for i := range 1024 {
		v.PalRam[i] = 0
	}
}

func (v *Video) DecodeTiles() {
	index := 0
	for _, op := range data.Maze {
		if (op & 0x80) == 0 {
			index += int(op) - 1
			continue
		}
		index++
		tile := op
		v.TileRam[index] = tile
		mirrorIndex := 31*32 - index + (index&31)*2
		mirrorTile := tile ^ 1
		v.TileRam[mirrorIndex] = mirrorTile
	}
}

// Tile co-ord conversion:
//
//	top (0 <= x < 32, y < 2) - note x=0,1,30,31 are invisible
//	index := (y+30)*32 + (31-x)				// 0x3c0-0x3ff
//
//	normal (0 <= x < 28, 2 <= y < 34)
//	index := (29-x)*32 + (y - 2) 	// 0x040-0x3bf
//
//	bottom (0 <= x < 32, y < 2) - note x=0,1,30,31 are invisible
//	index := y*32 + (31-x) 					// 0x000-0x03f
func TileIndex(x int, y int) int {
	switch {
	case y < 2:
		return 0x3c0 + y*32 + 31 - x
	case y >= 34:
		return 0x000 + (y-34)*32 + 31 - x
	default:
		return 0x40 + (27-x)*32 + (y - 2)
	}
}

func (v *Video) SetTile(x, y int, t byte) {
	v.TileRam[TileIndex(x, y)] = t
}

func (v *Video) ColorTile(x, y int, pal byte) {
	v.PalRam[TileIndex(x, y)] = pal
}

func (v *Video) GetTile(x, y int) byte {
	return v.TileRam[TileIndex(x, y)]
}

func (v *Video) SetStatusQuad(baseX int, baseTile byte, pal byte) {
	baseY := 34
	tile := baseTile

	for i := range 2 {
		for j := range 2 {
			x, y := baseX+1-j, baseY+i
			v.SetTile(x, y, tile)
			v.ColorTile(x, y, pal)
			tile += 1
		}
	}
}

func (v *Video) DrawTiles(screen *ebiten.Image) {
	for ty := 0; ty < 36; ty++ {
		for tx := range 28 {
			pos := geom.TilePos(tx, ty)
			op := colorm.DrawImageOptions{}
			op.GeoM.Translate(float64(hOffset+pos.X), float64(vOffset+pos.Y))
			op.GeoM.Scale(1, 1)
			index := TileIndex(tx, ty)
			colorm.DrawImage(screen, tile.Image[v.TileRam[index]], palette.ColorM[v.PalRam[index]], &op)
		}
	}
}

func (v *Video) Draw(screen *ebiten.Image) {
	v.DrawTiles(screen)
	v.DrawSprites(screen)
}
