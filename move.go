package main

func (g *GhostActor) GhostTunnel(pcm uint32) {
	m := &g.Motion
	tileX := m.X / 8
	tileY := m.Y / 8
	// TODO - constants
	if tileY == 17 && (tileX <= 5 || tileX >= 22) {
		if m.TunnelPcm == 0 {
			m.TunnelPcm = pcm
		}
	} else {
		m.TunnelPcm = 0
	}
}

func (g *GhostActor) MoveGhost() {
	m := &g.Motion

	nextX := m.X + m.Vx
	nextY := m.Y + m.Vy

	// account for tunnel:
	if nextX <= 4 && m.Vx < 0 {
		nextX += 215
	} else if nextX >= 220 && m.Vx > 0 {
		nextX -= 215
	}

	// NOTE
	// sanity to prevent ghosts falling off the top or bottom of the maze
	// this shouldn't be necessary if navigation is operating correctly
	if nextY < 12 {
		nextY = 12
	} else if nextY > 260 {
		nextY = 260
	}

	m.X = nextX
	m.Y = nextY
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

	msb := (*pcm) >> 31
	*pcm = (*pcm << 1) | msb
	return msb != 0
}
