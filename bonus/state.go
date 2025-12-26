package bonus

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/tile"
	"github.com/adrmcintyre/ebiman/video"
)

const (
	MaxStatusCount = 7
)

// Status represents the consumed bonuses to display in the status area.
type Status struct {
	Indicator [MaxStatusCount]Id // last few bonuses awarded, most recent first
	Count     int                // how many bonuses awarded so far
}

// ClearBonuses resets the list.
func (s *Status) ClearBonuses() {
	s.Count = 0
}

// AddBonus adds a new bonus to the status area, removing
// the oldest one to make room if there would be too many
// to display.
func (s *Status) AddBonus(bonus Id) {
	s.Count += 1
	copy(s.Indicator[1:], s.Indicator[:MaxStatusCount-1])
	s.Indicator[0] = bonus
}

// Write updates the tiles in the status area to reflect the
// current list, ready for the next frame.
func (s *Status) Write(v *video.Video) {
	tileBase := tile.SPACE_BASE
	pal := color.PAL_BLACK
	j := 0
	for i := range MaxStatusCount {
		if i+s.Count >= MaxStatusCount {
			info := &Info[s.Indicator[j]]
			j++
			tileBase = info.BaseTile
			pal = info.Pal
		}
		v.SetStatusQuad(12+i*2, tileBase, pal)
	}
}
