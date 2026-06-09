package actor

import (
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/video"
)

// An Bonus represents a fruity bonus.
type Bonus struct {
	Visible bool          // is the bonus visible?
	Pos     geom.Position // where to display it
}

// NewActor returns a new bonus actor, correctly positioned,
// but initially invisible.
func NewBonus() *Bonus {
	return &Bonus{
		Pos: geom.BonusPos,
	}
}

// Start gets the bonus actor ready for start of play
// (initially invisible).
func (b *Bonus) Start() {
	b.Visible = false
}

// Draw prepares the bonus's sprite for rendering at the next frame.
// An InfoEntry provides the look.
func (b *Bonus) Draw(v *video.Video, info data.BonusInfoEntry) {
	if b.Visible {
		v.AddSprite(b.Pos, info.Look, info.Pal)
	}
}
