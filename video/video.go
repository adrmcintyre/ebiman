package video

import (
	"github.com/adrmcintyre/poweraid/data"
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
	Cursor      TilePos                  // current cursor position for adding tiles
	Sprites     [MAX_SPRITES]SpriteState // attributes of each sprite
	SpriteCount int                      // how many sprites are active
	FlashCycle  int                      // control flashing of dots
	FlashOff    bool                     // """
}

func (v *Video) ColorMaze() {
	v.FlashMaze(false)
	for x := 11; x <= 16; x++ {
		v.ColorTile(TilePos{x, 14}, palette.PAL_26)
		v.ColorTile(TilePos{x, 26}, palette.PAL_26)
	}
	for y := 16; y <= 18; y++ {
		for x := 23; x <= 27; x++ {
			v.ColorTile(TilePos{x, y}, palette.TUNNEL)
		}
		for x := 0; x <= 4; x++ {
			v.ColorTile(TilePos{x, y}, palette.TUNNEL)
		}
	}
	v.ColorTile(TilePos{13, 15}, palette.GATE)
	v.ColorTile(TilePos{14, 15}, palette.GATE)
}

func (v *Video) FlashMaze(flash bool) {
	pal := palette.MAZE
	if flash {
		pal = palette.MAZE_FLASH
	}
	for y := 2; y <= 33; y++ {
		for x := range 28 {
			v.ColorTile(TilePos{x, y}, pal)
		}
	}
	for y := range 2 {
		for x := range 32 {
			v.ColorTile(TilePos{x, y}, palette.SCORE)
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

func (v *Video) SetTile(pos TilePos, t byte) {
	v.TileRam[pos.tileIndex()] = t
}

func (v *Video) ColorTile(pos TilePos, pal byte) {
	v.PalRam[pos.tileIndex()] = pal
}

func (v *Video) GetTile(pos TilePos) byte {
	return v.TileRam[pos.tileIndex()]
}

func (v *Video) SetStatusQuad(xPos int, baseTile byte, pal byte) {
	yPos := 34
	tile := baseTile

	for i := range 2 {
		for j := range 2 {
			v.SetTile(TilePos{xPos + 1 - j, yPos + i}, tile)
			v.ColorTile(TilePos{xPos + 1 - j, yPos + i}, pal)
			tile += 1
		}
	}
}

func (v *Video) DrawTiles(screen *ebiten.Image) {
	for ty := 0; ty < 36; ty++ {
		for tx := range 28 {
			pos := TilePos{tx, ty}
			scrPos := pos.ToScreenPos()
			op := colorm.DrawImageOptions{}
			op.GeoM.Translate(float64(hOffset+scrPos.X), float64(vOffset+scrPos.Y))
			op.GeoM.Scale(1, 1)
			index := pos.tileIndex()
			colorm.DrawImage(screen, tile.Image[v.TileRam[index]], palette.ColorM[v.PalRam[index]], &op)
		}
	}
}

func (v *Video) Draw(screen *ebiten.Image) {
	v.DrawTiles(screen)
	v.DrawSprites(screen)
}
