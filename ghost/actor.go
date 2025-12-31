package ghost

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/geom"
	"github.com/adrmcintyre/ebiman/pacman"
	"github.com/adrmcintyre/ebiman/sprite"
	"github.com/adrmcintyre/ebiman/video"
)

// A Mode identifies a ghost's relationship with its home.
type Mode int

const (
	ModeHome      Mode = iota // confined at home
	ModeLeaving               // leaving home (being released)
	ModePlaying               // in active play
	ModeReturning             // returning home
)

// A SubMode identifies a ghost's behaviour.
type SubMode int

const (
	SubModeScatter SubMode = iota // seeking preferred area of the maze
	SubModeChase                  // actively hunting pacman
	SubModeScared                 // fleeing pacman
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
	Blinky Id = iota
	Pinky
	Inky
	Clyde
)

// NewBlinky returns a new Actor configured to represent blinky.
// It takes a reference to pacman for target seeking.
func NewBlinky(pacman *pacman.Actor) *Actor {
	return &Actor{
		Id:                Blinky,
		Pal:               color.PalBlinky,
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
func NewPinky(pacman *pacman.Actor) *Actor {
	return &Actor{
		Id:                Pinky,
		Pal:               color.PalPinky,
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
func NewInky(pacman *pacman.Actor, blinky *Actor) *Actor {
	return &Actor{
		Id:                Inky,
		Pal:               color.PalInky,
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
func NewClyde(pacman *pacman.Actor) *Actor {
	return &Actor{
		Id:                Clyde,
		Pal:               color.PalClyde,
		HomePos:           geom.ClydeHome,
		StartPos:          geom.ClydeHome,
		ScatterPos:        geom.ClydeScatter,
		AllDotLimit:       32,
		DotsAtHomeCounter: 0,
		Pacman:            pacman,
	}
}

// Start puts the actor into its initial state ready for playing.
func (g *Actor) Start(pcmBlinky data.PCM, maxGhosts int, dotLimits *data.DotLimitEntry) {
	switch g.Id {
	case Blinky:
		g.Mode = ModePlaying
		g.SubMode = SubModeScatter
		g.DotLimit = 0
		g.Pcm = pcmBlinky
		g.Dir = geom.Left

	case Pinky:
		if maxGhosts <= 1 {
			g.Mode = ModeHome
		} else {
			g.Mode = ModeLeaving
		}
		g.SubMode = SubModeScatter
		g.DotLimit = dotLimits.Pinky
		g.Pcm = data.PCM50
		g.Dir = geom.Down

	case Inky:
		g.Mode = ModeHome
		g.SubMode = SubModeScatter
		g.DotLimit = dotLimits.Inky
		g.Pcm = data.PCM50
		g.Dir = geom.Up

	case Clyde:
		g.Mode = ModeHome
		g.SubMode = SubModeScatter
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
func (g *Actor) SetLeaveState() {
	g.Mode = ModeLeaving
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

	case SubModeChase:
		if subMode == SubModeScared || subMode == SubModeScatter {
			g.ReversePending = true
		}

	case SubModeScatter:
		if subMode == SubModeScared || subMode == SubModeChase {
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
// based on its current heading. Returns true if
// ghost has moved to a different tile.
func (g *Actor) Move() bool {
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
func (g *Actor) Draw(v *video.Video, isWhite bool, wobble bool) {
	var look sprite.Look
	var pal color.Palette
	if g.Visible {
		switch {
		case g.Dir.IsUp():
			look = sprite.GhostUp1
		case g.Dir.IsLeft():
			look = sprite.GhostDown2
		case g.Dir.IsDown():
			look = sprite.GhostDown1
		case g.Dir.IsRight():
			look = sprite.GhostRight1
		}
		pal = g.Pal
		if g.ScoreLook > 0 {
			look = g.ScoreLook
		} else {
			switch {
			case g.Mode == ModeReturning:
				pal = color.PalEyes
			case g.SubMode == SubModeScared:
				look = sprite.GhostScared1
				pal = color.PalScared
				if isWhite {
					pal = color.PalScaredFlash
				}
			}
			if wobble {
				look += 1
			}
		}
		v.AddSprite(g.Pos, look, pal)
	}
}
