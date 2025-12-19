package ghost

import (
	"math/rand"

	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/pacman"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

// TODO inject speeds on ghost construction?
func (g *Actor) Steer(v *video.Video, pm *pacman.Actor, blinky *Actor, speeds *data.Speeds, ghostAi bool) {
	switch g.Mode {
	case MODE_HOME:
		reachedTop := g.Dir.IsUp() && g.Pos.Y <= geom.HOME_TOP
		reachedBot := g.Dir.IsDown() && g.Pos.Y >= geom.HOME_BOTTOM
		if reachedTop || reachedBot {
			// bounce
			g.Dir = g.Dir.Reverse()
		}
		return

	case MODE_LEAVING:
		//         <--+
		//            |
		// ;---------G|G--------;
		// ;   x      |         ;
		// ;   -------+         ;
		// ;          +------x  ;
		// ;--------------------;
		if g.Pos.X < geom.HOME_CENTRE.X {
			g.Dir = geom.RIGHT
		} else if g.Pos.X > geom.HOME_CENTRE.X {
			g.Dir = geom.LEFT
		} else if g.Pos.Y == geom.HOME_EXITED_Y {
			g.Mode = MODE_PLAYING
			g.Dir = geom.LEFT
			if g.SubMode == SUBMODE_SCARED {
				g.Pcm = speeds.GhostBlue
			} else {
				g.Pcm = speeds.Ghost
			}
			// TODO apply submode rules???
		} else {
			g.Dir = geom.UP
		}
		return

	case MODE_RETURNING:
		if g.Pos == g.HomePos {
			g.Mode = MODE_HOME
			g.SetSubMode(SUBMODE_SCATTER)
			g.Pcm = data.PCM_40 // move at slowest speed when home (1 pixel every other frame)
			g.Dir = geom.UP
			return
		}
	}

	// TODO could add these as utility methods
	hCentred := g.Pos.X&7 == 0
	vCentred := g.Pos.Y&7 == 0

	if !(hCentred && vCentred) {
		// take care of reversals when transitioning between tiles
		hEntering := g.Pos.X&7 == 4
		vEntering := g.Pos.Y&7 == 4
		if (hEntering && vCentred) || (vEntering && hCentred) {
			if g.ReversePending {
				g.Dir = g.Dir.Reverse()
				g.ReversePending = false
			}
		}
		return
	}

	// decision time - we're at the centre of a tile
	g.UpdateTarget(pm, blinky)

	exits := g.ComputeExits(v)
	g.Dir = g.ChooseExitDirection(exits, ghostAi)
}

type exitResult struct {
	Dir     geom.Delta
	NextPos geom.Position
}

func (g *Actor) ComputeExits(v *video.Video) []exitResult {
	// TODO: heap allocation - to avoid this the caller could supply
	// a reusable buffer to write to instead
	exits := make([]exitResult, 0, 3)

	// anti clockwise of current heading
	dir := g.Dir.TurnLeft()

	for range 3 {
		nextPos := g.Pos.Add(dir.Scale(8))
		nextTile := v.GetTile(nextPos.TileXY())

		viable := nextTile.IsTraversable()
		gateOpen := g.Mode == MODE_RETURNING
		onGate := nextTile == tile.GATE_LEFT || nextTile == tile.GATE_RIGHT
		onHome := nextTile == tile.HOME_LEFT || nextTile == tile.HOME_RIGHT

		if gateOpen && (onGate || onHome) {
			// open the gate for returning ghosts
			viable = true
		} else if g.SubMode != SUBMODE_SCARED {
			// cannot turn UP at one of 4 special tiles
			x, y := g.Pos.TileXY()
			specialTile := (x == 12 || x == 15) && (y == 12 || y == 24)
			if dir.IsUp() && specialTile {
				viable = false
			}
		}

		if viable {
			exits = append(exits, exitResult{
				Dir:     dir,
				NextPos: nextPos,
			})
		}

		// try one turn clockwise
		dir = dir.TurnRight()
	}

	return exits
}

func (g *Actor) ChooseExitDirection(exits []exitResult, ai bool) geom.Delta {
	n := len(exits)
	if n == 0 {
		return g.Dir
	}
	if n == 1 {
		return exits[0].Dir
	}

	if g.Mode == MODE_PLAYING && (g.SubMode == SUBMODE_SCARED || !ai) {
		return exits[rand.Intn(n)].Dir
	}

	bestExit := -1
	bestD2 := 32767
	for i := range n {
		if d2 := g.TargetPos.TileDistSq(exits[i].NextPos); d2 < bestD2 {
			// TODO - ties should be broken in order up,left,down
			bestD2 = d2
			bestExit = i
		}
	}

	return exits[bestExit].Dir
}
