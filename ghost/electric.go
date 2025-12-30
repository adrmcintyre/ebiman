package ghost

import (
	"math/rand"

	"github.com/adrmcintyre/ebiman/tile"
	"github.com/adrmcintyre/ebiman/video"
)

// CheckModifyCharge gives a ghost the chance to change the charge on the
// pill beneath it. The returned value is the net change in charge.
func (g *Actor) CheckModifyCharge(v *video.Video, frameCounter int) int {
	// percentage chances of modifying charge under the ghost
	const (
		scaredPct = 0
		inkyPct   = 70
		pinkyPct  = 70
		blinkyPct = 60
		clydePct  = 60
	)

	x, y := g.Pos.TileXY()
	t := v.GetTile(x, y)
	if !t.IsPill() {
		return 0
	}

	charge := t.Charge()
	newCharge := charge

	r := rand.Intn(100)
	switch {
	// scared ghosts bring pills one unit closer to neutral
	case g.Mode == MODE_PLAYING && g.SubMode == SUBMODE_SCARED:
		if r < scaredPct {
			if charge < 0 {
				newCharge += 1
			} else if charge > 0 {
				newCharge -= 1
			}
		}
	case g.Id == BLINKY:
		// blinky gives neutral pills one negative unit of charge
		if r < blinkyPct && charge == 0 {
			newCharge -= 1
		}
	case g.Id == PINKY:
		// pinky gives neutral pills one positive unit of charge
		if r < pinkyPct && charge == 0 {
			newCharge += 1
		}
	case g.Id == INKY:
		// inky switches between leaving negative and positive charges
		// approx every 8 seconds
		if r < inkyPct && charge == 0 {
			if frameCounter&512 == 0 {
				newCharge += 1
			} else {
				newCharge -= 1
			}
		}
	case g.Id == CLYDE:
		// clyde doubles any unit charges he passes over
		if r < clydePct && (charge == -1 || charge == 1) {
			newCharge = 2 * charge
		}
	}
	if newCharge != charge {
		v.SetTile(x, y, tile.FromCharge(newCharge))
	}
	return newCharge - charge
}
