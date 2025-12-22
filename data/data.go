package data

import "github.com/adrmcintyre/poweraid/sprite"

// Miscellaneous constants
const (
	FPS = 60 // frames-per-sec
	UPS = 30 // updates-per-sec

	// EXTRA_LIFE_SCORE defines how many points must be scored for an extra
	// life to be awarded. Note this is a one-time only award!
	EXTRA_LIFE_SCORE = 10000

	DOT_SCORE   = 10 // score for eating a dot
	POWER_SCORE = 50 // score for eating a power pill

	// Pacman pauses briefly when eating (but the ghosts continue moving).
	// These constants specifies for how long.
	DOT_STALL   = 1 // how long pacman stalls after eating a dot
	POWER_STALL = 4 // how long pacman stalls after eating a power pill

	// DISPLAY_GHOST_SCORE_MS defines how many milliseconds to display
	// its points value after a ghost is consumed.
	DISPLAY_GHOST_SCORE_MS = 1000

	// WHITE_BLUE_PERIOD defines how many updates between
	// ghosts flashing white and blue.
	WHITE_BLUE_PERIOD = 14
)

// A GhostScoreEntry describes how many points are awarded
// for consuming a ghost, and what to display.
type GhostScoreEntry struct {
	Score int         // points to award
	Look  sprite.Look // sprite to display
}

// GhostScore defines a GhostScoreEntry for each consecutive ghost
// consumed during the same period of panic.
var GhostScore = [4]GhostScoreEntry{
	{200, sprite.SCORE_200},
	{400, sprite.SCORE_400},
	{800, sprite.SCORE_800},
	{1600, sprite.SCORE_1600},
}
