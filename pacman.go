package main

import (
	"github.com/adrmcintyre/ebiman/video"
)

// PacmanStart gets pacman ready at the start of a level.
func (g *Game) PacmanStart() {
	g.Pacman.Start(g.LevelConfig.Speeds.Pacman)
}

// IsPacmanIdle returns true if pacman has failed to consume any
// dots before the idle timeout has expired.
func (g *Game) IsPacmanIdle() bool {
	return g.Level.FrameCounter >= g.Level.IdleAfter
}

// PacmanResetIdleTimer resets the expiry time of the idleness timer.
func (g *Game) PacmanResetIdleTimer() {
	g.Level.IdleAfter = g.Level.FrameCounter + g.LevelConfig.IdleLimit
}

// PacmanRevertSpeed returns pacman to his normal speed.
func (g *Game) PacmanRevertSpeed() {
	g.Pacman.Pcm = g.LevelConfig.Speeds.Pacman
}

// PacmanPulse advances pacman's pulse train, and returns true if
// he pulsed, and so is due for a movement update.
func (g *Game) PacmanPulse() bool {
	return g.Pacman.Pulse()
}

// PacmanSteer alters pacman's current heading in accordance with
// joystick input and the constraints of the maze.
func (g *Game) PacmanSteer(pulsed bool) {
	inDir := g.Input.JoystickDirection()
	g.Pacman.Steer(g.Video, inDir)
}

// PacmanMove advances pacman is he has just pulsed.
func (g *Game) PacmanMove(pulsed bool) {
	if pulsed {
		g.Pacman.Move(g.Video)
	}
}

// PacmanCollide checks for collisions with pills, power pills, the
// bonus and vulnerable ghosts and takes the appropriate action.
// If pacman collides with an invulnerable ghost, true is returned.
func (g *Game) PacmanCollide() bool {
	v := g.Video

	pacPos := g.Pacman.Pos
	x, y := pacPos.TileXY()

	t := v.GetTile(x, y)
	switch {
	case t.IsPill():
		g.EatPill(t)
		v.SetTile(x, y, video.TileSpace)
	case t.IsPower():
		g.EatPower()
		v.SetTile(x, y, video.TileSpace)
	}

	if g.Level.BonusTimeout > 0 && pacPos.TileEq(g.BonusActor.Pos) {
		g.EatBonus()
	}

	collided := false
	for _, gh := range g.Ghosts {
		if pacPos.TileEq(gh.Pos) {
			switch {
			case gh.IsVulnerable():
				g.PacmanEatsGhost(gh)
			case gh.IsDangerous():
				collided = true
			}
		}
	}

	return collided
}
