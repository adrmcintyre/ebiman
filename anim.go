package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
)

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
		return next, 1 * data.FPS

	case 1:
		v.SetCursor(9, 14)
		v.WriteString("          ", palette.BLACK)
		g.RenderFrame()
		return next, 1 * data.FPS

	case 2:
		v.SetCursor(9, 20)
		v.WriteString("          ", palette.BLACK)
		return 0, 0

	default:
		return 0, 0
	}
}

func (g *Game) AnimEndLevel(frame int) (nextFrame int, delay int) {
	next := frame + 1

	switch frame {
	case 0:
		return next, data.FPS

	case 1:
		for i := range 4 {
			g.Ghosts[i].Motion.Visible = false
		}
		return next, 0

	case 2, 4:
		g.Video.ColorMaze(2)
		return next, data.FPS / 4

	case 3, 5:
		g.Video.ColorMaze(1)
		return next, data.FPS / 4

	default:
		return 0, 0
	}
}

func (g *Game) AnimPacmanDie(frame int) (nextFrame int, delay int) {
	next := frame + 1

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

	if frame >= 1 && frame <= 12 {
		g.Pacman.DyingFrame = frame
	}

	tm := func(t1, t2 int) (int, int) {
		return next, (t2 - t1) * data.FPS / 120
	}

	switch frame {
	case 0:
		// everything continues to animate, but ghosts and pacman stop moving
		return tm(0, 120)
	case 1:
		// hide all ghosts and pills (and fruit)
		for j := range 4 {
			g.Ghosts[j].Motion.Visible = false
		}
		g.HideBonus()
		return tm(120, 180)
	case 2:
		return tm(180, 195)
	case 3:
		return tm(195, 210)
	case 4:
		return tm(210, 225)
	case 5:
		return tm(225, 240)
	case 6:
		return tm(240, 255)
	case 7:
		return tm(255, 270)
	case 8:
		return tm(270, 285)
	case 9:
		return tm(285, 300)
	case 10:
		return tm(300, 315)
	case 11:
		return tm(315, 345)
	case 12:
		return tm(345, 440)
	default:
		return 0, 0
	}
}

func (g *Game) AnimGameOver(frame int) (nextFrame int, delay int) {
	next := frame + 1
	v := &g.Video

	switch frame {
	case 0:
		v.SetCursor(9, 20)
		v.WriteString("GAME  OVER", palette.PAL_29) // red
		return next, 2 * data.FPS

	default:
		return 0, 0
	}
}
