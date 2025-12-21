package pacman

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

// An Actor describes the state and look of pacman
type Actor struct {
	// configuration fields
	StartPos geom.Position // position at start of play

	// state fields
	Visible    bool          // is pacman visible?
	Pos        geom.Position // current screen position
	Dir        geom.Delta    // current heading
	Pcm        data.PCM      // current speed
	TunnelPcm  data.PCM      // speed when tunneling
	StallTimer int           // how many frames to stall for
	DyingFrame int           // dying animation frame to show if non-zero
}

// NewActor returns an Actor representing pacman at its start position.
func NewActor() *Actor {
	return &Actor{
		StartPos: geom.PACMAN_START,
	}
}

// Start gets the actor ready for start of play.
func (p *Actor) Start(pcm data.PCM) {
	p.Visible = true
	p.Pos = p.StartPos
	p.Dir = geom.LEFT
	p.Pcm = pcm
	p.TunnelPcm = 0
	p.StallTimer = 0
	p.DyingFrame = 0

}

// Steer adjusts pacman's heading based on the joystick input
// and the constraints of the maze.
func (p *Actor) Steer(v *video.Video, inDir int) {
	dir, ok := input.JoyDirection[inDir]
	if !ok {
		return
	}

	// direction can be taken if pacman is "lined up"
	if (dir.IsVertical() && (p.Pos.X&7) == 0) || (dir.IsHorizontal() && (p.Pos.Y&7) == 0) {
		nextPos := p.Pos.Add(dir.ScaleUp(8)).WrapTunnel()
		nextTile := v.GetTile(nextPos.TileXY())
		if nextTile.IsTraversable() {
			p.Dir = dir
		}
	}
}

// Pulse advances pacman's pulse train, and returns true if
// a movement update is due. If pacman is currently stalled,
// false is returned.
func (p *Actor) Pulse() bool {
	if p.Pcm.Pulse() {
		// TODO not clear if he should stall for a specified number of frames, updates, or pulses
		// let's go with pulses for now
		if p.StallTimer <= 0 {
			return true
		}
		p.StallTimer -= 1
	}
	return false
}

// Move moves pacman to its next screen position based on
// the current heading.
func (p *Actor) Move(v *video.Video) {
	viable := true

	if (p.Pos.X&7) == 0 && (p.Pos.Y&7) == 0 {
		nextPos := p.Pos.Add(p.Dir.ScaleUp(8)).WrapTunnel()
		nextTile := v.GetTile(nextPos.TileXY())
		viable = nextTile.IsTraversable()
	}

	if viable {
		p.Pos = p.Pos.Add(p.Dir)
		if p.Pos.X <= 4 && p.Dir.IsLeft() {
			p.Pos.X += 215
		} else if p.Pos.X >= 220 && p.Dir.IsRight() {
			p.Pos.X -= 215
		}
	}
}

// An anim describes a cycle of sprites for animating pacman
// opening and closing its mouth.
type anim [4]sprite.Look

// anims defines the animations for each of pacman's possible headings.
var anims = struct {
	Up    anim
	Left  anim
	Down  anim
	Right anim
}{
	anim{sprite.PACMAN_SHUT, sprite.PACMAN_UP2, sprite.PACMAN_UP1, sprite.PACMAN_UP2},
	anim{sprite.PACMAN_SHUT, sprite.PACMAN_LEFT2, sprite.PACMAN_LEFT1, sprite.PACMAN_LEFT2},
	anim{sprite.PACMAN_SHUT, sprite.PACMAN_DOWN2, sprite.PACMAN_DOWN1, sprite.PACMAN_DOWN2},
	anim{sprite.PACMAN_SHUT, sprite.PACMAN_RIGHT2, sprite.PACMAN_RIGHT1, sprite.PACMAN_RIGHT2},
}

// Draw schedules a sprite to render pacman in the next frame.
//
// The playerNumber allows for the look of each player's pacman to differ.
func (p *Actor) Draw(v *video.Video, playerNumber int) {
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
				look = anims.Up[frame]
			case p.Dir.IsLeft():
				look = anims.Left[frame]
			case p.Dir.IsDown():
				look = anims.Down[frame]
			case p.Dir.IsRight():
				look = anims.Right[frame]
			}
		}
		v.AddSprite(p.Pos, look, pal)
	}
}
