package main

import "github.com/adrmcintyre/poweraid/data"

// Current level 'constants'
type LevelConfig struct {
	BlueTime       int                // how many game updates ghosts remain blue
	WhiteBlueCount int                // number of white-blue flashes before ghosts revert
	IdleLimit      int                // max number of frames without eating before pacman is considered idle
	DotLimits      data.DotLimitEntry // dot limits for inky, pinky and clyde
	Speeds         data.Speeds        // various PCM trains for pacman and ghosts
	SwitchTactics  [7]int             // frames counts (as offsets) for ghosts to switch between scatter and chase
	ElroyPills1    int                // blinky's first speed boost when this number of pills left
	ElroyPills2    int                // blinky's second speed boost
	BonusType      int
	BonusInfo      data.BonusInfoEntry
}

func DefaultLevelConfig() LevelConfig {
	return LevelConfig{}
}
