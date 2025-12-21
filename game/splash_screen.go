package game

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/ghost"
	"github.com/adrmcintyre/poweraid/option"
	"github.com/adrmcintyre/poweraid/tile"
)

// Splash Screen Cue Sheet
//
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

// RosterEntry lists an antagonist
type RosterEntry struct {
	Name string
	Nick string
	Pal  color.Palette
}

// roster collects all of the ghosts
var roster = [4]RosterEntry{
	{"-SHADOW", "\"BLINKY\"", color.PAL_BLINKY},
	{"-SPEEDY", "\"PINKY\"", color.PAL_PINKY},
	{"-BASHFUL", "\"INKY\"", color.PAL_INKY},
	{"-POKEY", "\"CLYDE\"", color.PAL_CLYDE},
}

// SplashScreen is an animator coroutine for the splash screen.
func (g *Game) SplashScreen(frame int) (nextFrame int, delay int) {
	v := &g.Video
	next := frame + 1

	g.LevelState.DemoMode = true

	switch frame {
	case 0:
		g.LevelConfig.Init(0, option.DIFFICULTY_MEDIUM)
		g.LevelState.Init(0)
		// TODO - ResetPlayer?
		g.LevelState.ClearScores()
		g.LevelState.BonusStatus.ClearBonuses()
		g.LevelState.LevelStart()

		g.Audio.Mute()
		g.HideActors()

		v.ClearTiles()
		v.ClearPalette()
		v.ColorMaze()
		v.Write1Up()
		v.SetCursor(6, 5)
		v.WriteString("CHARACTER / NICKNAME", color.PAL_SCORE)
		g.LevelState.WriteScores(v, option.GAME_MODE_1P)
		return next, 56

	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12:
		i, step := (frame-1)/3, (frame-1)%3
		y := 8 + i*3
		pal := roster[i].Pal
		switch step {
		case 0:
			t := tile.GHOST_BASE
			for j := range 3 {
				v.SetCursor(3, y+j-1)
				for range 2 {
					v.WriteTile(t, pal)
					t += 1
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
		v.WriteTiles([]tile.Tile{tile.PILL}, color.PAL_MAZE)
		v.WriteTiles([]tile.Tile{tile.SPACE, tile.SCORE_1000, tile.SPACE, tile.PTS, tile.PTS + 1, tile.PTS + 2},
			color.PAL_SCORE)

		v.SetCursor(10, 25)
		v.WriteTiles([]tile.Tile{tile.POWER}, color.PAL_MAZE)
		v.WriteTiles([]tile.Tile{tile.SCORE_5000_1, tile.SCORE_5000_2, tile.SPACE, tile.PTS, tile.PTS + 1, tile.PTS + 2},
			color.PAL_SCORE)
		return next, 60

	case 14:
		v.SetCursor(3, 20)
		v.WriteTile(tile.POWER, color.PAL_MAZE)

		return next, 60

	case 15:
		g.RunningGame = true

		// where the pacman vs ghost animation occurs
		y := 20

		p := g.Pacman
		p.Pos = geom.TilePos(26, y)
		p.Dir = geom.LEFT
		p.Pcm = data.PCM_80
		p.Visible = true

		for i, gh := range g.Ghosts {
			gh.Mode = ghost.MODE_PLAYING
			gh.SubMode = ghost.SUBMODE_CHASE
			gh.Visible = true
			gh.Pos = geom.TilePos(p.Pos.TileX()+3+2*i, y)
			gh.Dir = geom.LEFT
			gh.Pcm = data.PCM_85
		}
		return next, 0

	case 16:
		if g.LevelState.BlueTimeout == 0 {
			endPos := geom.TilePos(30, 20)
			for _, gh := range g.Ghosts {
				gh.Visible = gh.Pos.X <= endPos.X
			}

			g.RenderFrame()
			for range 4 {
				g.UpdateState()
			}

			return frame, 1
		}
		return next, 0

	case 17:
		for _, gh := range g.Ghosts {
			gh.Dir = geom.RIGHT
		}
		return next, 0

		// pacman continues for a few frames before turning...
	case 18, 19, 20, 21:
		g.RenderFrame()
		for range 4 {
			g.UpdateState()
		}

		return next, 1

	case 22:
		g.Pacman.Dir = geom.RIGHT
		return next, 0

	case 23:
		if g.LevelState.GhostsEaten < 4 {
			for _, gh := range g.Ghosts {
				gh.Visible = gh.Mode == ghost.MODE_PLAYING
			}

			g.RenderFrame()
			for range 4 {
				g.UpdateState()
			}

			return frame, 1
		}
		g.RunningGame = false
		g.Pacman.Visible = false
		g.Ghosts[3].Visible = false
		return 0, 0

	default:
		return 0, 0
	}
}
