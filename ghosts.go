package main

import (
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/palette"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

// TODO - move to a common file
const (
	HOME_CENTRE_X = 108
	HOME_CENTRE_Y = 136

	PACMAN_START_X = HOME_CENTRE_X
	PACMAN_START_Y = 208
)

const (
	GHOST_HOME_CENTRE_X = HOME_CENTRE_X
	GHOST_HOME_CENTRE_Y = HOME_CENTRE_Y
	GHOST_HOME_TOP      = HOME_CENTRE_Y - 4
	GHOST_HOME_BOTTOM   = HOME_CENTRE_Y + 4
	GHOST_HOME_EXITED_Y = 112
)

// Control ghost behaviour
const (
	MODE_HOME      = 0
	MODE_LEAVING   = 1
	MODE_PLAYING   = 2
	MODE_RETURNING = 3
)

const (
	SUBMODE_SCATTER = 0
	SUBMODE_CHASE   = 1
	SUBMODE_SCARED  = 2
)

// TODO - move to a separate file
type Motion struct {
	X, Y           int
	Pcm, TunnelPcm uint32
	Vx, Vy         int
	Visible        bool
}

type GhostActor struct {
	Id                 int
	Pal                byte
	HomeX, HomeY       int
	StartX, StartY     int
	GlobalDotLimit     int
	ScatterX, ScatterY int

	Motion           Motion
	Mode             int
	SubMode          int
	ScoreSprite      byte
	TargetX, TargetY int
	DotCounter       int
	DotLimit         int
	ReversePending   bool
}

// Ghost identities
const (
	BLINKY = 0
	PINKY  = 1
	INKY   = 2
	CLYDE  = 3
)

func MakeBlinky() GhostActor {
	return GhostActor{
		Id:             BLINKY,
		Pal:            palette.BLINKY,
		HomeX:          GHOST_HOME_CENTRE_X,
		HomeY:          GHOST_HOME_CENTRE_Y,
		StartX:         GHOST_HOME_CENTRE_X,
		StartY:         GHOST_HOME_EXITED_Y,
		ScatterX:       25,
		ScatterY:       0,
		GlobalDotLimit: 0,
		DotCounter:     0,
	}
}

func MakePinky() GhostActor {
	return GhostActor{
		Id:             PINKY,
		Pal:            palette.PINKY,
		HomeX:          GHOST_HOME_CENTRE_X,
		HomeY:          GHOST_HOME_CENTRE_Y,
		StartX:         GHOST_HOME_CENTRE_X,
		StartY:         GHOST_HOME_CENTRE_Y,
		ScatterX:       2,
		ScatterY:       2,
		GlobalDotLimit: 7,
		DotCounter:     0,
	}
}

func MakeInky() GhostActor {
	return GhostActor{
		Id:             INKY,
		Pal:            palette.INKY,
		HomeX:          GHOST_HOME_CENTRE_X - 16,
		HomeY:          GHOST_HOME_CENTRE_Y,
		StartX:         GHOST_HOME_CENTRE_X - 16,
		StartY:         GHOST_HOME_CENTRE_Y,
		ScatterX:       25,
		ScatterY:       36,
		GlobalDotLimit: 17,
		DotCounter:     0,
	}
}

func MakeClyde() GhostActor {
	return GhostActor{
		Id:             CLYDE,
		Pal:            palette.CLYDE,
		HomeX:          GHOST_HOME_CENTRE_X + 16,
		HomeY:          GHOST_HOME_CENTRE_Y,
		StartX:         GHOST_HOME_CENTRE_X + 16,
		StartY:         GHOST_HOME_CENTRE_Y,
		ScatterX:       0,
		ScatterY:       36,
		GlobalDotLimit: 32,
		DotCounter:     0,
	}
}

func MakeGhosts() [4]GhostActor {
	return [4]GhostActor{
		MakeBlinky(),
		MakePinky(),
		MakeInky(),
		MakeClyde(),
	}
}

// TODO - move this somewhere better
func (g *Game) GhostsStart() {
	for i := range 4 {
		g.Ghosts[i].Start(
			g.LevelConfig.Speeds.Ghost,
			g.Options.MaxGhosts,
			&g.LevelConfig.DotLimits,
		)
	}
}

func (g *GhostActor) Start(pcmBlinky uint32, maxGhosts int, dotLimits *data.DotLimitEntry) {
	switch g.Id {
	case BLINKY:
		g.Mode = MODE_PLAYING
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = 0
		g.Motion.Pcm = pcmBlinky
		g.Motion.Vx = -1
		g.Motion.Vy = 0

	case PINKY:
		if maxGhosts <= 1 {
			g.Mode = MODE_HOME
		} else {
			g.Mode = MODE_LEAVING
		}
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Pinky
		g.Motion.Pcm = data.PCM_50
		g.Motion.Vx = 0
		g.Motion.Vy = 1

	case INKY:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Inky
		g.Motion.Pcm = data.PCM_50
		g.Motion.Vx = 0
		g.Motion.Vy = -1

	case CLYDE:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Clyde
		g.Motion.Pcm = data.PCM_50
		g.Motion.Vx = 0
		g.Motion.Vy = -1
	}

	g.ReversePending = false
	g.ScoreSprite = 0
	g.DotCounter = 0

	m := &g.Motion
	m.X = g.StartX
	m.Y = g.StartY
	m.TunnelPcm = 0
	m.Visible = true
}

func (g *GhostActor) DrawGhost(v *video.Video, isWhite bool, wobble bool) {
	var look byte
	var pal byte
	m := &g.Motion
	if m.Visible {
		switch {
		case m.Vx > 0:
			look = sprite.GHOST_RIGHT1
		case m.Vx < 0:
			look = sprite.GHOST_LEFT1
		case m.Vy > 0:
			look = sprite.GHOST_DOWN1
		case m.Vy < 0:
			look = sprite.GHOST_UP1
		}
		pal = byte(g.Pal)
		if g.ScoreSprite > 0 {
			look = g.ScoreSprite
		} else {
			switch {
			case g.Mode == MODE_RETURNING:
				pal = palette.EYES
			case g.SubMode == SUBMODE_SCARED:
				look = sprite.GHOST_SCARED1
				pal = palette.SCARED
				if isWhite {
					pal = palette.SCARED_FLASH
				}
			}
			if wobble {
				look += 1
			}
		}
		v.AddSprite(m.X-4, m.Y-4-MAZE_TOP, look, pal)
	}
}
