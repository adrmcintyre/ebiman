package bonus

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

type State struct {
	Indicator [7]Id // last 7 bonuses awarded, most recent first
	Count     int   // how many bonuses awarded so far
}

func (s *State) ClearBonuses() {
	s.Count = 0
}

func (s *State) AddBonus(bonus Id) {
	s.Count += 1
	copy(s.Indicator[1:], s.Indicator[:6])
	s.Indicator[0] = bonus
}

func (s *State) WriteBonuses(v *video.Video) {
	tileBase := tile.SPACE_BASE
	pal := color.PAL_BLACK
	j := 0
	for i := range 7 {
		if i+s.Count >= 7 {
			info := &Info[s.Indicator[j]]
			j++
			tileBase = info.BaseTile
			pal = info.Pal
		}
		v.SetStatusQuad(12+i*2, tileBase, pal)
	}
}
