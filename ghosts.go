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

// maybe an enum of up,down,left,right instead, plus a converter?
type Velocity struct {
	Vx int
	Vy int
}

type GhostActor struct {
	// configuration fields, these don't change once set
	Id          int
	Pal         byte
	HomePos     video.ScreenPos
	StartPos    video.ScreenPos
	AllDotLimit int
	ScatterPos  video.TilePos

	// state fields
	Visible           bool
	Mode              Mode
	SubMode           SubMode
	ScoreSprite       byte
	TargetPos         video.TilePos
	DotsAtHomeCounter int
	DotLimit          int
	ReversePending    bool
	Pos               video.ScreenPos
	Vel               Velocity
	Pcm               data.PCM
	TunnelPcm         data.PCM
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
		Id:                BLINKY,
		Pal:               palette.BLINKY,
		HomePos:           video.ScreenPos{GHOST_HOME_CENTRE_X, GHOST_HOME_CENTRE_Y},
		StartPos:          video.ScreenPos{GHOST_HOME_CENTRE_X, GHOST_HOME_EXITED_Y},
		ScatterPos:        video.TilePos{25, 0},
		AllDotLimit:       0,
		DotsAtHomeCounter: 0,
	}
}

func MakePinky() GhostActor {
	return GhostActor{
		Id:                PINKY,
		Pal:               palette.PINKY,
		HomePos:           video.ScreenPos{GHOST_HOME_CENTRE_X, GHOST_HOME_CENTRE_Y},
		StartPos:          video.ScreenPos{GHOST_HOME_CENTRE_X, GHOST_HOME_CENTRE_Y},
		ScatterPos:        video.TilePos{2, 2},
		AllDotLimit:       7,
		DotsAtHomeCounter: 0,
	}
}

func MakeInky() GhostActor {
	return GhostActor{
		Id:                INKY,
		Pal:               palette.INKY,
		HomePos:           video.ScreenPos{GHOST_HOME_CENTRE_X - 16, GHOST_HOME_CENTRE_Y},
		StartPos:          video.ScreenPos{GHOST_HOME_CENTRE_X - 16, GHOST_HOME_CENTRE_Y},
		ScatterPos:        video.TilePos{25, 36},
		AllDotLimit:       17,
		DotsAtHomeCounter: 0,
	}
}

func MakeClyde() GhostActor {
	return GhostActor{
		Id:                CLYDE,
		Pal:               palette.CLYDE,
		HomePos:           video.ScreenPos{GHOST_HOME_CENTRE_X + 16, GHOST_HOME_CENTRE_Y},
		StartPos:          video.ScreenPos{GHOST_HOME_CENTRE_X + 16, GHOST_HOME_CENTRE_Y},
		ScatterPos:        video.TilePos{0, 36},
		AllDotLimit:       32,
		DotsAtHomeCounter: 0,
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

func (g *GhostActor) Start(pcmBlinky data.PCM, maxGhosts int, dotLimits *data.DotLimitEntry) {
	switch g.Id {
	case BLINKY:
		g.Mode = MODE_PLAYING
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = 0
		g.Pcm = pcmBlinky
		g.Vel = Velocity{-1, 0}

	case PINKY:
		if maxGhosts <= 1 {
			g.Mode = MODE_HOME
		} else {
			g.Mode = MODE_LEAVING
		}
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Pinky
		g.Pcm = data.PCM_50
		g.Vel = Velocity{0, 1}

	case INKY:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Inky
		g.Pcm = data.PCM_50
		g.Vel = Velocity{0, -1}

	case CLYDE:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Clyde
		g.Pcm = data.PCM_50
		g.Vel = Velocity{0, -1}
	}

	g.ReversePending = false
	g.ScoreSprite = 0
	g.DotsAtHomeCounter = 0

	g.Visible = true
	g.Pos = g.StartPos
	g.TunnelPcm = 0
}

func (g *GhostActor) DrawGhost(v *video.Video, isWhite bool, wobble bool) {
	var look byte
	var pal byte
	if g.Visible {
		switch {
		case g.Vel.Vx > 0:
			look = sprite.GHOST_RIGHT1
		case g.Vel.Vx < 0:
			look = sprite.GHOST_LEFT1
		case g.Vel.Vy > 0:
			look = sprite.GHOST_DOWN1
		case g.Vel.Vy < 0:
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
		v.AddSprite(g.Pos.X-4, g.Pos.Y-4-MAZE_TOP, look, pal)
	}
}
