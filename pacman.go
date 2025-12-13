package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

type PacmanActor struct {
	StartX, StartY int
	Motion         Motion
	StallTimer     byte
	DyingFrame     int
}

func MakePacman() PacmanActor {
	return PacmanActor{
		StartX: PACMAN_START_X,
		StartY: PACMAN_START_Y,
	}
}

func (p *PacmanActor) Start(pcm uint32) {
	p.StallTimer = 0
	p.DyingFrame = 0

	m := &p.Motion
	m.X = p.StartX
	m.Y = p.StartY
	m.Pcm = pcm
	m.TunnelPcm = 0
	m.Vx = -1
	m.Vy = 0
	m.Visible = true
}

func (p *PacmanActor) SteerPacman(v *video.Video, inDir int) {
	m := &p.Motion

	dir, ok := data.JoyDirection[inDir]
	if !ok {
		return
	}
	x := m.X
	y := m.Y
	dx := dir.Dx
	dy := dir.Dy

	// direction can be taken if pacman is "lined up"
	if (dx == 0 && (x&7) == 0) || (dy == 0 && (y&7) == 0) {
		tileX := x / 8
		tileY := y / 8
		nextX := (tileX + dx + 28) % 28 // wrap left<->right (tunnel)
		nextY := tileY + dy
		next := v.GetTile(nextX, nextY)
		if IsTraversableTile(next) {
			m.Vx = dx
			m.Vy = dy
		}
	}
}

func (p *PacmanActor) Pulse() bool {
	m := &p.Motion
	pcm := m.Pcm
	msb := pcm >> 31
	m.Pcm = (pcm << 1) | msb
	return msb != 0
}

func (p *PacmanActor) MovePacman(v *video.Video) {
	m := &p.Motion

	vx := m.Vx
	vy := m.Vy

	x := m.X
	y := m.Y
	ok := true

	if (x&7) == 0 && (y&7) == 0 {
		nextX := (x/8 + vx + 28) % 28 // wrap left<->right (tunnel)
		nextY := y/8 + vy
		next := v.GetTile(nextX, nextY)
		ok = IsTraversableTile(next)
	}

	if ok {
		x += vx
		y += vy
		if x <= 4 && vx < 0 {
			x += 215
		} else if x >= 220 && vx > 0 {
			x -= 215
		}

		m.X = x
		m.Y = y
	}
}

func (pm *PacmanActor) DrawPacman(v *video.Video, playerNumber int) {
	var look byte
	var pal byte = palette.PACMAN

	m := &pm.Motion
	if m.Visible {
		if playerNumber == 0 {
			pal = palette.PACMAN2
		}
		if pm.DyingFrame > 0 {
			look = sprite.PACMAN_DEAD1 + byte(pm.DyingFrame-1)
		} else {
			pos := m.X
			if m.Vy != 0 {
				pos = m.Y
			}
			j := ((pos + 5) & 7) >> 1
			dir := 0
			switch {
			case m.Vx > 0:
				dir = 0
			case m.Vx < 0:
				dir = 1
			case m.Vy > 0:
				dir = 2
			case m.Vy < 0:
				dir = 3
			}
			look = PacmanAnims[dir][j]
		}
		v.AddSprite(m.X-4, m.Y-4-MAZE_TOP, look, pal)
	}
}
