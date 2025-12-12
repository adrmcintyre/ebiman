package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/tile"
)

// GAME START:
// +0 1 UP HIGH SCORE 2 UP CHARACTER / NICKNAME
// +56 sprite.blinky
// +60               SHADOW
// +30                      "BLINKY"
// +30 sprite.pinky
// +60               SPEEDY
// +30                       "PINKY"
// +30 sprite.inky
// +60               BASHFUL
// +30                       "INKY"
// +30 sprite.clyde
// +60               POKEY
// +30                       "CLYDE"
// +60 DOT 10 PTS POWER 50 PTS
// +60 POWER + COPYRIGHT
// +60 pacman appears from right persued approx 8 pixels later by BLINKY, PINKY, INKY, CLYDE
// +168 pacman consumes POWER in state CLOSED when at far left, with 0 pixels gap from ghosts, ghosts turn blue and head right
// +2 pacman has paused for two frames
// +6 pacman turns round (mouth is fully open at this point) and persues ghosts
// +32 pacman consumes first ghost 200 appears, pacman vanishes (pacman is midway between P and O of POKEY)
// +54 score vanishes pacman reappears (fully closed)
// +30 second ghost consumed (400) (midway between Y and <space>)
// +54 score vanishes pacman reappears
// +30 third ghost (800)
// +54 score vanishes pacman reappears
// +30 third ghost (1600)
// +30 empty maze with 1 PLAYER ONLY below home
// +2 dots fill maze
// +4 maze vanishes, PUSH START BUTTON, etc.

type Roster struct {
	Name string
	Nick string
	Pal  byte
}

var roster = [4]Roster{
	{"-SHADOW", "\"BLINKY\"", palette.BLINKY},
	{"-SPEEDY", "\"PINKY\"", palette.PINKY},
	{"-BASHFUL", "\"INKY\"", palette.INKY},
	{"-POKEY", "\"CLYDE\"", palette.CLYDE},
}

func (g *Game) SplashScreen(frame int) (nextFrame int, delay int) {
	v := &g.Video
	next := frame + 1

	g.LevelState.DemoMode = true

	switch frame {
	case 0:
		v.SetCursor(6, 5)
		v.WriteString("CHARACTER / NICKNAME", palette.SCORE)
		return next, 56

	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12:
		i, step := (frame-1)/3, (frame-1)%3
		y := i*3 + 8
		pal := roster[i].Pal
		switch step {
		case 0:
			tile := tile.GHOST_BASE
			for j := range 3 {
				v.SetCursor(3, y+j-1)
				for range 2 {
					v.WriteTile(tile, pal)
					tile += 1
				}
			}
			return next, 60
		case 1:
			v.SetCursor(6, y)
			v.WriteString(roster[i].Name, pal)
		case 2:
			v.SetCursor(17, y)
			v.WriteString(roster[i].Nick, pal)
		}
		return next, 30

	case 13:
		v.SetCursor(10, 23)
		v.WriteTiles([]byte{tile.PILL}, palette.MAZE)
		v.WriteTiles([]byte{tile.SPACE, tile.SCORE_1000, tile.SPACE, tile.PTS, tile.PTS + 1, tile.PTS + 2},
			palette.SCORE)

		v.SetCursor(10, 25)
		v.WriteTiles([]byte{tile.POWER}, palette.MAZE)
		v.WriteTiles([]byte{tile.SCORE_5000_1, tile.SCORE_5000_2, tile.SPACE, tile.PTS, tile.PTS + 1, tile.PTS + 2},
			palette.SCORE)
		return next, 60

	case 14:
		v.SetCursor(3, 20)
		v.WriteTile(tile.POWER, palette.MAZE)
		return next, 60

	case 15:
		g.RunningGame = true

		// where the pacman vs ghost animation occurs
		y := 20 * 8

		pm := &g.Pacman.Motion
		pm.X = 208
		pm.Y = y
		pm.Vx = -1
		pm.Vy = 0
		pm.Pcm = data.PCM_80
		pm.Visible = true

		for i := range 4 {
			ghost := &g.Ghosts[i]
			gm := &ghost.Motion
			ghost.Mode = MODE_PLAYING
			ghost.SubMode = SUBMODE_CHASE
			gm.X = pm.X + 24 + 16*i
			gm.Y = y
			gm.Vx = -1
			gm.Vy = 0
			gm.Pcm = data.PCM_85
			gm.Visible = true
		}
		return next, 0

	case 16:
		if g.LevelState.BlueTimeout == 0 {
			for i := range 4 {
				ghost := &g.Ghosts[i]
				gm := &ghost.Motion
				gm.Visible = gm.X <= 240
			}

			g.RenderFrame()
			g.UpdateState()
			g.UpdateState()

			return frame, 1
		}
		return next, 0

	case 17:
		for i := range 4 {
			g.Ghosts[i].Motion.Vx = 1
		}
		return next, 0

		// pacman continues for a few frames before turning...
	case 18, 19, 20, 21, 22, 23, 24, 25:
		g.RenderFrame()
		g.UpdateState()
		g.UpdateState()

		return next, 1

	case 26:
		g.Pacman.Motion.Vx = 1
		return next, 0

	case 27:
		if g.LevelState.GhostsEaten < 4 {
			for i := range 4 {
				g.Ghosts[i].Motion.Visible = g.Ghosts[i].Mode == MODE_PLAYING
			}

			g.RenderFrame()
			g.UpdateState()
			g.UpdateState()

			return frame, 1
		}
		g.RunningGame = false
		g.Pacman.Motion.Visible = false
		g.Ghosts[3].Motion.Visible = false
		return 0, 0

	default:
		return 0, 0
	}
}
