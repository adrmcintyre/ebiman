package game

import (
	"github.com/adrmcintyre/poweraid/audio"
	"github.com/adrmcintyre/poweraid/message"
)

// AnimReady is an coroutine for managing the "READY" prompt.
func (g *Game) AnimReady(coro *Coro) bool {
	switch coro.Step() {
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
		return coro.WaitNext(1000)

	case 1:
		g.GhostsStart()                             // reset ghosts to starting positions
		g.Pacman.Start(g.LevelConfig.Speeds.Pacman) // reset pacman to starting position
		g.BonusActor.Start()
		g.PlayerMsg = message.ClearPlayer
		return coro.WaitNext(1000)

	case 2:
		g.StatusMsg = message.ClearStatus
		return coro.Stop()

	default:
		g.PlayerMsg = message.None
		g.StatusMsg = message.None
		return coro.Stop()
	}
}

// AnimEndLevel is an coroutine for managing the end-of-level effect.
func (g *Game) AnimEndLevel(coro *Coro) bool {
	switch step := coro.Step(); step {
	case 0:
		g.Audio.StopAllTransientEffects()
		g.Audio.StopAllBackgroundEffects()
		g.Audio.StopAllPacmanEffects()
		return coro.WaitNext(1000)

	case 1:
		for _, gh := range g.Ghosts {
			gh.Visible = false
		}
		return coro.Next()

	case 2, 3, 4, 5:
		g.Video.FlashMaze(step == 2 || step == 4)
		return coro.WaitNext(250)

	default:
		return coro.Stop()
	}
}

// AnimPacmanDie is an coroutine for animating pacman's death throes.
func (g *Game) AnimPacmanDie(coro *Coro) bool {
	step := coro.Step()

	// Override pacman sprite with 12 frames of animation
	if step >= 1 && step <= 12 {
		g.Pacman.DyingFrame = step
	} else {
		g.Pacman.DyingFrame = 0
	}

	switch step {
	case 0:
		// everything continues to animate, but ghosts and pacman stop moving
		g.Audio.StopAllBackgroundEffects()
		return coro.WaitNext(1000)
	case 1:
		// hide all ghosts and pills (and fruit)
		for _, gh := range g.Ghosts {
			gh.Visible = false
		}
		g.HideBonus()
		return coro.WaitNext(500)
	case 2:
		// start dying audio
		g.Audio.PlayPacmanEffect(audio.PacmanDead)
		return coro.WaitNext(125)
	case 3, 4, 5, 6, 7, 8, 9, 10:
		return coro.WaitNext(125)
	case 11:
		// clear sound / pop sound (pacman)
		g.Audio.PlayPacmanEffect(audio.PacmanPop)
		return coro.WaitNext(250)
	case 12:
		return coro.WaitNext(800)
	default:
		return coro.Stop()
	}
}

// AnimGameOver is a coroutine for managing the "GAME OVER" prompt.
func (g *Game) AnimGameOver(coro *Coro) bool {
	switch coro.Step() {
	case 0:
		g.HideActors()
		g.StatusMsg = message.GameOver
		return coro.WaitNext(2000)

	case 1:
		g.StatusMsg = message.ClearStatus
		return coro.Next()

	default:
		g.StatusMsg = message.None
		return coro.Stop()
	}
}
