package ghost

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/pacman"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/video"
)

// A Mode identifies a ghost's relationship with its home.
type Mode int

const (
	MODE_HOME      Mode = iota // confined at home
	MODE_LEAVING               // leaving home (being released)
	MODE_PLAYING               // in active play
	MODE_RETURNING             // returning home
)

// A SubMode identifies a ghost's behaviour.
type SubMode int

const (
	SUBMODE_SCATTER SubMode = iota // seeking preferred area of the maze
	SUBMODE_CHASE                  // actively hunting pacman
	SUBMODE_SCARED                 // fleeing pacman
)

// An Actor describes the look and behaviour of a ghost.
type Actor struct {
	// configuration fields, these don't change once set
	Id         Id            // the ghost's identity
	Pal        color.Palette // its colouring
	StartPos   geom.Position // where it starts
	HomePos    geom.Position // its preferred spot at home
	ScatterPos geom.Position // its preferred spot in scatter mode
	// AllDotLimit specifies how many dots pacman can eat
	// after dying for the first time before the ghost no
	// longer stays at home.
	AllDotLimit int
	Pacman      *pacman.Actor // reference to pacman for targetting
	Blinky      *Actor        // reference to blinky for coordination

	// state fields
	Visible           bool          // is it visible?
	Pos               geom.Position // its screen position
	Dir               geom.Delta    // its current heading
	Pcm               data.PCM      // its current speed
	TunnelPcm         data.PCM      // its speed when tunneling
	Mode              Mode          // its current mode
	SubMode           SubMode       // its current submode
	TargetPos         geom.Position // target location to seek
	DotsAtHomeCounter int           // dots eaten while ghost is home
	DotLimit          int           // how many dots before release
	ReversePending    bool          // change direction entering next tile?
	ScoreLook         sprite.Look   // render as this score sprite if non-zero
}

// An Id identifies a specific Ghost.
type Id int

// The identities of each ghost.
const (
	BLINKY Id = iota
	PINKY
	INKY
	CLYDE
)

// NewBlinky returns a new Actor configured to represent blinky.
// It takes a reference to pacman for target seeking.
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

// NewPinky returns a new Actor configure to represent pinky.
// It takes a reference to pacman for target seeking.
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

// NewInky returns a new Actor configured to represent inky.
// It takes a reference to pacman for target seeking, and blinky
// with whom it co-ordinates behavior.
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

// NewClyde returns a new Actor configured to represent clyde.
// It takes a reference to pacman for target seeking.
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

// Start puts the actor into its initial state ready for playing.
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

// SetLeaveState tells the ghost to leave its home.
func (g *Actor) SetLeaveState() {
	g.Mode = MODE_LEAVING
}

// SetSubMode changes the ghost's submode.
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

// CheckTunnelSpeed ensures the ghost moves at the correct
// speed if in the tunnel.
func (g *Actor) CheckTunnelSpeed(pcm data.PCM) {
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

// Move moves the ghost to its next screen location
// based on its current heading.
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

// Draw schedules the ghost's sprite for display at the next frame.
// If isWhite is set, the "scared" palette is applied.
// One of two looks are selected by wobble, giving the ghost its
// distinctive animation.
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
