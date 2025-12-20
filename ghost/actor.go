package ghost

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/pacman"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
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

type Actor struct {
	// configuration fields, these don't change once set
	Id          Id
	Pal         color.Palette
	StartPos    geom.Position
	HomePos     geom.Position
	ScatterPos  geom.Position
	AllDotLimit int
	Pacman      *pacman.Actor
	Blinky      *Actor

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
type Id int

const (
	BLINKY Id = iota
	PINKY
	INKY
	CLYDE
)

func NewBlinky(pacman *pacman.Actor) *Actor {
	return &Actor{
		Id:                BLINKY,
		Pal:               color.PAL_BLINKY,
		HomePos:           geom.BLINKY_HOME,
		StartPos:          geom.BLINKY_START,
		ScatterPos:        geom.BLINKY_SCATTER,
		AllDotLimit:       0,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
	}
}

func NewPinky(pacman *pacman.Actor) *Actor {
	return &Actor{
		Id:                PINKY,
		Pal:               color.PAL_PINKY,
		HomePos:           geom.PINKY_HOME,
		StartPos:          geom.PINKY_HOME,
		ScatterPos:        geom.PINKY_SCATTER,
		AllDotLimit:       7,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
	}
}

func NewInky(pacman *pacman.Actor, blinky *Actor) *Actor {
	return &Actor{
		Id:                INKY,
		Pal:               color.PAL_INKY,
		HomePos:           geom.INKY_HOME,
		StartPos:          geom.INKY_HOME,
		ScatterPos:        geom.INKY_SCATTER,
		AllDotLimit:       17,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
		Blinky:            blinky,
	}
}

func NewClyde(pacman *pacman.Actor) *Actor {
	return &Actor{
		Id:                CLYDE,
		Pal:               color.PAL_CLYDE,
		HomePos:           geom.CLYDE_HOME,
		StartPos:          geom.CLYDE_HOME,
		ScatterPos:        geom.CLYDE_SCATTER,
		AllDotLimit:       32,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
	}
}

func (g *Actor) Start(pcmBlinky data.PCM, maxGhosts int, dotLimits *data.DotLimitEntry) {
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

func (g *Actor) SetLeaveState() {
	g.Mode = MODE_LEAVING
}

func (g *Actor) SetSubMode(subMode SubMode) {
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

func (g *Actor) Tunnel(pcm data.PCM) {
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

func (g *Actor) Move() {
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

func (g *Actor) Draw(v *video.Video, isWhite bool, wobble bool) {
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
