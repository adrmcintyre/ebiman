package game

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/video"
)

var (
	BONUS_POS = geom.Position{HOME_CENTRE.X, 160}
)

type BonusActor struct {
	Visible bool
	Pos     geom.Position
}

func MakeBonus() BonusActor {
	return BonusActor{
		Pos: BONUS_POS,
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
