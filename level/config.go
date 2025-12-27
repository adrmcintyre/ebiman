package level

import (
	"github.com/adrmcintyre/ebiman/bonus"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/option"
)

// Config describes the static configuration of an in-play level.
type Config struct {
	BlueTime       int                // how many game updates ghosts remain blue
	WhiteBlueCount int                // number of white-blue flashes before ghosts revert
	IdleLimit      int                // max number of frames without eating before pacman is considered idle
	DotLimits      data.DotLimitEntry // dot limits for inky, pinky and clyde
	Speeds         data.Speeds        // various PCM trains for pacman and ghosts
	SwitchTactics  [7]int             // frames counts (as offsets) for ghosts to switch between scatter and chase
	ElroyPills1    int                // blinky's first speed boost when this number of pills left
	ElroyPills2    int                // blinky's second speed boost
	BonusType      bonus.Id           // identifies the bonus for the level
	BonusInfo      bonus.InfoEntry    // description of the bonus's attributes
}

// DefaultConfig returns an uninitalised configuration.
func DefaultConfig() Config {
	return Config{}
}

// Init initialises the configuration based on the given level number and difficulty level.
func (cfg *Config) Init(levelNumber int, difficulty int) {
	levelIndex := min(levelNumber, len(data.Level)-1)
	level := data.Level[levelIndex]

	speeds := data.SpeedData[level.SpeedIndex-3]

	switch difficulty {
	case option.DIFFICULTY_EASY:
		cfg.Speeds = speeds.Easy
	case option.DIFFICULTY_MEDIUM:
		cfg.Speeds = speeds.Medium
	case option.DIFFICULTY_HARD:
		cfg.Speeds = speeds.Hard
	}

	cfg.SwitchTactics = speeds.SwitchTactics
	cfg.DotLimits = data.DotLimit[level.DotLimitIndex]
	elroy := data.Elroy[level.ElroyIndex]
	cfg.ElroyPills1 = elroy.Pills1
	cfg.ElroyPills2 = elroy.Pills2

	blueControl := data.BlueControl[level.BlueIndex]
	cfg.BlueTime = blueControl.BlueTime
	cfg.WhiteBlueCount = blueControl.WhiteBlueCount

	if difficulty == option.DIFFICULTY_EASY {
		cfg.BlueTime *= 2
	}

	cfg.IdleLimit = data.IdleLimit[level.IdleIndex]
	cfg.BonusType = bonus.LevelBonus[levelIndex]
	cfg.BonusInfo = bonus.Info[cfg.BonusType]
}
