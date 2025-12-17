package tile

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var Image [256]*ebiten.Image

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
