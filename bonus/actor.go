package bonus

import (
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/video"
)

type Actor struct {
	Visible bool
	Pos     geom.Position
}

func MakeBonusActor() Actor {
	return Actor{
		Pos: geom.BONUS_POS,
	}
}

func (b *Actor) BonusInit() {
	b.Visible = false
}

func (b *Actor) BonusStart() {
	b.Visible = false
}

func (b *Actor) Draw(v *video.Video, bonusInfo InfoEntry) {
	if b.Visible {
		v.AddSprite(b.Pos, bonusInfo.Look, bonusInfo.Pal)
	}
}
