package video

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
	// The maximum number of sprites supported by the simulated hardware.
	// This is enough for pacman, 4 ghosts, and the bonus fruit.
	maxSprites = 6
)

var (
	// An adjustment to make sure sprites are correctly centred.
	centreOffset = geom.Delta{-4, -4}
)

// A spriteState holds the appearance of a sprite.
type spriteState struct {
	Look  sprite.Look   // identifies the bitmap to render
	Pal   color.Palette // identifies the palette to apply
	Pos   geom.Position // the position in screen co-ordinates
	FlipX bool          // is the sprite mirrored horizontally?
	FlipY bool          // is the sprite mirrored vertically?
}

// ClearSprites removes all sprites from display for the next frame.
// Call AddSprite() for each sprite to display before the next Draw() call.
func (v *Video) ClearSprites() {
	v.spriteCount = 0
}

// AddSprite specifies a sprite to display.
func (v *Video) AddSprite(pos geom.Position, look sprite.Look, pal color.Palette) {
	if v.spriteCount < maxSprites {
		v.sprites[v.spriteCount] = spriteState{
			Pos:   pos.Add(centreOffset),
			FlipX: false,
			FlipY: false,
			Look:  look,
			Pal:   pal,
		}
		v.spriteCount = v.spriteCount + 1
	}
}

// DrawSprites paints all the sprites established for this frame onto the supplied bitmap.
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
		op.GeoM.Translate(float64(hOffset+s.Pos.X), float64(vOffset+s.Pos.Y))
		op.GeoM.Scale(scaleX, scaleY)
		colorm.DrawImage(screen, sprite.Image[s.Look], color.ColorM[s.Pal], &op)
	}
}
