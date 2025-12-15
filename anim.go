package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/video"
)

func (g *Game) AnimReady(frame int) (nextFrame int, delay int) {
	next := frame + 1
	v := &g.Video

	switch frame {
	case 0:
		v.SetCursor(video.TilePos{9, 14})
		if g.PlayerNumber == 0 {
			v.WriteString("PLAYER ONE", palette.INKY)
		} else {
			v.WriteString("PLAYER TWO", palette.INKY)
		}

		v.SetCursor(video.TilePos{9, 20})
		v.WriteString("  READY!  ", palette.PACMAN)

		// at this point sprites should be hidden
		return next, 1 * data.FPS

	case 1:
		v.SetCursor(video.TilePos{9, 14})
		v.WriteString("          ", palette.BLACK)
		g.RenderFrame()
		return next, 1 * data.FPS

	case 2:
		v.SetCursor(video.TilePos{9, 20})
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

	case 2, 3, 4, 5:
		g.Video.FlashMaze(frame == 2 || frame == 4)
		return next, data.FPS / 4

	default:
		return 0, 0
	}
}

func (g *Game) AnimPacmanDie(frame int) (int, int) {
	next := frame + 1

	// Override pacman sprite with 12 frames of animation
	if frame >= 1 && frame <= 12 {
		g.Pacman.DyingFrame = frame
	}

	// Timing units: 120 = 1 second
	delay := func(t int) (int, int) {
		return next, t * data.FPS / 120
	}

	switch frame {
	case 0:
		// everything continues to animate, but ghosts and pacman stop moving
		return delay(120)
	case 1:
		// hide all ghosts and pills (and fruit)
		for j := range 4 {
			g.Ghosts[j].Motion.Visible = false
		}
		g.HideBonus()
		return delay(60)
	case 2:
		// start dying audio
		return delay(15)
	case 3, 4, 5, 6, 7, 8, 9, 10:
		return delay(15)
	case 11:
		// clear sound / pop sound (pacman)
		return delay(30)
	case 12:
		return delay(95)
	default:
		return 0, 0
	}
}

func (g *Game) AnimGameOver(frame int) (nextFrame int, delay int) {
	next := frame + 1
	v := &g.Video

	switch frame {
	case 0:
		v.SetCursor(video.TilePos{9, 20})
		v.WriteString("GAME  OVER", palette.PAL_29) // red
		return next, 2 * data.FPS

	default:
		return 0, 0
	}
}
