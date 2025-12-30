package game

const (
	WarningCharge  = 8
	DangerCharge   = 14
	OverloadCharge = 20
)

func (g *Game) DrawElectricStatus() {
	if !g.Options.IsElectric() {
		return
	}
	v := g.Video
	charge := g.LevelState.PillState.NetCharge
	shift := max(-1, min(float64(charge)/OverloadCharge, 1))
	v.SetChromaShift(shift)
	absCharge := charge
	if absCharge < 0 {
		absCharge = -absCharge
	}
	switch {
	case absCharge >= OverloadCharge:
		if g.LevelState.FrameCounter&8 == 0 {
			v.WriteAlert(" FATAL ", charge*5)
		} else {
			v.WriteAlert("       ", charge*5)
		}
	case absCharge >= DangerCharge:
		if g.LevelState.FrameCounter&16 == 0 {
			v.WriteAlert("DANGER ", charge*5)
		} else {
			v.WriteAlert("       ", charge*5)
		}
	case absCharge >= WarningCharge:
		if g.LevelState.FrameCounter&32 == 0 {
			v.WriteAlert("WARNING", charge*5)
		} else {
			v.WriteAlert("       ", charge*5)
		}
	default:
		v.WriteAlert("NORMAL ", charge*5)
	}
}

func (g *Game) ElectricOverload() bool {
	if !g.Options.IsElectric() {
		return false
	}
	absCharge := g.LevelState.PillState.NetCharge
	if absCharge < 0 {
		absCharge = -absCharge
	}
	return absCharge >= OverloadCharge
}
