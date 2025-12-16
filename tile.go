package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

const MAZE_TOP = 16

func (ds *DotState) ResetPellets() {
	for i := range 30 {
		ds.PillBits[i] = 0xff
	}
	for i := range 4 {
		ds.PowerPills[i] = tile.POWER
	}
}

func (ds *DotState) DrawPellets(v *video.Video) {
	src := 0
	dst := 0

	// FIXME poking directly into TileRam - not very nice
	for b := range 30 {
		a := ds.PillBits[b]
		for mask := byte(0x80); mask > 0; mask >>= 1 {
			dst += int(data.Pill[src])
			src++
			if a&mask != 0 {
				v.TileRam[dst] = tile.PILL
			} else {
				v.TileRam[dst] = tile.SPACE
			}
		}
	}
	v.TileRam[3*32+4] = ds.PowerPills[0]
	v.TileRam[3*32+24] = ds.PowerPills[1]
	v.TileRam[28*32+4] = ds.PowerPills[2]
	v.TileRam[28*32+24] = ds.PowerPills[3]
}

func (ds *DotState) SavePellets(v *video.Video) {
	pillIndex := 0
	tileIndex := 0

	// FIXME peeking directly into TileRam - not very nice
	for i := range 30 {
		a := byte(0)
		for mask := byte(0x80); mask != 0; mask >>= 1 {
			tileIndex += int(data.Pill[pillIndex])
			pillIndex += 1
			if v.TileRam[tileIndex] == tile.PILL {
				a |= mask
			}
		}
		ds.PillBits[i] = a
	}

	ds.PowerPills[0] = v.TileRam[3*32+4]
	ds.PowerPills[1] = v.TileRam[3*32+24]
	ds.PowerPills[2] = v.TileRam[28*32+4]
	ds.PowerPills[3] = v.TileRam[28*32+24]
}

var PacmanAnims = [4][4]byte{
	{sprite.PACMAN_SHUT, sprite.PACMAN_RIGHT2, sprite.PACMAN_RIGHT1, sprite.PACMAN_RIGHT2},
	{sprite.PACMAN_SHUT, sprite.PACMAN_LEFT2, sprite.PACMAN_LEFT1, sprite.PACMAN_LEFT2},
	{sprite.PACMAN_SHUT, sprite.PACMAN_DOWN2, sprite.PACMAN_DOWN1, sprite.PACMAN_DOWN2},
	{sprite.PACMAN_SHUT, sprite.PACMAN_UP2, sprite.PACMAN_UP1, sprite.PACMAN_UP2},
}

// TODO - move to another file
func (g *Game) DrawGhosts() {
	for j := range 4 {
		g.Ghosts[j].DrawGhost(&g.Video, g.LevelState.IsWhite, g.LevelState.FrameCounter&8 > 0)
	}
}

func (b *BonusActor) DrawBonus(v *video.Video, bonusInfo data.BonusInfoEntry) {
	look := bonusInfo.Sprite
	pal := bonusInfo.Pal
	m := &b.Motion

	if m.Visible {
		v.AddSprite(m.Pos.X-4, m.Pos.Y-4-MAZE_TOP, look, pal)
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
}

func (g *Game) FlashPlayerUp() {
	switch g.LevelState.FrameCounter & 31 {
	case 0:
		g.WritePlayerUp(&g.Video)
	case 16:
		g.ClearPlayerUp(&g.Video)
	}
}

func (g *Game) RenderFrameUncounted() {
	g.LevelState.WriteScores(&g.Video, g.Options.GameMode)
	g.Video.WriteLives(g.LevelState.Lives)
	g.LevelState.BonusState.WriteBonuses(&g.Video)

	if g.LevelState.BonusScoreTimeout > 0 {
		g.Video.SetCursor(video.TilePos{12, 20})
		g.Video.WriteTiles(g.LevelConfig.BonusInfo.Tiles, palette.SCORE)
	} else {
		// TODO need to avoid clearing when READY! is visible
		g.Video.SetCursor(video.TilePos{12, 20})
		for range 4 {
			g.Video.WriteTile(tile.SPACE, palette.BLACK)
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
