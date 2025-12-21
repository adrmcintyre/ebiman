package game

import (
	"github.com/adrmcintyre/poweraid/audio"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/message"
)

// AnimReady is an coroutine for managing the "READY" prompt.
func (g *Game) AnimReady(step int) (nextStep int, delay int) {
	next := step + 1

	switch step {
	case 0:
		v := g.Video
		v.ClearTiles()
		v.ClearPalette()
		v.ColorMaze()
		v.DecodeTiles()

		g.HideActors()
		g.HideBonus()
		g.HideBonusScore()

		if g.PlayerNumber == 0 {
			v.Write1Up()
		} else {
			v.Write2Up()
		}

		g.LevelState.PillState.Draw(v)

		if g.PlayerNumber == 0 {
			g.PlayerMsg = message.Player1
		} else {
			g.PlayerMsg = message.Player2
		}

		g.StatusMsg = message.Ready

		// at this point sprites should be hidden
		return next, 1 * data.FPS

	case 1:
		g.GhostsStart()                             // reset ghosts to starting positions
		g.Pacman.Start(g.LevelConfig.Speeds.Pacman) // reset pacman to starting position
		g.BonusActor.Start()
		g.PlayerMsg = message.ClearPlayer
		return next, 1 * data.FPS

	case 2:
		g.StatusMsg = message.ClearStatus
		return 0, 0

	default:
		g.PlayerMsg = message.None
		g.StatusMsg = message.None
		return 0, 0
	}
}

// AnimEndLevel is an coroutine for managing the end-of-level effect.
func (g *Game) AnimEndLevel(step int) (nextStep int, delay int) {
	next := step + 1

	switch step {
	case 0:
		g.Audio.StopAllTransientEffects()
		g.Audio.StopAllBackgroundEffects()
		g.Audio.StopAllPacmanEffects()
		return next, data.FPS

	case 1:
		for _, gh := range g.Ghosts {
			gh.Visible = false
		}
		return next, 0

	case 2, 3, 4, 5:
		g.Video.FlashMaze(step == 2 || step == 4)
		return next, data.FPS / 4

	default:
		return 0, 0
	}
}

// AnimPacmanDie is an coroutine for animating pacman's death throes.
func (g *Game) AnimPacmanDie(step int) (int, int) {
	next := step + 1

	// Override pacman sprite with 12 frames of animation
	if step >= 1 && step <= 12 {
		g.Pacman.DyingFrame = step
	} else {
		g.Pacman.DyingFrame = 0
	}

	// Timing units: 120 = 1 second
	delay := func(t int) (int, int) {
		return next, t * data.FPS / 120
	}

	switch step {
	case 0:
		// everything continues to animate, but ghosts and pacman stop moving
		g.Audio.StopAllBackgroundEffects()
		return delay(120)
	case 1:
		// hide all ghosts and pills (and fruit)
		for _, gh := range g.Ghosts {
			gh.Visible = false
		}
		g.HideBonus()
		return delay(60)
	case 2:
		// start dying audio
		g.Audio.PlayPacmanEffect(audio.PacmanDead)
		return delay(15)
	case 3, 4, 5, 6, 7, 8, 9, 10:
		return delay(15)
	case 11:
		// clear sound / pop sound (pacman)
		g.Audio.PlayPacmanEffect(audio.PacmanPop)
		return delay(30)
	case 12:
		return delay(95)
	default:
		return 0, 0
	}
}

// AnimGameOver is a coroutine for managing the "GAME OVER" prompt.
func (g *Game) AnimGameOver(step int) (nextStep int, delay int) {
	next := step + 1

	switch step {
	case 0:
		g.HideActors()
		g.StatusMsg = message.GameOver
		return next, 2 * data.FPS

	case 1:
		g.StatusMsg = message.ClearStatus
		return next, 0

	default:
		g.StatusMsg = message.None
		return 0, 0
	}
}
