package data

// Timing constants
const (
	FPS = 60 // frames-per-sec
	UPS = 30 // updates-per-sec

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
