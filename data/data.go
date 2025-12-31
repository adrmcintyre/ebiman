package data

import "github.com/adrmcintyre/ebiman/sprite"

// Miscellaneous constants
const (
	FPS = 60 // frames-per-sec
	UPS = 30 // updates-per-sec

	// ExtraLifeScore defines how many points must be scored for an extra
	// life to be awarded. Note this is a one-time only award!
	ExtraLifeScore = 10000

	DotScore        = 10  // score for eating a dot
	DotScoreCharge1 = 20  // score for eating a dot with unit charge
	DotScoreCharge2 = 50  // score for eating a dot with two charges
	PowerScore      = 200 // score for eating a power pill

	// Pacman pauses briefly when eating (but the ghosts continue moving).
	// These constants specifies for how long.
	DotStall   = 1 // how long pacman stalls after eating a dot
	PowerStall = 4 // how long pacman stalls after eating a power pill

	// DisplayGhostScoreMs defines how many milliseconds to display
	// its points value after a ghost is consumed.
	DisplayGhostScoreMs = 1000

	// WhiteBluePeriod defines how many updates between
	// ghosts flashing white and blue.
	WhiteBluePeriod = 14
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
	{200, sprite.Score200},
	{400, sprite.Score400},
	{800, sprite.Score800},
	{1600, sprite.Score1600},
}
