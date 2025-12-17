package main

import (
	"github.com/adrmcintyre/poweraid/data"
)

func (g *GhostActor) Tunnel(pcm data.PCM) {
	x, y := g.Pos.TileXY()
	// TODO - constants
	if y == 17 && (x <= 5 || x >= 22) {
		if g.TunnelPcm == 0 {
			g.TunnelPcm = pcm
		}
	} else {
		g.TunnelPcm = 0
	}
}

func (g *GhostActor) MoveGhost() {
	nextPos := g.Pos.Add(g.Dir)

	// account for tunnel:
	if nextPos.X <= 4 && g.Dir.IsLeft() {
		nextPos.X += 215
	} else if nextPos.X >= 220 && g.Dir.IsRight() {
		nextPos.X -= 215
	}

	// NOTE
	// sanity to prevent ghosts falling off the top or bottom of the maze
	// this shouldn't be necessary if navigation is operating correctly
	if nextPos.Y < 12 {
		nextPos.Y = 12
	} else if nextPos.Y > 260 {
		nextPos.Y = 260
	}

	g.Pos = nextPos
}

func (g *Game) Pulse(i int) bool {
	ghost := &g.Ghosts[i]

	pcm := &ghost.Pcm

	isBlinky := ghost.Id == BLINKY
	isHunting := ghost.Mode == MODE_PLAYING && ghost.SubMode != SUBMODE_SCARED
	isClydeOut := g.Ghosts[CLYDE].Mode != MODE_HOME

	if ghost.TunnelPcm != 0 {
		pcm = &ghost.TunnelPcm
	} else if isBlinky && isHunting && isClydeOut {
		if g.LevelState.DotsRemaining <= g.LevelConfig.ElroyPills2 {
			pcm = &g.LevelConfig.Speeds.Elroy2
		} else if g.LevelState.DotsRemaining <= g.LevelConfig.ElroyPills1 {
			pcm = &g.LevelConfig.Speeds.Elroy1
		}
	}

	return pcm.Pulse()
}
