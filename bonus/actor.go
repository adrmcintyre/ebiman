package bonus

import (
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/video"
)

// An Actor represents a fruity bonus.
type Actor struct {
	Visible bool          // is the bonus visible?
	Pos     geom.Position // where to display it
}

// NewActor returns a new bonus actor, correctly positioned,
// but initially invisible.
func NewActor() *Actor {
	return &Actor{
		Pos: geom.BONUS_POS,
	}
}

// Start gets the bonus actor ready for start of play
// (initially invisible).
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
