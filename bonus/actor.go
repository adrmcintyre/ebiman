package bonus

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/video"
)

type Actor struct {
	Visible bool
	Pos     geom.Position
}

func MakeActor() Actor {
	return Actor{
		Pos: geom.BONUS_POS,
	}
}

func (b *Actor) Init() {
	b.Visible = false
}

func (b *Actor) Start() {
	b.Visible = false
}

func (b *Actor) Draw(v *video.Video, info InfoEntry) {
	if b.Visible {
		v.AddSprite(b.Pos, info.Look, info.Pal)
	}
}
