package game

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/tile"
)

func (g *Game) DrawGhosts() {
	for j := range 4 {
		g.Ghosts[j].DrawGhost(&g.Video, g.LevelState.IsWhite, g.LevelState.FrameCounter&8 > 0)
	}
}

func (g *Game) DrawSprites() {
	g.Video.ClearSprites()
	g.BonusActor.DrawBonus(&g.Video, g.LevelConfig.BonusInfo)
	if g.LevelState.BlueTimeout == 0 {
		g.Pacman.DrawPacman(&g.Video, g.PlayerNumber)
		g.DrawGhosts()
	} else {
		g.DrawGhosts()
		g.Pacman.DrawPacman(&g.Video, g.PlayerNumber)
	}
	g.PlayerMsg.Draw(&g.Video)
	g.StatusMsg.Draw(&g.Video)
}

func (g *Game) FlashPlayerUp() {
	switch g.LevelState.FrameCounter & 31 {
	case 0:
		g.WritePlayerUp()
	case 16:
		g.ClearPlayerUp()
	}
}

func (g *Game) RenderFrameUncounted() {
	g.LevelState.WriteScores(&g.Video, g.Options.GameMode)
	g.Video.WriteLives(g.LevelState.Lives)
	g.LevelState.BonusState.WriteBonuses(&g.Video)

	if g.LevelState.BonusScoreTimeout > 0 {
		g.Video.SetCursor(12, 20)
		g.Video.WriteTiles(g.LevelConfig.BonusInfo.Tiles, color.PAL_SCORE)
	} else {
		// TODO need to avoid clearing when READY! is visible
		g.Video.SetCursor(12, 20)
		for range 4 {
			g.Video.WriteTile(tile.SPACE, color.PAL_BLACK)
		}
	}

	g.DrawSprites()
	g.Video.FlashPills()
	g.FlashPlayerUp()
}

func (g *Game) RenderFrame() {
	g.RenderFrameUncounted()
	g.LevelState.FrameCounter += 1
}
