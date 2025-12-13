package video

import (
	"fmt"

	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
)

var charToTile = map[rune]byte{
	'-': tile.MINUS,
	'*': tile.POWER_SMALL,
	'.': tile.POINT,
	' ': tile.SPACE,
	'"': tile.QUOTES,
	'/': tile.SLASH,
	'!': tile.EXCLAM,
}

func toTile(ch rune) byte {
	switch {
	case ch >= '0' && ch <= '9':
		return tile.DIGIT_BASE + byte(ch-'0')
	case ch >= 'A' && ch <= 'Z':
		return tile.ALPHA_BASE + byte(ch-'A')
	default:
		if maybeTile, ok := charToTile[ch]; ok {
			return maybeTile
		}
	}
	return tile.PILL
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

func (v *Video) WriteChar(ch rune, pal byte) {
	v.WriteTile(toTile(ch), pal)
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
	v.SetTopTile(3, 0, toTile('1'))
	v.SetTopTile(4, 0, toTile('U'))
	v.SetTopTile(5, 0, toTile('P'))
}

func (v *Video) Clear1Up() {
	v.SetTopTile(3, 0, tile.SPACE)
	v.SetTopTile(4, 0, tile.SPACE)
	v.SetTopTile(5, 0, tile.SPACE)
}

func (v *Video) Write2Up() {
	v.SetTopTile(22, 0, toTile('2'))
	v.SetTopTile(23, 0, toTile('U'))
	v.SetTopTile(24, 0, toTile('P'))
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
	for i, ch := range buf {
		v.SetTopTile(x+i, y, toTile(ch))
	}
}

func (v *Video) WriteHighScore(score int) {
	txt := "HIGH SCORE"
	for i, ch := range txt {
		v.SetTopTile(9+i, 0, toTile(ch))
	}
	v.WriteScoreAt(11, 1, score)
}
