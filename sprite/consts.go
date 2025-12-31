package sprite

// A Look identifies a specific sprite bitmap.
type Look byte

// The identifiers for each sprite.
const (
	// The bonus sprites
	Cherry     Look = 0x00
	Strawberry Look = 0x01
	Orange     Look = 0x02
	Bell       Look = 0x03
	Apple      Look = 0x04
	Pineapple  Look = 0x05
	Galaxian   Look = 0x06
	Key        Look = 0x07

	// animations of pacman turning into a puddle
	GhostStick1  Look = 0x08
	GhostStick2  Look = 0x09
	GhostWorm1   Look = 0x0a
	GhostWorm2   Look = 0x0b
	GhostPuddle1 Look = 0x0c
	GhostPuddle2 Look = 0x0d

	// pacman looking up, mouth moving
	PacmanUp1 Look = 0x0e
	PacmanUp2 Look = 0x0f

	// scared ghost, wobbling
	GhostScared1 Look = 0x1c
	GhostScared2 Look = 0x1d

	// pacman looking left, mouth moving
	PacmanLeft1 Look = 0x1e
	PacmanLeft2 Look = 0x1f

	// normal ghost looking right/down/left/up, wobbling
	GhostRight1 Look = 0x20
	GhostRight2 Look = 0x21
	GhostDown1  Look = 0x22
	GhostDown2  Look = 0x23
	GhostLeft1  Look = 0x24
	GhostLeft2  Look = 0x25
	GhostUp1    Look = 0x26
	GhostUp2    Look = 0x27

	// scores for each consecutive ghost eaten
	Score200  Look = 0x28
	Score400  Look = 0x29
	Score800  Look = 0x2a
	Score1600 Look = 0x2b

	// more of pacman's movement animations
	PacmanRight1 Look = 0x2c
	PacmanDown1  Look = 0x2d
	PacmanRight2 Look = 0x2e
	PacmanDown2  Look = 0x2f

	// pacman with mouth shut - good for any direction
	PacmanShut Look = 0x30

	// pacman going "pop" on death
	PacmanExplode Look = 0x31

	// frames 1-12 of pacman's demise (pop actually comes last)
	PacmanDead1  Look = 0x34
	PacmanDead2  Look = 0x35
	PacmanDead3  Look = 0x36
	PacmanDead4  Look = 0x37
	PacmanDead5  Look = 0x38
	PacmanDead6  Look = 0x39
	PacmanDead7  Look = 0x3a
	PacmanDead8  Look = 0x3b
	PacmanDead9  Look = 0x3c
	PacmanDead10 Look = 0x3d
	PacmanDead11 Look = 0x3e
	PacmanDead12 Look = 0x3f
)
