package data

import "github.com/adrmcintyre/ebiman/video"

// Score related constants
const (
	// ExtraLifeScore defines how many points must be scored for an extra
	// life to be awarded. Note this is a one-time only award!
	ExtraLifeScore = 10000

	DotScore        = 10  // score for eating a dot
	DotScoreCharge1 = 20  // score for eating a dot with unit charge
	DotScoreCharge2 = 50  // score for eating a dot with two charges
	PowerScore      = 200 // score for eating a power pill
)

// A GhostScoreEntry describes how many points are awarded
// for consuming a ghost, and what to display.
type GhostScoreEntry struct {
	Score int          // points to award
	Look  video.Sprite // sprite to display
}

// GhostScore defines a GhostScoreEntry for each consecutive ghost
// consumed during the same period of panic.
var GhostScore = [4]GhostScoreEntry{
	{200, video.SpriteScore200},
	{400, video.SpriteScore400},
	{800, video.SpriteScore800},
	{1600, video.SpriteScore1600},
}
