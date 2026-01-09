package game

import (
	"github.com/adrmcintyre/ebiman/ghost"
	"github.com/adrmcintyre/ebiman/tile"
)

// PacmanStart gets pacman ready at the start of a level.
func (g *Game) PacmanStart() {
	g.Pacman.Start(g.LevelConfig.Speeds.Pacman)
}

// IsPacmanIdle returns true if pacman has failed to consume any
// dots before the idle timeout has expired.
func (g *Game) IsPacmanIdle() bool {
	return g.LevelState.FrameCounter >= g.LevelState.IdleAfter
}

// PacmanResetIdleTimer resets the expiry time of the idleness timer.
func (g *Game) PacmanResetIdleTimer() {
	g.LevelState.IdleAfter = g.LevelState.FrameCounter + g.LevelConfig.IdleLimit
}

// PacmanRevert returns pacman to his normal speed if
// the revert flag is set.
func (g *Game) PacmanRevert(revert bool) {
	if revert {
		g.Pacman.Pcm = g.LevelConfig.Speeds.Pacman
	}
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
		v.SetTile(x, y, tile.Space)
	case t.IsPower():
		g.EatPower()
		v.SetTile(x, y, tile.Space)
	}

	if g.LevelState.BonusTimeout > 0 &&
		pacPos.TileEq(g.BonusActor.Pos) {
		g.EatBonus()
	}

	for _, gh := range g.Ghosts {
		if (gh.Mode == ghost.ModePlaying) &&
			(gh.SubMode == ghost.SubModeScared) &&
			pacPos.TileEq(gh.Pos) {
			g.PacmanEatsGhost(gh)
		}
	}

	for _, gh := range g.Ghosts {
		if (gh.Mode == ghost.ModePlaying) &&
			(gh.SubMode != ghost.SubModeScared) &&
			pacPos.TileEq(gh.Pos) {
			return true
		}
	}

	return false
}
