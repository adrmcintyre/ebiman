package sprite

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Image contains an ebiten Image for each sprite identifier.
var Image [64]*ebiten.Image

// MakeImages initialises the Image cache from the 2-bpp source data.
func MakeImages() {
	for i := range 64 {
		img := ebiten.NewImage(16, 16)
		for y := range 16 {
			u32 := uint32(bitmapData[i][y*2+1])<<16 | uint32(bitmapData[i][y*2])
			for x := range 16 {
				c := color.RGBA{}
				switch (u32 >> (x * 2)) & 0b11 {
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
