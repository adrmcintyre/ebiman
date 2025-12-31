package bonus

import (
	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/sprite"
	"github.com/adrmcintyre/ebiman/tile"
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
	Look sprite.Look
	// BaseTile is the first of 4 consecutive tiles to use when displaying
	// a consumed bonus in the status area.
	BaseTile tile.Tile
	// Pal is the palette to apply for the bonus sprite,
	// as well as the status area tiles.
	Pal color.Palette
	// Score is how much the bonus is worth.
	Score int
	// Tiles is a sequence of tiles to display the score briefly when
	// the bonus is consumed.
	Tiles []tile.Tile
}

// Info provides an InfoEntry for each bonus item (one for each Id).
var Info = [bonusCount]InfoEntry{
	{
		sprite.Cherry, tile.CherryBase, color.PalCherry, 100,
		[]tile.Tile{tile.Space, tile.Score100, tile.ScoreX00, tile.Space},
	},
	{
		sprite.Strawberry, tile.StrawberryBase, color.PalStrawberry, 300,
		[]tile.Tile{tile.Space, tile.Score300, tile.ScoreX00, tile.Space},
	},
	{
		sprite.Orange, tile.OrangeBase, color.PalOrange, 500,
		[]tile.Tile{tile.Space, tile.Score500, tile.ScoreX00, tile.Space},
	},
	{
		sprite.Bell, tile.BellBase, color.PalBell, 700,
		[]tile.Tile{tile.Space, tile.Score700, tile.ScoreX00, tile.Space},
	},
	{
		sprite.Apple, tile.AppleBase, color.PalApple, 1000,
		[]tile.Tile{tile.Space, tile.Score1000, tile.ScoreX000_1, tile.ScoreX000_2},
	},
	{
		sprite.Pineapple, tile.PineappleBase, color.PalPineapple, 2000,
		[]tile.Tile{tile.Score2000_1, tile.Score2000_2, tile.ScoreX000_1, tile.ScoreX000_2},
	},
	{
		sprite.Galaxian, tile.GalaxianBase, color.PalGalaxian, 3000,
		[]tile.Tile{tile.Score3000_1, tile.Score3000_2, tile.ScoreX000_1, tile.ScoreX000_2},
	},
	{
		sprite.Key, tile.KeyBase, color.PalKey, 5000,
		[]tile.Tile{tile.Score5000_1, tile.Score5000_2, tile.ScoreX000_1, tile.ScoreX000_2},
	},
}
