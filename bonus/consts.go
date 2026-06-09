package bonus

import (
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/video"
)

// when to drop bonuses
const (
	FirstBonusDots  = 70
	SecondBonusDots = 170
)

// How long to keep bonus visible - when a bonus drops, it is made
// visible for a random time between MinBonusTime and MaxBonusTime.
const (
	MinBonusTime = 9 * data.FPS
	MaxBonusTime = 10 * data.FPS
)

// An Id identifies a particular bonus item.
type Id int

// Bonus items
const (
	Cherry Id = iota
	Strawberry
	Orange
	Bell
	Apple
	Pineapple
	Galaxian
	Key
	bonusCount
)

// LevelBonus contains the bonus appearing at each level.
var LevelBonus = [21]Id{
	Cherry,     // level 1
	Strawberry, // level 2
	Orange,     // level 3
	Orange,     // level 4
	Apple,      // level 5
	Apple,      // level 6
	Pineapple,  // level 7
	Pineapple,  // level 8
	Galaxian,   // level 9
	Galaxian,   // level 10
	Bell,       // level 11
	Bell,       // level 12
	Key,        // level 13
	Key,        // level 14
	Key,        // level 15
	Key,        // level 16
	Key,        // level 17
	Key,        // level 18
	Key,        // level 19
	Key,        // level 20
	Key,        // level 21+
}

// An InfoEntry describes the appearance and value of a bonus.
type InfoEntry struct {
	// Look is the sprite to use when displaying a bonus drop.
	Look video.Sprite
	// BaseTile is the first of 4 consecutive tiles to use when displaying
	// a consumed bonus in the status area.
	BaseTile video.Tile
	// Pal is the palette to apply for the bonus sprite,
	// as well as the status area tiles.
	Pal video.Palette
	// Score is how much the bonus is worth.
	Score int
	// Tiles is a sequence of tiles to display the score briefly when
	// the bonus is consumed.
	Tiles []video.Tile
}

// Info provides an InfoEntry for each bonus item (one for each Id).
var Info = [bonusCount]InfoEntry{
	{
		video.SpriteCherry, video.TileCherryBase, video.PalCherry, 100,
		[]video.Tile{video.TileSpace, video.TileScore100, video.TileScoreX00, video.TileSpace},
	},
	{
		video.SpriteStrawberry, video.TileStrawberryBase, video.PalStrawberry, 300,
		[]video.Tile{video.TileSpace, video.TileScore300, video.TileScoreX00, video.TileSpace},
	},
	{
		video.SpriteOrange, video.TileOrangeBase, video.PalOrange, 500,
		[]video.Tile{video.TileSpace, video.TileScore500, video.TileScoreX00, video.TileSpace},
	},
	{
		video.SpriteBell, video.TileBellBase, video.PalBell, 700,
		[]video.Tile{video.TileSpace, video.TileScore700, video.TileScoreX00, video.TileSpace},
	},
	{
		video.SpriteApple, video.TileAppleBase, video.PalApple, 1000,
		[]video.Tile{video.TileSpace, video.TileScore1000, video.TileScoreX000_1, video.TileScoreX000_2},
	},
	{
		video.SpritePineapple, video.TilePineappleBase, video.PalPineapple, 2000,
		[]video.Tile{video.TileScore2000_1, video.TileScore2000_2, video.TileScoreX000_1, video.TileScoreX000_2},
	},
	{
		video.SpriteGalaxian, video.TileGalaxianBase, video.PalGalaxian, 3000,
		[]video.Tile{video.TileScore3000_1, video.TileScore3000_2, video.TileScoreX000_1, video.TileScoreX000_2},
	},
	{
		video.SpriteKey, video.TileKeyBase, video.PalKey, 5000,
		[]video.Tile{video.TileScore5000_1, video.TileScore5000_2, video.TileScoreX000_1, video.TileScoreX000_2},
	},
}
