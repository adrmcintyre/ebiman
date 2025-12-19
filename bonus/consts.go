package bonus

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
)

// when to drop bonuses
const FIRST_BONUS_DOTS = 70
const SECOND_BONUS_DOTS = 170

// how long to keep bonus visible
const MIN_BONUS_TIME = 9 * data.FPS
const MAX_BONUS_TIME = 10 * data.FPS

type Id int

// bonus types
const (
	CHERRY Id = iota
	STRAWBERRY
	ORANGE
	BELL
	APPLE
	PINEAPPLE
	GALAXIAN
	KEY
	bonusCount
)

// These are the bonuses appearing in each level
var BonusType = [21]Id{
	CHERRY,     // level 1
	STRAWBERRY, // level 2
	ORANGE,     // level 3
	ORANGE,     // level 4
	APPLE,      // level 5
	APPLE,      // level 6
	PINEAPPLE,  // level 7
	PINEAPPLE,  // level 8
	GALAXIAN,   // level 9
	GALAXIAN,   // level 10
	BELL,       // level 11
	BELL,       // level 12
	KEY,        // level 13
	KEY,        // level 14
	KEY,        // level 15
	KEY,        // level 16
	KEY,        // level 17
	KEY,        // level 18
	KEY,        // level 19
	KEY,        // level 20
	KEY,        // level 21+
}

type InfoEntry struct {
	Look     sprite.Look
	BaseTile tile.Tile
	Pal      color.Palette
	Score    int
	Tiles    []tile.Tile
}

var Info = [bonusCount]InfoEntry{
	{
		sprite.CHERRY, tile.CHERRY_BASE, color.PAL_CHERRY, 100,
		[]tile.Tile{tile.SPACE, tile.SCORE_100, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.STRAWBERRY, tile.STRAWBERRY_BASE, color.PAL_STRAWBERRY, 300,
		[]tile.Tile{tile.SPACE, tile.SCORE_300, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.ORANGE, tile.ORANGE_BASE, color.PAL_ORANGE, 500,
		[]tile.Tile{tile.SPACE, tile.SCORE_500, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.BELL, tile.BELL_BASE, color.PAL_BELL, 700,
		[]tile.Tile{tile.SPACE, tile.SCORE_700, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.APPLE, tile.APPLE_BASE, color.PAL_APPLE, 1000,
		[]tile.Tile{tile.SPACE, tile.SCORE_1000, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
	{
		sprite.PINEAPPLE, tile.PINEAPPLE_BASE, color.PAL_PINEAPPLE, 2000,
		[]tile.Tile{tile.SCORE_2000_1, tile.SCORE_2000_2, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
	{
		sprite.GALAXIAN, tile.GALAXIAN_BASE, color.PAL_GALAXIAN, 3000,
		[]tile.Tile{tile.SCORE_3000_1, tile.SCORE_3000_2, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
	{
		sprite.KEY, tile.KEY_BASE, color.PAL_KEY, 5000,
		[]tile.Tile{tile.SCORE_5000_1, tile.SCORE_5000_2, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
}
