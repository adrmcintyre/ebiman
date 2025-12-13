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
	Vx, Vy       int
	NextX, NextY int
}

func (g *GhostActor) ComputeExits(v *video.Video) []ExitResult {
	// TODO: heap allocation - to avoid this the caller could supply
	// a reusable buffer to write to instead
	exits := make([]ExitResult, 0, 3)
	m := &g.Motion

	tileX := m.X / 8
	tileY := m.Y / 8

	// anti clockwise of current heading
	vx := m.Vy
	vy := -m.Vx

	for range 3 {
		nextX := tileX + vx
		nextY := tileY + vy
		// tunnel
		if nextX == -1 {
			nextX = 27
		}
		if nextX == 28 {
			nextX = 0
		}
		next := v.GetTile(nextX, nextY)

		ok := IsTraversableTile(next)
		gateOpen := g.Mode == MODE_RETURNING
		onGate := next == tile.GATE_LEFT || next == tile.GATE_RIGHT
		onHome := next == tile.HOME_LEFT || next == tile.HOME_RIGHT

		if gateOpen && (onGate || onHome) {
			// open the gate for returning ghosts
			ok = true
		} else if g.SubMode != SUBMODE_SCARED {
			// cannot turn UP at one of 4 special tiles
			goingUp := vx == 0 && vy == -1
			specialTile := (tileX == 12 || tileX == 15) && (tileY == 12 || tileY == 24)
			if goingUp && specialTile {
				ok = false
			}
		}

		if ok {
			exits = append(exits, ExitResult{
				Vx:    vx,
				Vy:    vy,
				NextX: nextX,
				NextY: nextY,
			})
		}

		// try one turn clockwise
		vx, vy = -vy, vx
	}

	return exits
}

// TODO inject speeds on ghost construction?
func (g *GhostActor) SteerGhost(v *video.Video, pacman *PacmanActor, blinky *GhostActor, speeds *data.Speeds, ghostAi bool) {
	m := &g.Motion

	switch g.Mode {
	case MODE_HOME:
		reachedTop := m.Vy < 0 && m.Y <= GHOST_HOME_TOP
		reachedBot := m.Vy > 0 && m.Y >= GHOST_HOME_BOTTOM
		if reachedTop || reachedBot {
			// bounce
			m.Vy = -m.Vy
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
		if m.X < GHOST_HOME_CENTRE_X {
			m.Vx = 1
			m.Vy = 0
		} else if m.X > GHOST_HOME_CENTRE_X {
			m.Vx = -1
			m.Vy = 0
		} else if m.Y == GHOST_HOME_EXITED_Y {
			g.Mode = MODE_PLAYING
			m.Vx = -1
			m.Vy = 0
			if g.SubMode == SUBMODE_SCARED {
				m.Pcm = speeds.GhostBlue
			} else {
				m.Pcm = speeds.Ghost
			}
			// TODO apply submode rules???
		} else {
			m.Vx = 0
			m.Vy = -1
		}
		return

	case MODE_RETURNING:
		if m.X == g.HomeX && m.Y == g.HomeY {
			g.Mode = MODE_HOME
			g.GhostSetSubmode(SUBMODE_SCATTER)
			m.Pcm = data.PCM_40 // move at slowest speed when home (1 pixel every other frame)
			m.Vx = 0
			m.Vy = -1
			return
		}
	}

	hCentred := m.X&7 == 0
	vCentred := m.Y&7 == 0

	if !(hCentred && vCentred) {
		// take care of reversals when transitioning between tiles
		hEntering := m.X&7 == 4
		vEntering := m.Y&7 == 4
		if (hEntering && vCentred) || (vEntering && hCentred) {
			if g.ReversePending {
				m.Vx = -m.Vx
				m.Vy = -m.Vy
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
		m.Vx = exits[0].Vx
		m.Vy = exits[0].Vy
		return
	}

	if g.Mode == MODE_PLAYING && (g.SubMode == SUBMODE_SCARED || !ghostAi) {
		i := rand.Intn(n)
		m.Vx = exits[i].Vx
		m.Vy = exits[i].Vy
		return
	}

	bestExit := -1
	bestD2 := 32767
	for i := range n {
		dx := g.TargetX - exits[i].NextX
		dy := g.TargetY - exits[i].NextY
		d2 := dx*dx + dy*dy
		if d2 < bestD2 {
			// TODO - ties should be broken in order up,left,down
			bestD2 = d2
			bestExit = i
		}
	}

	m.Vx = exits[bestExit].Vx
	m.Vy = exits[bestExit].Vy
}
