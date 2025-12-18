package game

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/palette"
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
		StartPos: geom.Position{PACMAN_START_X, PACMAN_START_Y},
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
		if IsTraversableTile(nextTile) {
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
		ok = IsTraversableTile(nextTile)
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

var PacmanAnims = struct {
	Up, Left, Down, Right [4]byte
}{
	[4]byte{sprite.PACMAN_SHUT, sprite.PACMAN_UP2, sprite.PACMAN_UP1, sprite.PACMAN_UP2},
	[4]byte{sprite.PACMAN_SHUT, sprite.PACMAN_LEFT2, sprite.PACMAN_LEFT1, sprite.PACMAN_LEFT2},
	[4]byte{sprite.PACMAN_SHUT, sprite.PACMAN_DOWN2, sprite.PACMAN_DOWN1, sprite.PACMAN_DOWN2},
	[4]byte{sprite.PACMAN_SHUT, sprite.PACMAN_RIGHT2, sprite.PACMAN_RIGHT1, sprite.PACMAN_RIGHT2},
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
			if p.Dir.IsVertical() {
				delta = (p.Pos.Y + 5) % 8
			}
			frame := delta >> 1

			// which way are we facing?
			switch {
			case p.Dir.IsUp():
				look = PacmanAnims.Up[frame]
			case p.Dir.IsLeft():
				look = PacmanAnims.Left[frame]
			case p.Dir.IsDown():
				look = PacmanAnims.Down[frame]
			case p.Dir.IsRight():
				look = PacmanAnims.Right[frame]
			}
		}
		offset := geom.Delta{-4, -4 - MAZE_TOP}
		v.AddSprite(p.Pos.Add(offset), look, pal)
	}
}
