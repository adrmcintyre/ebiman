package tile

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// imageCache defines an ebiten Image for each tile identifier.
var imageCache [count]*ebiten.Image

// Init initialises the Image cache for each tile.
func Init() {
	for i, bitmap := range bitmapData {
		img := ebiten.NewImage(width, height)
		for y, row := range bitmap {
			for x := range width {
				img.Set(x, y, color.Channel[row&0b11])
				row >>= 2
			}
		}
		imageCache[i] = img
	}
}

// Draw paints the tile onto img.
func (t Tile) Draw(img *ebiten.Image, x, y int, pal color.Palette) {
	op := colorm.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Scale(1, 1)
	colorm.DrawImage(img, imageCache[t], color.ColorM[pal], &op)

}

// IsTraversable returns true if the tile can be passed over (i.e. not a maze barrier).
func (t Tile) IsTraversable() bool {
	return t == Space || t.IsPill() || t.IsPower() || t >= ScoreMin && t <= ScoreMax
}

// IsPill returns true if the tile is a pill of some kind.
func (t Tile) IsPill() bool {
	switch t {
	case Pill, PillMinus, PillPlus, PillMinus2, PillPlus2:
		return true
	default:
		return false
	}
}

// IsPower returns true if the tile is a power pill.
func (t Tile) IsPower() bool {
	switch t {
	case Power, PowerSmall:
		return true
	default:
		return false
	}
}

// Charge returns the net charge on a tile.
func (t Tile) Charge() int {
	switch t {
	case Pill:
		return 0
	case PillMinus:
		return -1
	case PillPlus:
		return 1
	case PillMinus2:
		return -2
	case PillPlus2:
		return 2
	default:
		return 0
	}
}

// TileFromCharge returns a tile with the specified charge.
func FromCharge(c int) Tile {
	switch c {
	case -1:
		return PillMinus
	case 1:
		return PillPlus
	case -2:
		return PillMinus2
	case 2:
		return PillPlus2
	default:
		return Pill
	}
}
