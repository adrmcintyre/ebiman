package game

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/tile"
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
	g.Video.ClearSprites()
	g.BonusActor.Draw(&g.Video, g.LevelConfig.BonusInfo)
	if g.LevelState.BlueTimeout == 0 {
		g.Pacman.Draw(&g.Video, g.PlayerNumber)
		g.DrawGhosts()
	} else {
		g.DrawGhosts()
		g.Pacman.Draw(&g.Video, g.PlayerNumber)
	}
	// TODO - er - these aren't sprites, we is this here?
	g.PlayerMsg.Draw(&g.Video)
	g.StatusMsg.Draw(&g.Video)
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
	g.LevelState.WriteScores(&g.Video, g.Options.GameMode)
	g.Video.WriteLives(g.LevelState.Lives)
	g.LevelState.BonusStatus.Write(&g.Video)

	if g.LevelState.BonusScoreTimeout > 0 {
		g.Video.SetCursor(12, 20)
		g.Video.WriteTiles(g.LevelConfig.BonusInfo.Tiles, color.PAL_SCORE)
	} else {
		g.Video.SetCursor(12, 20)
		for range 4 {
			g.Video.WriteTile(tile.SPACE, color.PAL_BLACK)
		}
	}

	g.DrawSprites()
	g.Video.FlashPills()
	g.FlashPlayerUp()
}

// RenderFrameUncounted performs all necessary status tile and sprite updates
// ready for the next frame, and updates the frame counter.
func (g *Game) RenderFrame() {
	g.RenderFrameUncounted()
	g.LevelState.FrameCounter += 1
}
