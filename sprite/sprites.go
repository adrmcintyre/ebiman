package sprite

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// imageCache contains an ebiten Image for each sprite identifier.
var imageCache [count]*ebiten.Image

// Init initialises the Image cache from the 2-bpp source data.
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

// Draw paints the sprite onto img.
func (look Look) Draw(img *ebiten.Image, x, y int, flipX, flipY bool, pal color.Palette) {
	scaleX, scaleY := 1.0, 1.0
	if flipX {
		scaleX = -1.0
	}
	if flipY {
		scaleY = -1.0
	}
	// draw centred
	op := colorm.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Scale(scaleX, scaleY)
	colorm.DrawImage(img, imageCache[look], color.ColorM[pal], &op)
}
