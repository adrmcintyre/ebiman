package game

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

// -------------------------------------------------------------------------
// Bonuses (fruit)
// -------------------------------------------------------------------------

type BonusState struct {
	BonusIndicator [7]int // last 7 bonuses awarded, most recent first
	BonusCount     int    // how many bonuses awarded so far
}

func (bs *BonusState) ClearBonuses() {
	bs.BonusCount = 0
}

func (bs *BonusState) AddBonus(bonus int) {
	bs.BonusCount += 1
	copy(bs.BonusIndicator[1:], bs.BonusIndicator[:6])
	bs.BonusIndicator[0] = bonus
}

func (bs *BonusState) WriteBonuses(v *video.Video) {
	tileBase := tile.SPACE_BASE
	pal := color.PAL_BLACK
	j := 0
	for i := range 7 {
		if i+bs.BonusCount >= 7 {
			info := &data.BonusInfo[bs.BonusIndicator[j]]
			j++
			tileBase = info.BaseTile
			pal = info.Pal
		}
		v.SetStatusQuad(12+i*2, tileBase, pal)
	}
}
