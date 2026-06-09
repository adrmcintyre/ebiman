package player

import (
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/video"
)

const (
	MaxBonusStatusCount = 7
)

// BonusStatus represents the consumed bonuses to display in the status area.
type BonusStatus struct {
	Indicator [MaxBonusStatusCount]data.BonusId // last few bonuses awarded, most recent first
	Count     int                               // how many bonuses awarded so far
}

// ClearBonuses resets the list.
func (s *BonusStatus) ClearBonuses() {
	s.Count = 0
}

// AddBonus adds a new bonus to the status area, removing
// the oldest one to make room if there would be too many
// to display.
func (s *BonusStatus) AddBonus(id data.BonusId) {
	s.Count += 1
	copy(s.Indicator[1:], s.Indicator[:MaxBonusStatusCount-1])
	s.Indicator[0] = id
}

// Write updates the tiles in the status area to reflect the
// current list, ready for the next frame.
func (s *BonusStatus) Write(v *video.Video) {
	tileBase := video.TileSpaceBase
	pal := video.PalBlack
	j := 0
	for i := range MaxBonusStatusCount {
		if i+s.Count >= MaxBonusStatusCount {
			info := &data.BonusInfo[s.Indicator[j]]
			j++
			tileBase = info.BaseTile
			pal = info.Pal
		}
		v.SetStatusQuad(12+i*2, tileBase, pal)
	}
}
