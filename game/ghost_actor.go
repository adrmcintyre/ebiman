package game

import (
	"math/rand"

	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/geom"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
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

type GhostActor struct {
	// configuration fields, these don't change once set
	Id          GhostId
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
type GhostId int

const (
	BLINKY GhostId = iota
	PINKY
	INKY
	CLYDE
)

func MakeBlinky() GhostActor {
	return GhostActor{
		Id:                BLINKY,
		Pal:               color.PAL_BLINKY,
		HomePos:           geom.BLINKY_HOME,
		StartPos:          geom.BLINKY_START,
		ScatterPos:        geom.BLINKY_SCATTER,
		AllDotLimit:       0,
		DotsAtHomeCounter: 0,
	}
}

func MakePinky() GhostActor {
	return GhostActor{
		Id:                PINKY,
		Pal:               color.PAL_PINKY,
		HomePos:           geom.PINKY_HOME,
		StartPos:          geom.PINKY_HOME,
		ScatterPos:        geom.PINKY_SCATTER,
		AllDotLimit:       7,
		DotsAtHomeCounter: 0,
	}
}

func MakeInky() GhostActor {
	return GhostActor{
		Id:                INKY,
		Pal:               color.PAL_INKY,
		HomePos:           geom.INKY_HOME,
		StartPos:          geom.INKY_HOME,
		ScatterPos:        geom.INKY_SCATTER,
		AllDotLimit:       17,
		DotsAtHomeCounter: 0,
	}
}

func MakeClyde() GhostActor {
	return GhostActor{
		Id:                CLYDE,
		Pal:               color.PAL_CLYDE,
		HomePos:           geom.CLYDE_HOME,
		StartPos:          geom.CLYDE_HOME,
		ScatterPos:        geom.CLYDE_SCATTER,
		AllDotLimit:       32,
		DotsAtHomeCounter: 0,
	}
}

func MakeGhostActors() [4]GhostActor {
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

func (g *GhostActor) UpdateTarget(pm *PacmanActor, blinky *GhostActor) {
	switch g.Mode {
	case MODE_RETURNING:
		g.TargetPos = g.HomePos
	case MODE_PLAYING:
		switch g.SubMode {
		case SUBMODE_SCATTER:
			g.TargetPos = g.ScatterPos
		case SUBMODE_CHASE:
			g.TargetPos = g.GetChaseTarget(pm, blinky)
		}
	}
}

func (g *GhostActor) GetChaseTarget(pm *PacmanActor, blinky *GhostActor) geom.Position {
	switch g.Id {
	case PINKY:
		targetPos := pm.Pos.Add(pm.Dir.Scale(4 * 8))
		if pm.Dir.IsUp() {
			targetPos.X -= 4 * 8
		}
		return targetPos
	case INKY:
		return pm.Pos.Add(pm.Dir.Scale(4 * 8)).Add(pm.Pos.Sub(blinky.Pos))
	case CLYDE:
		if g.Pos.TileDistSq(pm.Pos) < 64 {
			return g.ScatterPos
		}
	}

	return pm.Pos
}

// TODO inject speeds on ghost construction?
func (g *GhostActor) Steer(v *video.Video, pacman *PacmanActor, blinky *GhostActor, speeds *data.Speeds, ghostAi bool) {
	switch g.Mode {
	case MODE_HOME:
		reachedTop := g.Dir.IsUp() && g.Pos.Y <= geom.HOME_TOP
		reachedBot := g.Dir.IsDown() && g.Pos.Y >= geom.HOME_BOTTOM
		if reachedTop || reachedBot {
			// bounce
			g.Dir = g.Dir.Reverse()
		}
		return

	case MODE_LEAVING:
		//         <--+
		//            |
		// ;---------G|G--------;
		// ;   x      |         ;
		// ;   -------+         ;
		// ;          +------x  ;
		// ;--------------------;
		if g.Pos.X < geom.HOME_CENTRE.X {
			g.Dir = geom.RIGHT
		} else if g.Pos.X > geom.HOME_CENTRE.X {
			g.Dir = geom.LEFT
		} else if g.Pos.Y == geom.HOME_EXITED_Y {
			g.Mode = MODE_PLAYING
			g.Dir = geom.LEFT
			if g.SubMode == SUBMODE_SCARED {
				g.Pcm = speeds.GhostBlue
			} else {
				g.Pcm = speeds.Ghost
			}
			// TODO apply submode rules???
		} else {
			g.Dir = geom.UP
		}
		return

	case MODE_RETURNING:
		if g.Pos == g.HomePos {
			g.Mode = MODE_HOME
			g.SetSubMode(SUBMODE_SCATTER)
			g.Pcm = data.PCM_40 // move at slowest speed when home (1 pixel every other frame)
			g.Dir = geom.UP
			return
		}
	}

	// TODO could add these as utility methods
	hCentred := g.Pos.X&7 == 0
	vCentred := g.Pos.Y&7 == 0

	if !(hCentred && vCentred) {
		// take care of reversals when transitioning between tiles
		hEntering := g.Pos.X&7 == 4
		vEntering := g.Pos.Y&7 == 4
		if (hEntering && vCentred) || (vEntering && hCentred) {
			if g.ReversePending {
				g.Dir = g.Dir.Reverse()
				g.ReversePending = false
			}
		}
		return
	}

	// decision time - we're at the centre of a tile
	g.UpdateTarget(pacman, blinky)

	exits := g.ComputeExits(v)
	g.Dir = g.ChooseExitDirection(exits, ghostAi)
}

type exitResult struct {
	Dir     geom.Delta
	NextPos geom.Position
}

func (g *GhostActor) ComputeExits(v *video.Video) []exitResult {
	// TODO: heap allocation - to avoid this the caller could supply
	// a reusable buffer to write to instead
	exits := make([]exitResult, 0, 3)

	// anti clockwise of current heading
	dir := g.Dir.TurnLeft()

	for range 3 {
		nextPos := g.Pos.Add(dir.Scale(8))
		nextTile := v.GetTile(nextPos.TileXY())

		viable := nextTile.IsTraversable()
		gateOpen := g.Mode == MODE_RETURNING
		onGate := nextTile == tile.GATE_LEFT || nextTile == tile.GATE_RIGHT
		onHome := nextTile == tile.HOME_LEFT || nextTile == tile.HOME_RIGHT

		if gateOpen && (onGate || onHome) {
			// open the gate for returning ghosts
			viable = true
		} else if g.SubMode != SUBMODE_SCARED {
			// cannot turn UP at one of 4 special tiles
			x, y := g.Pos.TileXY()
			specialTile := (x == 12 || x == 15) && (y == 12 || y == 24)
			if dir.IsUp() && specialTile {
				viable = false
			}
		}

		if viable {
			exits = append(exits, exitResult{
				Dir:     dir,
				NextPos: nextPos,
			})
		}

		// try one turn clockwise
		dir = dir.TurnRight()
	}

	return exits
}

func (g *GhostActor) ChooseExitDirection(exits []exitResult, ai bool) geom.Delta {
	n := len(exits)
	if n == 0 {
		return g.Dir
	}
	if n == 1 {
		return exits[0].Dir
	}

	if g.Mode == MODE_PLAYING && (g.SubMode == SUBMODE_SCARED || !ai) {
		return exits[rand.Intn(n)].Dir
	}

	bestExit := -1
	bestD2 := 32767
	for i := range n {
		if d2 := g.TargetPos.TileDistSq(exits[i].NextPos); d2 < bestD2 {
			// TODO - ties should be broken in order up,left,down
			bestD2 = d2
			bestExit = i
		}
	}

	return exits[bestExit].Dir
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
