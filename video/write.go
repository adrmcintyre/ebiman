package video

import (
	"fmt"

	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/tile"
)

// puncTile is a lookup table mapping some special characters to tiles.
var puncTile = map[rune]tile.Tile{
	'-': tile.MINUS,
	'*': tile.POWER_SMALL,
	'.': tile.POINT,
	' ': tile.SPACE,
	'"': tile.QUOTES,
	'/': tile.SLASH,
	'!': tile.EXCLAM,
}

// runeTile returns the tile corresponding to a given rune,
// or the PILL tile if there is no equivalent tile.
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

// SetCursor sets the position of the next tile to be placed
// by any of the WriteXXX() calls.
func (v *Video) SetCursor(x int, y int) {
	v.cursorX = x
	v.cursorY = y
}

// WriteTile places a single tile together with its palette
// at the current cursor position, and advances the cursor;
func (v *Video) WriteTile(t tile.Tile, pal color.Palette) {
	v.SetTile(v.cursorX, v.cursorY, t)
	v.ColorTile(v.cursorX, v.cursorY, pal)
	v.cursorX += 1
}

// WriteTiles places an array of tiles in sequence starting
// at the current cursor position, and advances the cursor.
// The tiles' palettes are set at the same time.
func (v *Video) WriteTiles(tiles []tile.Tile, pal color.Palette) {
	for _, t := range tiles {
		v.WriteTile(t, pal)
	}
}

// WriteChar places the tile corresponding to the given rune
// at the current cursor position, and advances the cursor.
// The tile's palette is set at the same time.
func (v *Video) WriteChar(ch rune, pal color.Palette) {
	v.WriteTile(runeTile(ch), pal)
}

// WriteString places a sequence of tiles corresponding to
// the string's runes starting at the current cursor position,
// and advances the cursor. The tiles' palettes are set at
// the same time.
func (v *Video) WriteString(s string, pal color.Palette) {
	for _, ch := range s {
		v.WriteChar(ch, pal)
	}
}

// ClearRight replaces all tiles at and to the right of the
// current cursor position with blanks, and sets their
// palettes to black. The cursor is advanced.
func (v *Video) ClearRight() {
	for v.cursorX < 28 {
		v.WriteChar(' ', color.PAL_BLACK)
	}
}

// WritePlayerUp writes "1UP" or "2UP" in the appropriate location
// in the top status area. The cursor is not affected.
func (v *Video) WritePlayerUp(i int) {
	if i == 0 {
		v.Write1Up()
	} else {
		v.Write2Up()
	}
}

// ClearPlayerUp blanks the tiles for the "1UP" or "2UP" messages.
// Their palettes are not changed. The cursor is not affected.
func (v *Video) ClearPlayerUp(i int) {
	if i == 0 {
		v.Clear1Up()
	} else {
		v.Clear2Up()
	}
}

// Write1Up writes the "1UP" message to the top status area.
// The cursor is not affected.
func (v *Video) Write1Up() {
	v.SetTile(3, 0, runeTile('1'))
	v.SetTile(4, 0, runeTile('U'))
	v.SetTile(5, 0, runeTile('P'))
}

// Clear1Up clears the "1UP" message from the top status area.
// The cursor is not affected.
func (v *Video) Clear1Up() {
	v.SetTile(3, 0, tile.SPACE)
	v.SetTile(4, 0, tile.SPACE)
	v.SetTile(5, 0, tile.SPACE)
}

// Write2Up writes the "2UP" message to the top status area.
// The cursor is not affected.
func (v *Video) Write2Up() {
	v.SetTile(22, 0, runeTile('2'))
	v.SetTile(23, 0, runeTile('U'))
	v.SetTile(24, 0, runeTile('P'))
}

// Clear2Up clears the "2UP" message from the top status area.
// The cursor is not affected.
func (v *Video) Clear2Up() {
	v.SetTile(22, 0, tile.SPACE)
	v.SetTile(23, 0, tile.SPACE)
	v.SetTile(24, 0, tile.SPACE)
}

// WriteLives updates the bottom status areas with tiles and palette
// representing the specified number of lives. Any excess lives
// already present are blanked out. The cursor is not affected.
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

// WriteScoreAt updates the specified location in the top
// status area with a sequence of tiles representing the
// given score. The cursor is not affected.
func (v *Video) WriteScoreAt(x, y int, value int) {
	buf := fmt.Sprintf("%5d%d", (value/10)%100000, value%10)
	for i, ch := range buf {
		v.SetTile(x+i, y, runeTile(ch))
	}
}

// WriteHighScore updates the high-score location in the
// top status area with a sequence of tiles representing
// "HIGH SCORE" and the given high score value. The
// cursor is not affected.
func (v *Video) WriteHighScore(score int) {
	txt := "HIGH SCORE"
	for i, ch := range txt {
		v.SetTile(9+i, 0, runeTile(ch))
	}
	v.WriteScoreAt(11, 1, score)
}
