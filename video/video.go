package video

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/tile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
	hOffset int = 8
	vOffset int = 8
)

// Video abstracts the video hardware.
type Video struct {
	TileRam     [1024]tile.Tile         // tiles
	palRam      [1024]color.Palette     // per-tile colour palettes
	cursorX     int                     // current cursor position for adding tiles
	cursorY     int                     // current cursor position for adding tiles
	sprites     [maxSprites]spriteState // attributes of each sprite
	spriteCount int                     // how many sprites are active
	flashCycle  int                     // control flashing of dots
	flashOff    bool                    // """
	shader      *ebiten.Shader          // shader for output filtering
	offsetX     int
	offsetY     int
}

func (v *Video) SetOffset(x int, y int) {
	v.offsetX = x
	v.offsetY = y
}

// ColorMaze establishes the proper colour palettes for the maze area of the screen.
func (v *Video) ColorMaze() {
	v.FlashMaze(false)
	for x := 11; x <= 16; x++ {
		v.ColorTile(x, 14, color.PAL_26)
		v.ColorTile(x, 26, color.PAL_26)
	}
	for y := 16; y <= 18; y++ {
		for x := 23; x <= 27; x++ {
			v.ColorTile(x, y, color.PAL_TUNNEL)
		}
		for x := 0; x <= 4; x++ {
			v.ColorTile(x, y, color.PAL_TUNNEL)
		}
	}
	v.ColorTile(13, 15, color.PAL_GATE)
	v.ColorTile(14, 15, color.PAL_GATE)
}

// FlashMaze switches the maze colour palettes to/from an alternate bright version.
// This is used for signalling the end of a level.
func (v *Video) FlashMaze(flash bool) {
	pal := color.PAL_MAZE
	if flash {
		pal = color.PAL_MAZE_FLASH
	}
	for y := 2; y <= 33; y++ {
		for x := range 28 {
			v.ColorTile(x, y, pal)
		}
	}
	for y := range 2 {
		for x := range 32 {
			v.ColorTile(x, y, color.PAL_SCORE)
		}
	}
}

// ClearTiles replaces all tiles in the display with a blank SPACE tile.
func (v *Video) ClearTiles() {
	for i := range 1024 {
		v.TileRam[i] = tile.SPACE
	}
}

// ClearPalettes replaces all the tile palettes with black.
func (v *Video) ClearPalette() {
	for i := range 1024 {
		v.palRam[i] = color.PAL_BLACK
	}
}

// DecodeTiles sets up all of the maze tiles from an encoded representation.
// The data only contains the left hand part of the maze; the right hard part
// is inferred by mirroring the layout and placing mirror image tiles.
func (v *Video) DecodeTiles() {
	index := 0
	for _, op := range data.Maze {
		if (op & 0x80) == 0 {
			index += int(op) - 1
			continue
		}
		index++
		mirrorIndex := 31*32 - index + (index&31)*2
		v.TileRam[index] = tile.Tile(op)
		v.TileRam[mirrorIndex] = tile.Tile(op ^ 1)
	}
}

// tileIndex converts tile co-ordinates to an index into
// tileRam / paletteRam.
//
//	top (0 <= x < 32, y < 2) - note x=0,1,30,31 are invisible
//	index := (y+30)*32 + (31-x)				// 0x3c0-0x3ff
//
//	normal (0 <= x < 28, 2 <= y < 34)
//	index := (29-x)*32 + (y - 2) 	// 0x040-0x3bf
//
//	bottom (0 <= x < 32, y < 2) - note x=0,1,30,31 are invisible
//	index := y*32 + (31-x) 					// 0x000-0x03f
func tileIndex(x int, y int) int {
	switch {
	case y < 2:
		return 0x3c0 + y*32 + 31 - x
	case y >= 34:
		return 0x000 + (y-34)*32 + 31 - x
	default:
		return 0x40 + (27-x)*32 + (y - 2)
	}
}

// SetTile places the specified tile into tile ram.
func (v *Video) SetTile(x, y int, t tile.Tile) {
	v.TileRam[tileIndex(x, y)] = t
}

// ColorTile sets the palette for a specific tile.
func (v *Video) ColorTile(x, y int, pal color.Palette) {
	v.palRam[tileIndex(x, y)] = pal
}

// GetTile returns the tile at the given co-ordinates.
func (v *Video) GetTile(x, y int) tile.Tile {
	return v.TileRam[tileIndex(x, y)]
}

// SetStatusQuad places 4 tiles (baseTile and 3 consecutively
// numbered tiles) in a 2x2 arrangement at the given posiiton
// in the lower status area of the display.
func (v *Video) SetStatusQuad(baseX int, baseTile tile.Tile, pal color.Palette) {
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

// DrawTiles paints the supplied bitmap with the contents
// of tile ram mixed with the colours from palette ram.
func (v *Video) DrawTiles(screen *ebiten.Image) {
	for ty := 0; ty < 36; ty++ {
		for tx := range 28 {
			pos := geom.TilePos(tx, ty)
			op := colorm.DrawImageOptions{}
			op.GeoM.Translate(float64(hOffset+pos.X), float64(vOffset+pos.Y))
			op.GeoM.Scale(1, 1)
			index := tileIndex(tx, ty)
			colorm.DrawImage(screen, tile.Image[v.TileRam[index]], color.ColorM[v.palRam[index]], &op)
		}
	}
}

// Draw paints the supplied bitmap with tiles, with all sprites
// established for this frame rendered on top.
func (v *Video) Draw(screen *ebiten.Image) {
	v.DrawTiles(screen)
	v.DrawSprites(screen)
}
