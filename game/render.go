package game

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/tile"
)

// HideActors turns all actors' visibility off.
func (g *Game) HideActors() {
	g.Pacman.Visible = false
	for _, gh := range g.Ghosts {
		gh.Visible = false
	}
	g.BonusActor.Visible = false
}

// DrawSprites schedules all relevant sprites for drawing
func (g *Game) DrawSprites() {
	v := g.Video
	v.ClearSprites()
	g.BonusActor.Draw(v, g.LevelConfig.BonusInfo)
	if g.LevelState.BlueTimeout == 0 {
		g.Pacman.Draw(v, g.PlayerNumber)
		g.DrawGhosts()
	} else {
		g.DrawGhosts()
		g.Pacman.Draw(v, g.PlayerNumber)
	}
	// TODO - er - these aren't sprites, why is this here?
	g.PlayerMsg.Draw(v)
	g.StatusMsg.Draw(v)
}

// FlashPlayerUp maintains the flashing of 1UP / 2UP messages as appropriate.
func (g *Game) FlashPlayerUp() {
	switch g.LevelState.FrameCounter & 31 {
	case 0:
		g.WritePlayerUp()
	case 16:
		g.ClearPlayerUp()
	}
}

// WritePlayerUp places the "1UP" or "2UP" tiles for the current player.
func (g *Game) WritePlayerUp() {
	g.Video.WritePlayerUp(g.PlayerNumber)
}

// ClearPlayerUp blanks the tiles for the current player.
func (g *Game) ClearPlayerUp() {
	g.Video.ClearPlayerUp(g.PlayerNumber)
}

// RenderFrameUncounted performs all necessary status tile and sprite updates
// ready for the next frame. The frame counter is not updated.
func (g *Game) RenderFrameUncounted() {
	v := g.Video

	g.FlashPlayerUp()
	g.LevelState.WriteScores(v, g.Options.NumPlayers())

	v.FlashPills()
	v.WriteLives(g.LevelState.Lives)

	g.LevelState.BonusStatus.Write(v)
	if g.LevelState.BonusScoreTimeout > 0 {
		v.SetCursor(12, 20)
		v.WriteTiles(g.LevelConfig.BonusInfo.Tiles, color.PAL_SCORE)
	} else {
		v.SetCursor(12, 20)
		for range 4 {
			v.WriteTile(tile.SPACE, color.PAL_BLACK)
		}
	}

	g.DrawSprites()

	if g.Options.IsElectric() {
		charge := float64(g.LevelState.PillState.NetCharge) / 40.0
		shift := max(-.8, min(charge, .8))
		v.SetChromaShift(shift)
	}
}

// RenderFrameUncounted performs all necessary status tile and sprite updates
// ready for the next frame, and updates the frame counter.
func (g *Game) RenderFrame() {
	g.RenderFrameUncounted()
	g.LevelState.FrameCounter += 1
}
