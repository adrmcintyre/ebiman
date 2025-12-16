package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

type PacmanActor struct {
	StartPos video.ScreenPos

	StallTimer byte
	DyingFrame int
	Visible    bool
	Pos        video.ScreenPos
	Pcm        data.PCM
	TunnelPcm  data.PCM
	Vel        Velocity
}

func MakePacman() PacmanActor {
	return PacmanActor{
		StartPos: video.ScreenPos{PACMAN_START_X, PACMAN_START_Y},
	}
}

func (p *PacmanActor) Start(pcm data.PCM) {
	p.StallTimer = 0
	p.DyingFrame = 0

	p.Pos = p.StartPos
	p.Pcm = pcm
	p.TunnelPcm = 0
	p.Vel = Velocity{-1, 0}
	p.Visible = true
}

func (p *PacmanActor) Steer(v *video.Video, inDir int) {
	dir, ok := data.JoyDirection[inDir]
	if !ok {
		return
	}

	// direction can be taken if pacman is "lined up"
	if (dir.Dx == 0 && (p.Pos.X&7) == 0) || (dir.Dy == 0 && (p.Pos.Y&7) == 0) {
		tilePos := p.Pos.ToTilePos()
		nextPos := video.TilePos{
			(tilePos.X + dir.Dx + 28) % 28, // wrap left<->right (tunnel)
			tilePos.Y + dir.Dy,
		}
		nextTile := v.GetTile(nextPos)
		if IsTraversableTile(nextTile) {
			p.Vel = Velocity{dir.Dx, dir.Dy}
		}
	}
}

func (p *PacmanActor) Pulse() bool {
	return p.Pcm.Pulse()
}

func (p *PacmanActor) MovePacman(v *video.Video) {
	ok := true

	if (p.Pos.X&7) == 0 && (p.Pos.Y&7) == 0 {
		tilePos := p.Pos.ToTilePos()
		nextPos := video.TilePos{
			(tilePos.X + p.Vel.Vx + 28) % 28, // wrap left<->right (tunnel)
			tilePos.Y + p.Vel.Vy,
		}
		nextTile := v.GetTile(nextPos)
		ok = IsTraversableTile(nextTile)
	}

	if ok {
		p.Pos = video.ScreenPos{
			p.Pos.X + p.Vel.Vx,
			p.Pos.Y + p.Vel.Vy,
		}
		if p.Pos.X <= 4 && p.Vel.Vx < 0 {
			p.Pos.X += 215
		} else if p.Pos.X >= 220 && p.Vel.Vx > 0 {
			p.Pos.X -= 215
		}
	}
}

func (p *PacmanActor) DrawPacman(v *video.Video, playerNumber int) {
	var look byte
	var pal byte = palette.PACMAN

	if p.Visible {
		if playerNumber == 0 {
			pal = palette.PACMAN2
		}
		if p.DyingFrame > 0 {
			look = sprite.PACMAN_DEAD1 + byte(p.DyingFrame-1)
		} else {
			// how far into the tile are we?
			delta := (p.Pos.X + 5) % 8
			if p.Vel.Vy != 0 {
				delta = (p.Pos.Y + 5) % 8
			}
			frame := delta >> 1

			// which way are we facing?
			dir := 0
			switch {
			case p.Vel.Vx > 0:
				dir = 0
			case p.Vel.Vx < 0:
				dir = 1
			case p.Vel.Vy > 0:
				dir = 2
			case p.Vel.Vy < 0:
				dir = 3
			}
			look = PacmanAnims[dir][frame]
		}
		v.AddSprite(p.Pos.X-4, p.Pos.Y-4-MAZE_TOP, look, pal)
	}
}
