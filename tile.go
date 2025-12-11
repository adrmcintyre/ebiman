package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
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

func (ds *DotState) SavePellets(v *Video) {
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

func (b *BonusActor) DrawBonus(v *Video, bonusInfo data.BonusInfoEntry) {
	look := bonusInfo.Sprite
	pal := bonusInfo.Pal
	m := &b.Motion

	if m.Visible {
		v.AddSprite(m.X-4, m.Y-4-MAZE_TOP, look, pal)
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

func (g *Game) RenderFrame() {
	g.DrawSprites()
	g.Video.FlashPills()
	g.FlashPlayerUp()
}

func (g *Game) AnimReady(frame int) (nextFrame int, delay int) {
	next := frame + 1
	v := &g.Video

	switch frame {
	case 0:
		v.SetCursor(9, 14)
		if g.PlayerNumber == 0 {
			v.WriteString("PLAYER ONE", palette.INKY)
		} else {
			v.WriteString("PLAYER TWO", palette.INKY)
		}

		v.SetCursor(9, 20)
		v.WriteString("  READY!  ", palette.PACMAN)

		// at this point sprites should be hidden
		//v.RenderMaze()
		//render_status(true)
		return next, 2 * data.FPS

	case 1:
		v.SetCursor(9, 14)
		v.WriteString("          ", palette.BLACK)
		//RenderMaze()
		g.RenderFrame()
		return next, 2 * data.FPS

	case 2:
		v.SetCursor(9, 20)
		v.WriteString("          ", palette.BLACK)
		//RenderMaze()
		return 0, 0

	default:
		return 0, 0
	}
}

func (g *Game) AnimGameOver() {
	v := &g.Video
	v.SetCursor(9, 20)
	v.WriteString("GAME  OVER", palette.PAL_29) // red

	//RenderMaze()
	//RenderStatus(true)
	delay(2000)
}

func (g *Game) AnimPacmanDie() {
	pal := palette.PACMAN

	/*
	  Timing units: 120 = 1 second

	  120 sprite 0x34 SPRITE_PACMAN_DEAD1 + hide ghosts
	  180 sprite 0x35 SPRITE_PACMAN_DEAD2 + start dying audio
	  195 sprite 0x36 SPRITE_PACMAN_DEAD3
	  210 sprite 0x37 SPRITE_PACMAN_DEAD4
	  225 sprite 0x38 SPRITE_PACMAN_DEAD5
	  240 sprite 0x39 SPRITE_PACMAN_DEAD6
	  255 sprite 0x3a SPRITE_PACMAN_DEAD7
	  270 sprite 0x3b SPRITE_PACMAN_DEAD8
	  285 sprite 0x3c SPRITE_PACMAN_DEAD9
	  300 sprite 0x3d SPRITE_PACMAN_DEAD10
	  315 sprite 0x3e SPRITE_PACMAN_DEAD11 + clear sound / pop sound (pacman)
	  345 sprite 0x3f SPRITE_PACMAN_DEAD12
	  440 done - decrement lives at this point
	*/
	lastTicks := 0
	for step := range 12 {
		var ticks int
		switch step {
		case 0:
			// TODO everything should continues to animate,
			// but ghosts and pacman stop moving
			ticks = 120

		case 1:
			// hide all ghosts and pills (and fruit)
			for j := range 4 {
				g.Ghosts[j].Motion.Visible = false
			}
			g.HideBonus()
			//g.EraseSprites()
			ticks = 180

		case 2:
			ticks = 195
		case 3:
			ticks = 210
		case 4:
			ticks = 225
		case 5:
			ticks = 240
		case 6:
			ticks = 255
		case 7:
			ticks = 270
		case 8:
			ticks = 285
		case 9:
			ticks = 300
		case 10:
			ticks = 315
		case 11:
			ticks = 345
		case 12:
			ticks = 440
		}

		delay(25 * (ticks - lastTicks) / 3)
		lastTicks = ticks

		if step < 12 {
			_ = pal
			// TODO
			//blit_sprite(pacman.motion.x-4, pacman.motion.y-4, SPRITE_PACMAN_DEAD1+j, pal)
		}
	}
}

// TODO
func delay(x int) {}
