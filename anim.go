package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/message"
)

func (g *Game) AnimReady(frame int) (nextFrame int, delay int) {
	next := frame + 1

	switch frame {
	case 0:
		if g.PlayerNumber == 0 {
			g.PlayerMsg = message.Player1
		} else {
			g.PlayerMsg = message.Player2
		}

		g.StatusMsg = message.Ready

		// at this point sprites should be hidden
		return next, 1 * data.FPS

	case 1:
		g.PlayerMsg = message.NoPlayer
		return next, 1 * data.FPS

	case 2:
		g.StatusMsg = message.NoStatus
		return 0, 0

	default:
		g.PlayerMsg = message.None
		g.StatusMsg = message.None
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
			g.Ghosts[i].Visible = false
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
	} else {
		g.Pacman.DyingFrame = 0
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
			g.Ghosts[j].Visible = false
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

	switch frame {
	case 0:
		g.StatusMsg = message.GameOver
		return next, 2 * data.FPS

	case 1:
		g.StatusMsg = message.NoStatus
		return next, 0

	default:
		g.StatusMsg = message.None
		return 0, 0
	}
}
