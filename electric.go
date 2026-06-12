package main

import "github.com/adrmcintyre/ebiman/data"

func (g *Game) DrawElectricStatus() {
	if !g.Options.IsElectric() {
		return
	}
	v := g.Video
	charge := g.Player.Pills.NetCharge * data.UnitCharge

	shift := max(-1, min(float64(charge)/data.OverloadCharge, 1))
	v.SetChromaShift(shift)

	absCharge := charge
	if absCharge < 0 {
		absCharge = -absCharge
	}
	if g.IsWasmBuild {
		// glow effect causes artifacts if too intense under wasm
		v.SetPhosphorGlow(float64(absCharge)/data.OverloadCharge*0.2 + 0.2)
	} else {
		v.SetPhosphorGlow(float64(absCharge)/data.OverloadCharge*0.4 + 0.5)
	}

	switch {
	case absCharge >= data.OverloadCharge:
		if g.Level.FrameCounter&8 == 0 {
			v.WriteAlert(" FATAL ", charge)
		} else {
			v.WriteAlert("       ", charge)
		}
	case absCharge >= data.DangerCharge:
		if g.Level.FrameCounter&16 == 0 {
			v.WriteAlert("DANGER ", charge)
		} else {
			v.WriteAlert("       ", charge)
		}
	case absCharge >= data.WarningCharge:
		if g.Level.FrameCounter&32 == 0 {
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
	absCharge := g.Player.Pills.NetCharge * data.UnitCharge
	if absCharge < 0 {
		absCharge = -absCharge
	}
	return absCharge >= data.OverloadCharge
}
