package sprite

// A Look identifies a specific sprite bitmap.
type Look byte

// The identifiers for each sprite.
const (
	// The bonus sprites
	CHERRY     Look = 0x00
	STRAWBERRY Look = 0x01
	ORANGE     Look = 0x02
	BELL       Look = 0x03
	APPLE      Look = 0x04
	PINEAPPLE  Look = 0x05
	GALAXIAN   Look = 0x06
	KEY        Look = 0x07

	// animations of pacman turning into a puddle
	GHOST_STICK1  Look = 0x08
	GHOST_STICK2  Look = 0x09
	GHOST_WORM1   Look = 0x0a
	GHOST_WORM2   Look = 0x0b
	GHOST_PUDDLE1 Look = 0x0c
	GHOST_PUDDLE2 Look = 0x0d

	// pacman looking up, mouth moving
	PACMAN_UP1 Look = 0x0e
	PACMAN_UP2 Look = 0x0f

	// scared ghost, wobbling
	GHOST_SCARED1 Look = 0x1c
	GHOST_SCARED2 Look = 0x1d

	// pacman looking left, mouth moving
	PACMAN_LEFT1 Look = 0x1e
	PACMAN_LEFT2 Look = 0x1f

	// normal ghost looking right/down/left/up, wobbling
	GHOST_RIGHT1 Look = 0x20
	GHOST_RIGHT2 Look = 0x21
	GHOST_DOWN1  Look = 0x22
	GHOST_DOWN2  Look = 0x23
	GHOST_LEFT1  Look = 0x24
	GHOST_LEFT2  Look = 0x25
	GHOST_UP1    Look = 0x26
	GHOST_UP2    Look = 0x27

	// scores for each consecutive ghost eaten
	SCORE_200  Look = 0x28
	SCORE_400  Look = 0x29
	SCORE_800  Look = 0x2a
	SCORE_1600 Look = 0x2b

	// more of pacman's movement animations
	PACMAN_RIGHT1 Look = 0x2c
	PACMAN_DOWN1  Look = 0x2d
	PACMAN_RIGHT2 Look = 0x2e
	PACMAN_DOWN2  Look = 0x2f

	// pacman with mouth shut - good for any direction
	PACMAN_SHUT Look = 0x30

	// pacman going "pop" on death
	PACMAN_EXPLODE Look = 0x31

	// frames 1-12 of pacman's demise (pop actually comes last)
	PACMAN_DEAD1  Look = 0x34
	PACMAN_DEAD2  Look = 0x35
	PACMAN_DEAD3  Look = 0x36
	PACMAN_DEAD4  Look = 0x37
	PACMAN_DEAD5  Look = 0x38
	PACMAN_DEAD6  Look = 0x39
	PACMAN_DEAD7  Look = 0x3a
	PACMAN_DEAD8  Look = 0x3b
	PACMAN_DEAD9  Look = 0x3c
	PACMAN_DEAD10 Look = 0x3d
	PACMAN_DEAD11 Look = 0x3e
	PACMAN_DEAD12 Look = 0x3f
)
