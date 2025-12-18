package bonus

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/video"
)

type BonusActor struct {
	Visible bool
	Pos     geom.Position
}

func MakeBonusActor() BonusActor {
	return BonusActor{
		Pos: geom.BONUS_POS,
	}
}

func (b *BonusActor) BonusInit() {
	b.Visible = false
}

func (b *BonusActor) BonusStart() {
	b.Visible = false
}

func (b *BonusActor) DrawBonus(v *video.Video, bonusInfo data.BonusInfoEntry) {
	if b.Visible {
		v.AddSprite(b.Pos, bonusInfo.Look, bonusInfo.Pal)
	}
}
