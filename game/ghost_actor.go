package game

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

var (
	HOME_CENTRE         = geom.Position{108, 136}
	PACMAN_START        = geom.Position{108, 208}
	BLINKY_START        = geom.Position{108, GHOST_HOME_EXITED_Y}
	GHOST_HOME_CENTRE   = HOME_CENTRE
	GHOST_HOME_TOP      = HOME_CENTRE.Y - 4
	GHOST_HOME_BOTTOM   = HOME_CENTRE.Y + 4
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

type GhostActor struct {
	// configuration fields, these don't change once set
	Id          int
	Pal         color.Palette
	StartPos    geom.Position
	HomePos     geom.Position
	ScatterPos  geom.Position
	AllDotLimit int

	// state fields
	Visible           bool
	Pos               geom.Position
	Dir               geom.Delta
	Pcm               data.PCM
	TunnelPcm         data.PCM
	Mode              Mode
	SubMode           SubMode
	TargetPos         geom.Position
	DotsAtHomeCounter int
	DotLimit          int
	ReversePending    bool
	ScoreLook         sprite.Look
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
		Pal:               color.PAL_BLINKY,
		HomePos:           GHOST_HOME_CENTRE,
		StartPos:          BLINKY_START,
		ScatterPos:        geom.TilePos(25, 0),
		AllDotLimit:       0,
		DotsAtHomeCounter: 0,
	}
}

func MakePinky() GhostActor {
	return GhostActor{
		Id:                PINKY,
		Pal:               color.PAL_PINKY,
		HomePos:           GHOST_HOME_CENTRE,
		StartPos:          GHOST_HOME_CENTRE,
		ScatterPos:        geom.TilePos(2, 2),
		AllDotLimit:       7,
		DotsAtHomeCounter: 0,
	}
}

func MakeInky() GhostActor {
	homePos := GHOST_HOME_CENTRE.Add(geom.Delta{-16, 0})
	return GhostActor{
		Id:                INKY,
		Pal:               color.PAL_INKY,
		HomePos:           homePos,
		StartPos:          homePos,
		ScatterPos:        geom.TilePos(25, 36),
		AllDotLimit:       17,
		DotsAtHomeCounter: 0,
	}
}

func MakeClyde() GhostActor {
	homePos := GHOST_HOME_CENTRE.Add(geom.Delta{16, 0})
	return GhostActor{
		Id:                CLYDE,
		Pal:               color.PAL_CLYDE,
		HomePos:           homePos,
		StartPos:          homePos,
		ScatterPos:        geom.TilePos(0, 36),
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
		g.Dir = geom.LEFT

	case PINKY:
		if maxGhosts <= 1 {
			g.Mode = MODE_HOME
		} else {
			g.Mode = MODE_LEAVING
		}
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Pinky
		g.Pcm = data.PCM_50
		g.Dir = geom.DOWN

	case INKY:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Inky
		g.Pcm = data.PCM_50
		g.Dir = geom.UP

	case CLYDE:
		g.Mode = MODE_HOME
		g.SubMode = SUBMODE_SCATTER
		g.DotLimit = dotLimits.Clyde
		g.Pcm = data.PCM_50
		g.Dir = geom.UP
	}

	g.ReversePending = false
	g.ScoreLook = 0
	g.DotsAtHomeCounter = 0

	g.Visible = true
	g.Pos = g.StartPos
	g.TunnelPcm = 0
}

func (g *GhostActor) SetLeaveState() {
	g.Mode = MODE_LEAVING
}

func (g *GhostActor) SetSubMode(subMode SubMode) {
	// Ghosts are forced to reverse direction by the system anytime the mode
	// changes from: chase-to-scatter, chase-to-frightened, scatter-to-chase,
	// and scatter-to-frightened.
	// Ghosts do not reverse direction when changing back from frightened to
	// chase or scatter modes.
	switch g.SubMode {
	case subMode:
		return

	case SUBMODE_CHASE:
		if subMode == SUBMODE_SCARED || subMode == SUBMODE_SCATTER {
			g.ReversePending = true
		}

	case SUBMODE_SCATTER:
		if subMode == SUBMODE_SCARED || subMode == SUBMODE_CHASE {
			g.ReversePending = true
		}
	}
	g.SubMode = subMode
}

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

func (g *GhostActor) Move() {
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

func (g *GhostActor) DrawGhost(v *video.Video, isWhite bool, wobble bool) {
	var look sprite.Look
	var pal color.Palette
	if g.Visible {
		switch {
		case g.Dir.IsUp():
			look = sprite.GHOST_UP1
		case g.Dir.IsLeft():
			look = sprite.GHOST_LEFT1
		case g.Dir.IsDown():
			look = sprite.GHOST_DOWN1
		case g.Dir.IsRight():
			look = sprite.GHOST_RIGHT1
		}
		pal = g.Pal
		if g.ScoreLook > 0 {
			look = g.ScoreLook
		} else {
			switch {
			case g.Mode == MODE_RETURNING:
				pal = color.PAL_EYES
			case g.SubMode == SUBMODE_SCARED:
				look = sprite.GHOST_SCARED1
				pal = color.PAL_SCARED
				if isWhite {
					pal = color.PAL_SCARED_FLASH
				}
			}
			if wobble {
				look += 1
			}
		}
		v.AddSprite(g.Pos, look, pal)
	}
}
