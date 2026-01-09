package game

const (
	UnitCharge     = 5
	WarningCharge  = 45
	DangerCharge   = 70
	OverloadCharge = 95
)

func (g *Game) DrawElectricStatus() {
	if !g.Options.IsElectric() {
		return
	}
	v := g.Video
	charge := g.LevelState.PillState.NetCharge * UnitCharge

	shift := max(-1, min(float64(charge)/OverloadCharge, 1))
	v.SetChromaShift(shift)

	absCharge := charge
	if absCharge < 0 {
		absCharge = -absCharge
	}
	if g.IsWasmBuild {
		// glow effect causes artifacts if too intense under wasm
		v.SetPhosphorGlow(float64(absCharge)/OverloadCharge*0.2 + 0.2)
	} else {
		v.SetPhosphorGlow(float64(absCharge)/OverloadCharge*0.4 + 0.5)
	}

	switch {
	case absCharge >= OverloadCharge:
		if g.LevelState.FrameCounter&8 == 0 {
			v.WriteAlert(" FATAL ", charge)
		} else {
			v.WriteAlert("       ", charge)
		}
	case absCharge >= DangerCharge:
		if g.LevelState.FrameCounter&16 == 0 {
			v.WriteAlert("DANGER ", charge)
		} else {
			v.WriteAlert("       ", charge)
		}
	case absCharge >= WarningCharge:
		if g.LevelState.FrameCounter&32 == 0 {
			v.WriteAlert("WARNING", charge)
		} else {
			v.WriteAlert("       ", charge)
		}
	default:
		v.WriteAlert("NORMAL ", charge)
	}
}

func (g *Game) ElectricOverload() bool {
	if !g.Options.IsElectric() {
		return false
	}
	absCharge := g.LevelState.PillState.NetCharge * UnitCharge
	if absCharge < 0 {
		absCharge = -absCharge
	}
	return absCharge >= OverloadCharge
}
