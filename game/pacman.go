package game

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

type PacmanActor struct {
	// configuration fields
	StartPos geom.Position

	// state fields
	Visible    bool
	Pos        geom.Position
	Dir        geom.Delta
	Pcm        data.PCM
	TunnelPcm  data.PCM
	StallTimer int
	DyingFrame int
}

func MakePacman() PacmanActor {
	return PacmanActor{
		StartPos: PACMAN_START,
	}
}

func (p *PacmanActor) Start(pcm data.PCM) {
	p.Visible = true
	p.Pos = p.StartPos
	p.Dir = geom.LEFT
	p.Pcm = pcm
	p.TunnelPcm = 0
	p.StallTimer = 0
	p.DyingFrame = 0

}

func (p *PacmanActor) Steer(v *video.Video, inDir int) {
	dir, ok := input.JoyDirection[inDir]
	if !ok {
		return
	}

	// direction can be taken if pacman is "lined up"
	if (dir.IsVertical() && (p.Pos.X&7) == 0) || (dir.IsHorizontal() && (p.Pos.Y&7) == 0) {
		nextPos := p.Pos.Add(dir.Scale(8)).WrapTunnel()
		nextTile := v.GetTile(nextPos.TileXY())
		if nextTile.IsTraversable() {
			p.Dir = dir
		}
	}
}

func (p *PacmanActor) Pulse() bool {
	return p.Pcm.Pulse()
}

func (p *PacmanActor) MovePacman(v *video.Video) {
	ok := true

	if (p.Pos.X&7) == 0 && (p.Pos.Y&7) == 0 {
		nextPos := p.Pos.Add(p.Dir.Scale(8)).WrapTunnel()
		nextTile := v.GetTile(nextPos.TileXY())
		ok = nextTile.IsTraversable()
	}

	if ok {
		p.Pos = p.Pos.Add(p.Dir)
		if p.Pos.X <= 4 && p.Dir.IsLeft() {
			p.Pos.X += 215
		} else if p.Pos.X >= 220 && p.Dir.IsRight() {
			p.Pos.X -= 215
		}
	}
}

var pacmanAnims = struct {
	Up, Left, Down, Right [4]sprite.Look
}{
	[4]sprite.Look{sprite.PACMAN_SHUT, sprite.PACMAN_UP2, sprite.PACMAN_UP1, sprite.PACMAN_UP2},
	[4]sprite.Look{sprite.PACMAN_SHUT, sprite.PACMAN_LEFT2, sprite.PACMAN_LEFT1, sprite.PACMAN_LEFT2},
	[4]sprite.Look{sprite.PACMAN_SHUT, sprite.PACMAN_DOWN2, sprite.PACMAN_DOWN1, sprite.PACMAN_DOWN2},
	[4]sprite.Look{sprite.PACMAN_SHUT, sprite.PACMAN_RIGHT2, sprite.PACMAN_RIGHT1, sprite.PACMAN_RIGHT2},
}

func (p *PacmanActor) DrawPacman(v *video.Video, playerNumber int) {
	if p.Visible {
		var pal = color.PAL_PACMAN
		if playerNumber == 1 {
			pal = color.PAL_PACMAN2
		}

		var look sprite.Look
		if p.DyingFrame > 0 {
			look = sprite.PACMAN_DEAD1 + sprite.Look(p.DyingFrame-1)
		} else {
			// how far into the tile are we?
			delta := (p.Pos.X + 5) % 8
			if p.Dir.IsVertical() {
				delta = (p.Pos.Y + 5) % 8
			}
			frame := delta >> 1

			// which way are we facing?
			switch {
			case p.Dir.IsUp():
				look = pacmanAnims.Up[frame]
			case p.Dir.IsLeft():
				look = pacmanAnims.Left[frame]
			case p.Dir.IsDown():
				look = pacmanAnims.Down[frame]
			case p.Dir.IsRight():
				look = pacmanAnims.Right[frame]
			}
		}
		v.AddSprite(p.Pos, look, pal)
	}
}
