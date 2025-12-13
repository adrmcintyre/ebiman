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
	NextPos Position
}

func (g *GhostActor) ComputeExits(v *video.Video) []ExitResult {
	// TODO: heap allocation - to avoid this the caller could supply
	// a reusable buffer to write to instead
	exits := make([]ExitResult, 0, 3)
	m := &g.Motion

	tilePos := Position{m.Pos.X / 8, m.Pos.Y / 8}

	// anti clockwise of current heading
	vel := Velocity{m.Vel.Vy, -m.Vel.Vx}

	for range 3 {
		nextPos := Position{
			(tilePos.X + vel.Vx + 28) % 28, // wrap for tunnel
			tilePos.Y + vel.Vy,
		}
		next := v.GetTile(nextPos.X, nextPos.Y)

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
func (g *GhostActor) SteerGhost(v *video.Video, pacman *PacmanActor, blinky *GhostActor, speeds *data.Speeds, ghostAi bool) {
	m := &g.Motion

	switch g.Mode {
	case MODE_HOME:
		reachedTop := m.Vel.Vy < 0 && m.Pos.Y <= GHOST_HOME_TOP
		reachedBot := m.Vel.Vy > 0 && m.Pos.Y >= GHOST_HOME_BOTTOM
		if reachedTop || reachedBot {
			// bounce
			m.Vel.Vy = -m.Vel.Vy
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
		if m.Pos.X < GHOST_HOME_CENTRE_X {
			m.Vel = Velocity{1, 0}
		} else if m.Pos.X > GHOST_HOME_CENTRE_X {
			m.Vel = Velocity{-1, 0}
		} else if m.Pos.Y == GHOST_HOME_EXITED_Y {
			g.Mode = MODE_PLAYING
			m.Vel = Velocity{-1, 0}
			if g.SubMode == SUBMODE_SCARED {
				m.Pcm = speeds.GhostBlue
			} else {
				m.Pcm = speeds.Ghost
			}
			// TODO apply submode rules???
		} else {
			m.Vel = Velocity{0, -1}
		}
		return

	case MODE_RETURNING:
		if m.Pos == g.HomePos {
			g.Mode = MODE_HOME
			g.GhostSetSubmode(SUBMODE_SCATTER)
			m.Pcm = data.PCM_40 // move at slowest speed when home (1 pixel every other frame)
			m.Vel = Velocity{0, -1}
			return
		}
	}

	hCentred := m.Pos.X&7 == 0
	vCentred := m.Pos.Y&7 == 0

	if !(hCentred && vCentred) {
		// take care of reversals when transitioning between tiles
		hEntering := m.Pos.X&7 == 4
		vEntering := m.Pos.Y&7 == 4
		if (hEntering && vCentred) || (vEntering && hCentred) {
			if g.ReversePending {
				m.Vel = Velocity{-m.Vel.Vx, -m.Vel.Vy}
				g.ReversePending = false
			}
		}
		return
	}

	// TODO - we could split this out into a separate function

	// decision time - we're at the centre of a tile

	g.UpdateTarget(pacman, blinky)

	exits := g.ComputeExits(v)
	n := len(exits)
	if n == 0 {
		return
	}
	if n == 1 {
		m.Vel = exits[0].Vel
		return
	}

	if g.Mode == MODE_PLAYING && (g.SubMode == SUBMODE_SCARED || !ghostAi) {
		i := rand.Intn(n)
		m.Vel = exits[i].Vel
		return
	}

	bestExit := -1
	bestD2 := 32767
	for i := range n {
		dx := g.TargetPos.X - exits[i].NextPos.X
		dy := g.TargetPos.Y - exits[i].NextPos.Y
		d2 := dx*dx + dy*dy
		if d2 < bestD2 {
			// TODO - ties should be broken in order up,left,down
			bestD2 = d2
			bestExit = i
		}
	}

	m.Vel = exits[bestExit].Vel
}
