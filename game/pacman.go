package game

import (
	"github.com/adrmcintyre/poweraid/ghost"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/tile"
)

func (g *Game) PacmanStart() {
	g.Pacman.Start(g.LevelConfig.Speeds.Pacman)
}

func (g *Game) IsPacmanIdle() bool {
	return g.LevelState.FrameCounter >= g.LevelState.IdleAfter
}

func (g *Game) PacmanResetIdleTimer() {
	g.LevelState.IdleAfter = g.LevelState.FrameCounter + g.LevelConfig.IdleLimit
}

func (g *Game) PacmanRevert(revert bool) {
	if revert {
		g.Pacman.Pcm = g.LevelConfig.Speeds.Pacman
	}
}

func (g *Game) PacmanPulse() bool {
	return g.Pacman.Pulse()
}

func (g *Game) PacmanSteer(pulsed bool) {
	//	if pulsed {
	inDir := input.GetJoystickDirection()
	g.Pacman.Steer(&g.Video, inDir)
	// }
}

func (g *Game) PacmanMove(pulsed bool) {
	if pulsed {
		g.Pacman.Move(&g.Video)
	}
}

// returns true if pacman dies
func (g *Game) PacmanCollide() bool {
	v := &g.Video

	pacPos := g.Pacman.Pos
	x, y := pacPos.TileXY()

	switch v.GetTile(x, y) {
	case tile.PILL:
		v.SetTile(x, y, tile.SPACE)
		g.EatPill()
	case tile.POWER, tile.POWER_SMALL:
		v.SetTile(x, y, tile.SPACE)
		g.EatPower()
	}

	if g.LevelState.BonusTimeout > 0 &&
		pacPos.TileEq(g.BonusActor.Pos) {
		g.EatBonus()
	}

	for _, gh := range g.Ghosts {
		if (gh.Mode == ghost.MODE_PLAYING) &&
			(gh.SubMode == ghost.SUBMODE_SCARED) &&
			pacPos.TileEq(gh.Pos) {
			g.GhostConsume(gh)
		}
	}

	for _, gh := range g.Ghosts {
		if (gh.Mode == ghost.MODE_PLAYING) &&
			(gh.SubMode != ghost.SUBMODE_SCARED) &&
			pacPos.TileEq(gh.Pos) {
			return true
		}
	}

	return false
}
