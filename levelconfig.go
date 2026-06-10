package main

import (
	"github.com/adrmcintyre/ebiman/data"
)

// LevelConfig describes the static configuration of an in-play level.
type LevelConfig struct {
	ScaredTime     int                 // how many game updates ghosts remain scared
	WhiteBlueCount int                 // number of white-blue flashes before ghosts revert
	IdleLimit      int                 // max number of frames without eating before pacman is considered idle
	DotLimits      data.DotLimitEntry  // dot limits for inky, pinky and clyde
	Speeds         data.Speeds         // various PCM trains for pacman and ghosts
	SwitchTactics  [7]int              // frames counts (as offsets) for ghosts to switch between scatter and chase
	ElroyPills1    int                 // blinky's first speed boost when this number of pills left
	ElroyPills2    int                 // blinky's second speed boost
	BonusType      data.BonusId        // identifies the bonus for the level
	BonusInfo      data.BonusInfoEntry // description of the bonus's attributes
	Electric       data.ElectricEntry  // how each ghost manipulates charge
}

// NewLevelConfig returns a new configuration based on the given level number and difficulty level.
func NewLevelConfig(levelNumber int, difficulty int) *LevelConfig {
	levelIndex := min(levelNumber, len(data.Level)-1)
	level := data.Level[levelIndex]

	speeds := data.SpeedData[level.SpeedIndex-3]
	elroy := data.Elroy[level.ElroyIndex]
	blueControl := data.BlueControl[level.BlueIndex]
	bonusType := data.LevelBonus[levelIndex]

	cfg := LevelConfig{
		SwitchTactics:  speeds.SwitchTactics,
		DotLimits:      data.DotLimit[level.DotLimitIndex],
		ElroyPills1:    elroy.Pills1,
		ElroyPills2:    elroy.Pills2,
		ScaredTime:     blueControl.BlueTime,
		WhiteBlueCount: blueControl.WhiteBlueCount,
		IdleLimit:      data.IdleLimit[level.IdleIndex],
		BonusType:      bonusType,
		BonusInfo:      data.BonusInfo[bonusType],
	}

	switch difficulty {
	case DifficultyEasy:
		cfg.ScaredTime *= 2
		cfg.Speeds = speeds.Easy
		cfg.Electric = data.Electric.Easy
	case DifficultyMedium:
		cfg.Speeds = speeds.Medium
		cfg.Electric = data.Electric.Medium
	case DifficultyHard:
		cfg.Speeds = speeds.Hard
		cfg.Electric = data.Electric.Hard
	}

	return &cfg
}
