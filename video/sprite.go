package video

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const MAX_SPRITES = 6

type SpriteState struct {
	Look         byte
	Pal          byte
	Pos          geom.Position
	FlipX, FlipY bool
}

func (v *Video) ClearSprites() {
	v.SpriteCount = 0
}

func (v *Video) AddSprite(pos geom.Position, sprite byte, pal byte) {
	if v.SpriteCount < MAX_SPRITES {
		v.Sprites[v.SpriteCount] = SpriteState{
			Pos:   pos,
			FlipX: false,
			FlipY: false,
			Look:  sprite,
			Pal:   pal,
		}
		v.SpriteCount = v.SpriteCount + 1
	}
}

func (v *Video) DrawSprites(screen *ebiten.Image) {
	for i := range v.SpriteCount {
		s := v.Sprites[i]
		if s.Pos.X <= 0 && s.Pos.Y <= 0 {
			continue
		}
		scaleX, scaleY := 1.0, 1.0
		if s.FlipX {
			scaleX = -1.0
		}
		if s.FlipY {
			scaleY = -1.0
		}
		op := colorm.DrawImageOptions{}
		op.GeoM.Translate(float64(hOffset+s.Pos.X), float64(vOffset+s.Pos.Y+16))
		op.GeoM.Scale(scaleX, scaleY)
		colorm.DrawImage(screen, sprite.Image[s.Look], palette.ColorM[s.Pal], &op)
	}
}
