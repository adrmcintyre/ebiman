package data

import (
	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/sprite"
	"github.com/adrmcintyre/poweraid/tile"
)

// when to drop bonuses
const FIRST_BONUS_DOTS = 70
const SECOND_BONUS_DOTS = 170

// how long to keep bonus visible
const MIN_BONUS_TIME = 9 * FPS
const MAX_BONUS_TIME = 10 * FPS

// bonus types
const (
	BONUS_CHERRY = iota
	BONUS_STRAWBERRY
	BONUS_ORANGE
	BONUS_BELL
	BONUS_APPLE
	BONUS_PINEAPPLE
	BONUS_GALAXIAN
	BONUS_KEY
	bonusCount
)

// These are the bonuses appearing in each level
var BonusType = [21]int{
	BONUS_CHERRY,     // level 1
	BONUS_STRAWBERRY, // level 2
	BONUS_ORANGE,     // level 3
	BONUS_ORANGE,     // level 4
	BONUS_APPLE,      // level 5
	BONUS_APPLE,      // level 6
	BONUS_PINEAPPLE,  // level 7
	BONUS_PINEAPPLE,  // level 8
	BONUS_GALAXIAN,   // level 9
	BONUS_GALAXIAN,   // level 10
	BONUS_BELL,       // level 11
	BONUS_BELL,       // level 12
	BONUS_KEY,        // level 13
	BONUS_KEY,        // level 14
	BONUS_KEY,        // level 15
	BONUS_KEY,        // level 16
	BONUS_KEY,        // level 17
	BONUS_KEY,        // level 18
	BONUS_KEY,        // level 19
	BONUS_KEY,        // level 20
	BONUS_KEY,        // level 21+
}

type BonusTiles []tile.Tile

type BonusInfoEntry struct {
	Look     sprite.Look
	BaseTile tile.Tile
	Pal      color.Palette
	Score    int
	Tiles    BonusTiles
}

var BonusInfo = [bonusCount]BonusInfoEntry{
	{
		sprite.CHERRY, tile.CHERRY_BASE, color.PAL_CHERRY, 100,
		BonusTiles{tile.SPACE, tile.SCORE_100, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.STRAWBERRY, tile.STRAWBERRY_BASE, color.PAL_STRAWBERRY, 300,
		BonusTiles{tile.SPACE, tile.SCORE_300, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.ORANGE, tile.ORANGE_BASE, color.PAL_ORANGE, 500,
		BonusTiles{tile.SPACE, tile.SCORE_500, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.BELL, tile.BELL_BASE, color.PAL_BELL, 700,
		BonusTiles{tile.SPACE, tile.SCORE_700, tile.SCORE_X00, tile.SPACE},
	},
	{
		sprite.APPLE, tile.APPLE_BASE, color.PAL_APPLE, 1000,
		BonusTiles{tile.SPACE, tile.SCORE_1000, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
	{
		sprite.PINEAPPLE, tile.PINEAPPLE_BASE, color.PAL_PINEAPPLE, 2000,
		BonusTiles{tile.SCORE_2000_1, tile.SCORE_2000_2, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
	{
		sprite.GALAXIAN, tile.GALAXIAN_BASE, color.PAL_GALAXIAN, 3000,
		BonusTiles{tile.SCORE_3000_1, tile.SCORE_3000_2, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
	{
		sprite.KEY, tile.KEY_BASE, color.PAL_KEY, 5000,
		BonusTiles{tile.SCORE_5000_1, tile.SCORE_5000_2, tile.SCORE_X000_1, tile.SCORE_X000_2},
	},
}
