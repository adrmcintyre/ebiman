package data

import (
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
	MinBonusTime = 9 * FPS
	MaxBonusTime = 10 * FPS
)

// An BonusId identifies a particular bonus item.
type BonusId int

// Bonus items
const (
	BonusCherry BonusId = iota
	BonusStrawberry
	BonusOrange
	BonusBell
	BonusApple
	BonusPineapple
	BonusGalaxian
	BonusKey
	bonusCount
)

// LevelBonus contains the bonus appearing at each level.
var LevelBonus = [21]BonusId{
	BonusCherry,     // level 1
	BonusStrawberry, // level 2
	BonusOrange,     // level 3
	BonusOrange,     // level 4
	BonusApple,      // level 5
	BonusApple,      // level 6
	BonusPineapple,  // level 7
	BonusPineapple,  // level 8
	BonusGalaxian,   // level 9
	BonusGalaxian,   // level 10
	BonusBell,       // level 11
	BonusBell,       // level 12
	BonusKey,        // level 13
	BonusKey,        // level 14
	BonusKey,        // level 15
	BonusKey,        // level 16
	BonusKey,        // level 17
	BonusKey,        // level 18
	BonusKey,        // level 19
	BonusKey,        // level 20
	BonusKey,        // level 21+
}

// An BonusInfoEntry describes the appearance and value of a bonus.
type BonusInfoEntry struct {
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

// BonusInfo provides an InfoEntry for each bonus item (one for each Id).
var BonusInfo = [bonusCount]BonusInfoEntry{
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
