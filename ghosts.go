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
type Mode int

const (
	MODE_HOME Mode = iota
	MODE_LEAVING
	MODE_PLAYING
	MODE_RETURNING
)

type SubMode int

const (
	SUBMODE_SCATTER SubMode = iota
	SUBMODE_CHASE
	SUBMODE_SCARED
)

type Position struct {
	X int
	Y int
}

type Velocity struct {
	Vx int
	Vy int
}

// TODO - move to a separate file
type Motion struct {
	Pos       Position
	Vel       Velocity
	Pcm       uint32
	TunnelPcm uint32
	Visible   bool
}

type GhostActor struct {
	Id             int
	Pal            byte
	HomePos        Position
	StartPos       Position
	GlobalDotLimit int
	ScatterPos     Position

	Motion         Motion
	Mode           Mode
	SubMode        SubMode
	ScoreSprite    byte
	TargetPos      Position
	DotCounter     int
	DotLimit       int
	ReversePending bool
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
		HomePos:        Position{GHOST_HOME_CENTRE_X, GHOST_HOME_CENTRE_Y},
		StartPos:       Position{GHOST_HOME_CENTRE_X, GHOST_HOME_EXITED_Y},
		ScatterPos:     Position{25, 0},
		GlobalDotLimit: 0,
		DotCounter:     0,
	}
}

func MakePinky() GhostActor {
	return GhostActor{
		Id:             PINKY,
		Pal:            palette.PINKY,
		HomePos:        Position{GHOST_HOME_CENTRE_X, GHOST_HOME_CENTRE_Y},
		StartPos:       Position{GHOST_HOME_CENTRE_X, GHOST_HOME_CENTRE_Y},
		ScatterPos:     Position{2, 2},
		GlobalDotLimit: 7,
		DotCounter:     0,
	}
}

func MakeInky() GhostActor {
	return GhostActor{
		Id:             INKY,
		Pal:            palette.INKY,
		HomePos:        Position{GHOST_HOME_CENTRE_X - 16, GHOST_HOME_CENTRE_Y},
		StartPos:       Position{GHOST_HOME_CENTRE_X - 16, GHOST_HOME_CENTRE_Y},
		ScatterPos:     Position{25, 36},
		GlobalDotLimit: 17,
		DotCounter:     0,
	}
}

func MakeClyde() GhostActor {
	return GhostActor{
		Id:             CLYDE,
		Pal:            palette.CLYDE,
		HomePos:        Position{GHOST_HOME_CENTRE_X + 16, GHOST_HOME_CENTRE_Y},
		StartPos:       Position{GHOST_HOME_CENTRE_X + 16, GHOST_HOME_CENTRE_Y},
		ScatterPos:     Position{0, 36},
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
		g.Motion.Vel = Velocity{-1, 0}

	case PINKY:
		if maxGhosts <= 1 {
			g.Mode = MODE_HOME
		} else {
			g.Mode = MODE_LEAVING
		}
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Pinky
		g.Motion.Pcm = data.PCM_50
		g.Motion.Vel = Velocity{0, 1}

	case INKY:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Inky
		g.Motion.Pcm = data.PCM_50
		g.Motion.Vel = Velocity{0, -1}

	case CLYDE:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Clyde
		g.Motion.Pcm = data.PCM_50
		g.Motion.Vel = Velocity{0, -1}
	}

	g.ReversePending = false
	g.ScoreSprite = 0
	g.DotCounter = 0

	m := &g.Motion
	m.Pos = g.StartPos
	m.TunnelPcm = 0
	m.Visible = true
}

func (g *GhostActor) DrawGhost(v *video.Video, isWhite bool, wobble bool) {
	var look byte
	var pal byte
	m := &g.Motion
	if m.Visible {
		switch {
		case m.Vel.Vx > 0:
			look = sprite.GHOST_RIGHT1
		case m.Vel.Vx < 0:
			look = sprite.GHOST_LEFT1
		case m.Vel.Vy > 0:
			look = sprite.GHOST_DOWN1
		case m.Vel.Vy < 0:
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
		v.AddSprite(m.Pos.X-4, m.Pos.Y-4-MAZE_TOP, look, pal)
	}
}
