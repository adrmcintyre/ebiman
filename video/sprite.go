package video

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const maxSprites = 6

type spriteState struct {
	Look         sprite.Look
	Pal          color.Palette
	Pos          geom.Position
	FlipX, FlipY bool
}

func (v *Video) ClearSprites() {
	v.spriteCount = 0
}

func (v *Video) AddSprite(pos geom.Position, look sprite.Look, pal color.Palette) {
	if v.spriteCount < maxSprites {
		v.sprites[v.spriteCount] = spriteState{
			Pos:   pos,
			FlipX: false,
			FlipY: false,
			Look:  look,
			Pal:   pal,
		}
		v.spriteCount = v.spriteCount + 1
	}
}

func (v *Video) DrawSprites(screen *ebiten.Image) {
	for i := range v.spriteCount {
		s := v.sprites[i]
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
		colorm.DrawImage(screen, sprite.Image[s.Look], color.ColorM[s.Pal], &op)
	}
}
