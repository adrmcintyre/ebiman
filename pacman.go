package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

type PacmanActor struct {
	StartPos   video.ScreenPos
	Motion     Motion
	StallTimer byte
	DyingFrame int
}

func MakePacman() PacmanActor {
	return PacmanActor{
		StartPos: video.ScreenPos{PACMAN_START_X, PACMAN_START_Y},
	}
}

func (p *PacmanActor) Start(pcm uint32) {
	p.StallTimer = 0
	p.DyingFrame = 0

	m := &p.Motion
	m.Pos = p.StartPos
	m.Pcm = pcm
	m.TunnelPcm = 0
	m.Vel = Velocity{-1, 0}
	m.Visible = true
}

func (p *PacmanActor) SteerPacman(v *video.Video, inDir int) {
	m := &p.Motion

	dir, ok := data.JoyDirection[inDir]
	if !ok {
		return
	}
	pos := m.Pos

	// direction can be taken if pacman is "lined up"
	if (dir.Dx == 0 && (pos.X&7) == 0) || (dir.Dy == 0 && (pos.Y&7) == 0) {
		tilePos := pos.ToTilePos()
		nextPos := video.TilePos{
			(tilePos.X + dir.Dx + 28) % 28, // wrap left<->right (tunnel)
			tilePos.Y + dir.Dy,
		}
		nextTile := v.GetTile(nextPos)
		if IsTraversableTile(nextTile) {
			m.Vel = Velocity{dir.Dx, dir.Dy}
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

	vel := m.Vel

	ok := true
	pos := m.Pos

	if (m.Pos.X&7) == 0 && (m.Pos.Y&7) == 0 {
		tilePos := pos.ToTilePos()
		nextPos := video.TilePos{
			(tilePos.X + vel.Vx + 28) % 28, // wrap left<->right (tunnel)
			tilePos.Y + vel.Vy,
		}
		nextTile := v.GetTile(nextPos)
		ok = IsTraversableTile(nextTile)
	}

	if ok {
		m.Pos = video.ScreenPos{
			pos.X + vel.Vx,
			pos.Y + vel.Vy,
		}
		if m.Pos.X <= 4 && vel.Vx < 0 {
			m.Pos.X += 215
		} else if m.Pos.X >= 220 && vel.Vx > 0 {
			m.Pos.X -= 215
		}
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
			// how far into the tile are we?
			delta := (m.Pos.X + 5) % 8
			if m.Vel.Vy != 0 {
				delta = (m.Pos.Y + 5) % 8
			}
			frame := delta >> 1

			// which way are we facing?
			dir := 0
			switch {
			case m.Vel.Vx > 0:
				dir = 0
			case m.Vel.Vx < 0:
				dir = 1
			case m.Vel.Vy > 0:
				dir = 2
			case m.Vel.Vy < 0:
				dir = 3
			}
			look = PacmanAnims[dir][frame]
		}
		v.AddSprite(m.Pos.X-4, m.Pos.Y-4-MAZE_TOP, look, pal)
	}
}
