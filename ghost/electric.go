package ghost

import (
	"math/rand"

	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/tile"
	"github.com/adrmcintyre/ebiman/video"
)

// CheckModifyCharge gives a ghost the chance to change the charge on the
// pill beneath it. The returned value is the net change in charge.
func (g *Actor) CheckModifyCharge(v *video.Video, frameCounter int, electric data.ElectricEntry) int {
	x, y := g.Pos.TileXY()
	t := v.GetTile(x, y)
	if !t.IsPill() {
		return 0
	}

	charge := t.Charge()
	newCharge := charge

	r := rand.Intn(100)
	switch {
	// scared ghosts bring charged pills towards neutral
	case g.Mode == ModePlaying && g.SubMode == SubModeScared:
		if r < electric.ScaredPct {
			if charge < 0 {
				newCharge += 1
			} else if charge > 0 {
				newCharge -= 1
			}
		}
	case g.Id == Blinky:
		// blinky gives neutral pills a -ve charge
		if r < electric.BlinkyPct && charge == 0 {
			newCharge -= 1
		}
	case g.Id == Pinky:
		// pinky gives neutral pills a +ve charge
		if r < electric.PinkyPct && charge == 0 {
			newCharge += 1
		}
	case g.Id == Inky:
		// inky switches between leaving -ve and +ve charges approx every 8 seconds
		if r < electric.InkyPct && charge == 0 {
			if frameCounter&512 == 0 {
				newCharge += 1
			} else {
				newCharge -= 1
			}
		}
	case g.Id == Clyde:
		// clyde doubles any single charged pill he passes over
		if r < electric.ClydePct && (charge == -1 || charge == 1) {
			newCharge = 2 * charge
		}
	}
	if newCharge != charge {
		v.SetTile(x, y, tile.FromCharge(newCharge))
	}
	return newCharge - charge
}
