package video

import (
	"fmt"

	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/tile"
)

var puncTile = map[rune]tile.Tile{
	'-': tile.MINUS,
	'*': tile.POWER_SMALL,
	'.': tile.POINT,
	' ': tile.SPACE,
	'"': tile.QUOTES,
	'/': tile.SLASH,
	'!': tile.EXCLAM,
}

func runeTile(ch rune) tile.Tile {
	switch {
	case ch >= '0' && ch <= '9':
		return tile.DIGIT_BASE + tile.Tile(ch-'0')
	case ch >= 'A' && ch <= 'Z':
		return tile.ALPHA_BASE + tile.Tile(ch-'A')
	default:
		if maybeTile, ok := puncTile[ch]; ok {
			return maybeTile
		}
	}
	return tile.PILL
}

func (v *Video) SetCursor(x int, y int) {
	v.CursorX = x
	v.CursorY = y
}

func (v *Video) WriteTile(t tile.Tile, pal color.Palette) {
	v.SetTile(v.CursorX, v.CursorY, t)
	v.ColorTile(v.CursorX, v.CursorY, pal)
	v.CursorX += 1
}

func (v *Video) WriteTiles(tiles []tile.Tile, pal color.Palette) {
	for _, t := range tiles {
		v.WriteTile(t, pal)
	}
}

func (v *Video) WriteChar(ch rune, pal color.Palette) {
	v.WriteTile(runeTile(ch), pal)
}

func (v *Video) WriteString(s string, pal color.Palette) {
	for _, ch := range s {
		v.WriteChar(ch, pal)
	}
}

func (v *Video) ClearRight() {
	for v.CursorX < 28 {
		v.WriteChar(' ', color.PAL_BLACK)
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
	v.SetTile(3, 0, runeTile('1'))
	v.SetTile(4, 0, runeTile('U'))
	v.SetTile(5, 0, runeTile('P'))
}

func (v *Video) Clear1Up() {
	v.SetTile(3, 0, tile.SPACE)
	v.SetTile(4, 0, tile.SPACE)
	v.SetTile(5, 0, tile.SPACE)
}

func (v *Video) Write2Up() {
	v.SetTile(22, 0, runeTile('2'))
	v.SetTile(23, 0, runeTile('U'))
	v.SetTile(24, 0, runeTile('P'))
}

func (v *Video) Clear2Up() {
	v.SetTile(22, 0, tile.SPACE)
	v.SetTile(23, 0, tile.SPACE)
	v.SetTile(24, 0, tile.SPACE)
}

func (v *Video) WriteLives(lives int) {
	// lives are added on the right - a maximum of 5 are displayed
	for i := range 5 {
		baseTile := tile.SPACE_BASE
		if lives > i {
			baseTile = tile.PACMAN_BASE
		}
		v.SetStatusQuad(2+i*2, baseTile, color.PAL_PACMAN)
	}
}

func (v *Video) WriteScoreAt(x, y int, value int) {
	buf := fmt.Sprintf("%5d%d", (value/10)%100000, value%10)
	for i, ch := range buf {
		v.SetTile(x+i, y, runeTile(ch))
	}
}

func (v *Video) WriteHighScore(score int) {
	txt := "HIGH SCORE"
	for i, ch := range txt {
		v.SetTile(9+i, 0, runeTile(ch))
	}
	v.WriteScoreAt(11, 1, score)
}
