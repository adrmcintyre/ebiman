package tile

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Image defines an ebiten Image for each tile identifier.
var Image [256]*ebiten.Image

// MakeImages initialises the Image cache for each tile.
func MakeImages() {
	for i := range 256 {
		img := ebiten.NewImage(8, 8)
		for y := range 8 {
			u16 := bitmapData[i][y]
			for x := range 8 {
				c := color.RGBA{}
				switch (u16 >> (x * 2)) & 0b11 {
				case 0b10:
					c = color.RGBA{0xff, 0x00, 0x00, 0xff} // colour 1
				case 0b01:
					c = color.RGBA{0x00, 0xff, 0x00, 0xff} // colour 2
				case 0b11:
					c = color.RGBA{0x00, 0x00, 0xff, 0xff} // colour 3
				}
				img.Set(x, y, c)
			}
		}
		Image[i] = img
	}
}

// IsTraversable returns true if the tile can be passed over (i.e. not a maze barrier).
func (t Tile) IsTraversable() bool {
	return t == SPACE || t.IsPill() || t.IsPower() || t >= SCORE_MIN && t <= SCORE_MAX
}

// IsPill returns true if the tile is a pill of some kind.
func (t Tile) IsPill() bool {
	switch t {
	case PILL, PILL_MINUS, PILL_PLUS, PILL_MINUS2, PILL_PLUS2:
		return true
	default:
		return false
	}
}

// IsPower returns true if the tile is a power pill.
func (t Tile) IsPower() bool {
	switch t {
	case POWER, POWER_SMALL:
		return true
	default:
		return false
	}
}

// Charge returns the net charge on a tile.
func (t Tile) Charge() int {
	switch t {
	case PILL:
		return 0
	case PILL_MINUS:
		return -1
	case PILL_PLUS:
		return 1
	case PILL_MINUS2:
		return -2
	case PILL_PLUS2:
		return 2
	default:
		return 0
	}
}

// TileFromCharge returns a tile with the specified charge.
func FromCharge(c int) Tile {
	switch c {
	case -1:
		return PILL_MINUS
	case 1:
		return PILL_PLUS
	case -2:
		return PILL_MINUS2
	case 2:
		return PILL_PLUS2
	default:
		return PILL
	}
}
