package main

import (
	"math/rand"

	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/tile"
	"github.com/adrmcintyre/poweraid/video"
)

func IsTraversableTile(t byte) bool {
	switch t {
	case tile.SPACE, tile.PILL, tile.POWER, tile.POWER_SMALL:
		return true
	}
	return t >= tile.BONUS_SCORE_MIN && t <= tile.BONUS_SCORE_MAX
}

type ExitResult struct {
	Vel     Velocity
	NextPos video.TilePos
}

func (g *GhostActor) ComputeExits(v *video.Video) []ExitResult {
	// TODO: heap allocation - to avoid this the caller could supply
	// a reusable buffer to write to instead
	exits := make([]ExitResult, 0, 3)

	tilePos := g.Pos.ToTilePos()

	// anti clockwise of current heading
	vel := Velocity{g.Vel.Vy, -g.Vel.Vx}

	for range 3 {
		nextPos := video.TilePos{
			(tilePos.X + vel.Vx + 28) % 28, // wrap for tunnel
			tilePos.Y + vel.Vy,
		}
		next := v.GetTile(nextPos)

		ok := IsTraversableTile(next)
		gateOpen := g.Mode == MODE_RETURNING
		onGate := next == tile.GATE_LEFT || next == tile.GATE_RIGHT
		onHome := next == tile.HOME_LEFT || next == tile.HOME_RIGHT

		if gateOpen && (onGate || onHome) {
			// open the gate for returning ghosts
			ok = true
		} else if g.SubMode != SUBMODE_SCARED {
			// cannot turn UP at one of 4 special tiles
			goingUp := vel.Vx == 0 && vel.Vy == -1
			specialTile := (tilePos.X == 12 || tilePos.X == 15) && (tilePos.Y == 12 || tilePos.Y == 24)
			if goingUp && specialTile {
				ok = false
			}
		}

		if ok {
			exits = append(exits, ExitResult{
				Vel:     vel,
				NextPos: nextPos,
			})
		}

		// try one turn clockwise
		vel = Velocity{-vel.Vy, vel.Vx}
	}

	return exits
}

// TODO inject speeds on ghost construction?
func (g *GhostActor) Steer(v *video.Video, pacman *PacmanActor, blinky *GhostActor, speeds *data.Speeds, ghostAi bool) {
	switch g.Mode {
	case MODE_HOME:
		reachedTop := g.Vel.Vy < 0 && g.Pos.Y <= GHOST_HOME_TOP
		reachedBot := g.Vel.Vy > 0 && g.Pos.Y >= GHOST_HOME_BOTTOM
		if reachedTop || reachedBot {
			// bounce
			g.Vel.Vy = -g.Vel.Vy
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
		if g.Pos.X < GHOST_HOME_CENTRE_X {
			g.Vel = Velocity{1, 0}
		} else if g.Pos.X > GHOST_HOME_CENTRE_X {
			g.Vel = Velocity{-1, 0}
		} else if g.Pos.Y == GHOST_HOME_EXITED_Y {
			g.Mode = MODE_PLAYING
			g.Vel = Velocity{-1, 0}
			if g.SubMode == SUBMODE_SCARED {
				g.Pcm = speeds.GhostBlue
			} else {
				g.Pcm = speeds.Ghost
			}
			// TODO apply submode rules???
		} else {
			g.Vel = Velocity{0, -1}
		}
		return

	case MODE_RETURNING:
		if g.Pos == g.HomePos {
			g.Mode = MODE_HOME
			g.SetSubMode(SUBMODE_SCATTER)
			g.Pcm = data.PCM_40 // move at slowest speed when home (1 pixel every other frame)
			g.Vel = Velocity{0, -1}
			return
		}
	}

	hCentred := g.Pos.X&7 == 0
	vCentred := g.Pos.Y&7 == 0

	if !(hCentred && vCentred) {
		// take care of reversals when transitioning between tiles
		hEntering := g.Pos.X&7 == 4
		vEntering := g.Pos.Y&7 == 4
		if (hEntering && vCentred) || (vEntering && hCentred) {
			if g.ReversePending {
				g.Vel = Velocity{-g.Vel.Vx, -g.Vel.Vy}
				g.ReversePending = false
			}
		}
		return
	}

	// decision time - we're at the centre of a tile
	g.UpdateTarget(pacman, blinky)

	exits := g.ComputeExits(v)
	g.Vel = g.ChooseExitDirection(exits, ghostAi)
}

func (g *GhostActor) ChooseExitDirection(exits []ExitResult, ai bool) Velocity {
	n := len(exits)
	if n == 0 {
		return g.Vel
	}
	if n == 1 {
		return exits[0].Vel
	}

	if g.Mode == MODE_PLAYING && (g.SubMode == SUBMODE_SCARED || !ai) {
		return exits[rand.Intn(n)].Vel
	}

	bestExit := -1
	bestD2 := 32767
	for i := range n {
		if d2 := g.TargetPos.DistSq(exits[i].NextPos); d2 < bestD2 {
			// TODO - ties should be broken in order up,left,down
			bestD2 = d2
			bestExit = i
		}
	}

	return exits[bestExit].Vel
}
