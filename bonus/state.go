package bonus

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

type State struct {
	Indicator [7]int // last 7 bonuses awarded, most recent first
	Count     int    // how many bonuses awarded so far
}

func (bs *State) ClearBonuses() {
	bs.Count = 0
}

func (bs *State) AddBonus(bonus int) {
	bs.Count += 1
	copy(bs.Indicator[1:], bs.Indicator[:6])
	bs.Indicator[0] = bonus
}

func (bs *State) WriteBonuses(v *video.Video) {
	tileBase := tile.SPACE_BASE
	pal := color.PAL_BLACK
	j := 0
	for i := range 7 {
		if i+bs.Count >= 7 {
			info := &Info[bs.Indicator[j]]
			j++
			tileBase = info.BaseTile
			pal = info.Pal
		}
		v.SetStatusQuad(12+i*2, tileBase, pal)
	}
}
