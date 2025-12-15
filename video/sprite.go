package video

import (
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const MAX_SPRITES = 6

type SpriteState struct {
	Sprite       byte
	Palette      byte
	X, Y         int
	FlipX, FlipY bool
}

func (v *Video) ClearSprites() {
	v.SpriteCount = 0
}

func (v *Video) AddSprite(x, y int, sprite byte, pal byte) {
	if v.SpriteCount < MAX_SPRITES {
		v.Sprites[v.SpriteCount] = SpriteState{
			X:       x,
			Y:       y,
			FlipX:   false,
			FlipY:   false,
			Sprite:  sprite,
			Palette: pal,
		}
		v.SpriteCount = v.SpriteCount + 1
	}
}

func (v *Video) DrawSprites(screen *ebiten.Image) {
	for i := range v.SpriteCount {
		s := v.Sprites[i]
		if s.X <= 0 && s.Y <= 0 {
			continue
		}
		screenPos := ScreenPos{s.X, s.Y + 16}
		scaleX, scaleY := 1.0, 1.0
		if s.FlipX {
			scaleX = -1.0
		}
		if s.FlipY {
			scaleY = -1.0
		}
		op := colorm.DrawImageOptions{}
		op.GeoM.Translate(float64(hOffset+screenPos.X), float64(vOffset+screenPos.Y))
		op.GeoM.Scale(scaleX, scaleY)
		colorm.DrawImage(screen, sprite.Image[s.Sprite], palette.ColorM[s.Palette], &op)
	}
}
