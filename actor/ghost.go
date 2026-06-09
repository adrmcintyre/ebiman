package actor

import (
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/video"
)

// A GhostMode identifies a ghost's relationship with its home.
type GhostMode int

const (
	GhostModeHome      GhostMode = iota // confined at home
	GhostModeLeaving                    // leaving home (being released)
	GhostModePlaying                    // in active play
	GhostModeReturning                  // returning home
)

// A GhostSubMode identifies a ghost's behaviour.
type GhostSubMode int

const (
	GhostSubModeScattering GhostSubMode = iota // seeking preferred area of the maze
	GhostSubModeChasing                        // actively hunting pacman
	GhostSubModeScared                         // fleeing pacman
)

// An Ghost describes the look and behaviour of a ghost.
type Ghost struct {
	// configuration fields, these don't change once set
	Id         Id            // the ghost's identity
	Pal        video.Palette // its colouring
	StartPos   geom.Position // where it starts
	HomePos    geom.Position // its preferred spot at home
	ScatterPos geom.Position // its preferred spot in scatter mode
	// AllDotLimit specifies how many dots pacman can eat
	// after dying for the first time before the ghost no
	// longer stays at home.
	AllDotLimit int
	Pacman      *Pacman // reference to pacman for targetting
	Blinky      *Ghost  // reference to blinky for coordination

	// state fields
	Visible           bool          // is it visible?
	Pos               geom.Position // its screen position
	Dir               geom.Delta    // its current heading
	Pcm               data.PCM      // its current speed
	TunnelPcm         data.PCM      // its speed when tunneling
	Mode              GhostMode     // its current mode
	SubMode           GhostSubMode  // its current submode
	TargetPos         geom.Position // target location to seek
	DotsAtHomeCounter int           // dots eaten while ghost is home
	DotLimit          int           // how many dots before release
	ReversePending    bool          // change direction entering next tile?
	ScoreLook         video.Sprite  // render as this score sprite if non-zero
}

// An Id identifies a specific Ghost.
type Id int

// The identities of each ghost.
const (
	Blinky Id = iota
	Pinky
	Inky
	Clyde
)

// NewBlinky returns a new Actor configured to represent blinky.
// It takes a reference to pacman for target seeking.
func NewBlinky(pacman *Pacman) *Ghost {
	return &Ghost{
		Id:                Blinky,
		Pal:               video.PalBlinky,
		HomePos:           geom.BlinkyHome,
		StartPos:          geom.BlinkyStart,
		ScatterPos:        geom.BlinkyScatter,
		AllDotLimit:       0,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
	}
}

// NewPinky returns a new Actor configure to represent pinky.
// It takes a reference to pacman for target seeking.
func NewPinky(pacman *Pacman) *Ghost {
	return &Ghost{
		Id:                Pinky,
		Pal:               video.PalPinky,
		HomePos:           geom.PinkyHome,
		StartPos:          geom.PinkyHome,
		ScatterPos:        geom.PinkyScatter,
		AllDotLimit:       7,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
	}
}

// NewInky returns a new Actor configured to represent inky.
// It takes a reference to pacman for target seeking, and blinky
// with whom it co-ordinates behavior.
func NewInky(pacman *Pacman, blinky *Ghost) *Ghost {
	return &Ghost{
		Id:                Inky,
		Pal:               video.PalInky,
		HomePos:           geom.InkyHome,
		StartPos:          geom.InkyHome,
		ScatterPos:        geom.InkyScatter,
		AllDotLimit:       17,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
		Blinky:            blinky,
	}
}

// NewClyde returns a new Actor configured to represent clyde.
// It takes a reference to pacman for target seeking.
func NewClyde(pacman *Pacman) *Ghost {
	return &Ghost{
		Id:                Clyde,
		Pal:               video.PalClyde,
		HomePos:           geom.ClydeHome,
		StartPos:          geom.ClydeHome,
		ScatterPos:        geom.ClydeScatter,
		AllDotLimit:       32,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
	}
}

// Start puts the actor into its initial state ready for playing.
func (g *Ghost) Start(pcmBlinky data.PCM, maxGhosts int, dotLimits *data.DotLimitEntry) {
	switch g.Id {
	case Blinky:
		g.Mode = GhostModePlaying
		g.SubMode = GhostSubModeScattering
		g.DotLimit = 0
		g.Pcm = pcmBlinky
		g.Dir = geom.Left

	case Pinky:
		if maxGhosts <= 1 {
			g.Mode = GhostModeHome
		} else {
			g.Mode = GhostModeLeaving
		}
		g.SubMode = GhostSubModeScattering
		g.DotLimit = dotLimits.Pinky
		g.Pcm = data.PCM50
		g.Dir = geom.Down

	case Inky:
		g.Mode = GhostModeHome
		g.SubMode = GhostSubModeScattering
		g.DotLimit = dotLimits.Inky
		g.Pcm = data.PCM50
		g.Dir = geom.Up

	case Clyde:
		g.Mode = GhostModeHome
		g.SubMode = GhostSubModeScattering
		g.DotLimit = dotLimits.Clyde
		g.Pcm = data.PCM50
		g.Dir = geom.Up
	}

	g.ReversePending = false
	g.ScoreLook = 0
	g.DotsAtHomeCounter = 0

	g.Visible = true
	g.Pos = g.StartPos
	g.TunnelPcm = 0
}

// SetLeaveState tells the ghost to leave its home.
func (g *Ghost) SetLeaveState() {
	g.Mode = GhostModeLeaving
}

// SetSubMode changes the ghost's submode.
func (g *Ghost) SetSubMode(subMode GhostSubMode) {
	// Ghosts are forced to reverse direction by the system anytime the mode
	// changes from: chase-to-scatter, chase-to-frightened, scatter-to-chase,
	// and scatter-to-frightened.
	// Ghosts do not reverse direction when changing back from frightened to
	// chase or scatter modes.
	switch g.SubMode {
	case subMode:
		return

	case GhostSubModeChasing:
		if subMode == GhostSubModeScared || subMode == GhostSubModeScattering {
			g.ReversePending = true
		}

	case GhostSubModeScattering:
		if subMode == GhostSubModeScared || subMode == GhostSubModeChasing {
			g.ReversePending = true
		}
	}
	g.SubMode = subMode
}

// CheckTunnelSpeed ensures the ghost moves at the correct
// speed if in the tunnel.
func (g *Ghost) CheckTunnelSpeed(pcm data.PCM) {
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
// based on its current heading. Returns true if
// ghost has moved to a different tile.
func (g *Ghost) Move() bool {
	x0, y0 := g.Pos.TileXY()
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
	x1, y1 := nextPos.TileXY()
	return !(x0 == x1 && y0 == y1)
}

// Draw schedules the ghost's sprite for display at the next frame.
// If isWhite is set, the "scared" palette is applied.
// One of two looks are selected by wobble, giving the ghost its
// distinctive animation.
func (g *Ghost) Draw(v *video.Video, isWhite bool, wobble bool) {
	var look video.Sprite
	var pal video.Palette
	if g.Visible {
		switch {
		case g.Dir.IsUp():
			look = video.SpriteGhostUp1
		case g.Dir.IsLeft():
			look = video.SpriteGhostDown2
		case g.Dir.IsDown():
			look = video.SpriteGhostDown1
		case g.Dir.IsRight():
			look = video.SpriteGhostRight1
		}
		pal = g.Pal
		if g.ScoreLook > 0 {
			look = g.ScoreLook
		} else {
			switch {
			case g.Mode == GhostModeReturning:
				pal = video.PalEyes
			case g.SubMode == GhostSubModeScared:
				look = video.SpriteGhostScared1
				pal = video.PalScared
				if isWhite {
					pal = video.PalScaredFlash
				}
			}
			if wobble {
				look += 1
			}
		}
		v.AddSprite(g.Pos, look, pal)
	}
}
