package video

import (
	"fmt"

	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
)

var puncTile = map[rune]byte{
	'-': tile.MINUS,
	'*': tile.POWER_SMALL,
	'.': tile.POINT,
	' ': tile.SPACE,
	'"': tile.QUOTES,
	'/': tile.SLASH,
	'!': tile.EXCLAM,
}

func runeTile(ch rune) byte {
	switch {
	case ch >= '0' && ch <= '9':
		return tile.DIGIT_BASE + byte(ch-'0')
	case ch >= 'A' && ch <= 'Z':
		return tile.ALPHA_BASE + byte(ch-'A')
	default:
		if maybeTile, ok := puncTile[ch]; ok {
			return maybeTile
		}
	}
	return tile.PILL
}

func (v *Video) SetCursor(pos TilePos) {
	v.Cursor = pos
}

func (v *Video) WriteTile(t byte, pal byte) {
	v.SetTile(v.Cursor, t)
	v.ColorTile(v.Cursor, pal)
	v.Cursor.X += 1
}

func (v *Video) WriteTiles(tiles []byte, pal byte) {
	for _, t := range tiles {
		v.WriteTile(t, pal)
	}
}

func (v *Video) WriteChar(ch rune, pal byte) {
	v.WriteTile(runeTile(ch), pal)
}

func (v *Video) WriteString(s string, pal byte) {
	for _, ch := range s {
		v.WriteChar(ch, pal)
	}
}

func (v *Video) ClearRight() {
	for v.Cursor.X < 28 {
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
	v.SetTile(TilePos{3, 0}, runeTile('1'))
	v.SetTile(TilePos{4, 0}, runeTile('U'))
	v.SetTile(TilePos{5, 0}, runeTile('P'))
}

func (v *Video) Clear1Up() {
	v.SetTile(TilePos{3, 0}, tile.SPACE)
	v.SetTile(TilePos{4, 0}, tile.SPACE)
	v.SetTile(TilePos{5, 0}, tile.SPACE)
}

func (v *Video) Write2Up() {
	v.SetTile(TilePos{22, 0}, runeTile('2'))
	v.SetTile(TilePos{23, 0}, runeTile('U'))
	v.SetTile(TilePos{24, 0}, runeTile('P'))
}

func (v *Video) Clear2Up() {
	v.SetTile(TilePos{22, 0}, tile.SPACE)
	v.SetTile(TilePos{23, 0}, tile.SPACE)
	v.SetTile(TilePos{24, 0}, tile.SPACE)
}

func (v *Video) WriteLives(lives int) {
	// lives are added on the right - a maximum of 5 are displayed
	for i := range 5 {
		baseTile := tile.SPACE_BASE
		if lives > i {
			baseTile = tile.PACMAN_BASE
		}
		v.SetStatusQuad(2+i*2, baseTile, palette.PACMAN)
	}
}

func (v *Video) WriteScoreAt(x, y int, value int) {
	buf := fmt.Sprintf("%5d%d", (value/10)%100000, value%10)
	for i, ch := range buf {
		v.SetTile(TilePos{x + i, y}, runeTile(ch))
	}
}

func (v *Video) WriteHighScore(score int) {
	txt := "HIGH SCORE"
	for i, ch := range txt {
		v.SetTile(TilePos{9 + i, 0}, runeTile(ch))
	}
	v.WriteScoreAt(11, 1, score)
}
