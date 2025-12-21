package bonus

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/video"
)

// An Actor represents a fruity bonus.
type Actor struct {
	Visible bool
	Pos     geom.Position
}

// MakeActor creates a new bonus actor, correctly positioned,
// but initially invisible.
func MakeActor() Actor {
	return Actor{
		Pos: geom.BONUS_POS,
	}
}

// TODO
func (b *Actor) Init() {
	b.Visible = false
}

// TODO
func (b *Actor) Start() {
	b.Visible = false
}

// Draw prepares the bonus's sprite for rendering at the next frame.
// An InfoEntry provides the look.
func (b *Actor) Draw(v *video.Video, info InfoEntry) {
	if b.Visible {
		v.AddSprite(b.Pos, info.Look, info.Pal)
	}
}
