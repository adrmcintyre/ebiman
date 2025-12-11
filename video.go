package main

import (
	"fmt"

	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type SpriteState struct {
	Sprite       byte
	Palette      byte
	X, Y         int
	FlipX, FlipY bool
}

const MAX_SPRITES = 6

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

func (v *Video) DecodePellets(dots *DotState) {
	src := 0
	dst := 0
	for b := range 30 {
		a := dots.PillBits[b]
		for mask := byte(0x80); mask > 0; mask >>= 1 {
			dst += int(data.Pill[src])
			src++
			if a&mask != 0 {
				v.TileRam[dst] = tile.PILL
			} else {
				v.TileRam[dst] = tile.SPACE
			}
		}
	}
	v.TileRam[3*32+4] = dots.PowerPills[0]
	v.TileRam[3*32+24] = dots.PowerPills[1]
	v.TileRam[28*32+4] = dots.PowerPills[2]
	v.TileRam[28*32+24] = dots.PowerPills[3]
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

func (v *Video) SetCursor(x, y int) {
	v.X = x
	v.Y = y
}

func (v *Video) WriteTile(t byte, pal byte) {
	i := tileIndex(v.X, v.Y)
	v.TileRam[i] = t
	v.PalRam[i] = pal
	v.X += 1
}

func (v *Video) WriteTiles(tiles []byte, pal byte) {
	for _, t := range tiles {
		v.WriteTile(t, pal)
	}
}

var charToTile = map[rune]byte{
	'-': tile.MINUS,
	'*': tile.POWER_SMALL,
	'.': tile.POINT,
	' ': tile.SPACE,
	'"': tile.QUOTES,
	'/': tile.SLASH,
	'!': tile.EXCLAM,
}

func (v *Video) WriteChar(ch rune, pal byte) {
	t := tile.PILL
	switch {
	case ch >= '0' && ch <= '9':
		t = tile.DIGIT_BASE + byte(ch-'0')
	case ch >= 'A' && ch <= 'Z':
		t = tile.ALPHA_BASE + byte(ch-'A')
	default:
		if maybeTile, ok := charToTile[ch]; ok {
			t = maybeTile
		}
	}
	v.WriteTile(t, pal)
}

func (v *Video) WriteString(s string, pal byte) {
	for _, ch := range s {
		v.WriteChar(ch, pal)
	}
}

func (v *Video) ClearRight() {
	for v.X < 28 {
		v.WriteChar(' ', palette.BLACK)
	}
}

func (v *Video) WritePlayerUp(i int) {
	if i == 0 {
		v.Write1Up()
	} else {
		v.Write2Up()
	}
}

func (v *Video) ClearPlayerUp(i int) {
	if i == 0 {
		v.Clear1Up()
	} else {
		v.Clear2Up()
	}
}

func (v *Video) Write1Up() {
	v.SetTopTile(3, 0, tile.DIGIT_BASE+1)
	v.SetTopTile(4, 0, tile.ALPHA_BASE+('U'-'A'))
	v.SetTopTile(5, 0, tile.ALPHA_BASE+('P'-'A'))
}

func (v *Video) Clear1Up() {
	v.SetTopTile(3, 0, tile.SPACE)
	v.SetTopTile(4, 0, tile.SPACE)
	v.SetTopTile(5, 0, tile.SPACE)
}

func (v *Video) Write2Up() {
	v.SetTopTile(22, 0, tile.DIGIT_BASE+2)
	v.SetTopTile(23, 0, tile.ALPHA_BASE+('U'-'A'))
	v.SetTopTile(24, 0, tile.ALPHA_BASE+('P'-'A'))
}

func (v *Video) Clear2Up() {
	v.SetTopTile(22, 0, tile.SPACE)
	v.SetTopTile(23, 0, tile.SPACE)
	v.SetTopTile(24, 0, tile.SPACE)
}

func (v *Video) WriteLives(lives int) {
	// lives are added on the right - a maximum of 5 are displayed
	for i := range 5 {
		x := i*2 + 2
		baseTile := tile.SPACE_BASE
		if lives > i {
			baseTile = tile.PACMAN_BASE
		}
		pal := palette.PACMAN
		v.SetBottomTile(x+1, 0, baseTile+0, pal)
		v.SetBottomTile(x+0, 0, baseTile+1, pal)
		v.SetBottomTile(x+1, 1, baseTile+2, pal)
		v.SetBottomTile(x+0, 1, baseTile+3, pal)
	}
}

func (v *Video) WriteScoreAt(x, y int, value int) {
	buf := fmt.Sprintf("%5d%d", (value/10)%100000, value%10)
	for i := range 6 {
		ch := buf[i]
		t := tile.SPACE
		if ch >= '0' && ch <= '9' {
			t = tile.DIGIT_BASE + (ch - '0')
		}
		v.SetTopTile(x+i, y, t)
	}
}

func (v *Video) WriteHighScore(score int) {
	txt := "HIGH SCORE"
	for i := range 10 {
		t := tile.SPACE
		if txt[i] != ' ' {
			t = tile.ALPHA_BASE + (txt[i] - 'A')
		}
		v.SetTopTile(9+i, 0, t)
	}
	v.WriteScoreAt(11, 1, score)
}

var pillCoords = [4]struct{ X, Y int }{
	{1, 6},
	{26, 6},
	{1, 26},
	{26, 26},
}

func (v *Video) FlashPills() {
	const FLASH_FRAMES = 10

	v.FlashCycle += 1
	if v.FlashCycle > FLASH_FRAMES {
		v.FlashCycle = 0
		v.FlashOff = !v.FlashOff

		pal := palette.BLACK
		if v.FlashOff {
			pal = palette.MAZE
		}
		for _, coords := range pillCoords {
			index := tileIndex(coords.X, coords.Y)
			if v.TileRam[index] == tile.POWER {
				v.PalRam[index] = pal
			}
		}
	}
}

func (v *Video) ClearSprites() {
	v.SpriteCount = 0
}

func (v *Video) AddSprite(x, y int, sprite byte, pal byte) {
	if v.SpriteCount < MAX_SPRITES {
		v.Sprites[v.SpriteCount] = SpriteState{
			X:       x,
			Y:       y,
			FlipX:   false,
			FlipY:   false,
			Sprite:  sprite,
			Palette: pal,
		}
		v.SpriteCount = v.SpriteCount + 1
	}
}

const hOffset = 8.0
const vOffset = 8.0

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

func (v *Video) DrawSprites(screen *ebiten.Image) {
	for i := range v.SpriteCount {
		s := v.Sprites[i]
		if s.X <= 0 && s.Y <= 0 {
			continue
		}
		screenX := s.X
		screenY := s.Y + 16
		scaleX, scaleY := 1.0, 1.0
		if s.FlipX {
			scaleX = -1.0
		}
		if s.FlipY {
			scaleY = -1.0
		}
		op := colorm.DrawImageOptions{}
		op.GeoM.Translate(hOffset+float64(screenX), vOffset+float64(screenY))
		op.GeoM.Scale(scaleX, scaleY)
		colorm.DrawImage(screen, sprite.Image[s.Sprite], palette.ColorM[s.Palette], &op)
	}
}

func (v *Video) Draw(screen *ebiten.Image) {
	v.DrawMaze(screen)
	v.DrawStatus(screen)
	v.DrawSprites(screen)
}
