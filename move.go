package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/video"
)

func (g *GhostActor) GhostTunnel(pcm data.PCM) {
	m := &g.Motion
	tilePos := m.Pos.ToTilePos()
	// TODO - constants
	if tilePos.Y == 17 && (tilePos.X <= 5 || tilePos.X >= 22) {
		if m.TunnelPcm == 0 {
			m.TunnelPcm = pcm
		}
	} else {
		m.TunnelPcm = 0
	}
}

func (g *GhostActor) MoveGhost() {
	m := &g.Motion

	nextPos := video.ScreenPos{
		m.Pos.X + m.Vel.Vx,
		m.Pos.Y + m.Vel.Vy,
	}

	// account for tunnel:
	if nextPos.X <= 4 && m.Vel.Vx < 0 {
		nextPos.X += 215
	} else if nextPos.X >= 220 && m.Vel.Vx > 0 {
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

	m.Pos = nextPos
}

func (g *Game) GhostPulse(i int) bool {
	ghost := &g.Ghosts[i]
	m := &ghost.Motion

	pcm := &m.Pcm

	isBlinky := ghost.Id == BLINKY
	isHunting := ghost.Mode == MODE_PLAYING && ghost.SubMode != SUBMODE_SCARED
	isClydeOut := g.Ghosts[CLYDE].Mode != MODE_HOME

	if m.TunnelPcm != 0 {
		pcm = &m.TunnelPcm
	} else if isBlinky && isHunting && isClydeOut {
		if g.LevelState.DotsRemaining <= g.LevelConfig.ElroyPills2 {
			pcm = &g.LevelConfig.Speeds.Elroy2
		} else if g.LevelState.DotsRemaining <= g.LevelConfig.ElroyPills1 {
			pcm = &g.LevelConfig.Speeds.Elroy1
		}
	}

	return pcm.Pulse()
}
