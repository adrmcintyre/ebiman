package sprite

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/hajimehoshi/ebiten/v2"
)

// Image contains an ebiten Image for each sprite identifier.
var Image [count]*ebiten.Image

// MakeImages initialises the Image cache from the 2-bpp source data.
func MakeImages() {
	for i, bitmap := range bitmapData {
		img := ebiten.NewImage(width, height)
		for y, row := range bitmap {
			for x := range width {
				img.Set(x, y, color.Channel[row&0b11])
				row >>= 2
			}
		}
		Image[i] = img
	}
}
