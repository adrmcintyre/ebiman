package video

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const hOffset = 8.0
const vOffset = 8.0

type Video struct {
	TileRam     [1024]byte               // tiles
	PalRam      [1024]byte               // per-tile colour palettes
	X, Y        int                      // current cursor position for adding tiles
	Sprites     [MAX_SPRITES]SpriteState // attributes of each sprite
	SpriteCount int                      // how many sprites are active
	FlashCycle  int                      // control flashing of dots
	FlashOff    bool                     // """
}

func tileIndex(tileX, tileY int) int {
	return 0x40 + (27-tileX)*32 + (tileY - 2)
}

func (v *Video) ColorMaze(mode int) {
	pal := palette.MAZE
	if mode == 2 {
		pal = palette.MAZE_FLASH
	}
	for y := 2; y <= 33; y++ {
		for x := range 28 {
			v.ColorTile(x, y, pal)
		}
	}
	for y := range 2 {
		for x := range 32 {
			v.PalRam[0x3c0+y*32+x] = palette.SCORE
		}
	}
	if mode != 0 {
		return
	}
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
	i := 0
	src := 0
	for {
		a := data.Maze[src]
		src++
		if a == 0 {
			break
		}
		if (a & 0x80) == 0 {
			i += int(a) - 1
			a = data.Maze[src]
			src++
		}
		i++
		v.TileRam[i] = a
		v.TileRam[31*32-i+(i&31)*2] = a ^ 1
	}
}

func (v *Video) SetTile(x, y int, t byte) {
	v.TileRam[tileIndex(x, y)] = t
}

func (v *Video) ColorTile(x, y int, pal byte) {
	v.PalRam[tileIndex(x, y)] = pal
}

func (v *Video) GetTile(x, y int) byte {
	return v.TileRam[tileIndex(x, y)]
}

func (v *Video) SetTopTile(x, y int, t byte) {
	if y < 2 && x < 30 {
		off := y*32 + (29 - x) + 2
		v.TileRam[0x3c0+off] = t
	}
}

func (v *Video) SetBottomTile(x, y int, t byte, pal byte) {
	if y < 2 && x < 30 {
		off := y*32 + (29 - x) + 2
		v.TileRam[0x000+off] = t
		v.PalRam[0x000+off] = pal
	}
}

func (v *Video) DrawMaze(screen *ebiten.Image) {
	for ty := range 32 {
		for tx := range 28 {
			screenX := tx * 8
			screenY := ty*8 + 16
			op := colorm.DrawImageOptions{}
			op.GeoM.Translate(hOffset+float64(screenX), vOffset+float64(screenY))
			op.GeoM.Scale(1, 1)
			i := tileIndex(tx, ty+2)
			colorm.DrawImage(screen, tile.Image[v.TileRam[i]], palette.ColorM[v.PalRam[i]], &op)
		}
	}
}

func (v *Video) DrawStatus(screen *ebiten.Image) {
	for section := range 2 {
		for y := range 2 {
			for x := range 30 {
				off := y*32 + (29 - x) + 2
				i := 0x000
				if section == 0 {
					i = 0x3c0
				}
				i += off

				screenX := x * 8
				screenY := (y + section*34) * 8

				op := colorm.DrawImageOptions{}
				op.GeoM.Translate(hOffset+float64(screenX), vOffset+float64(screenY))
				op.GeoM.Scale(1, 1)
				colorm.DrawImage(screen, tile.Image[v.TileRam[i]], palette.ColorM[v.PalRam[i]], &op)
			}
		}
	}
}

func (v *Video) Draw(screen *ebiten.Image) {
	v.DrawMaze(screen)
	v.DrawStatus(screen)
	v.DrawSprites(screen)
}
