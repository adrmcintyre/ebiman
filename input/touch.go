package input

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// A TouchLayout associates Regions of the screen with keys.
type TouchLayout struct {
	bind map[ebiten.Key]Region
}

// NewTouchLayout returns a new TouchLayout.
func NewTouchLayout() *TouchLayout {
	return &TouchLayout{
		bind: make(map[ebiten.Key]Region),
	}
}

// BindRegion associates a touch region with a key.
func (tl *TouchLayout) BindRegion(r Region, key ebiten.Key, visible bool) {
	if !visible {
		r = invisibleRegion{r}
	}
	tl.bind[key] = r
}

// IsJustTouched returns true if the region corresponding to key was just touched.
func (i *Input) IsJustTouched(key ebiten.Key) bool {
	if r, ok := i.touchLayout.bind[key]; ok {
		for _, id := range inpututil.JustPressedTouchIDs() {
			if r.Contains(ebiten.TouchPosition(id)) {
				return true
			}
		}
	}
	return false
}

// IsTouched returns true if the region corresponding to key is currently being touched.
func (i *Input) IsTouched(key ebiten.Key) bool {
	if r, ok := i.touchLayout.bind[key]; ok {
		for _, id := range ebiten.TouchIDs() {
			if r.Contains(ebiten.TouchPosition(id)) {
				return true
			}
		}
	}
	return false
}

// Draw renders all bound touch regions to img.
func (i *Input) Draw(img *ebiten.Image) {
	for k, r := range i.touchLayout.bind {
		r.Draw(img)

		// TODO - this is tacky: we should keep the sprite
		// (or potentially tile) with the binding.
		var look sprite.Look
		var pal color.Palette
		switch k {
		case ebiten.KeyLeft:
			look = sprite.PacmanLeft1
			pal = color.PalPacman
		case ebiten.KeyRight:
			look = sprite.PacmanRight1
			pal = color.PalPacman
		case ebiten.KeyUp:
			look = sprite.PacmanUp1
			pal = color.PalPacman
		case ebiten.KeyDown:
			look = sprite.PacmanDown1
			pal = color.PalPacman
		}
		cx, cy := r.Centre()
		if cx > 0 && cy > 0 {
			look.Draw(img, cx-8, cy-8, false, false, pal)
		}
	}
}
